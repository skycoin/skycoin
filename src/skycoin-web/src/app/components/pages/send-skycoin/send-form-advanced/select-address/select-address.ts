import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

import { Wallet } from '../../../../../app.datatypes';
import { WalletService } from '../../../../../services/wallet/wallet.service';
import { BaseCoin } from '../../../../../coins/basecoin';
import { CoinService } from '../../../../../services/coin.service';

@Component({
  selector: 'app-select-address',
  templateUrl: './select-address.html',
  styleUrls: ['./select-address.scss'],
})
export class SelectAddressComponent {

  wallets: Wallet[] = [];
  currentCoin: BaseCoin;

  constructor(
    public dialogRef: MatDialogRef<SelectAddressComponent>,
    public walletService: WalletService,
    private coinService: CoinService,
  ) {
    this.walletService.currentWallets.first().subscribe(wallets => this.wallets = wallets);

    this.coinService.currentCoin.first().subscribe((coin: BaseCoin) => this.currentCoin = coin);
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value: string) {
    this.dialogRef.close(value);
  }
}
