import { Component, OnInit } from '@angular/core';
import { BlockchainService } from '../../../services/blockchain.service';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.css']
})
export class BlockComponent implements OnInit {

  block: any;

  constructor(
    private blockchainService: BlockchainService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.params.switchMap(params => this.blockchainService.block(params.block)).subscribe(response => {
      this.block = response;
    });
  }
}
