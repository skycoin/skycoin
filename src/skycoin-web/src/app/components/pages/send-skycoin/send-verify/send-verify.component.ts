import { Component, EventEmitter, Input, OnDestroy, Output, ViewChild } from '@angular/core';

import { PreviewTransaction } from '../../../../app.datatypes';
import { BalanceService } from '../../../../services/wallet/balance.service';
import { SpendingService } from '../../../../services/wallet/spending.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { parseResponseMessage } from './../../../../utils/errors';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
  selector: 'app-send-verify',
  templateUrl: './send-verify.component.html',
  styleUrls: ['./send-verify.component.scss'],
})
export class SendVerifyComponent implements OnDestroy {
  @ViewChild('sendButton') sendButton: ButtonComponent;
  @ViewChild('backButton') backButton: ButtonComponent;
  @Input() transaction: PreviewTransaction;
  @Output() onBack = new EventEmitter<boolean>();

  constructor(
    private balanceService: BalanceService,
    private spendingService: SpendingService,
    private msgBarService: MsgBarService,
  ) {}

  ngOnDestroy() {
    this.msgBarService.hide();
  }

  send() {
    this.msgBarService.hide();
    this.sendButton.resetState();
    this.sendButton.setLoading();
    this.backButton.setDisabled();

    this.spendingService.injectTransaction(this.transaction.encoded)
      .subscribe(
        () => this.onSuccess(),
        (error) => this.onError(error)
      );
  }

  back() {
    this.onBack.emit(false);
  }

  private onSuccess() {
    setTimeout(() => this.msgBarService.showDone('send.sent'));
    this.balanceService.startGettingBalances();
    this.onBack.emit(true);
  }

  private onError(error) {
    const errorMessage = error.message ? error.message : parseResponseMessage(error['_body']);
    this.msgBarService.showError(errorMessage);
    this.sendButton.resetState();
  }
}
