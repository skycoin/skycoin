import { Component, Input, OnDestroy } from '@angular/core';
import { PriceService } from '../../../../../services/price.service';
import { SubscriptionLike } from 'rxjs';
import { BigNumber } from 'bignumber.js';
import { MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { ChangeNoteComponent } from './change-note/change-note.component';
import { GeneratedTransaction, OldTransaction } from '../../../../../services/wallet-operations/transaction-objects';

@Component({
  selector: 'app-transaction-info',
  templateUrl: './transaction-info.component.html',
  styleUrls: ['./transaction-info.component.scss'],
})
export class TransactionInfoComponent implements OnDestroy {
  @Input() transaction: GeneratedTransaction|OldTransaction;
  @Input() isPreview: boolean;
  price: number;
  showInputsOutputs = false;

  private subscription: SubscriptionLike;

  constructor(private priceService: PriceService, private dialog: MatDialog) {
    this.subscription = this.priceService.price.subscribe(price => this.price = price);
  }

  get hoursText(): string {
    if (!this.transaction) {
      return '';
    }

    if (!this.isPreview) {
      if ((this.transaction as any).coinsMovedInternally) {
        return 'tx.hours-moved';
      } else if ((this.transaction as OldTransaction).balance.isGreaterThan(0)) {
        return 'tx.hours-received';
      }
    }

    return 'tx.hours-sent';
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  toggleInputsOutputs(event) {
    event.preventDefault();

    this.showInputsOutputs = !this.showInputsOutputs;
  }

  editNote() {
    ChangeNoteComponent.openDialog(this.dialog, this.transaction as OldTransaction).afterClosed().subscribe(newNote => {
      if (newNote || newNote === '') {
        this.transaction.note = newNote;
      }
    });
  }
}
