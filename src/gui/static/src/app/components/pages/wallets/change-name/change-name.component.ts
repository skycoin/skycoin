import { Component, OnInit, Inject, ViewChild, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { SubscriptionLike } from 'rxjs';

import { ButtonComponent } from '../../../layout/button/button.component';
import { MessageIcons } from '../../../layout/hardware-wallet/hw-message/hw-message.component';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { processServiceError } from '../../../../utils/errors';
import { MsgBarService } from '../../../../services/msg-bar.service';
import { AppConfig } from '../../../../app.config';
import { SoftwareWalletService } from '../../../../services/wallet-operations/software-wallet.service';
import { WalletBase } from '../../../../services/wallet-operations/wallet-objects';
import { HardwareWalletService } from '../../../../services/wallet-operations/hardware-wallet.service';

/**
 * States in which ChangeNameComponent can be.
 */
enum States {
  /**
   * Showing the form.
   */
  Initial = 'Initial',
  /**
   * Asking the user to confirm the operation on the hw wallet.
   */
  WaitingForConfirmation = 'WaitingForConfirmation',
}

/**
 * Settings for ChangeNameComponent.
 */
export class ChangeNameData {
  /**
   * Wallet whose label will be changed.
   */
  wallet: WalletBase;
  /**
   * New label. If provided, the form for entering the new label will not be shown and the
   * procedure will start immediately. NOTE: only for hw wallets. If the operation fails, the
   * error will be returned after closing the modal window.
   */
  newName: string;
}

/**
 * Response returned if ChangeNameComponent was opened with the new label included in the
 * configuration and the operation failed .
 */
export class ChangeNameErrorResponse {
  errorMsg: string;
}

/**
 * Modal window for changing the label of a software or hardware wallet. If the label is changed,
 * the modal window is closed and the new label is returned in the "afterClosed" event. If the new
 * label is provided in the configuration when opening the window, the form is not shown. In that
 * case, if there was a problem changing the label the modal window is closed and a
 * ChangeNameErrorResponse instance is returned in the "afterClosed" event
 */
@Component({
  selector: 'app-change-name',
  templateUrl: './change-name.component.html',
  styleUrls: ['./change-name.component.scss'],
})
export class ChangeNameComponent implements OnInit, OnDestroy {
  // Confirmation button.
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;
  maxHwWalletLabelLength = HwWalletService.maxLabelLength;
  // Deactivates the form while the system is busy.
  working = false;

  // Vars with the validation error messages.
  inputErrorMsg = '';

  private hwConnectionSubscription: SubscriptionLike;
  private operationSubscription: SubscriptionLike;

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog, data: ChangeNameData, smallSize: boolean): MatDialogRef<ChangeNameComponent, any> {
    const config = new MatDialogConfig();
    config.data = data;
    config.autoFocus = true;
    config.width = smallSize ? '400px' : AppConfig.mediumModalWidth;

    return dialog.open(ChangeNameComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<ChangeNameComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ChangeNameData,
    private formBuilder: FormBuilder,
    private hwWalletService: HwWalletService,
    private msgBarService: MsgBarService,
    private softwareWalletService: SoftwareWalletService,
    private hardwareWalletService: HardwareWalletService,
    private changeDetector: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    if (!this.data.newName) {
      // If the configuration did not include the new label, show the form so the user
      // can enter it.
      this.form = this.formBuilder.group({
        label: [this.data.wallet.label],
      });

      this.form.setValidators(this.validateForm.bind(this));
    } else {
      // If the configuration included the new label, start the change immediately.
      this.finishRenaming(this.data.newName);
    }

    if (this.data.wallet.isHardware) {
      this.hwConnectionSubscription = this.hwWalletService.walletConnectedAsyncEvent.subscribe(connected => {
        if (!connected) {
          this.closePopup();
        }
      });
    }
  }

  ngOnDestroy() {
    this.msgBarService.hide();
    if (this.hwConnectionSubscription) {
      this.hwConnectionSubscription.unsubscribe();
    }
    if (this.operationSubscription) {
      this.operationSubscription.unsubscribe();
    }
  }

  closePopup() {
    this.dialogRef.close();
  }

  // Changes the label with the data entered on the form.
  rename() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.button.setLoading();

    this.finishRenaming(this.form.value.label);

    this.changeDetector.detectChanges();
  }

  /**
   * Performs the procedure for changing the label.
   * @param newLabel New label.
   */
  private finishRenaming(newLabel: string) {
    this.working = true;

    if (!this.data.wallet.isHardware) {
      this.operationSubscription = this.softwareWalletService.renameWallet(this.data.wallet, newLabel)
        .subscribe(() => {
          this.working = false;
          this.dialogRef.close(newLabel);
          setTimeout(() => this.msgBarService.showDone('common.changes-made'));
        }, e => {
          this.working = false;
          this.msgBarService.showError(e);
          if (this.button) {
            this.button.resetState();
          }
        });
    } else {
      this.currentState = States.WaitingForConfirmation;

      this.operationSubscription = this.hardwareWalletService.changeLabel(this.data.wallet, newLabel).subscribe(
        () => {
          this.working = false;
          this.dialogRef.close(newLabel);

          // Don't show a confirmation msg if the new label was included in the configuration.
          if (!this.data.newName) {
            setTimeout(() => this.msgBarService.showDone('common.changes-made'));
          }
        },
        err => {
          this.working = false;
          if (this.data.newName) {
            // Return the error to the caller.
            const response = new ChangeNameErrorResponse();
            response.errorMsg = processServiceError(err).translatableErrorMsg;
            this.dialogRef.close(response);
          } else {
            // Show the form again.
            this.msgBarService.showError(err);
            this.currentState = States.Initial;
            if (this.button) {
              this.button.resetState();
            }
          }
        },
      );
    }
  }

  /**
   * Validates the form and updates the vars with the validation errors.
   */
  validateForm() {
    this.inputErrorMsg = '';

    let valid = true;

    if (!this.form.get('label').value || (this.form.get('label').value as string).trim().length === 0) {
      valid = false;
      if (this.form.get('label').touched) {
        this.inputErrorMsg = 'wallet.rename.label-error-info';
      }
    }

    return valid ? null : { Invalid: true };
  }
}
