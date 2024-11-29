const maxSecretValueLength = 3500;

// Add event listener for setting base URL.
document.addEventListener('DOMContentLoaded', setBaseUrl);
document.addEventListener('htmx:load', setBaseUrl);

// Handle events after htmx swap for secret form.
document.addEventListener('htmx:afterSwap', (event) => {
  const maskedLength = 40;
  const maskedValue = '\u2022'.repeat(maskedLength);

  const target = event.target;
  if (target.id == 'secret-form-container') {
    const secretForm = document.getElementById('secret-form');
    secretForm.reset();
    disableElement('secret-form-fields');

    const passphraseField = document.getElementById('secret-passphrase');
    if (passphraseField) {
      passphraseField.value = maskedValue;
    }

    const copySecretFullLink = document.getElementById('copy-secret-full-link');
    if (copySecretFullLink) {
      copySecretFullLink.addEventListener('click', () => {
        copyToClipboard('secret-full-link', 'copy-secret-full-link');
      });
    }

    const copySecretPartialLink = document.getElementById('copy-secret-partial-link');
    if (copySecretPartialLink) {
      copySecretPartialLink.addEventListener('click', () => {
        copyToClipboard('secret-partial-link', 'copy-secret-partial-link');
      });
    }

    const copySecretPassphrase = document.getElementById('copy-secret-passphrase');
    if (copySecretPassphrase) {
      copySecretPassphrase.addEventListener('click', () => {
        copyToClipboard('secret-passphrase', 'copy-secret-passphrase');
      });
    }

    const secretFormTextAreaCounter = document.getElementById('secret-form-textarea-counter');
    if (secretFormTextAreaCounter) {
      secretFormTextAreaCounter.textContent = '0/' + maxSecretValueLength;
    }

    const errorOverlayCloseButton = document.getElementById('error-overlay-close-button');
    if (errorOverlayCloseButton) {
      errorOverlayCloseButton.addEventListener('click', () => {
        const overlay = document.getElementById('error-overlay');
        overlay.remove();
        enableElement('secret-form-fields');
      });
    }

    if (event.detail.requestConfig.verb == 'get') {
      enableElement('secret-form-fields');
    }
  }

  if (target.id == 'secret-result-container') {
    const errorOverlayCloseButton = document.getElementById('error-overlay-close-button');
    if (errorOverlayCloseButton) {
      errorOverlayCloseButton.addEventListener('click', () => {
        const overlay = document.getElementById('error-overlay');
        overlay.remove();
      });
    }
  }
});

// Handle events for secret result.
document.addEventListener('DOMContentLoaded', () => {
  const copySecretResultValue = document.getElementById('copy-secret-result-value');
  if (copySecretResultValue) {
    copySecretResultValue.addEventListener('click', () => {
      copyToClipboard('secret-result-value', 'copy-secret-result-value');
    });
  }
});

// Handle events before htmx swap for secret result. This to be able to update the swap
// type to beforeend in case of an error for the overlay.
document.addEventListener('htmx:beforeSwap', (event) => {
  const detail = event.detail;

  if (detail.xhr.status == 400) {
    detail.shouldSwap = true;
  }

  if (detail.xhr.status == 500) {
    detail.shouldSwap = true;
    const secretResultForm = document.getElementById('secret-result-form');
    if (secretResultForm) {
      secretResultForm.setAttribute('hx-swap', 'beforeend');
    }
  }
});

// Handle events after htmx swap for secret result.
document.addEventListener('htmx:afterSwap', (event) => {
  const target = event.target;
  if (target.id == 'secret-result-container') {
    const copySecretResultValue = document.getElementById('copy-secret-result-value');
    if (copySecretResultValue) {
      copySecretResultValue.addEventListener('click', () => {
        copyToClipboard('secret-result-value', 'copy-secret-result-value');
      });
    }
  }
});

// Event listener for htmx response error.
document.addEventListener('htmx:responseError', (event) => {
  const detail = event.detail;
  if (detail.target.id == 'secret-result-container') {
    const status = detail.xhr.status;
    const secretResultFormPassphrase = document.getElementById('secret-result-form-passphrase');

    if (status == 401) {
      secretResultFormPassphrase.addEventListener('click', () => {
        secretResultFormPassphrase.setAttribute('placeholder', 'Passphrase');
        secretResultFormPassphrase.classList.remove('placeholder-red-500');
        secretResultFormPassphrase.classList.add('placeholder-gray-400');
      });

      secretResultFormPassphrase.setAttribute('placeholder', 'Invalid passphrase');
      secretResultFormPassphrase.classList.remove('placeholder-gray-400');
      secretResultFormPassphrase.classList.add('placeholder-red-500');
    }
  }
});

// Event listener for secret form. Validate fields and handle character counter.
document.addEventListener('input', () => {
  const secretForm = document.getElementById('secret-form');
  if (!secretForm) {
    return;
  }

  const secretFormTextArea = document.getElementById('secret-form-textarea');
  if (secretFormTextArea) {
    const validate = validateSecretValue(secretFormTextArea.value);
    if (validate && !validate.valid) {
      secretFormTextArea.value = '';
      secretFormTextArea.setAttribute('placeholder', validate.message );
      secretFormTextArea.classList.remove('placeholder-gray-400');
      secretFormTextArea.classList.add('placeholder-red-500');
      disableElement('secret-form-submit');
    } else {
      secretFormTextArea.setAttribute('placeholder', 'Secret value...');
      secretFormTextArea.classList.remove('placeholder-red-500');
      secretFormTextArea.classList.add('placeholder-gray-400');
      enableElement('secret-form-submit');
    }
  }

  const secretFormTextAreaCounter = document.getElementById('secret-form-textarea-counter')
  if (secretFormTextArea && secretFormTextAreaCounter) {
    const length = secretFormTextArea.value.length;
    secretFormTextAreaCounter.textContent = length + '/' + maxSecretValueLength;
  }
});

// Event listener for secret result form (passphrase) to reset the form after submission.
document.addEventListener('submit', () => {
  const secretResultForm = document.getElementById('secret-result-form');
  if (secretResultForm) {
    secretResultForm.reset();
  }
});

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
      }, 1500);
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

// isEmpty checks if a value is empty including null, undefined, empty string, and null byte.
function isEmpty(v) {
  return v == null || v === undefined || v.trim() === '' || v === '\\x00';
}

// validateSecretValue validates the secret value.
function validateSecretValue(v) {
  if (isEmpty(v)) {
    return { valid: false, message: 'Secret value cannot be empty.' };
  }
  if (v.length > maxSecretValueLength) {
    return { valid: false, message: 'Secret value exceeds the maximum length of '+ maxSecretValueLength + ' characters.' };
  }
  return { valid: true };
}

window.setBaseUrl = setBaseUrl;
window.copyToClipboard = copyToClipboard;
window.disableElement = disableElement;
window.enableElement = enableElement;
