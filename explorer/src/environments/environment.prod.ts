export const environment = {
  production: true
};

// Create a "link" element for activating the integration with the browser search bar.
const searchLink = document.createElement('link');
searchLink.type = 'application/opensearchdescription+xml';
searchLink.rel = 'search';
searchLink.href = 'search.xml';
searchLink.title = 'Skycoin Explorer';
document.head.appendChild(searchLink);
