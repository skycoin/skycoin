import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs/Observable';
import { Bip39WordListService } from '../../../../services/bip39-word-list.service';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { TranslateService } from '@ngx-translate/core';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { HwDialogBaseComponent } from '../hw-dialog-base.component';

@Component({
  selector: 'app-hw-seed-word-dialog',
  templateUrl: './hw-seed-word-dialog.component.html',
  styleUrls: ['./hw-seed-word-dialog.component.scss'],
})
export class HwSeedWordDialogComponent extends HwDialogBaseComponent<HwSeedWordDialogComponent> implements OnInit, OnDestroy {
  form: FormGroup;
  filteredOptions: Observable<string[]>;

  private sendingWord = false;

  constructor(
    public dialogRef: MatDialogRef<HwSeedWordDialogComponent>,
    private formBuilder: FormBuilder,
    private bip38WordList: Bip39WordListService,
    private snackbar: MatSnackBar,
    private translateService: TranslateService,
    hwWalletService: HwWalletService,
  ) {
    super(hwWalletService, dialogRef);
  }

  ngOnInit() {
    this.form = this.formBuilder.group({
      word: ['', Validators.required],
    });

    this.filteredOptions = this.form.controls.word.valueChanges.map(value => {
      if ((value as string).trim() !== '') {
        const filterValue = value.trim().toLowerCase();

        return this.bip38WordList.wordList.filter(option => option.startsWith(filterValue));
      } else {
        return [];
      }
    });
  }

  ngOnDestroy() {
    super.ngOnDestroy();
    this.snackbar.dismiss();
  }

  sendWord() {
    if (!this.sendingWord) {
      this.sendingWord = true;
      this.snackbar.dismiss();

      setTimeout(() => {
        if (this.form.valid) {
          if (this.validateWord(this.form.value.word)) {
            this.dialogRef.close((this.form.value.word as string).trim().toLowerCase());
          } else {
            const config = new MatSnackBarConfig();
            config.duration = 5000;
            this.snackbar.open(this.translateService.instant('hardware-wallet.seed-word.error-invalid-word'), null, config);
          }
        }
        this.sendingWord = false;
      }, 32);
    }
  }

  private validateWord(word: string): boolean {
    if (!this.bip38WordList.wordList.includes(word.trim().toLowerCase())) {
      return false;
    }

    return true;
  }
}
