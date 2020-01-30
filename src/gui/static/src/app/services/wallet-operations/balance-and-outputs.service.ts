import { forkJoin as observableForkJoin, of, timer, Observable, ReplaySubject, Subscription, BehaviorSubject } from 'rxjs';
import { mergeMap, map, switchMap, tap, delay } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { ApiService } from '../api.service';
import { BigNumber } from 'bignumber.js';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletWithBalance, walletWithBalanceFromBase, Output, WalletBase, walletWithOutputsFromBase, WalletWithOutputs } from './wallet-objects';

@Injectable()
export class BalanceAndOutputsService {
  private walletsWithBalanceList: WalletWithBalance[];
  private walletsWithBalanceSubject: ReplaySubject<WalletWithBalance[]> = new ReplaySubject<WalletWithBalance[]>(1);
  private hasPendingTransactionsSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  private firstFullUpdateMadeSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private dataRefreshSubscription: Subscription;
  private gettingBalanceSubscription: Subscription;

  private savedBalanceData = new Map<string, any>();
  private temporalSavedBalanceData = new Map<string, any>();
  private savedWalletsList: WalletBase[];

  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
    private ngZone: NgZone,
  ) {
    this.walletsAndAddressesService.allWallets.subscribe(wallets => {
      this.savedWalletsList = wallets;
      this.startDataRefreshSubscription(0);
    });
  }

  get walletsWithBalance(): Observable<WalletWithBalance[]> {
    return this.walletsWithBalanceSubject.asObservable();
  }

  get hasPendingTransactions(): Observable<boolean> {
    return this.hasPendingTransactionsSubject.asObservable();
  }

  get firstFullUpdateMade(): Observable<boolean> {
    return this.firstFullUpdateMadeSubject.asObservable();
  }

  outputsWithWallets(): Observable<WalletWithOutputs[]> {
    return this.walletsWithBalance.pipe(switchMap(wallets => {
      const addresses = wallets.map(wallet => wallet.addresses.map(address => address.address).join(',')).join(',');

      return this.getOutputs(addresses);
    }), map(outputs => {
      const walletsList: WalletWithOutputs[] = [];
      this.walletsWithBalanceList.forEach(wallet => walletsList.push(walletWithOutputsFromBase(wallet)));

      walletsList.forEach(wallet => {
        wallet.addresses.forEach(address => {
          address.outputs = outputs.filter(output => output.address === address.address);
        });
      });

      return walletsList;
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
    this.startDataRefreshSubscription(0);
  }

  private startDataRefreshSubscription(delayMs: number) {
    if (this.dataRefreshSubscription) {
      this.dataRefreshSubscription.unsubscribe();
    }

    if (this.savedWalletsList) {
      this.ngZone.runOutsideAngular(() => {
        this.dataRefreshSubscription = of(0).pipe(delay(delayMs), mergeMap(() => {
          return this.refreshBalances(this.savedWalletsList, true);
        }), mergeMap(() => {
          return this.refreshBalances(this.savedWalletsList, false);
        })).subscribe(() => this.startDataRefreshSubscription(10000), () => this.startDataRefreshSubscription(0));
      });
    }
  }

  private refreshBalances(wallets: WalletBase[], forceQuickCompleteArrayUpdate: boolean): Observable<any> {
    if (this.gettingBalanceSubscription) {
      this.gettingBalanceSubscription.unsubscribe();
    }

    const temporalWallets: WalletWithBalance[] = [];
    wallets.forEach(wallet => {
      temporalWallets.push(walletWithBalanceFromBase(wallet));
    });

    if (!forceQuickCompleteArrayUpdate) {
      this.temporalSavedBalanceData = new Map<string, any>();
    }

    return observableForkJoin(temporalWallets.map(wallet => this.retrieveWalletBalance(wallet, forceQuickCompleteArrayUpdate))).pipe(tap(walletHasPendingTx => {
      this.hasPendingTransactionsSubject.next(walletHasPendingTx.some(value => value));

      if (!forceQuickCompleteArrayUpdate) {
        this.savedBalanceData = this.temporalSavedBalanceData;
        if (!this.firstFullUpdateMadeSubject.value) {
          this.firstFullUpdateMadeSubject.next(true);
        }
      }

      if (!this.walletsWithBalanceList || forceQuickCompleteArrayUpdate || this.walletsWithBalanceList.length !== temporalWallets.length) {
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
    }));
  }

  private retrieveWalletBalance(wallet: WalletWithBalance, useSavedBalanceData: boolean): Observable<boolean> {
    let query: Observable<any>;

    if (!useSavedBalanceData) {
      if (!wallet.isHardware) {
        query = this.apiService.get('wallet/balance', { id: wallet.id });
      } else {
        const formattedAddresses = wallet.addresses.map(a => a.address).join(',');
        query = this.apiService.post('balance', { addrs: formattedAddresses });
      }
    } else {
      if (this.savedBalanceData.has(wallet.id)) {
        query = of(this.savedBalanceData.get(wallet.id));
      } else {
        query = of({ addresses: [] });
      }
    }

    return query.pipe(map(balance => {
      this.temporalSavedBalanceData.set(wallet.id, balance);

      if (balance.confirmed) {
        wallet.coins = new BigNumber(balance.confirmed.coins).dividedBy(1000000);
        wallet.hours = new BigNumber(balance.confirmed.hours);
      } else {
        wallet.coins = new BigNumber(0);
        wallet.hours = new BigNumber(0);
      }

      wallet.addresses.forEach(address => {
        if (balance.addresses[address.address]) {
          address.coins = new BigNumber(balance.addresses[address.address].confirmed.coins).dividedBy(1000000);
          address.hours = new BigNumber(balance.addresses[address.address].confirmed.hours);
        } else {
          address.coins = new BigNumber(0);
          address.hours = new BigNumber(0);
        }
      });

      if (!useSavedBalanceData) {
        return !(new BigNumber(balance.predicted.coins).dividedBy(1000000)).isEqualTo(wallet.coins) ||
          !(new BigNumber(balance.predicted.hours)).isEqualTo(wallet.hours);
      } else {
        return this.hasPendingTransactionsSubject.value;
      }
    }));
  }

  private informDataUpdated() {
    this.ngZone.run(() => {
      this.walletsWithBalanceSubject.next(this.walletsWithBalanceList);
    });
  }
}
