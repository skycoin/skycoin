// Runs the ordinary send browser spec. ci-scripts/web-wallet-send-e2e.sh
// starts the chain and the web wallet and runs this with explorer's WebdriverIO.
const chromeOptions: any = {
  args: ['--no-sandbox', '--headless', '--disable-gpu', 'window-size=1400,1000'],
};
if (process.env.CHROME_BIN) {
  chromeOptions.binary = process.env.CHROME_BIN;
}

exports.config = {
  runner: 'local',
  specs: ['./send.e2e-spec.ts'],
  maxInstances: 1,
  capabilities: [{
    browserName: 'chrome',
    'goog:chromeOptions': chromeOptions,
  }],
  logLevel: 'warn',
  strictSelectors: false,
  bail: 1,
  baseUrl: 'http://127.0.0.1:8001/angular/',
  waitforTimeout: 20000,
  connectionRetryTimeout: 120000,
  connectionRetryCount: 3,
  services: [],
  framework: 'jasmine',
  reporters: ['spec'],
  jasmineOpts: {
    defaultTimeoutInterval: 180000,
  },
};
