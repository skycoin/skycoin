/* eslint-disable prefer-arrow/prefer-arrow-functions */

/**
 * This file contains general helper functions.
 */

/**
 * Indicates if an update is needed by comparing two version numbers.
 * @param from Current version number. Must be a 3 parts number using the SemVer format.
 * Each part must contain numbers only, but the last part may end with '-rc', in which
 * case this version will be considered older than the one on the 'to' param if both
 * versions differ only by the '-rc' part.
 * @param to Number of the lastest version. Must be a 3 parts number using the SemVer format.
 * Each part must contain numbers only.
 * @returns true if 'from' is smaller than 'to'.
 */
export function shouldUpgradeVersion(from: string, to: string): boolean {
  const toParts = to.split('.');
  const fromSplit = from.split('-');
  const fromParts = fromSplit[0].split('.');

  for (let i = 0; i < 3; i++) {
    const toNumber = Number(toParts[i]);
    const fromNumber = Number(fromParts[i]);

    if (toNumber > fromNumber) {
      return true;
    }

    if (toNumber < fromNumber) {
      return false;
    }
  }

  return fromSplit.length === 2 && fromSplit[1].indexOf('rc') !== -1;
}

/**
 * Copies a text to the clipboard.
 * @param text Text to be copied.
 */
export function copyTextToClipboard(text: string) {
  const selBox = document.createElement('textarea');

  selBox.style.position = 'fixed';
  selBox.style.left = '0';
  selBox.style.top = '0';
  selBox.style.opacity = '0';
  selBox.value = text;

  document.body.appendChild(selBox);
  selBox.focus();
  selBox.select();

  document.execCommand('copy');
  document.body.removeChild(selBox);
}

/**
 * Response returned by parseRequestLink. Only the values found in the link will have a value.
 */
export class RequestLinkParams {
  address: string;
  coins: string;
  hours: string;
  message: string;
}

/**
 * Gets the data of a link requesting coins. If the link is not valid, returns null. The
 * function does not validate the requested values, just the structure of the link.
 * @param linkText Link to check.
 */
export function parseRequestLink(linkText: string): RequestLinkParams {
  let parsed: URL;
  try {
    parsed = new URL(linkText);
  } catch (e) {}

  if (!parsed || !parsed.protocol || parsed.protocol.toLowerCase() !== 'skycoin:' || !parsed.pathname) {
    return null;
  }

  const response = new RequestLinkParams();

  response.address = parsed.pathname;
  response.coins = parsed.searchParams.get('amount');
  response.hours = parsed.searchParams.get('hours');
  response.message = parsed.searchParams.get('message');

  return response;
}

/**
 * Removes the commas from the string representation of a number.
 */
export function removeCommas(formattedNumber: string): string {
  if (formattedNumber) {
    const commaReg = new RegExp(',', 'g');

    return formattedNumber.replace(commaReg, '');
  }

  return formattedNumber;
}
