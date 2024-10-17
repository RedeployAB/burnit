document.addEventListener('DOMContentLoaded', function(){
  const baseUrl = window.location.protocol + '//' + window.location.hostname;
  
  secretCreateBaseUrl = document.getElementById('secret-create-base-url');
  if (secretCreateBaseUrl) {
    secretCreateBaseUrl.value = baseUrl;
  }
});
