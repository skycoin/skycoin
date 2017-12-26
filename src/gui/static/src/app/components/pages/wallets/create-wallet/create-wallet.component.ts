import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../../services/wallet.service';
import { MdDialogRef } from '@angular/material';

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss']
})
export class CreateWalletComponent implements OnInit {

  form: FormGroup;
  seed: string;

  constructor(
    public dialogRef: MdDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    this.walletService.create(this.form.value.label, this.form.value.seed)
      .subscribe(() => this.dialogRef.close());
  }

  generateSeed() {
    this.walletService.generateSeed().subscribe(seed => this.form.controls.seed.setValue(seed));
  }

  private initForm() {
    this.form = new FormGroup({});
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('confirm_seed', new FormControl('', [
      Validators.compose([Validators.required, this.validateAreEqual.bind(this)])
    ]));

    this.generateSeed();
  }

  private validateAreEqual(fieldControl: FormControl){
    return fieldControl.value === this.form.get('seed').value ? null : { NotEqual: true };
  }
}
