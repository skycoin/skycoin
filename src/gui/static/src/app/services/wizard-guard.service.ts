import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { first } from 'rxjs/operators';
import { WalletsAndAddressesService } from './wallet-operations/wallets-and-addresses.service';

@Injectable()
export class WizardGuardService implements CanActivate {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private router: Router,
  ) { }

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
