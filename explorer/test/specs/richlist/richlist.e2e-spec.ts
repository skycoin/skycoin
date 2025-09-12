import { RichlistPage } from './richlist.po';
import { GeneralPageFunctions } from '../general.po';

describe('skycoin-explorer Rich List Page', () => {
  const page = new RichlistPage();
  const generalFunctions = new GeneralPageFunctions();

  beforeEach(() => { });

  it('should display the title', async () => {
    await generalFunctions.navigateTo('/app/richlist');
    await generalFunctions.waitUntilNoSpinners();
    
    expect(await page.getPageTitle()).toEqual('Rich List');
  });

  it('should contain 20 Entries', async () => {
    await generalFunctions.waitUntilNoSpinners();
    expect(await page.getEntriesCount()).toEqual(20);
  });

  it('should have valid entry numbers', async () => {
    for (let i = 0; i < 20; i++) { expect(await page.getEntryNumber(i)).toBeGreaterThan(0); }
  });

  it('should have valid addresses', async () => {
    for (let i = 0; i < 20; i++) {
      const AddressLength = await page.getAddressLength(i);
      expect(AddressLength).toBeGreaterThan(26);
    }
  });

  it('should have valid amounts', async () => {
    for (let i = 0; i < 20; i++) { expect(await page.getAmount(i)).toBeGreaterThan(0); }
  });

});
