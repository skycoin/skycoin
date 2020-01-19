import { TransactionsPage } from './transactions.po';

describe('Transactions', () => {
  const page = new TransactionsPage();

  it('should display title', () => {
    page.navigateTo();
    expect<any>(page.getHeaderText()).toEqual('History');
  });

  it('should contain transactions', () => {
    expect<any>(page.getTransactionsCount()).toBeGreaterThan(0);
  });

  it('should show transaction detail modal', () => {
    expect<any>(page.showTransactionsModal()).toBeTruthy();
  });

  it('should hide transaction detail modal', () => {
    expect<any>(page.hideTransactionModal()).toBeFalsy();
  });
});
