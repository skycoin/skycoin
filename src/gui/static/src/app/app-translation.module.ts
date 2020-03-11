import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { from, Observable } from 'rxjs';
import { NgModule } from '@angular/core';

// Loads the translation files, with cache busting.
export class TranslationModuleLoader implements TranslateLoader {
  getTranslation(lang: string): Observable<any> {
    return from(import(`../assets/i18n/${lang}.json`));
  }
}

@NgModule({
  imports: [TranslateModule.forRoot({
    loader: {
      provide: TranslateLoader,
      useClass: TranslationModuleLoader,
    },
  })],
  exports: [TranslateModule],
})
export class AppTranslationModule { }
