import { Component } from '@angular/core';
import { MdDialog, MdDialogConfig } from '@angular/material';
import { AddDepositAddressComponent } from './add-deposit-address/add-deposit-address.component';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.css']
})
export class BuyComponent {

  orders = [];
  constructor(
    public purchaseService: PurchaseService,
    private dialog: MdDialog,
  ) { }

  addDepositAddress() {
    const config = new MdDialogConfig();
    config.width = '500px';
    this.dialog.open(AddDepositAddressComponent, config);
  }
}
