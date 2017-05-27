"use strict";
require('./polyfills.ts');
require('zone.js/dist/long-stack-trace-zone');
require('zone.js/dist/proxy.js');
require('zone.js/dist/sync-test');
require('zone.js/dist/jasmine-patch');
require('zone.js/dist/async-test');
require('zone.js/dist/fake-async-test');
// Prevent Karma from running prematurely.
__karma__.loaded = function () { };
Promise.all([
    System.import('@angular/core/testing'),
    System.import('@angular/platform-browser-dynamic/testing')
])
    .then(function (_a) {
    var testing = _a[0], testingBrowser = _a[1];
    testing.getTestBed().initTestEnvironment(testingBrowser.BrowserDynamicTestingModule, testingBrowser.platformBrowserDynamicTesting());
})
    .then(function () { return require.context('./', true, /\.spec\.ts/); })
    .then(function (context) { return context.keys().map(context); })
    .then(__karma__.start, __karma__.error);
//# sourceMappingURL=test.js.map