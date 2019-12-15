import { Component, OnInit, OnDestroy } from '@angular/core';
import { StoredExchangeOrder } from '../../../app.datatypes';
import { ExchangeService } from '../../../services/exchange.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ExchangeHistoryComponent } from './exchange-history/exchange-history.component';
import { SubscriptionLike } from 'rxjs';

@Component({
  selector: 'app-exchange',
  templateUrl: './exchange.component.html',
  styleUrls: ['./exchange.component.scss'],
})
export class ExchangeComponent implements OnInit, OnDestroy {
  currentOrderDetails: StoredExchangeOrder;
  hasHistory = false;
  loading = true;

  private lastViewedSubscription: SubscriptionLike;
  private historySubscription: SubscriptionLike;

  constructor(
    private exchangeService: ExchangeService,
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
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

    this.historySubscription = this.exchangeService.history().subscribe(() => this.hasHistory = true);
  }

  ngOnDestroy() {
    this.lastViewedSubscription.unsubscribe();
    this.historySubscription.unsubscribe();
  }

  showStatus(order) {
    this.currentOrderDetails = order;
    this.hasHistory = true;
  }

  showHistory(event) {
    event.preventDefault();

    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;

    this.dialog.open(ExchangeHistoryComponent, config).afterClosed().subscribe((oldOrder: StoredExchangeOrder) => {
      if (oldOrder) {
        this.currentOrderDetails = oldOrder;
      }
    });
  }

  goBack() {
    this.currentOrderDetails = null;
    this.exchangeService.lastViewedOrder = null;
  }
}
