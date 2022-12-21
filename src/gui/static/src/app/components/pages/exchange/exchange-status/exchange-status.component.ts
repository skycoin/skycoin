import { Component, Input, OnDestroy, Output, EventEmitter } from '@angular/core';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';
import { SubscriptionLike, of } from 'rxjs';
import { delay, mergeMap } from 'rxjs/operators';

import { ExchangeService, ExchangeOrder, StoredExchangeOrder } from '../../../../services/exchange.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { environment } from '../../../../../environments/environment';
import { AppService } from '../../../../services/app.service';
import { ConfirmationParams, DefaultConfirmationButtons, ConfirmationComponent } from '../../../layout/confirmation/confirmation.component';

/**
 * Shows the current status of a previously saved order. The data is updated periodically.
 * It needs a previously saved order because, depending on the current status of the order,
 * the backend may not return all the data, data which is normally added to a locally saved
 * order at creation time.
 */
@Component({
  selector: 'app-exchange-status',
  templateUrl: './exchange-status.component.html',
  styleUrls: ['./exchange-status.component.scss'],
})
export class ExchangeStatusComponent implements OnDestroy {
  // If the page must work in test mode. If true, the page will use the sandbox API endpoints
  // and show simulated state updates.
  private readonly TEST_MODE = environment.swaplab.activateTestMode;
  // If true and in test mode, the page will show an error as the last state of the order.
  private readonly TEST_ERROR = environment.swaplab.endStatusInError;
  // Known states the backend return.
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

  // True while the page is getting the data of the order for the first time.
  loading = true;
  // If there was an error when trying to get the data from the backend for the first time.
  showError = false;
  // If the UI must show the detailed information panel openned.
  expanded = false;

  private subscription: SubscriptionLike;
  // If in test mode, the index inside the statuses var of the status currently being shown.
  private testStatusIndex = 0;
  // Data obtained from the backend.
  private order: ExchangeOrder;

  // Previously saved order for which the data will be shown.
  _orderDetails: StoredExchangeOrder;
  @Input() set orderDetails(val: StoredExchangeOrder) {
    const oldOrderDetails = this._orderDetails;
    this._orderDetails = val;

    // As the value can be updated without recreaating the page, the state of the page is
    // resetted, but only if the order was changed.
    if (val !== null && (!oldOrderDetails || oldOrderDetails.id !== val.id)) {
      this.exchangeService.lastViewedOrder = this._orderDetails;
      this.testStatusIndex = 0;
      this.loading = true;
      this.getStatus();
    }
  }

  // Emits when the user presses the button for returning to the form.
  @Output() goBack = new EventEmitter<void>();

  // Gets the "from" part of the name of a trading pair.
  get fromCoin(): string {
    return this.order.pair.split('/')[0].toUpperCase();
  }

  // Gets the "to" part of the name of a trading pair.
  get toCoin(): string {
    return this.order.pair.split('/')[1].toUpperCase();
  }

  // Gets the params needed for using the translate pipe for showing the name and info
  // of the current state.
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
      params: params,
    };
  }

  // Gets the icon for the current status.
  get statusIcon(): string {
    if (this.order.status === this.statuses[5]) {
      return 'done';
    }

    if (this.order.status === this.statuses[6] || this.order.status === this.statuses[7]) {
      return 'close';
    }

    return 'refresh';
  }

  // Gets the completion percentage of the order, according to its current status.
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
    this.goBack.complete();
  }

  // Opens or closes the detailed information panel.
  toggleDetails() {
    this.expanded = !this.expanded;
  }

  // Sends the event for returning to the form. If the order has not been finished, the
  // user is asked for confirmation.
  close() {
    if (this.loading || this.exchangeService.isOrderFinished(this.order)) {
      this.goBack.emit();
    } else {
      const confirmationParams: ConfirmationParams = {
        text: 'exchange.details.back-alert',
        defaultButtons: DefaultConfirmationButtons.YesNo,
      };

      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.goBack.emit();
        }
      });
    }
  }

  // Periodically updates the status of the order.
  private getStatus(delayTime = 0) {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }

    const fromAmount = this._orderDetails.fromAmount;

    // Go to the correct status if using the test mode.
    if (this.TEST_MODE && this.TEST_ERROR && this.testStatusIndex === this.statuses.length - 2) {
      this.testStatusIndex = this.statuses.length - 1;
    }

    /* eslint-disable arrow-body-style */
    this.subscription = of(0).pipe(delay(delayTime), mergeMap(() => {
      // Orders created using the sandbox methods are not saved on the backend, so a
      // predefined ID is provided if the tests mode is being used, and a simulated status.
      return this.exchangeService.status(
        !this.TEST_MODE ? this._orderDetails.id : '4729821d-390d-4ef8-a31e-2465d82a142f',
        !this.TEST_MODE ? null : this.statuses[this.testStatusIndex++],
      );
    })).subscribe(order => {
      // Restore the amount of coins the user must send, as the backend may have not included it.
      this.order = { ...order, fromAmount: fromAmount };
      this._orderDetails.id = order.id;
      // Remember the last viewed order.
      this.exchangeService.lastViewedOrder = this._orderDetails;

      if (!this.exchangeService.isOrderFinished(order)) {
        this.getStatus(this.TEST_MODE ? 3000 : 30000);
      } else {
        // If the order is already finished, forget about it as the last viewed order, to avoid
        // restoring it when returning to this section of the app, instead of showing the form.
        this.exchangeService.lastViewedOrder = null;
      }

      this.loading = false;
    }, () => {
      if (this.loading) {
        this.showError = true;
      } else {
        this.getStatus(this.TEST_MODE ? 3000 : 30000);
      }
    });
  }
}
