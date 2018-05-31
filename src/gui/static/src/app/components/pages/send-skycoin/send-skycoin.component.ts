import { Component } from '@angular/core';

@Component({
  selector: 'app-send-skycoin',
  templateUrl: './send-skycoin.component.html',
  styleUrls: ['./send-skycoin.component.scss'],
})
export class SendSkycoinComponent {
  showForm = true;
  formData: any;

  onFormSubmitted(data) {
    this.formData = data;
    this.showForm = false;
  }

  onBack(deleteFormData) {
    if (deleteFormData) {
      this.formData = null;
    }

    this.showForm = true;
  }

  get transaction() {
    const transaction = this.formData.transaction;

    transaction.from = this.formData.wallet.label;
    transaction.to = this.formData.address;
    transaction.balance = this.formData.amount;

    return transaction;
  }
}
