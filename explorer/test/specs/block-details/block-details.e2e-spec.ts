import { BlockDetailsPage } from './block-details.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Block Details Page', () => {
  const page = new BlockDetailsPage();
  const generalFunctions = new GeneralPageFunctions();

  beforeEach(() => { });

  it('should display the title', async () => {
    await generalFunctions.navigateTo('/app/block/5');
    await generalFunctions.waitUntilNoSpinners();
    
    expect(await generalFunctions.getPageTitle()).toEqual('Block Details');
  });

  it('should display 6 block details rows', async () => {
    expect(await generalFunctions.getDetailsRowCount()).toEqual(6);
  });

  it('should show the correct block hight', async () => {
    expect(await page.getBlockHeight()).toBe(5);
  });

  it('should show a valid timestamp', async () => {
    expect(await page.getTimestampValidity()).toBeTruthy();
  });

  it('should show the correct size', async () => {
    expect(await page.getSize()).toBe(317);
  });

  it('should show the correct block hash', async () => {
    expect(await page.getBlockHash()).toEqual('114fe60587a158428a47e0f9571d764f495912c299aa4e67fc88004cf21b0c24');
  });

  it('should show the correct parent block hash', async () => {
    expect(await page.getParentHash()).toEqual('415e47348a1e642cb2e31d00ee500747d3aed0336aabfff7d783ed21465251c7');
  });

  it('should show the correct amount', async () => {
    expect(await page.getAmount()).toBe(999990);
  });

  it('should show the correct transaction ID', async () => {
    expect(await generalFunctions.getTransactionId(0)).toEqual('0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69');
  });

  it('should show the correct transaction inputs', async () => {
    expect(await generalFunctions.getTransactionInputs(0)).toBe('R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ,R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ');
  });

  it('should show the correct transaction outputs', async () => {
    expect(await generalFunctions.getTransactionOutputs(0)).toBe('R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ,2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF');
  });

  it('should have the correct coins amount', async () => {
    expect(await generalFunctions.getTransactionInputsAndOutputsTotalCoins()).toBe(2013866);
  });

  it('should display the "Previous block" button', async () => {
    expect(await page.getNavigationButtonText(0)).toBe('Previous block'.toUpperCase());
  });

  it('should display the "Next block" button', async () => {
    expect(await page.getNavigationButtonText(1)).toBe('Next block'.toUpperCase());
  });

  it('should navigate to the previous block using the nav button', async () => {
    await page.clickNavigationButton(0).then(async () => {
      await generalFunctions.waitUntilNoSpinners();
      expect(await page.getBlockHash()).toEqual('415e47348a1e642cb2e31d00ee500747d3aed0336aabfff7d783ed21465251c7');
    });
  });

  it('should navigate to the next block using the nav button', async () => {
    await page.clickNavigationButton(1).then(async () => {
      await generalFunctions.waitUntilNoSpinners();
      expect(await page.getBlockHash()).toEqual('114fe60587a158428a47e0f9571d764f495912c299aa4e67fc88004cf21b0c24');
    });
  });

  it('should show the error message', async () => {
    await generalFunctions.navigateTo('/app/block/-1');
    await generalFunctions.waitUntilErrorMessage();

    expect(await generalFunctions.getErrorMessage()).toBeDefined();
  });
});
