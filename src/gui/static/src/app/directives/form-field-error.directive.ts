import { Directive, ElementRef, Renderer2, Input } from '@angular/core';
import { MatTooltip } from '@angular/material/tooltip';
import { TranslateService } from '@ngx-translate/core';


/**
 * Makes a form field show red boders and text, as well as a tooltip, if there is a
 * validation error. For making it work, set a valid error msg (may be a var for translation)
 * using the directive like this: '[appFormFieldError]="Msg"'.
 */
@Directive({
    selector: '[appFormFieldError]',
    standalone: false
})
export class FormFieldErrorDirective extends MatTooltip {
  // Error msg.
  @Input() set appFormFieldError(val: string) {
    if (val) {
      this.message = this.translate.instant(val);
    } else {
      this.message = null;
    }

    this.updateField();
  }

  // As of Angular Material 22 MatTooltip resolves its own dependencies with
  // inject() and its constructor takes no arguments, so we only inject what this
  // directive itself needs and call super() with no positional dependencies.
  constructor(
    private translate: TranslateService,
    private renderer: Renderer2,
    private elementRef: ElementRef,
  ) {
    super();

    this.tooltipClass = 'error-tooltip';
  }

  // Activates or disables the highlight effect.
  private updateField() {
    if (this.message) {
      this.renderer.addClass(this.elementRef.nativeElement, 'red-field');

      // If the fied is focussed, show the tooltip, enven if the mouse cursor is not over the field.
      if (document.activeElement === this.elementRef.nativeElement) {
        this.show();
      }
    } else {
      this.renderer.removeClass(this.elementRef.nativeElement, 'red-field');
    }
  }
}
