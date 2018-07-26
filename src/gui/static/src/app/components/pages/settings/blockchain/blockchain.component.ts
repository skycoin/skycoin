import { Component, OnInit } from '@angular/core';
import { BlockchainService } from '../../../../services/blockchain.service';

@Component({
  selector: 'app-blockchain',
  templateUrl: './blockchain.component.html',
  styleUrls: ['./blockchain.component.scss'],
})
export class BlockchainComponent implements OnInit {
  block: any;
  coinSupply: any;

  constructor(
    private blockchainService: BlockchainService,
  ) { }

  ngOnInit() {
    this.blockchainService.lastBlock().subscribe(block => this.block = block);
    this.blockchainService.coinSupply().subscribe(coinSupply => this.coinSupply = coinSupply);
  }
}
