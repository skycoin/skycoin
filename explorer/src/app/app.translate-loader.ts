import { from,  Observable } from 'rxjs';
import { TranslateLoader } from '@ngx-translate/core';

// Loads the translation files, with cache busting.
export class AppTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    // Unwrap the module's default export rather than handing the module
    // namespace object straight to ngx-translate.
    //
    // esbuild turns a JSON import into a module with one named export per
    // top-level key, plus `default`. Keys that are not valid JS identifiers —
    // anything hyphenated, such as "offline-transactions" — cannot become named
    // exports, so they are missing from the namespace object and every string
    // under them silently renders as its raw key. webpack exposed them, which
    // is why this only appears after the esbuild migration. `default` is always
    // the fully parsed JSON, whatever the bundler.
    return from(import(`../assets/i18n/${lang}.json`).then(m => (m as any).default ?? m));
  }
}
