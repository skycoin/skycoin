import { Component, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MatSnackBar } from '@angular/material';
import { showSnackbarError } from '../../../../utils/errors';
import { PreviewTransaction } from '../../../../app.datatypes';

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

  constructor(
    private walletService: WalletService,
    private snackbar: MatSnackBar,
  ) {}

  ngOnDestroy() {
    this.snackbar.dismiss();
  }

  send() {
    if (this.sendButton.isLoading()) {
      return;
    }

    this.snackbar.dismiss();
    this.sendButton.resetState();
    this.sendButton.setLoading();
    this.backButton.setDisabled();

    this.walletService.injectTransaction(this.transaction.encoded).subscribe(() => {
      this.sendButton.setSuccess();
      this.sendButton.setDisabled();

      this.walletService.startDataRefreshSubscription();

      setTimeout(() => {
        this.onBack.emit(true);
      }, 3000);
    }, error => {
      showSnackbarError(this.snackbar, error);

      this.sendButton.setError(error);
      this.backButton.setEnabled();
    });
  }

  back() {
    this.onBack.emit(false);
  }
}
