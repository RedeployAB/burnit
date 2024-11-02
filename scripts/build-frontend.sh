#!/bin/bash
expected_esbuild_sha256=$ESBUILD_SHA256
expected_htmx_sha256=$HTMX_SHA256
expected_tailwindcss_sha256=$TAILWINDCSS_SHA256
curr_dir=$(pwd)

# Download esbuild if it doesn't exist and verify the SHA256 checksum.
if [ -z "$(command -v esbuild)" ]; then
    curl -fsSL https://esbuild.github.io/dl/v0.24.0 | sh
    if [ $? -ne 0 ]; then
        echo "Failed to download esbuild"
        exit 1
    fi
    if [ "$expected_esbuild_sha256" != "$(sha256sum esbuild)" ]; then
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

  if [ "$expected_tailwindcss_sha256" != "$(sha256sum tailwindcss-linux-x64)" ]; then
    echo "SHA256 checksum mismatch for tailwindcss"
    exit 1
  fi
fi


# Download htmx.min.js.
curl -so internal/frontend/static/js/htmx.min.js https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js

# Verify the SHA256 checksum.
actual_htmx_sha256=$(sha256sum internal/frontend/static/js/htmx.min.js | awk '{print $1}')
if [ "$expected_htmx_sha256" != "$actual_htmx_sha256" ]; then
  echo "SHA256 checksum mismatch for htmx.min.js"
  exit 1
fi

# Build CSS with tailwindcss.
cd internal/frontend
tailwindcss -i ./static/css/tailwind.css -o ./static/css/main.css
cd $curr_dir

# Bundle JS and CSS.
esbuild internal/frontend/static/js/main.js --bundle --minify --outfile=internal/frontend/static/js/main.min.js
esbuild internal/frontend/static/css/main.css --bundle --minify --outfile=internal/frontend/static/css/main.min.css
