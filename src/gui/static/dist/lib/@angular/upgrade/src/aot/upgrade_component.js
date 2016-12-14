/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { EventEmitter } from '@angular/core';
import * as angular from '../angular_js';
import { looseIdentical } from '../facade/lang';
import { controllerKey } from '../util';
import { $COMPILE, $CONTROLLER, $HTTP_BACKEND, $INJECTOR, $SCOPE, $TEMPLATE_CACHE } from './constants';
var /** @type {?} */ REQUIRE_PREFIX_RE = /^(\^\^?)?(\?)?(\^\^?)?/;
var /** @type {?} */ NOT_SUPPORTED = 'NOT_SUPPORTED';
var /** @type {?} */ INITIAL_VALUE = {
    __UNINITIALIZED__: true
};
var Bindings = (function () {
    function Bindings() {
        this.twoWayBoundProperties = [];
        this.twoWayBoundLastValues = [];
        this.expressionBoundProperties = [];
        this.propertyToOutputMap = {};
    }
    return Bindings;
}());
function Bindings_tsickle_Closure_declarations() {
    /** @type {?} */
    Bindings.prototype.twoWayBoundProperties;
    /** @type {?} */
    Bindings.prototype.twoWayBoundLastValues;
    /** @type {?} */
    Bindings.prototype.expressionBoundProperties;
    /** @type {?} */
    Bindings.prototype.propertyToOutputMap;
}
/**
 *  *
  * *Part of the [upgrade/static](/docs/ts/latest/api/#!?query=upgrade%2Fstatic)
  * library for hybrid upgrade apps that support AoT compilation*
  * *
  * Allows an Angular 1 component to be used from Angular 2+.
  * *
  * *
  * Let's assume that you have an Angular 1 component called `ng1Hero` that needs
  * to be made available in Angular 2+ templates.
  * *
  * {@example upgrade/static/ts/module.ts region="ng1-hero"}
  * *
  * We must create a {@link Directive} that will make this Angular 1 component
  * available inside Angular 2+ templates.
  * *
  * {@example upgrade/static/ts/module.ts region="ng1-hero-wrapper"}
  * *
  * In this example you can see that we must derive from the {@link UpgradeComponent}
  * base class but also provide an {@link Directive `@Directive`} decorator. This is
  * because the AoT compiler requires that this information is statically available at
  * compile time.
  * *
  * Note that we must do the following:
  * * specify the directive's selector (`ng1-hero`)
  * * specify all inputs and outputs that the Angular 1 component expects
  * * derive from `UpgradeComponent`
  * * call the base class from the constructor, passing
  * * the Angular 1 name of the component (`ng1Hero`)
  * * the {@link ElementRef} and {@link Injector} for the component wrapper
  * *
  * *
  * A helper class that should be used as a base class for creating Angular directives
  * that wrap Angular 1 components that need to be "upgraded".
  * *
 */
