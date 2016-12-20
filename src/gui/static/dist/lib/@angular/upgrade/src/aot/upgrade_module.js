/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { Injector, NgModule, NgZone, Testability } from '@angular/core';
import * as angular from '../angular_js';
import { controllerKey } from '../util';
import { angular1Providers, setTempInjectorRef } from './angular1_providers';
import { $$TESTABILITY, $DELEGATE, $INJECTOR, $PROVIDE, INJECTOR_KEY, UPGRADE_MODULE_NAME } from './constants';
/**
 *  *
  * *Part of the [upgrade/static](/docs/ts/latest/api/#!?query=upgrade%2Fstatic)
  * library for hybrid upgrade apps that support AoT compilation*
  * *
  * Allows Angular 1 and Angular 2+ components to be used together inside a hybrid upgrade
  * application, which supports AoT compilation.
  * *
  * Specifically, the classes and functions in the `upgrade/static` module allow the following:
  * 1. Creation of an Angular 2+ directive that wraps and exposes an Angular 1 component so
  * that it can be used in an Angular 2 template. See {@link UpgradeComponent}.
  * 2. Creation of an Angular 1 directive that wraps and exposes an Angular 2+ component so
  * that it can be used in an Angular 1 template. See {@link downgradeComponent}.
  * 3. Creation of an Angular 2+ root injector provider that wraps and exposes an Angular 1
  * service so that it can be injected into an Angular 2+ context. See
  * {@link UpgradeModule#upgrading-an-angular-1-service Upgrading an Angular 1 service} below.
  * 4. Creation of an Angular 1 service that wraps and exposes an Angular 2+ injectable
  * so that it can be injected into an Angular 1 context. See {@link downgradeInjectable}.
  * 3. Bootstrapping of a hybrid Angular application which contains both of the frameworks
  * coexisting in a single application. See the
  * {@link UpgradeModule#example example} below.
  * *
  * ## Mental Model
  * *
  * When reasoning about how a hybrid application works it is useful to have a mental model which
  * describes what is happening and explains what is happening at the lowest level.
  * *
  * 1. There are two independent frameworks running in a single application, each framework treats
  * the other as a black box.
  * 2. Each DOM element on the page is owned exactly by one framework. Whichever framework
  * instantiated the element is the owner. Each framework only updates/interacts with its own
  * DOM elements and ignores others.
  * 3. Angular 1 directives always execute inside the Angular 1 framework codebase regardless of
  * where they are instantiated.
  * 4. Angular 2+ components always execute inside the Angular 2+ framework codebase regardless of
  * where they are instantiated.
  * 5. An Angular 1 component can be "upgraded"" to an Angular 2+ component. This is achieved by
  * defining an Angular 2+ directive, which bootstraps the Angular 1 component at its location
  * in the DOM. See {@link UpgradeComponent}.
  * 6. An Angular 2+ component can be "downgraded"" to an Angular 1 component. This is achieved by
  * defining an Angular 1 directive, which bootstraps the Angular 2+ component at its location
  * in the DOM. See {@link downgradeComponent}.
  * 7. Whenever an "upgraded"/"downgraded" component is instantiated the host element is owned by
  * the framework doing the instantiation. The other framework then instantiates and owns the
  * view for that component.
  * a. This implies that the component bindings will always follow the semantics of the
  * instantiation framework.
  * b. The DOM attributes are parsed by the framework that owns the current template. So
  * attributes
  * in Angular 1 templates must use kebab-case, while Angular 1 templates must use camelCase.
  * c. However the template binding syntax will always use the Angular 2+ style, e.g. square
  * brackets (`[...]`) for property binding.
  * 8. Angular 1 is always bootstrapped first and owns the root component.
  * 9. The new application is running in an Angular 2+ zone, and therefore it no longer needs calls
  * to
  * `$apply()`.
  * *
  * *
  * `import {UpgradeModule} from '@angular/upgrade/static';`
  * *
  * ## Example
  * Import the {@link UpgradeModule} into your top level {@link NgModule Angular 2+ `NgModule`}.
  * *
  * {@example upgrade/static/ts/module.ts region='ng2-module'}
  * *
  * Then bootstrap the hybrid upgrade app's module, get hold of the {@link UpgradeModule} instance
  * and use it to bootstrap the top level [Angular 1
  * module](https://docs.angularjs.org/api/ng/type/angular.Module).
  * *
  * {@example upgrade/static/ts/module.ts region='bootstrap'}
  * *
  * *
  * ## Upgrading an Angular 1 service
  * *
  * There is no specific API for upgrading an Angular 1 service. Instead you should just follow the
  * following recipe:
  * *
  * Let's say you have an Angular 1 service:
  * *
  * {@example upgrade/static/ts/module.ts region="ng1-title-case-service"}
  * *
  * Then you should define an Angular 2+ provider to be included in your {@link NgModule} `providers`
  * property.
  * *
  * {@example upgrade/static/ts/module.ts region="upgrade-ng1-service"}
  * *
  * Then you can use the "upgraded" Angular 1 service by injecting it into an Angular 2 component
  * or service.
  * *
  * {@example upgrade/static/ts/module.ts region="use-ng1-upgraded-service"}
  * *
  * *
  * This class is an `NgModule`, which you import to provide Angular 1 core services,
  * and has an instance method used to bootstrap the hybrid upgrade application.
  * *
  * ## Core Angular 1 services
  * Importing this {@link NgModule} will add providers for the core
  * [Angular 1 services](https://docs.angularjs.org/api/ng/service) to the root injector.
  * *
  * ## Bootstrap
  * The runtime instance of this class contains a {@link UpgradeModule#bootstrap `bootstrap()`}
  * method, which you use to bootstrap the top level Angular 1 module onto an element in the
  * DOM for the hybrid upgrade app.
  * *
  * It also contains properties to access the {@link UpgradeModule#injector root injector}, the
  * bootstrap {@link NgZone} and the
  * [Angular 1 $injector](https://docs.angularjs.org/api/auto/service/$injector).
  * *
 */
