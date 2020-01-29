/*
  IMPORTANT: Unused for a long time, it may need changes to work properly.
*/
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { PurchaseService } from '../../../../services/purchase.service';
import { MatDialogRef } from '@angular/material/dialog';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';

@Component({
  selector: 'app-add-deposit-address',
  templateUrl: './add-deposit-address.component.html',
  styleUrls: ['./add-deposit-address.component.css'],
})
export class AddDepositAddressComponent implements OnInit {

  form: FormGroup;

  constructor(
    public walletsAndAddressesService: WalletsAndAddressesService,
    private dialogRef: MatDialogRef<AddDepositAddressComponent>,
    private formBuilder: FormBuilder,
    private purchaseService: PurchaseService,
  ) {}

  ngOnInit() {
    this.initForm();
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
