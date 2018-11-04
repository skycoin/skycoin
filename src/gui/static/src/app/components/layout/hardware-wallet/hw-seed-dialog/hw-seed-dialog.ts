import { Component, OnDestroy, ViewChild } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ButtonComponent } from '../../button/button.component';

export enum States {
  Initial,
  Processing,
  Finalized,
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

  constructor(
    public dialogRef: MatDialogRef<HwSeedDialogComponent>,
    private hwWalletService: HwWalletService,
    private formBuilder: FormBuilder,
  ) {
    this.form = this.formBuilder.group({
      seed: ['cloud flower upset remain green metal below cup stem infant art thank', Validators.required],
    });
  }

  closePopup() {
    this.dialogRef.close();
  }

  setSeed() {
    this.currentState = States.Processing;
    this.operationSubscription = this.hwWalletService.setMnemonic(this.form.value.seed).subscribe(
      () => {
        this.currentState = States.Finalized;
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
  }
}
