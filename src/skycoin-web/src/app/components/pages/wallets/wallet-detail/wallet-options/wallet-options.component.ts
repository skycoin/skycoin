import { Component, Inject, ChangeDetectionStrategy } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';

import { environment } from '../../../../../../environments/environment';
import { Wallet } from '../../../../../app.datatypes';
import { BaseCoin } from '../../../../../coins/basecoin';
import { CoinService } from '../../../../../services/coin.service';
import { WalletService } from '../../../../../services/wallet/wallet.service';

export enum WalletOptionsResponses {
  UnlockWallet,
  AddNewAddress,
  EditWallet,
  DeleteWallet,
  EncryptWallet,
  DecryptWallet,
  ScanAddresses,
}

@Component({
    selector: 'app-wallet-options',
    templateUrl: './wallet-options.component.html',
    styleUrls: ['./wallet-options.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class WalletOptionsComponent {

  showUnlockOption: boolean;
  showEncryptOption: boolean;
  showScanOption: boolean;
  isEncrypted: boolean;

  constructor(
    @Inject(MAT_DIALOG_DATA) data: Wallet,
    public dialogRef: MatDialogRef<WalletOptionsComponent>,
    coinService: CoinService,
    walletService: WalletService,
  ) {
    this.showUnlockOption = !environment.production && !data.seed;
    const coin = coinService.currentCoin.getValue();
    // Show encrypt option for server-managed wallets OR browser-only wallets (not hardware)
    this.showEncryptOption = !data.temporary && (!!(coin && coin.serverWallets && data.filename) || (!data.isHardware && !(coin && coin.serverWallets)));
    this.isEncrypted = !!data.encrypted || !!data.encryptedSeed;
    this.showScanOption = walletService.canFindUsedAddresses(data);
  }

  onUnlockWallet() {
    this.closePopup(WalletOptionsResponses.UnlockWallet);
  }

  onAddNewAddress() {
    this.closePopup(WalletOptionsResponses.AddNewAddress);
  }

  onScanAddresses() {
    this.closePopup(WalletOptionsResponses.ScanAddresses);
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

  closePopup(response: WalletOptionsResponses | null = null) {
    this.dialogRef.close(response);
  }
}
