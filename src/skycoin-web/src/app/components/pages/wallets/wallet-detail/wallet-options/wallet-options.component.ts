import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { environment } from '../../../../../../environments/environment';
import { Wallet } from '../../../../../app.datatypes';
import { BaseCoin } from '../../../../../coins/basecoin';
import { CoinService } from '../../../../../services/coin.service';

export enum WalletOptionsResponses {
  UnlockWallet,
  AddNewAddress,
  EditWallet,
  DeleteWallet,
  EncryptWallet,
  DecryptWallet,
}

@Component({
    selector: 'app-wallet-options',
    templateUrl: './wallet-options.component.html',
    styleUrls: ['./wallet-options.component.scss'],
    standalone: false
})
export class WalletOptionsComponent {

  showUnlockOption: boolean;
  showEncryptOption: boolean;
  isEncrypted: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) data: Wallet,
    public dialogRef: MatDialogRef<WalletOptionsComponent>,
    coinService: CoinService,
  ) {
    this.showUnlockOption = !environment.production && !data.seed;
    const coin = coinService.currentCoin.getValue();
    this.showEncryptOption = !!(coin && coin.serverWallets && data.filename);
    this.isEncrypted = !!data.encrypted;
  }

  onUnlockWallet() {
    this.closePopup(WalletOptionsResponses.UnlockWallet);
  }

  onAddNewAddress() {
    this.closePopup(WalletOptionsResponses.AddNewAddress);
  }

  onEditWallet() {
    this.closePopup(WalletOptionsResponses.EditWallet);
  }

  onDeleteWallet() {
    this.closePopup(WalletOptionsResponses.DeleteWallet);
  }

  onEncryptWallet() {
    this.closePopup(WalletOptionsResponses.EncryptWallet);
  }

  onDecryptWallet() {
    this.closePopup(WalletOptionsResponses.DecryptWallet);
  }

  closePopup(response: WalletOptionsResponses = null) {
    this.dialogRef.close(response);
  }
}
