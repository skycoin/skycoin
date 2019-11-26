import { Component, Inject, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material';
import { FormControl, FormGroup } from '@angular/forms';
import { ButtonComponent } from '../button/button.component';
import { parseResponseMessage } from '../../../utils/errors';
import { Subject } from 'rxjs/Subject';
import { ISubscription } from 'rxjs/Subscription';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../services/msg-bar.service';

@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.scss'],
})
export class PasswordDialogComponent implements OnInit, OnDestroy {
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  passwordSubmit = new Subject<any>();
  working = false;

  private subscriptions: ISubscription[] = [];
  private errors: any;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    private msgBarService: MsgBarService,
    private translateService: TranslateService,
  ) {
    this.data = Object.assign({
      confirm: false,
      description: null,
      warning: false,
      title: null,
      wallet: null,
    }, data || {});

    this.translateService.get(['errors.incorrect-password', 'errors.api-disabled', 'errors.no-wallet']).subscribe(res => {
      this.errors = res;
    });
  }

  ngOnInit() {
    this.form = new FormGroup({}, this.validateForm.bind(this));
    this.form.addControl('password', new FormControl(''));
    this.form.addControl('confirm_password', new FormControl(''));

    ['password', 'confirm_password'].forEach(control => {
      this.subscriptions.push(this.form.get(control).valueChanges.subscribe(() => {
        if (this.button.state === 2) {
          this.button.resetState();
        }
      }));
    });

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

    this.subscriptions.forEach(sub => sub.unsubscribe());
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
          error = this.errors['errors.incorrect-password'];
          break;
        case 403:
          error = this.errors['errors.api-disabled'];
          break;
        case 404:
          error = this.errors['errors.no-wallet'];
          break;
        default:
            error = this.errors['errors.error-decrypting'];
        }
      } else {
        error = this.errors['errors.error-decrypting'];
      }
    }

    error = error ? error : this.errors['errors.error-decrypting'];

    this.msgBarService.showError(error);
    this.button.resetState();
    this.working = false;
  }
}
