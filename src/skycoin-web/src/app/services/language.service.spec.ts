import { TestBed } from '@angular/core/testing';
import { TranslateService } from '@ngx-translate/core';

import { LanguageService } from './language.service';
import { MockTranslateService } from '../utils/test-mocks';

describe('LanguageService', () => {
  let service: LanguageService;
  let translate: MockTranslateService;
  let savedLang: string | null;
  let savedCookie: string | null;

  function readCookie(): string | null {
    const part = document.cookie.split(';').map(p => p.trim()).find(p => p.startsWith('ui2-lang='));
    return part ? part.substring('ui2-lang='.length) : null;
  }

  function writeCookie(value: string | null) {
    document.cookie = value === null
      ? 'ui2-lang=;path=/;max-age=0'
      : 'ui2-lang=' + value + ';path=/';
  }

  beforeEach(() => {
    savedLang = localStorage.getItem('lang');
    savedCookie = readCookie();
    localStorage.removeItem('lang');
    writeCookie(null);

    TestBed.configureTestingModule({
      providers: [
        LanguageService,
        { provide: TranslateService, useClass: MockTranslateService }
      ]
    });
    service = TestBed.inject(LanguageService);
    translate = TestBed.inject(TranslateService) as unknown as MockTranslateService;
    spyOn(translate, 'use').and.callThrough();
  });

  afterEach(() => {
    if (savedLang === null) {
      localStorage.removeItem('lang');
    } else {
      localStorage.setItem('lang', savedLang);
    }
    writeCookie(savedCookie);
  });

  it('starts in English when nothing was chosen', () => {
    service.loadLanguageSettings();
    expect(translate.use).toHaveBeenCalledWith('en');
  });

  it('starts in the language the dashboard saved', () => {
    localStorage.setItem('lang', 'zh');
    service.loadLanguageSettings();
    expect(translate.use).toHaveBeenCalledWith('zh');
  });

  it('falls back to the ui2 cookie and keeps a dashboard language it lacks', () => {
    localStorage.setItem('lang', 'de');
    writeCookie('zh');
    service.loadLanguageSettings();
    expect(translate.use).toHaveBeenCalledWith('zh');
    expect(localStorage.getItem('lang')).toBe('de');
  });

  it('saves a language the user picks', () => {
    service.loadLanguageSettings();
    service.changeLanguage('es');
    expect(localStorage.getItem('lang')).toBe('es');
  });
});
