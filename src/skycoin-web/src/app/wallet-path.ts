import { Inject, InjectionToken, Pipe, PipeTransform } from '@angular/core';

/**
 * Where the wallet's routes are mounted. The standalone app mounts them at the
 * root; an app that hosts the wallet as one of its sections provides its own
 * base, e.g. '/wallet', and every link and navigation inside the wallet then
 * stays under it.
 */
export const WALLET_BASE_PATH = new InjectionToken<string>('WALLET_BASE_PATH', {
  providedIn: 'root',
  factory: () => '/',
});

/** walletPath('/wallet', '/history') is '/wallet/history'; at the root base it is '/history'. */
export function walletPath(base: string, path: string): string {
  return base.replace(/\/+$/, '') + '/' + path.replace(/^\/+/, '');
}

/** The template form of walletPath: [routerLink]="'/history' | walletPath". */
@Pipe({ name: 'walletPath', standalone: false })
export class WalletPathPipe implements PipeTransform {
  constructor(@Inject(WALLET_BASE_PATH) private base: string) {}

  transform(path: string): string {
    return walletPath(this.base, path);
  }
}
