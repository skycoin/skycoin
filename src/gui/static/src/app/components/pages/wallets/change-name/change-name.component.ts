import { Component, OnInit, Inject, ViewChild, OnDestroy } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { FormBuilder, Validators, FormGroup } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Wallet } from '../../../../app.datatypes';
import { ButtonComponent } from '../../../layout/button/button.component';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { showSnackbarError, getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { TranslateService } from '@ngx-translate/core';
import { MatSnackBar } from '@angular/material';
import { MessageIcons } from '../../../layout/hardware-wallet/hw-message/hw-message.component';

enum States {
  Initial,
  WaitingForConfirmation,
}

export class ChangeNameData {
  wallet: Wallet;
  newName: string;
}

export class ChangeNameErrorResponse {
  errorMsg: string;
}

@Component({
  selector: 'app-change-name',
  templateUrl: './change-name.component.html',
  styleUrls: ['./change-name.component.css'],
})
export class ChangeNameComponent implements OnInit, OnDestroy {
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;
  currentState: States = States.Initial;
  states = States;
  msgIcons = MessageIcons;

  private newLabel: string;

  constructor(
    public dialogRef: MatDialogRef<ChangeNameComponent>,
    @Inject(MAT_DIALOG_DATA) private data: ChangeNameData,
    private formBuilder: FormBuilder,
    private walletService: WalletService,
    private hwWalletService: HwWalletService,
    private translateService: TranslateService,
    private snackbar: MatSnackBar,
  ) {}

  ngOnInit() {
    if (!this.data.newName) {
      this.form = this.formBuilder.group({
        label: [this.data.newName ? this.data.newName : this.data.wallet.label, Validators.required],
      });
    } else {
      this.finishRenaming(this.data.newName);
    }
  }

  ngOnDestroy() {
    this.snackbar.dismiss();
  }

  closePopup() {
    this.dialogRef.close();
  }

  rename() {
    if (!this.form.valid || this.button.isLoading()) {
      return;
    }

    this.snackbar.dismiss();
    this.button.setLoading();

    this.finishRenaming(this.form.value.label);
  }

  private finishRenaming(newLabel) {
    this.newLabel = newLabel;

    if (!this.data.wallet.isHardware) {
      this.walletService.renameWallet(this.data.wallet, this.newLabel)
        .subscribe(() => this.dialogRef.close(this.newLabel));
    } else {
      if (this.data.newName) {
        this.currentState = States.WaitingForConfirmation;
      }

      this.hwWalletService.checkIfCorrectHwConnected(this.data.wallet.addresses[0].address)
        .flatMap(() => {
          this.currentState = States.WaitingForConfirmation;

          return this.hwWalletService.changeLabel(this.newLabel);
        })
        .subscribe(
          () => {
            this.data.wallet.label = this.newLabel;
            this.walletService.saveHardwareWallets();
            this.dialogRef.close(this.newLabel);
          },
          err => {
            if (this.data.newName) {
              const response = new ChangeNameErrorResponse();
              response.errorMsg = getHardwareWalletErrorMsg(this.hwWalletService, this.translateService, err);
              this.dialogRef.close(response);
            } else {
              showSnackbarError(this.snackbar, getHardwareWalletErrorMsg(this.hwWalletService, this.translateService, err));
              this.currentState = States.Initial;
              if (this.button) {
                this.button.resetState();
              }
            }
          },
        );
    }
  }
}
