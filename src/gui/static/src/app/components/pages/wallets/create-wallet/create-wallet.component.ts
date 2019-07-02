import { Component, Inject, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { WalletService } from '../../../../services/wallet.service';
import { MatDialogRef, MatDialog } from '@angular/material/dialog';
import { ButtonComponent } from '../../../layout/button/button.component';
import { MAT_DIALOG_DATA } from '@angular/material';
import { CreateWalletFormComponent } from './create-wallet-form/create-wallet-form.component';
import { ISubscription } from 'rxjs/Subscription';
import { BlockchainService } from '../../../../services/blockchain.service';
import { ConfirmationData } from '../../../../app.datatypes';
import { showConfirmationModal } from '../../../../utils';

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.scss'],
})
export class CreateWalletComponent implements OnDestroy {
  @ViewChild('formControl') formControl: CreateWalletFormComponent;
  @ViewChild('createButton') createButton: ButtonComponent;
  @ViewChild('cancelButton') cancelButton: ButtonComponent;

  scan: Number;
  disableDismiss = false;

  private synchronized = true;
  private synchronizedSubscription: ISubscription;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data,
    public dialogRef: MatDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
    private dialog: MatDialog,
    blockchainService: BlockchainService,
  ) {
    this.synchronizedSubscription = blockchainService.synchronized.subscribe(value => this.synchronized = value);
  }

  ngOnDestroy() {
    this.synchronizedSubscription.unsubscribe();
  }

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    if (!this.formControl.isValid || this.createButton.isLoading()) {
      return;
    }

    if (this.synchronized || this.data.create) {
      this.continueCreating();
    } else {
      const confirmationData: ConfirmationData = {
        headerText: 'wallet.new.synchronizing-warning-title',
        text: 'wallet.new.synchronizing-warning-text',
        confirmButtonText: 'wallet.new.synchronizing-warning-continue',
        cancelButtonText: 'wallet.new.synchronizing-warning-cancel',
      };

      showConfirmationModal(this.dialog, confirmationData).afterClosed().subscribe(confirmationResult => {
        if (confirmationResult) {
          this.continueCreating();
        }
      });
    }
  }

  private continueCreating() {
    const data = this.formControl.getData();

    this.createButton.resetState();
    this.createButton.setLoading();
    this.cancelButton.setDisabled();
    this.disableDismiss = true;

    this.walletService.create(data.label, data.seed, this.scan, data.password)
      .subscribe(() => this.dialogRef.close(), e => {
        this.createButton.setError(e);
        this.cancelButton.disabled = false;
        this.disableDismiss = false;
      });
  }
}
