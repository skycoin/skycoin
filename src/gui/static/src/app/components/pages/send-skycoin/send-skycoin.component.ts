import { Component, OnDestroy, ChangeDetectorRef } from '@angular/core';
import { NavBarService } from '../../../services/nav-bar.service';
import { ISubscription } from 'rxjs/Subscription';
import { DoubleButtonActive } from '../../layout/double-button/double-button.component';

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

  private subscription: ISubscription;

  constructor(
    navbarService: NavBarService,
    private changeDetector: ChangeDetectorRef,
  ) {
    navbarService.setActiveComponent(DoubleButtonActive.LeftButton);
    this.subscription = navbarService.activeComponent.subscribe(value => {
      this.activeForm = value;
      this.formData = null;
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
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

  get transaction() {
    const transaction = this.formData.transaction;

    transaction.wallet = this.formData.form.wallet;
    transaction.from = this.formData.form.wallet.label;
    transaction.to = this.formData.to;
    transaction.balance = this.formData.amount;
    transaction.note = this.formData.form.note;

    return transaction;
  }
}
