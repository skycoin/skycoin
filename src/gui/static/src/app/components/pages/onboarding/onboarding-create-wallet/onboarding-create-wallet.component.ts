import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormBuilder, FormControl, FormGroup, Validators } from '@angular/forms';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { WalletService } from '../../../../services/wallet.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { OnboardingSafeguardComponent } from './onboarding-safeguard/onboarding-safeguard.component';
import { MatDialogRef } from '@angular/material';

@Component({
  selector: 'app-onboarding-create-wallet',
  templateUrl: './onboarding-create-wallet.component.html',
  styleUrls: ['./onboarding-create-wallet.component.scss'],
})
export class OnboardingCreateWalletComponent implements OnInit {
  @Input() fill = null;
  @Output() onLabelAndSeedCreated = new EventEmitter<[string, string]>();
  form: FormGroup;
  doubleButtonActive = DoubleButtonActive.LeftButton;

  constructor(
    private dialog: MatDialog,
    private walletService: WalletService,
    private formBuilder: FormBuilder,
  ) { }

  ngOnInit() {
    this.initForm();
  }

  initForm() {
    this.form = this.formBuilder.group({
        label: new FormControl('', Validators.compose([
          Validators.required, Validators.minLength(2),
        ])),
        seed: new FormControl('', Validators.compose([
          Validators.required, Validators.minLength(2),
        ])),
        confirm_seed: new FormControl('',
          Validators.compose(this.showCreateForm ? [Validators.required, Validators.minLength(2)] : [])
        ),
      },
      this.showCreateForm ? { validator: this.seedMatchValidator.bind(this) } : {},
    );

    if (this.fill) {
      this.form.get('label').setValue(this.fill['label']);
      this.form.get('seed').setValue(this.fill['seed']);
      this.form.get('confirm_seed').setValue(this.fill['seed']);
      this.doubleButtonActive = this.fill['create'] ? DoubleButtonActive.LeftButton : DoubleButtonActive.RightButton;
    } else if (this.showCreateForm) {
      this.generateSeed();
    }
  }

  changeForm(newState) {
    this.doubleButtonActive = newState;
    this.fill = null;
    this.initForm();
  }

  createWallet() {
    this.showSafe().afterClosed().subscribe(result => {
      if (result) {
        this.emitCreatedData();
      }
    });
  }

  loadWallet() {
    this.emitCreatedData();
  }

  generateSeed() {
    this.walletService.generateSeed().subscribe(seed => {
      this.form.get('seed').setValue(seed);
    });
  }

  get showCreateForm() {
    return this.doubleButtonActive === DoubleButtonActive.LeftButton;
  }

  private emitCreatedData() {
    this.onLabelAndSeedCreated.emit([
      this.form.get('label').value,
      this.form.get('seed').value,
      this.doubleButtonActive === DoubleButtonActive.LeftButton,
    ]);
  }

  private seedMatchValidator(g: FormGroup) {
    return g.get('seed').value === g.get('confirm_seed').value
      ? null : { mismatch: true };
  }

  private showSafe(): MatDialogRef<OnboardingSafeguardComponent> {
    const config = new MatDialogConfig();
    config.width = '450px';
    return this.dialog.open(OnboardingSafeguardComponent, config);
  }
}
