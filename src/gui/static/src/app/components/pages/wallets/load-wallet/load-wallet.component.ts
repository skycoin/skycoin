import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-load-wallet',
  templateUrl: './load-wallet.component.html',
  styleUrls: ['./load-wallet.component.scss']
})
export class LoadWalletComponent implements OnInit {

  form: FormGroup;
  seed: string;
  scan: Number;

  constructor(
    public dialogRef: MatDialogRef<LoadWalletComponent>,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  closePopup() {
    this.dialogRef.close();
  }

  generateSeed() {
    this.walletService.generateSeed().subscribe(seed => this.form.controls.seed.setValue(seed));
  }

  loadWallet() {
    this.walletService.create(this.form.value.label, this.form.value.seed, this.scan)
      .subscribe(() => this.dialogRef.close());
  }

  private initForm() {
    this.form = new FormGroup({});
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.generateSeed();
    this.scan = 100;
  }
}
