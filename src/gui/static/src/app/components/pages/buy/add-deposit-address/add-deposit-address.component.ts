/*
  IMPORTANT: Unused for a long time, it may need changes to work properly.
*/
import { Component, OnInit, OnDestroy } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { PurchaseService } from '../../../../services/purchase.service';
import { MatDialogRef } from '@angular/material/dialog';
import { WalletsAndAddressesService } from '../../../../services/wallet-operations/wallets-and-addresses.service';
import { Subscription } from 'rxjs';
import { AddressBase } from '../../../../services/wallet-operations/wallet-objects';

@Component({
  selector: 'app-add-deposit-address',
  templateUrl: './add-deposit-address.component.html',
  styleUrls: ['./add-deposit-address.component.css'],
})
export class AddDepositAddressComponent implements OnInit, OnDestroy {

  form: UntypedFormGroup;
  addresses: AddressBase[] = [];

  private getWalletsSubscription: Subscription;

  constructor(
    public walletsAndAddressesService: WalletsAndAddressesService,
    private dialogRef: MatDialogRef<AddDepositAddressComponent>,
    private formBuilder: UntypedFormBuilder,
    private purchaseService: PurchaseService,
  ) {}

  ngOnInit() {
    this.initForm();

    this.getWalletsSubscription = this.walletsAndAddressesService.allWallets.subscribe(wallets => {
      this.addresses = wallets.reduce((array, wallet) => array.concat(wallet.addresses), [] as AddressBase[]);
    });
  }

  ngOnDestroy() {
    this.getWalletsSubscription.unsubscribe();
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
