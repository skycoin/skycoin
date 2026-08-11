import { Injectable } from '@angular/core';
import { Observable, forkJoin, of } from 'rxjs';
import { mergeMap, map, first } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { toBigNumber } from '../../utils/converters';

import { ApiService } from '../api.service';
import { CoinService } from '../coin.service';
import { Address, NormalTransaction, Wallet } from '../../app.datatypes';
import { BaseCoin } from '../../coins/basecoin';
import { WalletService } from './wallet.service';
import { GlobalsService } from '../globals.service';
import { isEqualOrSuperiorVersion } from '../../utils/semver';

@Injectable()
export class HistoryService {
  private currentCoin!: BaseCoin;

  constructor(
    private apiService: ApiService,
    private walletService: WalletService,
    private globalsService: GlobalsService,
    coinService: CoinService
  ) {
    coinService.currentCoin.subscribe(coin => this.currentCoin = coin);
  }

  transactions(): Observable<any[]> {
    let wallets: Wallet[];
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();


    return this.walletService.wallets.pipe(first(), mergeMap(w => {
      wallets = w;

      return this.walletService.addresses.pipe(first());
    }), mergeMap(addresses => {
        if (addresses.length === 0) {
          return of([]);
        }

        addresses.map(add => addressesMap.set(add.address, true));

        // Bitcoin uses its own history endpoint
        if (this.currentCoin && this.currentCoin.isBitcoin()) {
          return this.retrieveBitcoinHistory(addresses, addressesMap);
        }

        return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
          let TxObsv: Observable<any>;
          if (isEqualOrSuperiorVersion(version, '0.25.0')) {
            TxObsv = this.retrieveAddressesTransactions(addresses).pipe(map(transactions => {
              return transactions.sort((a: any, b: any) =>  b.timestamp - a.timestamp);
            }));
          } else {
            TxObsv = forkJoin(addresses.map(address => this.retrieveAddressTransactions(address))).pipe(
              map(transactions => {
                return ([] as any[]).concat(...transactions)
                  .reduce((array, item) => {
                    if (!array.find((trans: any) => trans.txid === item.txid)) {
                      array.push(item);
                    }

                    return array;
                  }, [])
                  .sort((a: any, b: any) =>  b.timestamp - a.timestamp);
              }));
          }

          return TxObsv.pipe(map(transactions => {
            return transactions.map((transaction: any) => {
              const outgoing = transaction.inputs.some((input: any) => addressesMap.has(input.owner));

              const relevantAddresses: Map<string, boolean> = new Map<string, boolean>();
              transaction.balance = new BigNumber('0');
              transaction.hoursSent = new BigNumber('0');

              if (!outgoing) {
                transaction.outputs.map((output: any) => {
                  if (addressesMap.has(output.dst)) {
                    relevantAddresses.set(output.dst, true);
                    transaction.balance = transaction.balance.plus(output.coins);
                    transaction.hoursSent = transaction.hoursSent.plus(toBigNumber(output.hours));
                  }
                });
              } else {
                const possibleReturnAddressesMap: Map<string, boolean> = new Map<string, boolean>();
                transaction.inputs.map((input: any) => {
                  if (addressesMap.has(input.owner)) {
                    relevantAddresses.set(input.owner, true);
                    wallets.map(wallet => {
                      if (wallet.addresses.some(add => add.address === input.owner)) {
                        wallet.addresses.map(add => possibleReturnAddressesMap.set(add.address, true));
                      }
                    });
                  }
                });

                transaction.outputs.map((output: any) => {
                  if (!possibleReturnAddressesMap.has(output.dst)) {
                    transaction.balance = transaction.balance.minus(output.coins);
                    transaction.hoursSent = transaction.hoursSent.plus(toBigNumber(output.hours));
                  }
                });

                if (transaction.balance.isEqualTo(0)) {
                  transaction.coinsMovedInternally = true;
                  const inputAddressesMap: Map<string, boolean> = new Map<string, boolean>();

                  transaction.inputs.map((input: any) => {
                    inputAddressesMap.set(input.owner, true);
                  });

                  transaction.outputs.map((output: any) => {
                    if (!inputAddressesMap.has(output.dst)) {
                      relevantAddresses.set(output.dst, true);
                      transaction.balance = transaction.balance.plus(output.coins);
                      transaction.hoursSent = transaction.hoursSent.plus(toBigNumber(output.hours));
                    }
                  });
                }
              }

              relevantAddresses.forEach((value, key) => {
                transaction.addresses.push(key);
              });

              let inputsHours = new BigNumber('0');
              transaction.inputs.map((input: any) => inputsHours = inputsHours.plus(toBigNumber(input.calculated_hours)));
              let outputsHours = new BigNumber('0');
              transaction.outputs.map((output: any) => outputsHours = outputsHours.plus(toBigNumber(output.hours)));
              transaction.hoursBurned = inputsHours.minus(outputsHours);

              return transaction;
            });
          }));
        }));
      }));
  }

  private retrieveBitcoinHistory(addresses: Address[], addressesMap: Map<string, boolean>): Observable<any[]> {
    const formattedAddresses = addresses.map(a => a.address).join(',');
    return this.apiService.get('btc/history', { addrs: formattedAddresses } as any).pipe(
      map((transactions: any[]) => {
        return (transactions || []).map(tx => {
          const outgoing = tx.inputs && tx.inputs.some((input: any) => addressesMap.has(input.address));

          let balance = new BigNumber('0');
          if (!outgoing) {
            (tx.outputs || []).forEach((output: any) => {
              if (addressesMap.has(output.address)) {
                balance = balance.plus(new BigNumber(output.value));
              }
            });
          } else {
            (tx.outputs || []).forEach((output: any) => {
              if (!addressesMap.has(output.address)) {
                balance = balance.minus(new BigNumber(output.value));
              }
            });
          }

          return {
            txid: tx.txid,
            addresses: [],
            balance: balance,
            timestamp: tx.timestamp || 0,
            block: tx.block_height || 0,
            confirmed: (tx.confirmations || 0) > 0,
            inputs: tx.inputs || [],
            outputs: tx.outputs || [],
            hoursSent: new BigNumber('0'),
            hoursBurned: new BigNumber('0'),
            fee: tx.fee || 0,
          };
        }).sort((a: any, b: any) => b.timestamp - a.timestamp);
      })
    );
  }

  retrieveAddressTransactions(address: Address): Observable<NormalTransaction[]> {
    return this.apiService.get('explorer/address', { address: address.address } as any).pipe(
      map(transactions => transactions.map((transaction: any) => ({
        addresses: [],
        balance: new BigNumber('0'),
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.timestamp,
        txid: transaction.txid,
        inputs: transaction.inputs,
        outputs: transaction.outputs,
      }))));
  }

  retrieveAddressesTransactions(addresses: Address[]): Observable<NormalTransaction[]> {
    const formattedAddresses = addresses.map(a => a.address).join(',');

    return this.apiService.post('transactions', { addrs: formattedAddresses, verbose: true }).pipe(
      map(transactions => transactions.map((transaction: any) => ({
        addresses: [],
        balance: new BigNumber('0'),
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.txn.timestamp,
        txid: transaction.txn.txid,
        inputs: transaction.txn.inputs,
        outputs: transaction.txn.outputs,
      }))));
  }

  getAllPendingTransactions(): Observable<any> {
    return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
      if (isEqualOrSuperiorVersion(version, '0.25.0')) {
        return this.apiService.get('pendingTxs', { verbose: true } as any);
      } else {
        return this.apiService.get('pendingTxs');
      }
    }));
  }

  deletePendingTransaction(txid: string): Observable<any> {
    return this.apiService.delete('pendingTxs', { txid } as any);
  }

  getTransactionDetails(uxid: string): Observable<any> {
    return this.apiService.get('uxout', { uxid: uxid } as any);
  }
}
