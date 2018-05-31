import { Component } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { CreateWalletComponent } from './create-wallet/create-wallet.component';
import { Wallet } from '../../../app.datatypes';

@Component({
  selector: 'app-wallets',
  templateUrl: './wallets.component.html',
  styleUrls: ['./wallets.component.scss'],
})
export class WalletsComponent {
  constructor(
    public walletService: WalletService,
    private dialog: MatDialog,
  ) {}

  addWallet(create) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = { create };
    this.dialog.open(CreateWalletComponent, config);
  }

  toggleWallet(wallet: Wallet) {
    wallet.opened = !wallet.opened;
  }
}
