import { SearchPage } from './search.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Search Component', () => {
  const page = new SearchPage();
  const generalFunctions = new GeneralPageFunctions();

  beforeEach(() => { });

  it('should find a block by hash', async () => {
    await generalFunctions.navigateTo('/');
    await page.search('114fe60587a158428a47e0f9571d764f495912c299aa4e67fc88004cf21b0c24');
    await browser.keys("\uE007");
    await generalFunctions.waitUntilPageTitleShown('Block Details');
  });

  it('should find a transaction by hash', async () => {
    await page.search('0579e7727627cd9815a8a8b5e1df86124f45a4132cc0dbd00d2f110e4f409b69');
    await browser.keys("\uE007");
    await generalFunctions.waitUntilPageTitleShown('Transaction');
  });

  it('should find a block by number', async () => {
    await page.search('5');
    await browser.keys("\uE007");
    await generalFunctions.waitUntilPageTitleShown('Block Details');
  });

  it('should find an address', async () => {
    await page.search('24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x7');
    await browser.keys("\uE007");
    await generalFunctions.waitUntilPageTitleShown('24ooGeabUGQLnJmoyviqA8h2y7Cgz2CY4x7');
  });
});
