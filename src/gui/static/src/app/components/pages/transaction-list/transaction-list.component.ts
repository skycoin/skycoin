import { delay, mergeMap, first } from 'rxjs/operators';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { PriceService } from '../../../services/price.service';
import { SubscriptionLike, of } from 'rxjs';
import { MatDialog } from '@angular/material/dialog';
import { TransactionDetailComponent } from './transaction-detail/transaction-detail.component';
import { NormalTransaction } from '../../../app.datatypes';
import { FormGroup, FormBuilder } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';

export class Wallet {
  id: string;
  label: string;
  coins: string;
  hours: string;
  addresses: Address[];
  allAddressesSelected: boolean;
}

export class Address {
  address: string;
  coins: string;
  hours: string;
  showingWholeWallet: boolean;
}

@Component({
  selector: 'app-transaction-list',
  templateUrl: './transaction-list.component.html',
  styleUrls: ['./transaction-list.component.scss'],
})
export class TransactionListComponent implements OnInit, OnDestroy {
  allTransactions: NormalTransaction[];
  transactions: NormalTransaction[];
  wallets: Wallet[];
  form: FormGroup;

  readonly maxInitialElements = 40;
  viewAll = false;
  viewingTruncatedList = false;
  totalElements: number;

  private price: number;
  private requestedFilters: string[];
  private transactionsLoaded = false;
  private priceSubscription: SubscriptionLike;
  private filterSubscription: SubscriptionLike;
  private walletsSubscription: SubscriptionLike;
  private routeSubscription: SubscriptionLike;

  constructor(
    private dialog: MatDialog,
    private priceService: PriceService,
    private walletService: WalletService,
    private formBuilder: FormBuilder,
    route: ActivatedRoute,
  ) {

    this.form = this.formBuilder.group({
      filter: [[]],
    });

    this.routeSubscription = route.queryParams.subscribe(params => {
      let Addresses = params['addr'] ? (params['addr'] as string).split(',') : [];
      let Wallets = params['wal'] ? (params['wal'] as string).split(',') : [];
      Addresses = Addresses.map(element => 'a-' + element);
      Wallets = Wallets.map(element => 'w-' + element);
      this.viewAll = false;

      this.requestedFilters = Addresses.concat(Wallets);
      this.showRequestedFilters();
    });

    this.walletsSubscription = walletService.all().pipe(delay(1), mergeMap(wallets => {
      if (!this.wallets) {
        this.wallets = [];
        let incompleteData = false;

        // A local copy of the data is created to avoid problems after updating the
        // wallet addresses while updating the balance.
        wallets.forEach(wallet => {
          if (!wallet.coins || !wallet.hours || incompleteData) {
            incompleteData = true;

            return;
          }

          this.wallets.push({
            id: wallet.filename,
            label: wallet.label,
            coins: wallet.coins.decimalPlaces(6).toString(),
            hours: wallet.hours.decimalPlaces(0).toString(),
            addresses: [],
            allAddressesSelected: false,
          });

          wallet.addresses.forEach(address => {
            if (!address.coins || !address.hours || incompleteData) {
              incompleteData = true;

              return;
            }

            this.wallets[this.wallets.length - 1].addresses.push({
              address: address.address,
              coins: address.coins.decimalPlaces(6).toString(),
              hours: address.hours.decimalPlaces(0).toString(),
              showingWholeWallet: false,
            });
          });
        });

        if (incompleteData) {
          this.wallets = null;

          return of(null);
        } else {
          return this.walletService.transactions().pipe(first());
        }
      } else {
        return this.walletService.transactions().pipe(first());
      }
    })).subscribe(transactions => {
      if (transactions) {
        this.allTransactions = transactions;

        this.transactionsLoaded = true;
        this.showRequestedFilters();

        this.filterTransactions();
      }
    });
  }

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
    this.filterSubscription = this.form.get('filter').valueChanges.subscribe(() => {
      this.viewAll = false;
      this.filterTransactions();
    });
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
    this.filterSubscription.unsubscribe();
    this.walletsSubscription.unsubscribe();
    this.routeSubscription.unsubscribe();
  }

  showAll() {
    if (!this.viewAll) {
      this.viewAll = true;
      this.filterTransactions();
    }
  }

  showTransaction(transaction: NormalTransaction) {
    TransactionDetailComponent.openDialog(this.dialog, transaction);
  }

  removeFilters() {
    this.form.get('filter').setValue([]);
  }

  private filterTransactions() {
    const selectedfilters: (Wallet|Address)[] = this.form.get('filter').value;
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
          (filter as Wallet).addresses.forEach(address => selectedAddresses.set(address.address, true));
          (filter as Wallet).allAddressesSelected = true;
          (filter as Wallet).addresses.forEach(address => address.showingWholeWallet = true);
        } else {
          selectedAddresses.set((filter as Address).address, true);
        }
      });

      this.transactions = this.allTransactions.filter(tx =>
        tx.inputs.some(input => selectedAddresses.has(input.owner)) || tx.outputs.some(output => selectedAddresses.has(output.dst)),
      );
    }

    this.totalElements = this.transactions.length;

    if (!this.viewAll && this.totalElements > this.maxInitialElements) {
      this.transactions = this.transactions.slice(0, this.maxInitialElements);
      this.viewingTruncatedList = true;
    } else {
      this.viewingTruncatedList = false;
    }
  }

  private showRequestedFilters() {
    if (!this.transactionsLoaded || !this.wallets || this.wallets.length === 0 || this.requestedFilters === null || this.requestedFilters === undefined) {
      return;
    }

    if (this.requestedFilters.length > 0) {
      const filters: (Wallet|Address)[] = [];

      this.requestedFilters.forEach(filter => {
        const filterContent = filter.substr(2, filter.length - 2);
        this.wallets.forEach(wallet => {
          if (filter.startsWith('w-')) {
            if (filterContent === wallet.id) {
              filters.push(wallet);
            }
          } else if (filter.startsWith('a-')) {
            wallet.addresses.forEach(address => {
              if (filterContent === address.address) {
                filters.push(address);
              }
            });
          }
        });
      });

      this.form.get('filter').setValue(filters);
    } else {
      this.form.get('filter').setValue([]);
    }

    this.requestedFilters = null;
  }
}
