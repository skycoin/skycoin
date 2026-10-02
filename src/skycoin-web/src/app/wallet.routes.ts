import { Routes } from '@angular/router';

import { AppRoutes } from './app.routes';
import { WalletShellComponent } from './components/layout/wallet-shell/wallet-shell.component';

/** The wallet's pages (app.routes.ts), framed by its shell. */
export const WALLET_ROUTES: Routes = [
  {
    path: '',
    component: WalletShellComponent,
    children: AppRoutes,
  },
];
