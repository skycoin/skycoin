import { browser } from 'protractor';

export class Page {
  dontWait() {
    browser.waitForAngularEnabled(false);
  }
}
