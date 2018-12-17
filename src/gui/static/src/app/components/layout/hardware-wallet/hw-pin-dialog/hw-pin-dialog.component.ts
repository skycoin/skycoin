import { Component, OnInit, OnDestroy, HostListener, Inject } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef, MatDialog, MatDialogConfig, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService, ChangePinStates } from '../../../../services/hw-wallet.service';
import { HwPinHelpDialogComponent } from '../hw-pin-help-dialog/hw-pin-help-dialog.component';

export interface HwPinDialogParams {
  signingTx: boolean;
  currentSignature: number;
  totalSignatures: number;
  changingPin: boolean;
  changePinState: ChangePinStates;
}

@Component({
  selector: 'app-hw-pin-dialog',
  templateUrl: './hw-pin-dialog.component.html',
  styleUrls: ['./hw-pin-dialog.component.scss'],
})
export class HwPinDialogComponent implements OnInit, OnDestroy {
  form: FormGroup;
  changePinStates = ChangePinStates;

  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: HwPinDialogParams,
    public dialogRef: MatDialogRef<HwPinDialogComponent>,
    private formBuilder: FormBuilder,
    private hwWalletService: HwWalletService,
    private dialog: MatDialog,
  ) {}

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
