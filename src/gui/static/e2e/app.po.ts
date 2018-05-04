import { browser, by, element } from 'protractor';

export class DesktopwalletPage {
  navigateTo() {
    return browser.get('/');
  }

  getParagraphText() {
    return element(by.css('.title')).getText();
  }
}
