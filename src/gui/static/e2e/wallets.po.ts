import { browser, by, element, protractor } from 'protractor';

export class WalletsPage {
  navigateTo() {
    return browser.get('/#/wallets');
  }

  getHeaderText() {
    return element(by.css('.title')).getText();
  }

  showAddWallet() {
    const btnAdd = element(by.buttonText('New Wallet'));

    return btnAdd.click().then(() => {
      return element(by.css('app-create-wallet')).isPresent();
    });
  }

  showLoadWallet() {
    const btnLoad = element(by.buttonText('Load Wallet'));

    return btnLoad.click().then(() => {
      return element(by.css('app-create-wallet')).isPresent();
    });
  }

  getWalletModalTitle() {
    return element(by.css('app-create-wallet .-header')).getText();
  }

  fillWalletForm(label: string, seed: string, confirm: string|null, goToManualSeedMode = true) {

    if (goToManualSeedMode) {
      element(by.css('.seed-type-button >div')).click();
      if (confirm !== null) {
        element(by.css('.e2e-confirm-checkbox')).click();
      }
      element(by.buttonText('Continue')).click();
    }

    const labelEl = element(by.css('[formcontrolname="label"]'));
    const seedEl = element(by.css('[formcontrolname="seed"]'));
    const btn = element(by.buttonText(confirm ? 'Create' : 'Load'));
    const encrypt = element(by.css('.mat-checkbox-label'));

    encrypt.click();
    labelEl.clear();
    seedEl.clear();
    labelEl.sendKeys(label);
    seedEl.sendKeys(seed);

    if (confirm) {
      const confirmEl = element(by.css('[formcontrolname="confirm_seed"]'));
      confirmEl.clear();
      confirmEl.sendKeys(confirm);
    }

    if (label !== '' && (seed === confirm || (!confirm && seed !== ''))) {
      browser.sleep(1000);
      const seedValidationCheckBox = element(by.css('.alert-box .mat-checkbox-inner-container'));
      seedValidationCheckBox.click();
    }

    return btn.isEnabled().then(status => {
      if (status) {
        btn.click();
      }

      return status;
    });
  }

  expandWallet() {
    return this.getWalletWithName('Test create wallet').click().then(() => {
      return element(by.css('app-wallet-detail')).isPresent();
    });
  }

  showQrDialog() {
    return browser.sleep(1000).then(() => element(by.css('app-qr-code-button')).click().then(() => {
      return element(by.css('app-qr-code')).isPresent();
    }));
  }

  hideQrDialog() {
    return browser.sleep(1000).then(() => element(by.css('app-modal .-header img')).click().then(() => {
      return element(by.css('app-qr-code')).isPresent();
    }));
  }

  addAddress() {
    return element.all(by.css('.-record')).count().then(originalCount => {
      return element(by.css('.-address-options')).click().then(() => {
        return browser.sleep(2000).then(() => {
          return element(by.css('.top-line')).click().then(() => {
            return browser.sleep(2000).then(() => {
              return element(by.buttonText('Create')).click().then(() => {
                return browser.sleep(2000).then(() => {
                  return element.all(by.css('.-record')).count().then(newCount => {
                    return newCount > originalCount;
                  });
                });
              });
            });
          });
        });
      });
    });
  }

  getCountOfEmptyAddresses(clickSelector: string) {
    return element(by.css(clickSelector)).click().then(() => {
      return element.all(by.css('.-record > div:nth-child(3)')).filter((address) => {
        return address.getText().then(value => {
          return value === '0';
        });
      }).count();
    });
  }

  showChangeWalletName() {
    return element(by.css('.-rename-wallet')).click().then(() => {
      return element(by.css('app-change-name')).isPresent();
    });
  }

  changeWalletName() {
    const label = element(by.css('[formcontrolname="label"]'));
    const btn = element(by.buttonText('Rename'));

    return label.clear().then(() => {
      return label.sendKeys('New Wallet Name').then(() => {
        return btn.click().then(() => {
          return browser.sleep(1000).then(() => {
            return this.getWalletWithName('New Wallet Name').isPresent();
          });
        });
      });
    });
  }

  canEncrypt() {
    return element(by.css('.-enable-encryption')).click().then(() => {
      const p1 = element(by.css('[formcontrolname="password"]'));
      const p2 = element(by.css('[formcontrolname="confirm_password"]'));
      const btn = element(by.buttonText('Proceed'));

      p1.sendKeys('password');
      p2.sendKeys('password');

      return btn.click().then(() => {
        return browser.wait(
          protractor.ExpectedConditions.stalenessOf(element(by.css('app-password-dialog'))),
          30000,
          'Can not encrypt wallet',
        ).then(() => {
          return element(by.css('.-disable-encryption')).isPresent();
        });
      });
    });
  }

  canDecrypt() {
    return element(by.css('.-disable-encryption')).click().then(() => {
      const p1 = element(by.css('[formcontrolname="password"]'));
      const btn = element(by.buttonText('Proceed'));

      p1.clear();
      p1.sendKeys('password');

      return btn.click().then(() => {
        return browser.wait(
          protractor.ExpectedConditions.stalenessOf(element(by.css('app-password-dialog'))),
          30000,
          'Can not decrypt wallet',
        ).then(() => {
          return element(by.css('.-enable-encryption')).isPresent();
        });
      });
    });
  }

  showPriceInformation() {
    return element(by.css('.balance p.dollars')).getText().then(text => {
      return this.checkHeaderPriceFormat(text);
    });
  }

  waitForWalletToBeCreated() {
    browser.wait(
      protractor.ExpectedConditions.stalenessOf(element(by.css('app-create-wallet'))),
      10000,
      'Wallet was not created',
    );
  }

  private getWalletWithName(name: string) {
    return element.all(by.css('.e2e-wallets.ng-star-inserted'))
      .filter(wallet => wallet.element(by.css('.e2e-label')).getText().then(text => text === name))
      .first();
  }

  private checkHeaderPriceFormat(price: string) {
    const reg = /^\$[0-9,]+.[0-9]{2}\s\(\$[0-9,]+.[0-9]{2}\)$/;

    return !!price.match(reg);
  }
}
