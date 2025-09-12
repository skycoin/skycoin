export class AddressDetailPage {
  getPageTitleForSmallScreens() {
    return $('.element-details-wrapper > h2:nth-child(2)').getText();
  }

  getAddressForSmallScreens() {
    return $('.element-details > div:nth-of-type(1) > span:nth-of-type(2)').getText().then(address => address.trim());
  }

  getNumberOfTransactions() {
    return $('.element-details > div:nth-of-type(2) > div').getText()
      .then(text => Number(text));
  }

  getTotalReceived() {
    return $('.element-details > div:nth-of-type(3) > div').getText()
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }

  getTotalSent() {
    return $('.element-details > div:nth-of-type(4) > div').getText()
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }

  getCurrentBalance() {
    return $('.balance-container > div:nth-of-type(1)').getText()
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }

  getInitialBalance(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$('.-balance-variation > div:nth-of-type(2) > div:nth-of-type(2)').getText()
      .then(text => Number(text.replace(new RegExp(',', 'g'), '')));
  }

  getFinalBalance(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$('.-balance-variation > div:nth-of-type(3) > div:nth-of-type(2)').getText()
      .then(text => Number(text.replace(new RegExp(',', 'g'), '')));
  }

  getBalanceChange(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$('.-label').getText()
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }
}
