import { isEqualOrSuperiorVersion, shouldUpgradeVersion } from './semver';

describe('semver', () => {
  it('correctly compares versions', () => {
    expect(isEqualOrSuperiorVersion('0.23.0', '0.22.0')).toBeTruthy();
    expect(isEqualOrSuperiorVersion('0.23.0', '0.23.0')).toBeTruthy();
    expect(isEqualOrSuperiorVersion('0.23.0', '0.23.1')).toBeTruthy();
    expect(isEqualOrSuperiorVersion('0.23.1', '0.24.0')).toBeFalsy();
    expect(isEqualOrSuperiorVersion('0.24.0', '1.0.0')).toBeFalsy();
    expect(isEqualOrSuperiorVersion('0.24', '0.23')).toBeTruthy();
    expect(isEqualOrSuperiorVersion('0.24', '0.24')).toBeTruthy();
    expect(isEqualOrSuperiorVersion('0.24', '0.25')).toBeFalsy();
  });

  it('tells when a newer release is out', () => {
    expect(shouldUpgradeVersion('0.28.6', '0.28.7')).toBeTrue();
    expect(shouldUpgradeVersion('0.27.9', '0.28.0')).toBeTrue();
    expect(shouldUpgradeVersion('0.28.7-rc1', '0.28.7')).toBeTrue();
    expect(shouldUpgradeVersion('0.28.7', '0.28.7')).toBeFalse();
    expect(shouldUpgradeVersion('0.28.8', '0.28.7')).toBeFalse();
    expect(shouldUpgradeVersion('0.29.0-dev', '0.28.7')).toBeFalse();
    expect(shouldUpgradeVersion('', '0.28.7')).toBeFalse();
    expect(shouldUpgradeVersion('0.28.6', 'garbage')).toBeFalse();
  });
});
