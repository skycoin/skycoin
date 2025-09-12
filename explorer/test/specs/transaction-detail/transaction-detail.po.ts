export class TransactionDetailPage {
  getSize() {
    return $('.element-details > div:nth-of-type(2) > div').getText()
      .then(text => Number(text.split(' ')[0]));
  }

  getBlockNumber() {
    return $('.element-details > div:nth-of-type(3) > div > a').getText()
      .then(text => Number(text.replace(new RegExp(',', 'g'), '')));
  }
}
