import { Injectable } from '@angular/core';
import { Observable, forkJoin, of } from 'rxjs';
import { mergeMap, map, first } from 'rxjs';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api.service';
import { Address, NormalTransaction, Wallet } from '../../app.datatypes';
import { WalletService } from './wallet.service';
import { GlobalsService } from '../globals.service';
import { isEqualOrSuperiorVersion } from '../../utils/semver';

@Injectable()
export class HistoryService {

  constructor(
    private apiService: ApiService,
    private walletService: WalletService,
    private globalsService: GlobalsService
  ) { }

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

        return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
          let TxObsv: Observable<any>;
          if (isEqualOrSuperiorVersion(version, '0.25.0')) {
            TxObsv = this.retrieveAddressesTransactions(addresses).pipe(map(transactions => {
              return transactions.sort((a, b) =>  b.timestamp - a.timestamp);
            }));
          } else {
            TxObsv = forkJoin(addresses.map(address => this.retrieveAddressTransactions(address))).pipe(
              map(transactions => {
                return []
                  .concat.apply([], transactions)
                  .reduce((array, item) => {
                    if (!array.find(trans => trans.txid === item.txid)) {
                      array.push(item);
                    }

                    return array;
                  }, [])
                  .sort((a, b) =>  b.timestamp - a.timestamp);
              }));
          }

          return TxObsv.pipe(map(transactions => {
            return transactions.map(transaction => {
              const outgoing = transaction.inputs.some(input => addressesMap.has(input.owner));

              const relevantAddresses: Map<string, boolean> = new Map<string, boolean>();
              transaction.balance = new BigNumber('0');
              transaction.hoursSent = new BigNumber('0');

              if (!outgoing) {
                transaction.outputs.map(output => {
                  if (addressesMap.has(output.dst)) {
                    relevantAddresses.set(output.dst, true);
                    transaction.balance = transaction.balance.plus(output.coins);
                    transaction.hoursSent = transaction.hoursSent.plus(output.hours);
                  }
                });
              } else {
                const possibleReturnAddressesMap: Map<string, boolean> = new Map<string, boolean>();
                transaction.inputs.map(input => {
                  if (addressesMap.has(input.owner)) {
                    relevantAddresses.set(input.owner, true);
                    wallets.map(wallet => {
                      if (wallet.addresses.some(add => add.address === input.owner)) {
                        wallet.addresses.map(add => possibleReturnAddressesMap.set(add.address, true));
                      }
                    });
                  }
                });

                transaction.outputs.map(output => {
                  if (!possibleReturnAddressesMap.has(output.dst)) {
                    transaction.balance = transaction.balance.minus(output.coins);
                    transaction.hoursSent = transaction.hoursSent.plus(output.hours);
                  }
                });

                if (transaction.balance.isEqualTo(0)) {
                  transaction.coinsMovedInternally = true;
                  const inputAddressesMap: Map<string, boolean> = new Map<string, boolean>();

                  transaction.inputs.map(input => {
                    inputAddressesMap.set(input.owner, true);
                  });

                  transaction.outputs.map(output => {
                    if (!inputAddressesMap.has(output.dst)) {
                      relevantAddresses.set(output.dst, true);
                      transaction.balance = transaction.balance.plus(output.coins);
                      transaction.hoursSent = transaction.hoursSent.plus(output.hours);
                    }
                  });
                }
              }

              relevantAddresses.forEach((value, key) => {
                transaction.addresses.push(key);
              });

              let inputsHours = new BigNumber('0');
              transaction.inputs.map(input => inputsHours = inputsHours.plus(new BigNumber(input.calculated_hours)));
              let outputsHours = new BigNumber('0');
              transaction.outputs.map(output => outputsHours = outputsHours.plus(new BigNumber(output.hours)));
              transaction.hoursBurned = inputsHours.minus(outputsHours);

              return transaction;
            });
          }));
        }));
      }));
  }

  retrieveAddressTransactions(address: Address): Observable<NormalTransaction[]> {
    return this.apiService.get('explorer/address', { address: address.address }).pipe(
      map(transactions => transactions.map(transaction => ({
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
      map(transactions => transactions.map(transaction => ({
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
        return this.apiService.get('pendingTxs', { verbose: true });
      } else {
        return this.apiService.get('pendingTxs');
      }
    }));
  }

  deletePendingTransaction(txid: string): Observable<any> {
    return this.apiService.delete('pendingTxs', { txid });
  }

  getTransactionDetails(uxid: string): Observable<any> {
    return this.apiService.get('uxout', { uxid: uxid });
  }
}
