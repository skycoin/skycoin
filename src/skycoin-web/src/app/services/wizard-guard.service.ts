import { Injectable, inject } from '@angular/core';
import { WALLET_BASE_PATH, walletPath } from '../wallet-path';
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { Observable, filter, first, map, merge } from 'rxjs';

import { WalletService } from './wallet/wallet.service';

@Injectable()
export class WizardGuardService  {
  private walletBase = inject(WALLET_BASE_PATH);

  constructor(
    private walletService: WalletService,
    private router: Router,
  ) {
  }

  canActivate(next: ActivatedRouteSnapshot,
              state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return new Promise<boolean>((resolve, reject) => {
      // While the wallets are locked the page opens under the lock screen, which
      // has to render before any wallet can load.
      merge(
        this.walletService.locked.pipe(filter(locked => locked), map(() => null)),
        this.walletService.haveWallets,
      ).pipe(first()).subscribe(result => {
        if (result === false) {
          this.router.navigate([walletPath(this.walletBase, '/wizard')]);
          return resolve(false);
        }
        return resolve(true);
      });
    });
  }
}
