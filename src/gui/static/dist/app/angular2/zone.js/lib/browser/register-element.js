/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
System.register(['../common/utils', './define-property'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var utils_1, define_property_1;
    function registerElementPatch(_global) {
        if (!utils_1.isBrowser || !('registerElement' in _global.document)) {
            return;
        }
        const _registerElement = document.registerElement;
        const callbacks = ['createdCallback', 'attachedCallback', 'detachedCallback', 'attributeChangedCallback'];
        document.registerElement = function (name, opts) {
            if (opts && opts.prototype) {
                callbacks.forEach(function (callback) {
                    const source = 'Document.registerElement::' + callback;
                    if (opts.prototype.hasOwnProperty(callback)) {
                        const descriptor = Object.getOwnPropertyDescriptor(opts.prototype, callback);
                        if (descriptor && descriptor.value) {
                            descriptor.value = Zone.current.wrap(descriptor.value, source);
                            define_property_1._redefineProperty(opts.prototype, callback, descriptor);
                        }
                        else {
                            opts.prototype[callback] = Zone.current.wrap(opts.prototype[callback], source);
                        }
                    }
                    else if (opts.prototype[callback]) {
                        opts.prototype[callback] = Zone.current.wrap(opts.prototype[callback], source);
                    }
                });
            }
            return _registerElement.apply(document, [name, opts]);
        };
    }
    exports_1("registerElementPatch", registerElementPatch);
    return {
        setters:[
            function (utils_1_1) {
                utils_1 = utils_1_1;
            },
            function (define_property_1_1) {
                define_property_1 = define_property_1_1;
            }],
        execute: function() {
        }
    }
});

//# sourceMappingURL=register-element.js.map
