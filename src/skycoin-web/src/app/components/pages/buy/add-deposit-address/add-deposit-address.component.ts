import { Component, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';

import { PurchaseService } from '../../../../services/purchase.service';
import { WalletService } from '../../../../services/wallet/wallet.service';

@Component({
    selector: 'app-add-deposit-address',
    templateUrl: './add-deposit-address.component.html',
    styleUrls: ['./add-deposit-address.component.css'],
    standalone: false
})
export class AddDepositAddressComponent implements OnInit {

  form: UntypedFormGroup;
  addresses: any[];

  constructor(
    private walletService: WalletService,
    private dialogRef: MatDialogRef<AddDepositAddressComponent>,
    private formBuilder: UntypedFormBuilder,
    private purchaseService: PurchaseService,
  ) {}

  ngOnInit() {
    this.initForm();
    this.walletService.addresses.subscribe( (addresses) => {
      this.addresses = addresses;
    });
  }

  generate() {
    this.purchaseService.generate(this.form.value.address).subscribe(() => this.dialogRef.close());
  }

  private initForm() {
    this.form = this.formBuilder.group({
      address: ['', Validators.required],
    });
  }
}
