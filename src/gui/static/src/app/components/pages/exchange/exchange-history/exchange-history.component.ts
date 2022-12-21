import { Component, OnInit } from '@angular/core';
import { MatLegacyDialogRef as MatDialogRef, MatLegacyDialog as MatDialog, MatLegacyDialogConfig as MatDialogConfig } from '@angular/material/legacy-dialog';

import { ExchangeService, StoredExchangeOrder } from '../../../../services/exchange.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { AppConfig } from '../../../../app.config';

/**
 * Modal window for checking the lists of previously saved orders. If the user selects an order,
 * the modal window is closed and the selected order is returned in the "afterClosed" event.
 */
@Component({
  selector: 'app-exchange-history',
  templateUrl: './exchange-history.component.html',
  styleUrls: ['./exchange-history.component.scss'],
})
export class ExchangeHistoryComponent implements OnInit {
  // List of saved orders.
  orders: StoredExchangeOrder[];

  /**
   * Opens the modal window. Please use this function instead of opening the window "by hand".
   */
  public static openDialog(dialog: MatDialog): MatDialogRef<ExchangeHistoryComponent, StoredExchangeOrder> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(ExchangeHistoryComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<ExchangeHistoryComponent>,
    public blockchainService: BlockchainService,
    private exchangeService: ExchangeService,
  ) { }

  ngOnInit() {
    // Get the saved transactions.
    this.exchangeService.history().subscribe(
      (orders: StoredExchangeOrder[]) => this.orders = orders.reverse(),
      () => this.orders = [],
    );
  }

  closePopup() {
    this.dialogRef.close();
  }

  // Closes the modal window and returns the selected order.
  select(value: StoredExchangeOrder) {
    this.dialogRef.close(value);
  }

  // Gets the "from" part of the name of a trading pair.
  getFromCoin(pair: string) {
    return pair.split('/')[0].toUpperCase();
  }

  // Gets the "to" part of the name of a trading pair.
  getToCoin(pair: string) {
    return pair.split('/')[1].toUpperCase();
  }
}
