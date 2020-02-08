import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { first } from 'rxjs/operators';

import { WalletsAndAddressesService } from './wallet-operations/wallets-and-addresses.service';

/**
 * Forces the app to always redirect the user to the wizard if there are no
 * registered wallets.
 */
@Injectable()
export class WizardGuardService implements CanActivate {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private router: Router,
  ) { }

  /**
   * Function called by the system to know if it is possible to navigate to a page.
   */
  canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot): Promise<boolean> {
    return new Promise<boolean>(resolve => {
      this.walletsAndAddressesService.allWallets.pipe(first()).subscribe(wallets => {
        if (wallets.length === 0) {
          this.router.navigate(['/wizard']);

          return resolve(false);
        }

        return resolve(true);
      });
    });
  }
}
