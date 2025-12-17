import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { environment } from '../../../../../../environments/environment';
import { Wallet } from '../../../../../app.datatypes';

export enum WalletOptionsResponses {
  UnlockWallet,
  AddNewAddress,
  EditWallet,
  DeleteWallet,
}

@Component({
  selector: 'app-wallet-options',
  templateUrl: './wallet-options.component.html',
  styleUrls: ['./wallet-options.component.scss'],
})
export class WalletOptionsComponent {


  showUnlockOption: boolean;
  constructor(
    @Inject(MAT_DIALOG_DATA) data: Wallet,
    public dialogRef: MatDialogRef<WalletOptionsComponent>
  ) {
    this.showUnlockOption = !environment.production && !data.seed;
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

  closePopup(response: WalletOptionsResponses = null) {
    this.dialogRef.close(response);
  }
}
