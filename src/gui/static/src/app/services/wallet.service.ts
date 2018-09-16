import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/forkJoin';
import 'rxjs/add/observable/of';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/first';
import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/observable/timer';
import 'rxjs/add/observable/zip';
import { Address, NormalTransaction, PreviewTransaction, Wallet } from '../app.datatypes';
import { ReplaySubject } from 'rxjs/ReplaySubject';
import { Subscription } from 'rxjs/Subscription';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class WalletService {
  addresses: Address[];
  wallets: Subject<Wallet[]> = new ReplaySubject<Wallet[]>();
  pendingTxs: Subject<any[]> = new ReplaySubject<any[]>();
  dataRefreshSubscription: Subscription;

  initialLoadFailed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor(
    private apiService: ApiService,
    private ngZone: NgZone,
  ) {
    this.loadData();
    this.startDataRefreshSubscription();
  }

  addressesAsString(): Observable<string> {
    return this.allAddresses().map(addrs => addrs.map(addr => addr.address)).map(addrs => addrs.join(','));
  }

  addAddress(wallet: Wallet, num: number, password?: string) {
    return this.apiService.postWalletNewAddress(wallet, num, password)
      .do(addresses => {
        addresses.forEach(value => wallet.addresses.push(value));
        this.refreshBalances();
      });
  }

  all(): Observable<Wallet[]> {
    return this.wallets.asObservable();
  }

  allAddresses(): Observable<any[]> {
    return this.all().map(wallets => wallets.reduce((array, wallet) => array.concat(wallet.addresses), []));
  }

  create(label, seed, scan, password) {
    seed = seed.replace(/(\n|\r\n)$/, '');

    return this.apiService.postWalletCreate(label ? label : 'undefined', seed, scan ? scan : 100, password)
      .do(wallet => {
        console.log(wallet);
        this.wallets.first().subscribe(wallets => {
          wallets.push(wallet);
          this.wallets.next(wallets);
          this.refreshBalances();
        });
      });
  }

  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').map(response => response.address);
  }

  outputs(): Observable<any> {
    return this.addressesAsString()
      .first()
      .filter(addresses => !!addresses)
      .flatMap(addresses => this.apiService.get('outputs', {addrs: addresses}));
  }

  outputsWithWallets(): Observable<any> {
    return Observable.zip(this.all(), this.outputs(), (wallets, outputs) => {
      return wallets.map(wallet => {
        wallet.addresses = wallet.addresses.map(address => {
          address.outputs = outputs.head_outputs.filter(output => output.address === address.address);

          return address;
        });

        return wallet;
      });
    });
  }

  allPendingTransactions(): Observable<any> {
    return Observable.timer(0, 10000).flatMap(() => this.apiService.get('pendingTxs'));
  }

  pendingTransactions(): Observable<any> {
    return this.pendingTxs.asObservable();
  }

  refreshBalances() {
    this.wallets.first().subscribe(wallets => {
      Observable.forkJoin(wallets.map(wallet => this.retrieveWalletBalance(wallet).map(response => {
        wallet.coins = response.coins;
        wallet.hours = response.hours;
        wallet.addresses = wallet.addresses.map(address => {
          return response.addresses.find(addr => addr.address === address.address);
        });

        return wallet;
      })))
        .subscribe(newWallets => this.wallets.next(newWallets));
    });
  }

  renameWallet(wallet: Wallet, label: string): Observable<Wallet> {
    return this.apiService.post('wallet/update', { id: wallet.filename, label: label })
      .do(() => {
        wallet.label = label;
        this.updateWallet(wallet);
      });
  }

  toggleEncryption(wallet: Wallet, password: string): Observable<Wallet> {
    return this.apiService.postWalletToggleEncryption(wallet, password)
      .do(w => {
        wallet.encrypted = w.meta.encrypted;
        this.updateWallet(w);
      });
  }

  getWalletSeed(wallet: Wallet, password: string): Observable<string> {
    return this.apiService.getWalletSeed(wallet, password);
  }

  createTransaction(wallet: Wallet, addresses: string[]|null, destinations: any[], hoursSelection: any, changeAddress: string|null, password: string|null): Observable<PreviewTransaction> {
    return this.apiService.post(
      'wallet/transaction',
      {
        hours_selection: hoursSelection,
        wallet: {
          id: wallet.filename,
          password,
          addresses,
        },
        to: destinations,
        change_address: changeAddress,
      },
      {
        json: true,
      },
    ).map(response => {
      return {
        ...response.transaction,
        hoursBurned: response.transaction.fee,
        encoded: response.encoded_transaction,
      };
    });
  }

  injectTransaction(encodedTx: string) {
    return this.apiService.post('injectTransaction', { rawtx: encodedTx }, { json: true });
  }

  transaction(txid: string): Observable<any> {
    return this.apiService.get('transaction', {txid: txid}).flatMap(transaction => {
      if (transaction.txn.inputs && !transaction.txn.inputs.length) {
        return Observable.of(transaction);
      }

      return Observable.forkJoin(transaction.txn.inputs.map(input => this.retrieveInputAddress(input).map(response => {
        return response.owner_address;
      }))).map(inputs => {
        transaction.txn.inputs = inputs;

        return transaction;
      });
    });
  }

  transactions(): Observable<NormalTransaction[]> {
    return this.allAddresses().first().flatMap(addresses => {
      this.addresses = addresses;

      return Observable.forkJoin(addresses.map(address => this.apiService.getExplorerAddress(address)));
    }).map(transactions => {
      return []
        .concat.apply([], transactions)
        .reduce((array, item) => {
          if (!array.find(trans => trans.txid === item.txid)) {
            array.push(item);
          }

          return array;
        }, [])
        .sort((a, b) =>  b.timestamp - a.timestamp)
        .map(transaction => {
          const outgoing = this.addresses.some(address => {
            return transaction.inputs.some(input => input.owner === address.address);
          });

          const relevantOutputs = transaction.outputs.reduce((array, output) => {
            const isMyOutput = this.addresses.some(address => address.address === output.dst);

            if ((outgoing && !isMyOutput) || (!outgoing && isMyOutput)) {
              array.push(output);
            }

            return array;
          }, []);

          const calculatedOutputs = (outgoing && relevantOutputs.length === 0)
          || (!outgoing && relevantOutputs.length === transaction.outputs.length)
            ? transaction.outputs
            : relevantOutputs;

          transaction.addresses.push(
            ...calculatedOutputs
              .map(output => output.dst)
              .filter((dst, i, self) => self.indexOf(dst) === i),
          );

          transaction.balance += calculatedOutputs.reduce((a, b) => a + parseFloat(b.coins), 0) * (outgoing ? -1 : 1);
          transaction.hoursSent = calculatedOutputs.reduce((a, b) => a + b.hours, 0);

          const inputsHours = transaction.inputs.reduce((a, b) => a + b.calculated_hours, 0);
          const outputsHours = transaction.outputs.reduce((a, b) => a + b.hours, 0);
          transaction.hoursBurned = inputsHours - outputsHours;

          return transaction;
        });
    });
  }

  startDataRefreshSubscription() {
    if (this.dataRefreshSubscription) {
      this.dataRefreshSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.dataRefreshSubscription = Observable.timer(0, 10000)
        .subscribe(() => this.ngZone.run(() => {
          this.refreshBalances();
          this.refreshPendingTransactions();
        }));
    });
  }

  private loadData(): void {
    this.apiService.getWallets().first().subscribe(
      wallets => this.wallets.next(wallets),
      () => this.initialLoadFailed.next(true),
    );
  }

  private retrieveInputAddress(input: string) {
    return this.apiService.get('uxout', {uxid: input});
  }

  private retrieveWalletBalance(wallet: Wallet): Observable<any> {
    return this.apiService.get('wallet/balance', { id: wallet.filename }).map(balance => {
      return {
        coins: balance.confirmed.coins / 1000000,
        hours: balance.confirmed.hours,
        addresses: Object.keys(balance.addresses).map(address => ({
          address,
          coins: balance.addresses[address].confirmed.coins / 1000000,
          hours: balance.addresses[address].confirmed.hours,
        })),
      };
    });
  }

  private updateWallet(wallet: Wallet) {
    this.wallets.first().subscribe(wallets => {
      const index = wallets.findIndex(w => w.filename === wallet.filename);
      wallets[index] = wallet;
      this.wallets.next(wallets);
    });
  }

  private refreshPendingTransactions() {
    this.wallets.first().subscribe(wallets => {
      Observable.forkJoin(wallets.map(wallet => this.apiService.get('wallet/transactions', { id: wallet.filename })))
        .subscribe(pending => {
          this.pendingTxs.next([].concat.apply(
            [],
            pending
              .filter(response => response.transactions.length > 0)
              .map(response => response.transactions),
          ).reduce((txs, tx) => {
            if (!txs.find(t => t.transaction.txid === tx.transaction.txid)) {
              txs.push(tx);
            }

            return txs;
          }, []));
        });
    });
  }
}
