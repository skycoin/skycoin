import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { OnboardingSafeguardComponent } from './onboarding-safeguard/onboarding-safeguard.component';
import { MatDialogRef } from '@angular/material';
import { CreateWalletFormComponent } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';

@Component({
  selector: 'app-onboarding-create-wallet',
  templateUrl: './onboarding-create-wallet.component.html',
  styleUrls: ['./onboarding-create-wallet.component.scss'],
})
export class OnboardingCreateWalletComponent implements OnInit {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  @Input() fill = null;
  @Output() onLabelAndSeedCreated = new EventEmitter<[string, string, boolean]>();

  showNewForm = true;
  doubleButtonActive = DoubleButtonActive.LeftButton;

  constructor(
    private dialog: MatDialog,
  ) { }

  ngOnInit() {
    setTimeout(() => { this.formControl.initForm(null, this.fill); });
    if (this.fill) {
      this.doubleButtonActive = this.fill['create'] ? DoubleButtonActive.LeftButton : DoubleButtonActive.RightButton;
      this.showNewForm = this.fill['create'];
    }
  }

  changeForm(newState) {
    newState === DoubleButtonActive.RightButton ? this.showNewForm = false : this.showNewForm = true;

    this.doubleButtonActive = newState;
    this.fill = null;
    this.formControl.initForm(this.showNewForm, this.fill);
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

  private emitCreatedData() {

    const data = this.formControl.getData();

    this.onLabelAndSeedCreated.emit([
      data.label,
      data.seed,
      this.doubleButtonActive === DoubleButtonActive.LeftButton,
    ]);
  }

  private showSafe(): MatDialogRef<OnboardingSafeguardComponent> {
    const config = new MatDialogConfig();
    config.width = '450px';

    return this.dialog.open(OnboardingSafeguardComponent, config);
  }
}
