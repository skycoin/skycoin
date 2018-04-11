import { Component, Inject, OnInit, ViewChild } from '@angular/core';
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
export class PasswordDialogComponent implements OnInit {

  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  passwordSubmit: Observable;
  private passwordChanged;

  constructor(
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    @Inject(MAT_DIALOG_DATA) private data: any,
  ) {
    this.passwordSubmit = Observable.create(observer => {
      this.passwordChanged = password => {
        observer.next({
          password,
          close: this.close,
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

  proceed() {
    this.passwordChanged(this.form.get('password').value);
  }

  private close() {
    this.dialogRef.close();
  }

  private error(error: any) {
    this.button.setError(error ? error : 'Incorret password');
  }
}
