import { Component, EventEmitter, Input, OnInit, Output, ViewChild } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { OnboardingSafeguardComponent } from './onboarding-safeguard/onboarding-safeguard.component';
import { MatDialogRef } from '@angular/material';
import { CreateWalletFormComponent } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { HwOptionsDialogComponent } from '../../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { Router } from '@angular/router';
import { HwWalletService } from '../../../../services/hw-wallet.service';

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
  hwCompatibilityActivated = false;

  constructor(
    private dialog: MatDialog,
    private router: Router,
    hwWalletService: HwWalletService,
  ) {
    this.hwCompatibilityActivated = hwWalletService.hwWalletCompatibilityActivated;
  }

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

  useHardwareWallet() {
    const config = new MatDialogConfig();
    config.width = '566px';
    config.autoFocus = false;
    config.data = true;
    this.dialog.open(HwOptionsDialogComponent, config).afterClosed().subscribe(result => {
      if (result) {
        this.router.navigate(['/wallets']);
      }
    });
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
