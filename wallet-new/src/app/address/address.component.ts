import { Component, OnInit } from '@angular/core';
import {WalletService} from "../services/wallet.service";
import * as _ from "underscore";
@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.css']
})
export class AddressComponent implements OnInit {

  selectedAddress:string="";
  addresses:string[];

  options = [{id: 0, name:"First"}, {id: 1, name:"Second"}];

  selected1;
  selected2 = "";

  handleChange(event) {
    console.log(event.target.value);
    this.selected1 = this.options.filter((option) => {
      return option.id == event.target.value;
    })[0];

  }

  constructor(private _walletService:WalletService) { }

  ngOnInit() {
    this.addresses = [];
    this._walletService.getWallets().subscribe(wallets=>{
      _.each(wallets, (wallet)=>{
        let walletTemp:any={};
        _.each(wallet.entries,(entry)=>{
          this.addresses.push(entry.address);
        });
      });
    });
  }

  handleTap(address):void{
    this.selectedAddress = address;
  }

}
