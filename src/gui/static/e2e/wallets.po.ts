export class WalletsPage {
  navigateTo() {
    return browser.url('/#/wallets');
  }

  getHeaderText() {
    return $('.title').getText();
  }

  showAddWallet() {
    const btnAdd = $("button=New Wallet");

    return btnAdd.click().then(async () => {
      await browser.pause(500);
      return $('app-create-wallet').isExisting();
    });
  }

  showLoadWallet() {
    const btnLoad = $("button=Load Wallet");

    return btnLoad.click().then(() => {
      return $('app-create-wallet').isExisting();
    });
  }

  getWalletModalTitle() {
    return $('app-create-wallet .-header').getText();
  }

  async fillWalletForm(label, seed, confirm, goToManualSeedMode = true) {

    if (goToManualSeedMode) {
      await $('.seed-type-button >span').click();
      if (confirm !== null) {
        await $('.e2e-confirm-checkbox').click();
      }
      await $("button=Continue").click();
    }

    const labelEl = $('[formcontrolname="label"]');
    const seedEl = $('[formcontrolname="seed"]');
    const btn = $(`button=${confirm ? 'Create' : 'Load'}`);
    const encrypt = $('.mat-mdc-checkbox');

    await encrypt.click();
    await labelEl.clearValue();
    await seedEl.clearValue();
    await labelEl.setValue(label);
    await seedEl.setValue(seed);

    if (confirm) {
      const confirmEl = $('[formcontrolname="confirm_seed"]');
      await confirmEl.clearValue();
      await confirmEl.setValue(confirm);
    }

    if (label !== '' && (seed === confirm || (!confirm && seed !== ''))) {
      browser.pause(1000);
      const seedValidationCheckBox = $('.alert-box .mat-mdc-checkbox');
      await seedValidationCheckBox.click();
    }

    return btn.isEnabled().then(async status => {
      if (status) {
        await btn.click();
      }

      return status;
    });
  }

  expandWallet() {
    return this.getWalletWithName('Test create wallet').click().then(() => {
      return $('app-wallet-detail').isExisting();
    });
  }

  showQrDialog() {
    return browser.pause(1000).then(() => $('app-qr-code-button').click().then(() => {
      return $('app-qr-code').isExisting();
    }));
  }

  hideQrDialog() {
    return browser.pause(1000).then(() => $('app-modal .-header img').click().then(() => {
      return $('app-qr-code').isExisting();
    }));
  }

  addAddress() {
    return $$('.-record').length.then(originalCount => {
      return $('.-address-options').click().then(() => {
        return browser.pause(2000).then(() => {
          return $('.top-line').click().then(() => {
            return browser.pause(2000).then(() => {
              return $("button=Create").click().then(() => {
                return browser.pause(2000).then(() => {
                  return $$('.-record').length.then(newCount => {
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

  getCountOfEmptyAddresses(clickSelector) {
    return $(clickSelector).click().then(async () => {
      return (await $$('.-record > div:nth-child(3)').filter((address) => {
        return address.getText().then(value => {
          return value === '0';
        });
      })).length;
    });
  }

  showChangeWalletName() {
    return $('.-rename-wallet').click().then(() => {
      return $('app-change-name').isExisting();
    });
  }

  changeWalletName() {
    const label = $('[formcontrolname="label"]');
    const btn = $("button=Rename");

    return label.clearValue().then(() => {
      return label.setValue('New Wallet Name').then(() => {
        return btn.click().then(() => {
          return browser.pause(1000).then(() => {
            return this.getWalletWithName('New Wallet Name').isExisting();
          });
        });
      });
    });
  }

  canEncrypt() {
    return $('.-enable-encryption').click().then(async () => {
      const p1 = $('[formcontrolname="password"]');
      const p2 = $('[formcontrolname="confirm_password"]');
      const btn = $("button=Proceed");

      await p1.click();
      await p1.setValue('password');
      await p2.click();
      await p2.setValue('password');

      return btn.click().then(() => {
        return browser.waitUntil(require("wdio-wait-for").stalenessOf($('app-password-dialog')), {
          timeout: 30000,
          timeoutMsg: 'Can not encrypt wallet'
        }).then(() => {
          return $('.-disable-encryption').isExisting();
        });
      });
    });
  }

  canDecrypt() {
    return $('.-disable-encryption').click().then(async () => {
      const p1 = $('[formcontrolname="password"]');
      const btn = $("button=Proceed");

      await p1.click();
      await p1.clearValue();
      await p1.setValue('password');

      return btn.click().then(() => {
        return browser.waitUntil(require("wdio-wait-for").stalenessOf($('app-password-dialog')), {
          timeout: 30000,
          timeoutMsg: 'Can not decrypt wallet'
        }).then(() => {
          return $('.-enable-encryption').isExisting();
        });
      });
    });
  }

  async waitForWalletToBeCreated() {
    await browser.waitUntil(require("wdio-wait-for").stalenessOf($('app-create-wallet')), {
      timeout: 10000,
      timeoutMsg: 'Wallet was not created'
    });
  }

  getWalletWithName(name) {
    return $$('.e2e-wallets.ng-star-inserted')
      .filter(wallet => wallet.$('.e2e-label').getText().then(text => text === name))[0];
  }

  checkHeaderPriceFormat(price) {
    const reg = /^\$[0-9,]+.[0-9]{2}\s\(\$[0-9,]+.[0-9]{2}\)$/;

    return !!price.match(reg);
  }
}
