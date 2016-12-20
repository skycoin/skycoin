"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var core_1 = require('@angular/core');
var common_1 = require('@angular/common');
var dom_adapter_1 = require('../../dom/dom_adapter');
var BrowserPlatformLocation = (function (_super) {
    __extends(BrowserPlatformLocation, _super);
    function BrowserPlatformLocation() {
        _super.call(this);
        this._init();
    }
    // This is moved to its own method so that `MockPlatformLocationStrategy` can overwrite it
    /** @internal */
    BrowserPlatformLocation.prototype._init = function () {
        this._location = dom_adapter_1.getDOM().getLocation();
        this._history = dom_adapter_1.getDOM().getHistory();
    };
    Object.defineProperty(BrowserPlatformLocation.prototype, "location", {
        /** @internal */
        get: function () { return this._location; },
        enumerable: true,
        configurable: true
    });
    BrowserPlatformLocation.prototype.getBaseHrefFromDOM = function () { return dom_adapter_1.getDOM().getBaseHref(); };
    BrowserPlatformLocation.prototype.onPopState = function (fn) {
        dom_adapter_1.getDOM().getGlobalEventTarget('window').addEventListener('popstate', fn, false);
    };
    BrowserPlatformLocation.prototype.onHashChange = function (fn) {
        dom_adapter_1.getDOM().getGlobalEventTarget('window').addEventListener('hashchange', fn, false);
    };
    Object.defineProperty(BrowserPlatformLocation.prototype, "pathname", {
        get: function () { return this._location.pathname; },
        set: function (newPath) { this._location.pathname = newPath; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(BrowserPlatformLocation.prototype, "search", {
        get: function () { return this._location.search; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(BrowserPlatformLocation.prototype, "hash", {
        get: function () { return this._location.hash; },
        enumerable: true,
        configurable: true
    });
    BrowserPlatformLocation.prototype.pushState = function (state, title, url) {
        this._history.pushState(state, title, url);
    };
    BrowserPlatformLocation.prototype.replaceState = function (state, title, url) {
        this._history.replaceState(state, title, url);
    };
    BrowserPlatformLocation.prototype.forward = function () { this._history.forward(); };
    BrowserPlatformLocation.prototype.back = function () { this._history.back(); };
    BrowserPlatformLocation.decorators = [
        { type: core_1.Injectable },
    ];
    BrowserPlatformLocation.ctorParameters = [];
    return BrowserPlatformLocation;
}(common_1.PlatformLocation));
exports.BrowserPlatformLocation = BrowserPlatformLocation;
//# sourceMappingURL=browser_platform_location.js.map