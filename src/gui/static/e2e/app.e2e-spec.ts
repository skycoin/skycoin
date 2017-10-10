import { DesktopwalletPage } from './app.po';

describe('desktopwallet App', () => {
  let page: DesktopwalletPage;

  beforeEach(() => {
    page = new DesktopwalletPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});