export var UpgradeComponent = (function () {
    /**
     *  Create a new `UpgradeComponent` instance. You should not normally need to do this.
      * Instead you should derive a new class from this one and call the super constructor
      * from the base class.
      * *
      * {@example upgrade/static/ts/module.ts region="ng1-hero-wrapper" }
      * *
      * * The `name` parameter should be the name of the Angular 1 directive.
      * * The `elementRef` and `injector` parameters should be acquired from Angular by dependency
      * injection into the base class constructor.
      * *
      * Note that we must manually implement lifecycle hooks that call through to the super class.
      * This is because, at the moment, the AoT compiler is not able to tell that the
      * `UpgradeComponent`
      * already implements them and so does not wire up calls to them at runtime.
     * @param {?} name
     * @param {?} elementRef
     * @param {?} injector
     */
    function UpgradeComponent(name, elementRef, injector) {
        this.name = name;
        this.elementRef = elementRef;
        this.injector = injector;
        this.controllerInstance = null;
        this.bindingDestination = null;
        this.$injector = injector.get($INJECTOR);
        this.$compile = this.$injector.get($COMPILE);
        this.$templateCache = this.$injector.get($TEMPLATE_CACHE);
        this.$httpBackend = this.$injector.get($HTTP_BACKEND);
        this.$controller = this.$injector.get($CONTROLLER);
        this.element = elementRef.nativeElement;
        this.$element = angular.element(this.element);
        this.directive = this.getDirective(name);
        this.bindings = this.initializeBindings(this.directive);
        this.linkFn = this.compileTemplate(this.directive);
        // We ask for the Angular 1 scope from the Angular 2+ injector, since
        // we will put the new component scope onto the new injector for each component
        var $parentScope = injector.get($SCOPE);
        // QUESTION 1: Should we create an isolated scope if the scope is only true?
        // QUESTION 2: Should we make the scope accessible through `$element.scope()/isolateScope()`?
        this.$componentScope = $parentScope.$new(!!this.directive.scope);
        var controllerType = this.directive.controller;
        var bindToController = this.directive.bindToController;
        if (controllerType) {
            this.controllerInstance = this.buildController(controllerType, this.$componentScope, this.$element, this.directive.controllerAs);
        }
        else if (bindToController) {
            throw new Error("Upgraded directive '" + name + "' specifies 'bindToController' but no controller.");
        }
        this.bindingDestination = bindToController ? this.controllerInstance : this.$componentScope;
        this.setupOutputs();
    }
    /**
     * @return {?}
     */
    UpgradeComponent.prototype.ngOnInit = function () {
        var _this = this;
        var /** @type {?} */ attrs = NOT_SUPPORTED;
        var /** @type {?} */ transcludeFn = NOT_SUPPORTED;
        var /** @type {?} */ directiveRequire = this.getDirectiveRequire(this.directive);
        var /** @type {?} */ requiredControllers = this.resolveRequire(this.directive.name, this.$element, directiveRequire);
        if (this.directive.bindToController && isMap(directiveRequire)) {
            var /** @type {?} */ requiredControllersMap_1 = (requiredControllers);
            Object.keys(requiredControllersMap_1).forEach(function (key) {
                _this.controllerInstance[key] = requiredControllersMap_1[key];
            });
        }
        this.callLifecycleHook('$onInit', this.controllerInstance);
        var /** @type {?} */ link = this.directive.link;
        var /** @type {?} */ preLink = (typeof link == 'object') && ((link)).pre;
        var /** @type {?} */ postLink = (typeof link == 'object') ? ((link)).post : link;
        if (preLink) {
            preLink(this.$componentScope, this.$element, attrs, requiredControllers, transcludeFn);
        }
        var /** @type {?} */ childNodes = [];
        var /** @type {?} */ childNode;
        while (childNode = this.element.firstChild) {
            this.element.removeChild(childNode);
            childNodes.push(childNode);
        }
        var /** @type {?} */ attachElement = function (clonedElements, scope) { _this.$element.append(clonedElements); };
        var /** @type {?} */ attachChildNodes = function (scope, cloneAttach) { return cloneAttach(childNodes); };
        this.linkFn(this.$componentScope, attachElement, { parentBoundTranscludeFn: attachChildNodes });
        if (postLink) {
            postLink(this.$componentScope, this.$element, attrs, requiredControllers, transcludeFn);
        }
        this.callLifecycleHook('$postLink', this.controllerInstance);
    };
    /**
     * @param {?} changes
     * @return {?}
     */
    UpgradeComponent.prototype.ngOnChanges = function (changes) {
        var _this = this;
        // Forward input changes to `bindingDestination`
        Object.keys(changes).forEach(function (propName) { return _this.bindingDestination[propName] = changes[propName].currentValue; });
        this.callLifecycleHook('$onChanges', this.bindingDestination, changes);
    };
    /**
     * @return {?}
     */
    UpgradeComponent.prototype.ngDoCheck = function () {
        var _this = this;
        var /** @type {?} */ twoWayBoundProperties = this.bindings.twoWayBoundProperties;
        var /** @type {?} */ twoWayBoundLastValues = this.bindings.twoWayBoundLastValues;
        var /** @type {?} */ propertyToOutputMap = this.bindings.propertyToOutputMap;
        twoWayBoundProperties.forEach(function (propName, idx) {
            var /** @type {?} */ newValue = _this.bindingDestination[propName];
            var /** @type {?} */ oldValue = twoWayBoundLastValues[idx];
            if (!looseIdentical(newValue, oldValue)) {
                var /** @type {?} */ outputName = propertyToOutputMap[propName];
                var /** @type {?} */ eventEmitter = ((_this))[outputName];
                eventEmitter.emit(newValue);
                twoWayBoundLastValues[idx] = newValue;
            }
        });
    };
    /**
     * @return {?}
     */
    UpgradeComponent.prototype.ngOnDestroy = function () {
        this.callLifecycleHook('$onDestroy', this.controllerInstance);
        this.$componentScope.$destroy();
    };
    /**
     * @param {?} method
     * @param {?} context
     * @param {?=} arg
     * @return {?}
     */
    UpgradeComponent.prototype.callLifecycleHook = function (method, context, arg) {
        if (context && typeof context[method] === 'function') {
            context[method](arg);
        }
    };
    /**
     * @param {?} name
     * @return {?}
     */
    UpgradeComponent.prototype.getDirective = function (name) {
        var /** @type {?} */ directives = this.$injector.get(name + 'Directive');
        if (directives.length > 1) {
            throw new Error('Only support single directive definition for: ' + this.name);
        }
        var /** @type {?} */ directive = directives[0];
        if (directive.replace)
            this.notSupported('replace');
        if (directive.terminal)
            this.notSupported('terminal');
        if (directive.compile)
            this.notSupported('compile');
        var /** @type {?} */ link = directive.link;
        // QUESTION: why not support link.post?
        if (typeof link == 'object') {
            if (((link)).post)
                this.notSupported('link.post');
        }
        return directive;
    };
    /**
     * @param {?} directive
     * @return {?}
     */
    UpgradeComponent.prototype.getDirectiveRequire = function (directive) {
        var /** @type {?} */ require = directive.require || (directive.controller && directive.name);
        if (isMap(require)) {
            Object.keys(require).forEach(function (key) {
                var /** @type {?} */ value = require[key];
                var /** @type {?} */ match = value.match(REQUIRE_PREFIX_RE);
                var /** @type {?} */ name = value.substring(match[0].length);
                if (!name) {
                    require[key] = match[0] + key;
                }
            });
        }
        return require;
    };
    /**
     * @param {?} directive
     * @return {?}
     */
    UpgradeComponent.prototype.initializeBindings = function (directive) {
        var _this = this;
        var /** @type {?} */ btcIsObject = typeof directive.bindToController === 'object';
        if (btcIsObject && Object.keys(directive.scope).length) {
            throw new Error("Binding definitions on scope and controller at the same time is not supported.");
        }
        var /** @type {?} */ context = (btcIsObject) ? directive.bindToController : directive.scope;
        var /** @type {?} */ bindings = new Bindings();
        if (typeof context == 'object') {
            Object.keys(context).forEach(function (propName) {
                var /** @type {?} */ definition = context[propName];
                var /** @type {?} */ bindingType = definition.charAt(0);
                // QUESTION: What about `=*`? Ignore? Throw? Support?
                switch (bindingType) {
                    case '@':
                    case '<':
                        // We don't need to do anything special. They will be defined as inputs on the
                        // upgraded component facade and the change propagation will be handled by
                        // `ngOnChanges()`.
                        break;
                    case '=':
                        bindings.twoWayBoundProperties.push(propName);
                        bindings.twoWayBoundLastValues.push(INITIAL_VALUE);
                        bindings.propertyToOutputMap[propName] = propName + 'Change';
                        break;
                    case '&':
                        bindings.expressionBoundProperties.push(propName);
                        bindings.propertyToOutputMap[propName] = propName;
                        break;
                    default:
                        var /** @type {?} */ json = JSON.stringify(context);
                        throw new Error("Unexpected mapping '" + bindingType + "' in '" + json + "' in '" + _this.name + "' directive.");
                }
            });
        }
        return bindings;
    };
    /**
     * @param {?} directive
     * @return {?}
     */
    UpgradeComponent.prototype.compileTemplate = function (directive) {
        if (this.directive.template !== undefined) {
            return this.compileHtml(getOrCall(this.directive.template));
        }
        else if (this.directive.templateUrl) {
            var /** @type {?} */ url = getOrCall(this.directive.templateUrl);
            var /** @type {?} */ html = (this.$templateCache.get(url));
            if (html !== undefined) {
                return this.compileHtml(html);
            }
            else {
                throw new Error('loading directive templates asynchronously is not supported');
            }
        }
        else {
            throw new Error("Directive '" + this.name + "' is not a component, it is missing template.");
        }
    };
    /**
     * @param {?} controllerType
     * @param {?} $scope
     * @param {?} $element
     * @param {?} controllerAs
     * @return {?}
     */
    UpgradeComponent.prototype.buildController = function (controllerType, $scope, $element, controllerAs) {
        // TODO: Document that we do not pre-assign bindings on the controller instance
        var /** @type {?} */ locals = { $scope: $scope, $element: $element };
        var /** @type {?} */ controller = this.$controller(controllerType, locals, null, controllerAs);
        $element.data(controllerKey(this.directive.name), controller);
        return controller;
    };
    /**
     * @param {?} directiveName
     * @param {?} $element
     * @param {?} require
     * @return {?}
     */
    UpgradeComponent.prototype.resolveRequire = function (directiveName, $element, require) {
        var _this = this;
        if (!require) {
            return null;
        }
        else if (Array.isArray(require)) {
            return require.map(function (req) { return _this.resolveRequire(directiveName, $element, req); });
        }
        else if (typeof require === 'object') {
            var /** @type {?} */ value_1 = {};
            Object.keys(require).forEach(function (key) { return value_1[key] = _this.resolveRequire(directiveName, $element, require[key]); });
            return value_1;
        }
        else if (typeof require === 'string') {
            var /** @type {?} */ match = require.match(REQUIRE_PREFIX_RE);
            var /** @type {?} */ inheritType = match[1] || match[3];
            var /** @type {?} */ name_1 = require.substring(match[0].length);
            var /** @type {?} */ isOptional = !!match[2];
            var /** @type {?} */ searchParents = !!inheritType;
            var /** @type {?} */ startOnParent = inheritType === '^^';
            var /** @type {?} */ ctrlKey = controllerKey(name_1);
            if (startOnParent) {
                $element = $element.parent();
            }
            var /** @type {?} */ value = searchParents ? $element.inheritedData(ctrlKey) : $element.data(ctrlKey);
            if (!value && !isOptional) {
                throw new Error("Unable to find required '" + require + "' in upgraded directive '" + directiveName + "'.");
            }
            return value;
        }
        else {
            throw new Error("Unrecognized require syntax on upgraded directive '" + directiveName + "': " + require);
        }
    };
    /**
     * @return {?}
     */
    UpgradeComponent.prototype.setupOutputs = function () {
        var _this = this;
        // Set up the outputs for `=` bindings
        this.bindings.twoWayBoundProperties.forEach(function (propName) {
            var /** @type {?} */ outputName = _this.bindings.propertyToOutputMap[propName];
            ((_this))[outputName] = new EventEmitter();
        });
        // Set up the outputs for `&` bindings
        this.bindings.expressionBoundProperties.forEach(function (propName) {
            var /** @type {?} */ outputName = _this.bindings.propertyToOutputMap[propName];
            var /** @type {?} */ emitter = ((_this))[outputName] = new EventEmitter();
            // QUESTION: Do we want the ng1 component to call the function with `<value>` or with
            //           `{$event: <value>}`. The former is closer to ng2, the latter to ng1.
            _this.bindingDestination[propName] = function (value) { return emitter.emit(value); };
        });
    };
    /**
     * @param {?} feature
     * @return {?}
     */
    UpgradeComponent.prototype.notSupported = function (feature) {
        throw new Error("Upgraded directive '" + this.name + "' contains unsupported feature: '" + feature + "'.");
    };
    /**
     * @param {?} html
     * @return {?}
     */
    UpgradeComponent.prototype.compileHtml = function (html) {
        var /** @type {?} */ div = document.createElement('div');
        div.innerHTML = html;
        return this.$compile(div.childNodes);
    };
    return UpgradeComponent;
}());
function UpgradeComponent_tsickle_Closure_declarations() {
    /** @type {?} */
    UpgradeComponent.prototype.$injector;
    /** @type {?} */
    UpgradeComponent.prototype.$compile;
    /** @type {?} */
    UpgradeComponent.prototype.$templateCache;
    /** @type {?} */
    UpgradeComponent.prototype.$httpBackend;
    /** @type {?} */
    UpgradeComponent.prototype.$controller;
    /** @type {?} */
    UpgradeComponent.prototype.element;
    /** @type {?} */
    UpgradeComponent.prototype.$element;
    /** @type {?} */
    UpgradeComponent.prototype.$componentScope;
    /** @type {?} */
    UpgradeComponent.prototype.directive;
    /** @type {?} */
    UpgradeComponent.prototype.bindings;
    /** @type {?} */
    UpgradeComponent.prototype.linkFn;
    /** @type {?} */
    UpgradeComponent.prototype.controllerInstance;
    /** @type {?} */
    UpgradeComponent.prototype.bindingDestination;
    /** @type {?} */
    UpgradeComponent.prototype.name;
    /** @type {?} */
    UpgradeComponent.prototype.elementRef;
    /** @type {?} */
    UpgradeComponent.prototype.injector;
}
/**
 * @param {?} property
 * @return {?}
 */
function getOrCall(property) {
    return typeof (property) === 'function' ? property() : property;
}
/**
 * @param {?} value
 * @return {?}
 */
function isMap(value) {
    return value && !Array.isArray(value) && typeof value === 'object';
}
//# sourceMappingURL=upgrade_component.js.map