import { Component, OnInit } from '@angular/core';
import {BlockChainService} from "../block-chain-table/block-chain.service";
import {Block} from "../block-chain-table/block";
import {Router} from "@angular/router";
import {UxOutputsService} from "../address-detail/UxOutputs.service";
import {isNumber} from "util";

@Component({
  selector: 'app-skycoin-search-bar',
  templateUrl: './skycoin-search-bar.component.html',
  styleUrls: ['./skycoin-search-bar.component.css']
})
export class SkycoinSearchBarComponent implements OnInit {

  private block:Block;

  constructor(private service:BlockChainService,private router: Router) { }

  ngOnInit() {
  }

  searchBlockHistory(hashVal:string){
    if(hashVal.length ==34){
      this.router.navigate(['/address', hashVal]);
      return;
    }
    if(hashVal.length ==64){
      this.router.navigate(['/transaction', hashVal]);
      return;
    }
    this.router.navigate(['/block', hashVal]);
    return;

    // this.service.getBlockByHash(hashVal).subscribe((block)=>{
    //   this.block = block;
    //   this.router.navigate(['/block', block.header.seq]);
    //
    // });
  }

}
