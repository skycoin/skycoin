import { Component, OnInit } from '@angular/core';
import {BlockChainService} from "./block-chain.service";
import {Observable} from "rxjs";
import {Block, Transaction} from "./block";
import * as moment from 'moment';
import {Router} from "@angular/router";
declare var _: any;

@Component({
  selector: 'app-block-chain-table',
  templateUrl: './block-chain-table.component.html',
  styleUrls: ['./block-chain-table.component.css']
})
export class BlockChainTableComponent implements OnInit {

  private blocks:Block[];
  private totalBlocks:number;
  private loading:boolean;
  constructor(private blockService:BlockChainService,private router: Router) {
    this.loading = false;
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
  }

  showDetails(block: Block) {
    this.router.navigate(['/block', block.header.seq]);
  }

  handlePageChange(pagesData:number[]){
    this.totalBlocks= pagesData[1];
    let currentPage =pagesData[0];

    let blockStart = this.totalBlocks - currentPage*10 +1;
    let blockEnd = blockStart +9;

    if(blockEnd>=this.totalBlocks){
      blockEnd = this.totalBlocks;
    }

    if(blockStart <=1 ){
      blockStart = 1;
    }
    this.loading = true;

    this.blockService.getBlocks(blockStart,blockEnd).subscribe(
      (data)=>{
        let newData= _.sortBy(data, function (block) {return block.header.seq}).reverse();
        this.blocks= newData;
        this.loading = false;
      },(err)=>{
        this.loading = false;
      }
    );
  }

}
