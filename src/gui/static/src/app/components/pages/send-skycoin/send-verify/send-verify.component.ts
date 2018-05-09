import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MatSnackBar, MatSnackBarConfig } from '@angular/material';
import { parseResponseMessage } from '../../../../utils/index';
import { PriceService } from '../../../../services/price.service';

@Component({
  selector: 'app-send-verify',
  templateUrl: './send-verify.component.html',
  styleUrls: ['./send-verify.component.scss']
})
export class SendVerifyComponent implements OnInit {
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @ViewChild('backButton') backButton: ButtonComponent;
  @Input() transaction: any;
  @Input() encodedTransaction: string;
  @Input() formData: any;
  @Output() onBack = new EventEmitter<boolean>();

  price: number;
  showInputsOutputs = false;
  sentHours = 0;

  constructor(
    private walletService: WalletService,
    private snackbar: MatSnackBar,
    private priceService: PriceService,
  ) {
    this.priceService.price.subscribe(price => this.price = price);
  }

  ngOnInit() {
    this.sentHours = this.transaction.outputs
      .filter(o => this.transaction.inputs.find(i => i.address !== o.address))
      .map(o => parseInt(o.hours, 10))
      .reduce((a, b) => a + b, 0);

    if (this.sentHours === 0 && this.transaction.outputs.length === 1) {
      this.sentHours = this.transaction.outputs[0].hours;
    }
  }

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
      const errorMessage = parseResponseMessage(error['_body']);
      const config = new MatSnackBarConfig();
      config.duration = 300000;
      this.snackbar.open(errorMessage, null, config);
      this.sendButton.setError(errorMessage);
    });
  }

  back() {
    this.onBack.emit(false);
  }

  toggleInputsOutputs(event) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }
}
