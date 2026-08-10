import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { UntypedFormBuilder, UntypedFormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';

import { PurchaseService } from '../../../../services/purchase.service';
import { WalletService } from '../../../../services/wallet/wallet.service';

@Component({
    selector: 'app-add-deposit-address',
    templateUrl: './add-deposit-address.component.html',
    styleUrls: ['./add-deposit-address.component.css'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class AddDepositAddressComponent implements OnInit {

  form!: UntypedFormGroup;
  addresses!: any[];

  constructor(
    private walletService: WalletService,
    private dialogRef: MatDialogRef<AddDepositAddressComponent>,
    private formBuilder: UntypedFormBuilder,
    private purchaseService: PurchaseService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {}

  ngOnInit() {
    this.initForm();
    this.walletService.addresses.subscribe( (addresses) => {
      this.addresses = addresses;
      this.changeDetectorRef.markForCheck();
    });
  }

  generate() {
    this.purchaseService.generate(this.form.value.address).subscribe(() => {
      this.dialogRef.close();
      this.changeDetectorRef.markForCheck();
    });
  }

  private initForm() {
    this.form = this.formBuilder.group({
      address: ['', Validators.required],
    });
  }
}
