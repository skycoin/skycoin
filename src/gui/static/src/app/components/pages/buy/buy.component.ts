import { Component } from '@angular/core';
import { PurchaseService } from '../../../services/purchase.service';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { AddDepositAddressComponent } from './add-deposit-address/add-deposit-address.component';
import { config } from '../../../app.config';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../services/wallet.service';

@Component({
  selector: 'app-buy',
  templateUrl: './buy.component.html',
  styleUrls: ['./buy.component.scss']
})
export class BuyComponent {

  form: FormGroup;

  constructor(
    public walletService: WalletService,
    private dialog: MatDialog,
    private formBuilder: FormBuilder,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  // addDepositAddress() {
  //   const config = new MatDialogConfig();
  //   config.width = '500px';
  //   this.dialog.open(AddDepositAddressComponent, config);
  // }
  //
  // searchDepositAddress(address: string) {
  //   this.purchaseService.scan(address).subscribe(() => {
  //     this.disableScanning();
  //   }, () => {
  //     this.disableScanning();
  //   });
  // }

  private disableScanning()
  {
    // setTimeout(() => this.scanning = false, 1000);
  }

  private initForm() {
    this.form = this.formBuilder.group({
      address: ['', Validators.required],
    });
  }
}
