import { Component, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { NavBarService } from '../../../services/nav-bar.service';
import { SubscriptionLike } from 'rxjs';
import { DoubleButtonActive } from '../../layout/double-button/double-button.component';
import { MatDialog } from '@angular/material/dialog';
import { SignRawTxComponent } from './offline-dialogs/implementations/sign-raw-tx.component';
import { BroadcastRawTxComponent } from './offline-dialogs/implementations/broadcast-raw-tx.component';
import { SendCoinsFormComponent } from './send-coins-form/send-coins-form.component';

@Component({
  selector: 'app-send-skycoin',
  templateUrl: './send-skycoin.component.html',
  styleUrls: ['./send-skycoin.component.scss'],
})
export class SendSkycoinComponent implements OnDestroy {
  showForm = true;
  formData: any;
  activeForm: DoubleButtonActive;
  activeForms = DoubleButtonActive;

  private subscription: SubscriptionLike;

  constructor(
    private navbarService: NavBarService,
    private changeDetector: ChangeDetectorRef,
    private dialog: MatDialog,
  ) {
    this.navbarService.showSwitch('send.simple-form-button', 'send.advanced-form-button', DoubleButtonActive.LeftButton);
    this.subscription = navbarService.activeComponent.subscribe(value => {
      if (this.activeForm !== value) {
        SendCoinsFormComponent.lastShowForManualUnsignedValue = false;
        this.activeForm = value;
        this.formData = null;
      }
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    this.navbarService.hideSwitch();
  }

  onFormSubmitted(data) {
    this.formData = data;
    this.showForm = false;
  }

  onBack(deleteFormData) {
    if (deleteFormData) {
      this.formData = null;
    }

    this.showForm = true;
    this.changeDetector.detectChanges();
  }

  signTransaction() {
    SignRawTxComponent.openDialog(this.dialog);
  }

  broadcastTransaction() {
    BroadcastRawTxComponent.openDialog(this.dialog);
  }

  get transaction() {
    return this.formData.transaction;
  }
}
