export class SendPage {
  navigateTo() {
    return browser.url('/#/send');
  }

  getHeaderText() {
    return $('.title').getText();
  }

  async getWalletsCount() {
    await browser.keys("Escape");
    
    return $('.mat-mdc-select').click().then(() => {
      return $$('.mat-mdc-select-panel mat-option span').length;
    });
  }

  async getWalletsWithCoins() {
    await browser.keys("Escape");

    return $('.mat-mdc-select').click().then(() => {
      return $$('.mat-mdc-select-panel mat-option span').filter((opt) => {
        return opt.getText().then((v) => {
          return this.getCoinsFromOptionString(v) > 0;
        });
      });
    });
  }

  async getValidWallets() {
    await browser.keys("Escape");

    return $('.mat-mdc-select').click().then(() => {
      return $$('.mat-mdc-select-panel mat-option span').filter((opt) => {
        return opt.getText().then((v) => {
          return this.getCoinsFromOptionString(v) > 0;
        });
      });
    });
  }

  selectValidWallet() {
    return this.getValidWallets().then(wallets => {
      return wallets[0].click().then(() => {
        return true;
      });
    });
  }

  async fillFormWithCoins(coins) {
    const dest = $('[formcontrolname="address"]');
    const amount = $('[formcontrolname="coins"]');
    const btnSend = $("button=Send");

    await dest.click;
    await dest.clearValue();
    await amount.click;
    await amount.clearValue();

    await dest.click;
    return dest.setValue('2e1erPpaxNVC37PkEv3n8PESNw2DNr5aJNy').then(() => {
      return this.getValidWallets().then(wallets => {
        return wallets[0].click().then(async () => {
          await amount.click;
          return amount.setValue(coins).then(() => {
            return btnSend.isEnabled();
          });
        });
      });
    });
  }

  getCoinsFromOptionString(option) {
    const value = option.slice(option.indexOf('-') + 1, option.indexOf(' SKY'));

    return parseFloat(value);
  }
}
