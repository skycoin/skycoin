import { Component, OnInit, HostListener, Inject } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef, MatDialog, MatDialogConfig, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwPinHelpDialogComponent } from '../hw-pin-help-dialog/hw-pin-help-dialog.component';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ChangePinStates } from '../../../../services/hw-wallet-pin.service';

export interface HwPinDialogParams {
  signingTx: boolean;
  changingPin: boolean;
  changePinState: ChangePinStates;
}

@Component({
  selector: 'app-hw-pin-dialog',
  templateUrl: './hw-pin-dialog.component.html',
  styleUrls: ['./hw-pin-dialog.component.scss'],
})
export class HwPinDialogComponent extends HwDialogBaseComponent<HwPinDialogComponent> implements OnInit {
  form: FormGroup;
  changePinStates = ChangePinStates;
  buttonsContent = '•';

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: HwPinDialogParams,
    public dialogRef: MatDialogRef<HwPinDialogComponent>,
    private formBuilder: FormBuilder,
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
    this.dialog.open(HwPinHelpDialogComponent, <MatDialogConfig> {
      width: '450px',
    });
  }

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

  addNumber(number: string) {
    const currentValue: string = this.form.value.pin;
    if (currentValue.length < 8) {
      this.form.get('pin').setValue(currentValue + number);
    }
  }

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
