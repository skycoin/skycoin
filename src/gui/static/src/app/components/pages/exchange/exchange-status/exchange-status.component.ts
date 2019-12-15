import { Component, Input, OnDestroy, Output, EventEmitter } from '@angular/core';
import { ExchangeOrder, StoredExchangeOrder, ConfirmationData } from '../../../../app.datatypes';
import { ExchangeService } from '../../../../services/exchange.service';
import { QrCodeComponent, QrDialogConfig } from '../../../layout/qr-code/qr-code.component';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { SubscriptionLike, Observable, of } from 'rxjs';
import { showConfirmationModal } from '../../../../utils';
import { BlockchainService } from '../../../../services/blockchain.service';
import { environment } from '../../../../../environments/environment';
import { AppService } from '../../../../services/app.service';
import { delay, mergeMap } from 'rxjs/operators';

@Component({
  selector: 'app-exchange-status',
  templateUrl: './exchange-status.component.html',
  styleUrls: ['./exchange-status.component.scss'],
})
export class ExchangeStatusComponent implements OnDestroy {
  private readonly TEST_MODE = environment.swaplab.activateTestMode;
  private readonly TEST_ERROR = environment.swaplab.endStatusInError;
  readonly statuses = [
    'user_waiting',
    'market_waiting_confirmations',
    'market_confirmed',
    'market_exchanged',
    'market_withdraw_waiting',
    'complete',
    'error',
    'user_deposit_timeout',
  ];

  loading = true;
  showError = false;
  expanded = false;

  private subscription: SubscriptionLike;
  private testStatusIndex = 0;
  private order: ExchangeOrder;

  _orderDetails: StoredExchangeOrder;
  @Input() set orderDetails(val: StoredExchangeOrder) {
    const oldOrderDetails = this._orderDetails;
    this._orderDetails = val;

    if (val !== null && (!oldOrderDetails || oldOrderDetails.id !== val.id)) {
      this.exchangeService.lastViewedOrder = this._orderDetails;
      this.testStatusIndex = 0;
      this.loading = true;
      this.getStatus();
    }
  }

  @Output() goBack = new EventEmitter<void>();

  get fromCoin() {
    return this.order.pair.split('/')[0].toUpperCase();
  }

  get toCoin() {
    return this.order.pair.split('/')[1].toUpperCase();
  }

  get translatedStatus() {
    const status = this.order.status.replace(/_/g, '-');
    const params = {
      from: this.fromCoin,
      amount: this.order.fromAmount,
      to: this.toCoin,
    };

    return {
      text: `exchange.statuses.${status}`,
      info: `exchange.statuses.${status}-info`,
      params,
    };
  }

  get statusIcon() {
    if (this.order.status === this.statuses[5]) {
      return 'done';
    }

    if (this.order.status === this.statuses[6] || this.order.status === this.statuses[7]) {
      return 'close';
    }

    return 'refresh';
  }

  get progress() {
    let index = this.statuses.indexOf(this.order.status);

    index = Math.min(index, 5) + 1;

    return Math.ceil((100 / 6) * index);
  }

  constructor(
    private exchangeService: ExchangeService,
    private dialog: MatDialog,
    public blockchainService: BlockchainService,
    public appService: AppService,
  ) { }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  showQrCode(address) {
    const config: QrDialogConfig = {
      address: address,
      hideCoinRequestForm: true,
      ignoreCoinPrefix: true,
    };
    QrCodeComponent.openDialog(this.dialog, config);
  }

  toggleDetails() {
    this.expanded = !this.expanded;
  }

  close() {
    if (this.loading || this.exchangeService.isOrderFinished(this.order)) {
      this.goBack.emit();
    } else {
      const confirmationData: ConfirmationData = {
        text: 'exchange.details.back-alert',
        headerText: 'confirmation.header-text',
        confirmButtonText: 'confirmation.confirm-button',
        cancelButtonText: 'confirmation.cancel-button',
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.goBack.emit();
        }
      });
    }
  }

  private getStatus(delayTime = 0) {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }

    const fromAmount = this._orderDetails.fromAmount;

    if (this.TEST_MODE && this.TEST_ERROR && this.testStatusIndex === this.statuses.length - 2) {
      this.testStatusIndex = this.statuses.length - 1;
    }

    this.subscription = of(0).pipe(delay(delayTime), mergeMap(() => {
      return this.exchangeService.status(
        !this.TEST_MODE ? this._orderDetails.id : '4729821d-390d-4ef8-a31e-2465d82a142f',
        !this.TEST_MODE ? null : this.statuses[this.testStatusIndex++],
      );
    })).subscribe(order => {
      this.order = { ...order, fromAmount };
      this._orderDetails.id = order.id;
      this.exchangeService.lastViewedOrder = this._orderDetails;

      if (!this.exchangeService.isOrderFinished(order)) {
        this.getStatus(this.TEST_MODE ? 3000 : 30000);
      } else {
        this.exchangeService.lastViewedOrder = null;
      }

      this.loading = false;
    }, () => this.showError = true);
  }
}
