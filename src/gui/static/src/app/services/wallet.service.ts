import { forkJoin as observableForkJoin, throwError as observableThrowError, zip, of, timer, Subject, Observable, ReplaySubject, Subscription, BehaviorSubject } from 'rxjs';
import { concat, delay, filter, retryWhen, first, take, tap, mergeMap, catchError, map } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Address, NormalTransaction, PreviewTransaction, Wallet, Output } from '../app.datatypes';
import { BigNumber } from 'bignumber.js';
import { HwWalletService, HwOutput, HwInput } from './hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';
import { HttpClient } from '@angular/common/http';
import { StorageService, StorageType } from './storage.service';
import { TxEncoder } from '../utils/tx-encoder';

export enum HwSecurityWarnings {
  NeedsBackup,
  NeedsPin,
  FirmwareVersionNotVerified,
  OutdatedFirmware,
}

export interface HwFeaturesResponse {
  features: any;
  securityWarnings: HwSecurityWarnings[];
  walletNameUpdated: boolean;
}

export interface PendingTransactions {
  user: any[];
  all: any[];
}

@Injectable()
export class WalletService {

  addresses: Address[];
  wallets: Subject<Wallet[]> = new ReplaySubject<Wallet[]>(1);
  pendingTxs: Subject<PendingTransactions> = new ReplaySubject<PendingTransactions>(1);
  dataRefreshSubscription: Subscription;

  initialLoadFailed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor(
    private apiService: ApiService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private ngZone: NgZone,
    private http: HttpClient,
    private storageService: StorageService,
  ) {
    this.loadData();
    this.startDataRefreshSubscription();
  }

  addressesAsString(): Observable<string> {
    return this.allAddresses().pipe(map(addrs => addrs.map(addr => addr.address)), map(addrs => addrs.join(',')));
  }

  addAddress(wallet: Wallet, num: number, password?: string) {
    if (!wallet.isHardware) {
      return this.apiService.postWalletNewAddress(wallet, num, password).pipe(
        tap(addresses => {
          addresses.forEach(value => wallet.addresses.push(value));
          this.refreshBalances();
        }));
    } else {
      return this.hwWalletService.getAddresses(num, wallet.addresses.length).pipe(mergeMap(response => {
        (response.rawResponse as any[]).forEach(value => wallet.addresses.push({
          address: value,
          coins: null,
          hours: null,
        }));
        this.saveHardwareWallets();
        this.refreshBalances();

        return of(response.rawResponse);
      }));
    }
  }

  scanAddresses(wallet: Wallet, password?: string) {
    if (!wallet.isHardware) {
      return this.apiService.postWalletScan(wallet, password).pipe(
        map((addresses: any[]) => {
          if (addresses && addresses.length > 0) {
            addresses.forEach(address => {
              const currentAddress = new Address();
              currentAddress.address = address;
              wallet.addresses.push(currentAddress);
            });

            this.refreshBalances();

            return true;
          } else {
            return false;
          }
        }),
      );
    } else {
      // Not implemented.
    }
  }

  all(): Observable<Wallet[]> {
    return this.wallets.asObservable();
  }

  allAddresses(): Observable<any[]> {
    return this.all().pipe(map(wallets => wallets.reduce((array, wallet) => array.concat(wallet.addresses), [])));
  }

  create(label, seed, scan, password) {
    seed = seed.replace(/(\n|\r\n)$/, '');

    return this.apiService.postWalletCreate(label ? label : 'undefined', seed, scan ? scan : 100, password, 'deterministic').pipe(
      tap(wallet => {
        this.wallets.pipe(first()).subscribe(wallets => {
          wallets.push(wallet);
          this.wallets.next(wallets);
          this.refreshBalances();
        });
      }));
  }

