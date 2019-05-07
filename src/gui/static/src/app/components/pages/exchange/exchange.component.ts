import * as moment from 'moment';
import { Component, OnInit } from '@angular/core';
import { StoredExchangeOrder } from '../../../app.datatypes';
import { ExchangeService } from '../../../services/exchange.service';
import { MatDialog, MatDialogConfig } from '@angular/material';
import { ExchangeHistoryComponent } from './exchange-history/exchange-history.component';
import { ISubscription } from 'rxjs/Subscription';

@Component({
  selector: 'app-exchange',
  templateUrl: './exchange.component.html',
  styleUrls: ['./exchange.component.scss'],
})
export class ExchangeComponent implements OnInit {
  previewOrderDetails: StoredExchangeOrder;
  hasHistory = false;
  loading = true;

  constructor(
    private exchangeService: ExchangeService,
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
    const sub = this.exchangeService.lastViewedOrderLoaded.subscribe(response => {
      if (response) {
        const lastViewedOrder = this.exchangeService.lastViewedOrder;
        if (lastViewedOrder) {
          this.previewOrderDetails = lastViewedOrder;
        }

        setTimeout(() => sub.unsubscribe());
        this.loading = false;
      }
    });

    this.exchangeService.history().subscribe(() => this.hasHistory = true);
  }

  showStatus(order) {
    this.previewOrderDetails = order;
  }

  showHistory(event) {
    event.preventDefault();

    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;

    this.dialog.open(ExchangeHistoryComponent, config).afterClosed().subscribe((oldOrder: StoredExchangeOrder) => {
      if (oldOrder) {
        this.previewOrderDetails = oldOrder;
      }
    });
  }

  goBack() {
    this.previewOrderDetails = null;
    this.exchangeService.lastViewedOrder = null;
  }
}
