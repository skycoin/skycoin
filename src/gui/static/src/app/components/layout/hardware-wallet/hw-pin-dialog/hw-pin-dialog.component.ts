import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { HwWalletService } from '../../../../services/hw-wallet.service';

@Component({
  selector: 'app-hw-pin-dialog',
  templateUrl: './hw-pin-dialog.component.html',
  styleUrls: ['./hw-pin-dialog.component.scss'],
})
export class HwPinDialogComponent implements OnInit, OnDestroy {
  form: FormGroup;

  private hwConnectionSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwPinDialogComponent>,
    private formBuilder: FormBuilder,
    private hwWalletService: HwWalletService,
  ) {}

  ngOnInit() {
    this.form = this.formBuilder.group({
      pin: ['', Validators.required],
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

  closeModal() {
    this.dialogRef.close();
  }

  sendPin() {
    this.dialogRef.close(this.form.value.pin);
  }
}
