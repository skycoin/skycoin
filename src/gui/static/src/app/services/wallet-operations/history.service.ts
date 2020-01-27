import { of, Observable } from 'rxjs';
import { first, mergeMap, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { ApiService } from '../api.service';
import { NormalTransaction } from '../../app.datatypes';
import { BigNumber } from 'bignumber.js';
import { StorageService, StorageType } from '../storage.service';
import { WalletBase, AddressBase } from './wallet-objects';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';

export interface PendingTransactions {
  user: any[];
  all: any[];
}

@Injectable()
export class HistoryService {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
    private storageService: StorageService,
  ) { }

  getTransactionsHistory(): Observable<NormalTransaction[]> {
    let wallets: WalletBase[];
    let transactions: NormalTransaction[];
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();

    return this.walletsAndAddressesService.allWallets.pipe(first(), mergeMap(w => {
      wallets = w;

      const addresses: AddressBase[] = [];
      wallets.forEach(wallet => {
        wallet.addresses.map(add => {
          addresses.push(add);
          addressesMap.set(add.address, true);
        });
      });

      const formattedAddresses = addresses.map(a => a.address).join(',');

      return this.apiService.post('transactions', {addrs: formattedAddresses, verbose: true});
    }), map((response: any[]) => {
      return response.map<NormalTransaction>(transaction => ({
        addresses: [],
        balance: new BigNumber(0),
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.txn.timestamp,
        txid: transaction.txn.txid,
        inputs: transaction.txn.inputs,
        outputs: transaction.txn.outputs,
      }));
    }), mergeMap((recoveredTransactions: NormalTransaction[]) => {
      transactions = recoveredTransactions;

      return this.storageService.get(StorageType.NOTES, null);
    }), map(notes => {
      const notesMap: Map<string, string> = new Map<string, string>();
      Object.keys(notes.data).forEach(key => {
        notesMap.set(key, notes.data[key]);
      });

      return transactions
        .sort((a, b) =>  b.timestamp - a.timestamp)
        .map(transaction => {
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

          const txNote = notesMap.get(transaction.txid);
          if (txNote) {
            transaction.note = txNote;
          }

          return transaction;
        });
    }));
  }

  getPendingTransactions(): Observable<PendingTransactions> {
    return this.apiService.get('pendingTxs', { verbose: true }).pipe(
      mergeMap((transactions: any) => {
        if (transactions.length === 0) {
          return of({
            user: [],
            all: [],
          });
        }

        return this.walletsAndAddressesService.allWallets.pipe(first(), map((wallets: WalletBase[]) => {
          const walletAddresses = new Set<string>();
          wallets.forEach(wallet => {
            wallet.addresses.forEach(address => walletAddresses.add(address.address));
          });

          const userTransactions = transactions.filter(tran => {
            return tran.transaction.inputs.some(input => walletAddresses.has(input.owner)) ||
            tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
          });

          return {
            user: userTransactions,
            all: transactions,
          };
        }));
      }));
  }
}
