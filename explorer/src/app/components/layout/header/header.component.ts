import {Component, OnInit} from '@angular/core';
import { Router } from '@angular/router';
import disableScroll from 'disable-scroll';
import menuData from './menu';

/**
 * Skycoin header. To activate it, FooterConfig.useGenericFooter (in app.config.ts)
 * must be false. Read the docs for more information.
 */
@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  /**
   * If the window is too small, indicates if the menu must me shown.
   */
  menuVisible: boolean;
  /**
   * Data for building the navigation menu.
   */
  menu: any;

  constructor(public router: Router) {
    this.setupMenu();
  }

  /**
   * Resets the menu to the default values. This makes any open submenu to be closed.
   */
  private setupMenu() {
    this.menu = menuData.map(obj => ({...obj}));
  }

  ngOnInit() {
    // If the window is too small, the menu should be hidden until the user clicks the menu button.
    this.menuVisible = false;
  }

  /**
   * If the window is too small, shows of hides the menu.
   *
   * @param value Indicates if the menu should be visible or not.
   */
  toggleMenu(value: boolean) {
    this.menuVisible = value;

    // Deactivate scrolling in the main area while the menu is shown.
    if (value) {
      disableScroll.on();
    } else {
      disableScroll.off();
    }

    // Reset the menu to the default values.
    this.setupMenu();
  }

  /**
   * If the window is too small, opens or closes a submenu.
   *
   * @param item Submenu that was clicked.
   */
  onClickMenuDropdown(item) {
    item.open = !item.open;
  }
}
