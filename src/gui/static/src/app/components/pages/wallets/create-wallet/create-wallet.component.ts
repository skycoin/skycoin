import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormControl, FormGroup, Validators } from '@angular/forms';
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
    this.walletService.create(this.form.value.label, this.form.value.seed, this.scan)
      .subscribe(() => this.dialogRef.close());
  }

  generateSeed(entropy: number) {
    this.walletService.generateSeed(entropy).subscribe(seed => this.form.controls.seed.setValue(seed));
  }

  private initForm() {
    this.form = new FormGroup({}, this.validateSeedsAreEqual);
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('confirm_seed', new FormControl('', [Validators.required]));

    this.generateSeed(128);

    this.scan = 100;
  }

  private validateSeedsAreEqual(control: AbstractControl) {
    return control && control.get('seed') && control.get('confirm_seed')
      && control.get('seed').value !== control.get('confirm_seed').value
      ? { NotEqual: true }
      : null;
  }
}
