import {Component, OnInit} from "@angular/core";
import {WalletService} from "../../services/wallet.service";
import {AddressBalance} from "../models/models";
import * as _ from "underscore";

@Component({
  selector: 'home-page-app-coinbalance',
  templateUrl: './coinbalance.component.html',
  styleUrls: ['./coinbalance.component.css']
})
export class CoinbalanceComponent implements OnInit {

  constructor(private _walletService:WalletService) { }

  walletFinal:any;

  ngOnInit() {
    this.walletFinal = {};
    this._walletService.getWallets().subscribe(wallets=>{
      _.each(wallets, (wallet)=>{
        let walletTemp:any={};
        this.walletFinal.coin = wallet.meta.coin;
        _.each(wallet.entries,(entry)=>{
          this.walletFinal.balance = 0;
          this._walletService.getCurrentBalanceOfAddress(entry.address).
            subscribe((addressBalance: AddressBalance) =>{
            this.walletFinal.balance = this.walletFinal.balance  + addressBalance.confirmed.coins/1000000;
          });
        });
      });
    });
  }



}
