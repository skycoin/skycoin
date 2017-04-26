import { Component, OnInit } from '@angular/core';
import {Router, ActivatedRoute, Params} from "@angular/router";
import {Block, Transaction} from "../block-chain-table/block";
import {Observable} from "rxjs";
import {BlockChainService} from "../block-chain-table/block-chain.service";
declare var _: any;

import * as moment from 'moment';

@Component({
  selector: 'app-block-details',
  templateUrl: './block-details.component.html',
  styleUrls: ['./block-details.component.css']
})
export class BlockDetailsComponent implements OnInit {

  private blocksObservable:Observable<Block[]>;

  private block:Block;

  constructor(   private service:BlockChainService,
                  private route: ActivatedRoute,
                  private router: Router) {
    this.block=null;
  }

  ngOnInit() {
    this.blocksObservable= this.route.params
      .switchMap((params: Params) => {
        let selectedBlock = +params['id'];
        return this.service.getBlocks(selectedBlock,selectedBlock);
      });

    this.blocksObservable.subscribe((blocks)=>{
      this.block = blocks[0];
    })

  }

  getTime(time:number){
    return moment.unix(time).format();
  }

  getAmount(block:Block){
    var ret = [];
    let txns:Transaction[] = block.body.txns;
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

}
