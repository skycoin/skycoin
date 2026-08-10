// This file is required by karma.conf.js and loads recursively all the .spec and framework files

import 'zone.js/testing';
import { getTestBed } from '@angular/core/testing';
import {
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting
} from '@angular/platform-browser-dynamic/testing';

// Unfortunately there's no typing for the `__karma__` variable. Just declare it as any.
declare const __karma__: any;

// Prevent Karma from running prematurely.
__karma__.loaded = function () {};

// First, initialize the Angular testing environment.
getTestBed().initTestEnvironment(
  BrowserDynamicTestingModule,
  platformBrowserDynamicTesting(), {
    teardown: { destroyAfterEach: false }
}
);
// Then we find all the tests. require.context is a webpack extension, so it is
// not part of the NodeJS Require typings; reach it through a local alias rather
// than pulling in @types/webpack-env just for this file.
interface WebpackRequireContext {
  (id: string): unknown;
  keys(): string[];
}
const webpackRequire = require as unknown as {
  context(path: string, deep?: boolean, filter?: RegExp): WebpackRequireContext;
};

let context: WebpackRequireContext;

if (__karma__.config.cipher) {
  context = webpackRequire.context('./', true, /cipher\.provider\.lib\.spec\.ts$/);
} else {
  context = webpackRequire.context('./', true, /(?!.*?cipher\.provider\.lib\.spec\.ts$)(^.*\.spec\.ts$)/);
}

// Set cipher mode to have ability to read this value in .spec file
process.argv = [ __karma__.config.cipher ];

// And load the modules.
context.keys().map(context);
// Finally, start Karma to run the tests.
__karma__.start();
