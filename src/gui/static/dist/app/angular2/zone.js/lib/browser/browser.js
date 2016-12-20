/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
System.register(['../zone', '../common/timers', '../common/utils', './define-property', './event-target', './property-descriptor', './register-element'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var timers_1, utils_1, define_property_1, event_target_1, property_descriptor_1, register_element_1;
    var set, clear, blockingMethods, _global, i, name, XHR_TASK, XHR_SYNC;
    function patchXHR(window) {
        function findPendingTask(target) {
            var pendingTask = target[XHR_TASK];
            return pendingTask;
        }
        function scheduleTask(task) {
            var data = task.data;
            data.target.addEventListener('readystatechange', () => {
                if (data.target.readyState === data.target.DONE) {
                    if (!data.aborted) {
                        task.invoke();
                    }
                }
            });
            var storedTask = data.target[XHR_TASK];
            if (!storedTask) {
                data.target[XHR_TASK] = task;
            }
            sendNative.apply(data.target, data.args);
            return task;
        }
        function placeholderCallback() { }
        function clearTask(task) {
            var data = task.data;
            // Note - ideally, we would call data.target.removeEventListener here, but it's too late
            // to prevent it from firing. So instead, we store info for the event listener.
            data.aborted = true;
            return abortNative.apply(data.target, data.args);
        }
        var openNative = utils_1.patchMethod(window.XMLHttpRequest.prototype, 'open', () => function (self, args) {
            self[XHR_SYNC] = args[2] == false;
            return openNative.apply(self, args);
        });
        var sendNative = utils_1.patchMethod(window.XMLHttpRequest.prototype, 'send', () => function (self, args) {
            var zone = Zone.current;
            if (self[XHR_SYNC]) {
                // if the XHR is sync there is no task to schedule, just execute the code.
                return sendNative.apply(self, args);
            }
            else {
                var options = { target: self, isPeriodic: false, delay: null, args: args, aborted: false };
                return zone.scheduleMacroTask('XMLHttpRequest.send', placeholderCallback, options, scheduleTask, clearTask);
            }
        });
        var abortNative = utils_1.patchMethod(window.XMLHttpRequest.prototype, 'abort', (delegate) => function (self, args) {
            var task = findPendingTask(self);
            if (task && typeof task.type == 'string') {
                // If the XHR has already completed, do nothing.
                if (task.cancelFn == null) {
                    return;
                }
                task.zone.cancelTask(task);
            }
            // Otherwise, we are trying to abort an XHR which has not yet been sent, so there is no task
            // to cancel. Do nothing.
        });
    }
    return {
        setters:[
            function (_1) {},
            function (timers_1_1) {
                timers_1 = timers_1_1;
            },
            function (utils_1_1) {
                utils_1 = utils_1_1;
            },
            function (define_property_1_1) {
                define_property_1 = define_property_1_1;
            },
            function (event_target_1_1) {
                event_target_1 = event_target_1_1;
            },
            function (property_descriptor_1_1) {
                property_descriptor_1 = property_descriptor_1_1;
            },
            function (register_element_1_1) {
                register_element_1 = register_element_1_1;
            }],
        execute: function() {
            set = 'set';
            clear = 'clear';
            blockingMethods = ['alert', 'prompt', 'confirm'];
            _global = typeof window === 'object' && window || typeof self === 'object' && self || global;
            timers_1.patchTimer(_global, set, clear, 'Timeout');
            timers_1.patchTimer(_global, set, clear, 'Interval');
            timers_1.patchTimer(_global, set, clear, 'Immediate');
            timers_1.patchTimer(_global, 'request', 'cancel', 'AnimationFrame');
            timers_1.patchTimer(_global, 'mozRequest', 'mozCancel', 'AnimationFrame');
            timers_1.patchTimer(_global, 'webkitRequest', 'webkitCancel', 'AnimationFrame');
            for (i = 0; i < blockingMethods.length; i++) {
                name = blockingMethods[i];
                utils_1.patchMethod(_global, name, (delegate, symbol, name) => {
                    return function (s, args) {
                        return Zone.current.run(delegate, _global, args, name);
                    };
                });
            }
            event_target_1.eventTargetPatch(_global);
            property_descriptor_1.propertyDescriptorPatch(_global);
            utils_1.patchClass('MutationObserver');
            utils_1.patchClass('WebKitMutationObserver');
            utils_1.patchClass('FileReader');
            define_property_1.propertyPatch();
            register_element_1.registerElementPatch(_global);
            // Treat XMLHTTPRequest as a macrotask.
            patchXHR(_global);
            XHR_TASK = utils_1.zoneSymbol('xhrTask');
            XHR_SYNC = utils_1.zoneSymbol('xhrSync');
            /// GEO_LOCATION
            if (_global['navigator'] && _global['navigator'].geolocation) {
                utils_1.patchPrototype(_global['navigator'].geolocation, ['getCurrentPosition', 'watchPosition']);
            }
        }
    }
});

//# sourceMappingURL=browser.js.map
