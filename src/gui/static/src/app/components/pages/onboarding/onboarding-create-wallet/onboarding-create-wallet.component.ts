import { Component, EventEmitter, Input, OnInit, Output, ViewChild, OnDestroy } from '@angular/core';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { OnboardingSafeguardComponent } from './onboarding-safeguard/onboarding-safeguard.component';
import { MatDialogRef } from '@angular/material';
import { CreateWalletFormComponent, WalletFormData } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { HwOptionsDialogComponent } from '../../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { Router } from '@angular/router';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { ISubscription } from 'rxjs/Subscription';
import { BlockchainService } from '../../../../services/blockchain.service';
import { ConfirmationData } from '../../../../app.datatypes';
import { showConfirmationModal } from '../../../../utils';

@Component({
  selector: 'app-onboarding-create-wallet',
  templateUrl: './onboarding-create-wallet.component.html',
  styleUrls: ['./onboarding-create-wallet.component.scss'],
})
export class OnboardingCreateWalletComponent implements OnInit, OnDestroy {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  @Input() fill: WalletFormData = null;
  @Output() onLabelAndSeedCreated = new EventEmitter<WalletFormData>();

  showNewForm = true;
  doubleButtonActive = DoubleButtonActive.LeftButton;
  hwCompatibilityActivated = false;

  private synchronized = true;
  private synchronizedSubscription: ISubscription;

  constructor(
    private dialog: MatDialog,
    private router: Router,
    hwWalletService: HwWalletService,
    blockchainService: BlockchainService,
  ) {
    this.hwCompatibilityActivated = hwWalletService.hwWalletCompatibilityActivated;
    this.synchronizedSubscription = blockchainService.synchronized.subscribe(value => this.synchronized = value);
  }

  ngOnInit() {
    setTimeout(() => { this.formControl.initForm(null, this.fill); });
    if (this.fill) {
      this.doubleButtonActive = this.fill.creatingNewWallet ? DoubleButtonActive.LeftButton : DoubleButtonActive.RightButton;
      this.showNewForm = this.fill.creatingNewWallet;
    }
  }

  ngOnDestroy() {
    this.synchronizedSubscription.unsubscribe();
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
    if (this.synchronized) {
      this.emitCreatedData();
    } else {
      const confirmationData: ConfirmationData = {
        headerText: 'wallet.new.synchronizing-warning-title',
        text: 'wallet.new.synchronizing-warning-text',
        confirmButtonText: 'wallet.new.synchronizing-warning-continue',
        cancelButtonText: 'wallet.new.synchronizing-warning-cancel',
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.emitCreatedData();
        }
      });
    }
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
    this.onLabelAndSeedCreated.emit(this.formControl.getData());
  }

  private showSafe(): MatDialogRef<OnboardingSafeguardComponent> {
    const config = new MatDialogConfig();
    config.width = '450px';

    return this.dialog.open(OnboardingSafeguardComponent, config);
  }
}
