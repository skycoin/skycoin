/**
 * Compares two SemVer versions, returns true if 'from' is smaller than 'to'.
 * Special cases with 'rc' suffix are described in spec file.
 *
 * @returns {boolean}
 * @param from
 * @param to
 */
export function shouldUpgradeVersion(from: string, to: string): boolean {
  const fromParts = to.split('.');
  const toSplit = from.split('-');
  const toParts = toSplit[0].split('.');

  for (let i = 0; i < 3; i++) {
    const fromNumber = Number(fromParts[i]);
    const toNumber = Number(toParts[i]);

    if (fromNumber > toNumber) {
      return true;
    }

    if (fromNumber < toNumber) {
      return false;
    }
  }

  return toSplit.length === 2 && toSplit[1].indexOf('rc') !== -1;
}
