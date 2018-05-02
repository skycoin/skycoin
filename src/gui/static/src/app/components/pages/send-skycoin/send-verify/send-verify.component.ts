import { Component, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';

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
  ) {}

  send() {
    if (this.sendButton.isLoading()) {
      return;
    }

    this.sendButton.setLoading();
    this.backButton.setDisabled();

    this.walletService.injectTransaction(this.encodedTransaction).subscribe(() => {
      this.sendButton.setSuccess();
      this.sendButton.setDisabled();

      this.walletService.startDataRefreshSubscription();

      setTimeout(() => {
        this.onBack.emit(true);
      }, 3000);
    }, error => this.sendButton.setError(error));
  }

  back() {
    this.onBack.emit(false);
  }
}
