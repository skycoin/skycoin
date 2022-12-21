import { Component, Inject } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA } from '@angular/material/legacy-dialog';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';

/**
 * Allows to wipe the device and to remove the wallet from the list. This modal window was
 * created for being oppenend by the hw wallet options modal window.
 */
@Component({
  selector: 'app-hw-wipe-dialog',
  templateUrl: './hw-wipe-dialog.component.html',
  styleUrls: ['./hw-wipe-dialog.component.scss'],
})
export class HwWipeDialogComponent extends HwDialogBaseComponent<HwWipeDialogComponent> {
  showDeleteFromList = true;
  deleteFromList = true;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwWipeDialogComponent>,
    private hwWalletService: HwWalletService,
    private walletsAndAddressesService: WalletsAndAddressesService,
  ) {
    super(hwWalletService, dialogRef);

    // If no wallet was send as part of the data, the option for removing the wallet from
    // the wallet list is not shown and no wallet is removed, as there is no way to know which
    // wallet to remove. This is for wiping unknown devices, mainly ones blocked with a PIN code.
    if (!data.wallet) {
      this.showDeleteFromList = false;
      this.deleteFromList = false;
    }
  }

  // Changes the option which indicates if the wallet must be removed from the wallet list.
  setDeleteFromList(event) {
    this.deleteFromList = event.checked;
  }

  // Wipes the device.
  requestWipe() {
    this.currentState = this.states.Processing;

    this.operationSubscription = this.hwWalletService.wipe().subscribe(
      () => {
        this.showResult({
          text: 'hardware-wallet.general.completed',
          icon: this.msgIcons.Success,
        });

        // Request the data and state of the hw wallet options modal window to be refreshed.
        this.data.requestOptionsComponentRefresh();

        // Remove the wallet from the list, if requested.
        if (this.deleteFromList) {
          this.walletsAndAddressesService.deleteHardwareWallet(this.data.wallet.id);
        }
      },
      err => this.processHwOperationError(err),
    );
  }
}
