import { Component, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { ActivatedRoute } from '@angular/router';
import 'rxjs/add/operator/switchMap';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {

  total: number;
  transaction: any;

  constructor(
    private route: ActivatedRoute,
    private walletService: WalletService
  ) { }

  ngOnInit() {
    this.route.params.switchMap(params => this.walletService.transaction(params.transaction)).subscribe(transaction => {
      this.transaction = transaction;
      this.total = transaction.txn.outputs.reduce((a , b) => a + parseInt(b.coins), 0);
    });
  }

}
