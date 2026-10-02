import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NoopAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { SkycoinWalletModule } from './wallet.module';

@NgModule({
    declarations: [AppComponent],
    bootstrap: [AppComponent],
    imports: [
        BrowserModule,
        NoopAnimationsModule,
        RouterModule.forRoot([], { useHash: true }),
        SkycoinWalletModule,
    ],
})
export class AppModule { }
