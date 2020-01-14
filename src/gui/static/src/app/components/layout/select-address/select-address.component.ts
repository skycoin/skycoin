import { Component } from '@angular/core';
import { MatDialogRef, MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { WalletService } from '../../../services/wallet.service';
import { first } from 'rxjs/operators';
import { AppConfig } from '../../../app.config';
import BigNumber from 'bignumber.js';

class ListElement {
  label: string;
  addresses: ElementAddress[] = [];
}

class ElementAddress {
  address: string;
  coins: BigNumber;
  hours: BigNumber;
}

@Component({
  selector: 'app-select-address',
  templateUrl: './select-address.component.html',
  styleUrls: ['./select-address.component.scss'],
})
export class SelectAddressComponent {

  listElements: ListElement[] = [];

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
    this.walletService.all().pipe(first()).subscribe(wallets => {
      wallets.forEach(wallet => {
        const element = new ListElement();
        element.label = wallet.label;

        wallet.addresses.forEach(address => {
          if (!wallet.isHardware || address.confirmed) {
            element.addresses.push({
              address: address.address,
              coins: address.coins,
              hours: address.hours,
            });
          }
        });

        this.listElements.push(element);
      });
    });
  }

  closePopup() {
    this.dialogRef.close();
  }

  select(value: string) {
    this.dialogRef.close(value);
  }
}
