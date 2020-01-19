import { shouldUpgradeVersion } from './semver';

describe('semver', () => {
  it('correctly compares versions', () => {
    expect(shouldUpgradeVersion('0.23.0', '0.22.0')).toBeFalsy();
    expect(shouldUpgradeVersion('0.23.0', '0.23.0')).toBeFalsy();
    expect(shouldUpgradeVersion('0.23.0', '0.23.1')).toBeTruthy();
    expect(shouldUpgradeVersion('0.23.1', '0.24.0')).toBeTruthy();
    expect(shouldUpgradeVersion('0.24.0', '1.0.0')).toBeTruthy();
  });

  it('correctly handles rc versions', () => {
    expect(shouldUpgradeVersion('0.23.1-rc.1', '0.23.0')).toBeFalsy();
    expect(shouldUpgradeVersion('0.23.1-rc.1', '0.23.1')).toBeTruthy();
    expect(shouldUpgradeVersion('0.23.1-rc.1', '0.23.2')).toBeTruthy();
  });
});
