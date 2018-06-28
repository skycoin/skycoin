import { browser, by, element } from 'protractor';

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
    return element(by.css('app-transaction-detail .-header img')).click().then(() => {
      return this.getTransactionDetailIsShow();
    });
  }
}
