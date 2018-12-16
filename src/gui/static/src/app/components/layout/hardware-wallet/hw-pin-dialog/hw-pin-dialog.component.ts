import { Component, OnInit, OnDestroy, HostListener } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, ChangePinStates } from '../../../../services/hw-wallet.service';
import { HwPinHelpComponent } from '../hw-pin-help/hw-pin-help.component';

@Component({
  selector: 'app-hw-pin-dialog',
  templateUrl: './hw-pin-dialog.component.html',
  styleUrls: ['./hw-pin-dialog.component.scss'],
})
export class HwPinDialogComponent implements OnInit, OnDestroy {
  static showForSigningTx = false;
  static currentSignature = 1;
  static totalSignatures = 2;
  static showForChangingPin = false;
  static changePinState = ChangePinStates.RequestingCurrentPin;

  form: FormGroup;
  showForSigning: boolean;
  current: number;
  total: number;
  changingPin: boolean;
  changeState = ChangePinStates.RequestingCurrentPin;
  changePinStates = ChangePinStates;

  private hwConnectionSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwPinDialogComponent>,
    private formBuilder: FormBuilder,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
  ) {
    this.showForSigning = HwPinDialogComponent.showForSigningTx;
    this.current = HwPinDialogComponent.currentSignature;
    this.total = HwPinDialogComponent.totalSignatures;

    this.changingPin = HwPinDialogComponent.showForChangingPin;
    this.changeState = HwPinDialogComponent.changePinState;
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      pin: ['', Validators.compose([Validators.required, Validators.minLength(4)])],
    });

    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.dialogRef.close();
      }
    });
  }

  ngOnDestroy() {
    this.hwConnectionSubscription.unsubscribe();

    if (HwPinDialogComponent.showForChangingPin) {
      if (HwPinDialogComponent.changePinState === ChangePinStates.RequestingCurrentPin) {
        HwPinDialogComponent.changePinState = ChangePinStates.RequestingNewPin;
      } else if (HwPinDialogComponent.changePinState === ChangePinStates.RequestingNewPin) {
        HwPinDialogComponent.changePinState = ChangePinStates.ConfirmingNewPin;
      }
    }
  }

  openHelp() {
    this.dialog.open(HwPinHelpComponent, <MatDialogConfig> {
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
