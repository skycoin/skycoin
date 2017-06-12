import { WalletNewPage } from './app.po';

describe('wallet-new App', () => {
  let page: WalletNewPage;

  beforeEach(() => {
    page = new WalletNewPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
