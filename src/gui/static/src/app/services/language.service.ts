import { Injectable } from '@angular/core';
import { TranslateService, LangChangeEvent } from '@ngx-translate/core';
import { ReplaySubject, SubscriptionLike, Observable } from 'rxjs';

import { AppConfig } from '../app.config';
import { StorageService, StorageType } from './storage.service';

/**
 * Represents a language the app can use for the UI.
 */
export class LanguageData {
  /**
   * ID of the language.
   */
  code: string;
  /**
   * Language name.
   */
  name: string;
  /**
   * Name of the file containing the flag which is used for identifying the language.
   */
  iconName: string;

  constructor(langObj) {
    Object.assign(this, langObj);
  }
}

/**
 * Allows to know which language is the UI using and to change it.
 */
@Injectable()
export class LanguageService {
  /**
   * Allows to know the currently selected language for the UI.
   */
  get currentLanguage(): Observable<LanguageData> {
    return this.currentLanguageInternal.asObservable();
  }
  currentLanguageInternal = new ReplaySubject<LanguageData>(1);

  /**
   * Allows to know if the service was able to load from persistent storage what the lastest
   * language selected by the user was.
   */
  get savedSelectedLanguageLoaded(): Observable<boolean> {
    return this.savedSelectedLanguageLoadedInternal.asObservable();
  }
  savedSelectedLanguageLoadedInternal = new ReplaySubject<boolean>(1);

  /**
   * List with the languages availabe for the app to use. Values should not be overwritten.
   */
  get languages(): LanguageData[] {
    return this.languagesInternal;
  }
  private languagesInternal: LanguageData[] = [];

  private readonly storageKey = 'lang';
  private subscription: SubscriptionLike;

  constructor(
    private translate: TranslateService,
    private storageService: StorageService,
  ) { }

  /**
   * Makes the service initialize itself and ngx-translate. Should be called when
   * starting the app.
   */
  initialize() {
    // Load the list of available languages from the configuration file.
    const langs: string[] = [];
    AppConfig.languages.forEach(lang => {
      const LangObj = new LanguageData(lang);
      this.languagesInternal.push(LangObj);
      langs.push(LangObj.code);
    });

    // Initialize ngx-translate.
    this.translate.addLangs(langs);
    this.translate.setDefaultLang(AppConfig.defaultLanguage);
    this.translate.onLangChange.subscribe((event: LangChangeEvent) => this.onLanguageChanged(event));

    this.loadCurrentLanguage();
  }

  /**
   * Changes the language used by the UI.
   * @param langCode Code which identifies the desired language.
   */
  changeLanguage(langCode: string) {
    this.translate.use(langCode);
  }

  /**
   * Function called when the language used by ngx-translate is changed. Informs about the
   * change and saves the code of the new selection.
   */
  private onLanguageChanged(event: LangChangeEvent) {
    this.currentLanguageInternal.next(this.languages.find(val => val.code === event.lang));

    if (this.subscription) {
      this.subscription.unsubscribe();
    }
    this.subscription = this.storageService.store(StorageType.CLIENT, this.storageKey, event.lang).subscribe();
  }

  /**
   * Loads from persistent storage what the last language selected by the user was and
   * makes ngx-translate use it.
   */
  private loadCurrentLanguage() {
    this.storageService.get(StorageType.CLIENT, this.storageKey).subscribe(response => {
      if (response.data && this.translate.getLangs().indexOf(response.data) !== -1) {
        setTimeout(() => { this.changeLanguage(response.data); }, 16);
        this.savedSelectedLanguageLoadedInternal.next(true);
      } else {
        this.savedSelectedLanguageLoadedInternal.next(false);
      }
    }, () => {
      this.savedSelectedLanguageLoadedInternal.next(false);
    });
  }
}
