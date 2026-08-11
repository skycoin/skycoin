import { BigNumber } from 'bignumber.js';

export function convertAsciiToHexa(str: any): string {
  const arr1: string[] = [];
  for (let n = 0, l = str.length; n < l; n ++) {
    const hex = Number(str.charCodeAt(n)).toString(16);
    arr1.push(hex.length !== 1 ? hex : '0' + hex);
  }
  return arr1.join('');
}

/**
 * Builds a BigNumber, yielding NaN for anything it cannot parse.
 *
 * bignumber.js threw this behaviour away in v11: `new BigNumber(undefined)`
 * used to produce NaN and now raises. API responses legitimately omit numeric
 * fields — an unconfirmed transaction has no calculated_hours on its inputs —
 * so constructing directly turns a missing field into an exception that
 * escapes into the view.
 */
export function toBigNumber(value: any): BigNumber {
  try {
    return new BigNumber(value);
  } catch {
    return new BigNumber(NaN);
  }
}
