/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { ReflectiveInjector } from '@angular/core';
import { PropertyBinding } from './component_info';
import { $SCOPE } from './constants';
var /** @type {?} */ INITIAL_VALUE = {
    __UNINITIALIZED__: true
};
export var DowngradeComponentAdapter = (function () {
    /**
     * @param {?} id
     * @param {?} info
     * @param {?} element
     * @param {?} attrs
     * @param {?} scope
     * @param {?} parentInjector
     * @param {?} parse
     * @param {?} componentFactory
     */
    function DowngradeComponentAdapter(id, info, element, attrs, scope, parentInjector, parse, componentFactory) {
        this.id = id;
        this.info = info;
        this.element = element;
        this.attrs = attrs;
        this.scope = scope;
        this.parentInjector = parentInjector;
        this.parse = parse;
        this.componentFactory = componentFactory;
        this.component = null;
        this.inputChangeCount = 0;
        this.inputChanges = null;
        this.componentRef = null;
        this.changeDetector = null;
        this.contentInsertionPoint = null;
        this.element[0].id = id;
        this.componentScope = scope.$new();
        this.childNodes = element.contents();
    }
    /**
     * @return {?}
     */
    DowngradeComponentAdapter.prototype.createComponent = function () {
        var /** @type {?} */ childInjector = ReflectiveInjector.resolveAndCreate([{ provide: $SCOPE, useValue: this.componentScope }], this.parentInjector);
        this.contentInsertionPoint = document.createComment('ng1 insertion point');
        this.componentRef = this.componentFactory.create(childInjector, [[this.contentInsertionPoint]], this.element[0]);
        this.changeDetector = this.componentRef.changeDetectorRef;
        this.component = this.componentRef.instance;
    };
    /**
     * @return {?}
     */
    DowngradeComponentAdapter.prototype.setupInputs = function () {
        var _this = this;
        var /** @type {?} */ attrs = this.attrs;
        var /** @type {?} */ inputs = this.info.inputs || [];
        for (var /** @type {?} */ i = 0; i < inputs.length; i++) {
            var /** @type {?} */ input = new PropertyBinding(inputs[i]);
            var /** @type {?} */ expr = null;
            if (attrs.hasOwnProperty(input.attr)) {
                var /** @type {?} */ observeFn = (function (prop /** TODO #9100 */) {
                    var /** @type {?} */ prevValue = INITIAL_VALUE;
                    return function (value /** TODO #9100 */) {
                        if (_this.inputChanges !== null) {
                            _this.inputChangeCount++;
                            _this.inputChanges[prop] =
                                new Ng1Change(value, prevValue === INITIAL_VALUE ? value : prevValue);
                            prevValue = value;
                        }
                        _this.component[prop] = value;
                    };
                })(input.prop);
                attrs.$observe(input.attr, observeFn);
            }
            else if (attrs.hasOwnProperty(input.bindAttr)) {
                expr = ((attrs) /** TODO #9100 */)[input.bindAttr];
            }
            else if (attrs.hasOwnProperty(input.bracketAttr)) {
                expr = ((attrs) /** TODO #9100 */)[input.bracketAttr];
            }
            else if (attrs.hasOwnProperty(input.bindonAttr)) {
                expr = ((attrs) /** TODO #9100 */)[input.bindonAttr];
            }
            else if (attrs.hasOwnProperty(input.bracketParenAttr)) {
                expr = ((attrs) /** TODO #9100 */)[input.bracketParenAttr];
            }
            if (expr != null) {
                var /** @type {?} */ watchFn = (function (prop /** TODO #9100 */) {
                    return function (value /** TODO #9100 */, prevValue /** TODO #9100 */) {
                        if (_this.inputChanges != null) {
                            _this.inputChangeCount++;
                            _this.inputChanges[prop] = new Ng1Change(prevValue, value);
                        }
                        _this.component[prop] = value;
                    };
                })(input.prop);
                this.componentScope.$watch(expr, watchFn);
            }
        }
        var /** @type {?} */ prototype = this.info.component.prototype;
        if (prototype && ((prototype)).ngOnChanges) {
            // Detect: OnChanges interface
            this.inputChanges = {};
            this.componentScope.$watch(function () { return _this.inputChangeCount; }, function () {
                var /** @type {?} */ inputChanges = _this.inputChanges;
                _this.inputChanges = {};
                ((_this.component)).ngOnChanges(inputChanges);
            });
        }
        this.componentScope.$watch(function () { return _this.changeDetector && _this.changeDetector.detectChanges(); });
    };
    /**
     * @return {?}
     */
    DowngradeComponentAdapter.prototype.projectContent = function () {
        var /** @type {?} */ childNodes = this.childNodes;
        var /** @type {?} */ parent = this.contentInsertionPoint.parentNode;
        if (parent) {
            for (var /** @type {?} */ i = 0, /** @type {?} */ ii = childNodes.length; i < ii; i++) {
                parent.insertBefore(childNodes[i], this.contentInsertionPoint);
            }
        }
    };
    /**
     * @return {?}
     */
    DowngradeComponentAdapter.prototype.setupOutputs = function () {
        var _this = this;
        var /** @type {?} */ attrs = this.attrs;
        var /** @type {?} */ outputs = this.info.outputs || [];
        for (var /** @type {?} */ j = 0; j < outputs.length; j++) {
            var /** @type {?} */ output = new PropertyBinding(outputs[j]);
            var /** @type {?} */ expr = null;
            var /** @type {?} */ assignExpr = false;
            var /** @type {?} */ bindonAttr = output.bindonAttr ? output.bindonAttr.substring(0, output.bindonAttr.length - 6) : null;
            var /** @type {?} */ bracketParenAttr = output.bracketParenAttr ?
                "[(" + output.bracketParenAttr.substring(2, output.bracketParenAttr.length - 8) + ")]" :
                null;
            if (attrs.hasOwnProperty(output.onAttr)) {
                expr = ((attrs) /** TODO #9100 */)[output.onAttr];
            }
            else if (attrs.hasOwnProperty(output.parenAttr)) {
                expr = ((attrs) /** TODO #9100 */)[output.parenAttr];
            }
            else if (attrs.hasOwnProperty(bindonAttr)) {
                expr = ((attrs) /** TODO #9100 */)[bindonAttr];
                assignExpr = true;
            }
            else if (attrs.hasOwnProperty(bracketParenAttr)) {
                expr = ((attrs) /** TODO #9100 */)[bracketParenAttr];
                assignExpr = true;
            }
            if (expr != null && assignExpr != null) {
                var /** @type {?} */ getter = this.parse(expr);
                var /** @type {?} */ setter = getter.assign;
                if (assignExpr && !setter) {
                    throw new Error("Expression '" + expr + "' is not assignable!");
                }
                var /** @type {?} */ emitter = (this.component[output.prop]);
                if (emitter) {
                    emitter.subscribe({
                        next: assignExpr ?
                            (function (setter) { return function (v /** TODO #9100 */) { return setter(_this.scope, v); }; })(setter) :
                            (function (getter) { return function (v /** TODO #9100 */) {
                                return getter(_this.scope, { $event: v });
                            }; })(getter)
                    });
                }
                else {
                    throw new Error("Missing emitter '" + output.prop + "' on component '" + this.info.component + "'!");
                }
            }
        }
    };
    /**
     * @return {?}
     */
    DowngradeComponentAdapter.prototype.registerCleanup = function () {
        var _this = this;
        this.element.bind('$destroy', function () {
            _this.componentScope.$destroy();
            _this.componentRef.destroy();
        });
    };
    return DowngradeComponentAdapter;
}());
function DowngradeComponentAdapter_tsickle_Closure_declarations() {
    /** @type {?} */
    DowngradeComponentAdapter.prototype.component;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.inputs;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.inputChangeCount;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.inputChanges;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.componentRef;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.changeDetector;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.componentScope;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.childNodes;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.contentInsertionPoint;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.id;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.info;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.element;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.attrs;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.scope;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.parentInjector;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.parse;
    /** @type {?} */
    DowngradeComponentAdapter.prototype.componentFactory;
}
var Ng1Change = (function () {
    /**
     * @param {?} previousValue
     * @param {?} currentValue
     */
    function Ng1Change(previousValue, currentValue) {
        this.previousValue = previousValue;
        this.currentValue = currentValue;
    }
    /**
     * @return {?}
     */
    Ng1Change.prototype.isFirstChange = function () { return this.previousValue === this.currentValue; };
    return Ng1Change;
}());
function Ng1Change_tsickle_Closure_declarations() {
    /** @type {?} */
    Ng1Change.prototype.previousValue;
    /** @type {?} */
    Ng1Change.prototype.currentValue;
}
//# sourceMappingURL=downgrade_component_adapter.js.map