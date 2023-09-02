import { SendPage } from './send.po';

describe('Send', () => {
  const page = new SendPage();

  it('should display title', async () => {
    await page.navigateTo();
    expect<any>(await page.getHeaderText()).toEqual('Send');
  });

  it('should have wallets', async () => {
    expect<any>(await page.getWalletsCount()).toBeGreaterThan(0);
  });

  it('should have coins in wallets', async () => {
    expect<any>(await page.getWalletsWithCoins().then(w => w.length)).toBeGreaterThan(0);
  });

  it('should have wallets enabled', async () => {
    expect<any>(await page.getValidWallets().then(w => w.length)).toBeGreaterThan(0);
  });

  it('should select valid wallet', async () => {
    expect<any>(await page.selectValidWallet()).toBeTruthy();
  });

  it('should not be able to send with wrong amount', async () => {
    expect<any>(await page.fillFormWithCoins('615701')).toBeFalsy();
    expect<any>(await page.fillFormWithCoins('0')).toBeFalsy();
    expect<any>(await page.fillFormWithCoins('a')).toBeFalsy();
  });

  it('should be able to send with correct amount', async () => {
    expect<any>(await page.fillFormWithCoins('615700')).toBeTruthy();
    expect<any>(await page.fillFormWithCoins('1')).toBeTruthy();
  });
});
