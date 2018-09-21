import { Component } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { WalletService } from '../../../../../services/wallet.service';
import { Wallet } from '../../../../../app.datatypes';

@Component({
  selector: 'app-select-address',
  templateUrl: './select-address.html',
  styleUrls: ['./select-address.scss'],
})
export class SelectAddressComponent {

  wallets: Wallet[] = [];

  constructor(
    public dialogRef: MatDialogRef<SelectAddressComponent>,
    public walletService: WalletService,
  ) {
    this.walletService.all().first().subscribe(wallets => this.wallets = wallets);
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value: string) {
    this.dialogRef.close(value);
  }
}
