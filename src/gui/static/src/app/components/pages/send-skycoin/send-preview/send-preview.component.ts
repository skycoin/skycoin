import { Component, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MatDialogConfig, MatDialog } from '@angular/material';
import { getHardwareWalletErrorMsg } from '../../../../utils/errors';
import { PreviewTransaction, Wallet } from '../../../../app.datatypes';
import { ISubscription } from 'rxjs/Subscription';
import { PasswordDialogComponent } from '../../../layout/password-dialog/password-dialog.component';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
  selector: 'app-send-preview',
  templateUrl: './send-preview.component.html',
  styleUrls: ['./send-preview.component.scss'],
})
export class SendVerifyComponent implements OnDestroy {
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @ViewChild('backButton') backButton: ButtonComponent;
  @Input() transaction: PreviewTransaction;
  @Output() onBack = new EventEmitter<boolean>();

  private sendSubscription: ISubscription;

  constructor(
    private walletService: WalletService,
    private msgBarService: MsgBarService,
    private dialog: MatDialog,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
  ) {}

  ngOnDestroy() {
    this.msgBarService.hide();

    if (this.sendSubscription) {
      this.sendSubscription.unsubscribe();
    }
  }

  back() {
    this.onBack.emit(false);
  }

  send() {
    if (this.sendButton.isLoading()) {
      return;
    }

    this.msgBarService.hide();
    this.sendButton.resetState();

    if (this.transaction.wallet.encrypted && !this.transaction.wallet.isHardware) {
      const config = new MatDialogConfig();
      config.data = {
        wallet: this.transaction.wallet,
      };

      this.dialog.open(PasswordDialogComponent, config).componentInstance.passwordSubmit
        .subscribe(passwordDialog => {
          this.finishSending(passwordDialog);
        });
    } else {
      if (!this.transaction.wallet.isHardware) {
        this.finishSending();
      } else {
        this.showBusy();
        this.sendSubscription = this.hwWalletService.checkIfCorrectHwConnected(this.transaction.wallet.addresses[0].address).subscribe(
          () => this.finishSending(),
          err => this.showError(getHardwareWalletErrorMsg(this.hwWalletService, this.translate, err)),
        );
      }
    }
  }

  private showBusy() {
    this.sendButton.setLoading();
    this.backButton.setDisabled();
  }

  private finishSending(passwordDialog?: any) {
    this.showBusy();

    const note = this.transaction.note.trim();

    this.sendSubscription = this.walletService.signTransaction(
      this.transaction.wallet,
      passwordDialog ? passwordDialog.password : null,
      this.transaction,
    ).flatMap(result => {
      if (passwordDialog) {
        passwordDialog.close();
      }

      return this.walletService.injectTransaction(result.encoded, note);
    }).subscribe(noteSaved => {
      if (note && !noteSaved) {
        this.msgBarService.showError(this.translate.instant('send.error-saving-note'));
      }

      this.sendButton.setSuccess();
      this.sendButton.setDisabled();

      this.walletService.startDataRefreshSubscription();

      setTimeout(() => {
        this.onBack.emit(true);
      }, 3000);
    }, error => {
      if (passwordDialog) {
        passwordDialog.error(error);
      }

      if (error && error.result) {
        this.showError(getHardwareWalletErrorMsg(this.hwWalletService, this.translate, error));
      } else {
        this.showError(error);
      }
    });
  }

  private showError(error) {
    this.msgBarService.showError(error);
    this.sendButton.resetState().setError(error);
    this.backButton.resetState().setEnabled();
  }
}
