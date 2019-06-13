import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs/Observable';
import { Bip39WordListService } from '../../../../services/bip39-word-list.service';
import { TranslateService } from '@ngx-translate/core';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';
import { ISubscription } from 'rxjs/Subscription';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
  selector: 'app-hw-seed-word-dialog',
  templateUrl: './hw-seed-word-dialog.component.html',
  styleUrls: ['./hw-seed-word-dialog.component.scss'],
})
export class HwSeedWordDialogComponent extends HwDialogBaseComponent<HwSeedWordDialogComponent> implements OnInit, OnDestroy {
  form: FormGroup;
  filteredOptions: Observable<string[]>;

  private sendingWord = false;
  protected valueChangeSubscription: ISubscription;

  constructor(
    public dialogRef: MatDialogRef<HwSeedWordDialogComponent>,
    private formBuilder: FormBuilder,
    private bip38WordList: Bip39WordListService,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      word: ['', Validators.required],
    });

    this.valueChangeSubscription = this.form.controls.word.valueChanges.subscribe(value => {
      this.bip38WordList.setSearchTerm(value.trim().toLowerCase());
    });

    this.filteredOptions = this.bip38WordList.searchResults.map(value => value);
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.msgBarService.hide();
    this.valueChangeSubscription.unsubscribe();
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
