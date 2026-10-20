import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';

import { CoinService } from '../../../services/coin.service';
import { DoubleButtonActive } from '../../layout/double-button/double-button.component';
import { NavBarService } from '../../../services/nav-bar.service';
import { BaseCoin } from '../../../coins/basecoin';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { MsgBarService } from '../../../services/msg-bar.service';
import { OfflineTxComponent, OfflineTxMode } from './offline-tx/offline-tx.component';

@Component({
    selector: 'app-send-skycoin',
    templateUrl: './send-skycoin.component.html',
    styleUrls: ['./send-skycoin.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class SendSkycoinComponent implements OnInit, OnDestroy {
  showForm = true;
  restarting = false;
  formData: any;
  activeForm!: DoubleButtonActive;
  activeForms = DoubleButtonActive;
  currentCoin?: BaseCoin;

  private subscription: Subscription;
  private coinSubscription!: Subscription;

  constructor(
    private coinService: CoinService,
    private dialog: CustomMatDialogService,
    private msgBarService: MsgBarService,
    navbarService: NavBarService,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.subscription = navbarService.activeComponent.subscribe(value => {
      this.activeForm = value;
      this.formData = null;
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnInit() {
    this.coinSubscription = this.coinService.currentCoin
      .subscribe(coin => {
        this.currentCoin = coin;
        this.onBack(true);
        this.changeDetectorRef.markForCheck();
      });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    this.coinSubscription.unsubscribe();
  }

  onFormSubmitted(data: any) {
    this.formData = data;
    this.showForm = false;

    this.goUp();
  }

  onBack(deleteFormData: any) {
    if (deleteFormData) {
      this.formData = null;
    }

    this.restarting = true;
    setTimeout(() => {
      this.restarting = false;
      this.changeDetectorRef.markForCheck();
    }, 0);

    this.showForm = true;

    this.goUp();
  }

  signTransaction() {
    OfflineTxComponent.open(this.dialog, { mode: OfflineTxMode.Sign });
  }

  broadcastTransaction() {
    OfflineTxComponent.open(this.dialog, { mode: OfflineTxMode.Broadcast }).afterClosed().subscribe(sent => {
      if (sent) {
        this.msgBarService.showDone('offline-tx.broadcast-done');
      }
    });
  }

  get transaction() {
    const transaction = this.formData.transaction;

    transaction.from = this.formData.form.wallet.label;
    transaction.to = this.formData.to;
    transaction.balance = this.formData.amount;

    return transaction;
  }

  private goUp() {
    window.scrollTo(0, 0);
  }
}
