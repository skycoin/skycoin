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
    this.walletService.create(this.form.value.label, this.form.value.seed, this.scan, this.form.password.value)
      .subscribe(
        () => this.dialogRef.close(),
        error => this.snackbar.open(error['_body'], null, { duration: 5000 })
      );
  }

  private initForm() {
    this.form = new FormGroup({}, [this.validatePassword.bind(this)]);
    this.form.addControl('label', new FormControl('', [Validators.required]));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('password', new FormControl('', []));
    this.form.addControl('confirm_password', new FormControl('', []));
    this.scan = 100;
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
