import { browser, by, element, protractor } from 'protractor';

export class WalletsPage {
  navigateTo() {
    return browser.get('/#/wallets');
  }

  getHeaderText() {
    return element(by.css('.title')).getText();
  }

  showAddWallet() {
    const btnAdd = element(by.buttonText('Add Wallet'));

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

  fillWalletForm(label: string, seed: string, confirm: string|null) {
    const labelEl = element(by.css('[formcontrolname="label"]'));
    const seedEl = element(by.css('[formcontrolname="seed"]'));
    const btn = element(by.buttonText(confirm ? 'Create' : 'Load'));

    labelEl.clear();
    seedEl.clear();
    labelEl.sendKeys(label);
    seedEl.sendKeys(seed);

    if (confirm) {
      const confirmEl = element(by.css('[formcontrolname="confirm_seed"]'));
      confirmEl.clear();
      confirmEl.sendKeys(seed);
    }

    return btn.isEnabled();
  }

  createWallet() {
    const label = element(by.css('[formcontrolname="label"]'));
    const seed = element(by.css('[formcontrolname="seed"]'));
    const confirm = element(by.css('[formcontrolname="confirm_seed"]'));
    const encrypt = element(by.css('.mat-checkbox-label'));
    const btnCreate = element(by.buttonText('Create'));

    label.clear();
    label.sendKeys('Test create wallet');
    seed.clear();
    seed.sendKeys('test create wallet');
    confirm.clear();
    confirm.sendKeys('test create wallet');
    encrypt.click();

    return btnCreate.isEnabled().then(status => {
      if (status) {
        btnCreate.click();
      }

      return status;
    });
  }

  loadWallet() {
    const label = element(by.css('[formcontrolname="label"]'));
    const seed = element(by.css('[formcontrolname="seed"]'));
    const encrypt = element(by.css('.mat-checkbox-label'));
    const btnLoad = element(by.buttonText('Load'));

    label.clear();
    label.sendKeys('Test load wallet');
    seed.clear();
    seed.sendKeys('test load wallet');
    encrypt.click();

    return btnLoad.isEnabled().then(status => {
      if (status) {
        btnLoad.click();
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
    return browser.sleep(1000).then(() => element(by.css('.qr-code-button')).click().then(() => {
      return element(by.css('app-qr-code')).isPresent();
    }));
  }

  hideQrDialog() {
    return browser.sleep(1000).then(() => element(by.css('app-modal .-header img')).click().then(() => {
      return element(by.css('app-qr-code')).isPresent();
    }));
  }

  addAddress() {
    return element.all(by.css('.-detail')).count().then(originalCount => {
      return element(by.css('.-new-address')).click().then(() => {
        return browser.sleep(2000).then(() => {
          return element.all(by.css('.-detail')).count().then(newCount => {
            return newCount > originalCount;
          });
        });
      });
    });
  }

  hideEmptyAddress() {
    return element(by.css('.-hide-empty')).click().then(() => {
      return element.all(by.css('.-detail > div:nth-child(3)')).filter((address) => {
        return address.getText().then(value => {
          return value === '0';
        });
      }).count();
    });
  }

  showEmptyAddress() {
    return element.all(by.css('.-show-empty')).first().click().then(() => {
      return element.all(by.css('.-detail')).count().then(count => {
        return count > 0;
      });
    });
  }

  showChangeWalletName() {
    return element(by.css('.-edit-wallet')).click().then(() => {
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
    return element.all(by.css('.-table.ng-star-inserted'))
      .filter(wallet => wallet.element(by.css('.-label')).getText().then(text => text === name))
      .first();
  }

  private checkHeaderPriceFormat(price: string) {
    const reg = /^\$[0-9,]+.[0-9]{2}\s\(\$[0-9,]+.[0-9]{2}\)$/;

    return !!price.match(reg);
  }
}
