import { OnboardingCreatePage } from './onboarding.po';

describe('Onboarding', () => {
  const page = new OnboardingCreatePage();

  it('should display title', async () => {
    await page.navigateTo();
    expect<any>(await page.getHeaderText()).toEqual('Create Wallet');
  });

  it('should select language', async () => {
    expect<any>(await page.selectLanguage()).toEqual(true);
  });

  it('should load wallet', async () => {
    expect<any>(await page.loadWallet()).toEqual(true);
  });

  it('should create wallet', async () => {
    expect<any>(await page.createWallet()).toEqual(true);
  });

  it('should show safeguard', async () => {
    expect<any>(await page.getSafeguardIsShown()).toEqual(true);
  });

  it('should hide accepted safeguard', async () => {
    expect<any>(await page.acceptSafeguard()).toEqual(false);
  });

  it('should be able to go back from wallet encryption', async () => {
    expect<any>(await page.goBack()).toEqual('Create Wallet');
    await page.createWallet(false);
    await page.acceptSafeguard();
  });

  it('should encrypt wallet by default', async () => {
    expect<any>(await page.getEncryptWalletCheckbox().isSelected()).toBeTruthy();
  });

  it('should be able to continue without encryption', async () => {
    expect<any>(await page.canContinueWithoutEncryption()).toEqual(true);
  });

  it('should encrypt wallet', async () => {
    expect<any>(await page.encryptWallet()).toEqual(true);
  });
});
