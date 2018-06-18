import { TranslateLoader } from '@ngx-translate/core';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/fromPromise';

export class AppTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return Observable.fromPromise(System.import(`../assets/i18n/${lang}.json`));
  }
}
