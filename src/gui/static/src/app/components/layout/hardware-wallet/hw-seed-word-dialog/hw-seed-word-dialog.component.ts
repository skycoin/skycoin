import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs/Observable';
import { MessageIcons } from '../hw-message/hw-message.component';
import { Bip39WordListService } from '../../../../services/bip39-word-list.service';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-hw-seed-word-dialog',
  templateUrl: './hw-seed-word-dialog.component.html',
  styleUrls: ['./hw-seed-word-dialog.component.scss'],
})
export class HwSeedWordDialogComponent implements OnInit, OnDestroy {
  msgIcons = MessageIcons;
  form: FormGroup;
  filteredOptions: Observable<string[]>;

  constructor(
    public dialogRef: MatDialogRef<HwSeedWordDialogComponent>,
    private formBuilder: FormBuilder,
    private bip38WordList: Bip39WordListService,
    private snackbar: MatSnackBar,
    private translateService: TranslateService,
  ) {}

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
    this.snackbar.dismiss();
  }

  sendWord() {
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
    }, 32);
  }

  private validateWord(word: string): boolean {
    if (!this.bip38WordList.wordList.includes(word)) {
      return false;
    }

    return true;
  }
}
