import { Component } from '@angular/core';
import { PurchaseService } from '../../../services/purchase.service';
import { MdDialog, MdDialogConfig } from '@angular/material';
import { AddDepositAddressComponent } from './add-deposit-address/add-deposit-address.component';
import { config } from '../../../app.config';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.css']
})
export class BuyComponent {

  orders = [];
  otcEnabled: boolean;
  scanning = false;

  constructor(
    public purchaseService: PurchaseService,
    private dialog: MdDialog,
  ) {
    this.otcEnabled = config.otcEnabled;
  }

  addDepositAddress() {
    const config = new MdDialogConfig();
    config.width = '500px';
    this.dialog.open(AddDepositAddressComponent, config);
  }

  searchDepositAddress(address: string) {
    this.scanning = true;
    this.purchaseService.scan(address).subscribe(() => {
      this.disableScanning();
    }, error => {
      this.disableScanning();
    });
  }

  private disableScanning()
  {
    setTimeout(() => this.scanning = false, 1000);
  }
}
