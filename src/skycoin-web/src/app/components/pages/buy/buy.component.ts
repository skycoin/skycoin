import { Component, OnInit } from '@angular/core';
import { MatDialogConfig } from '@angular/material/dialog';

import { config } from '../../../app.config';
import { PurchaseService } from '../../../services/purchase.service';
import { AddDepositAddressComponent } from './add-deposit-address/add-deposit-address.component';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.css'],
})
export class BuyComponent implements OnInit {

  orders: any[];
  otcEnabled: boolean;
  scanning = false;

  constructor(
    public purchaseService: PurchaseService,
    private dialog: CustomMatDialogService,
  ) {
    this.otcEnabled = config.otcEnabled;
  }

  ngOnInit() {
    this.purchaseService.all().subscribe((orders) => {
      this.orders = orders;
    });
  }

  addDepositAddress() {
    const dialogConfig = new MatDialogConfig();
    dialogConfig.width = '500px';
    this.dialog.open(AddDepositAddressComponent, dialogConfig);
  }

  searchDepositAddress(address: string) {
    this.scanning = true;
    this.purchaseService.scan(address).subscribe(() => {
      this.disableScanning();
    }, error => {
      this.disableScanning();
    });
  }

  private disableScanning() {
    setTimeout(() => this.scanning = false, 1000);
  }
}
