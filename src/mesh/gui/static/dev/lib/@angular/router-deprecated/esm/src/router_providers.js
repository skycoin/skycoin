import { ROUTER_PROVIDERS_COMMON } from './router_providers_common';
import { BrowserPlatformLocation } from '@angular/platform-browser';
import { PlatformLocation } from '@angular/common';
/**
 * A list of {@link Provider}s. To use the router, you must add this to your application.
 *
 * ### Example ([live demo](http://plnkr.co/edit/iRUP8B5OUbxCWQ3AcIDm))
 *
 * ```
 * import {Component} from '@angular/core';
 * import {
 *   ROUTER_DIRECTIVES,
 *   ROUTER_PROVIDERS,
 *   RouteConfig
 * } from '@angular/router-deprecated';
 *
 * @Component({directives: [ROUTER_DIRECTIVES]})
 * @RouteConfig([
 *  {...},
 * ])
 * class AppCmp {
 *   // ...
 * }
 *
 * bootstrap(AppCmp, [ROUTER_PROVIDERS]);
 * ```
 */
export const ROUTER_PROVIDERS = [
    ROUTER_PROVIDERS_COMMON,
    /*@ts2dart_const*/ (
    /* @ts2dart_Provider */ { provide: PlatformLocation, useClass: BrowserPlatformLocation }),
];
/**
 * Use {@link ROUTER_PROVIDERS} instead.
 *
 * @deprecated
 */
export const ROUTER_BINDINGS = ROUTER_PROVIDERS;
//# sourceMappingURL=router_providers.js.map