import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';

export class AddressConfirmationParams {
  wallet: WalletBase;
  addressIndex: number;
  showCompleteConfirmation: boolean;
}

@Component({
  selector: 'app-hw-confirm-address-dialog',
  templateUrl: './hw-confirm-address-dialog.component.html',
  styleUrls: ['./hw-confirm-address-dialog.component.scss'],
})
export class HwConfirmAddressDialogComponent extends HwDialogBaseComponent<HwConfirmAddressDialogComponent> {
  constructor(
    @Inject(MAT_DIALOG_DATA) public data: AddressConfirmationParams,
    public dialogRef: MatDialogRef<HwConfirmAddressDialogComponent>,
    private hardwareWalletService: HardwareWalletService,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = this.hardwareWalletService.confirmAddress(data.wallet, data.addressIndex).subscribe(
      () => {
        this.showResult({
          text: data.showCompleteConfirmation ? 'hardware-wallet.confirm-address.confirmation' : 'hardware-wallet.confirm-address.short-confirmation',
          icon: this.msgIcons.Success,
        });
      },
      err => this.processResult(err),
    );
  }
}
