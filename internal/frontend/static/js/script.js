// Add event listener for setting base URL.
document.addEventListener('DOMContentLoaded', setBaseUrl);
document.addEventListener('htmx:load', setBaseUrl);

// setBaseUrl sets the base URL for the secret form.
function setBaseUrl() {
  const port = window.location.port;
  let baseUrl = window.location.protocol + '//' + window.location.hostname;
  if (port && port !== '80' && port !== '443') {
    baseUrl += ':' + port;
  }
  
  
  secretFormBaseUrl = document.getElementById('secret-form-base-url');
  if (secretFormBaseUrl) {
    secretFormBaseUrl.value = baseUrl;
  }
}

// Mask the secret passphrase on the secret result page.
document.addEventListener('htmx:afterSwap', () => {
  const maskedLength = 40;
  const maskedValue = '\u2022'.repeat(maskedLength);

  const element = document.getElementById('secret-passphrase');
  if (!element) {
    return;
  }

  element.value = maskedValue;
});

// copyToClipboard copies the contents of an element to the clipboard.
function copyToClipboard(elementId, feedbackElementId) {
  const element = document.getElementById(elementId);
  if (!element) {
    return;
  }

  let text = element.innerText || element.textContent || element.value;
  // To handle the masked passphrase we need to check if the custom attribute is set.
  // This should override the text value.
  if (element.getAttribute('data-value')) {
    text = element.getAttribute('data-value');
  }

  navigator.clipboard.writeText(text).then(() => {
    if (feedbackElementId) {
      const feedback = document.getElementById(feedbackElementId);
      feedback.disabled = true;

      feedback.classList.remove('hover:text-gray-200');
      const innerHTML = feedback.innerHTML;

      feedback.innerHTML = `
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6 text-green-600">
          <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
        </svg>
      `;

      setTimeout(() => {
        feedback.innerHTML = innerHTML;
        feedback.disabled = false;
        feedback.classList.add('hover:text-gray-200');
      }, 3000);
    }
  });
}

// disableElement disables an element by ID.
function disableElement(elementId) {
  const element = document.getElementById(elementId);
  if (!element) {
    return;
  }
  element.disabled = true;
}

// enableElement enables an element by ID.
function enableElement(elementId) {
  const element = document.getElementById(elementId);
  if (!element) {
    return;
  }
  element.disabled = false;
}
