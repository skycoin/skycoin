import { UnconfirmedTransactionsPage } from './unconfirmed-transactions.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Unconfirmed Transactions Page', () => {
  const page = new UnconfirmedTransactionsPage();
  const generalFunctions = new GeneralPageFunctions();

  let for180BlockchainTesting = true;

  beforeEach(() => { });

  it('should display Unconfirmed Transactions text', async () => {
    await generalFunctions.navigateTo('/app/unconfirmedtransactions');
    await generalFunctions.waitUntilNoSpinners();

    for180BlockchainTesting = !((await browser.getUrl()).includes('localhost:4200'));

    expect(await generalFunctions.getPageTitle()).toEqual('Unconfirmed Transactions');
  });

  it('should display 4 Unconfirmed Transactions details', async () => {
    expect(await generalFunctions.getDetailsRowCount()).toEqual(4);
  });

  it('should show the correct number of transactions', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getNumberOfTransactions()).toBe(1);
    }
  });

  it('should show the correct total size', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getSize()).toBe(220);
    }
  });

  it('should show a valid newest date', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getTimestampValidity(3)).toBeTruthy();
    }
  });

  it('should show a valid oldest date', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getTimestampValidity(4)).toBeTruthy();
    }
  });

  it('should show the correct transaction ID', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionId(0)).toEqual('701d23fd513bad325938ba56869f9faba19384a8ec3dd41833aff147eac53947');
    }
  });

  it('should show a valid transaction date', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionDateValidity(0)).toBeTruthy();
    }
  });

  it('should show the correct transaction inputs', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionInputs(0)).toBe('R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ');
    }
  });

  it('should show the correct transaction outputs', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionOutputs(0)).toBe('R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ,212mwY3Dmey6vwnWpiph99zzCmopXTqeVEN');
    }
  });

  it('should have the correct coins amount', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionInputsAndOutputsTotalCoins()).toBe(63915681);
    }
  });
});
