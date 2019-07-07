import { WalletsPage } from './wallets.po';

describe('Wallets', () => {
  const page = new WalletsPage();

  it('should display title', () => {
    page.navigateTo();
    expect<any>(page.getHeaderText()).toEqual('Wallets');
  });

  it('should show create wallet', () => {
    expect<any>(page.showAddWallet()).toEqual(true);
    expect<any>(page.getWalletModalTitle()).toEqual('Create Wallet');
  });

  it('should validate create wallet, seed mismatch', () => {
    expect<any>(page.fillWalletForm('Test', 'seed', 'seed2')).toEqual(false);
  });

  it('should validate create wallet, empty label', () => {
    expect<any>(page.fillWalletForm('', 'seed', 'seed', false)).toEqual(false);
  });

  it('should create wallet', () => {
    expect<any>(page.fillWalletForm('Test create wallet', 'test create wallet', 'test create wallet', false)).toEqual(true);
    page.waitForWalletToBeCreated();
  });

  it('should show load wallet', () => {
    expect<any>(page.showLoadWallet()).toEqual(true);
    expect<any>(page.getWalletModalTitle()).toEqual('Load Wallet');
  });

  it('should validate load wallet, seed', () => {
    expect<any>(page.fillWalletForm('Test', '', null)).toEqual(false);
  });

  it('should validate load wallet, empty label', () => {
    expect<any>(page.fillWalletForm('', 'seed', null, false)).toEqual(false);
  });

  it('should load wallet', () => {
    expect<any>(page.fillWalletForm('Test load wallet', 'test load wallet', null, false)).toEqual(true);
    page.waitForWalletToBeCreated();
  });

  it('should expand wallet', () => {
    expect<any>(page.expandWallet()).toEqual(true);
  });

  it('should show wallet QR modal', () => {
    expect<any>(page.showQrDialog()).toEqual(true);
  });

  it('should hide wallet QR modal', () => {
    expect<any>(page.hideQrDialog()).toEqual(false);
  });

  it('should add address to wallet', () => {
    expect<any>(page.addAddress()).toEqual(true);
  });

  it('should hide empty address', () => {
    expect<any>(page.getCountOfEmptyAddresses('.-hide-empty')).toEqual(0);
  });

  it('should show empty address', () => {
    expect<any>(page.getCountOfEmptyAddresses('.-show-empty')).toBeGreaterThan(0);
  });

  it('should show change wallet name modal', () => {
    expect<any>(page.showChangeWalletName()).toEqual(true);
  });

  it('should change wallet name', () => {
    expect<any>(page.changeWalletName()).toEqual(true);
  });

  it('should encrypt wallet', () => {
    expect<any>(page.canEncrypt()).toEqual(true);
  });

  it('should decrypt wallet', () => {
    expect<any>(page.canDecrypt()).toEqual(true);
  });

  it('should display price information', () => {
    expect<any>(page.showPriceInformation()).toEqual(true);
  });
});
