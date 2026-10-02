import { WalletPathPipe, walletPath } from './wallet-path';

describe('walletPath', () => {
  it('leaves paths unchanged at the root, where the standalone app mounts the wallet', () => {
    expect(walletPath('/', '/history')).toBe('/history');
    expect(walletPath('/', 'settings/node')).toBe('/settings/node');
  });

  it('puts paths under the base a hosting app provides', () => {
    expect(walletPath('/wallet', '/history')).toBe('/wallet/history');
    expect(walletPath('/wallet/', '/settings/outputs')).toBe('/wallet/settings/outputs');
  });

  it('is what the walletPath pipe returns', () => {
    expect(new WalletPathPipe('/wallet').transform('/send')).toBe('/wallet/send');
  });
});
