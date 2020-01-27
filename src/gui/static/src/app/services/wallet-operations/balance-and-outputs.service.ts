import { forkJoin as observableForkJoin, of, timer, Observable, ReplaySubject, Subscription } from 'rxjs';
import { first, mergeMap, map, switchMap, tap } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { ApiService } from '../api.service';
import { BigNumber } from 'bignumber.js';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletWithBalance, walletWithBalanceFromBase, Output, WalletBase, walletWithOutputsFromBase, WalletWithOutputs } from './wallet-objects';

@Injectable()
export class BalanceAndOutputsService {
  private walletsWithBalanceList: WalletWithBalance[];
  private walletsWithBalanceSubject: ReplaySubject<WalletWithBalance[]> = new ReplaySubject<WalletWithBalance[]>(1);
  private hasPendingTransactionsSubject: ReplaySubject<boolean> = new ReplaySubject<boolean>(1);

  private dataRefreshSubscription: Subscription;
  private gettingBalanceSubscription: Subscription;

  private forceCompleteBalanceArrayUpdate = false;

  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
    private ngZone: NgZone,
  ) {
    this.startDataRefreshSubscription();
  }

  get balance(): Observable<WalletWithBalance[]> {
    return this.walletsWithBalanceSubject.asObservable();
  }

  get hasPendingTransactions(): Observable<boolean> {
    return this.hasPendingTransactionsSubject.asObservable();
  }

  outputsWithWallets(): Observable<WalletWithOutputs[]> {
    return this.balance.pipe(switchMap(wallets => {
      const addresses = wallets.map(wallet => wallet.addresses.map(address => address.address).join(',')).join(',');

      return this.getOutputs(addresses);
    }), map(outputs => {
      const walletsList: WalletWithOutputs[] = [];
      this.walletsWithBalanceList.forEach(result => walletsList.push(walletWithOutputsFromBase(result)));

      return walletsList.map(wallet => {
        wallet.addresses = wallet.addresses.map(address => {
          address.outputs = outputs.filter(output => output.address === address.address);

          return address;
        });

        return wallet;
      });
    }));
  }

  getOutputs(addresses: string): Observable<Output[]> {
    if (!addresses) {
      return of([]);
    } else {
      return this.apiService.post('outputs', { addrs: addresses }).pipe(map((response) => {
        const outputs: Output[] = [];
        response.head_outputs.forEach(output => {
          const processedOutput = new Output();
          processedOutput.address = output.address;
          processedOutput.coins = new BigNumber(output.coins),
          processedOutput.hash = output.hash,
          processedOutput.calculated_hours = new BigNumber(output.calculated_hours),

          outputs.push(processedOutput);
        });

        return outputs;
      }));
    }
  }

  getWalletUnspentOutputs(wallet: WalletBase): Observable<Output[]> {
    const addresses = wallet.addresses.map(a => a.address).join(',');

    return this.getOutputs(addresses);
  }

  refreshBalance() {
    this.startDataRefreshSubscription();
  }

  private startDataRefreshSubscription() {
    if (this.dataRefreshSubscription) {
      this.dataRefreshSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.dataRefreshSubscription = this.walletsAndAddressesService.allWallets.pipe(
        tap(() => this.forceCompleteBalanceArrayUpdate = true),
        mergeMap(wallets => timer(0, 10000).pipe(map(() => wallets))),
      ).subscribe(wallets => {
        this.ngZone.run(() => this.refreshBalances(wallets));
      });
    });
  }

  private refreshBalances(wallets: WalletBase[]) {
    if (this.gettingBalanceSubscription) {
      this.gettingBalanceSubscription.unsubscribe();
    }

    const temporalWallets: WalletWithBalance[] = [];
    wallets.forEach(wallet => {
      temporalWallets.push(walletWithBalanceFromBase(wallet));
    });

    this.gettingBalanceSubscription = observableForkJoin(temporalWallets.map(wallet => this.retrieveWalletBalance(wallet))).subscribe(walletHasPendingTx => {
      this.hasPendingTransactionsSubject.next(walletHasPendingTx.some(value => value));

      if (!this.walletsWithBalanceList || this.forceCompleteBalanceArrayUpdate || this.walletsWithBalanceList.length !== temporalWallets.length) {
        this.forceCompleteBalanceArrayUpdate = false;
        this.walletsWithBalanceList = temporalWallets;
        this.informDataUpdated();
      } else {
        let changeDetected = false;
        this.walletsWithBalanceList.forEach((currentWallet, i) => {
          if (currentWallet.id !== temporalWallets[i].id) {
            changeDetected = true;
          }
        });

        if (changeDetected) {
          this.walletsWithBalanceList = temporalWallets;
          this.informDataUpdated();
        } else {
          this.walletsWithBalanceList.forEach((currentWallet, i) => {
            if (!currentWallet.coins.isEqualTo(temporalWallets[i].coins) || !currentWallet.hours.isEqualTo(temporalWallets[i].hours)) {
              currentWallet.coins = temporalWallets[i].coins;
              currentWallet.hours = temporalWallets[i].hours;
              changeDetected = true;
            }

            if (currentWallet.addresses.length !== temporalWallets[i].addresses.length) {
              currentWallet.addresses = temporalWallets[i].addresses;
              changeDetected = true;
            } else {
              currentWallet.addresses.forEach((currentAddress, j) => {
                if (!currentAddress.coins.isEqualTo(temporalWallets[i].addresses[j].coins) || !currentAddress.hours.isEqualTo(temporalWallets[i].addresses[j].hours)) {
                  currentAddress.coins = temporalWallets[i].addresses[j].coins;
                  currentAddress.hours = temporalWallets[i].addresses[j].hours;
                  changeDetected = true;
                }
              });
            }
          });

          if (changeDetected) {
            this.informDataUpdated();
          }
        }
      }
    });
  }

  private retrieveWalletBalance(wallet: WalletWithBalance): Observable<boolean> {
    let query: Observable<any>;
    if (!wallet.isHardware) {
      query = this.apiService.get('wallet/balance', { id: wallet.id });
    } else {
      const formattedAddresses = wallet.addresses.map(a => a.address).join(',');
      query = this.apiService.post('balance', { addrs: formattedAddresses });
    }

    return query.pipe(map(balance => {
      wallet.coins = new BigNumber(balance.confirmed.coins).dividedBy(1000000);
      wallet.hours = new BigNumber(balance.confirmed.hours);

      wallet.addresses.forEach(address => {
        if (balance.addresses[address.address]) {
          address.coins = new BigNumber(balance.addresses[address.address].confirmed.coins).dividedBy(1000000);
          address.hours = new BigNumber(balance.addresses[address.address].confirmed.hours);
        } else {
          address.coins = new BigNumber(0);
          address.hours = new BigNumber(0);
        }
      });

      return !(new BigNumber(balance.predicted.coins).dividedBy(1000000)).isEqualTo(wallet.coins) ||
        !(new BigNumber(balance.predicted.hours)).isEqualTo(wallet.hours);
    }));
  }

  private informDataUpdated() {
    this.walletsWithBalanceSubject.next(this.walletsWithBalanceList);
  }
}
