import { Component, HostListener, OnInit } from '@angular/core';

/**
 * Info about a value saved using the functions from PageBase.
 */
export class LocalValueData {
  /**
   * Actual value.
   */
  value: any;
  /**
   * Moment in which the value was updated for the last time. It is the value
   * of new Date().getTime().
   */
  date: number;
}

/**
 * Base class for the pages. It restores the scroll position when navigating back to the page
 * and provides functions for saving data on the browser history state.
 */
@Component({
  selector: 'app-page-base',
  template: '',
  styles: [],
})
export class PageBaseComponent implements OnInit {
  // Key for saving the scroll position.
  private readonly persistentScrollPosKey = 'scroll-pos';

  // Property to make mandatory calling super.ngOnInit on child classes.
  private static mustCallNgOnInitSuper = Symbol('You must call super.ngOnInit.');
  ngOnInit() {
    // Restore the scroll position. If there is no saved value, go to the top.
    let lastScrollPos = this.getLocalValue(this.persistentScrollPosKey);
    lastScrollPos = lastScrollPos ? lastScrollPos.value : '0';
    window.scrollTo(0, Number(lastScrollPos));
    setTimeout(() => window.scrollTo(0, Number(lastScrollPos)), 1);

    return undefined as typeof PageBaseComponent.mustCallNgOnInitSuper & never;
  }

  // Saves the scroll position after each change.
  @HostListener('window:scroll', ['$event'])
  saveScrollPosition(event: any) {
    this.saveLocalValue(this.persistentScrollPosKey, window.scrollY + '');
  }

  /**
   * Saves a value on the state inside window.history.
   *
   * @param key Key for identifying the value.
   * @param value Value to save.
   */
  saveLocalValue(key: string, value: any) {
    const state = window.history.state;
    state[key] = value;
    state[key + '_time'] = new Date().getTime();
    window.history.replaceState(state, '', window.location.pathname);
  }

  /**
   * Gets a value from window.history. If the value is not foun, it returns null.
   *
   * @param key Key that identifies the value.
   */
  getLocalValue(key: string): LocalValueData | null {
    if (!window.history.state || window.history.state[key] === undefined) {
      return null;
    }

    const response = new LocalValueData();
    response.value = window.history.state[key];
    response.date = window.history.state[key + '_time'];
    return response;
  }
}
