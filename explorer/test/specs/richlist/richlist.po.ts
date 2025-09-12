export class RichlistPage {
  getPageTitle() {
    return $$('app-root h2')[0].getText();
  }

  getEntriesCount() {
    return $$('.table a.-row').length;
  }

  getTableCellValue(row, column) {
    return $$('.table a.-row')[row].$$('div')[column].getText();
  }

  getEntryNumber(row) {
    return this.getTableCellValue(row, 1)
      .then(text => Number(text));
  }

  getAddressLength(row) {
    return this.getTableCellValue(row, 2)
    .then(text => text.length);
  }

  getAmount(row) {
    return this.getTableCellValue(row, 3)
      .then(text => Number(text.split(' ')[0].replace(new RegExp(',', 'g'), '')));
  }
}
