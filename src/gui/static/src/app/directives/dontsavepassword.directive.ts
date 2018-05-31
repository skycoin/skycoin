import { Directive, ElementRef, HostListener } from '@angular/core';

@Directive({
  selector: '[appDontSavePassword]',
})
export class DontsavepasswordDirective {
  constructor(
    private el: ElementRef,
  ) {
    el.nativeElement.autocomplete = 'new-password';
    el.nativeElement.readOnly = true;
  }

  @HostListener('focus') onFocus() {
    this.el.nativeElement.readOnly = false;
  }
}
