export class GeneralPageFunctions {

  static processAndCheckDate(text) {
    let numberText = text;
    if (isNaN(Number(numberText.substr(0, 1)))) {
      numberText = numberText.substr(1, numberText.length - 1);
    }
    return !isNaN((new Date(numberText)).getTime());
  }

  async navigateTo(route) {
    return browser.url(route);
  }

  goBack() {
    return browser.back();
  }

  getPageTitle() {
    return $('.element-details-wrapper h2').getText();
  }

  getDetailsRowCount() {
    return $$('.element-details .-row').length
      .then(count => {
        return count;
      });
  }

  getTransactionId(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$('.-row a').getText();
  }

  getTransactionDateValidity(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$('.-date').getText()
      .then(text => GeneralPageFunctions.processAndCheckDate(text));
  }

  getTransactionInputs(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$$('.-data > .row > div:nth-of-type(1) a')
      .map((elem, i) => elem.getText())
      .then(texts => texts.join(','));
  }

  getTransactionOutputs(transsactionIndex) {
    return $$('.transaction')[transsactionIndex].$$('.-data > .row > div:nth-of-type(2) a')
      .map((elem, i) => elem.getText())
      .then(texts => texts.join(','));
  }

  getTransactionInputsAndOutputsTotalCoins() {
    return $$('.-balance > div:nth-of-type(2)')
      .map((elem, i) => elem.getText())
      .then(texts => texts.map(text => Number(text.replace(new RegExp(',', 'g'), ''))).reduce((total, val) => total + val, 0))
      .then(result => Math.round(result * 1000000) / 1000000);
  }

  getErrorMessage() {
    return $('div.row.-msg-container > div').getText();
  }

  async waitUntilErrorMessage() {
    await browser.waitUntil(
      async () => (await $$('div.row.-msg-container > div').length) > 0,
      {
        timeout: 5000,
        timeoutMsg: 'No error after 5s'
      }
    );
  }

  async waitUntilNoSpinners() {
    await browser.waitUntil(
      async () => (await $$('.fa-spinner').length) === 0,
      {
          timeout: 5000,
          timeoutMsg: 'There are still spinners after 5s'
      }
    );
  }

  async waitUntilPageTitleShown(pageTitle) {
    await browser.waitUntil(
      async () => (await this.getPageTitle()) === pageTitle,
      {
          timeout: 2000,
          timeoutMsg: 'Expected title not found after 2s'
      }
    );
  }
}
