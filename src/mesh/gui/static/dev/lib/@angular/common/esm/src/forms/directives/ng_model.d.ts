import { OnChanges, SimpleChange } from '@angular/core';
import { EventEmitter } from '../../../src/facade/async';
import { ControlValueAccessor } from './control_value_accessor';
import { NgControl } from './ng_control';
import { Control } from '../model';
import { ValidatorFn, AsyncValidatorFn } from './validators';
export declare const formControlBinding: any;
/**
 * Binds a domain model to a form control.
 *
 * ### Usage
 *
 * `ngModel` binds an existing domain model to a form control. For a
 * two-way binding, use `[(ngModel)]` to ensure the model updates in
 * both directions.
 *
 * ### Example ([live demo](http://plnkr.co/edit/R3UX5qDaUqFO2VYR0UzH?p=preview))
 *  ```typescript
 * @Component({
 *      selector: "search-comp",
 *      directives: [FORM_DIRECTIVES],
 *      template: `<input type='text' [(ngModel)]="searchQuery">`
 *      })
 * class SearchComp {
 *  searchQuery: string;
 * }
 *  ```
 */
export declare class NgModel extends NgControl implements OnChanges {
    private _validators;
    private _asyncValidators;
    /** @internal */
    _control: Control;
    /** @internal */
    _added: boolean;
    update: EventEmitter<{}>;
    model: any;
    viewModel: any;
    constructor(_validators: any[], _asyncValidators: any[], valueAccessors: ControlValueAccessor[]);
    ngOnChanges(changes: {
        [key: string]: SimpleChange;
    }): void;
    readonly control: Control;
    readonly path: string[];
    readonly validator: ValidatorFn;
    readonly asyncValidator: AsyncValidatorFn;
    viewToModelUpdate(newValue: any): void;
}
