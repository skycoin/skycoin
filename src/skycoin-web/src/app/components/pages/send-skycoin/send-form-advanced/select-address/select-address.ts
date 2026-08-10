import { ChangeDetectionStrategy, ChangeDetectorRef, Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';

import { first } from 'rxjs';

import { Wallet } from '../../../../../app.datatypes';
import { WalletService } from '../../../../../services/wallet/wallet.service';
import { BaseCoin } from '../../../../../coins/basecoin';
import { CoinService } from '../../../../../services/coin.service';

@Component({
    selector: 'app-select-address',
    templateUrl: './select-address.html',
    styleUrls: ['./select-address.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class SelectAddressComponent {

  wallets: Wallet[] = [];
  currentCoin!: BaseCoin;

  constructor(
    public dialogRef: MatDialogRef<SelectAddressComponent>,
    public walletService: WalletService,
    private coinService: CoinService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.walletService.currentWallets.pipe(first()).subscribe(wallets => {
      this.wallets = wallets;
      this.changeDetectorRef.markForCheck();
    });

    this.coinService.currentCoin.pipe(first()).subscribe((coin: BaseCoin) => {
      this.currentCoin = coin;
      this.changeDetectorRef.markForCheck();
    });
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value: string) {
    this.dialogRef.close(value);
  }
}
