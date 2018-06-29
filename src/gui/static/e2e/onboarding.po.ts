import { browser, by, element } from 'protractor';

export class OnboardingCreatePage {
  navigateTo() {
    return browser.get('/#/wizard');
  }

  getHeaderText() {
    return element(by.css('.-header span')).getText();
  }

  getSafeguardIsShown() {
    return element(by.css('app-onboarding-safeguard')).isPresent();
  }

  acceptSafeguard() {
   return element.all(by.css('.mat-checkbox-label')).first().click().then(() => {
     return element(by.buttonText('Continue')).click().then(() => {
        return this.getSafeguardIsShown();
      });
    });
  }

  createWallet() {
    element(by.buttonText('New')).click();

    const label = element(by.css('[formcontrolname="label"]'));
    const seed = element(by.css('[formcontrolname="seed"]'));
    const confirm = element(by.css('[formcontrolname="confirm_seed"]'));
    const btnCreate = element(by.buttonText('Create'));

    label.clear();
    label.sendKeys('Test onboarding wallet');
    seed.clear();
    seed.sendKeys('test test');
    confirm.clear();
    confirm.sendKeys('test test');

    return btnCreate.isEnabled().then(status => {
      if (status) {
        btnCreate.click();
      }

      return status;
    });
  }

  loadWallet() {
    element(by.buttonText('Load')).click();

    const label = element(by.css('[formcontrolname="label"]'));
    const seed = element(by.css('[formcontrolname="seed"]'));
    const btnLoad = element(by.buttonText('Create'));

    label.clear();
    label.sendKeys('Test wallet');
    seed.clear();
    seed.sendKeys('test test');

    return btnLoad.isEnabled();
  }

  goBack() {
    return element(by.buttonText('Back')).click().then(() => {
      return this.getHeaderText();
    });
  }

  getEncryptWalletCheckbox() {
    return element(by.css('.mat-checkbox-input'));
  }

  canContinueWithoutEncryption() {
    return element(by.css('.mat-checkbox-label')).click().then(() => {
      return element(by.buttonText('Finish')).isEnabled();
    });
  }

  encryptWallet() {
    const password = element(by.css('[formcontrolname="password"]'));
    const confirm = element(by.css('[formcontrolname="confirm"]'));
    const button = element(by.buttonText('Finish'));

    return element(by.css('.mat-checkbox-label')).click().then(() => {
      password.clear();
      password.sendKeys('password');
      confirm.clear();
      confirm.sendKeys('password');

      return button.isEnabled();
    });
  }
}
