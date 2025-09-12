import { AddressDetailPage } from './address-detail.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Address Page', () => {
  const page = new AddressDetailPage();
  const generalFunctions = new GeneralPageFunctions();
  
  let for180BlockchainTesting = true;
  let address = '';

  beforeEach(() => { });

  it('should display the correct title', async () => {
    await generalFunctions.navigateTo('/');

    for180BlockchainTesting = !((await browser.getUrl()).includes('localhost:4200'));
    address = for180BlockchainTesting ? '2A2YC8kxWnUDbscpzZ6UPfNAmx5ddKBeYNs' : '24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x7';

    await generalFunctions.navigateTo('/app/address/' + address);
    await generalFunctions.waitUntilNoSpinners();

    expect(await generalFunctions.getPageTitle()).toBe(address);
  });

  it('should have the title for small screens', async () => {
    await browser.setWindowRect(0, 0, 400, 500);
    expect(await page.getPageTitleForSmallScreens()).toBe('Address');
  });

  it('should have 5 address details rows', async () => {
    await browser.setWindowRect(0, 0, 1920, 1080);
    expect(await generalFunctions.getDetailsRowCount()).toEqual(5);
  });

  it('should have the correct address for small screens', async () => {
    await browser.setWindowRect(0, 0, 400, 500);
    expect(await page.getAddressForSmallScreens()).toEqual(address);
  });

  it('should show the correct # of transactions', async () => {
    await browser.setWindowRect(0, 0, 1920, 1080);
    expect(await page.getNumberOfTransactions()).toBe(1);
  });

  it('should show the correct received amount', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getTotalReceived()).toBe(100);
    } else {
      expect(await page.getTotalReceived()).toBe(0.001);
    }
  });

  it('should show the correct sent amount', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getTotalSent()).toBe(0);
    } else {
      expect(await page.getTotalSent()).toBe(0);
    }
  });

  it('should show the correct current balance', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getCurrentBalance()).toBe(100);
    } else {
      expect(await page.getCurrentBalance()).toBe(0.001);
    }
  });

  it('should show the correct transaction ID', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionId(0)).toEqual('b29222c08f10b8bc4ea18981519a3b0e02b9c9cec63ee28d9ffa2efcaf2a8e5a');
    } else {
      expect(await generalFunctions.getTransactionId(0)).toEqual('f375dfb0cea3f082daa03341f6283a483a9857206442feabc94f2b05b1df1fab');
    }
  });

  it('should show a valid transaction date', async () => {
    expect(await generalFunctions.getTransactionDateValidity(0)).toBeTruthy();
  });

  it('should show the correct transaction inputs', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionInputs(0)).toBe('R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ');
    } else {
      expect(await generalFunctions.getTransactionInputs(0)).toBe('TXXyvkvs7Lt3EFtvTh51i4BagiJaGARsva');
    }
  });

  it('should show the correct transaction outputs', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionOutputs(0)).toBe('R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ,2A2YC8kxWnUDbscpzZ6UPfNAmx5ddKBeYNs');
    } else {
      expect(await generalFunctions.getTransactionOutputs(0)).toBe('24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x7,TXXyvkvs7Lt3EFtvTh51i4BagiJaGARsva');
    }
  });

  it('should have the correct coins amount', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionInputsAndOutputsTotalCoins()).toBe(1655151);
    } else {
      expect(await generalFunctions.getTransactionInputsAndOutputsTotalCoins()).toBe(158.002);
    }
  });

  it('should show the correct initial address balance', async () => {
    expect(await page.getInitialBalance(0)).toBe(0);
  });

  it('should show the correct final address balance', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getFinalBalance(0)).toBe(100);
    } else {
      expect(await page.getFinalBalance(0)).toBe(0.001);
    }
  });

  it('should show the correct balance variation', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getBalanceChange(0)).toBe(100);
    } else {
      expect(await page.getBalanceChange(0)).toBe(0.001);
    }
  });

  it('should show the error message', async () => {
    await generalFunctions.navigateTo('/app/address/24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x8');
    await generalFunctions.waitUntilErrorMessage();

    expect(await generalFunctions.getErrorMessage()).toBeDefined();
  });
});
