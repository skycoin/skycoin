/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
/**
 *  A `PropertyBinding` represents a mapping between a property name
  * and an attribute name. It is parsed from a string of the form
  * `"prop: attr"`; or simply `"propAndAttr" where the property
  * and attribute have the same identifier.
 */
export var PropertyBinding = (function () {
    /**
     * @param {?} binding
     */
    function PropertyBinding(binding) {
        this.binding = binding;
        this.parseBinding();
    }
    /**
     * @return {?}
     */
    PropertyBinding.prototype.parseBinding = function () {
        var /** @type {?} */ parts = this.binding.split(':');
        this.prop = parts[0].trim();
        this.attr = (parts[1] || this.prop).trim();
        this.bracketAttr = "[" + this.attr + "]";
        this.parenAttr = "(" + this.attr + ")";
        this.bracketParenAttr = "[(" + this.attr + ")]";
        var /** @type {?} */ capitalAttr = this.attr.charAt(0).toUpperCase() + this.attr.substr(1);
        this.onAttr = "on" + capitalAttr;
        this.bindAttr = "bind" + capitalAttr;
        this.bindonAttr = "bindon" + capitalAttr;
    };
    return PropertyBinding;
}());
function PropertyBinding_tsickle_Closure_declarations() {
    /** @type {?} */
    PropertyBinding.prototype.prop;
    /** @type {?} */
    PropertyBinding.prototype.attr;
    /** @type {?} */
    PropertyBinding.prototype.bracketAttr;
    /** @type {?} */
    PropertyBinding.prototype.bracketParenAttr;
    /** @type {?} */
    PropertyBinding.prototype.parenAttr;
    /** @type {?} */
    PropertyBinding.prototype.onAttr;
    /** @type {?} */
    PropertyBinding.prototype.bindAttr;
    /** @type {?} */
    PropertyBinding.prototype.bindonAttr;
    /** @type {?} */
    PropertyBinding.prototype.binding;
}
//# sourceMappingURL=component_info.js.map