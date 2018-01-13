import { Component, Input } from '@angular/core';
import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet.service';
import { MdDialog, MdDialogConfig } from '@angular/material';
import { ChangeNameComponent } from '../change-name/change-name.component';

@Component({
  selector: 'app-wallet-detail',
  templateUrl: './wallet-detail.component.html',
  styleUrls: ['./wallet-detail.component.scss']
})
export class WalletDetailComponent {
  @Input() wallet: Wallet;

  constructor(
    private dialog: MdDialog,
    private walletService: WalletService,
  ) { }

  editWallet() {
    const config = new MdDialogConfig();
    config.width = '566px';
    config.data = this.wallet;
    this.dialog.open(ChangeNameComponent, config);
  }

  newAddress() {
    this.walletService.addAddress(this.wallet).subscribe();
  }

  toggleEmpty() {
    this.wallet.hideEmpty = !this.wallet.hideEmpty;
  }
}
