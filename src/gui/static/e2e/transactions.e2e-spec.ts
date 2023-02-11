import { TransactionsPage } from './transactions.po';

describe('Transactions', () => {
  const page = new TransactionsPage();

  it('should display title', async () => {
    await page.navigateTo();
    expect<any>(await page.getHeaderText()).toEqual('History');
  });

  it('should contain transactions', async () => {
    await browser.waitUntil(
      async () => await $$('.mdc-circular-progress__spinner-layer').length === 0,
      { timeout: 10000, timeoutMsg: 'History data not loaded after 10s' }
    );

    expect<any>(await page.getTransactionsCount()).toBeGreaterThan(0);
  });

  it('should show transaction detail modal', async () => {
    expect<any>(await page.showTransactionsModal()).toBeTruthy();
  });

  it('should hide transaction detail modal', async () => {
    expect<any>(await page.hideTransactionModal()).toBeFalsy();
  });
});
