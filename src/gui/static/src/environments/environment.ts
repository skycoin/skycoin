// The file contents for the current environment will overwrite these during build.
// The build system defaults to the dev environment which uses `environment.ts`, but if you do
// `ng build --env=prod` then `environment.prod.ts` will be used instead.
// The list of which env maps to which file can be found in `.angular-cli.json`.

export const environment = {
  nodeUrl: '/api/',
  production: false,
  tellerUrl: '/teller/',
  isInE2eMode: false,

  swaplab: {
    apiKey: 'w4bxe2tbf9beb72r', // if set to null, integration will be disabled
    activateTestMode: true,
    endStatusInError: false,
  },
};
