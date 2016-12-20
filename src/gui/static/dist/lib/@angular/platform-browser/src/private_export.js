/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import * as browser from './browser';
import * as browserDomAdapter from './browser/browser_adapter';
import * as location from './browser/location/browser_platform_location';
import * as testability from './browser/testability';
import * as ng_probe from './dom/debug/ng_probe';
import * as dom_adapter from './dom/dom_adapter';
import * as dom_renderer from './dom/dom_renderer';
import * as dom_events from './dom/events/dom_events';
import * as hammer_gesture from './dom/events/hammer_gestures';
import * as key_events from './dom/events/key_events';
import * as shared_styles_host from './dom/shared_styles_host';
import { WebAnimationsDriver } from './dom/web_animations_driver';
export var __platform_browser_private__ = {
    BrowserPlatformLocation: location.BrowserPlatformLocation,
    DomAdapter: dom_adapter.DomAdapter,
    BrowserDomAdapter: browserDomAdapter.BrowserDomAdapter,
    BrowserGetTestability: testability.BrowserGetTestability,
    getDOM: dom_adapter.getDOM,
    setRootDomAdapter: dom_adapter.setRootDomAdapter,
    DomRootRenderer_: dom_renderer.DomRootRenderer_,
    DomRootRenderer: dom_renderer.DomRootRenderer,
    NAMESPACE_URIS: dom_renderer.NAMESPACE_URIS,
    shimContentAttribute: dom_renderer.shimContentAttribute,
    shimHostAttribute: dom_renderer.shimHostAttribute,
    flattenStyles: dom_renderer.flattenStyles,
    splitNamespace: dom_renderer.splitNamespace,
    isNamespaced: dom_renderer.isNamespaced,
    DomSharedStylesHost: shared_styles_host.DomSharedStylesHost,
    SharedStylesHost: shared_styles_host.SharedStylesHost,
    ELEMENT_PROBE_PROVIDERS: ng_probe.ELEMENT_PROBE_PROVIDERS,
    DomEventsPlugin: dom_events.DomEventsPlugin,
    KeyEventsPlugin: key_events.KeyEventsPlugin,
    HammerGesturesPlugin: hammer_gesture.HammerGesturesPlugin,
    initDomAdapter: browser.initDomAdapter,
    INTERNAL_BROWSER_PLATFORM_PROVIDERS: browser.INTERNAL_BROWSER_PLATFORM_PROVIDERS,
    BROWSER_SANITIZATION_PROVIDERS: browser.BROWSER_SANITIZATION_PROVIDERS,
    WebAnimationsDriver: WebAnimationsDriver
};
//# sourceMappingURL=private_export.js.map