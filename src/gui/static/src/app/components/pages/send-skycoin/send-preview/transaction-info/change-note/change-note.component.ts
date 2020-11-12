import { Component, OnInit, Inject, ViewChild, OnDestroy } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { SubscriptionLike } from 'rxjs';

import { ButtonComponent } from '../../../../../layout/button/button.component';
import { StorageService, StorageType } from '../../../../../../services/storage.service';
import { MsgBarService } from '../../../../../../services/msg-bar.service';
import { AppConfig } from '../../../../../../app.config';
import { OldTransaction } from '../../../../../../services/wallet-operations/transaction-objects';

/**
 * Modal window used for changing the note of a transaction from the history. If the note is
 * changed, the modal window is closed and new note is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-change-note',
  templateUrl: './change-note.component.html',
  styleUrls: ['./change-note.component.scss'],
})
export class ChangeNoteComponent implements OnInit, OnDestroy {
  // Max chars the note can have.
  public static readonly MAX_NOTE_CHARS = 64;

  @ViewChild('button', { static: false }) button: ButtonComponent;
  form: FormGroup;
  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;
  // Deactivates the form while the system is busy.
  busy = false;

  private OperationSubscription: SubscriptionLike;
  private originalNote: string;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, transaction: OldTransaction): MatDialogRef<ChangeNoteComponent, any> {
    const config = new MatDialogConfig();
    config.data = transaction;
    config.autoFocus = true;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(ChangeNoteComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<ChangeNoteComponent>,
    @Inject(MAT_DIALOG_DATA) private data: OldTransaction,
    private formBuilder: FormBuilder,
    private msgBarService: MsgBarService,
    private storageService: StorageService,
  ) {}

  ngOnInit() {
    this.originalNote = this.data.note ? this.data.note : '';

    this.form = this.formBuilder.group({
      note: [this.data.note],
    });
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    if (this.OperationSubscription) {
      this.OperationSubscription.unsubscribe();
    }
  }

  closePopup() {
    this.dialogRef.close();
  }

  // Changes the note.
  changeNote() {
    if (this.busy) {
      return;
    }

    const newNote = this.form.value.note ? this.form.value.note.trim() : '';

    // If the new note is equal to the original one, do nothing.
    if (this.originalNote === newNote) {
      this.closePopup();

      return;
    }

    this.busy = true;
    this.msgBarService.hide();
    this.button.setLoading();

    // Change the note.
    this.OperationSubscription = this.storageService.store(StorageType.NOTES, this.data.id, newNote).subscribe(() => {
      this.busy = false;
      this.dialogRef.close(newNote);
    }, error => {
      this.busy = false;
      this.msgBarService.showError(error);
      this.button.resetState().setEnabled();
    });
  }
}
