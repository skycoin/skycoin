import { Component, OnInit } from '@angular/core';
import { BlockchainService } from '../../../services/blockchain.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.css']
})
export class AddressComponent implements OnInit {

  transactions: any[];
  balance: any;
  id: string;

  constructor(
    private blockchainService: BlockchainService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.params.switchMap(params => {
      this.id = params.address;
      return this.blockchainService.addressTransactions(params.address);
    }).subscribe(response => {
      this.transactions = response;
      console.log(response);
    });
    this.route.params.switchMap(params => this.blockchainService.addressBalance(params.address)).subscribe(response => {
      this.balance = response;
      console.log(response);

    });
  }

}
