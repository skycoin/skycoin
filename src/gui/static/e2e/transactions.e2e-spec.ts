import { TransactionsPage } from './transactions.po';
import { browser } from 'protractor';

describe('Transactions', () => {
  const page = new TransactionsPage();

  it('should display title', () => {
    page.navigateTo();
    browser.sleep(1000);
    expect<any>(page.getHeaderText()).toEqual('Transactions');
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
