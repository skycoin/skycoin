import { Component, OnDestroy, ViewChild, Inject } from '@angular/core';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ButtonComponent } from '../../button/button.component';

enum States {
  Initial,
  Processing,
  ReturnedSuccess,
  ReturnedRefused,
  Failed,
}

@Component({
  selector: 'app-hw-seed-dialog',
  templateUrl: './hw-seed-dialog.html',
  styleUrls: ['./hw-seed-dialog.scss'],
})
export class HwSeedDialogComponent implements OnDestroy {

  @ViewChild('button') button: ButtonComponent;

  form: FormGroup;
  currentState: States = States.Initial;
  states = States;

  private operationSubscription: ISubscription;
  private hwConnectionSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public requestRecheck: any,
    public dialogRef: MatDialogRef<HwSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    private formBuilder: FormBuilder,
  ) {
    this.form = this.formBuilder.group({
      seed: ['cloud flower upset remain green metal below cup stem infant art thank', Validators.required],
    });

    this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
      if (!connected) {
        this.dialogRef.close();
      }
    });
  }

  closePopup() {
    this.dialogRef.close();
  }

  setSeed() {
    this.currentState = States.Processing;
    this.operationSubscription = this.hwWalletService.setMnemonic(this.form.value.seed).subscribe(
      response => {
        if (response.success) {
          this.requestRecheck();
          this.currentState = States.ReturnedSuccess;
        } else {
          this.currentState = States.ReturnedRefused;
        }
      },
      () => {
        this.currentState = States.Failed;
      },
    );
  }

  ngOnDestroy() {
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
    this.hwConnectionSubscription.unsubscribe();
  }
}
