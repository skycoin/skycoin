import { UnspentOutputsPage } from './unspent-outputs.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Unspent Outputs Page', () => {
  const page = new UnspentOutputsPage();
  const generalFunctions = new GeneralPageFunctions();

  let for180BlockchainTesting = true;
  let address = '';

  beforeEach(() => { });

  it('should display the title', async () => {
    await generalFunctions.navigateTo('/');

    for180BlockchainTesting = !((await browser.getUrl()).includes('localhost:4200'));
    address = for180BlockchainTesting ? '2A2YC8kxWnUDbscpzZ6UPfNAmx5ddKBeYNs' : '24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x7';

    await generalFunctions.navigateTo('/app/unspent/' + address);
    await generalFunctions.waitUntilNoSpinners();

    expect(await generalFunctions.getPageTitle()).toBe('Unspent Outputs');
  });

  it('should display 4 details rows', async () => {
    expect(await generalFunctions.getDetailsRowCount()).toEqual(4);
  });

  it('should show the correct address', async () => {
    expect(await page.getAddressText()).toBe(address);
  });

  it('should show the correct number of outputs', async () => {
    expect(await page.getNumberOfOutputs()).toBe(1);
  });

  it('should show the correct number of coins', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getTotalCoins()).toBe(100);
    } else {
      expect(await page.getTotalCoins()).toBe(0.001);
    }
  });

  it('should show the correct origin transaction ID', async () => {
    if (for180BlockchainTesting) {
      expect(await generalFunctions.getTransactionId(0)).toEqual('b29222c08f10b8bc4ea18981519a3b0e02b9c9cec63ee28d9ffa2efcaf2a8e5a');
    } else {
      expect(await generalFunctions.getTransactionId(0)).toEqual('f375dfb0cea3f082daa03341f6283a483a9857206442feabc94f2b05b1df1fab');
    }
  });

  it('should show a valid transaction date', async () => {
    expect(await generalFunctions.getTransactionDateValidity(0)).toBeTruthy();
  });

  it('should show the correct output hash', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getOutputId(0)).toBe('0b5e5259c276ac949de97062492ea6dc93ae6215c8dd1615862907e3c3ae9cf0');
    } else {
      expect(await page.getOutputId(0)).toBe('74ef6da8c13ec55f28a99aab6bcf890b8decea8c7c56ed55f0b9b4b8cef069da');
    }
  });

  it('should have the correct coins amount', async () => {
    if (for180BlockchainTesting) {
      expect(await page.getOutputsTotalCoins()).toBe(100);
    } else {
      expect(await page.getOutputsTotalCoins()).toBe(0.001);
    }
  });

  it('should show the error message', async () => {
    await generalFunctions.navigateTo('/app/unspent/24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x8');
    expect(await generalFunctions.getErrorMessage()).toBeDefined();
  });
});
