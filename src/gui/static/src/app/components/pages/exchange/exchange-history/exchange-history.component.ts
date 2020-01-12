import { Component, OnInit } from '@angular/core';
import { StoredExchangeOrder } from '../../../../app.datatypes';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { ExchangeService } from '../../../../services/exchange.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { AppConfig } from '../../../../app.config';

@Component({
  selector: 'app-exchange-history',
  templateUrl: './exchange-history.component.html',
  styleUrls: ['./exchange-history.component.scss'],
})
export class ExchangeHistoryComponent implements OnInit {
  orders: StoredExchangeOrder[];

  public static openDialog(dialog: MatDialog): MatDialogRef<ExchangeHistoryComponent, any> {
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
    this.exchangeService.history().subscribe(
      (orders: StoredExchangeOrder[]) => this.orders = orders.reverse(),
      () => this.orders = [],
    );
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value) {
    this.dialogRef.close(value);
  }

  getFromCoin(pair: string) {
    return pair.split('/')[0].toUpperCase();
  }

  getToCoin(pair: string) {
    return pair.split('/')[1].toUpperCase();
  }
}
