import { ChangeDetectionStrategy, Component } from '@angular/core';

/** The standalone app is only the wallet, mounted at the root (wallet.module.ts). */
@Component({
    selector: 'app-root',
    template: '<router-outlet></router-outlet>',
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class AppComponent {}
