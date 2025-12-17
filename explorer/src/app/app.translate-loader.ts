import { from,  Observable } from 'rxjs';
import { TranslateLoader } from '@ngx-translate/core';

// Loads the translation files, with cache busting.
export class AppTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(import(`../assets/i18n/${lang}.json`));
  }
}
