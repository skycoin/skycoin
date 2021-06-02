import { SendPage } from './send.po';

describe('Send', () => {
  const page = new SendPage();

  it('should display title', () => {
    page.navigateTo();
    expect<any>(page.getHeaderText()).toEqual('Send');
  });

  it('should have wallets', () => {
    expect<any>(page.getWalletsCount()).toBeGreaterThan(0);
  });

  it('should have coins in wallets', () => {
    expect<any>(page.getWalletsWithCoins().then(w => w.length)).toBeGreaterThan(0);
  });

  it('should have wallets enabled', () => {
    expect<any>(page.getValidWallets().then(w => w.length)).toBeGreaterThan(0);
  });

  it('should select valid wallet', () => {
    expect<any>(page.selectValidWallet()).toBeTruthy();
  });

  it('should not be able to send with wrong amount', () => {
    expect<any>(page.fillFormWithCoins('615701')).toBeFalsy();
    expect<any>(page.fillFormWithCoins('0')).toBeFalsy();
    expect<any>(page.fillFormWithCoins('a')).toBeFalsy();
  });

  it('should be able to send with correct amount', () => {
    expect<any>(page.fillFormWithCoins('615700')).toBeTruthy();
    expect<any>(page.fillFormWithCoins('1')).toBeTruthy();
  });
});
