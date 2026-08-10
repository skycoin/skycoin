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
// Which specs run is controlled by the "include" list of the karma builder in
// angular.json, not by a glob over the whole project.
//
// The ~60 other specs here still import `async` from @angular/core/testing
// (removed in Angular 12), the @angular/material barrel (removed in v9) and
// @angular/http (removed in v8), together with app/utils/test-mocks.ts which
// supports them. Compiling them fails the suite before a single assertion runs,
// which is why this project had no working tests at all. They stay in the tree
// for whoever revives them; add each one back to that "include" list, and to
// src/tsconfig.spec.json, as it is rewritten against current testing APIs.

import './app/components/layout/msg-bar/msg-bar.component.spec';

__karma__.start();
