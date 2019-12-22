import { browser, by, element } from 'protractor';

export class SendPage {
  navigateTo() {
    return browser.get('/#/send');
  }

  getHeaderText() {
    return element(by.css('.title')).getText();
  }

  getWalletsCount() {
    return element.all(by.css('#wallet option')).count();
  }

  getWalletsWithCoins() {
    return element.all(by.tagName('#wallet option'))
      .filter((opt) => {
        return opt.getText().then((v) => {
          return this.getCoinsFromOptionString(v) > 0;
        });
      });
  }

  getValidWallets() {
    return element.all(by.tagName('#wallet option'))
      .filter((opt) => {
        return opt.getText().then((v) => {
          return opt.getAttribute('disabled').then(status => {
            return status === null && this.getCoinsFromOptionString(v) > 0;
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