  createHardwareWallet(): Observable<Wallet> {
    let addresses: string[];
    let lastAddressWithTx = 0;
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();
    const addressesWithTxMap: Map<string, boolean> = new Map<string, boolean>();

    return this.hwWalletService.getAddresses(AppConfig.maxHardwareWalletAddresses, 0).pipe(mergeMap(response => {
      addresses = response.rawResponse;
      addresses.forEach(address => {
        addressesMap.set(address, true);
      });

      const addressesString = addresses.join(',');

      return this.apiService.post('transactions', { addrs: addressesString });
    }), mergeMap(response => {
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

      return this.wallets.pipe(first(), map(wallets => {
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
      }));
    }));
  }

  getHwFeaturesAndUpdateData(wallet: Wallet): Observable<HwFeaturesResponse> {
    if (!wallet || wallet.isHardware) {

      let lastestFirmwareVersion: string;

      return this.http.get(AppConfig.urlForHwWalletVersionChecking, { responseType: 'text' }).pipe(
      catchError(() => of(null)),
      mergeMap((res: any) => {
        if (res) {
          lastestFirmwareVersion = res;
        } else {
          lastestFirmwareVersion = null;
        }

        return this.hwWalletService.getFeatures();
      }),
      map(result => {
        let lastestFirmwareVersionReaded = false;
        let firmwareUpdated = false;

        if (lastestFirmwareVersion) {
          lastestFirmwareVersion = lastestFirmwareVersion.trim();
          const versionParts = lastestFirmwareVersion.split('.');

          if (versionParts.length === 3) {
            lastestFirmwareVersionReaded = true;

            const numVersionParts = versionParts.map(value => Number.parseInt(value.replace(/\D/g, ''), 10));

            const devMajorVersion = result.rawResponse.fw_major;
            const devMinorVersion = result.rawResponse.fw_minor;
            const devPatchVersion = result.rawResponse.fw_patch;

            if (devMajorVersion > numVersionParts[0]) {
              firmwareUpdated = true;
            } else {
              if (devMajorVersion === numVersionParts[0]) {
                if (devMinorVersion > numVersionParts[1]) {
                  firmwareUpdated = true;
                } else {
                  if (devMinorVersion === numVersionParts[1] && devPatchVersion >= numVersionParts[2]) {
                    firmwareUpdated = true;
                  }
                }
              }
            }
          }
        }

        const warnings: HwSecurityWarnings[] = [];
        let hasHwSecurityWarnings = false;

        if (result.rawResponse.needs_backup) {
          warnings.push(HwSecurityWarnings.NeedsBackup);
          hasHwSecurityWarnings = true;
        }
        if (!result.rawResponse.pin_protection) {
          warnings.push(HwSecurityWarnings.NeedsPin);
          hasHwSecurityWarnings = true;
        }

        if (!lastestFirmwareVersionReaded) {
          warnings.push(HwSecurityWarnings.FirmwareVersionNotVerified);
        } else {
          if (!firmwareUpdated) {
            warnings.push(HwSecurityWarnings.OutdatedFirmware);
            hasHwSecurityWarnings = true;
          }
        }

        let walletNameUpdated = false;

        if (wallet) {
          const deviceLabel = result.rawResponse.label ? result.rawResponse.label : (result.rawResponse.deviceId ? result.rawResponse.deviceId : result.rawResponse.device_id);
          if (wallet.label !== deviceLabel) {
            wallet.label = deviceLabel;
            walletNameUpdated = true;
          }
          wallet.hasHwSecurityWarnings = hasHwSecurityWarnings;
          this.saveHardwareWallets();
        }

        const response = {
          features: result.rawResponse,
          securityWarnings: warnings,
          walletNameUpdated: walletNameUpdated,
        };

        return response;
      }));
    } else {
      return null;
    }
  }

