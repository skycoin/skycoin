import { Component, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef, MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ButtonComponent } from '../button/button.component';
import { Observable } from 'rxjs/Observable';
import { parseResponseMessage } from '../../../utils/index';

@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.css']
})
export class PasswordDialogComponent implements OnInit, OnDestroy {

  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  passwordSubmit: Observable<any>;
  disableDismiss = false;
  private passwordChanged;

  constructor(
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    @Inject(MAT_DIALOG_DATA) private data: any,
    private snackbar: MatSnackBar,
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
    this.form = new FormGroup({});
    this.form.addControl('password', new FormControl('', [Validators.required]));
    this.form.get('password').valueChanges.subscribe(() => {
      if (this.button.state === 2) {
        this.button.resetState();
      }
    });
  }

  ngOnDestroy() {
    this.form.get('password').setValue('');
  }

  proceed() {
    this.button.setLoading();
    this.passwordChanged(this.form.get('password').value);
    this.disableDismiss = true;
  }

  private close() {
    this.dialogRef.close();
  }

  private error(error: any) {
    if (typeof error === 'object') {
      switch (error.status) {
        case 400:
          error = parseResponseMessage(error['_body']);
          break;
        case 401:
          error = 'Incorrect password';
          break;
        case 403:
          error = 'API Disabled';
          break;
        case 404:
          error = 'Wallet does not exist';
          break;
        default:
          const config = new MatSnackBarConfig();
          config.duration = 5000;
          this.snackbar.open(parseResponseMessage(error['_body']), null, config);
      }
    }

    this.button.setError(error ? error : 'Incorrect password');
    this.disableDismiss = false;
  }
}
