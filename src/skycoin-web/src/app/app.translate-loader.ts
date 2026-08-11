import { TranslateLoader } from '@ngx-translate/core';
import { Observable, from } from 'rxjs';

export class AppTranslateLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    // Unwrap the module's default export rather than handing the module
    // namespace object straight to ngx-translate. esbuild turns a JSON import
    // into a module with one named export per top-level key plus `default`, and
    // keys that are not valid JS identifiers — change-coin, pending-txs,
    // qr-code, hardware-wallet — cannot become named exports, so every string
    // under them would silently render as its raw key. webpack exposed them.
    return from(import(`../assets/i18n/${lang}.json`).then(m => (m as any).default ?? m));
  }
}
