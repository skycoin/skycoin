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
import { catchError, map } from 'rxjs/operators';

declare var Cipher: any;
declare var CipherExtras: any;

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

    return this.hwWalletService.getMaxAddresses().flatMap(response => {
      addresses = response;
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
        const newWallet = this.createHardwareWalletData(this.translate.instant('hardware-wallet.general.default-wallet-name'), addresses.slice(0, lastAddressWithTx + 1), true, false);
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
        wallet.addresses = wallet.addresses.map(address => {
          return response.addresses.find(addr => addr.address === address.address);
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

  createTransaction(wallet: Wallet, addresses: string[]|null, destinations: any[], hoursSelection: any, changeAddress: string|null, password: string|null): Observable<PreviewTransaction> {
    return this.apiService.post(
      'wallet/transaction',
      {
        hours_selection: hoursSelection,
        wallet: {
          id: wallet.filename,
          password,
          addresses,
        },
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
   let unburnedHoursRatio: BigNumber;

    const addresses = wallet.addresses.map(a => a.address).join(',');

    let totalHours = new BigNumber('0');
    let hoursToSend = new BigNumber('0');
    let calculatedHours = new BigNumber('0');

    let innerHash: string;
    let convertedOutputs: any[];

    const txOutputs = [];
    const txInputs = [];
    const txSignatures = [];

    return this.apiService.get('health').flatMap(response => {
        unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(response.user_verify_transaction.burn_factor));

        return this.getOutputs(addresses);
      }).flatMap((outputs: Output[]) => {
        const minRequiredOutputs =  this.getMinRequiredOutputs(amount, outputs);
        let totalCoins = new BigNumber('0');
        minRequiredOutputs.map(output => totalCoins = totalCoins.plus(output.coins));

        if (totalCoins.isLessThan(amount)) {
          throw new Error(this.translate.instant('service.wallet.not-enough-hours'));
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
        } else {
          hoursToSend = calculatedHours;
        }

        txOutputs.push({ address: address, coins: amount.toNumber(), hours: hoursToSend.toNumber() });

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
        });

        convertedOutputs = txOutputs.map(output => {
          return {
            ...output,
            coins: parseInt((output.coins * 1000000) + '', 10),
          };
        });

        innerHash = Cipher.GetTransactionInnerHash(JSON.stringify(txInputs), JSON.stringify(convertedOutputs));

        // Request signatures and add them to txSignatures sequentially.
        return this.addSignatures(txInputs.length - 1, txInputs, txSignatures, innerHash);

      }).flatMap(() => {
        const rawTransaction = Cipher.PrepareTransactionWithSignatures(
          JSON.stringify(txInputs),
          JSON.stringify(convertedOutputs),
          JSON.stringify(txSignatures),
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
    return this.allAddresses().first().flatMap(addresses => {
      this.addresses = addresses;

      return this.apiService.getTransactions(addresses);
    }).map(transactions => {
      return transactions
        .sort((a, b) =>  b.timestamp - a.timestamp)
        .map(transaction => {
          const outgoing = this.addresses.some(address => {
            return transaction.inputs.some(input => input.owner === address.address);
          });

          const relevantOutputs = transaction.outputs.reduce((array, output) => {
            const isMyOutput = this.addresses.some(address => address.address === output.dst);

            if ((outgoing && !isMyOutput) || (!outgoing && isMyOutput)) {
              array.push(output);
            }

            return array;
          }, []);

          const calculatedOutputs = (outgoing && relevantOutputs.length === 0)
          || (!outgoing && relevantOutputs.length === transaction.outputs.length)
            ? transaction.outputs
            : relevantOutputs;

          transaction.addresses.push(
            ...calculatedOutputs
              .map(output => output.dst)
              .filter((dst, i, self) => self.indexOf(dst) === i),
          );

          calculatedOutputs.map (output => transaction.balance = transaction.balance.plus(output.coins));
          transaction.balance = (outgoing ? transaction.balance.negated() : transaction.balance);

          transaction.hoursSent = new BigNumber('0');
          calculatedOutputs.map(output => transaction.hoursSent = transaction.hoursSent.plus(new BigNumber(output.hours)));

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
            wallet.addresses.map(address => address.address),
            wallet.hasHwSecurityWarnings,
            wallet.stopShowingHwSecurityPopup,
          ));
        }
      });

      this.hwWalletService.saveWalletsDataSync(JSON.stringify(hardwareWallets));

      this.wallets.next(wallets);
    });
  }

  verifyAddress(address: string) {
    return this.apiService.post('address/verify', { address }, {}, true)
      .pipe(map(() => true), catchError(() => Observable.of(false)));
  }

  private addSignatures(index: number, txInputs: any[], txSignatures: string[], txInnerHash: string): Observable<any> {
    let chain: Observable<any>;
    if (index > 0) {
      chain = this.addSignatures(index - 1, txInputs, txSignatures, txInnerHash).first();
    } else {
      chain = Observable.of(1);
    }

    chain = chain.flatMap(() => {
      const msgToSign = CipherExtras.AddSHA256(txInnerHash, txInputs[index].hash);

      return this.hwWalletService.signMessage(txInputs[index].address_index, msgToSign, index + 1, txInputs.length)
        .map(response => {
          txSignatures.push(response.rawResponse);
        });
    });

    return chain;
  }

  private createHardwareWalletData(label: string, addresses: string[], hasHwSecurityWarnings: boolean, stopShowingHwSecurityPopup: boolean): Wallet {
    return {
      label: label,
      filename: '',
      hasHwSecurityWarnings: hasHwSecurityWarnings,
      stopShowingHwSecurityPopup: stopShowingHwSecurityPopup,
      coins: null,
      hours: null,
      addresses: addresses.map(address => {
        return {
          address: address,
          coins: null,
          hours: null,
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
        if (window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated')) {
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

  private getMinRequiredOutputs(transactionAmount: BigNumber, outputs: Output[]): Output[] {
    outputs.sort( function(a, b) {
      return b.coins.minus(a.coins).toNumber();
    });

    const minRequiredOutputs: Output[] = [];
    let sumCoins: BigNumber = new BigNumber('0');

    outputs.forEach(output => {
      if (sumCoins.isLessThan(transactionAmount) && output.calculated_hours.isGreaterThan(0)) {
        minRequiredOutputs.push(output);
        sumCoins = sumCoins.plus(output.coins);
      }
    });

    return minRequiredOutputs;
  }
}
