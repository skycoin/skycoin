import { Component, OnInit, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { WalletService } from '../../../../services/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';

@Component({
  selector: 'app-onboarding-encrypt-wallet',
  templateUrl: './onboarding-encrypt-wallet.component.html',
  styleUrls: ['./onboarding-encrypt-wallet.component.scss'],
})
export class OnboardingEncryptWalletComponent implements OnInit {
  @ViewChild('button') button: ButtonComponent;
  form: FormGroup;

  constructor(
    private formBuilder: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
    private walletService: WalletService,
) { }

  ngOnInit() {
    this.initEncryptForm();
  }

  initEncryptForm() {
    this.form = this.formBuilder.group({
        password: new FormControl('', Validators.compose([Validators.required, Validators.minLength(2)])),
        confirm: new FormControl('',
          Validators.compose([
            Validators.required,
            Validators.minLength(2),
          ]),
        ),
      },
      {
        validator: this.passwordMatchValidator.bind(this),
      });

    this.form.disable();
  }

  setEncrypt(event) {
    event.checked ? this.form.enable() : this.form.disable();
  }

  encryptWallet() {
    this.button.setLoading();

    this.walletService.find(this.route.snapshot.queryParams['wallet']).subscribe(wallet => {
      // TODO: encrypt the wallet
      this.skip();
    });
  }

  skip() {
    this.router.navigate(['/wallets']);
  }

  private passwordMatchValidator(g: FormGroup) {
    return g.get('password').value === g.get('confirm').value
      ? null : { mismatch: true };
  }
}