  deleteHardwareWallet(wallet: Wallet): Observable<boolean> {
    if (wallet.isHardware) {
      return this.wallets.pipe(first(), map(wallets => {
        const index = wallets.indexOf(wallet);
        if (index !== -1) {
          wallets.splice(index, 1);

          this.saveHardwareWallets();
          this.refreshBalances();

          return true;
        }

        return false;
      }));
    }

    return null;
  }

  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').pipe(map(response => response.address));
  }

  outputs(): Observable<any> {
    return this.addressesAsString().pipe(
      first(),
      filter(addresses => !!addresses),
      mergeMap(addresses => this.apiService.post('outputs', {addrs: addresses})));
  }

  outputsWithWallets(): Observable<any> {
    return zip(this.all(), this.outputs(), (wallets, outputs) => {
      return wallets.map(wallet => {
        wallet.addresses = wallet.addresses.map(address => {
          address.outputs = outputs.head_outputs.filter(output => output.address === address.address);

          return address;
        });

        return wallet;
      });
    });
  }

  pendingTransactions(): Observable<PendingTransactions> {
    return this.pendingTxs.asObservable();
  }

  refreshBalances() {
    this.wallets.pipe(first()).subscribe(wallets => {
      observableForkJoin(wallets.map(wallet => this.retrieveWalletBalance(wallet).pipe(map(response => {
        wallet.coins = response.coins;
        wallet.hours = response.hours;
        wallet.addresses.map(address => {
          const balance = response.addresses.find(addr => addr.address === address.address);
          address.coins = balance ? balance.coins : new BigNumber(0);
          address.hours = balance ? balance.hours : new BigNumber(0);
        });

        return wallet;
      }))))
      .subscribe(newWallets => this.wallets.next(newWallets));
    });
  }

  renameWallet(wallet: Wallet, label: string): Observable<Wallet> {
    return this.apiService.post('wallet/update', { id: wallet.filename, label: label }).pipe(
      tap(() => {
        wallet.label = label;
        this.updateWallet(wallet);
      }));
  }

  toggleEncryption(wallet: Wallet, password: string): Observable<Wallet> {
    return this.apiService.postWalletToggleEncryption(wallet, password).pipe(
      tap(w => {
        wallet.encrypted = w.meta.encrypted;
        this.updateWallet(w);
      }));
  }

  resetPassword(wallet: Wallet, seed: string, password: string): Observable<Wallet> {
    const params = new Object();
    params['id'] = wallet.filename;
    params['seed'] = seed;
    if (password) {
      params['password'] = password;
    }

    return this.apiService.post('wallet/recover', params, {}, true).pipe(tap(w => {
      wallet.encrypted = w.data.meta.encrypted;
      this.updateWallet(w.data);
    }));
  }

  getWalletSeed(wallet: Wallet, password: string): Observable<string> {
    return this.apiService.getWalletSeed(wallet, password);
  }

  createTransaction(
    wallet: Wallet|null,
    addresses: string[]|null,
    unspents: string[]|null,
    destinations: any[],
    hoursSelection: any,
    changeAddress: string|null,
    password: string|null,
    unsigned: boolean): Observable<PreviewTransaction> {

    if (unspents) {
      addresses = null;
    }

    if (wallet && wallet.isHardware && !changeAddress) {
      changeAddress = wallet.addresses[0].address;
    }

    const useV2Endpoint = !wallet || !!wallet.isHardware;

    const params = {
      hours_selection: hoursSelection,
      wallet_id: wallet ? wallet.filename : null,
      password: password,
      addresses: addresses,
      unspents: unspents,
      to: destinations,
      change_address: changeAddress,
    };

    if (!useV2Endpoint) {
      params['unsigned'] = unsigned;
    }

    let response: Observable<PreviewTransaction> = this.apiService.post(
      useV2Endpoint ? 'transaction' : 'wallet/transaction',
      params,
      {
        json: true,
      },
      useV2Endpoint,
    ).pipe(map(transaction => {
      const data = useV2Endpoint ? transaction.data : transaction;

      if (wallet && wallet.isHardware) {
        if (data.transaction.inputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs-outputs'));
        }
        if (data.transaction.outputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs-outputs'));
        }
      }

      return {
        ...data.transaction,
        hoursBurned: new BigNumber(data.transaction.fee),
        encoded: data.encoded_transaction,
        innerHash: data.transaction.inner_hash,
      };
    }));

    if (wallet && wallet.isHardware && !unsigned) {
      let unsignedTx: PreviewTransaction;

      response = response.pipe(mergeMap(transaction => {
        unsignedTx = transaction;

        return this.signTransaction(wallet, null, transaction);
      })).pipe(map(signedTx => {
        unsignedTx.encoded = signedTx.encoded;

        return unsignedTx;
      }));
    }

    return response;
  }

  signTransaction(
    wallet: Wallet,
    password: string|null,
    transaction: PreviewTransaction,
    rawTransactionString = ''): Observable<PreviewTransaction> {

    if (!wallet.isHardware) {
      return this.apiService.post(
        'wallet/transaction/sign',
        {
          wallet_id: wallet ? wallet.filename : null,
          password: password,
          encoded_transaction: rawTransactionString ? rawTransactionString : transaction.encoded,
        },
        {
          json: true,
        },
        true,
      ).pipe(map(response => {
        return {
          ...response.data.transaction,
          hoursBurned: new BigNumber(response.data.transaction.fee),
          encoded: response.data.encoded_transaction,
        };
      }));

    } else {
      if (rawTransactionString) {
        throw new Error('Raw transactions not allowed.');
      }

      const txOutputs = [];
      const txInputs = [];
      const hwOutputs: HwOutput[] = [];
      const hwInputs: HwInput[] = [];

      transaction.outputs.forEach(output => {
        txOutputs.push({
          address: output.address,
          coins: parseInt(new BigNumber(output.coins).multipliedBy(1000000).toFixed(0), 10),
          hours: parseInt(output.hours, 10),
        });

        hwOutputs.push({
          address: output.address,
          coins: new BigNumber(output.coins).toString(),
          hours: new BigNumber(output.hours).toFixed(0),
        });
      });

      if (hwOutputs.length > 1) {
        for (let i = txOutputs.length - 1; i >= 0; i--) {
          if (hwOutputs[i].address === wallet.addresses[0].address) {
            hwOutputs[i].address_index = 0;
            break;
          }
        }
      }

      const addressesMap: Map<string, number> = new Map<string, number>();
      wallet.addresses.forEach((address, i) => addressesMap.set(address.address, i));

      transaction.inputs.forEach(input => {
        txInputs.push({
          hash: input.uxid,
          secret: '',
          address: input.address,
          address_index: addressesMap.get(input.address),
          calculated_hours: parseInt(input.calculated_hours, 10),
          coins: parseInt(input.coins, 10),
        });

        hwInputs.push({
          hash: input.uxid,
          index: addressesMap.get(input.address),
        });
      });

      return this.hwWalletService.signTransaction(hwInputs, hwOutputs).pipe(mergeMap(signatures => {
        const rawTransaction = TxEncoder.encode(
          hwInputs,
          hwOutputs,
          signatures.rawResponse,
          transaction.innerHash,
        );

        return of({
          ...transaction,
          encoded: rawTransaction,
        });
      }));
    }
  }

  injectTransaction(encodedTx: string, note: string): Observable<boolean> {
    return this.apiService.post('injectTransaction', { rawtx: encodedTx }, { json: true }).pipe(
      mergeMap(txId => {
        setTimeout(() => this.startDataRefreshSubscription(), 32);

        if (!note) {
          return of(false);
        } else {
          return this.storageService.store(StorageType.NOTES, txId, note).pipe(
            retryWhen(errors => errors.pipe(delay(1000), take(3), concat(observableThrowError(-1)))),
            catchError(err => err === -1 ? of(-1) : err),
            map(result => result === -1 ? false : true));
        }
      }));
  }

  transaction(txid: string): Observable<any> {
    return this.apiService.get('transaction', {txid: txid}).pipe(mergeMap(transaction => {
      if (transaction.txn.inputs && !transaction.txn.inputs.length) {
        return of(transaction);
      }

      return observableForkJoin(transaction.txn.inputs.map(input => this.retrieveInputAddress(input).pipe(map(response => {
        return response.owner_address;
      })))).pipe(map(inputs => {
        transaction.txn.inputs = inputs;

        return transaction;
      }));
    }));
  }

  transactions(): Observable<NormalTransaction[]> {
    let wallets: Wallet[];
    let transactions: NormalTransaction[];
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();


    return this.wallets.pipe(first(), mergeMap(w => {
      wallets = w;

      return this.allAddresses().pipe(first());
    }), mergeMap(addresses => {
      this.addresses = addresses;
      addresses.map(add => addressesMap.set(add.address, true));

      return this.apiService.getTransactions(addresses);
    }), mergeMap(recoveredTransactions => {
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

  startDataRefreshSubscription() {
    if (this.dataRefreshSubscription) {
      this.dataRefreshSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.dataRefreshSubscription = timer(0, 10000)
        .subscribe(() => this.ngZone.run(() => {
          this.refreshBalances();
          this.refreshPendingTransactions();
        }));
    });
  }

  saveHardwareWallets() {
    this.wallets.pipe(first()).subscribe(wallets => {
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

      this.hwWalletService.saveWalletsData(JSON.stringify(hardwareWallets));

      this.wallets.next(wallets);
    });
  }

  verifyAddress(address: string) {
    return this.apiService.post('address/verify', { address }, {}, true)
      .pipe(map(() => true), catchError(() => of(false)));
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
    let wallets: Wallet[] = [];
    let softwareWallets: Wallet[] = [];

    this.apiService.getWallets().pipe(first(), mergeMap(recoveredWallets => {
      softwareWallets = recoveredWallets;

      if (this.hwWalletService.hwWalletCompatibilityActivated) {
        return this.loadHardwareWallets(wallets);
      }

      return of(null);

    })).subscribe(() => {
      wallets = wallets.concat(softwareWallets);
      this.wallets.next(wallets);
    }, () => this.initialLoadFailed.next(true));
  }

  private loadHardwareWallets(wallets: Wallet[]): Observable<any> {
    return this.hwWalletService.getSavedWalletsData().pipe(map(storedWallets => {
      if (storedWallets) {
        const loadedWallets: Wallet[] = JSON.parse(storedWallets);
        loadedWallets.map(wallet => wallets.push(wallet));
      }

      return null;
    }));
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

    return query.pipe(map(balance => {
      return {
        coins: new BigNumber(balance.confirmed.coins).dividedBy(1000000),
        hours: new BigNumber(balance.confirmed.hours),
        addresses: Object.keys(balance.addresses).map(address => ({
          address,
          coins: new BigNumber(balance.addresses[address].confirmed.coins).dividedBy(1000000),
          hours: new BigNumber(balance.addresses[address].confirmed.hours),
        })),
      };
    }));
  }

  private updateWallet(wallet: Wallet) {
    this.wallets.pipe(first()).subscribe(wallets => {
      const index = wallets.findIndex(w => w.filename === wallet.filename);
      wallets[index] = wallet;
      this.wallets.next(wallets);
    });
  }

  private refreshPendingTransactions() {
    this.apiService.get('pendingTxs', { verbose: true }).pipe(
      mergeMap((transactions: any) => {
        if (transactions.length === 0) {
          return of({
            user: [],
            all: [],
          });
        }

        return this.wallets.pipe(first(), map((wallets: Wallet[]) => {
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
      }))
      .subscribe(transactions => this.pendingTxs.next(transactions));
  }

  getOutputs(addresses): Observable<Output[]> {
    if (!addresses) {
      return of([]);
    } else {
      return this.apiService.post('outputs', { addrs: addresses }).pipe(map((response) => {
        const outputs = [];
        response.head_outputs.forEach(output => outputs.push({
          address: output.address,
          coins: new BigNumber(output.coins),
          hash: output.hash,
          calculated_hours: new BigNumber(output.calculated_hours),
        }));

        return outputs;
      }));
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
