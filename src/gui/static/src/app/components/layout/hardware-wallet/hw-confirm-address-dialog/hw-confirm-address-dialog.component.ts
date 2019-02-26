import { Component, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService, OperationResults } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { Address } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet.service';

enum States {
  Initial,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

export class AddressConfirmationParams {
  address: Address;
  addressIndex: number;
  showCompleteConfirmation: boolean;
}

@Component({
  selector: 'app-hw-confirm-address-dialog',
  templateUrl: './hw-confirm-address-dialog.component.html',
  styleUrls: ['./hw-confirm-address-dialog.component.scss'],
})
export class HwConfirmAddressDialogComponent extends HwDialogBaseComponent<HwConfirmAddressDialogComponent> {

  currentState: States = States.Initial;
  states = States;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: AddressConfirmationParams,
    public dialogRef: MatDialogRef<HwConfirmAddressDialogComponent>,
    private hwWalletService: HwWalletService,
    private walletService: WalletService,
  ) {
    super(hwWalletService, dialogRef);
    this.operationSubscription = this.hwWalletService.confirmAddress(data.addressIndex).subscribe(
      () => {
        this.currentState = States.ReturnedSuccess;
        this.data.address.confirmed = true;
        this.walletService.saveHardwareWallets();
      },
      err => {
        if (err.result && err.result === OperationResults.FailedOrRefused) {
          this.currentState = States.ReturnedRefused;
        } else {
          this.currentState = States.Failed;
        }
      },
    );
  }
}
