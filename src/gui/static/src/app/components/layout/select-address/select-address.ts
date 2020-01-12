import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { WalletService } from '../../../services/wallet.service';
import { Wallet } from '../../../app.datatypes';
import { first } from 'rxjs/operators';
import { AppConfig } from '../../../app.config';

@Component({
  selector: 'app-select-address',
  templateUrl: './select-address.html',
  styleUrls: ['./select-address.scss'],
})
export class SelectAddressComponent {

  wallets: Wallet[] = [];

  public static openDialog(dialog: MatDialog): MatDialogRef<SelectAddressComponent, any> {
    const config = new MatDialogConfig();
    config.autoFocus = false;
    config.width = AppConfig.mediumModalWidth;

    return dialog.open(SelectAddressComponent, config);
  }

  constructor(
    public dialogRef: MatDialogRef<SelectAddressComponent>,
    public walletService: WalletService,
  ) {
    this.walletService.all().pipe(first()).subscribe(wallets => this.wallets = wallets);
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value: string) {
    this.dialogRef.close(value);
  }
}
