import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';
import 'rxjs/add/observable/forkJoin';
import 'rxjs/add/observable/of';
import 'rxjs/add/operator/do';
import 'rxjs/add/operator/filter';
import 'rxjs/add/operator/first';
import 'rxjs/add/operator/mergeMap';
import { Address, Wallet } from '../app.datatypes';
import { ReplaySubject } from 'rxjs/ReplaySubject';

@Injectable()
export class WalletService {
  addresses: Address[];
  wallets: Subject<Wallet[]> = new ReplaySubject<Wallet[]>();
  pendingTxs: Subject<any[]> = new ReplaySubject<any[]>();

  constructor(
    private apiService: ApiService
  ) {
    this.loadData();

    IntervalObservable
      .create(30000)
      .subscribe(() => this.refreshBalances());

    IntervalObservable.create(10000)
      .subscribe(() => this.refreshPendingTransactions());
  }

  addressesAsString(): Observable<string> {
    return this.allAddresses().map(addrs => addrs.map(addr => addr.address)).map(addrs => addrs.join(','));
  }

  addAddress(wallet: Wallet, password?: string) {
    return this.apiService.postWalletNewAddress(wallet, password)
      .do(address => {
        wallet.addresses.push(address);
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
    seed = seed.replace(/\r?\n|\r/g, ' ').replace(/ +/g, ' ').trim();

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

  find(filename: string): Observable<Wallet> {
    return this.all().map(wallets => wallets.find(wallet => wallet.filename === filename));
  }

  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').map(response => response.address);
  }

  generateSeed(): Observable<string> {
    return this.apiService.getWalletNewSeed();
  }

  outputs(): Observable<any> {
    return this.addressesAsString()
      .filter(addresses => !!addresses)
      .flatMap(addresses => this.apiService.get('outputs', {addrs: addresses}));
  }

  allPendingTransactions(): Observable<any> {
    return this.apiService.get('pendingTxs');
  }

  pendingTransactions(): Observable<any> {
    return this.pendingTxs.asObservable();
  }

  refreshPendingTransactions() {
    this.wallets.first().subscribe(wallets => {
      Observable.forkJoin(wallets.map(wallet => this.apiService.get('wallet/transactions', { id: wallet.filename })))
        .subscribe(pending => {
          this.pendingTxs.next([].concat.apply(
            [],
            pending
              .filter(response => response.transactions.length > 0)
              .map(response => response.transactions)
          ).reduce((txs, tx) => {
            if (!txs.find(t => t.transaction.txid === tx.transaction.txid)) {
              txs.push(tx);
            }
            return txs;
          }, []));
        });
    });
  }

  refreshBalances() {
    this.wallets.first().subscribe(wallets => {
      Observable.forkJoin(wallets.map(wallet => this.retrieveWalletBalance(wallet).map(response => {
        wallet.addresses = response;
        wallet.coins = response.map(address => address.coins >= 0 ? address.coins : 0).reduce((a , b) => a + b, 0);
        wallet.hours = response.map(address => address.hours >= 0 ? address.hours : 0).reduce((a , b) => a + b, 0);
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

  sendSkycoin(wallet: Wallet, address: string, amount: number, password: string|null) {
    return this.apiService.post('wallet/spend', {id: wallet.filename, dst: address, coins: amount, password});
  }

  sum(): Observable<number> {
    return this.all().map(wallets => wallets.map(wallet => wallet.coins >= 0 ? wallet.coins : 0).reduce((a , b) => a + b, 0));
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

  transactions(): Observable<any[]> {
    return this.allAddresses().filter(addresses => !!addresses.length).first().flatMap(addresses => {
      this.addresses = addresses;
      return Observable.forkJoin(addresses.map(address => this.apiService.getExplorerAddress(address)));
    }).map(transactions => [].concat.apply([], transactions).sort((a, b) =>  b.timestamp - a.timestamp))
      .map(transactions => transactions.reduce((array, item) => {
        if (!array.find(trans => trans.txid === item.txid)) {
          array.push(item);
        }
        return array;
      }, []))
      .map(transactions => transactions.map(transaction => {
        const outgoing = !!this.addresses.find(address => transaction.inputs[0].owner === address.address);
        transaction.outputs.forEach(output => {
          if (outgoing && !this.addresses.find(address => output.dst === address.address)) {
            transaction.addresses.push(output.dst);
            transaction.balance = transaction.balance - parseFloat(output.coins);
          }
          if (!outgoing && this.addresses.find(address => output.dst === address.address)) {
            transaction.addresses.push(output.dst);
            transaction.balance = transaction.balance + parseFloat(output.coins);
          }
          return transaction;
        });

        return transaction;
      }));
  }

  private loadData(): void {
    this.apiService.getWallets().first().subscribe(wallets => {
      this.wallets.next(wallets);
      this.refreshBalances();
      this.refreshPendingTransactions();
    });
  }

  private retrieveAddressBalance(address: any|any[]) {
    const addresses = Array.isArray(address) ? address.map(addr => addr.address).join(',') : address.address;
    return this.apiService.get('balance', {addrs: addresses});
  }

  private retrieveInputAddress(input: string) {
    return this.apiService.get('uxout', {uxid: input});
  }

  private retrieveWalletBalance(wallet: Wallet): Observable<any> {
    return Observable.forkJoin(wallet.addresses.map(address => this.retrieveAddressBalance(address).map(balance => {
      address.coins = balance.confirmed.coins / 1000000;
      address.hours = balance.confirmed.hours;
      return address;
    })));
  }

  private updateWallet(wallet: Wallet) {
    this.wallets.first().subscribe(wallets => {
      const index = wallets.findIndex(w => w.filename === wallet.filename);
      wallets[index] = wallet;
      this.wallets.next(wallets);
    });
  }
}
