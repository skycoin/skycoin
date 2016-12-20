/**
 * @license
 * Copyright Google Inc. All Rights Reserved.
 *
 * Use of this source code is governed by an MIT-style license that can be
 * found in the LICENSE file at https://angular.io/license
 */
import { ChangeDetectorRef, ComponentFactory, ComponentRef, Injector, SimpleChanges } from '@angular/core';
import * as angular from '../angular_js';
import { ComponentInfo } from './component_info';
export declare class DowngradeComponentAdapter {
    private id;
    private info;
    private element;
    private attrs;
    private scope;
    private parentInjector;
    private parse;
    private componentFactory;
    component: any;
    inputs: Attr;
    inputChangeCount: number;
    inputChanges: SimpleChanges;
    componentRef: ComponentRef<any>;
    changeDetector: ChangeDetectorRef;
    componentScope: angular.IScope;
    childNodes: Node[];
    contentInsertionPoint: Node;
    constructor(id: string, info: ComponentInfo, element: angular.IAugmentedJQuery, attrs: angular.IAttributes, scope: angular.IScope, parentInjector: Injector, parse: angular.IParseService, componentFactory: ComponentFactory<any>);
    createComponent(): void;
    setupInputs(): void;
    projectContent(): void;
    setupOutputs(): void;
    registerCleanup(): void;
}
