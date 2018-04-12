import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss']
})
export class CreateWalletComponent implements OnInit {

  form: FormGroup;
  seed: string;
  scan: Number;

  constructor(
    public dialogRef: MatDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    this.walletService.create(this.form.value.label, this.form.value.seed, this.scan, this.form.value.password)
      .subscribe(() => this.dialogRef.close());
  }

  generateSeed() {
    this.walletService.generateSeed().subscribe(seed => this.form.get('seed').setValue(seed));
  }

  private initForm() {
    this.form = new FormGroup({}, [this.validatePassword.bind(this)]);
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('confirm_seed', new FormControl('', [
      Validators.compose([Validators.required, this.validateAreEqual.bind(this)])
    ]));
    this.form.addControl('password', new FormControl('', []));
    this.form.addControl('confirm_password', new FormControl('', []));

    this.generateSeed();

    this.scan = 100;
  }

  private validateAreEqual(fieldControl: FormControl) {
    return fieldControl.value.trim() === this.form.get('seed').value.trim() ? null : { NotEqual: true };
  }

  private validatePassword() {
    if (this.form && this.form.get('password') && this.form.get('confirm_password')) {
      if (this.form.get('password').value) {
        if (this.form.get('password').value !== this.form.get('confirm_password').value) {
          return { NotEqual: true };
        }
      }
    }

    return null;
  }
}
