import { Component, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { SubscriptionLike } from 'rxjs';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

import { NavBarSwitchService } from '../../../services/nav-bar-switch.service';
import { DoubleButtonActive } from '../../layout/double-button/double-button.component';
import { SignRawTxComponent } from './offline-dialogs/implementations/sign-raw-tx.component';
import { BroadcastRawTxComponent } from './offline-dialogs/implementations/broadcast-raw-tx.component';
import { SendCoinsData } from './send-coins-form/send-coins-form.component';

/**
 * Shows the form which allows the user to send coins.
 */
@Component({
  selector: 'app-send-skycoin',
  templateUrl: './send-skycoin.component.html',
  styleUrls: ['./send-skycoin.component.scss'],
})
export class SendSkycoinComponent implements OnDestroy {
  // If true, the form for sending coins is shown. If false, the tx preview is shown.
  showForm = true;
  // Saves the last data entered on the form.
  formData: SendCoinsData;
  // If the page must show the simple form (left) or the advanced one (right).
  activeForm: DoubleButtonActive;
  activeForms = DoubleButtonActive;

  private subscription: SubscriptionLike;

  constructor(
    private navBarSwitchService: NavBarSwitchService,
    private changeDetector: ChangeDetectorRef,
    private dialog: MatDialog,
  ) {
    // Show the switch for changing the form and react to its event.
    this.navBarSwitchService.showSwitch('send.simple-form-button', 'send.advanced-form-button', DoubleButtonActive.LeftButton);
    this.subscription = navBarSwitchService.activeComponent.subscribe(value => {
      if (this.activeForm !== value) {
        this.activeForm = value;
        this.formData = null;
      }
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
    this.navBarSwitchService.hideSwitch();
  }

  // Called when the form requests to show the preview. The param includes the data entered
  // on the form.
  onFormSubmitted(data: SendCoinsData) {
    this.formData = data;
    this.showForm = false;
  }

  // Returns from the tx preview to the form.
  onBack(deleteFormData) {
    // Erase the form data if requested.
    if (deleteFormData) {
      this.formData = null;
    }

    this.showForm = true;
    this.changeDetector.detectChanges();
  }

  // Opens the modal window for signing raw unsigned transactions.
  signTransaction() {
    SignRawTxComponent.openDialog(this.dialog);
  }

  // Opens the modal window for sending raw signed transactions.
  broadcastTransaction() {
    BroadcastRawTxComponent.openDialog(this.dialog);
  }
}
