import { Component, OnDestroy, ViewChild } from '@angular/core';
import { ISubscription } from 'rxjs/Subscription';
import { ButtonComponent } from '../../layout/button/button.component';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Params, ActivatedRoute, Router } from '@angular/router';
import { Observable } from 'rxjs/Observable';
import { WalletService } from '../../../services/wallet.service';
import { Wallet } from '../../../app.datatypes';
import { MatSnackBar } from '@angular/material';
import { showSnackbarError } from '../../../utils/errors';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnDestroy {
  @ViewChild('resetButton') resetButton: ButtonComponent;

  form: FormGroup;

  private subscription: ISubscription;
  private wallet: Wallet;
  private done = false;

  constructor(
    public formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private walletService: WalletService,
    private snackbar: MatSnackBar,
  ) {
    this.initForm('');
    this.subscription = Observable.zip(this.route.params, this.walletService.all(), (params: Params, wallets: Wallet[]) => {
      const wallet = wallets.find(w => w.filename === params['id']);
      if (!wallet) {
        setTimeout(() => this.router.navigate([''], {skipLocationChange: true}));

        return;
      }

      this.wallet = wallet;
      this.initForm(wallet.label);
    }).subscribe();
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    this.snackbar.dismiss();
  }

  initForm(walletName: string) {
    const validators = [];
    validators.push(this.passwordMatchValidator.bind(this));

    this.form = new FormGroup({}, validators);
    this.form.addControl('wallet', new FormControl(walletName));
    this.form.addControl('seed', new FormControl('', [Validators.required]));
    this.form.addControl('password', new FormControl());
    this.form.addControl('confirm', new FormControl());
  }

  reset() {
    if (!this.form.valid || this.resetButton.isLoading() || this.done) {
      return;
    }

    this.snackbar.dismiss();
    this.resetButton.setLoading();

    this.walletService.resetPassword(this.wallet, this.form.value.seed, this.form.value.password !== '' ? this.form.value.password : null)
      .subscribe(() => {
        this.resetButton.setSuccess();
        this.resetButton.setDisabled();
        this.done = true;

        setTimeout(() => {
          this.router.navigate(['']);
        }, 2000);
      }, error => {
        this.resetButton.setError(error);
        showSnackbarError(this.snackbar, error);
      });
  }

  private passwordMatchValidator() {
    if (this.form && this.form.get('password') && this.form.get('confirm')) {
      return this.form.get('password').value === this.form.get('confirm').value ? null : { NotEqual: true };
    } else {
      return { NotEqual: true };
    }
  }
}
