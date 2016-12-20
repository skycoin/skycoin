/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { ComponentFactoryResolver } from '@angular/core';
import { $INJECTOR, $PARSE, INJECTOR_KEY } from './constants';
import { DowngradeComponentAdapter } from './downgrade_component_adapter';
var /** @type {?} */ downgradeCount = 0;
/**
 *  *
  * *Part of the [upgrade/static](/docs/ts/latest/api/#!?query=upgrade%2Fstatic)
  * library for hybrid upgrade apps that support AoT compilation*
  * *
  * Allows an Angular 2+ component to be used from Angular 1.
  * *
  * *
  * Let's assume that you have an Angular 2+ component called `ng2Heroes` that needs
  * to be made available in Angular 1 templates.
  * *
  * {@example upgrade/static/ts/module.ts region="ng2-heroes"}
  * *
  * We must create an Angular 1 [directive](https://docs.angularjs.org/guide/directive)
  * that will make this Angular 2+ component available inside Angular 1 templates.
  * The `downgradeComponent()` function returns a factory function that we
  * can use to define the Angular 1 directive that wraps the "downgraded" component.
  * *
  * {@example upgrade/static/ts/module.ts region="ng2-heroes-wrapper"}
  * *
  * In this example you can see that we must provide information about the component being
  * "downgraded". This is because once the AoT compiler has run, all metadata about the
  * component has been removed from the code, and so cannot be inferred.
  * *
  * We must do the following:
  * * specify the Angular 2+ component class that is to be downgraded
  * * specify all inputs and outputs that the Angular 1 component expects
  * *
  * *
  * A helper function that returns a factory function to be used for registering an
  * Angular 1 wrapper directive for "downgrading" an Angular 2+ component.
  * *
  * The parameter contains information about the Component that is being downgraded:
  * *
  * * `component: Type<any>`: The type of the Component that will be downgraded
  * * `inputs: string[]`: A collection of strings that specify what inputs the component accepts.
  * * `outputs: string[]`: A collection of strings that specify what outputs the component emits.
  * *
  * The `inputs` and `outputs` are strings that map the names of properties to camelCased
  * attribute names. They are of the form `"prop: attr"`; or simply `"propAndAttr" where the
  * property and attribute have the same identifier.
  * *
 * @param {?} info
 * @return {?}
 */
export function downgradeComponent(info) {
    var /** @type {?} */ idPrefix = "NG2_UPGRADE_" + downgradeCount++ + "_";
    var /** @type {?} */ idCount = 0;
    var /** @type {?} */ directiveFactory = function ($injector, $parse) {
        return {
            restrict: 'E',
            require: '?^' + INJECTOR_KEY,
            link: function (scope, element, attrs, parentInjector, transclude) {
                if (parentInjector === null) {
                    parentInjector = $injector.get(INJECTOR_KEY);
                }
                var /** @type {?} */ componentFactoryResolver = parentInjector.get(ComponentFactoryResolver);
                var /** @type {?} */ componentFactory = componentFactoryResolver.resolveComponentFactory(info.component);
                if (!componentFactory) {
                    throw new Error('Expecting ComponentFactory for: ' + info.component);
                }
                var /** @type {?} */ facade = new DowngradeComponentAdapter(idPrefix + (idCount++), info, element, attrs, scope, parentInjector, $parse, componentFactory);
                facade.setupInputs();
                facade.createComponent();
                facade.projectContent();
                facade.setupOutputs();
                facade.registerCleanup();
            }
        };
    };
    directiveFactory.$inject = [$INJECTOR, $PARSE];
    return directiveFactory;
}
//# sourceMappingURL=downgrade_component.js.map