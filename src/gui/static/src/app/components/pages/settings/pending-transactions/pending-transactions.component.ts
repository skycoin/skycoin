import { Component, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import * as moment from 'moment';
import { Router } from '@angular/router';

@Component({
  selector: 'app-pending-transactions',
  templateUrl: './pending-transactions.component.html',
  styleUrls: ['./pending-transactions.component.css']
})
export class PendingTransactionsComponent implements OnInit {

  transactions: any[];

  constructor(
    public walletService: WalletService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.walletService.pendingTransactions().subscribe(transactions => {
      this.transactions = this.mapTransactions(transactions);
    });
  }

  onActivate(response) {
    if (response.row && response.row.txid) {
      this.router.navigate(['/history', response.row.txid]);
    }
  }

  private mapTransactions(transactions) {
    return transactions.map(transaction => {
      transaction.transaction.timestamp = moment(transaction.received).unix();
      return transaction.transaction;
    })
    .map(transaction => {
      transaction.amount = transaction.outputs.map(output => output.coins >= 0 ? output.coins : 0)
        .reduce((a , b) => a + parseInt(b), 0);
      return transaction;
    });
  }
}
