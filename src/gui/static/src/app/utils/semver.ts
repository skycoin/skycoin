/**
 * Compares two SemVer versions, returns true if 'from' is smaller than 'to'.
 * Special cases with 'rc' suffix are described in spec file.
 *
 * @returns {boolean}
 * @param from
 * @param to
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
