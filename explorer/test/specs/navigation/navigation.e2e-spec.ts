import { NavigarionPage } from './navigation.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Navigation', () => {
  const page = new NavigarionPage();
  const generalFunctions = new GeneralPageFunctions();

  beforeEach(() => { });

  it('should open the blocks page', async () => {
    await generalFunctions.navigateTo('/');
  });

  it('should navigate to the richlist page', async () => {
    page.goToRichlistPage();
    await browser.waitUntil(
      async () => await $('.wrapper h2').getText() === 'Rich List',
      { timeout: 2000, timeoutMsg: 'Title not found after 2s' }
    );
  });

  it('should navigate to the unconfirmed transaction page', async () => {
    await generalFunctions.goBack();
    page.goToUnconfirmedTransactions();
    await generalFunctions.waitUntilPageTitleShown('Unconfirmed Transactions');
  });

  it('should navigate to the block details page', async () => {
    await generalFunctions.goBack();
    await generalFunctions.waitUntilNoSpinners();
    page.goToBlockDetails();
    await generalFunctions.waitUntilPageTitleShown('Block Details');
  });

  it('should navigate to the transaction detail page', async () => {
    await generalFunctions.waitUntilNoSpinners();
    page.goToTransactionDetail();
    await generalFunctions.waitUntilPageTitleShown('Transaction');
  });

  it('should navigate to the address detail page', async () => {
    await generalFunctions.goBack();
    await generalFunctions.waitUntilNoSpinners();
    page.goToAddressDetail();
    await browser.waitUntil(
      async () => await (await generalFunctions.getPageTitle()).length > 0,
      { timeout: 2000, timeoutMsg: 'Title not found after 2s' }
    );
  });

  it('should navigate to the unspent outputs page', async () => {
    page.goToUnspentOutputs();
    await generalFunctions.waitUntilPageTitleShown('Unspent Outputs');
  });
});
