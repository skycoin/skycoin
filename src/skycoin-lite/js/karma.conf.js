// Karma configuration file, see link for more information
// https://karma-runner.github.io/0.13/config/configuration-file.html

module.exports = function (config) {
 
  var cipherParamIndex = process.argv.indexOf('--mode');
  // check if command line has cipher parameter with not empty value
  if (cipherParamIndex > -1 && (cipherParamIndex + 1) < process.argv.length && process.argv[cipherParamIndex + 1]) {
    var cipherMode = process.argv[cipherParamIndex + 1];
  }

  config.set({
    basePath: '',
    frameworks: ['jasmine', 'karma-typescript'],
    plugins: [
      require('karma-jasmine'),
      require('karma-chrome-launcher'),
      require('karma-jasmine-html-reporter'),
      require('karma-read-json'),
      require('karma-typescript')
    ],
    files: [
      'tests/*.spec.ts',
      { pattern: 'tests/test-fixtures/*.golden', included: false },
      { pattern: 'tests/*.ts', included: true },
      { pattern: 'skycoin.js', included: true }
    ],
    preprocessors: {
      "**/*.ts": "karma-typescript"
    },
    client: {
      clearContext: false // leave Jasmine Spec Runner output visible in browser
    },
    reporters: ['progress', 'kjhtml', 'karma-typescript'],
    karmaTypescriptConfig: {
      bundlerOptions: {
        constants: {
          "TESTING_MODE": cipherMode
        }
      }
    },
    port: 9876,
    colors: true,
    logLevel: config.LOG_INFO,
    autoWatch: true,
    browsers: ['ChromeHeadless', 'Chrome'],
    singleRun: false
  });
};
