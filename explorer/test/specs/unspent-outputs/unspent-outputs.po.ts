export class UnspentOutputsPage {

  getAddressText() {
    return $('.element-details-wrapper .element-details .-link').getText();
  }

  getNumberOfOutputs() {
    return $('.element-details > div:nth-of-type(2) > div').getText()
      .then(text => Number(text));
  }

  getTotalCoins() {
    return $('.element-details > div:nth-of-type(3) > div').getText()
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }

  getOutputId(outputIndex) {
    return $$('.transaction')[outputIndex].$('.-data > .row .-body div').getText();
  }

  getOutputsTotalCoins() {
    return $$('.-balance')[0].$$('div:nth-of-type(2)')
      .map((elem, i) => elem.getText())
      .then(texts => texts.map(text => Number(text.replace(new RegExp(',', 'g'), ''))).reduce((total, val) => total + val, 0))
      .then(result => Math.round(result * 1000000) / 1000000);
  }
}
