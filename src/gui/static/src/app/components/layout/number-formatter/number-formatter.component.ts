import { Component, Input } from '@angular/core';
import BigNumber from 'bignumber.js';

/**
 * Formats a number, to show the decimal part in a differen style.
 */
@Component({
  selector: 'app-number-formatter',
  templateUrl: 'number-formatter.component.html',
  styleUrls: ['number-formatter.component.scss'],
})
export class NumberFormatterComponent {
  @Input() set number(val: string | number) {
    // Remove all invalid characters and commas from the text.
    val = (val + '').replace(/[^0-9.]/gi, '');

    const number = new BigNumber(val);

    this.integerPart = '0';
    this.decimalPart = null;

    if (number.isNaN()) {
      return;
    }

    // Process the number and separate it in parts, if needed.
    if (number.isInteger()) {
      this.integerPart = number.integerValue(BigNumber.ROUND_DOWN).toFormat();

      return;
    } else {
      const integerNumber = number.integerValue(BigNumber.ROUND_DOWN);

      this.integerPart = integerNumber.toFormat();
      this.decimalPart = number.minus(integerNumber).toFormat().substr(2);
    }
  }

  /**
   * If true, the alternative style for the header will be used.
   */
  @Input() altDecimalsStyle = false;

  // Integer part to show.
  integerPart = '0';
  // Decimal part to show.
  decimalPart: string = null;
}
