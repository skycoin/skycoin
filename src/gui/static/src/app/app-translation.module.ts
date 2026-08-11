import { TranslateLoader, TranslateDirective, TranslatePipe, provideTranslateService } from '@ngx-translate/core';
import { from, Observable } from 'rxjs';
import { NgModule } from '@angular/core';

// Loads the translation files, with cache busting.
export class TranslationModuleLoader implements TranslateLoader {
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

@NgModule({
  imports: [TranslatePipe, TranslateDirective],
  exports: [TranslatePipe, TranslateDirective],
  providers: [
    provideTranslateService({
      // Declared as an explicit class provider rather than through
      // provideTranslateLoader(). That helper picks between useClass and
      // useFactory with isClass(), which tests Function.prototype.toString()
      // against /^class\s/. esbuild minifies `class TranslationModuleLoader {`
      // down to `class{`, leaving no space for that regex, so the loader would
      // be registered as a factory and Angular would call the constructor
      // without `new`: "Class constructor cannot be invoked without 'new'".
      loader: { provide: TranslateLoader, useClass: TranslationModuleLoader },
    }),
  ],
})
export class AppTranslationModule { }
