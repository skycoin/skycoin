// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

export const environment = {
  production: false
};

// Create a "link" element for activating the integration with the browser search bar.
const searchLink = document.createElement('link');
searchLink.type = 'application/opensearchdescription+xml';
searchLink.rel = 'search';
searchLink.href = 'search.dev.xml';
searchLink.title = 'Skycoin Explorer Test';
document.head.appendChild(searchLink);
