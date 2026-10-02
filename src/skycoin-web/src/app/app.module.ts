import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { MATERIAL_ANIMATIONS } from '@angular/material/core';
import { RouterModule } from '@angular/router';

import { AppComponent } from './app.component';
import { SkycoinWalletModule } from './wallet.module';

@NgModule({
    declarations: [AppComponent],
    bootstrap: [AppComponent],
    imports: [
        BrowserModule,
        RouterModule.forRoot([], { useHash: true }),
        SkycoinWalletModule,
    ],
    providers: [
        // The page has always run with Material's animations off.
        { provide: MATERIAL_ANIMATIONS, useValue: { animationsDisabled: true } },
    ],
})
export class AppModule { }
