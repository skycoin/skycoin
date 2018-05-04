import { DesktopwalletPage } from './app.po';

describe('desktopwallet App', () => {
  let page: DesktopwalletPage;

  beforeEach(() => {
    page = new DesktopwalletPage();
  });

  it('should show wallets page', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Wallets');
  });
});
