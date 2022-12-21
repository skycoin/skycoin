import { Component, OnInit, HostListener, Inject } from '@angular/core';
import { UntypedFormBuilder, Validators, UntypedFormGroup } from '@angular/forms';
import { MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig, MAT_LEGACY_DIALOG_DATA as MAT_DIALOG_DATA } from '@angular/material/legacy-dialog';

import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwPinHelpDialogComponent } from '../hw-pin-help-dialog/hw-pin-help-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ChangePinStates } from '../../../../services/hw-wallet-pin.service';

/**
 * Settings for HwPinDialogComponent.
 */
export interface HwPinDialogParams {
  /**
   * If the PIN code is being requested for signing a tx. Ignored if changingPin is true.
   */
  signingTx: boolean;
  /**
   * If the PIN code is being requested for setting or changing the PIN on the device.
   */
  changingPin: boolean;
  /**
   * State of the PIN changing operation if changingPin is true.
   */
  changePinState: ChangePinStates;
}

/**
 * Allows the user to enter the PIN code. If the user completes the operation, the modal window
 * is closed and the positions selected by the user on the PIN matrix are returned in the
 * "afterClosed" event.
 */
@Component({
  selector: 'app-hw-pin-dialog',
  templateUrl: './hw-pin-dialog.component.html',
  styleUrls: ['./hw-pin-dialog.component.scss'],
})
export class HwPinDialogComponent extends HwDialogBaseComponent<HwPinDialogComponent> implements OnInit {
  form: UntypedFormGroup;
  changePinStates = ChangePinStates;
  buttonsContent = '•';

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, params: HwPinDialogParams): MatDialogRef<HwPinDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = params;
    config.autoFocus = false;
    config.width = '350px';

    return dialog.open(HwPinDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: HwPinDialogParams,
    public dialogRef: MatDialogRef<HwPinDialogComponent>,
    private formBuilder: UntypedFormBuilder,
    private dialog: MatDialog,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      pin: ['', Validators.compose([Validators.required, Validators.minLength(4)])],
    });
  }

  get title(): string {
    if (!this.data.changingPin) {
      return 'hardware-wallet.enter-pin.title';
    } else if (this.data.changePinState === ChangePinStates.RequestingNewPin) {
      return 'hardware-wallet.enter-pin.title-change-new';
    } else if (this.data.changePinState === ChangePinStates.ConfirmingNewPin) {
      return 'hardware-wallet.enter-pin.title-change-confirm';
    } else {
      return 'hardware-wallet.enter-pin.title-change-current';
    }
  }

  openHelp() {
    HwPinHelpDialogComponent.openDialog(this.dialog);
  }

  /**
   * Allow to enter the PIN using the numeric keys to emulate the PIN matrix.
   */
  @HostListener('window:keyup', ['$event'])
  keyEvent(event: KeyboardEvent) {
    const key = parseInt(event.key, 10);
    if (key > 0 && key < 10) {
      this.addNumber(key.toString());
    } else if (event.keyCode === 8) {
      this.removeNumber();
    } else if (event.keyCode === 13) {
      this.sendPin();
    }
  }

  /**
   * Add a new number to the PIN.
   * @param number Position of the number of the new number on the matrix.
   */
  addNumber(number: string) {
    const currentValue: string = this.form.value.pin;
    if (currentValue.length < 8) {
      this.form.get('pin').setValue(currentValue + number);
    }
  }

  /**
   * Removes the last number from the PIN.
   */
  removeNumber() {
    const currentValue: string = this.form.value.pin;
    this.form.get('pin').setValue(currentValue.substring(0, currentValue.length - 1));
  }

  sendPin() {
    if (this.form.valid) {
      this.dialogRef.close(this.form.value.pin);
    }
  }
}
