import { Component } from '@angular/core';
import { WalletService } from './services/wallet.service';
import { BlockchainService } from './services/blockchain.service';
import 'rxjs/add/operator/takeWhile';
import { ApiService } from './services/api.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {

  current: number;
  highest: number;
  version: string;

  constructor(
    public walletService: WalletService,
    private apiService: ApiService,
    private blockchainService: BlockchainService,
  ) {}

  ngOnInit() {
    this.setVersion();
  }

  private setVersion() {
    return this.apiService.get('version')
      .subscribe(output => this.version = output.version);
  }
}
