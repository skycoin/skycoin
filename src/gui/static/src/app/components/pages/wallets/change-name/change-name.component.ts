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
    @Inject(MAT_DIALOG_DATA) private data: Wallet,
    private formBuilder: FormBuilder,
    private walletService: WalletService,
    private hwWalletService: HwWalletService,
    private translateService: TranslateService,
    private snackbar: MatSnackBar,
  ) {}

  ngOnInit() {
    this.form = this.formBuilder.group({
      label: [this.data.label, Validators.required],
    });
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

    this.newLabel = this.form.value.label;

    if (!this.data.isHardware) {
      this.walletService.renameWallet(this.data, this.newLabel)
        .subscribe(() => this.dialogRef.close(this.newLabel));
    } else {
      this.hwWalletService.checkIfCorrectHwConnected(this.data.addresses[0].address)
        .flatMap(() => {
          this.currentState = States.WaitingForConfirmation;

          return this.hwWalletService.changeLabel(this.newLabel);
        })
        .subscribe(
          () => {
            this.data.label = this.newLabel;
            this.walletService.saveHardwareWallets();
            this.dialogRef.close(this.newLabel);
          },
          err => {
            showSnackbarError(this.snackbar, getHardwareWalletErrorMsg(this.hwWalletService, this.translateService, err));
            this.currentState = States.Initial;
            if (this.button) {
              this.button.resetState();
            }
          },
        );
    }
  }
}
