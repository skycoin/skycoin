import { GeneralPageFunctions } from '../general.po';

export class BlocksPage {
  getPageTitle(titleNumber) {
    return $$('app-root h2')[titleNumber].getText();
  }

  getExplorerApiText() {
    return $$('.-link')[0].getText();
  }

  getUnconfirmedTransactionsText() {
    return $$('.-link')[1].getText();
  }

  getRichListText() {
    return $$('.-link')[2].getText();
  }

  getStatsLabelsCount() {
    return $$('.-label').length;
  }

  getStat(statIndex) {
    return $$('.-value')[statIndex].getText()
      .then(text => Number(text.replace(new RegExp(',', 'g'), '')));
  }

  getBlocksListRowCount() {
    return $$('.table a.-row').length;
  }

  async waitForHavingBlockRows() {
    await browser.waitUntil(
      async () => (await this.getBlocksListRowCount()) > 0,
      {
          timeout: 5000,
          timeoutMsg: 'No rows after 5s'
      }
    );
  }

  getTableCellValue(row, column) {
    return $$('.table a.-row')[row].$$('div')[column].getText();
  }

  getTimeValidity(row) {
    return this.getTableCellValue(row, 1)
      .then(text => GeneralPageFunctions.processAndCheckDate(text));
  }

  getBlockNumber(row) {
    return this.getTableCellValue(row, 5)
      .then(text => Number(text));
  }

  getTransactionCount(row) {
    return this.getTableCellValue(row, 6)
      .then(text => Number(text));
  }

  getBlockHashLength(row) {
    return this.getTableCellValue(row, 7)
      .then(text => text.length);
  }

  getIfPaginationControlExists() {
    return $('.pagination').isExisting();
  }

  navigateToTheNextPage() {
    return $('.-next').click();
  }

  navigateToThePreviousPage() {
    return $('.-previous').click();
  }

  navigateToTheLastPage() {
    return $('.-last').click();
  }

  navigateToTheFirstPage() {
    return $('.-first').click();
  }

  pressPageButton(index) {
    return $$('.-page')[index].click();
  }
}
