import { OnboardingCreatePage } from './onboarding.po';

describe('Onboarding', () => {
  const page = new OnboardingCreatePage();

  it('should display title', () => {
    page.navigateTo();
    expect<any>(page.getHeaderText()).toEqual('Create Wallet');
  });

  it('should select language', () => {
    expect<any>(page.selectLanguage()).toEqual(true);
  });

  it('should load wallet', () => {
    expect<any>(page.loadWallet()).toEqual(true);
  });

  it('should create wallet', () => {
    expect<any>(page.createWallet()).toEqual(true);
  });

  it('should show safeguard', () => {
    expect<any>(page.getSafeguardIsShown()).toEqual(true);
  });

  it('should hide accepted safeguard', () => {
    expect<any>(page.acceptSafeguard()).toEqual(false);
  });

  it('should be able to go back from wallet encryption', () => {
    expect<any>(page.goBack()).toEqual('Create Wallet');
    page.createWallet(false);
    page.acceptSafeguard();
  });

  it('should encrypt wallet by default', () => {
    expect<any>(page.getEncryptWalletCheckbox().isSelected()).toBeTruthy();
  });

  it('should be able to continue without encryption', () => {
    expect<any>(page.canContinueWithoutEncryption()).toEqual(true);
  });

  it('should encrypt wallet', () => {
    expect<any>(page.encryptWallet()).toEqual(true);
  });
});
