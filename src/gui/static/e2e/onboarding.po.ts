export class OnboardingCreatePage {
  navigateTo() {
    return browser.url('/#/wizard');
  }

  getHeaderText() {
    return $('.-header span').getText();
  }

  async selectLanguage() {
    await browser.pause(1000);

    return $('.e2e-language-modal').isExisting().then(languageSelectionIsShown => {
      if (!languageSelectionIsShown) {
        return true;
      }

      return $$('.e2e-language-modal .button')[0].click().then(() => {
        const el = $('.e2e-language-modal');

        return browser.waitUntil(require("wdio-wait-for").invisibilityOf(el), {
          timeout: 5000
        }).then(() => true);
      });
    });
  }

  getSafeguardIsShown() {
    return $('app-confirmation').isExisting();
  }

  acceptSafeguard() {
   return $$('app-modal .mat-mdc-checkbox')[0].click().then(() => {
     return $("button=Continue").click().then(() => {
        return this.getSafeguardIsShown();
      });
    });
  }

  async createWallet(goToManualSeedMode = true) {
    await $("button=New").click();

    if (goToManualSeedMode) {
      await $('.seed-type-button >span').click();
      await $('.e2e-confirm-checkbox').click();
      await $("button=Continue").click();
    }

    const label = $('[formcontrolname="label"]');
    const seed = $('[formcontrolname="seed"]');
    const confirm = $('[formcontrolname="confirm_seed"]');
    const btnCreate = $("button=Create");

    await label.clearValue();
    await label.setValue('Test onboarding wallet');
    await seed.clearValue();
    await seed.setValue('test test');
    await confirm.clearValue();
    await confirm.setValue('test test');

    await browser.pause(1000);
    const seedValidationCheckBox = $('.-check');
    await seedValidationCheckBox.click();

    return btnCreate.isEnabled().then(async status => {
      if (status) {
        await btnCreate.click();
      }

      return status;
    });
  }

  async loadWallet() {
    await $("button=Load").click();

    await $('.seed-type-button >span').click();
    await $("button=Continue").click();

    const label = $('[formcontrolname="label"]');
    const seed = $('[formcontrolname="seed"]');
    const btnLoad = $("button=Create");

    await label.clearValue();
    await label.setValue('Test wallet');
    await seed.clearValue();
    await seed.setValue('test test');

    await browser.pause(1000);
    const seedValidationCheckBox = $('.-check');
    await seedValidationCheckBox.click();

    return btnLoad.isEnabled();
  }

  goBack() {
    return $("button=Back").click().then(() => {
      return this.getHeaderText();
    });
  }

  getEncryptWalletCheckbox() {
    return $('.mat-mdc-checkbox .mdc-checkbox__native-control');
  }

  canContinueWithoutEncryption() {
    return $('.mat-mdc-checkbox').click().then(() => {
      return $("button=Finish").isEnabled();
    });
  }

  encryptWallet() {
    const password = $('[formcontrolname="password"]');
    const confirm = $('[formcontrolname="confirm"]');
    const button = $("button=Finish");

    return $('.mat-mdc-checkbox').click().then(async () => {
      await password.click();
      await password.setValue('password');
      await confirm.click();
      await confirm.setValue('password');

      return button.isEnabled();
    });
  }
}
