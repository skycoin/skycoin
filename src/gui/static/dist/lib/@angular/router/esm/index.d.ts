/**
 * @module
 * @description
 * Maps application URLs into application states, to support deep-linking and navigation.
 */
export { Router, RouterOutletMap } from './src/router';
export { RouteSegment, UrlSegment, Tree, UrlTree, RouteTree } from './src/segments';
export { Routes } from './src/metadata/decorators';
export { Route } from './src/metadata/metadata';
export { RouterUrlSerializer, DefaultRouterUrlSerializer } from './src/router_url_serializer';
export { OnActivate, CanDeactivate } from './src/interfaces';
export { ROUTER_PROVIDERS } from './src/router_providers';
/**
 * A list of directives. To use the router directives like {@link RouterOutlet} and
 * {@link RouterLink}, add this to your `directives` array in the {@link View} decorator of your
 * component.
 *
 * ```
 * import {Component} from '@angular/core';
 * import {ROUTER_DIRECTIVES, Routes} from '@angular/router-deprecated';
 *
 * @Component({directives: [ROUTER_DIRECTIVES]})
 * @RouteConfig([
 *  {...},
 * ])
 * class AppCmp {
 *    // ...
 * }
 *
 * bootstrap(AppCmp);
 * ```
 */
export declare const ROUTER_DIRECTIVES: any[];
