import { Directive, ElementRef, Renderer, forwardRef } from '@angular/core';
import { isBlank } from '../../../src/facade/lang';
import { NG_VALUE_ACCESSOR } from './control_value_accessor';
export const DEFAULT_VALUE_ACCESSOR = 
/* @ts2dart_Provider */ {
    provide: NG_VALUE_ACCESSOR,
    useExisting: forwardRef(() => DefaultValueAccessor),
    multi: true
};
export class DefaultValueAccessor {
    constructor(_renderer, _elementRef) {
        this._renderer = _renderer;
        this._elementRef = _elementRef;
        this.onChange = (_) => { };
        this.onTouched = () => { };
    }
    writeValue(value) {
        var normalizedValue = isBlank(value) ? '' : value;
        this._renderer.setElementProperty(this._elementRef.nativeElement, 'value', normalizedValue);
    }
    registerOnChange(fn) { this.onChange = fn; }
    registerOnTouched(fn) { this.onTouched = fn; }
}
DefaultValueAccessor.decorators = [
    { type: Directive, args: [{
                selector: 'input:not([type=checkbox])[ngControl],textarea[ngControl],input:not([type=checkbox])[ngFormControl],textarea[ngFormControl],input:not([type=checkbox])[ngModel],textarea[ngModel],[ngDefaultControl]',
                // TODO: vsavkin replace the above selector with the one below it once
                // https://github.com/angular/angular/issues/3011 is implemented
                // selector: '[ngControl],[ngModel],[ngFormControl]',
                host: { '(input)': 'onChange($event.target.value)', '(blur)': 'onTouched()' },
                bindings: [DEFAULT_VALUE_ACCESSOR]
            },] },
];
DefaultValueAccessor.ctorParameters = [
    { type: Renderer, },
    { type: ElementRef, },
];
//# sourceMappingURL=default_value_accessor.js.map