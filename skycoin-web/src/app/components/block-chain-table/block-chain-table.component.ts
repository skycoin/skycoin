import { Component, OnInit } from '@angular/core';
import {BlockChainService} from "./block-chain.service";
import {Observable} from "rxjs";
import {Block, Transaction} from "./block";
import * as moment from 'moment';

import * as _ from 'underscore';
import {Router} from "@angular/router";

@Component({
  selector: 'app-block-chain-table',
  templateUrl: './block-chain-table.component.html',
  styleUrls: ['./block-chain-table.component.css']
})
export class BlockChainTableComponent implements OnInit {

  private blocks:Block[];

  constructor(private blockService:BlockChainService,private router: Router) {
    this.blocks=[];
  }

  GetBlockAmount(txns:Transaction[]) {
    var ret = [];
    _.each(txns, function(o){
      if(o.outputs){
        _.each(o.outputs, function(_o){
          ret.push(_o.coins);
        })
      }

    })
    let totalCoins=ret.reduce(function(memo, coin) {
      return memo + parseInt(coin);
    }, 0);
    return totalCoins;
  }

  getTime(time){
    return moment.unix(time).format("YYYY-MM-DD HH:mm");
  }

  ngOnInit() {
    this.blockService.getBlocks(1,10).subscribe(
      (data)=>{
        this.blocks= data;
      }
    );
  }

  showDetails(block: Block) {
    this.router.navigate(['/block', block.header.seq]);
  }

  handlePageChange(currentPage:number){
    this.blockService.getBlocks((currentPage-1)*10+1,(currentPage-1)*10+10).subscribe(
      (data)=>{
        console.log(data);
        this.blocks= data;
      }
    );
  }

}
