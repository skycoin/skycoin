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
import { Address, NormalTransaction, PreviewTransaction, Wallet, Output } from '../app.datatypes';
import { ReplaySubject } from 'rxjs/ReplaySubject';
import { Subscription } from 'rxjs/Subscription';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { BigNumber } from 'bignumber.js';
import { HwWalletService } from './hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { AppService } from './app.service';
import { AppConfig } from '../app.config';

declare var Cipher: any;

export enum HwSecurityWarnings {
  NeedsBackup,
  NeedsPin,
}

@Injectable()
export class WalletService {

  addresses: Address[];
  wallets: Subject<Wallet[]> = new ReplaySubject<Wallet[]>();
  pendingTxs: Subject<any[]> = new ReplaySubject<any[]>();
  dataRefreshSubscription: Subscription;

  initialLoadFailed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor(
    private appService: AppService,
    private apiService: ApiService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private ngZone: NgZone,
  ) {
    this.loadData();
    this.startDataRefreshSubscription();
  }

  addressesAsString(): Observable<string> {
    return this.allAddresses().map(addrs => addrs.map(addr => addr.address)).map(addrs => addrs.join(','));
  }

  addAddress(wallet: Wallet, num: number, password?: string) {
    if (!wallet.isHardware) {
      return this.apiService.postWalletNewAddress(wallet, num, password)
        .do(addresses => {
          addresses.forEach(value => wallet.addresses.push(value));
          this.refreshBalances();
        });
    } else {
      return this.hwWalletService.getAddresses(num, wallet.addresses.length).flatMap(response => {
        (response.rawResponse as any[]).forEach(value => wallet.addresses.push({
          address: value,
          coins: null,
          hours: null,
        }));
        this.saveHardwareWallets();
        this.refreshBalances();

        return Observable.of(response.rawResponse);
      });
    }
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

  createHardwareWallet(): Observable<Wallet> {
    let addresses: string[];
    let lastAddressWithTx = 0;
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();
    const addressesWithTxMap: Map<string, boolean> = new Map<string, boolean>();

    return this.hwWalletService.getAddresses(AppConfig.maxHardwareWalletAddresses, 0).flatMap(response => {
      addresses = response.rawResponse;
      addresses.forEach(address => {
        addressesMap.set(address, true);
      });

      const addressesString = addresses.join(',');

      return this.apiService.post('transactions', { addrs: addressesString });
    }).flatMap(response => {
      response.forEach(tx => {
        tx.txn.outputs.forEach(output => {
          if (addressesMap.has(output.dst)) {
            addressesWithTxMap.set(output.dst, true);
          }
        });
      });

      addresses.forEach((address, i) => {
        if (addressesWithTxMap.has(address)) {
          lastAddressWithTx = i;
        }
      });

      return this.wallets.first().map(wallets => {
        const newWallet = this.createHardwareWalletData(
          this.translate.instant('hardware-wallet.general.default-wallet-name'),
          addresses.slice(0, lastAddressWithTx + 1).map(add => {
            return { address: add, confirmed: false };
          }), true, false,
        );

        let lastHardwareWalletIndex = wallets.length - 1;
        for (let i = 0; i < wallets.length; i++) {
          if (!wallets[i].isHardware) {
            lastHardwareWalletIndex = i - 1;
            break;
          }
        }
        wallets.splice(lastHardwareWalletIndex + 1, 0, newWallet);
        this.saveHardwareWallets();
        this.refreshBalances();

        return newWallet;
      });
    });
  }

  updateWalletHasHwSecurityWarnings(wallet: Wallet): Observable<HwSecurityWarnings[]> {
    if (wallet.isHardware) {
      return this.hwWalletService.getFeatures().map(result => {
        const warnings: HwSecurityWarnings[] = [];

        wallet.hasHwSecurityWarnings = false;
        if (result.rawResponse.needsBackup) {
          warnings.push(HwSecurityWarnings.NeedsBackup);
          wallet.hasHwSecurityWarnings = true;
        }
        if (!result.rawResponse.pinProtection) {
          warnings.push(HwSecurityWarnings.NeedsPin);
          wallet.hasHwSecurityWarnings = true;
        }
        this.saveHardwareWallets();

        return warnings;
      });
    } else {
      return Observable.of([]);
    }
  }

  deleteHardwareWallet(wallet: Wallet): Observable<boolean> {
    if (wallet.isHardware) {
      return this.wallets.first().map(wallets => {
        const index = wallets.indexOf(wallet);
        if (index !== -1) {
          wallets.splice(index, 1);

          this.saveHardwareWallets();
          this.refreshBalances();

          return true;
        }

        return false;
      });
    }

    return null;
  }

  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').map(response => response.address);
  }

