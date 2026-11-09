/**
 * Return true is current is equal or superior to target. Versions normally are composed of 3 parts,
 * separated by points, but this function only takes into account the first two parts, so 0.24.0 is
 * considered equal to 0.24.1.
 */
export function isEqualOrSuperiorVersion(current: string, target: string): boolean {
  const currentParts = current.split('.');
  if (currentParts.length < 2) { return false; }

  const targetParts = target.split('.');
  if (targetParts.length < 2) { return false; }

  for (let i = 0; i < 2; i++) {
    const currentNumber = Number(currentParts[i]);
    const targetNumber = Number(targetParts[i]);

    if (currentNumber > targetNumber) {
      return true;
    }

    if (currentNumber < targetNumber) {
      return false;
    }
  }

  return true;
}

/**
 * True when version to is newer than from, comparing all 3 parts. A release
 * candidate such as 0.28.7-rc1 is older than 0.28.7.
 */
export function shouldUpgradeVersion(from: string, to: string): boolean {
  if (!from || !to) {
    return false;
  }
  const toParts = to.split('-')[0].split('.');
  const fromSplit = from.split('-');
  const fromParts = fromSplit[0].split('.');

  for (let i = 0; i < 3; i++) {
    const toNumber = Number(toParts[i] || 0);
    const fromNumber = Number(fromParts[i] || 0);
    if (isNaN(toNumber) || isNaN(fromNumber)) {
      return false;
    }
    if (toNumber > fromNumber) {
      return true;
    }
    if (toNumber < fromNumber) {
      return false;
    }
  }

  return fromSplit.length === 2 && fromSplit[1].indexOf('rc') !== -1;
}
