/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { NgModule } from '@angular/core';
import { COMMON_DIRECTIVES } from './directives/index';
import { NgLocaleLocalization, NgLocalization } from './localization';
import { COMMON_PIPES } from './pipes/index';
// Note: This does not contain the location providers,
// as they need some platform specific implementations to work.
/**
 * The module that includes all the basic Angular directives like {@link NgIf}, {@link NgFor}, ...
 *
 * @stable
 */
export var CommonModule = (function () {
    function CommonModule() {
    }
    CommonModule.decorators = [
        { type: NgModule, args: [{
                    declarations: [COMMON_DIRECTIVES, COMMON_PIPES],
                    exports: [COMMON_DIRECTIVES, COMMON_PIPES],
                    providers: [
                        { provide: NgLocalization, useClass: NgLocaleLocalization },
                    ],
                },] },
    ];
    /** @nocollapse */
    CommonModule.ctorParameters = [];
    return CommonModule;
}());
//# sourceMappingURL=common_module.js.map