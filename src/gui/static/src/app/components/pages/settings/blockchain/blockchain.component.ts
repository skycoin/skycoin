import { Component, OnInit } from '@angular/core';
import { BlockchainService } from '../../../../services/blockchain.service';

@Component({
  selector: 'app-blockchain',
  templateUrl: './blockchain.component.html',
  styleUrls: ['./blockchain.component.scss']
})
export class BlockchainComponent implements OnInit {

  block: any;

  constructor(
    private blockchainService: BlockchainService,
  ) { }

  ngOnInit() {
    this.blockchainService.lastBlock().subscribe(block => this.block = block);
  }
}
