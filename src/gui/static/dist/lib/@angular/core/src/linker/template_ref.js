"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var lang_1 = require('../facade/lang');
var EMPTY_CONTEXT = new Object();
/**
 * Represents an Embedded Template that can be used to instantiate Embedded Views.
 *
 * You can access a `TemplateRef`, in two ways. Via a directive placed on a `<template>` element (or
 * directive prefixed with `*`) and have the `TemplateRef` for this Embedded View injected into the
 * constructor of the directive using the `TemplateRef` Token. Alternatively you can query for the
 * `TemplateRef` from a Component or a Directive via {@link Query}.
 *
 * To instantiate Embedded Views based on a Template, use
 * {@link ViewContainerRef#createEmbeddedView}, which will create the View and attach it to the
 * View Container.
 */
var TemplateRef = (function () {
    function TemplateRef() {
    }
    Object.defineProperty(TemplateRef.prototype, "elementRef", {
        /**
         * The location in the View where the Embedded View logically belongs to.
         *
         * The data-binding and injection contexts of Embedded Views created from this `TemplateRef`
         * inherit from the contexts of this location.
         *
         * Typically new Embedded Views are attached to the View Container of this location, but in
         * advanced use-cases, the View can be attached to a different container while keeping the
         * data-binding and injection context from the original location.
         *
         */
        // TODO(i): rename to anchor or location
        get: function () { return null; },
        enumerable: true,
        configurable: true
    });
    return TemplateRef;
}());
exports.TemplateRef = TemplateRef;
var TemplateRef_ = (function (_super) {
    __extends(TemplateRef_, _super);
    function TemplateRef_(_appElement, _viewFactory) {
        _super.call(this);
        this._appElement = _appElement;
        this._viewFactory = _viewFactory;
    }
    TemplateRef_.prototype.createEmbeddedView = function (context) {
        var view = this._viewFactory(this._appElement.parentView.viewUtils, this._appElement.parentInjector, this._appElement);
        if (lang_1.isBlank(context)) {
            context = EMPTY_CONTEXT;
        }
        view.create(context, null, null);
        return view.ref;
    };
    Object.defineProperty(TemplateRef_.prototype, "elementRef", {
        get: function () { return this._appElement.elementRef; },
        enumerable: true,
        configurable: true
    });
    return TemplateRef_;
}(TemplateRef));
exports.TemplateRef_ = TemplateRef_;
//# sourceMappingURL=template_ref.js.map