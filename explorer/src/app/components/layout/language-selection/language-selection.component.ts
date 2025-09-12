import { Component, OnDestroy, HostListener } from '@angular/core';

import { LanguageService, LanguageData } from '../../../services/language/language.service';
import { Subscription } from 'rxjs';

/**
 * Small button with a flag, normally shown at the right of the search bar, for changing
 * the language. It includes the language selection menu.
 */
@Component({
  selector: 'app-language-selection',
  templateUrl: './language-selection.component.html',
  styleUrls: ['./language-selection.component.scss']
})
export class LanguageSelectionComponent implements OnDestroy {
  /**
   * List of all available languages.
   */
  langs: LanguageData[];
  /**
   * Currently selected language.
   */
  currentLanguage: LanguageData;
  /**
   * If the language selection menu should be visible.
   */
  showLangs = false;

  /**
   * Temporary var for knowing if the user clicked inside or outside this control.
   */
  private clickedInside = false;
  private languageChangeSubscription: Subscription;

  constructor(
    private languageService: LanguageService,
  ) {
    this.langs = languageService.languages;

    // Update the current language if the user changes it.
    this.languageChangeSubscription = this.languageService.currentLanguage.subscribe(lang => this.currentLanguage = lang);
  }

  ngOnDestroy() {
    this.languageChangeSubscription.unsubscribe();
  }

  /**
   * Shows of hides the langugage selection menu when the user clicks the button.
   */
  langsButtonClick() {
    this.showLangs = !this.showLangs;
  }

  /**
   * Changes the selected language.
   *
   * @param language New selected language.
   */
  changeLanguage(language: LanguageData) {
    // Change the UI language.
    this.languageService.changeLanguage(language.code);
    // Close the menu.
    this.showLangs = false;
  }

  /**
   * Event for when the user clicks inside this control. It is called before click(), so it
   * allow to tell click() that the users clicked inside the control and the menu should not
   * be closed.
   */
  @HostListener('click')
  clickInside() {
    this.clickedInside = true;
  }

  /**
   * Event for when the user clicks anywere.
   */
  @HostListener('document:click')
  click() {
    // If the user clicked outside this control (clickInside() did not set this.showLangs
    // to true), the menu is closed.
    if (!this.clickedInside) {
      this.showLangs = false;
    }

    // Resset this.clickedInside for the process to continue working.
    this.clickedInside = false;
  }
}
