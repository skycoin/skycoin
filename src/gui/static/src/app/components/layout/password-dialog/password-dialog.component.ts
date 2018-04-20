import { Component, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ButtonComponent } from '../button/button.component';
import { Observable } from 'rxjs/Observable';

@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.css']
})
export class PasswordDialogComponent implements OnInit, OnDestroy {

  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  passwordSubmit: Observable<any>;
  private passwordChanged;

  constructor(
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,
  ) {
    this.passwordSubmit = Observable.create(observer => {
      this.passwordChanged = password => {
        observer.next({
          password,
          close: this.close.bind(this),
          error: this.error.bind(this),
        });
      };
    });
  }

  ngOnInit() {
    this.form = new FormGroup({}, this.validateForm.bind(this));
    this.form.addControl('password', new FormControl(''));
    this.form.addControl('confirm_password', new FormControl(''));

    ['password', 'confirm_password'].forEach(control => {
      this.form.get(control).valueChanges.subscribe(() => {
        if (this.button.state === 2) {
          this.button.resetState();
        }
      });
    });

    if (this.requiresConfirmation) {
      this.form.get('confirm_password').enable();
    } else {
      this.form.get('confirm_password').disable();
    }
  }

  ngOnDestroy() {
    this.form.get('password').setValue('');
    this.form.get('confirm_password').setValue('');
  }

  proceed() {
    this.button.setLoading();
    this.passwordChanged(this.form.get('password').value);
  }

  get requiresConfirmation() {
    return this.data && this.data.confirm === true;
  }

  private validateForm() {
    if (this.form && this.form.get('password') && this.form.get('confirm_password')) {
      if (this.form.get('password').value.length === 0) {
        return { Required: true };
      }

      if (this.requiresConfirmation && this.form.get('password').value !== this.form.get('confirm_password').value) {
        return { NotEqual: true };
      }
    }

    return null;
  }

  private close() {
    this.dialogRef.close();
  }

  private error(error: any) {
    this.button.setError(error ? error : 'Incorrect password');
  }
}
