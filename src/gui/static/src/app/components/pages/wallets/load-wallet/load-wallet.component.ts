import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';

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
    private snackbar: MatSnackBar,
    private walletService: WalletService,
  ) {}

  ngOnInit() {
    this.initForm();
  }

  closePopup() {
    this.dialogRef.close();
  }

  loadWallet() {
    this.walletService.create(this.form.value.label, this.form.value.seed, this.scan)
      .subscribe(
        () => this.dialogRef.close(),
        error => this.snackbar.open(error['_body'], null, { duration: 5000 })
      );
  }

  private initForm() {
    this.form = new FormGroup({});
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.scan = 100;
  }
}
