import { browser, by, element, protractor } from 'protractor';

export class SendPage {
  navigateTo() {
    return browser.get('/#/send');
  }

  getHeaderText() {
    return element(by.css('.title')).getText();
  }

  getWalletsCount() {
    browser.actions().sendKeys(protractor.Key.ESCAPE).perform();
    
    return element(by.css('.mat-select')).click().then(() => {
      return element.all(by.css('.mat-select-panel mat-option .mat-option-text')).count();
    });
  }

  getWalletsWithCoins() {
    browser.actions().sendKeys(protractor.Key.ESCAPE).perform();

    return element(by.css('.mat-select')).click().then(() => {
      return element.all(by.css('.mat-select-panel mat-option .mat-option-text')).filter((opt) => {
        return opt.getText().then((v) => {
          return this.getCoinsFromOptionString(v) > 0;
        });
      });
    });
  }

  getValidWallets() {
    browser.actions().sendKeys(protractor.Key.ESCAPE).perform();

    return element(by.css('.mat-select')).click().then(() => {
      return element.all(by.css('.mat-select-panel .mat-active .mat-option-text')).filter((opt) => {
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

  fillFormWithCoins(coins: string) {
    const dest = element(by.css('[formcontrolname="address"]'));
    const amount = element(by.css('[formcontrolname="coins"]'));
    const btnSend = element(by.buttonText('Send'));

    dest.clear();
    amount.clear();

    return dest.sendKeys('2e1erPpaxNVC37PkEv3n8PESNw2DNr5aJNy').then(() => {
      return this.getValidWallets().then(wallets => {
        return wallets[0].click().then(() => {
          return amount.sendKeys(coins).then(() => {
            return btnSend.isEnabled();
          });
        });
      });
    });
  }

  private getCoinsFromOptionString(option: string) {
    const value = option.slice(option.indexOf('-') + 1, option.indexOf(' SKY'));

    return parseFloat(value);
  }
}
