import { GeneralPageFunctions } from '../general.po';

export class BlockDetailsPage {

  getBlockHeight() {
    return $('.element-details > div:nth-of-type(1) > div').getText()
      .then(text => Number(text));
  }

  getTimestampValidity() {
    return $('.element-details > div:nth-of-type(2) > div').getText()
      .then(text => GeneralPageFunctions.processAndCheckDate(text));
  }

  getSize() {
    return $('.element-details > div:nth-of-type(3) > div').getText()
      .then(text => Number(text.split(' ')[0]));
  }

  getBlockHash() {
    return $$('.element-details .-row a.-link')[0].getText();
  }

  getParentHash() {
    return $$('.element-details .-row a.-link')[1].getText();
  }

  getAmount() {
    return $('.element-details > div:nth-of-type(6) > div').getText()
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }

  getNavigationButtonText(index) {
    return $$('.header-container > div > a')[index].$('.btn-text-e2e').getText();
  }

  clickNavigationButton(index) {
    return $$('.header-container > div > a')[index].$('.-not-xs').click();
  }
}
