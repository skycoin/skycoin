import { Component, OnDestroy, ViewChild } from '@angular/core';
import { SubscriptionLike, zip } from 'rxjs';
import { ButtonComponent } from '../../layout/button/button.component';
import { FormGroup, FormBuilder, FormControl, Validators } from '@angular/forms';
import { Params, ActivatedRoute, Router } from '@angular/router';
import { WalletService } from '../../../services/wallet.service';
import { Wallet } from '../../../app.datatypes';
import { MsgBarService } from '../../../services/msg-bar.service';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.scss'],
})
export class ResetPasswordComponent implements OnDestroy {
  @ViewChild('resetButton', { static: false }) resetButton: ButtonComponent;

  form: FormGroup;
  busy = false;

  private subscription: SubscriptionLike;
  private wallet: Wallet;
  private done = false;
  private hideBarWhenClosing = true;

  constructor(
    public formBuilder: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private walletService: WalletService,
    private msgBarService: MsgBarService,
  ) {
    this.initForm('');
    this.subscription = zip(this.route.params, this.walletService.all(), (params: Params, wallets: Wallet[]) => {
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
    if (this.hideBarWhenClosing) {
      this.msgBarService.hide();
    }
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
    if (!this.form.valid || this.busy || this.done) {
      return;
    }

    this.busy = true;
    this.msgBarService.hide();
    this.resetButton.setLoading();

    this.walletService.resetPassword(this.wallet, this.form.value.seed, this.form.value.password !== '' ? this.form.value.password : null)
      .subscribe(() => {
        this.resetButton.setSuccess();
        this.resetButton.setDisabled();
        this.done = true;

        this.hideBarWhenClosing = false;
        this.msgBarService.showDone('reset.done');

        setTimeout(() => {
          this.router.navigate(['']);
        }, 2000);
      }, error => {
        this.busy = false;
        this.resetButton.resetState();
        this.msgBarService.showError(error);
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