  outputs(): Observable<any> {
    return this.addressesAsString()
      .first()
      .filter(addresses => !!addresses)
      .flatMap(addresses => this.apiService.post('outputs', {addrs: addresses}));
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
    return Observable.timer(0, 10000).flatMap(() => this.apiService.get('pendingTxs', {verbose: 1}));
  }

  pendingTransactions(): Observable<any> {
    return this.pendingTxs.asObservable();
  }

  refreshBalances() {
    this.wallets.first().subscribe(wallets => {
      Observable.forkJoin(wallets.map(wallet => this.retrieveWalletBalance(wallet).map(response => {
        wallet.coins = response.coins;
        wallet.hours = response.hours;
        wallet.addresses.map(address => {
          const balance = response.addresses.find(addr => addr.address === address.address);
          address.coins = balance.coins;
          address.hours = balance.hours;
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

  resetPassword(wallet: Wallet, seed: string, password: string): Observable<Wallet> {
    const params = new Object();
    params['id'] = wallet.filename;
    params['seed'] = seed;
    if (password) {
      params['password'] = password;
    }

    return this.apiService.post('wallet/recover', params, {}, true).do(w => {
      wallet.encrypted = w.data.meta.encrypted;
      this.updateWallet(w.data);
    });
  }

  getWalletSeed(wallet: Wallet, password: string): Observable<string> {
    return this.apiService.getWalletSeed(wallet, password);
  }

  createTransaction(
    wallet: Wallet,
    addresses: string[]|null,
    unspents: string[]|null,
    destinations: any[],
    hoursSelection: any,
    changeAddress: string|null,
    password: string|null): Observable<PreviewTransaction> {

    if (unspents) {
      addresses = null;
    }

    return this.apiService.post(
      'wallet/transaction',
      {
        hours_selection: hoursSelection,
        wallet_id: wallet.filename,
        password: password,
        addresses: addresses,
        unspents: unspents,
        to: destinations,
        change_address: changeAddress,
      },
      {
        json: true,
      },
    ).map(response => {
      return {
        ...response.transaction,
        hoursBurned: new BigNumber(response.transaction.fee),
        encoded: response.encoded_transaction,
      };
    });
  }

  createHwTransaction(wallet: Wallet, address: string, amount: BigNumber): Observable<PreviewTransaction> {
    const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.appService.burnRate));
    const addresses = wallet.addresses.map(a => a.address).join(',');

    let totalHours = new BigNumber('0');
    let hoursToSend = new BigNumber('0');
    let calculatedHours = new BigNumber('0');

    let convertedOutputs: any[];

    const txOutputs = [];
    const txInputs = [];
    const hwOutputs = [];
    const hwInputs = [];

    return this.getOutputs(addresses).flatMap((outputs: Output[]) => {
        const minRequiredOutputs =  this.getMinRequiredOutputs(amount, outputs);
        let totalCoins = new BigNumber('0');
        minRequiredOutputs.map(output => totalCoins = totalCoins.plus(output.coins));

        if (totalCoins.isLessThan(amount)) {
          throw new Error(this.translate.instant('service.wallet.not-enough-hours'));
        }
        if (minRequiredOutputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs'));
        }

        minRequiredOutputs.map(output => totalHours = totalHours.plus(output.calculated_hours));
        hoursToSend = totalHours.multipliedBy(unburnedHoursRatio).dividedBy(2).decimalPlaces(0, BigNumber.ROUND_FLOOR);

        calculatedHours = totalHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);

        const changeCoins = totalCoins.minus(amount).decimalPlaces(6);

        if (changeCoins.isGreaterThan(0)) {
          txOutputs.push({
            address: wallet.addresses[0].address,
            coins: changeCoins.toNumber(),
            hours: calculatedHours.minus(hoursToSend).toNumber(),
          });

          hwOutputs.push({
            address: wallet.addresses[0].address,
            address_index: 0,
            coin: parseInt(changeCoins.multipliedBy(1000000).toFixed(0), 10),
            hour: calculatedHours.minus(hoursToSend).toNumber(),
          });
        } else {
          hoursToSend = calculatedHours;
        }

        txOutputs.push({ address: address, coins: amount.toNumber(), hours: hoursToSend.toNumber() });
        hwOutputs.push({ address: address, coin: parseInt(amount.multipliedBy(1000000).toFixed(0), 10), hour: hoursToSend.toNumber() });

        if (address === wallet.addresses[0].address) {
          hoursToSend = calculatedHours;
        }

        minRequiredOutputs.forEach(input => {
          txInputs.push({
            hash: input.hash,
            secret: '',
            address: input.address,
            address_index: wallet.addresses.findIndex(a => a.address === input.address),
            calculated_hours: input.calculated_hours.toNumber(),
            coins: input.coins.toNumber(),
          });

          hwInputs.push({
            hashIn: input.hash,
            index: txInputs[txInputs.length - 1].address_index,
          });
        });

        convertedOutputs = txOutputs.map(output => {
          return {
            ...output,
            coins: parseInt(new BigNumber(output.coins).multipliedBy(1000000).toFixed(0), 10),
          };
        });

        // Impossible at this time. Here waiting for when the possibility of sending to multiple addresses is added.
        if (convertedOutputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-outputs'));
        }

        return this.hwWalletService.signTransaction(hwInputs, hwOutputs);
      }).flatMap(signatures => {
        const rawTransaction = Cipher.PrepareTransactionWithSignatures(
          JSON.stringify(txInputs),
          JSON.stringify(convertedOutputs),
          JSON.stringify(signatures.rawResponse),
        );

        return Observable.of({
          balance: amount,
          inputs: txInputs,
          outputs: txOutputs,
          txid: null,
          from: '',
          to: [],
          hoursSent: hoursToSend,
          hoursBurned: totalHours.minus(calculatedHours),
          encoded: rawTransaction,
        });
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
    let wallets: Wallet[];
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();


    return this.wallets.first().flatMap(w => {
      wallets = w;

      return this.allAddresses().first();
    }).flatMap(addresses => {
      this.addresses = addresses;
      addresses.map(add => addressesMap.set(add.address, true));

      return this.apiService.getTransactions(addresses);
    }).map(transactions => {
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

  saveHardwareWallets() {
    this.wallets.first().subscribe(wallets => {
      const hardwareWallets: Wallet[] = [];

      wallets.map(wallet => {
        if (wallet.isHardware) {
          hardwareWallets.push(this.createHardwareWalletData(
            wallet.label,
            wallet.addresses.map(address => {
              return { address: address.address, confirmed: address.confirmed };
            }),
            wallet.hasHwSecurityWarnings,
            wallet.stopShowingHwSecurityPopup,
          ));
        }
      });

      this.hwWalletService.saveWalletsDataSync(JSON.stringify(hardwareWallets));

      this.wallets.next(wallets);
    });
  }

  getWalletUnspentOutputs(wallet: Wallet): Observable<Output[]> {
    const addresses = wallet.addresses.map(a => a.address).join(',');

    return this.getOutputs(addresses);
  }

  private createHardwareWalletData(label: string, addresses: {address: string, confirmed: boolean}[], hasHwSecurityWarnings: boolean, stopShowingHwSecurityPopup: boolean): Wallet {
    return {
      label: label,
      filename: '',
      hasHwSecurityWarnings: hasHwSecurityWarnings,
      stopShowingHwSecurityPopup: stopShowingHwSecurityPopup,
      coins: null,
      hours: null,
      addresses: addresses.map(address => {
        return {
          address: address.address,
          coins: null,
          hours: null,
          confirmed: address.confirmed,
        };
      }),
      encrypted: false,
      isHardware: true,
    };
  }

  private loadData(): void {
    this.apiService.getWallets().first().subscribe(
      recoveredWallets => {
        let wallets: Wallet[] = [];
        if (this.hwWalletService.hwWalletCompatibilityActivated) {
          this.loadHardwareWallets(wallets);
        }
        wallets = wallets.concat(recoveredWallets);
        this.wallets.next(wallets);
      },
      () => this.initialLoadFailed.next(true),
    );
  }

  private loadHardwareWallets(wallets: Wallet[]) {
    const storedWallets: string = this.hwWalletService.getSavedWalletsDataSync();
    if (storedWallets) {
      const loadedWallets: Wallet[] = JSON.parse(storedWallets);
      loadedWallets.map(wallet => wallets.push(wallet));
    }
  }

  private retrieveInputAddress(input: string) {
    return this.apiService.get('uxout', {uxid: input});
  }

  private retrieveWalletBalance(wallet: Wallet): Observable<any> {
    let query: Observable<any>;
    if (!wallet.isHardware) {
      query = this.apiService.get('wallet/balance', { id: wallet.filename });
    } else {
      const formattedAddresses = wallet.addresses.map(a => a.address).join(',');
      query = this.apiService.post('balance', { addrs: formattedAddresses });
    }

    return query.map(balance => {
      return {
        coins: new BigNumber(balance.confirmed.coins).dividedBy(1000000),
        hours: new BigNumber(balance.confirmed.hours),
        addresses: Object.keys(balance.addresses).map(address => ({
          address,
          coins: new BigNumber(balance.addresses[address].confirmed.coins).dividedBy(1000000),
          hours: new BigNumber(balance.addresses[address].confirmed.hours),
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
    this.apiService.get('pendingTxs', { verbose: true })
      .flatMap((transactions: any) => {
        if (transactions.length === 0) {
          return Observable.of([]);
        }

        return this.wallets.first().map((wallets: Wallet[]) => {
          const walletAddresses = new Set<string>();
          wallets.forEach(wallet => {
            wallet.addresses.forEach(address => walletAddresses.add(address.address));
          });

          return transactions.filter(tran => {
            return tran.transaction.inputs.some(input => walletAddresses.has(input.owner)) ||
            tran.transaction.outputs.some(output => walletAddresses.has(output.dst));
          });
        });
      })
      .subscribe(transactions => this.pendingTxs.next(transactions));
  }

  private getOutputs(addresses): Observable<Output[]> {
    if (!addresses) {
      return Observable.of([]);
    } else {
      return this.apiService.post('outputs', { addrs: addresses }).map((response) => {
        const outputs = [];
        response.head_outputs.forEach(output => outputs.push({
          address: output.address,
          coins: new BigNumber(output.coins),
          hash: output.hash,
          calculated_hours: new BigNumber(output.calculated_hours),
        }));

        return outputs;
      });
    }
  }

  private sortOutputs(outputs: Output[], highestToLowest: boolean) {
    outputs.sort((a, b) => {
      if (b.coins.isGreaterThan(a.coins)) {
        return highestToLowest ? 1 : -1;
      } else if (b.coins.isLessThan(a.coins)) {
        return highestToLowest ? -1 : 1;
      } else {
        if (b.calculated_hours.isGreaterThan(a.calculated_hours)) {
          return -1;
        } else if (b.calculated_hours.isLessThan(a.calculated_hours)) {
          return 1;
        } else {
          return 0;
        }
      }
    });
  }

  private getMinRequiredOutputs(transactionAmount: BigNumber, outputs: Output[]): Output[] {

    // Split the outputs into those with and without hours
    const outputsWithHours: Output[] = [];
    const outputsWitouthHours: Output[] = [];
    outputs.forEach(output => {
      if (output.calculated_hours.isGreaterThan(0)) {
        outputsWithHours.push(output);
      } else {
        outputsWitouthHours.push(output);
      }
    });

    // Abort if there are no outputs with non-zero coinhours.
    if (outputsWithHours.length === 0) {
      return [];
    }

    // Sort the outputs with hours by coins, from highest to lowest. If two items have the same amount of
    // coins, the one with the least hours is placed first.
    this.sortOutputs(outputsWithHours, true);

    // Use the first nonzero output.
    const minRequiredOutputs: Output[] = [outputsWithHours[0]];
    let sumCoins: BigNumber = new BigNumber(outputsWithHours[0].coins);

    // If it's enough, finish.
    if (sumCoins.isGreaterThanOrEqualTo(transactionAmount)) {
      return minRequiredOutputs;
    }

    // Sort the outputs without hours by coins, from lowest to highest.
    this.sortOutputs(outputsWitouthHours, false);

    // Add the outputs without hours, until having the necessary amount of coins.
    outputsWitouthHours.forEach(output => {
      if (sumCoins.isLessThan(transactionAmount)) {
        minRequiredOutputs.push(output);
        sumCoins = sumCoins.plus(output.coins);
      }
    });

    // If it's enough, finish.
    if (sumCoins.isGreaterThanOrEqualTo(transactionAmount)) {
      return minRequiredOutputs;
    }

    outputsWithHours.splice(0, 1);
    // Sort the outputs with hours by coins, from lowest to highest.
    this.sortOutputs(outputsWithHours, false);

    // Add the outputs with hours, until having the necessary amount of coins.
    outputsWithHours.forEach((output) => {
      if (sumCoins.isLessThan(transactionAmount)) {
        minRequiredOutputs.push(output);
        sumCoins = sumCoins.plus(output.coins);
      }
    });

    return minRequiredOutputs;
  }
}
