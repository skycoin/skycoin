import { Component, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { parseResponseMessage } from '../../../../utils/index';

@Component({
  selector: 'app-send-verify',
  templateUrl: './send-verify.component.html',
  styleUrls: ['./send-verify.component.scss']
})
export class SendVerifyComponent {
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @ViewChild('backButton') backButton: ButtonComponent;
  @Input() transaction: any;
  @Input() encodedTransaction: string;
  @Output() onBack = new EventEmitter<boolean>();

  constructor(
    private walletService: WalletService,
    private snackbar: MatSnackBar,
  ) {}

  send() {
    if (this.sendButton.isLoading()) {
      return;
    }

    this.snackbar.dismiss();
    this.sendButton.resetState();
    this.sendButton.setLoading();
    this.backButton.setDisabled();

    this.walletService.injectTransaction(this.encodedTransaction).subscribe(() => {
      this.sendButton.setSuccess();
      this.sendButton.setDisabled();

      this.walletService.startDataRefreshSubscription();

      setTimeout(() => {
        this.onBack.emit(true);
      }, 3000);
    }, error => {
      const errorMessage = parseResponseMessage(error);
      const config = new MatSnackBarConfig();
      config.duration = 300000;
      this.snackbar.open(errorMessage, null, config);
      this.sendButton.setError(errorMessage);
    });
  }

  back() {
    this.onBack.emit(false);
  }
}
