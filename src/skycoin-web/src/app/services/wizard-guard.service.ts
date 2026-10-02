import { Injectable, inject } from '@angular/core';
import { WALLET_BASE_PATH, walletPath } from '../wallet-path';
import { ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';
import { Observable } from 'rxjs';

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
      this.walletService.haveWallets.subscribe(result => {
        if (!result) {
          this.router.navigate([walletPath(this.walletBase, '/wizard')]);
          return resolve(false);
        }
        return resolve(true);
      });
    });
  }
}
