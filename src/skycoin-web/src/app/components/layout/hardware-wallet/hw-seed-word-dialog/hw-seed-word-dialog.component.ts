import { ChangeDetectionStrategy, Component } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup } from '@angular/forms';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';

import { Bip39WordListService } from '../../../../services/bip39-word-list.service';

/**
 * Asks for the seed word the hardware wallet shows it wants. Closes with the
 * word, or with undefined when cancelled.
 */
@Component({
    selector: 'app-hw-seed-word-dialog',
    templateUrl: './hw-seed-word-dialog.component.html',
    styleUrls: ['./hw-seed-word-dialog.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class HwSeedWordDialogComponent {
  form: UntypedFormGroup;
  error = '';

  static openDialog(dialog: MatDialog): MatDialogRef<HwSeedWordDialogComponent> {
    return dialog.open(HwSeedWordDialogComponent, { width: '450px', disableClose: true });
  }

  constructor(
    public dialogRef: MatDialogRef<HwSeedWordDialogComponent>,
    private bip39WordListService: Bip39WordListService,
    formBuilder: UntypedFormBuilder,
  ) {
    this.form = formBuilder.group({ word: [''] });
  }

  submit() {
    const word = (this.form.value.word || '').trim().toLowerCase();
    if (!word) {
      return;
    }
    if (this.bip39WordListService.validateWord(word) === false) {
      this.error = 'hardware-wallet.seed-word.invalid-word';

      return;
    }
    this.dialogRef.close(word);
  }

  cancel() {
    this.dialogRef.close();
  }
}
