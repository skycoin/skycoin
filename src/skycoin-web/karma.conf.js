// Karma configuration file, see link for more information
// https://karma-runner.github.io/0.13/config/configuration-file.html

module.exports = function (config) {
 
  if (process.argv.indexOf('--cipher') > -1) {
    var testCipher = true;
  }

  config.set({
    basePath: '',
    frameworks: ['jasmine'],
    plugins: [
      require('karma-jasmine'),
      require('karma-chrome-launcher'),
      require('karma-jasmine-html-reporter'),
      require('karma-coverage-istanbul-reporter')
    ],
    // The @angular/build:karma builder ignores "files" and "proxies"; anything
    // the specs need to fetch is declared as a script or an asset on the test
    // target in angular.json.
    client: {
      // cipher.provider.lib.spec.ts walks a seed chain: each address is derived
      // from the seed the previous one returned, so the specs have to run in
      // declaration order. Jasmine randomises by default.
      jasmine: { random: false },
      // this works only with `karma start`, not `karma run`.
      cipher: testCipher,
      args: ['--browserNoActivityTimeout', config.browserNoActivityTimeout],
      clearContext: false // leave Jasmine Spec Runner output visible in browser
    },
    coverageIstanbulReporter: {
      reports: [ 'html', 'lcovonly' ],
      fixWebpackSourcePaths: true
    },
    reporters: ['progress', 'kjhtml'],
    port: 9876,
    // cipher.provider.lib.spec.ts loads a ~5 MB wasm module and then runs the
    // full Go cipher testsuite vectors; the 30s default disconnects mid-run.
    browserNoActivityTimeout: 180000,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: true,
    browsers: ['ChromeHeadless', 'ChromeHeadlessNoSandbox', 'Chrome'],
    customLaunchers: {
      ChromeHeadlessNoSandbox: {
        base: 'ChromeHeadless',
        flags: ['--no-sandbox']
      }
    },
    singleRun: false
  });
};
