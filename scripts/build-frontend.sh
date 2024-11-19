#!/bin/bash
expected_esbuild_sha256=$ESBUILD_SHA256
expected_htmx_sha256=$HTMX_SHA256
expected_tailwindcss_sha256=$TAILWINDCSS_SHA256
curr_dir=$(pwd)

if [[ "$OSTYPE" == "darwin"* ]]; then
  sed_flags=(-i '')
else
  sed_flags=(-i)
fi

# Download esbuild if it doesn't exist and verify the SHA256 checksum.
if [ -z "$(command -v esbuild)" ]; then
    curl -fsSL https://esbuild.github.io/dl/v0.24.0 | sh
    if [ $? -ne 0 ]; then
        echo "Failed to download esbuild"
        exit 1
    fi

    $actual_esbuild_sha256=$(sha256sum esbuild | awk '{print $1}')
    if [ "$expected_esbuild_sha256" != "actual_esbuild_sha256" ]; then
        echo "SHA256 checksum mismatch for esbuild"
        exit 1
    fi
fi

if [ -z "$(command -v tailwindcss)" ]; then
  curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.14/tailwindcss-linux-x64
  if [ $? -ne 0 ]; then
    echo "Failed to download tailwindcss"
    exit 1
  fi
  
  chmod +x tailwindcss-linux-x64
  mv tailwindcss-linux-x64 tailwindcss

  $actual_tailwindcss_sha256=$(sha256sum tailwindcss | awk '{print $1}')
  if [ "$expected_tailwindcss_sha256" != "$actual_tailwindcss_sha256" ]; then
    echo "SHA256 checksum mismatch for tailwindcss"
    exit 1
  fi
fi


# Download htmx.esm.js.
curl -sLo internal/frontend/static/js/htmx.esm.js https://github.com/bigskysoftware/htmx/releases/download/v2.0.3/htmx.esm.js
sed "${sed_flags[@]}" 's/return eval(str)/return (0, eval)(str)/g' internal/frontend/static/js/htmx.esm.js

# Verify the SHA256 checksum.
actual_htmx_sha256=$(sha256sum internal/frontend/static/js/htmx.esm.js | awk '{print $1}')
if [ "$expected_htmx_sha256" != "$actual_htmx_sha256" ]; then
  echo "SHA256 checksum mismatch for htmx.esm.js"
  exit 1
fi

# Build CSS with tailwindcss.
cd internal/frontend
tailwindcss -i ./static/css/tailwind.css -o ./static/css/main.css
cd $curr_dir

# Bundle JS and CSS.
esbuild internal/frontend/static/js/main.js --bundle --minify --outfile=internal/frontend/static/js/main.min.js
esbuild internal/frontend/static/css/main.css --bundle --minify --outfile=internal/frontend/static/css/main.min.css

gzip -k -f internal/frontend/static/js/main.min.js
gzip -k -f internal/frontend/static/css/main.min.css
gzip -k -f internal/frontend/static/icons/*.png

sed "${sed_flags[@]}" 's/main\.css/main.min.css/g' internal/frontend/templates/base.html
sed "${sed_flags[@]}" '/unpkg.com/d' internal/frontend/templates/base.html
sed "${sed_flags[@]}" 's/script\.js/main.min.js/g' internal/frontend/templates/base.html
