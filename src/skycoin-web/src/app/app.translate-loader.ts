import { TranslateLoader } from '@ngx-translate/core';
import { Observable, from } from 'rxjs';

export class AppTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(import(`../assets/i18n/${lang}.json`));
  }
}
