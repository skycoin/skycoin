import { Directive, HostListener } from '@angular/core';

// --- Allows an input to accept only numbers. Currently, the paste option is failing, it needs some more testing 2.222
@Directive({
  selector: 'input[appNumberField]'
})
export class NumberFieldDirective {
  @HostListener('keydown', ['$event']) onKeyDown(event) {
    const e = <KeyboardEvent> event;
    if ([46, 8, 9, 27, 13, 110, 190].indexOf(e.keyCode) !== -1 ||
      // Allow: Ctrl+A
      (e.keyCode === 65 && (e.ctrlKey || e.metaKey)) ||
      // Allow: Ctrl+C
      (e.keyCode === 67 && (e.ctrlKey || e.metaKey)) ||
      // Allow: Ctrl+V
      (e.keyCode === 86 && (e.ctrlKey || e.metaKey)) ||
      // Allow: Ctrl+X
      (e.keyCode === 88 && (e.ctrlKey || e.metaKey)) ||
      // Allow: Shift+insert
      (e.keyCode === 45 && (e.shiftKey || e.metaKey)) ||
      // Allow: home, end, left, right
      (e.keyCode >= 35 && e.keyCode <= 39)) {
      // let it happen, don't do anything
      return;
    }
    if ((e.shiftKey || (e.keyCode < 48 || e.keyCode > 57)) && (e.keyCode < 96 || e.keyCode > 105)) {
      e.preventDefault();
    }
  }

  @HostListener('paste', ['$event']) blockPaste(event) {
    const pastedValue = event.clipboardData
      ? event.clipboardData.getData('text/plain')
      : window['clipboardData'].getData('text');

    if (!pastedValue || !pastedValue.match(/^\d+(\.\d+)*$/)) {
      event.preventDefault();
    }
  }
}
