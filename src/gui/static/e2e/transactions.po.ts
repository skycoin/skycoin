import { browser, by, element, protractor } from 'protractor';

export class TransactionsPage {
  navigateTo() {
    return browser.get('/#/transactions');
  }

  getHeaderText() {
    return element(by.css('.title')).getText();
  }

  getTransactions() {
    return element.all(by.css('.-transaction'));
  }

  getTransactionsCount() {
    return this.getTransactions().count();
  }

  getTransactionDetailIsShow() {
    return element(by.css('app-transaction-detail')).isPresent();
  }

  showTransactionsModal() {
    return this.getTransactions().first().click().then(() => {
      return this.getTransactionDetailIsShow();
    });
  }

  hideTransactionModal() {
    const el = element(by.css('app-transaction-detail .-header img'));

    return browser.wait(protractor.ExpectedConditions.visibilityOf(el), 5000).then(() => el.click().then(() => {
      return this.getTransactionDetailIsShow();
    }));
  }
}
