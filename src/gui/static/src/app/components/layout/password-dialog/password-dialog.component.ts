import { Component, Inject, OnDestroy, OnInit, ViewChild, ChangeDetectorRef } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { FormControl, FormGroup } from '@angular/forms';
import { ButtonComponent } from '../button/button.component';
import { parseResponseMessage } from '../../../utils/errors';
import { Subject } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../services/msg-bar.service';
import { AppConfig } from '../../../app.config';
import { Wallet } from '../../../app.datatypes';

export interface PasswordDialogParams {
  confirm?: boolean;
  description?: string;
  warning?: boolean;
  title?: string;
  wallet: Wallet;
}

@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.scss'],
})
export class PasswordDialogComponent implements OnInit, OnDestroy {
  @ViewChild('button', { static: false }) button: ButtonComponent;
  form: FormGroup;
  passwordSubmit = new Subject<any>();
  working = false;

  private errors: any;

  public static openDialog(dialog: MatDialog, params: PasswordDialogParams, smallSize = true): MatDialogRef<PasswordDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = params;
    config.autoFocus = true;
    config.width = smallSize ? '260px' : AppConfig.mediumModalWidth;

    return dialog.open(PasswordDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
    private changeDetector: ChangeDetectorRef,
  ) {
    this.data = Object.assign({
      confirm: false,
      description: null,
      warning: false,
      title: null,
      wallet: null,
    }, data || {});

    this.translateService.get(['password.incorrect-password-error', 'password.api-disabled-error', 'password.no-wallet-error', 'password.decrypting-error']).subscribe(res => {
      this.errors = res;
    });
  }

  ngOnInit() {
    this.form = new FormGroup({}, this.validateForm.bind(this));
    this.form.addControl('password', new FormControl(''));
    this.form.addControl('confirm_password', new FormControl(''));

    if (this.data.confirm) {
      this.form.get('confirm_password').enable();
    } else {
      this.form.get('confirm_password').disable();
    }

    if (this.data.description) {
      this.dialogRef.updateSize('400px');
    }
  }

  ngOnDestroy() {
    this.msgBarService.hide();

    this.form.get('password').setValue('');
    this.form.get('confirm_password').setValue('');

    this.passwordSubmit.complete();
  }

  proceed() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.msgBarService.hide();

    this.button.setLoading();
    this.working = true;

    this.passwordSubmit.next({
      password: this.form.get('password').value,
      close: this.close.bind(this),
      error: this.error.bind(this),
    });

    this.changeDetector.detectChanges();
  }

  private validateForm() {
    if (this.form && this.form.get('password') && this.form.get('confirm_password')) {
      if (this.form.get('password').value.length === 0) {
        return { Required: true };
      }

      if (this.data.confirm && this.form.get('password').value !== this.form.get('confirm_password').value) {
        return { NotEqual: true };
      }
    }

    return null;
  }

  private close() {
    this.dialogRef.close();
  }

  private error(error: any) {
    if (typeof error === 'object') {
      if (error.status) {
      switch (error.status) {
        case 400:
          error = parseResponseMessage(error['_body']);
          break;
        case 401:
          error = this.errors['password.incorrect-password-error'];
          break;
        case 403:
          error = this.errors['password.api-disabled-error'];
          break;
        case 404:
          error = this.errors['password.no-wallet-error'];
          break;
        default:
          error = this.errors['password.decrypting-error'];
        }
      } else {
        error = this.errors['password.decrypting-error'];
      }
    }

    error = error ? error : this.errors['password.decrypting-error'];

    this.msgBarService.showError(error);
    this.button.resetState();
    this.working = false;
  }
}
