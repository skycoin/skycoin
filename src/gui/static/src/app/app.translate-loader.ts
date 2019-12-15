import { from, Observable } from 'rxjs';
import { TranslateLoader } from '@ngx-translate/core';

export class AppTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(System.import(`../assets/i18n/${lang}.json`));
  }
}
