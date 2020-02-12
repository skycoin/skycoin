import { Component, EventEmitter, Input, OnInit, Output, ViewChild, OnDestroy } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { CreateWalletFormComponent, WalletFormData } from '../../wallets/create-wallet/create-wallet-form/create-wallet-form.component';
import { HwOptionsDialogComponent } from '../../../layout/hardware-wallet/hw-options-dialog/hw-options-dialog.component';
import { Router } from '@angular/router';
import { HwWalletService } from '../../../../services/hw-wallet.service';
import { SubscriptionLike } from 'rxjs';
import { BlockchainService } from '../../../../services/blockchain.service';
import { ConfirmationParams, ConfirmationComponent, DefaultConfirmationButtons } from '../../../layout/confirmation/confirmation.component';
import { AppService } from '../../../../services/app.service';

@Component({
  selector: 'app-onboarding-create-wallet',
  templateUrl: './onboarding-create-wallet.component.html',
  styleUrls: ['./onboarding-create-wallet.component.scss'],
})
export class OnboardingCreateWalletComponent implements OnInit, OnDestroy {
  @ViewChild('formControl', { static: false }) formControl: CreateWalletFormComponent;
  @Input() fill: WalletFormData = null;
  @Output() onLabelAndSeedCreated = new EventEmitter<WalletFormData>();

  showNewForm = true;
  doubleButtonActive = DoubleButtonActive.LeftButton;
  hwCompatibilityActivated = false;

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
    setTimeout(() => { this.formControl.initForm(null, this.fill); });
    if (this.fill) {
      this.doubleButtonActive = this.fill.creatingNewWallet ? DoubleButtonActive.LeftButton : DoubleButtonActive.RightButton;
      this.showNewForm = this.fill.creatingNewWallet;
    }
  }

  ngOnDestroy() {
    this.blockchainSubscription.unsubscribe();
  }

  changeForm(newState) {
    newState === DoubleButtonActive.RightButton ? this.showNewForm = false : this.showNewForm = true;

    this.doubleButtonActive = newState;
    this.fill = null;
    this.formControl.initForm(this.showNewForm, this.fill);
  }

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

  useHardwareWallet() {
    HwOptionsDialogComponent.openDialog(this.dialog, true).afterClosed().subscribe(result => {
      if (result) {
        this.router.navigate(['/wallets']);
      }
    });
  }

  private emitCreatedData() {
    this.onLabelAndSeedCreated.emit(this.formControl.getData());
  }
}
