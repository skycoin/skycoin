import { WalletsPage } from './wallets.po';

describe('Wallets', () => {
  const page = new WalletsPage();

  it('should display title', async () => {
    await page.navigateTo();
    expect<any>(await page.getHeaderText()).toEqual('Wallets');
  });

  it('should show create wallet', async () => {
    expect<any>(await page.showAddWallet()).toEqual(true);
    expect<any>(await page.getWalletModalTitle()).toEqual('Create Wallet');
  });

  it('should validate create wallet, seed mismatch', async () => {
    expect<any>(await page.fillWalletForm('Test', 'seed', 'seed2')).toEqual(false);
  });

  it('should validate create wallet, empty label', async () => {
    expect<any>(await page.fillWalletForm('', 'seed', 'seed', false)).toEqual(false);
  });

  it('should create wallet', async () => {
    expect<any>(await page.fillWalletForm('Test create wallet', 'test create wallet', 'test create wallet', false)).toEqual(true);
    await page.waitForWalletToBeCreated();
  });

  it('should show load wallet', async () => {
    expect<any>(await page.showLoadWallet()).toEqual(true);
    expect<any>(await page.getWalletModalTitle()).toEqual('Load Wallet');
  });

  it('should validate load wallet, seed', async () => {
    expect<any>(await page.fillWalletForm('Test', '', null)).toEqual(false);
  });

  it('should validate load wallet, empty label', async () => {
    expect<any>(await page.fillWalletForm('', 'seed', null, false)).toEqual(false);
  });

  it('should load wallet', async () => {
    expect<any>(await page.fillWalletForm('Test load wallet', 'test load wallet', null, false)).toEqual(true);
    await page.waitForWalletToBeCreated();
  });

  it('should expand wallet', async () => {
    expect<any>(await page.expandWallet()).toEqual(true);
  });

  it('should show wallet QR modal', async () => {
    expect<any>(await page.showQrDialog()).toEqual(true);
  });

  it('should hide wallet QR modal', async () => {
    expect<any>(await page.hideQrDialog()).toEqual(false);
  });

  it('should add address to wallet', async () => {
    expect<any>(await page.addAddress()).toEqual(true);
  });

  it('should hide empty address', async () => {
    expect<any>(await page.getCountOfEmptyAddresses('.-hide-empty')).toEqual(0);
  });

  it('should show empty address', async () => {
    expect<any>(await page.getCountOfEmptyAddresses('.-show-empty')).toBeGreaterThan(0);
  });

  it('should show change wallet name modal', async () => {
    expect<any>(await page.showChangeWalletName()).toEqual(true);
  });

  it('should change wallet name', async () => {
    expect<any>(await page.changeWalletName()).toEqual(true);
  });

  it('should encrypt wallet', async () => {
    expect<any>(await page.canEncrypt()).toEqual(true);
  });

  it('should decrypt wallet', async () => {
    expect<any>(await page.canDecrypt()).toEqual(true);
  });
});
