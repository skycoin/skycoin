import { of, Observable } from 'rxjs';
import { first, mergeMap, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { BigNumber } from 'bignumber.js';

import * as moment from 'moment';

import { ApiService } from '../api.service';
import { StorageService, StorageType } from '../storage.service';
import { WalletBase, AddressBase } from './wallet-objects';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { OldTransaction } from './transaction-objects';

export interface PendingTransactionsResponse {
  /**
   * Pending transactions affecting one or more of the user addresses.
   */
  user: PendingTransactionData[];
  /**
   * All pending transactions known by the node, including the ones affecting one
   * or more of the user addresses.
   */
  all: PendingTransactionData[];
}

export interface PendingTransactionData {
  /**
   * Transaction ID.
   */
  id: string;
  /**
   * How many coins are on the outputs.
   */
  coins: string;
  /**
   * How many hours are on the outputs.
   */
  hours: string;
  /**
   * Transaction timestamp, in Unix format.
   */
  timestamp: number;
}

/**
 * Allows to get the transaction history and pending transactions.
 */
@Injectable()
export class HistoryService {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private apiService: ApiService,
    private storageService: StorageService,
  ) { }

  /**
   * Gets the transaction history of all the wallets or a specific wallet.
   * @param wallet Specific wallet for which the transaction history will be returned. If null,
   * the transactions of all wallets will be returned.
   */
  getTransactionsHistory(wallet: WalletBase|null): Observable<OldTransaction[]> {
    let transactions: OldTransaction[];
    /**
     * Allows to easily know which addresses are part of the user wallets and also to know
     * which wallet the address belong to.
     */
    const addressesMap: Map<string, WalletBase> = new Map<string, WalletBase>();

    // Use the provided wallet or get all wallets.
    let initialRequest: Observable<WalletBase[]>;
    if (wallet) {
      initialRequest = of([wallet]);
    } else {
      initialRequest = this.walletsAndAddressesService.allWallets;
    }

    return initialRequest.pipe(first(), mergeMap(wallets => {
      const addresses: AddressBase[] = [];
      wallets.forEach(w => {
        w.addresses.map(add => {
          addresses.push(add);
          // There could be more than one wallet with the address. This would happen if the wallet is repeated
          // (like when using the same seed for a software and a hardware wallet). In that case, the wallet
          // with most addresses is considered "the most complete one" and is used.
          if (!addressesMap.has(add.address) || addressesMap.get(add.address).addresses.length < w.addresses.length) {
            addressesMap.set(add.address, w);
          }
        });
      });

      // Get the transactions for all addresses.
      const formattedAddresses = addresses.map(a => a.address).join(',');

      return this.apiService.post('transactions', {addrs: formattedAddresses, verbose: true});
    }), mergeMap((response: any[]) => {
      // Process the response and convert it into a known object.
      transactions = response.map<OldTransaction>(transaction => ({
        relevantAddresses: [],
        balance: new BigNumber(0),
        hoursBalance: new BigNumber(0),
        hoursBurned: new BigNumber(0),
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.txn.timestamp,
        id: transaction.txn.txid,
        inputs: (transaction.txn.inputs as any[]).map(input => {
          return {
            hash: input.uxid,
            address: input.owner,
            coins: new BigNumber(input.coins),
            hours: new BigNumber(input.calculated_hours),
          };
        }),
        outputs: (transaction.txn.outputs as any[]).map(output => {
          return {
            hash: output.uxid,
            address: output.dst,
            coins: new BigNumber(output.coins),
            hours: new BigNumber(output.hours),
          };
        }),
      }));

      // Get the transaction notes.
      return this.storageService.get(StorageType.NOTES, null);
    }), map(notes => {
      const notesMap: Map<string, string> = new Map<string, string>();
      Object.keys(notes.data).forEach(key => {
        notesMap.set(key, notes.data[key]);
      });

      return transactions
        // Sort the transactions by date.
        .sort((a, b) =>  b.timestamp - a.timestamp)
        .map(transaction => {
          // Consider the transaction outgoing if the user owns any of the inputs.
          // There is code below to check if the tx was not outgoing, but an internal one.
          const isOutgoingTransaction = transaction.inputs.some(input => addressesMap.has(input.address));

          // Saves the user address which sent or received coins, depending on the transaction type.
          const involvedLocalAddresses: Map<string, boolean> = new Map<string, boolean>();

          if (!isOutgoingTransaction) {
            transaction.outputs.map(output => {
              // If the transactions is an incoming one, all coins and hours on outputs
              // pointing to local addresses are considered received.
              if (addressesMap.has(output.address)) {
                involvedLocalAddresses.set(output.address, true);
                transaction.balance = transaction.balance.plus(output.coins);
                transaction.hoursBalance = transaction.hoursBalance.plus(output.hours);
              }
            });
          } else {
            // If the transaction is an outgoing one, all addresses of all wallets used for inputs
            // are considered potential return addresses, so all coins sent to those addresses
            // will be excluded when counting how many coins and hours were sent.
            const possibleReturnAddressesMap: Map<string, boolean> = new Map<string, boolean>();
            transaction.inputs.map(input => {
              if (addressesMap.has(input.address)) {
                involvedLocalAddresses.set(input.address, true);
                addressesMap.get(input.address).addresses.map(add => possibleReturnAddressesMap.set(add.address, true));
              }
            });

            // Sum all coins and hours that were sent.
            transaction.outputs.map(output => {
              if (!possibleReturnAddressesMap.has(output.address)) {
                transaction.balance = transaction.balance.minus(output.coins);
                transaction.hoursBalance = transaction.hoursBalance.plus(output.hours);
              }
            });

            // If the result is 0, all coins were sent to addrreses which are part of the same
            // wallets used to send the coins, so the transaction was not an outgoing one, but
            // just one for moving coins internally.
            if (transaction.balance.isEqualTo(0)) {
              transaction.coinsMovedInternally = true;
              const inputAddressesMap: Map<string, boolean> = new Map<string, boolean>();

              transaction.inputs.map(input => {
                inputAddressesMap.set(input.address, true);
              });

              // Sum how many coins and hours were moved to addresses different to the ones which
              // own the inputs.
              transaction.outputs.map(output => {
                if (!inputAddressesMap.has(output.address)) {
                  involvedLocalAddresses.set(output.address, true);
                  transaction.balance = transaction.balance.plus(output.coins);
                  transaction.hoursBalance = transaction.hoursBalance.plus(output.hours);
                }
              });
            }
          }

          // Create the list of addresses which received the coins and hours or the addresses
          // used for sending them, depending on whether the transaction was an incoming or
          // outgoing one.
          involvedLocalAddresses.forEach((value, key) => {
            transaction.relevantAddresses.push(key);
          });

          // Calculate how many hours were burned.
          let inputsHours = new BigNumber('0');
          transaction.inputs.map(input => inputsHours = inputsHours.plus(new BigNumber(input.hours)));
          let outputsHours = new BigNumber('0');
          transaction.outputs.map(output => outputsHours = outputsHours.plus(new BigNumber(output.hours)));
          transaction.hoursBurned = inputsHours.minus(outputsHours);

          const txNote = notesMap.get(transaction.id);
          if (txNote) {
            transaction.note = txNote;
          }

          return transaction;
        });
    }));
  }

  /**
   * Get the list of pending transactions currently on the node.
   */
  getPendingTransactions(): Observable<PendingTransactionsResponse> {
    return this.apiService.get('pendingTxs', { verbose: true }).pipe(
      mergeMap((transactions: any[]) => {
        // Default response if no transactions were found.
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

          // Build an array with the transactions affecting the user.
          const userTransactions = transactions.filter(tran => {
            return tran.transaction.inputs.some(input => walletAddresses.has(input.owner)) ||
              tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
          });

          return {
            user: userTransactions.map(tx => this.processTransactionData(tx)).sort((a, b) => b.timestamp - a.timestamp),
            all: transactions.map(tx => this.processTransactionData(tx)).sort((a, b) => b.timestamp - a.timestamp),
          };
        }));
      }));
  }

  /**
   * Converts a pending transaction returned by the server to a PendingTransactionData instance.
   * @param transaction Transaction returned by the server.
   */
  private processTransactionData(transaction: any): PendingTransactionData {
    let coins = new BigNumber('0');
    let hours = new BigNumber('0');
    transaction.transaction.outputs.map(output => {
      coins = coins.plus(output.coins);
      hours = hours.plus(output.hours);
    });

    return {
      coins: coins.toString(),
      hours: hours.toString(),
      timestamp: moment(transaction.received).unix(),
      id: transaction.transaction.txid,
    };
  }
}
