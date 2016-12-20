/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
System.register(['./utils'], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var utils_1;
    function patchTimer(window, setName, cancelName, nameSuffix) {
        var setNative = null;
        var clearNative = null;
        setName += nameSuffix;
        cancelName += nameSuffix;
        const tasksByHandleId = {};
        function scheduleTask(task) {
            const data = task.data;
            data.args[0] = function () {
                task.invoke.apply(this, arguments);
                delete tasksByHandleId[data.handleId];
            };
            data.handleId = setNative.apply(window, data.args);
            tasksByHandleId[data.handleId] = task;
            return task;
        }
        function clearTask(task) {
            delete tasksByHandleId[task.data.handleId];
            return clearNative(task.data.handleId);
        }
        setNative =
            utils_1.patchMethod(window, setName, (delegate) => function (self, args) {
                if (typeof args[0] === 'function') {
                    var zone = Zone.current;
                    var options = {
                        handleId: null,
                        isPeriodic: nameSuffix === 'Interval',
                        delay: (nameSuffix === 'Timeout' || nameSuffix === 'Interval') ? args[1] || 0 : null,
                        args: args
                    };
                    var task = zone.scheduleMacroTask(setName, args[0], options, scheduleTask, clearTask);
                    if (!task) {
                        return task;
                    }
                    // Node.js must additionally support the ref and unref functions.
                    var handle = task.data.handleId;
                    if (handle.ref && handle.unref) {
                        task.ref = handle.ref.bind(handle);
                        task.unref = handle.unref.bind(handle);
                    }
                    return task;
                }
                else {
                    // cause an error by calling it directly.
                    return delegate.apply(window, args);
                }
            });
        clearNative =
            utils_1.patchMethod(window, cancelName, (delegate) => function (self, args) {
                var task = typeof args[0] === 'number' ? tasksByHandleId[args[0]] : args[0];
                if (task && typeof task.type === 'string') {
                    if (task.cancelFn && task.data.isPeriodic || task.runCount === 0) {
                        // Do not cancel already canceled functions
                        task.zone.cancelTask(task);
                    }
                }
                else {
                    // cause an error by calling it directly.
                    delegate.apply(window, args);
                }
            });
    }
    exports_1("patchTimer", patchTimer);
    return {
        setters:[
            function (utils_1_1) {
                utils_1 = utils_1_1;
            }],
        execute: function() {
        }
    }
});

//# sourceMappingURL=timers.js.map
