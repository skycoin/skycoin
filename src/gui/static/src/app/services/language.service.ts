import { Injectable } from '@angular/core';
import { TranslateService, LangChangeEvent } from '@ngx-translate/core';
import { ReplaySubject } from 'rxjs/ReplaySubject';

import { AppConfig } from '../app.config';
import { ISubscription } from 'rxjs/Subscription';
import { StorageService, StorageType } from './storage.service';

export class LanguageData {
  code: string;
  name: string;
  iconName: string;

  constructor(langObj) {
    Object.assign(this, langObj);
  }
}

@Injectable()
export class LanguageService {
  currentLanguage = new ReplaySubject<LanguageData>(1);
  selectedLanguageLoaded = new ReplaySubject<boolean>(1);

  private readonly storageKey = 'lang';
  private subscription: ISubscription;
  private languagesInternal: LanguageData[] = [];

  get languages(): LanguageData[] {
    return this.languagesInternal;
  }

  constructor(
    private translate: TranslateService,
    private storageService: StorageService,
  ) { }

  loadLanguageSettings() {

    const langs: string[] = [];
    AppConfig.languages.forEach(lang => {
      const LangObj = new LanguageData(lang);
      this.languagesInternal.push(LangObj);
      langs.push(LangObj.code);
    });

    this.translate.addLangs(langs);
    this.translate.setDefaultLang(AppConfig.defaultLanguage);

    this.translate.onLangChange
      .subscribe((event: LangChangeEvent) => this.onLanguageChanged(event));

    this.loadCurrentLanguage();
  }

  changeLanguage(langCode: string) {
    this.translate.use(langCode);
  }

  private onLanguageChanged(event: LangChangeEvent) {
    this.currentLanguage.next(this.languages.find(val => val.code === event.lang));

    if (this.subscription) {
      this.subscription.unsubscribe();
    }
    this.subscription = this.storageService.store(StorageType.CLIENT, this.storageKey, event.lang).subscribe();
  }

  private loadCurrentLanguage() {
    this.storageService.get(StorageType.CLIENT, this.storageKey).subscribe(response => {
      if (response.data && this.translate.getLangs().indexOf(response.data) !== -1) {
        setTimeout(() => { this.translate.use(response.data); }, 16);
        this.selectedLanguageLoaded.next(true);
      } else {
        this.selectedLanguageLoaded.next(false);
      }
    }, () => {
      this.selectedLanguageLoaded.next(false);
    });
  }
}
