import { Component, OnInit } from '@angular/core';
import {BlockChainService} from "../block-chain-table/block-chain.service";
import {Block} from "../block-chain-table/block";
import {Router} from "@angular/router";

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
    this.service.getBlockByHash(hashVal).subscribe((block)=>{
      this.block = block;
      this.router.navigate(['/block', block.header.seq]);

    });
  }

}
