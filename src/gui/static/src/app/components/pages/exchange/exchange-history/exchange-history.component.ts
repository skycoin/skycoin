import { Component, OnInit } from '@angular/core';
import { StoredExchangeOrder, Wallet } from '../../../../app.datatypes';
import { MatDialogRef } from '@angular/material';
import { WalletService } from '../../../../services/wallet.service';
import { ExchangeService } from '../../../../services/exchange.service';

@Component({
  selector: 'app-exchange-history',
  templateUrl: './exchange-history.component.html',
  styleUrls: ['./exchange-history.component.scss'],
})
export class ExchangeHistoryComponent implements OnInit {
  orders: StoredExchangeOrder[];

  constructor(
    public dialogRef: MatDialogRef<ExchangeHistoryComponent>,
    private exchangeService: ExchangeService,
  ) { }

  ngOnInit() {
    this.exchangeService.history().subscribe(
      (orders: StoredExchangeOrder[]) => this.orders = orders,
      () => this.orders = [],
    );
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value) {
    this.dialogRef.close(value);
  }
}
