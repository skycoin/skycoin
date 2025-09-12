import { Injectable } from '@angular/core';
import { TranslateService, LangChangeEvent } from '@ngx-translate/core';
import { ReplaySubject, Observable } from 'rxjs';

import { languageConfig } from 'app/app.config';

/**
 * Properties of a language available in the application.
 */
export class LanguageData {
  /**
   * Language code, for TranslateService (ngx-translate).
   */
  code: string;
  /**
   * Languege name (written in the language itself).
   */
  name: string;
  /**
   * Name of the file with the icon flag, with the file extension. The file must be inside
   * the assets/img/lang folder.
   */
  iconName: string;

  constructor(langObj) {
    Object.assign(this, langObj);
  }
}

/**
 * Service responsible for controlling the language in which the UI is displayed.
 */
@Injectable()
export class LanguageService {

  private currentLanguageInternal = new ReplaySubject<LanguageData>(1);
  /**
   * Allows to know which language is currently selected.
   */
  get currentLanguage(): Observable<LanguageData> {
    return this.currentLanguageInternal.asObservable();
  }

  private languagesInternal: LanguageData[] = [];
  /**
   * List of available languages.
   */
  get languages(): LanguageData[] {
    return this.languagesInternal;
  }

  /**
   * Key for saving the currently selecteg langugage in localStorage.
   */
  private readonly storageKey = 'lang';
  /**
   * Lets know if the initialize() function has already been called.
   */
  private initialized = false;

  constructor(
    private translate: TranslateService
  ) { }

  initialize() {
    if (this.initialized) {
      return;
    }
    this.initialized = true;

    // Build the available languages list.
    const langs: string[] = [];
    languageConfig.languages.forEach(lang => {
      const langObj = new LanguageData(lang);
      this.languagesInternal.push(langObj);
      langs.push(langObj.code);
    });

    // Initialize ngx-translate.
    this.translate.addLangs(langs);
    this.translate.setDefaultLang(languageConfig.defaultLanguage);
    this.translate.onLangChange.subscribe(this.onLanguageChanged.bind(this));

    // Load the last saved language.
    const storedLang = localStorage.getItem(this.storageKey);
    const currentLang = !!storedLang ? storedLang : languageConfig.defaultLanguage; // eslint-disable-line no-extra-boolean-cast
    this.changeLanguage(currentLang);
  }

  /**
   * Changes the language of the UI
   *
   * @param langCode Language code for TranslateService (ngx-translate). Can be found
   * in LanguageData.code.
   */
  changeLanguage(langCode: string) {
    this.translate.use(langCode);
  }

  /**
   * Event dispatched when the language of the UI is changed.
   *
   * @param event Event data.
   */
  private onLanguageChanged(event: LangChangeEvent) {
    // Update the currently selected language and save it in localStorage.
    this.currentLanguageInternal.next(this.languages.find(val => val.code === event.lang));
    localStorage.setItem(this.storageKey, event.lang);
  }
}
