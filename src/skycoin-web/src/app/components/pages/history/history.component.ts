import { ChangeDetectionStrategy, ChangeDetectorRef, Component, OnDestroy, OnInit } from '@angular/core';
import { MatDialogConfig } from '@angular/material/dialog';
import { Subscription, delay, first } from 'rxjs';
import { UntypedFormGroup, UntypedFormBuilder } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';

import { PriceService } from '../../../services/price.service';
import { HistoryService } from '../../../services/wallet/history.service';
import { TransactionDetailComponent } from './transaction-detail/transaction-detail.component';
import { BaseCoin } from '../../../coins/basecoin';
import { CoinService } from '../../../services/coin.service';
import { openQrModal } from '../../../utils';
import { CustomMatDialogService } from '../../../services/custom-mat-dialog.service';
import { NormalTransaction } from '../../../app.datatypes';
import { WalletService } from '../../../services/wallet/wallet.service';

export class Wallet {
  label!: string;
  coins!: string;
  hours!: string;
  addresses!: Address[];
  allAddressesSelected!: boolean;
}

export class Address {
  address!: string;
  coins!: string;
  hours!: string;
  showingWholeWallet!: boolean;
}

@Component({
    selector: 'app-history',
    templateUrl: './history.component.html',
    styleUrls: ['./history.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: false
})
export class HistoryComponent implements OnInit, OnDestroy {
  currentCoin!: BaseCoin;
  showError = false;

  allTransactions!: NormalTransaction[] | null;
  transactions!: NormalTransaction[] | null;
  price!: number;
  wallets!: Wallet[];
  form: UntypedFormGroup;

  private requestedAddress!: string;
  private walletsLoaded = false;
  private transactionsLoaded = false;
  private subscriptionsGroup: Subscription[] = [];
  private transactionsSubscription!: Subscription;
  private filterSubscription!: Subscription;
  private walletsSubscription!: Subscription;
  private routeSubscription: Subscription;

  constructor(
    private historyService: HistoryService,
    private priceService: PriceService,
    private dialog: CustomMatDialogService,
    private coinService: CoinService,
    private formBuilder: UntypedFormBuilder,
    private walletService: WalletService,
    route: ActivatedRoute,
    private changeDetectorRef: ChangeDetectorRef,
  ) {
    this.form = this.formBuilder.group({
      filter: [[]],
    });

    this.routeSubscription = route.queryParams.subscribe(params => {
      this.requestedAddress = params['addr'];
      this.showRequestedAddress();
      this.changeDetectorRef.markForCheck();
    });
  }

  ngOnInit() {
    this.subscriptionsGroup.push(this.coinService.currentCoin.subscribe((coin: BaseCoin) => {
      this.allTransactions = null;
      this.transactions = null;
      this.currentCoin = coin;
      this.showError = false;
      this.transactionsLoaded = false;
      this.walletsLoaded = false;
      this.form.get('filter')!.setValue([]);

      this.loadWallets();

      this.closeTransactionsSubscription();
      this.transactionsSubscription = this.historyService.transactions().subscribe(transactions => {
          this.allTransactions = transactions;
          this.transactions = transactions;

          this.transactionsLoaded = true;
          this.showRequestedAddress();
          this.changeDetectorRef.markForCheck();
        },
        () => {
          this.showError = true;
          this.changeDetectorRef.markForCheck();
        }
      );
      this.changeDetectorRef.markForCheck();
    }));

    this.filterSubscription = this.form.get('filter')!.valueChanges.subscribe(() => {
      const selectedfilters: (Wallet|Address)[] = this.form.get('filter')!.value;
      this.wallets.forEach(wallet => {
        wallet.allAddressesSelected = false;
        wallet.addresses.forEach(address => address.showingWholeWallet = false);
      });

      if (selectedfilters.length === 0) {
        this.transactions = this.allTransactions;
      } else {
        const selectedAddresses: Map<string, boolean> = new Map<string, boolean>();
        selectedfilters.forEach(filter => {
          if ((filter as Wallet).addresses) {
            (filter as Wallet).addresses.forEach(address => {
              selectedAddresses.set(address.address, true);
              address.showingWholeWallet = true;
            });
            (filter as Wallet).allAddressesSelected = true;
          } else {
            selectedAddresses.set((filter as Address).address, true);
          }
        });

        this.transactions = this.allTransactions!.filter(tx =>
          tx.inputs.some(input => selectedAddresses.has(input.owner)) || tx.outputs.some(output => selectedAddresses.has(output.dst)),
        );
      }
      this.changeDetectorRef.markForCheck();
    });

    this.subscriptionsGroup.push(this.priceService.price.subscribe(price => {
      this.price = price;
      this.changeDetectorRef.markForCheck();
    }));
  }

  ngOnDestroy() {
    this.subscriptionsGroup.forEach(sub => sub.unsubscribe());
    this.closeTransactionsSubscription();
    this.filterSubscription.unsubscribe();
    this.walletsSubscription.unsubscribe();
    this.routeSubscription.unsubscribe();
  }

  showTransaction(transaction: any) {
    const config = new MatDialogConfig();
    config.width = '800px';
    config.data = transaction;
    this.dialog.open(TransactionDetailComponent, config);
  }

  showQr(event: any, address: any) {
    event.stopPropagation();
    openQrModal(this.dialog, address);
  }

  removeFilters() {
    this.form.get('filter')!.setValue([]);
  }

  private loadWallets() {
    if (this.walletsSubscription) {
      this.walletsSubscription.unsubscribe();
    }

    this.walletsSubscription = this.walletService.currentWallets.pipe(delay(100), first()).subscribe(wallets => {
      this.wallets = [];
      let incompleteData = false;

      // A local copy of the data is created to avoid problems after updating the
      // wallet addresses while updating the balance.
      wallets.forEach(wallet => {
        if (!wallet.balance || !wallet.hours || incompleteData) {
          incompleteData = true;

          return;
        }

        this.wallets.push({
          label: wallet.label,
          coins: wallet.balance.decimalPlaces(6).toString(),
          hours: wallet.hours.decimalPlaces(0).toString(),
          addresses: [],
          allAddressesSelected: false,
        });

        wallet.addresses.forEach(address => {
          if (!address.balance || !address.hours || incompleteData) {
            incompleteData = true;

            return;
          }

          this.wallets[this.wallets.length - 1].addresses.push({
            address: address.address,
            coins: address.balance.decimalPlaces(6).toString(),
            hours: address.hours.decimalPlaces(0).toString(),
            showingWholeWallet: false,
          });
        });
      });

      if (incompleteData) {
        this.wallets = [];
        this.loadWallets();
      } else {
        this.walletsSubscription.unsubscribe();
        if (!this.walletsLoaded) {
          this.walletsLoaded = true;
          this.showRequestedAddress();
        }
      }

      this.changeDetectorRef.markForCheck();
    });
  }

  private closeTransactionsSubscription() {
    if (this.transactionsSubscription && !this.transactionsSubscription.closed) {
      this.transactionsSubscription.unsubscribe();
    }
  }

  private showRequestedAddress() {
    if (!this.transactionsLoaded || !this.wallets || this.wallets.length === 0) {
      return;
    }

    if (this.requestedAddress) {
      let addressFound: Address | undefined = undefined;
      this.wallets.forEach(wallet => {
        const found = wallet.addresses.find(address => address.address === this.requestedAddress);
        if (found) {
          addressFound = found;
        }
      });

      if (addressFound) {
        this.form.get('filter')!.setValue([addressFound]);
      }
    } else {
      this.form.get('filter')!.setValue([]);
    }
  }
}
