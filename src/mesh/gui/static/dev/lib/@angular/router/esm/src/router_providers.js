import { ROUTER_PROVIDERS_COMMON } from './router_providers_common';
import { BrowserPlatformLocation } from '@angular/platform-browser';
import { PlatformLocation } from '@angular/common';
/**
 * A list of {@link Provider}s. To use the router, you must add this to your application.
 *
 * ```
 * import {Component} from '@angular/core';
 * import {
 *   ROUTER_DIRECTIVES,
 *   ROUTER_PROVIDERS,
 *   Routes
 * } from '@angular/router';
 *
 * @Component({directives: [ROUTER_DIRECTIVES]})
 * @Routes([
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
    /*@ts2dart_Provider*/ { provide: PlatformLocation, useClass: BrowserPlatformLocation },
];
//# sourceMappingURL=router_providers.js.map