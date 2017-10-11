import { Component, OnInit } from '@angular/core';
import { BlockchainService } from '../../../services/blockchain.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-explorer',
  templateUrl: './explorer.component.html',
  styleUrls: ['./explorer.component.css']
})
export class ExplorerComponent implements OnInit {

  blocks: any[];

  constructor(
    public blockchainService: BlockchainService,
    private router: Router,
  ) { }

  ngOnInit() {
    this.blockchainService.blocks().subscribe(data => this.blocks = data);
  }

  onActivate(response) {
    if (response.row && response.row.header) {
      this.router.navigate(['/explorer', response.row.header.seq]);
    }
  }
}
