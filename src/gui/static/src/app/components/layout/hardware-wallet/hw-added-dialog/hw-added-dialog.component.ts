import { Component, Inject, OnDestroy, ViewChild, ElementRef } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog } from '@angular/material/dialog';
import { UntypedFormGroup, UntypedFormBuilder, Validators } from '@angular/forms';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ChildHwDialogParams } from '../hw-options-dialog/hw-options-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ChangeNameComponent, ChangeNameData } from '../../../pages/wallets/change-name/change-name.component';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { OperationError } from '../../../../utils/operation-error';
import { processServiceError } from '../../../../utils/errors';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';

/**
 * Modal window used to add a new device to the wallet list. This modal window was created
 * for being oppenend by the hw wallet options modal window.
 */
@Component({
  selector: 'app-hw-added-dialog',
  templateUrl: './hw-added-dialog.component.html',
  styleUrls: ['./hw-added-dialog.component.scss'],
})
export class HwAddedDialogComponent extends HwDialogBaseComponent<HwAddedDialogComponent> implements OnDestroy {
  @ViewChild('input') input: ElementRef;
  wallet: WalletBase;
  form: UntypedFormGroup;
  maxHwWalletLabelLength = HwWalletService.maxLabelLength;

  // Vars with the validation error messages.
  inputErrorMsg = '';

  // Saves the initial label of the device, to know if the user tried to change it.
  private initialLabel: string;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data: ChildHwDialogParams,
    public dialogRef: MatDialogRef<HwAddedDialogComponent>,
    hwWalletService: HwWalletService,
    private formBuilder: UntypedFormBuilder,
    private dialog: MatDialog,
    private msgBarService: MsgBarService,
    private walletsAndAddressesService: WalletsAndAddressesService,
    private hardwareWalletService: HardwareWalletService,
  ) {
    super(hwWalletService, dialogRef);

    // Add the device to the wallets list.
    this.operationSubscription = this.walletsAndAddressesService.createHardwareWallet().subscribe(wallet => {
      // Update the security warnings.
      this.operationSubscription = this.hardwareWalletService.getFeaturesAndUpdateData(wallet).subscribe(() => {
        this.wallet = wallet;
        this.initialLabel = wallet.label;

        this.form = this.formBuilder.group({
          label: [wallet.label],
        });

        this.form.setValidators(this.validateForm.bind(this));

        this.currentState = this.states.Finished;

        // Request the data and state of the hw wallet options modal window to be refreshed.
        this.data.requestOptionsComponentRefresh();

        setTimeout(() => this.input.nativeElement.focus());
      }, err => this.processError(err));
    }, err => this.processError(err));
  }

  private processError(err: OperationError) {
    err = processServiceError(err);
    this.processHwOperationError(err);

    // Make the hw wallet options modal window show the error msg.
    this.data.requestOptionsComponentRefresh(err.translatableErrorMsg);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.msgBarService.hide();
  }

  saveNameAndCloseModal() {
    if (this.form.value.label === this.initialLabel) {
      // If no change was made to the label, just close the window.
      this.closeModal();
    } else {
      this.msgBarService.hide();

      // Open the appropiate component to change the device label.
      const data = new ChangeNameData();
      data.wallet = this.wallet;
      data.newName = this.form.value.label;
      ChangeNameComponent.openDialog(this.dialog, data, true).afterClosed().subscribe(result => {
        if (result && !result.errorMsg) {
          this.closeModal();
        } else if (result && result.errorMsg) {
          this.msgBarService.showError(result.errorMsg);
        }
      });
    }
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.inputErrorMsg = '';

    let valid = true;

    if (!this.form.get('label').value) {
      valid = false;
      if (this.form.get('label').touched) {
        this.inputErrorMsg = 'hardware-wallet.added.added-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
