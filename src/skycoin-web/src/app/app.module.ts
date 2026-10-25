import { NgModule, inject, provideAppInitializer } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { MATERIAL_ANIMATIONS } from '@angular/material/core';
import { RouterModule } from '@angular/router';
import { ServiceWorkerModule, SwUpdate } from '@angular/service-worker';

import { AppComponent } from './app.component';
import { useOwnServiceWorker } from './service-worker';
import { SkycoinWalletModule } from './wallet.module';

@NgModule({
    declarations: [AppComponent],
    bootstrap: [AppComponent],
    imports: [
        BrowserModule,
        RouterModule.forRoot([], { useHash: true }),
        SkycoinWalletModule,
        ServiceWorkerModule.register('ngsw-worker.js', {
            enabled: useOwnServiceWorker(),
            registrationStrategy: 'registerWhenStable:30000',
        }),
    ],
    providers: [
        // The page has always run with Material's animations off.
        { provide: MATERIAL_ANIMATIONS, useValue: { animationsDisabled: true } },
        // A cached version the server no longer has cannot be repaired in place.
        provideAppInitializer(() => {
            const updates = inject(SwUpdate);
            if (updates.isEnabled) {
                updates.unrecoverable.subscribe(() => document.location.reload());
            }
        }),
    ],
})
export class AppModule { }
