import { Component, OnDestroy, OnInit } from '@angular/core';
import { ExchangeOrder, StoredExchangeOrder } from '../../../app.datatypes';
import { ExchangeService } from '../../../services/exchange.service';
import { MatDialog, MatDialogConfig, MatSnackBar } from '@angular/material';
import { ExchangeHistoryComponent } from './exchange-history/exchange-history.component';
import { showSnackbarError } from '../../../utils/errors';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-exchange',
  templateUrl: './exchange.component.html',
  styleUrls: ['./exchange.component.scss'],
})
export class ExchangeComponent implements OnInit, OnDestroy {
  order: ExchangeOrder;
  hasHistory = false;

  constructor(
    private exchangeService: ExchangeService,
    private dialog: MatDialog,
    private snackbar: MatSnackBar,
    private translate: TranslateService,
  ) { }

  ngOnInit() {
    const lastOrder = this.exchangeService.lastOrder;

    if (lastOrder) {
      if (!this.exchangeService.isOrderFinished(lastOrder)) {
        this.showLast();
      }
    }

    this.exchangeService.history().subscribe(() => this.hasHistory = true);
  }

  ngOnDestroy() {
    this.snackbar.dismiss();
  }

  showLast() {
    this.order = this.exchangeService.lastOrder;
  }

  showStatus(lastOrder) {
    this.order = lastOrder;
  }

  showHistory(event) {
    event.preventDefault();

    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;

    this.dialog.open(ExchangeHistoryComponent, config).afterClosed().subscribe((oldOrder: StoredExchangeOrder) => {
      if (oldOrder) {
        this.exchangeService.status(oldOrder.id).first()
          .subscribe(
            order => this.order = { ...order, fromAmount: oldOrder.fromAmount },
            () => showSnackbarError(this.snackbar, this.translate.instant('exchange.order-not-found'), 3000),
          );
      }
    });
  }
}