export var UpgradeModule = (function () {
    /**
     * @param {?} injector
     * @param {?} ngZone
     */
    function UpgradeModule(injector, ngZone) {
        this.injector = injector;
        this.ngZone = ngZone;
    }
    /**
     *  Bootstrap an Angular 1 application from this NgModule
     * @param {?} element the element on which to bootstrap the Angular 1 application
     * @param {?=} modules
     * @param {?=} config
     * @return {?}
     */
    UpgradeModule.prototype.bootstrap = function (element, modules, config /*angular.IAngularBootstrapConfig*/) {
        var _this = this;
        if (modules === void 0) { modules = []; }
        // Create an ng1 module to bootstrap
        var /** @type {?} */ upgradeModule = angular
            .module(UPGRADE_MODULE_NAME, modules)
            .value(INJECTOR_KEY, this.injector)
            .config([
            $PROVIDE, $INJECTOR,
            function ($provide, $injector) {
                if ($injector.has($$TESTABILITY)) {
                    $provide.decorator($$TESTABILITY, [
                        $DELEGATE,
                        function (testabilityDelegate) {
                            var /** @type {?} */ originalWhenStable = testabilityDelegate.whenStable;
                            var /** @type {?} */ injector = _this.injector;
                            // Cannot use arrow function below because we need the context
                            var /** @type {?} */ newWhenStable = function (callback) {
                                originalWhenStable.call(this, function () {
                                    var /** @type {?} */ ng2Testability = injector.get(Testability);
                                    if (ng2Testability.isStable()) {
                                        callback.apply(this, arguments);
                                    }
                                    else {
                                        ng2Testability.whenStable(newWhenStable.bind(this, callback));
                                    }
                                });
                            };
                            testabilityDelegate.whenStable = newWhenStable;
                            return testabilityDelegate;
                        }
                    ]);
                }
            }
        ])
            .run([
            $INJECTOR,
            function ($injector) {
                _this.$injector = $injector;
                // Initialize the ng1 $injector provider
                setTempInjectorRef($injector);
                _this.injector.get($INJECTOR);
                // Put the injector on the DOM, so that it can be "required"
                angular.element(element).data(controllerKey(INJECTOR_KEY), _this.injector);
                // Wire up the ng1 rootScope to run a digest cycle whenever the zone settles
                var /** @type {?} */ $rootScope = $injector.get('$rootScope');
                _this.ngZone.onMicrotaskEmpty.subscribe(function () { return _this.ngZone.runOutsideAngular(function () { return $rootScope.$evalAsync(); }); });
            }
        ]);
        // Make sure resumeBootstrap() only exists if the current bootstrap is deferred
        var /** @type {?} */ windowAngular = ((window) /** TODO #???? */)['angular'];
        windowAngular.resumeBootstrap = undefined;
        // Bootstrap the angular 1 application inside our zone
        this.ngZone.run(function () { angular.bootstrap(element, [upgradeModule.name], config); });
        // Patch resumeBootstrap() to run inside the ngZone
        if (windowAngular.resumeBootstrap) {
            var /** @type {?} */ originalResumeBootstrap_1 = windowAngular.resumeBootstrap;
            var /** @type {?} */ ngZone_1 = this.ngZone;
            windowAngular.resumeBootstrap = function () {
                var _this = this;
                var /** @type {?} */ args = arguments;
                windowAngular.resumeBootstrap = originalResumeBootstrap_1;
                ngZone_1.run(function () { windowAngular.resumeBootstrap.apply(_this, args); });
            };
        }
    };
    UpgradeModule.decorators = [
        { type: NgModule, args: [{ providers: angular1Providers },] },
    ];
    /** @nocollapse */
    UpgradeModule.ctorParameters = function () { return [
        { type: Injector, },
        { type: NgZone, },
    ]; };
    return UpgradeModule;
}());
function UpgradeModule_tsickle_Closure_declarations() {
    /** @type {?} */
    UpgradeModule.decorators;
    /**
     * @nocollapse
     * @type {?}
     */
    UpgradeModule.ctorParameters;
    /**
     * The Angular 1 `$injector` for the upgrade application.
     * @type {?}
     */
    UpgradeModule.prototype.$injector;
    /**
     * The root {@link Injector} for the upgrade application.
     * @type {?}
     */
    UpgradeModule.prototype.injector;
    /**
     * The bootstrap zone for the upgrade application
     * @type {?}
     */
    UpgradeModule.prototype.ngZone;
}
//# sourceMappingURL=upgrade_module.js.map