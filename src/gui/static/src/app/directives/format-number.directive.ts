import { Directive, ElementRef, HostListener } from '@angular/core';
import BigNumber from 'bignumber.js';

/**
 * Makes an input field format a number, to make it easier to read. It also removed
 * invalid characters.
 */
@Directive({
  selector: '[appFormatNumber]',
})
export class FormatNumberDirective {
  /**
   * Last moment in which the backspace key was pressed.
   */
  private lastBackspaceDate = 0;

  constructor(
    private el: ElementRef,
  ) { }

  /**
   * Called when the value of the field changes.
   */
  @HostListener('ngModelChange') onChange() {
    let value = (this.el.nativeElement as HTMLInputElement).value;
    if (value) {
      // Get the current position of the caret.
      let cursorPos = (this.el.nativeElement as HTMLInputElement).selectionStart;

      // Return the caret as many positions as commas and invalid characters are found before it.
      let charactersRemovedBeforePos = 0;
      for (let i = 0; i < cursorPos; i++) {
        if (
          value.charAt(i) !== '0' &&
          value.charAt(i) !== '1' &&
          value.charAt(i) !== '2' &&
          value.charAt(i) !== '3' &&
          value.charAt(i) !== '4' &&
          value.charAt(i) !== '5' &&
          value.charAt(i) !== '6' &&
          value.charAt(i) !== '7' &&
          value.charAt(i) !== '8' &&
          value.charAt(i) !== '9' &&
          value.charAt(i) !== '.'
        ) {
          charactersRemovedBeforePos += 1;
        }
      }
      cursorPos -= charactersRemovedBeforePos;

      // Remove all invalid characters and commas from the text.
      value = value.replace(/[^0-9.]/gi, '');

      if (!(new BigNumber(value)).isNaN()) {
        // Separate de decimal part.
        const numberParts = value.split('.');
        if (numberParts.length > 0 && numberParts.length < 3) {
          // Add the commas.
          let leftPart = numberParts[0];
          for (let i = numberParts[0].length - 3; i > 0; i -= 3) {
            leftPart = leftPart.substr(0, i) + ',' + leftPart.substr(i);

            // Move the caret if the comma was added at the left.
            if (i <= cursorPos) {
              cursorPos += 1;
            }
          }

          // If the function was called because the user pressed the backspace key to delete a comma,
          // move the caret one space back, to allow it to pass the comma character..
          let spacesToReturn = 0;
          if (this.lastBackspaceDate && (new Date()).getTime() - this.lastBackspaceDate < 30 && leftPart.charAt(cursorPos - 1) === ',') {
            spacesToReturn += 1;
          }

          // Update the value.
          (this.el.nativeElement as HTMLInputElement).value = leftPart + (numberParts.length > 1 ? ('.' + numberParts[1]) : '');
          // Update the caret position.
          (this.el.nativeElement as HTMLInputElement).selectionStart = cursorPos - spacesToReturn;
          (this.el.nativeElement as HTMLInputElement).selectionEnd = cursorPos - spacesToReturn;
        }
      } else {
        // Just update the field with the filtered string.
        (this.el.nativeElement as HTMLInputElement).value = value;
      }
    }
  }

  /**
   * Called when the user presses a key.
   */
  @HostListener('keydown', ['$event']) p1(event: KeyboardEvent) {
    // If the user pressed the backspace key, save in which moment it happened.
    if (event.keyCode === 8) {
      this.lastBackspaceDate = (new Date()).getTime();
    } else {
      this.lastBackspaceDate = null;
    }
  }
}
