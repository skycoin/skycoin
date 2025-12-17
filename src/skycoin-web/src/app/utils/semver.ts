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
