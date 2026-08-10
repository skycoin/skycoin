import { ChangeDetectionStrategy, ChangeDetectorRef, Component, Inject, OnDestroy, ViewChild } from '@angular/core';
import { MatDialogRef } from '@angular/material/dialog';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { Subscription, of } from 'rxjs';
import { delay } from 'rxjs';

import { WalletService } from '../../../../services/wallet/wallet.service';
import { ButtonComponent } from '../../../layout/button/button.component';
import { CoinService } from '../../../../services/coin.service';
import { BaseCoin } from '../../../../coins/basecoin';
import { CreateWalletFormComponent } from './create-wallet-form/create-wallet-form.component';
import { Wallet } from '../../../../app.datatypes';
import { scanAddresses } from '../../../../utils';
import { BlockchainService } from '../../../../services/blockchain.service';
import { CustomMatDialogService } from '../../../../services/custom-mat-dialog.service';
import { config } from '../../../../app.config';
import { MsgBarService } from '../../../../services/msg-bar.service';

@Component({
    selector: 'app-create-wallet',
    templateUrl: './create-wallet.component.html',
    styleUrls: ['./create-wallet.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class CreateWalletComponent implements OnDestroy {
  @ViewChild('formControl') formControl!: CreateWalletFormComponent;
  @ViewChild('create') createButton!: ButtonComponent;

  showSlowMobileInfo = false;
  disableDismiss = false;

  private slowInfoSubscription!: Subscription;
  // Optional password captured from the create form; when set, the new wallet's
  // seed is encrypted right after it is created (before the dialog closes).
  private pendingPassword = '';

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: any,
    public dialogRef: MatDialogRef<CreateWalletComponent>,
    private walletService: WalletService,
    private coinService: CoinService,
    private blockchainService: BlockchainService,
    private translate: TranslateService,
    private dialog: CustomMatDialogService,
    private msgBarService: MsgBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) { }

  ngOnDestroy() {
    this.removeSlowInfoSubscription();
    this.msgBarService.hide();
  }

  closePopup() {
    this.dialogRef.close();
  }

  createWallet() {
    this.createButton.setLoading();
    this.disableDismiss = true;

    this.slowInfoSubscription = of(1).pipe(delay(config.timeBeforeSlowMobileInfo))
      .subscribe(() => {
        this.showSlowMobileInfo = true;
        this.changeDetectorRef.markForCheck();
      });

    const data = this.formControl.getData();
    this.pendingPassword = this.data.create ? (data.password || '') : '';

    this.walletService.create(data.label, data.seed, data.coin.id, this.data.create, data.walletType, data.seedPassphrase, data.segwit)
      .subscribe(
        wallet => {
          this.onCreateSuccess(wallet, data.coin);
          this.changeDetectorRef.markForCheck();
        },
        (error) => {
          this.onCreateError(error.message);
          this.changeDetectorRef.markForCheck();
        }
      );
  }

  private onCreateSuccess(wallet: Wallet, coin: BaseCoin) {
    const initialCoin = this.coinService.currentCoin.value;
    this.coinService.changeCoin(coin);

    if (!this.data.create) {
      this.showSlowMobileInfo = false;
      this.removeSlowInfoSubscription();

      scanAddresses(this.dialog, wallet, this.blockchainService, this.translate).subscribe(
        response => {
          this.processScanResponse(initialCoin, wallet, false, response);
          this.changeDetectorRef.markForCheck();
        },
        error => {
          this.processScanResponse(initialCoin, wallet, true, error);
          this.changeDetectorRef.markForCheck();
        }
      );
    } else if (this.pendingPassword) {
      // Optional creation-time password: encrypt the freshly-created wallet's
      // seed (setWalletPassword persists it), then finish. On error, surface it
      // but the wallet itself was already created.
      this.walletService.setWalletPassword(wallet, this.pendingPassword).subscribe(
        () => {
          this.finish();
          this.changeDetectorRef.markForCheck();
        },
        err => {
          this.onCreateError(err && err.message ? err.message : String(err));
          this.changeDetectorRef.markForCheck();
        },
      );
    } else {
      this.finish();
    }
  }

  private processScanResponse(initialCoin: BaseCoin, wallet: Wallet, isError: boolean, response: any) {
    if (isError || response !== null) {
      this.coinService.changeCoin(initialCoin);
      this.onCreateError(response.message ? response.message : response.toString());
    } else {
      wallet.needSeedConfirmation = false;
      this.walletService.add(wallet);
      this.finish();
    }
  }

  private finish() {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.createButton.setSuccess();
    this.dialogRef.close();
    setTimeout(() => {
      this.msgBarService.showDone('wallet.new.wallet-created');
      this.changeDetectorRef.markForCheck();
    });
  }

  private onCreateError(errorMesasge: string) {
    this.showSlowMobileInfo = false;
    this.removeSlowInfoSubscription();
    this.msgBarService.showError(errorMesasge);
    this.createButton.resetState();

    this.disableDismiss = false;
  }

  private removeSlowInfoSubscription() {
    if (this.slowInfoSubscription) {
      this.slowInfoSubscription.unsubscribe();
    }
  }
}
