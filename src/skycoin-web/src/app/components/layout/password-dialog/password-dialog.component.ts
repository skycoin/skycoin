import { ChangeDetectionStrategy, Component, Inject, OnInit } from '@angular/core';
import { AbstractControl, UntypedFormBuilder, UntypedFormGroup, ValidationErrors, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Observable } from 'rxjs';

export interface PasswordDialogData {
  /** Translation key or text for the headline. Defaults to password.title. */
  title?: string;
  /** Shown above the field, a translation key or text. */
  text?: string;
  /** Asks for the password twice, for setting a new one. */
  confirm?: boolean;
  /** Lets an empty password through, for an optional one. */
  allowEmpty?: boolean;
}

/**
 * Asks for a password in a masked field. Closes with the password, or with
 * undefined when cancelled. Use PasswordDialogComponent.ask.
 */
@Component({
    selector: 'app-password-dialog',
    templateUrl: './password-dialog.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class PasswordDialogComponent implements OnInit {
  form!: UntypedFormGroup;

  static ask(dialog: MatDialog, data: PasswordDialogData = {}): Observable<string | undefined> {
    return dialog.open(PasswordDialogComponent, { width: '450px', data }).afterClosed();
  }

  constructor(
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: PasswordDialogData,
    private formBuilder: UntypedFormBuilder,
  ) {}

  ngOnInit() {
    this.form = this.formBuilder.group({
      password: ['', this.data.allowEmpty ? [] : Validators.required],
      confirm: [''],
    }, { validators: this.data.confirm ? this.matches : [] });
  }

  get mismatch(): boolean {
    return !!this.form.errors?.['mismatch'] && this.form.get('confirm')!.dirty;
  }

  submit() {
    if (this.form.valid) {
      this.dialogRef.close(this.form.value.password);
    }
  }

  cancel() {
    this.dialogRef.close();
  }

  private matches(group: AbstractControl): ValidationErrors | null {
    return group.get('password')!.value === group.get('confirm')!.value ? null : { mismatch: true };
  }
}
