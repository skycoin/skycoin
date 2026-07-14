import { SendSkycoinComponent } from './components/pages/send-skycoin/send-skycoin.component';
import { WizardGuardService } from './services/wizard-guard.service';
import { WalletsComponent } from './components/pages/wallets/wallets.component';
import { PendingTransactionsComponent } from './components/pages/settings/pending-transactions/pending-transactions.component';
import { OutputsComponent } from './components/pages/settings/outputs/outputs.component';
import { BlockchainComponent } from './components/pages/settings/blockchain/blockchain.component';
import { HistoryComponent } from './components/pages/history/history.component';
import { OnboardingEncryptWalletComponent } from './components/pages/onboarding/onboarding-encrypt-wallet/onboarding-encrypt-wallet.component';
import { BuyComponent } from './components/pages/buy/buy.component';
import { OnboardingCreateWalletComponent } from './components/pages/onboarding/onboarding-create-wallet/onboarding-create-wallet.component';
import { NodesComponent } from './components/pages/settings/nodes/nodes.component';
import { Routes } from "@angular/router";

export const AppRoutes: Routes = [
  {
    path: '',
    redirectTo: 'wallets',
    pathMatch: 'full',
  },
  {
    path: 'wallets',
    component: WalletsComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'send',
    component: SendSkycoinComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'history',
    component: HistoryComponent,
    canActivate: [WizardGuardService],
  },
  {
    path: 'buy',
    component: BuyComponent,
    canActivate: [WizardGuardService],
  },
  // Node settings are intentionally NOT behind the wizard guard: choosing/
  // verifying the backing node is inherently a pre-wallet concern — you need a
  // reachable, correct node before you can create or scan a wallet, and if the
  // default node is down you would otherwise be deadlocked in the wizard. This
  // more-specific route is declared before the guarded `settings` block so it
  // wins the match.
  {
    path: 'settings/node',
    component: NodesComponent,
  },
  {
    path: 'settings',
    children: [
      {
        path: 'blockchain',
        component: BlockchainComponent,
      },
      {
        path: 'outputs',
        component: OutputsComponent,
      },
      {
        path: 'pending-transactions',
        component: PendingTransactionsComponent,
      },
    ],
    canActivate: [WizardGuardService],
  },
  {
    path: 'wizard',
    children: [
      {
        path: '',
        redirectTo: 'create',
        pathMatch: 'full',
      },
      {
        path: 'create',
        component: OnboardingCreateWalletComponent,
      },
      {
        path: 'encrypt',
        component: OnboardingEncryptWalletComponent,
      },
    ],
  },
];
