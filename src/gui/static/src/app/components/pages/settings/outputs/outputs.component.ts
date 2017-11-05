import { Component, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';

@Component({
  selector: 'app-outputs',
  templateUrl: './outputs.component.html',
  styleUrls: ['./outputs.component.css']
})
export class OutputsComponent implements OnInit {

  outputs: any[];

  constructor(
    public walletService: WalletService,
  ) { }

  ngOnInit() {
    this.walletService.outputs().subscribe(outputs => this.outputs = outputs.head_outputs);
  }
}
