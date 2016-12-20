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
    var fs, err, TO_PATCH;
    return {
        setters:[
            function (utils_1_1) {
                utils_1 = utils_1_1;
            }],
        execute: function() {
            try {
                fs = require('fs');
            }
            catch (err) {
            }
            // TODO(alxhub): Patch `watch` and `unwatchFile`.
            TO_PATCH = [
                'access', 'appendFile', 'chmod', 'chown', 'close', 'exists', 'fchmod',
                'fchown', 'fdatasync', 'fstat', 'fsync', 'ftruncate', 'futimes', 'lchmod',
                'lchown', 'link', 'lstat', 'mkdir', 'mkdtemp', 'open', 'read',
                'readdir', 'readFile', 'readlink', 'realpath', 'rename', 'rmdir', 'stat',
                'symlink', 'truncate', 'unlink', 'utimes', 'write', 'writeFile',
            ];
            if (fs) {
                TO_PATCH.filter(name => !!fs[name] && typeof fs[name] === 'function').forEach(name => {
                    fs[name] = ((delegate) => {
                        return function () {
                            return delegate.apply(this, utils_1.bindArguments(arguments, 'fs.' + name));
                        };
                    })(fs[name]);
                });
            }
        }
    }
});

//# sourceMappingURL=fs.js.map
