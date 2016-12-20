/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
System.register(['../common/utils'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var utils_1;
    var EE_ADD_LISTENER, EE_PREPEND_LISTENER, EE_REMOVE_LISTENER, EE_LISTENERS, EE_ON, zoneAwareAddListener, zoneAwarePrependListener, zoneAwareRemoveListener, zoneAwareListeners, events, err;
    function patchEventEmitterMethods(obj) {
        if (obj && obj.addListener) {
            utils_1.patchMethod(obj, EE_ADD_LISTENER, () => zoneAwareAddListener);
            utils_1.patchMethod(obj, EE_PREPEND_LISTENER, () => zoneAwarePrependListener);
            utils_1.patchMethod(obj, EE_REMOVE_LISTENER, () => zoneAwareRemoveListener);
            utils_1.patchMethod(obj, EE_LISTENERS, () => zoneAwareListeners);
            obj[EE_ON] = obj[EE_ADD_LISTENER];
            return true;
        }
        else {
            return false;
        }
    }
    exports_1("patchEventEmitterMethods", patchEventEmitterMethods);
    return {
        setters:[
            function (utils_1_1) {
                utils_1 = utils_1_1;
            }],
        execute: function() {
            // For EventEmitter
            EE_ADD_LISTENER = 'addListener';
            EE_PREPEND_LISTENER = 'prependListener';
            EE_REMOVE_LISTENER = 'removeListener';
            EE_LISTENERS = 'listeners';
            EE_ON = 'on';
            zoneAwareAddListener = utils_1.makeZoneAwareAddListener(EE_ADD_LISTENER, EE_REMOVE_LISTENER, false, true);
            zoneAwarePrependListener = utils_1.makeZoneAwareAddListener(EE_PREPEND_LISTENER, EE_REMOVE_LISTENER, false, true);
            zoneAwareRemoveListener = utils_1.makeZoneAwareRemoveListener(EE_REMOVE_LISTENER, false);
            zoneAwareListeners = utils_1.makeZoneAwareListeners(EE_LISTENERS);
            // EventEmitter
            try {
                events = require('events');
            }
            catch (err) {
            }
            if (events && events.EventEmitter) {
                patchEventEmitterMethods(events.EventEmitter.prototype);
            }
        }
    }
});

//# sourceMappingURL=events.js.map
