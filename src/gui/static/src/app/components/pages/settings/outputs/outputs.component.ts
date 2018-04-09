import { Component, OnInit } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';

@Component({
  selector: 'app-outputs',
  templateUrl: './outputs.component.html',
  styleUrls: ['./outputs.component.scss']
})
export class OutputsComponent implements OnInit {

  wallets: any[];

  constructor(
    public walletService: WalletService,
  ) { }

  ngOnInit() {
    this.walletService.outputsWithWallets().subscribe(wallets => this.wallets = wallets);
  }
}
