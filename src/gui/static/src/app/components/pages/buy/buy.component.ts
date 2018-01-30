import { Component } from '@angular/core';
import { PurchaseService } from '../../../services/purchase.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { AddDepositAddressComponent } from './add-deposit-address/add-deposit-address.component';
import { config } from '../../../app.config';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.css']
})
export class BuyComponent {

  otcEnabled: boolean;
  scanning = false;

  constructor(
    public purchaseService: PurchaseService,
    private dialog: MatDialog,
  ) {
    this.otcEnabled = config.otcEnabled;
  }

  addDepositAddress() {
    const config = new MatDialogConfig();
    config.width = '500px';
    this.dialog.open(AddDepositAddressComponent, config);
  }

  searchDepositAddress(address: string) {
    this.scanning = true;
    this.purchaseService.scan(address).subscribe(() => {
      this.disableScanning();
    }, () => {
      this.disableScanning();
    });
  }

  private disableScanning()
  {
    setTimeout(() => this.scanning = false, 1000);
  }
}
