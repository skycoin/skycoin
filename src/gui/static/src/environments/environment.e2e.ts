// This file is for the e2e tests ony.

export const environment = {
  nodeUrl: '/api/',
  production: true,
  tellerUrl: '/teller/',
  isInE2eMode: true,

  swaplab: {
    apiKey: 'w4bxe2tbf9beb72r', // if set to null, integration will be disabled
    activateTestMode: true,
    endStatusInError: false,
  },
};
