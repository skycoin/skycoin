import { Component, OnInit } from '@angular/core';
import {CoinSupplyService} from "./coin-supply.service";
import {CoinSupply} from "../block-chain-table/block";

@Component({
  selector: 'app-block-chain-coin-supply',
  templateUrl: './block-chain-coin-supply.component.html',
  styleUrls: ['./block-chain-coin-supply.component.css']
})
export class BlockChainCoinSupplyComponent implements OnInit {

  private coinSupply:number;
  private coinCap:number;

  constructor(private service:CoinSupplyService) {
    this.coinSupply=this.coinCap=0;
  }

  ngOnInit() {
    this.service.getCoinSupply().subscribe((supply:CoinSupply)=>{
      this.coinCap =supply.coinCap;
      this.coinSupply = supply.coinSupply;
    })
  }

}
