import { Component, EventEmitter, Inject, OnInit, Output, ViewChild, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup, Validators, FormControl } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Subscription } from 'rxjs';
import { Observable } from 'rxjs';

import { Wallet } from '../../../../app.datatypes';
import { WalletService } from '../../../../services/wallet/wallet.service';
import { config } from '../../../../app.config';
import { MsgBarService } from '../../../../services/msg-bar.service';

export class ConfirmSeedParams {
  wallet: Wallet;
}

@Component({
  selector: 'app-unlock-wallet',
  templateUrl: './unlock-wallet.component.html',
  styleUrls: ['./unlock-wallet.component.scss'],
})
export class UnlockWalletComponent implements OnInit, OnDestroy {
  @Output() onWalletUnlocked = new EventEmitter<void>();
  @Output() onDeleteClicked = new EventEmitter<void>();
  @ViewChild('unlock') unlockButton;
  form: FormGroup;
  disableDismiss = false;
  loadingProgress = 0;
  showConfirmSeedWarning;
  showSlowMobileInfo = false;

  private wallet: Wallet;
  private unlockSubscription: Subscription;
  private progressSubscription: Subscription;
  private slowInfoSubscription: Subscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) private data,
    public dialogRef: MatDialogRef<UnlockWalletComponent>,
    private formBuilder: FormBuilder,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
  ) {
    if (data.wallet) {
      this.showConfirmSeedWarning = true;
      this.wallet = data.wallet;
    } else {
      this.showConfirmSeedWarning = false;
      this.wallet = data;
    }
  }

  ngOnInit() {
    this.initForm();
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    this.removeProgressSubscriptions();
    this.removeSlowInfoSubscription();
  }

  closePopup() {
    this.dialogRef.close();
  }

  unlockWallet() {
    this.removeProgressSubscriptions();

    this.unlockButton.setLoading();
    this.disableDismiss = true;

    this.createSlowInfoSubscription();

    const onProgressChanged = new EventEmitter<number>();
    if (this.wallet.addresses.length > 1) {
      this.progressSubscription = onProgressChanged.subscribe((progress) => {
        this.createSlowInfoSubscription();
        this.loadingProgress = progress;
      });
    }

    this.unlockSubscription = this.walletService.unlockWallet(this.wallet, this.form.value.seed, onProgressChanged)
      .subscribe(
        () => this.onUnlockSuccess(),
        (error: Error) => this.onUnlockError(error)
      );
  }

  delete() {
    this.closePopup();
    this.onDeleteClicked.emit();
  }

  private initForm() {
    this.form = this.formBuilder.group({
      seed: ['', Validators.required],
    });
  }

  private onUnlockSuccess() {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.removeProgressSubscriptions();
    this.unlockButton.setSuccess();
    this.closePopup();
    this.onWalletUnlocked.emit();
  }

  private onUnlockError(error: Error) {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.removeProgressSubscriptions();
    this.disableDismiss = false;
    this.msgBarService.showError(error.message);
    this.unlockButton.resetState();
  }

  private removeProgressSubscriptions() {
    if (this.progressSubscription && !this.progressSubscription.closed) {
      this.progressSubscription.unsubscribe();
    }
    if (this.unlockSubscription && !this.unlockSubscription.closed) {
      this.unlockSubscription.unsubscribe();
    }
  }

  private createSlowInfoSubscription() {
    this.removeSlowInfoSubscription();

    this.slowInfoSubscription = Observable.of(1).delay(config.timeBeforeSlowMobileInfo)
      .subscribe(() => this.showSlowMobileInfo = true);
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
