import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { Observable, SubscriptionLike } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { map } from 'rxjs/operators';

import { Bip39WordListService } from '../../../services/bip39-word-list.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { MessageIcons } from '../hardware-wallet/hw-message/hw-message.component';

/**
 * Reasons for asking a seed word.
 */
export enum WordAskedReasons {
  HWWalletOperation = 'HWWalletOperation',
  CreatingSoftwareWallet = 'CreatingSoftwareWallet',
  RecoveringSoftwareWallet = 'RecoveringSoftwareWallet',
}

/**
 * Settings for SeedWordDialogComponent.
 */
export interface SeedWordDialogParams {
  /**
   * Reason why the word is being requested.
   */
  reason: WordAskedReasons;
  /**
   * Number of the requested word, if it is not for a hw wallet.
   */
  wordNumber?: number;
}

/**
 * Modal window used to ask the user to enter a word of a seed. If used to request a word for a
 * hw wallet, the modal window is automatically closed if the device is disconnected. It only
 * allows BIP38 words. In the "afterClosed" event it returns the word the user entered, or
 * null, if the operation was cancelled.
 */
@Component({
  selector: 'app-seed-word-dialog',
  templateUrl: './seed-word-dialog.component.html',
  styleUrls: ['./seed-word-dialog.component.scss'],
})
export class SeedWordDialogComponent implements OnInit, OnDestroy {
  form: FormGroup;
  filteredOptions: Observable<string[]>;

  msgIcons = MessageIcons;
  wordAskedReasons = WordAskedReasons;

  private sendingWord = false;
  private valueChangeSubscription: SubscriptionLike;
  private hwConnectionSubscription: SubscriptionLike;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, params: SeedWordDialogParams): MatDialogRef<SeedWordDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = params;
    config.autoFocus = true;
    config.width = '350px';

    return dialog.open(SeedWordDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: SeedWordDialogParams,
    public dialogRef: MatDialogRef<SeedWordDialogComponent>,
    private formBuilder: FormBuilder,
    private bip38WordList: Bip39WordListService,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
    hwWalletService: HwWalletService,
  ) {
    if (data.reason === WordAskedReasons.HWWalletOperation) {
      // Close the window if the device is disconnected.
      this.hwConnectionSubscription = hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
        if (!connected) {
          this.dialogRef.close();
        }
      });
    }
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      word: ['', Validators.required],
    });

    // Search for sugestions when the user changes the content of the word field.
    this.valueChangeSubscription = this.form.controls.word.valueChanges.subscribe(value => {
      this.bip38WordList.setSearchTerm(value.trim().toLowerCase());
    });

    // Get the sugestions.
    this.filteredOptions = this.bip38WordList.searchResults.pipe(map(value => value));
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    this.valueChangeSubscription.unsubscribe();
    if (this.hwConnectionSubscription) {
      this.hwConnectionSubscription.unsubscribe();
    }
  }

  sendWord() {
    if (!this.sendingWord) {
      this.sendingWord = true;
      this.msgBarService.hide();

      // If the user selected a sugested word, the wait allows time for the word to be
      // added to the form field.
      setTimeout(() => {
        if (this.form.valid) {
          const validation = this.bip38WordList.validateWord(this.form.value.word.trim().toLowerCase());
          if (validation) {
            this.dialogRef.close((this.form.value.word as string).trim().toLowerCase());
          } else {
            if (validation === null) {
              this.msgBarService.showError(this.translateService.instant('hardware-wallet.seed-word.error-loading-words'));
            } else {
              this.msgBarService.showError(this.translateService.instant('hardware-wallet.seed-word.error-invalid-word'));
            }
          }
        }
        this.sendingWord = false;
      }, 32);
    }
  }
}
