import { Component, OnInit, Inject, ViewChild, OnDestroy } from '@angular/core';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ISubscription } from 'rxjs/Subscription';
import { ButtonComponent } from '../../../../../layout/button/button.component';
import { Transaction } from '../../../../../../app.datatypes';
import { StorageService, StorageType } from '../../../../../../services/storage.service';
import { MsgBarService } from '../../../../../../services/msg-bar.service';

@Component({
  selector: 'app-change-note',
  templateUrl: './change-note.component.html',
  styleUrls: ['./change-note.component.scss'],
})
export class ChangeNoteComponent implements OnInit, OnDestroy {

  public static readonly MAX_NOTE_CHARS = 64;

  @ViewChild('button', { static: false }) button: ButtonComponent;
  form: FormGroup;
  maxNoteChars = ChangeNoteComponent.MAX_NOTE_CHARS;

  private OperationSubscription: ISubscription;
  private originalNote: string;

  constructor(
    public dialogRef: MatDialogRef<ChangeNoteComponent>,
    @Inject(MAT_DIALOG_DATA) private data: Transaction,
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

  changeNote() {
    if (this.button.isLoading()) {
      return;
    }

    const newNote = this.form.value.note ? this.form.value.note.trim() : '';

    if (this.originalNote === newNote) {
      this.closePopup();

      return;
    }

    this.msgBarService.hide();
    this.button.setLoading();

    this.OperationSubscription = this.storageService.store(StorageType.NOTES, this.data.txid, newNote).subscribe(() => {
      this.dialogRef.close(newNote);
    }, error => {
      this.msgBarService.showError(error);
      this.button.resetState().setEnabled();
    });
  }
}
