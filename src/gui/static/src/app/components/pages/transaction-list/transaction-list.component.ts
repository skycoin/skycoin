import { Component, OnDestroy, OnInit } from '@angular/core';
import { WalletService } from '../../../services/wallet.service';
import { PriceService } from '../../../services/price.service';
import { ISubscription } from 'rxjs/Subscription';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { TransactionDetailComponent } from './transaction-detail/transaction-detail.component';
import { NormalTransaction } from '../../../app.datatypes';
import { QrCodeComponent } from '../../layout/qr-code/qr-code.component';
import { FormGroup, FormBuilder } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';

export class Wallet {
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

  private price: number;
  private requestedAddress: string;
  private walletsLoaded = false;
  private transactionsLoaded = false;
  private priceSubscription: ISubscription;
  private filterSubscription: ISubscription;
  private walletsSubscription: ISubscription;
  private routeSubscription: ISubscription;

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
      this.requestedAddress = params['addr'];
      this.showRequestedAddress();
    });

    this.walletsSubscription = walletService.all().delay(1).subscribe(wallets => {
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
        this.wallets = [];
      } else {
        this.walletsSubscription.unsubscribe();
        if (!this.walletsLoaded) {
          this.walletsLoaded = true;
          this.showRequestedAddress();
        }
      }

    });
  }

  ngOnInit() {
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);

    this.walletService.transactions().first().subscribe(transactions => {
      this.allTransactions = transactions;
      this.transactions = transactions;

      this.transactionsLoaded = true;
      this.showRequestedAddress();
    });

    this.filterSubscription = this.form.get('filter').valueChanges.subscribe(() => {
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
    });
  }

  ngOnDestroy() {
    this.priceSubscription.unsubscribe();
    this.filterSubscription.unsubscribe();
    this.walletsSubscription.unsubscribe();
    this.routeSubscription.unsubscribe();
  }

  showTransaction(transaction: NormalTransaction) {
    const config = new MatDialogConfig();
    config.width = '800px';
    config.data = transaction;
    this.dialog.open(TransactionDetailComponent, config);
  }

  showQrCode(event: any, address: string) {
    event.stopPropagation();

    const config = new MatDialogConfig();
    config.data = { address };
    this.dialog.open(QrCodeComponent, config);
  }

  removeFilters() {
    this.form.get('filter').setValue([]);
  }

  private showRequestedAddress() {
    if (!this.transactionsLoaded || !this.wallets || this.wallets.length === 0) {
      return;
    }

    if (this.requestedAddress) {
      let addressFound: Address;
      this.wallets.forEach(wallet => {
        const found = wallet.addresses.find(address => address.address === this.requestedAddress);
        if (found) {
          addressFound = found;
        }
      });

      if (addressFound) {
        this.form.get('filter').setValue([addressFound]);
      }
    } else {
      this.form.get('filter').setValue([]);
    }
  }
}
