import { browser, by, element } from 'protractor';
import { Page } from './page';

export class DesktopwalletPage extends Page {
  navigateTo() {
    return browser.get('/');
  }

  getParagraphText() {
    return element(by.css('.-header')).getText();
  }
}
