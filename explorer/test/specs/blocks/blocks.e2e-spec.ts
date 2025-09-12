import { BlocksPage } from './blocks.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Blocks Page', () => {
  const page = new BlocksPage();
  const generalFunctions = new GeneralPageFunctions();

  beforeEach(() => { });

  it('should display the title for big screens', async () => {
    await generalFunctions.navigateTo('/');
    await generalFunctions.waitUntilNoSpinners();
    
    expect(await page.getPageTitle(0)).toBe('Blocks');
  });
  
  it('should display the first title for small screens', async () => {
    await browser.setWindowRect(0, 0, 400, 500);
    expect(await page.getPageTitle(1)).toBe('Stats');
  });
  
  it('should display the second title for small screens', async () => {
    expect(await page.getPageTitle(2)).toBe('Blocks');
  });

  it('should display the Explorer API link', async () => {
    await browser.setWindowRect(0, 0, 1920, 1080);
    expect(await page.getExplorerApiText()).toEqual('Explorer API');
  });

  it('should display the Unconfirmed Transactions link', async () => {
    expect(await page.getUnconfirmedTransactionsText()).toEqual('Unconfirmed Transactions');
  });

  it('should display the Rich List link', async () => {
    expect(await page.getRichListText()).toEqual('Rich List');
  });

  it('should contain 5 stats labels', async () => {
    expect(await page.getStatsLabelsCount()).toEqual(5);
  });

  it('should contain valid stat values', async () => {
    expect(await page.getStat(0)).toBeGreaterThan(0);
    expect(await page.getStat(1)).toBeGreaterThan(0);
    expect(await page.getStat(2)).toBeGreaterThan(0);
    expect(await page.getStat(3)).toBeGreaterThan(0);
    expect(await page.getStat(4)).toBeGreaterThan(0);
  });

  it('should contain 10 blocks items', async () => {
    await page.waitForHavingBlockRows();
    expect(await page.getBlocksListRowCount()).toEqual(10);
  });

  it('should have valid block times', async () => {
    for (let i = 0; i < 10; i++) { expect(await page.getTimeValidity(i)).toBeTruthy(); }
  });

  it('should have valid block numbers', async () => {
    for (let i = 0; i < 10; i++) { expect(await page.getBlockNumber(i)).toBeGreaterThan(0); }
  });

  it('should have valid transaction counts', async () => {
    for (let i = 0; i < 10; i++) { expect(await page.getTransactionCount(i)).toBeGreaterThan(0); }
  });

  it('should have valid block hashes', async () => {
    for (let i = 0; i < 10; i++) { expect(await page.getBlockHashLength(i)).toEqual(64); }
  });
  
  it('should contain a pagination control', async () => {
    expect(await page.getIfPaginationControlExists()).toBe(true);
  });

  it('should navigate to the next page and show 10 block items', async () => {
    const blockCount = await page.navigateToTheNextPage().then(async () => {
      await page.waitForHavingBlockRows();
      return page.getBlocksListRowCount();
    });

    expect(blockCount).toEqual(10);
  });

  it('should navigate to the previous page and show 10 block items', async () => {
    const blockCount = await page.navigateToThePreviousPage().then(async () => {
      await page.waitForHavingBlockRows();
      return page.getBlocksListRowCount();
    });

    expect(blockCount).toEqual(10);
  });

  it('should navigate to the last page and show at least 1 block item', async () => {
    const blockCount = await page.navigateToTheLastPage().then(async () => {
      await page.waitForHavingBlockRows();
      return page.getBlocksListRowCount();
    });

    expect(blockCount).toBeGreaterThan(0);
  });

  it('should navigate to the first page and show 10 block items', async () => {
    const blockCount = await page.navigateToTheFirstPage().then(async () => {
      await page.waitForHavingBlockRows();
      return page.getBlocksListRowCount();
    });

    expect(blockCount).toEqual(10);
  });

  it('should navigate to the second page and show 10 block items', async () => {
    const blockCount = await page.pressPageButton(2).then(async () => {
      await page.waitForHavingBlockRows();
      return page.getBlocksListRowCount();
    });

    expect(blockCount).toEqual(10);
  });
});
