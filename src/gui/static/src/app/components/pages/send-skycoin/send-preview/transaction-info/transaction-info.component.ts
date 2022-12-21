import { Component, Input, OnDestroy } from '@angular/core';
import { PriceService } from '../../../../../services/price.service';
import { SubscriptionLike } from 'rxjs';
import { BigNumber } from 'bignumber.js';
import { MatLegacyDialog as MatDialog } from '@angular/material/legacy-dialog';

import { ChangeNoteComponent } from './change-note/change-note.component';
import { GeneratedTransaction, OldTransaction } from '../../../../../services/wallet-operations/transaction-objects';

/**
 * Allows to view the details of a transaction which is about to be sent or a transaction
 * from the history.
 */
@Component({
  selector: 'app-transaction-info',
  templateUrl: './transaction-info.component.html',
  styleUrls: ['./transaction-info.component.scss'],
})
export class TransactionInfoComponent implements OnDestroy {
  // Transaction which is going to be shown.
  @Input() transaction: GeneratedTransaction|OldTransaction;
  // True if the provided transaction was created to be sent, false if it is from the history.
  @Input() isPreview: boolean;
  @Input() isForOfflineTransaction: boolean;
  // Current price per coin, in usd.
  price: number;
  showInputsOutputs = false;

  private subscription: SubscriptionLike;

  constructor(private priceService: PriceService, private dialog: MatDialog) {
    this.subscription = this.priceService.price.subscribe(price => this.price = price);
  }

  // Gets the text which says what was done with the moved coins (if were received, sent or moved).
  get hoursText(): string {
    if (!this.isPreview) {
      if ((this.transaction as OldTransaction).coinsMovedInternally) {
        return 'tx.hours-moved';
      } else if ((this.transaction as OldTransaction).balance.isGreaterThan(0)) {
        return 'tx.hours-received';
      }
    }

    return 'tx.hours-sent';
  }

  // Gets the amount of moved hours.
  get sentOrReceivedHours(): BigNumber {
    return this.isPreview ?
      (this.transaction as GeneratedTransaction).hoursToSend :
      (this.transaction as OldTransaction).hoursBalance;
  }

  // File to be used as transaction icon in the UI.
  get transactionIcon(): string {
    if ((this.transaction as OldTransaction).coinsMovedInternally) {
      return 'moved-grey.png';
    } else if (this.isPreview || (this.transaction as OldTransaction).balance.isLessThan(0)) {
      return 'sent-blue.png';
    }

    return 'received-blue.png';
  }

  // Returns how many coins were moved.
  get balanceToShow(): BigNumber {
    return this.isPreview ?
      (this.transaction as GeneratedTransaction).coinsToSend :
      (this.transaction as OldTransaction).balance;
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  // Makes visible the list of the inputs and outputs.
  toggleInputsOutputs(event) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }

  // Opens the modal window for editing the note of the transaction.
  editNote() {
    ChangeNoteComponent.openDialog(this.dialog, this.transaction as OldTransaction).afterClosed().subscribe(newNote => {
      if (newNote || newNote === '') {
        this.transaction.note = newNote;
      }
    });
  }
}
