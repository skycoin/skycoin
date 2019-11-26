import { browser, by, element, protractor } from 'protractor';

export class OnboardingCreatePage {
  navigateTo() {
    return browser.get('/#/wizard');
  }

  getHeaderText() {
    return element(by.css('.-header span')).getText();
  }

  selectLanguage() {
    browser.sleep(1000);

    return element(by.css('.e2e-language-modal')).isPresent().then(languageSelectionIsShown => {
      if (!languageSelectionIsShown) {
        return true;
      }

      return element.all(by.css('.e2e-language-modal .button')).first().click().then(() => {
        const el = element(by.css('.e2e-language-modal'));

        return browser.wait(protractor.ExpectedConditions.invisibilityOf(el), 5000).then(() => true);
      });
    });
  }

  getSafeguardIsShown() {
    return element(by.css('app-onboarding-safeguard')).isPresent();
  }

  acceptSafeguard() {
   return element.all(by.css('app-modal .mat-checkbox-label')).first().click().then(() => {
     return element(by.buttonText('Continue')).click().then(() => {
        return this.getSafeguardIsShown();
      });
    });
  }

  createWallet(goToManualSeedMode = true) {
    element(by.buttonText('New')).click();

    if (goToManualSeedMode) {
      element(by.css('.seed-type-button >div')).click();
      element(by.css('.e2e-confirm-checkbox')).click();
      element(by.buttonText('Continue')).click();
    }

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

    browser.sleep(1000);
    const seedValidationCheckBox = element(by.css('.-check'));
    seedValidationCheckBox.click();

    return btnCreate.isEnabled().then(status => {
      if (status) {
        btnCreate.click();
      }

      return status;
    });
  }

  loadWallet() {
    element(by.buttonText('Load')).click();

    element(by.css('.seed-type-button >div')).click();
    element(by.buttonText('Continue')).click();

    const label = element(by.css('[formcontrolname="label"]'));
    const seed = element(by.css('[formcontrolname="seed"]'));
    const btnLoad = element(by.buttonText('Create'));

    label.clear();
    label.sendKeys('Test wallet');
    seed.clear();
    seed.sendKeys('test test');

    browser.sleep(1000);
    const seedValidationCheckBox = element(by.css('.-check'));
    seedValidationCheckBox.click();

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
