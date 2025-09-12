import { TransactionDetailPage } from './transaction-detail.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Transaction Page', () => {
  const page = new TransactionDetailPage();
  const generalFunctions = new GeneralPageFunctions();

  beforeEach(() => { });

  it('should display the title', async () => {
    await generalFunctions.navigateTo('/app/transaction/0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69');
    await generalFunctions.waitUntilNoSpinners();
    
    expect(await generalFunctions.getPageTitle()).toBe('Transaction');
  });

  it('should display 3 transaction details rows', async () => {
    expect(await generalFunctions.getDetailsRowCount()).toEqual(3);
  });

  it('should show the correct size', async () => {
    expect(await page.getSize()).toBe(317);
  });

  it('should show the correct block number', async () => {
    expect(await page.getBlockNumber()).toBe(5);
  });

  it('should show the correct transaction ID', async () => {
    expect(await generalFunctions.getTransactionId(0)).toEqual('0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69');
  });

  it('should show a valid transaction date', async () => {
    expect(await generalFunctions.getTransactionDateValidity(0)).toBeTruthy();
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

  it('should show the error message', async () => {
    await generalFunctions.navigateTo('/app/transaction/0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b68');
    expect(await generalFunctions.getErrorMessage()).toBeDefined();
  });
});
