export class NavigarionPage {
  navigateTo() {
    return browser.url('/');
  }

  goToRichlistPage() {
    return $$('.-link')[2].click();
  }

  goToUnconfirmedTransactions() {
    return $$('.-link')[1].click();
  }

  goToBlockDetails() {
    return $$('.table a.-row')[0].click();
  }

  goToTransactionDetail() {
    return $('.transaction .-row a').click();
  }

  goToAddressDetail() {
    return $('.transaction > .-data > .row > div:nth-of-type(2) a').click();
  }

  goToUnspentOutputs() {
    return $('.element-details-wrapper .element-details .-link').click();
  }
}
