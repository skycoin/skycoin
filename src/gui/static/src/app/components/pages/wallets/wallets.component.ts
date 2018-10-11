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

  hwCompatibilityActivated = false;

  constructor(
    public walletService: WalletService,
    private dialog: MatDialog,
  ) {
    if (window['isElectron']) {
      this.hwCompatibilityActivated = window['ipcRenderer'].sendSync('hwCompatibilityActivated')
    }
  }

  addWallet(create) {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.data = { create };
    this.dialog.open(CreateWalletComponent, config);
  }

  adminHwWallet() {
    alert(window['ipcRenderer'].sendSync('hwGetDevice') === null ?
    'Error: no hardware wallet detected' :
    'Hardware wallet detected');
  }

  toggleWallet(wallet: Wallet) {
    wallet.opened = !wallet.opened;
  }
}
