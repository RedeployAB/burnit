document.addEventListener('DOMContentLoaded', ()=>{
  const baseUrl = window.location.protocol + '//' + window.location.hostname;
  
  secretCreateBaseUrl = document.getElementById('secret-create-base-url');
  if (secretCreateBaseUrl) {
    secretCreateBaseUrl.value = baseUrl;
  }
});

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
  }).catch((error) => {
    // Replace console.log with better error handling
    // and feedback.
    console.log('Error copying text: ', error);
  });
}
