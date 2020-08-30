import { Directive, ElementRef, Renderer2, Input, ViewContainerRef, NgZone, Inject, Optional } from '@angular/core';
import { MatTooltip, MAT_TOOLTIP_SCROLL_STRATEGY, MAT_TOOLTIP_DEFAULT_OPTIONS, MatTooltipDefaultOptions } from '@angular/material/tooltip';
import { Overlay, ScrollDispatcher } from '@angular/cdk/overlay';
import { Platform } from '@angular/cdk/platform';
import { AriaDescriber, FocusMonitor } from '@angular/cdk/a11y';
import { Directionality } from '@angular/cdk/bidi';
import { TranslateService } from '@ngx-translate/core';

/**
 * Makes a form field show red boders and text, as well as a tooltip, if there is a
 * validation error. For making it work, set a valid error msg (may be a var for translation)
 * using the directive like this: '[appFormFieldError]="Msg"'.
 */
@Directive({
  selector: '[appFormFieldError]',
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

  constructor(
    private translate: TranslateService,
    private renderer: Renderer2,
    private elementRef: ElementRef,
    overlay: Overlay,
    scrollDispatcher: ScrollDispatcher,
    viewContainerRef: ViewContainerRef,
    ngZone: NgZone,
    platform: Platform,
    ariaDescriber: AriaDescriber,
    focusMonitor: FocusMonitor,
    @Inject(MAT_TOOLTIP_SCROLL_STRATEGY) scrollStrategy: any,
    @Optional() dir: Directionality,
    @Optional() @Inject(MAT_TOOLTIP_DEFAULT_OPTIONS) defaultOptions: MatTooltipDefaultOptions,
  ) {
    super(
      overlay,
      elementRef,
      scrollDispatcher,
      viewContainerRef,
      ngZone,
      platform,
      ariaDescriber,
      focusMonitor,
      scrollStrategy,
      dir,
      defaultOptions,
    );

    this.tooltipClass = 'error-tooltip';
  }

  // Activates or disables the highlight effect.
  private updateField() {
    if (this.message) {
      this.renderer.addClass(this.elementRef.nativeElement, 'red-field');
    } else {
      this.renderer.removeClass(this.elementRef.nativeElement, 'red-field');
    }
  }
}
