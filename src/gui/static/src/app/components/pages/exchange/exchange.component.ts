import { Component, OnInit, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { SubscriptionLike } from 'rxjs';

import { ExchangeService, StoredExchangeOrder } from '../../../services/exchange.service';
import { ExchangeHistoryComponent } from './exchange-history/exchange-history.component';
import { AppService } from '../../../services/app.service';
import { environment } from '../../../../environments/environment';

/**
 * Allows to buy coins using the Swaplab service.
 */
@Component({
  selector: 'app-exchange',
  templateUrl: './exchange.component.html',
  styleUrls: ['./exchange.component.scss'],
})
export class ExchangeComponent implements OnInit, OnDestroy {
  // Order for which the status must be shown. If null, the form for creating a new order is shown.
  currentOrderDetails: StoredExchangeOrder;
  // If there are previously created orders saved on the persistent storage.
  hasHistory = false;
  // If the page is loading the initial data needed for starting to show anything.
  loading = true;
  // If the service is not available.
  unavailable = false;

  private lastViewedSubscription: SubscriptionLike;
  private historySubscription: SubscriptionLike;

  constructor(
    public appService: AppService,
    private exchangeService: ExchangeService,
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
    // The service is not available in the portable version, as the backend would return cors
    // errors and there is no dev proxy or Electron process for fixing the problem.
    if (environment.production && navigator.userAgent.toLowerCase().indexOf('electron') === -1) {
      this.unavailable = true;

      return;
    }

    // Check if there is a "last vieved" order saved and, if there is one, show its status.
    this.lastViewedSubscription = this.exchangeService.lastViewedOrderLoaded.subscribe(response => {
      if (response) {
        const lastViewedOrder = this.exchangeService.lastViewedOrder;
        if (lastViewedOrder) {
          this.currentOrderDetails = lastViewedOrder;
        }

        setTimeout(() => this.lastViewedSubscription.unsubscribe());
        this.loading = false;
      }
    });

    // Check if there is an previously created orders history.
    this.historySubscription = this.exchangeService.history().subscribe(() => this.hasHistory = true);
  }

  ngOnDestroy() {
    if (this.lastViewedSubscription) {
      this.lastViewedSubscription.unsubscribe();
    }
    if (this.historySubscription) {
      this.historySubscription.unsubscribe();
    }
  }

  // Shows the control for displaying the status of an order.
  showStatus(order: StoredExchangeOrder) {
    this.currentOrderDetails = order;
    this.hasHistory = true;
  }

  // Opens the modal window with the list of previously created orders.
  showHistory(event) {
    event.preventDefault();

    ExchangeHistoryComponent.openDialog(this.dialog).afterClosed().subscribe((oldOrder: StoredExchangeOrder) => {
      if (oldOrder) {
        this.currentOrderDetails = oldOrder;
      }
    });
  }

  // Returns to the form.
  goBack() {
    this.currentOrderDetails = null;
    // Prevent this page for showing the details of the order when returning to the
    // exchange section.
    this.exchangeService.lastViewedOrder = null;
  }
}
