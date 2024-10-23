// Create a variable to store the base URL. This will be set to:
// - Secret create form.
document.addEventListener('DOMContentLoaded', ()=>{
  const port = window.location.port;
  let baseUrl = window.location.protocol + '//' + window.location.hostname;
  if (port && port !== '80' && port !== '443') {
    baseUrl += ':' + port;
  }
  
  
  secretFormBaseUrl = document.getElementById('secret-form-base-url');
  if (secretFormBaseUrl) {
    secretFormBaseUrl.value = baseUrl;
  }
});

// copyToClipboard copies the contents of an element to the clipboard.
function copyToClipboard(elementId, feedbackElementId) {
  const element = document.getElementById(elementId);
  if (!element) {
    return;
  }

  const text = element.innerText || element.textContent || element.value;

  navigator.clipboard.writeText(text).then(() => {
    if (feedbackElementId) {
      const feedback = document.getElementById(feedbackElementId);
      feedback.innerText = "Copied to clipboard!";
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
