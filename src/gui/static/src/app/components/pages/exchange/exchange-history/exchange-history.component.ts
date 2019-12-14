import { Component, OnInit } from '@angular/core';
import { StoredExchangeOrder, Wallet } from '../../../../app.datatypes';
import { MatDialogRef } from '@angular/material/dialog';
import { ExchangeService } from '../../../../services/exchange.service';
import { BlockchainService } from '../../../../services/blockchain.service';

@Component({
  selector: 'app-exchange-history',
  templateUrl: './exchange-history.component.html',
  styleUrls: ['./exchange-history.component.scss'],
})
export class ExchangeHistoryComponent implements OnInit {
  orders: StoredExchangeOrder[];

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
