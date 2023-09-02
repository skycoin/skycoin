export class TransactionsPage {
  navigateTo() {
    return browser.url('/#/transactions');
  }

  getHeaderText() {
    return $('.title').getText();
  }

  getTransactions() {
    return $$('.-transaction');
  }

  getTransactionsCount() {
    return this.getTransactions().length;
  }

  getTransactionDetailIsShow() {
    return $('app-transaction-detail').isExisting();
  }

  showTransactionsModal() {
    return this.getTransactions()[0].click().then(() => {
      return this.getTransactionDetailIsShow();
    });
  }

  hideTransactionModal() {
    const el = $('app-transaction-detail .-header img');

    return browser.waitUntil(require("wdio-wait-for").visibilityOf(el), {
      timeout: 5000
    }).then(() => el.click().then(() => {
      return this.getTransactionDetailIsShow();
    }));
  }
}
