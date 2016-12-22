"use strict";
var __extends = (this && this.__extends) || function (d, b) {
    for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p];
    function __() { this.constructor = d; }
    d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
};
var core_1 = require('@angular/core');
var async_1 = require('../../../src/facade/async');
var control_value_accessor_1 = require('./control_value_accessor');
var ng_control_1 = require('./ng_control');
var model_1 = require('../model');
var validators_1 = require('../validators');
var shared_1 = require('./shared');
exports.formControlBinding = 
/*@ts2dart_const*/ /* @ts2dart_Provider */ {
    provide: ng_control_1.NgControl,
    useExisting: core_1.forwardRef(function () { return NgModel; })
};
var NgModel = (function (_super) {
    __extends(NgModel, _super);
    function NgModel(_validators, _asyncValidators, valueAccessors) {
        _super.call(this);
        this._validators = _validators;
        this._asyncValidators = _asyncValidators;
        /** @internal */
        this._control = new model_1.Control();
        /** @internal */
        this._added = false;
        this.update = new async_1.EventEmitter();
        this.valueAccessor = shared_1.selectValueAccessor(this, valueAccessors);
    }
    NgModel.prototype.ngOnChanges = function (changes) {
        if (!this._added) {
            shared_1.setUpControl(this._control, this);
            this._control.updateValueAndValidity({ emitEvent: false });
            this._added = true;
        }
        if (shared_1.isPropertyUpdated(changes, this.viewModel)) {
            this._control.updateValue(this.model);
            this.viewModel = this.model;
        }
    };
    Object.defineProperty(NgModel.prototype, "control", {
        get: function () { return this._control; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NgModel.prototype, "path", {
        get: function () { return []; },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NgModel.prototype, "validator", {
        get: function () { return shared_1.composeValidators(this._validators); },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(NgModel.prototype, "asyncValidator", {
        get: function () { return shared_1.composeAsyncValidators(this._asyncValidators); },
        enumerable: true,
        configurable: true
    });
    NgModel.prototype.viewToModelUpdate = function (newValue) {
        this.viewModel = newValue;
        async_1.ObservableWrapper.callEmit(this.update, newValue);
    };
    NgModel.decorators = [
        { type: core_1.Directive, args: [{
                    selector: '[ngModel]:not([ngControl]):not([ngFormControl])',
                    bindings: [exports.formControlBinding],
                    inputs: ['model: ngModel'],
                    outputs: ['update: ngModelChange'],
                    exportAs: 'ngForm'
                },] },
    ];
    NgModel.ctorParameters = [
        { type: undefined, decorators: [{ type: core_1.Optional }, { type: core_1.Self }, { type: core_1.Inject, args: [validators_1.NG_VALIDATORS,] },] },
        { type: undefined, decorators: [{ type: core_1.Optional }, { type: core_1.Self }, { type: core_1.Inject, args: [validators_1.NG_ASYNC_VALIDATORS,] },] },
        { type: undefined, decorators: [{ type: core_1.Optional }, { type: core_1.Self }, { type: core_1.Inject, args: [control_value_accessor_1.NG_VALUE_ACCESSOR,] },] },
    ];
    return NgModel;
}(ng_control_1.NgControl));
exports.NgModel = NgModel;
//# sourceMappingURL=ng_model.js.map