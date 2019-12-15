import { Component, OnInit, OnDestroy, Inject } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { Observable, SubscriptionLike } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { Bip39WordListService } from '../../../services/bip39-word-list.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { HwWalletService } from '../../../services/hw-wallet.service';
import { MessageIcons } from '../hardware-wallet/hw-message/hw-message.component';
import { map } from 'rxjs/operators';

export class SeedWordDialogParams {
  isForHwWallet: boolean;
  wordNumber: number;
  restoringSoftwareWallet: false;
}

@Component({
  selector: 'app-seed-word-dialog',
  templateUrl: './seed-word-dialog.component.html',
  styleUrls: ['./seed-word-dialog.component.scss'],
})
export class SeedWordDialogComponent implements OnInit, OnDestroy {
  form: FormGroup;
  filteredOptions: Observable<string[]>;
  msgIcons = MessageIcons;

  private sendingWord = false;
  private valueChangeSubscription: SubscriptionLike;
  private hwConnectionSubscription: SubscriptionLike;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: SeedWordDialogParams,
    public dialogRef: MatDialogRef<SeedWordDialogComponent>,
    private formBuilder: FormBuilder,
    private bip38WordList: Bip39WordListService,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
    hwWalletService: HwWalletService,
  ) {
    if (data.isForHwWallet) {
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

    this.valueChangeSubscription = this.form.controls.word.valueChanges.subscribe(value => {
      this.bip38WordList.setSearchTerm(value.trim().toLowerCase());
    });

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
