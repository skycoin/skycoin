import { TranslateLoader, TranslateDirective, TranslatePipe, provideTranslateLoader, provideTranslateService } from '@ngx-translate/core';
import { from, Observable } from 'rxjs';
import { NgModule } from '@angular/core';

// Loads the translation files, with cache busting.
export class TranslationModuleLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(import(`../assets/i18n/${lang}.json`));
  }
}

@NgModule({
  imports: [TranslatePipe, TranslateDirective],
  exports: [TranslatePipe, TranslateDirective],
  providers: [
    provideTranslateService({
      loader: provideTranslateLoader(TranslationModuleLoader),
    }),
  ],
})
export class AppTranslationModule { }
