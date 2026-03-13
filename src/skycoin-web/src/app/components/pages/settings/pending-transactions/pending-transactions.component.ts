import { Component, OnDestroy, OnInit } from '@angular/core';
import * as moment from 'moment';
import { Subscription, of } from 'rxjs';
import { delay, mergeMap, first } from 'rxjs/operators';
import { BigNumber } from 'bignumber.js';
import { MatDialog } from '@angular/material/dialog';

import { WalletService } from '../../../../services/wallet/wallet.service';
import { HistoryService } from '../../../../services/wallet/history.service';
import { NavBarService } from '../../../../services/nav-bar.service';
import { DoubleButtonActive } from '../../../layout/double-button/double-button.component';
import { Wallet, ConfirmationData } from '../../../../app.datatypes';
import { Observable, forkJoin } from 'rxjs';
import { BaseCoin } from '../../../../coins/basecoin';
import { CoinService } from '../../../../services/coin.service';
import { GlobalsService } from '../../../../services/globals.service';
import { isEqualOrSuperiorVersion } from '../../../../utils/semver';
import { ConfirmationComponent } from '../../../layout/confirmation/confirmation.component';

@Component({
    selector: 'app-pending-transactions',
    templateUrl: './pending-transactions.component.html',
    styleUrls: ['./pending-transactions.component.scss'],
    standalone: false
})
export class PendingTransactionsComponent implements OnInit, OnDestroy {
  isLoading = false;
  transactions: any[] = [];
  currentCoin: BaseCoin;
  showError = false;

  private navbarSubscription: Subscription;
  private coinSubscription: Subscription;
  private dataSubscription: Subscription;
  private selectedNavbarOption: number;

  private readonly updatePeriod = 10 * 1000;
  private readonly errorUpdatePeriod = 2 * 1000;

  constructor(
    private walletService: WalletService,
    private historyService: HistoryService,
    private navbarService: NavBarService,
    private coinService: CoinService,
    private globalsService: GlobalsService,
    private dialog: MatDialog
  ) { }

  ngOnInit() {
    this.coinSubscription = this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        this.currentCoin = coin;
        this.navbarService.setActiveComponent();
      });

    this.navbarService.showSwitch('pending-txs.my', 'pending-txs.all', DoubleButtonActive.LeftButton);

    this.navbarSubscription = this.navbarService.activeComponent.subscribe(value => {
      this.selectedNavbarOption = value;
      this.transactions = [];
      this.startDataRefreshSubscription(0);
    });
  }

  ngOnDestroy() {
    this.navbarSubscription.unsubscribe();
    this.coinSubscription.unsubscribe();

    this.closeDataSubscription();
    this.navbarService.hideSwitch();
  }

  deleteTransaction(txid: string) {
    const confirmationData: ConfirmationData = {
      text: 'pending-txs.delete-confirm',
      headerText: 'confirmation.header-text',
      confirmButtonText: 'confirmation.confirm-button',
      cancelButtonText: 'confirmation.cancel-button',
      redTitle: true,
    };

    const dialogRef = this.dialog.open(ConfirmationComponent, {
      width: '450px',
      data: confirmationData,
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.historyService.deletePendingTransaction(txid).subscribe(() => {
          this.transactions = [];
          this.startDataRefreshSubscription(0);
        });
      }
    });
  }

  private startDataRefreshSubscription(delayMs: number) {
    this.closeDataSubscription();

    this.dataSubscription = of(0).pipe(
      delay(delayMs),
      mergeMap(() => this.loadTransactions(this.selectedNavbarOption))
    ).subscribe(transactions => {
      this.transactions = transactions;
      this.isLoading = false;
      this.startDataRefreshSubscription(this.updatePeriod);
    }, () => {
      this.showError = true;
      this.startDataRefreshSubscription(this.errorUpdatePeriod);
    });
  }

  private loadTransactions(value: number): Observable<any[]> {
    this.isLoading = true;
    this.showError = false;

    const showAllTransactions = value === DoubleButtonActive.RightButton;

    return this.historyService.getAllPendingTransactions().pipe(
      delay(32),
      mergeMap((transactions: any) => {
        return showAllTransactions ? of(transactions) : this.getWalletsTransactions(transactions);
      }),
      mergeMap(transactions => of(this.mapTransactions(transactions)))
    );
  }

  private mapTransactions(transactions) {
    return transactions.map(transaction => {
      transaction.transaction.timestamp = moment(transaction.received).unix();
      return transaction.transaction;
    })
      .map(transaction => {
        transaction.amount = new BigNumber('0');
        transaction.hours = new BigNumber('0');
        transaction.outputs.map(output => {
          transaction.amount = transaction.amount.plus(output.coins);
          transaction.hours = transaction.hours.plus(output.hours);
        });

        return transaction;
      });
  }

  private getWalletsTransactions(transactions: any): Observable<any> {
    if (transactions.length === 0) {
      return of([]);
    }

    return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
      let allTransactions: Observable<any>;
      if (isEqualOrSuperiorVersion(version, '0.25.0')) {
        allTransactions = of(transactions);
      } else {
        allTransactions = this.getUpdatedTransactions(transactions);
      }

      return forkJoin([allTransactions, this.walletService.currentWallets.pipe(first())]).pipe(
        mergeMap(([trans, wallets]: [any, Wallet[]]) => {
          const walletAddresses = new Set<string>();
          wallets.forEach(wallet => {
            wallet.addresses.forEach(address => walletAddresses.add(address.address));
          });

          return of(trans.filter(tran => {
            if (isEqualOrSuperiorVersion(version, '0.25.0')) {
              return tran.transaction.inputs.some(input => walletAddresses.has(input.owner)) ||
              tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
            } else {
              return tran.owner_addressses.some(address => walletAddresses.has(address)) ||
              tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
            }
          }));
        })
      );
    }));
  }

  private getUpdatedTransactions(transactions: any): Observable<any> {
    return forkJoin(transactions.map((transaction: any) => {
      return forkJoin(transaction.transaction.inputs
        .map(input => this.historyService.getTransactionDetails(input))).pipe(
        mergeMap((inputDetails: any[]) => {
          transaction.owner_addressses = inputDetails.map(d => d.owner_address);
          return of(transaction);
        })
      );
    }));
  }

  private closeDataSubscription() {
    if (this.dataSubscription && !this.dataSubscription.closed) {
      this.dataSubscription.unsubscribe();
    }
  }
}
