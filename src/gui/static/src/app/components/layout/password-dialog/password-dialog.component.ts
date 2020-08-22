import { Component, Inject, OnDestroy, OnInit, ViewChild, ChangeDetectorRef } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { MatDialogRef } from '@angular/material/dialog';
import { FormControl, FormGroup } from '@angular/forms';
import { Subject } from 'rxjs';

import { ButtonComponent } from '../button/button.component';
import { processServiceError } from '../../../utils/errors';
import { MsgBarService } from '../../../services/msg-bar.service';
import { AppConfig } from '../../../app.config';
import { OperationError } from './../../../utils/operation-error';
import { WalletBase } from '../../../services/wallet-operations/wallet-objects';

export interface PasswordSubmitEvent {
  /**
   * Password entered by the user.
   */
  password: string;
  /**
   * Function for closing the modal window after completing the operation.
   */
  close(): void;
  /**
   * Function for informing the modal window about an error while completing the operation.
   */
  error(error: OperationError): void;
}

/**
 * Settings for PasswordDialogComponent.
 */
export interface PasswordDialogParams {
  /**
   * If true, the user will have to confirm the password in a second field.
   */
  confirm?: boolean;
  /**
   * Optional felp text.
   */
  description?: string;
  /**
   * Optional warning text.
   */
  warning?: boolean;
  /**
   * Custom title for the modal window.
   */
  title?: string;
  /**
   * Wallet to which the resquested password corresponds.
   */
  wallet: WalletBase;
}

/**
 * Modal window for requesting the password of a wallet. After the user enters the password,
 * it sends the passwordSubmit event, with the password, to let the code which openned this
 * modal window to continue with the operation. After finishing the operation, the code
 * must use the object returned by the passwordSubmit event to close the modal window
 * or for informing about an error.
 */
@Component({
  selector: 'app-password-dialog',
  templateUrl: './password-dialog.component.html',
  styleUrls: ['./password-dialog.component.scss'],
})
export class PasswordDialogComponent implements OnInit, OnDestroy {
  @ViewChild('button', { static: false }) button: ButtonComponent;
  form: FormGroup;
  passwordSubmit = new Subject<PasswordSubmitEvent>();
  working = false;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, params: PasswordDialogParams, smallSize = true): MatDialogRef<PasswordDialogComponent, any> {
    const config = new MatDialogConfig();
    config.data = params;
    config.autoFocus = true;
    config.width = smallSize ? '260px' : AppConfig.mediumModalWidth;

    return dialog.open(PasswordDialogComponent, config);
  }

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: PasswordDialogParams,
    public dialogRef: MatDialogRef<PasswordDialogComponent>,
    private msgBarService: MsgBarService,
    private changeDetector: ChangeDetectorRef,
  ) {
    // Se default values.
    this.data = Object.assign({
      confirm: false,
      description: null,
      warning: false,
      title: null,
      wallet: null,
    }, data || {});
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

    // Make the window bigger if a help msg is going to be shown.
    if (this.data.description) {
      this.dialogRef.updateSize('400px');
    }
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    this.passwordSubmit.complete();
  }

  /**
   * Deactivates the UI and sends the password.
   */
  proceed() {
    if (this.working || !this.form.valid) {
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
    if (this.form) {
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

  /**
   * Reactivates the UI and shows an error.
   */
  private error(error: OperationError) {
    if (!error.type) {
      error = processServiceError(error);
    }

    error.translatableErrorMsg = error.translatableErrorMsg ? error.translatableErrorMsg : 'password.decrypting-error';

    this.msgBarService.showError(error);
    this.button.resetState();
    this.working = false;
  }
}
