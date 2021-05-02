import { Component, EventEmitter, Input, OnInit, Output, ViewChild, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { SubscriptionLike } from 'rxjs';

import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { CreateWalletFormComponent, WalletFormData } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { HwOptionsDialogComponent } from '../../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { BlockchainService } from '../../../../services/blockchain.service';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../layout/confirmation/confirmation.component';
import { AppService } from '../../../../services/app.service';

/**
 * Shows the first step of the wizard, which allows the user to create a new wallet or load
 * a wallet using a seed.
 */
@Component({
  selector: 'app-onboarding-create-wallet',
  templateUrl: './onboarding-create-wallet.component.html',
  styleUrls: ['./onboarding-create-wallet.component.scss'],
})
export class OnboardingCreateWalletComponent implements OnInit, OnDestroy {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  // Data for filling the form just after loading it.
  @Input() fill: WalletFormData = null;
  // Emits when the user presses the button for going to the next step of the wizard, after
  // filling the form. Includes an object with the data entered on the form.
  @Output() onLabelAndSeedCreated = new EventEmitter<WalletFormData>();

  // Current selection on the double button for choosing if the form must be shown for
  // creating a new wallet (left) or for loading a wallet using a seed (right).
  currentFormSelection = DoubleButtonActive.LeftButton;
  // If the option for adding a hw wallet must be shown.
  hwCompatibilityActivated = false;

  doubleButtonActive = DoubleButtonActive;

  // If the blockchain is synchronized.
  private synchronized = true;
  private blockchainSubscription: SubscriptionLike;

  constructor(
    public appService: AppService,
    private dialog: MatDialog,
    private router: Router,
    hwWalletService: HwWalletService,
    blockchainService: BlockchainService,
  ) {
    this.hwCompatibilityActivated = hwWalletService.hwWalletCompatibilityActivated;
    this.blockchainSubscription = blockchainService.progress.subscribe(response => this.synchronized = response.synchronized);
  }

  ngOnInit() {
    // Fill the form.
    setTimeout(() => { this.formControl.initForm(null, this.fill); });
    // Show the correct form.
    if (this.fill) {
      this.currentFormSelection = this.fill.creatingNewWallet ? DoubleButtonActive.LeftButton : DoubleButtonActive.RightButton;
    }
  }

  ngOnDestroy() {
    this.blockchainSubscription.unsubscribe();
    this.onLabelAndSeedCreated.complete();
  }

  // Changes the form currently shown on the UI.
  changeForm(newState: DoubleButtonActive) {
    this.currentFormSelection = newState;
    // Resets the form.
    this.fill = null;
    this.formControl.initForm(this.currentFormSelection === DoubleButtonActive.LeftButton, this.fill);
  }

  // Shows an alert asking for confirmation and emits an event for going to the next step
  // of the wizard.
  createWallet() {
    const confirmationParams: ConfirmationParams = {
      headerText: 'wizard.confirm.title',
      redTitle: true,
      text: 'wizard.confirm.desc',
      checkboxText: 'wizard.confirm.checkbox',
      confirmButtonText: 'common.continue-button',
    };

    ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
      if (confirmationResult) {
        this.emitCreatedData();
      }
    });
  }

  // Emits an event for going to the next step of the wizard. If the blockchain is not
  // synchronized, it shows an alert first, asking for confirmation.
  loadWallet() {
    if (this.synchronized) {
      this.emitCreatedData();
    } else {
      const confirmationParams: ConfirmationParams = {
        headerText: 'common.warning-title',
        text: 'wallet.new.synchronizing-warning-text',
        defaultButtons: DefaultConfirmationButtons.ContinueCancel,
        redTitle: true,
      };

      ConfirmationComponent.openDialog(this.dialog, confirmationParams).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.emitCreatedData();
        }
      });
    }
  }

  // Opens the hw wallet options modal window, which will try to add the connected hw wallet
  // to the wallet list. If at the end the wallet is on the list, the wizard is closed.
  useHardwareWallet() {
    HwOptionsDialogComponent.openDialog(this.dialog, true).afterClosed().subscribe(result => {
      if (result) {
        this.router.navigate(['/wallets']);
      }
    });
  }

  // Emits an event for going to the next step of the wizard.
  private emitCreatedData() {
    this.onLabelAndSeedCreated.emit(this.formControl.getData());
  }
}
