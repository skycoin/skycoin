"use strict";
var core_1 = require('@angular/core');
var lang_1 = require('../../../src/facade/lang');
var control_value_accessor_1 = require('./control_value_accessor');
exports.DEFAULT_VALUE_ACCESSOR = 
/* @ts2dart_Provider */ {
    provide: control_value_accessor_1.NG_VALUE_ACCESSOR,
    useExisting: core_1.forwardRef(function () { return DefaultValueAccessor; }),
    multi: true
};
var DefaultValueAccessor = (function () {
    function DefaultValueAccessor(_renderer, _elementRef) {
        this._renderer = _renderer;
        this._elementRef = _elementRef;
        this.onChange = function (_) { };
        this.onTouched = function () { };
    }
    DefaultValueAccessor.prototype.writeValue = function (value) {
        var normalizedValue = lang_1.isBlank(value) ? '' : value;
        this._renderer.setElementProperty(this._elementRef.nativeElement, 'value', normalizedValue);
    };
    DefaultValueAccessor.prototype.registerOnChange = function (fn) { this.onChange = fn; };
    DefaultValueAccessor.prototype.registerOnTouched = function (fn) { this.onTouched = fn; };
    DefaultValueAccessor.decorators = [
        { type: core_1.Directive, args: [{
                    selector: 'input:not([type=checkbox])[ngControl],textarea[ngControl],input:not([type=checkbox])[ngFormControl],textarea[ngFormControl],input:not([type=checkbox])[ngModel],textarea[ngModel],[ngDefaultControl]',
                    // TODO: vsavkin replace the above selector with the one below it once
                    // https://github.com/angular/angular/issues/3011 is implemented
                    // selector: '[ngControl],[ngModel],[ngFormControl]',
                    host: { '(input)': 'onChange($event.target.value)', '(blur)': 'onTouched()' },
                    bindings: [exports.DEFAULT_VALUE_ACCESSOR]
                },] },
    ];
    DefaultValueAccessor.ctorParameters = [
        { type: core_1.Renderer, },
        { type: core_1.ElementRef, },
    ];
    return DefaultValueAccessor;
}());
exports.DefaultValueAccessor = DefaultValueAccessor;
//# sourceMappingURL=default_value_accessor.js.map