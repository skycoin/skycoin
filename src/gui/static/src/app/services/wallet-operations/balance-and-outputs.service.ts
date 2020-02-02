import { forkJoin as observableForkJoin, of, Observable, ReplaySubject, Subscription, BehaviorSubject } from 'rxjs';
import { mergeMap, map, switchMap, tap, delay } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api.service';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletWithBalance, walletWithBalanceFromBase, Output, WalletBase, walletWithOutputsFromBase, WalletWithOutputs } from './wallet-objects';

/**
 * Allows to get the balance of the wallets and is in chage of maintaining those balances updated.
 * It also allows to get the unspent outputs of the wallets and lists of addresses.
 */
@Injectable()
export class BalanceAndOutputsService {
  // The list of wallets with balance and the subject used for informing when the list has been modified.
  private walletsWithBalanceList: WalletWithBalance[];
  private walletsWithBalanceSubject: ReplaySubject<WalletWithBalance[]> = new ReplaySubject<WalletWithBalance[]>(1);

  private hasPendingTransactionsSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  private firstFullUpdateMadeSubject: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  private dataRefreshSubscription: Subscription;
  private gettingBalanceSubscription: Subscription;

  /**
   * After the service retrieves the balance of each wallet, the response returned by the node for each
   * wallet is saved here, accessible via the wallet id.
   */
  private savedBalanceData = new Map<string, any>();
  /**
   * Temporal map for updating savedBalanceData only after retrieving the data of all wallets, to avoid
   * problems when the balance update procedure is cancelled early.
   */
  private temporalSavedBalanceData = new Map<string, any>();
  /**
   * Saves the lastest, most updated, wallet list obtained from the wallets service.
   */
  private savedWalletsList: WalletBase[];

  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
    private ngZone: NgZone,
  ) {
    // React every time the wallet list is updated.
    this.walletsAndAddressesService.allWallets.subscribe(wallets => {
      this.savedWalletsList = wallets;
      this.startDataRefreshSubscription(0, true);
    });
  }

  /**
   * Gets the wallet list, with the balance of each wallet and address. It emits when the
   * wallet list is updated and when the balance changes. Please note that the list will
   * tell all the wallets have balance 0 util the service finishes connecting to the
   * backend node for the firts time. Also note that if any value of the returned wallets
   * is modified, the changes must be notified to the wallets service or the behavior will
   * be indeterminate.
   */
  get walletsWithBalance(): Observable<WalletWithBalance[]> {
    return this.walletsWithBalanceSubject.asObservable();
  }

  /**
   * Indicates if there are pending transactions affecting any of the wallets of the
   * wallet list.
   */
  get hasPendingTransactions(): Observable<boolean> {
    return this.hasPendingTransactionsSubject.asObservable();
  }

  /**
   * Indicates if the service already got the balances of the wallets from the node for
   * the first time. The wallets returned by walletsWithBalance will always show blance 0
   * until this property returns true.
   */
  get firstFullUpdateMade(): Observable<boolean> {
    return this.firstFullUpdateMadeSubject.asObservable();
  }

  /**
   * Gets the wallet list, with the unspent outputs of each address. It emits when the
   * wallet list is updated and when the balance changes.Please note that if any value
   * of the returned wallets is modified, the changes must be notified to the wallets
   * service or the behavior will be indeterminate.
   */
  get outputsWithWallets(): Observable<WalletWithOutputs[]> {
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

  /**
   * Gets the list of unspent outputs of a list of addresses.
   * @param addresses List of addresses, comma separated.
   * @returns Arrays with all the unspent outputs owned by any of the provide addresses.
   */
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
          processedOutput.hours = new BigNumber(output.calculated_hours),

          outputs.push(processedOutput);
        });

        return outputs;
      }));
    }
  }

  /**
   * Gets the list of unspent outputs owned by a wallet.
   * @param wallet Wallet to check.
   * @returns Arrays with all the unspent outputs owned by any of the addresses of the wallet.
   */
  getWalletUnspentOutputs(wallet: WalletBase): Observable<Output[]> {
    const addresses = wallet.addresses.map(a => a.address).join(',');

    return this.getOutputs(addresses);
  }

  /**
   * Asks the service to update the balance inmediatelly.
   */
  refreshBalance() {
    this.startDataRefreshSubscription(0, false);
  }

  /**
   * Makes the service start updating the balance periodically. If this function was called
   * before, the previous updating procedure is cancelled.
   * @param delayMs Delay before starting to update the balance.
   * @param updateWalletsFirst If true, after the delay the function will inmediatelly update
   * the wallet list with the data on savedWalletsList and using the last balance data obtained
   * from the node (or will set all the wallets to 0, if no data exists) and only after that will
   * try to get the balance data from the node and update the wallet list again. This allows to
   * inmediatelly reflect changes made to the wallet list, without having to wait for the node
   * to respond.
   */
  private startDataRefreshSubscription(delayMs: number, updateWalletsFirst: boolean) {
    if (this.dataRefreshSubscription) {
      this.dataRefreshSubscription.unsubscribe();
    }

    if (this.savedWalletsList) {
      this.ngZone.runOutsideAngular(() => {
        this.dataRefreshSubscription = of(0).pipe(delay(delayMs), mergeMap(() => {
          if (updateWalletsFirst) {
            return this.refreshBalances(this.savedWalletsList, true);
          } else {
            return of(0);
          }
        }), mergeMap(() => {
          return this.refreshBalances(this.savedWalletsList, false);
        })).subscribe(() => this.startDataRefreshSubscription(10000, false), () => this.startDataRefreshSubscription(2000, false));
      });
    }
  }

  /**
   * Refreshes the wallets on walletsWithBalanceList and their balances.
   * @param wallets The current wallet lists.
   * @param forceQuickCompleteArrayUpdate If true, the balance data saved on savedBalanceData
   * will be used to set the balance of the wallet list, instead of getting the data from
   * the node. If false, the balance data is obtained from the node and savedBalanceData is
   * updated.
   */
  private refreshBalances(wallets: WalletBase[], forceQuickCompleteArrayUpdate: boolean): Observable<any> {
    if (this.gettingBalanceSubscription) {
      this.gettingBalanceSubscription.unsubscribe();
    }

    // Create a copy of the wallet list.
    const temporalWallets: WalletWithBalance[] = [];
    wallets.forEach(wallet => {
      temporalWallets.push(walletWithBalanceFromBase(wallet));
    });

    // This will help to update savedBalanceData when finishing the procedure.
    if (!forceQuickCompleteArrayUpdate) {
      this.temporalSavedBalanceData = new Map<string, any>();
    }

    // Get the balance of each wallet.
    return observableForkJoin(temporalWallets.map(wallet => this.retrieveWalletBalance(wallet, forceQuickCompleteArrayUpdate))).pipe(tap(walletHasPendingTx => {
      this.hasPendingTransactionsSubject.next(walletHasPendingTx.some(value => value));

      if (!this.walletsWithBalanceList || forceQuickCompleteArrayUpdate || this.walletsWithBalanceList.length !== temporalWallets.length) {
        // Update the whole list.
        this.walletsWithBalanceList = temporalWallets;
        this.informDataUpdated();
      } else {
        // If there is a change in the IDs of the wallet list, update the whole list.
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
          // Update only the balances with changes.
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

          // If any of the balances changed, inform that there were changes.
          if (changeDetected) {
            this.informDataUpdated();
          }
        }
      }

      if (!forceQuickCompleteArrayUpdate) {
        this.savedBalanceData = this.temporalSavedBalanceData;
        if (!this.firstFullUpdateMadeSubject.value) {
          // Inform that the service already optained the balance from the node for the first time.
          this.firstFullUpdateMadeSubject.next(true);
        }
      }
    }));
  }

  /**
   * Gets from the node the balance of a wallet and used the retrieved data to update an imstamce
   * of WalletWithBalance. It also saves the retrieved data on temporalSavedBalanceData.
   * @param wallet Wallet to update.
   * @param useSavedBalanceData If true, the balance data saved on savedBalanceData
   * will be used instead of retrieving the data from the node.
   * @returns True if there are one or more pending transactions that will affect the balance of
   * the provided walled, false otherwise. If useSavedBalanceData is true, the value of
   * hasPendingTransactionsSubject will be returned.
   */
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

  /**
   * Makes walletsWithBalanceSubject emit, to inform that the wallet list has been updated.
   */
  private informDataUpdated() {
    this.ngZone.run(() => {
      this.walletsWithBalanceSubject.next(this.walletsWithBalanceList);
    });
  }
}
