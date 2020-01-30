import { of, Subject, Observable, ReplaySubject, BehaviorSubject, throwError as observableThrowError, Subscription } from 'rxjs';
import { tap, mergeMap, map, catchError } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { ApiService } from '../api.service';
import { HwWalletService } from '../hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../../app.config';
import { WalletBase, AddressBase } from './wallet-objects';
import { processServiceError } from '../../utils/errors';
import { StorageService, StorageType } from '../storage.service';
import { OperationError } from '../../utils/operation-error';

@Injectable()
export class WalletsAndAddressesService {
  private readonly hwWalletsDataStorageKey = 'hw-wallets';

  private walletsList: WalletBase[];
  private walletsSubject: Subject<WalletBase[]> = new ReplaySubject<WalletBase[]>(1);
  private initialLoadFailed: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  private savingHwWalletDataSubscription: Subscription;

  constructor(
    private apiService: ApiService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private storageService: StorageService,
  ) {
    this.loadWallets();
  }

  get errorDuringinitialLoad(): Observable<boolean> {
    return this.initialLoadFailed.asObservable();
  }

  get allWallets(): Observable<WalletBase[]> {
    return this.walletsSubject.asObservable();
  }

  get allAddresses(): Observable<AddressBase[]> {
    return this.allWallets.pipe(map(wallets => wallets.reduce((array, wallet) => array.concat(wallet.addresses), [] as AddressBase[])));
  }

  get addressesAsString(): Observable<string> {
    return this.allAddresses.pipe(map(addrs => addrs.map(addr => addr.address)), map(addrs => addrs.join(',')));
  }

  addAddressesToWallet(wallet: WalletBase, num: number, password?: string): Observable<AddressBase[]> {
    if (!wallet.isHardware) {
      const params = new Object();
      params['id'] = wallet.id;
      params['num'] = num;
      if (password) {
        params['password'] = password;
      }

      return this.apiService.post('wallet/newAddress', params).pipe(map((response: any) => {
        const affectedWallet = this.walletsList.find(w => w.id === wallet.id);
        const newAddresses: AddressBase[] = [];
        (response.addresses as any[]).forEach(value => {
          const newAddress: AddressBase = {address: value, confirmed: true};
          newAddresses.push(newAddress);
          affectedWallet.addresses.push(newAddress);
        });
        this.informDataUpdated();

        return newAddresses;
      }));
    } else {
      return this.hwWalletService.getAddresses(num, wallet.addresses.length).pipe(map(response => {
        const affectedWallet = this.walletsList.find(w => w.id === wallet.id);
        const newAddresses: AddressBase[] = [];
        (response.rawResponse as any[]).forEach(value => {
          const newAddress: AddressBase = {address: value, confirmed: false};
          newAddresses.push(newAddress);
          affectedWallet.addresses.push(newAddress);
        });
        this.saveHardwareWalletsAndInformUpdate();

        return newAddresses;
      }));
    }
  }

  scanAddresses(wallet: WalletBase, password?: string): Observable<boolean> {
    if (!wallet.isHardware) {
      const params = new Object();
      params['id'] = wallet.id;
      if (password) {
        params['password'] = password;
      }

      return this.apiService.post('wallet/scan', params).pipe(map((response: any) => {
        const affectedWallet = this.walletsList.find(w => w.id === wallet.id);
        const newAddresses: string[] = response.addresses;
        if (newAddresses && newAddresses.length > 0) {
          newAddresses.forEach(address => {
            affectedWallet.addresses.push({address: address, confirmed: true});
          });
          this.informDataUpdated();

          return true;
        } else {
          return false;
        }
      }));
    } else {
      // Not implemented.
      return of(false);
    }
  }

  informValuesUpdated(wallet: WalletBase) {
    const affectedWallet = this.walletsList.find(w => w.id === wallet.id);
    const referenceWallet = new WalletBase();
    Object.keys(referenceWallet).forEach(property => {
      if (property !== 'addresses') {
        affectedWallet[property] = wallet[property];
      }
    });

    if (affectedWallet.addresses.length !== wallet.addresses.length) {
      affectedWallet.addresses = [];
      for (let i = 0; i < wallet.addresses.length; i++) {
        affectedWallet.addresses.push(new AddressBase());
      }
    }

    const referenceAddress = new AddressBase();
    wallet.addresses.forEach((address, i) => {
      Object.keys(referenceAddress).forEach(property => {
        affectedWallet.addresses[i][property] = address[property];
      });
    });

    if (wallet.isHardware) {
      this.saveHardwareWalletsAndInformUpdate();
    } else {
      this.informDataUpdated();
    }
  }

  createSoftwareWallet(label: string, seed: string, password: string): Observable<void> {
    seed = seed.replace(/(\n|\r\n)$/, '');

    const params = {
      label: label ? label : 'undefined',
      seed: seed,
      scan: 100,
      type: 'deterministic',
    };

    if (password) {
      params['password'] = password;
      params['encrypt'] = true;
    }

    return this.apiService.post('wallet/create', params).pipe(tap(response => {
      const wallet: WalletBase = {
        label: response.meta.label,
        id: response.meta.filename,
        addresses: [],
        encrypted: response.meta.encrypted,
        isHardware: false,
        hasHwSecurityWarnings: false,
        stopShowingHwSecurityPopup: true,
      };

      (response.entries as any[]).forEach(entry => wallet.addresses.push({address: entry.address, confirmed: true}));

      this.walletsList.push(wallet);

      this.informDataUpdated();
    }));
  }

  createHardwareWallet(): Observable<WalletBase> {
    let addresses: string[];
    let lastAddressWithTx = 0;
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();
    const addressesWithTxMap: Map<string, boolean> = new Map<string, boolean>();

    return this.hwWalletService.getAddresses(AppConfig.maxHardwareWalletAddresses, 0).pipe(mergeMap(response => {
      addresses = response.rawResponse;
      addresses.forEach(address => {
        addressesMap.set(address, true);
      });

      let walletAlreadyExists = false;
      this.walletsList.forEach(wallet => {
        if (addressesMap.has(wallet.id)) {
          walletAlreadyExists = true;
        }
      });
      if (walletAlreadyExists) {
        return observableThrowError(processServiceError('The wallet already exists'));
      }

      const addressesString = addresses.join(',');

      return this.apiService.post('transactions', { addrs: addressesString });
    }), map(response => {
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

      const newWallet = this.createHardwareWalletData(
        this.translate.instant('hardware-wallet.general.default-wallet-name'),
        addresses.slice(0, lastAddressWithTx + 1).map(add => {
          return { address: add, confirmed: false };
        }), true, false,
      );

      newWallet.id = newWallet.addresses[0].address;

      let lastHardwareWalletIndex = this.walletsList.length - 1;
      for (let i = 0; i < this.walletsList.length; i++) {
        if (!this.walletsList[i].isHardware) {
          lastHardwareWalletIndex = i - 1;
          break;
        }
      }
      this.walletsList.splice(lastHardwareWalletIndex + 1, 0, newWallet);
      this.saveHardwareWalletsAndInformUpdate();

      return newWallet;
    }));
  }

  deleteHardwareWallet(wallet: WalletBase): boolean {
    if (wallet.isHardware) {
      const index = this.walletsList.findIndex(w => {
        return w.id === wallet.id;
      });

      if (index !== -1) {
        this.walletsList.splice(index, 1);

        this.saveHardwareWalletsAndInformUpdate();

        return true;
      }

      return false;
    }

    return null;
  }

  saveHardwareWalletsAndInformUpdate() {
    const hardwareWallets: WalletBase[] = [];

    this.walletsList.map(wallet => {
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

    if (this.savingHwWalletDataSubscription) {
      this.savingHwWalletDataSubscription.unsubscribe();
    }

    this.savingHwWalletDataSubscription =
      this.storageService.store(StorageType.CLIENT, this.hwWalletsDataStorageKey, JSON.stringify(hardwareWallets)).subscribe();

    this.informDataUpdated();
  }

  private createHardwareWalletData(label: string, addresses: AddressBase[], hasHwSecurityWarnings: boolean, stopShowingHwSecurityPopup: boolean): WalletBase {
    return {
      label: label,
      id: '',
      hasHwSecurityWarnings: hasHwSecurityWarnings,
      stopShowingHwSecurityPopup: stopShowingHwSecurityPopup,
      addresses: addresses,
      encrypted: false,
      isHardware: true,
    };
  }

  private loadWallets(): void {
    const softwareWallets: WalletBase[] = [];
    this.apiService.get('wallets').pipe(mergeMap((response: any[]) => {
      response.forEach(wallet => {
        const processedWallet: WalletBase = {
          label: wallet.meta.label,
          id: wallet.meta.filename,
          addresses: [],
          encrypted: wallet.meta.encrypted,
          isHardware: false,
          hasHwSecurityWarnings: false,
          stopShowingHwSecurityPopup: true,
        };

        if (wallet.entries) {
          processedWallet.addresses = (wallet.entries as any[]).map<AddressBase>((entry: any) => {
            return {
              address: entry.address,
              confirmed: true,
            };
          });
        }

        softwareWallets.push(processedWallet);
      });

      if (this.hwWalletService.hwWalletCompatibilityActivated) {
        return this.loadHardwareWallets();
      }

      return of(null);
    })).subscribe((hardwareWallets: WalletBase[]) => {
      if (hardwareWallets) {
        this.walletsList = hardwareWallets.concat(softwareWallets);
      } else {
        this.walletsList = softwareWallets;
      }

      this.informDataUpdated();
    }, () => this.initialLoadFailed.next(true));
  }

  private loadHardwareWallets(): Observable<WalletBase[]> {

    return this.storageService.get(StorageType.CLIENT, this.hwWalletsDataStorageKey).pipe(
      map(result => result.data),
      catchError((err: OperationError) => {
        err = processServiceError(err);
        try {
          if (err && err.originalError &&  err.originalError.status && err.originalError.status === 404) {
            return of(null);
          }
        } catch (e) {}

        return observableThrowError(err);
      }),
      map(storedWallets => {
        if (storedWallets) {
          const loadedWallets: WalletBase[] = JSON.parse(storedWallets);

          const knownPropertiesMap = new Map<string, boolean>();
          const referenceObject = new WalletBase();
          Object.keys(referenceObject).forEach(property => {
            knownPropertiesMap.set(property, true);
          });

          loadedWallets.forEach(wallet => {
            const propertiesToRemove: string[] = [];
            Object.keys(wallet).forEach(property => {
              if (!knownPropertiesMap.has(property)) {
                propertiesToRemove.push(property);
              }
            });

            propertiesToRemove.forEach(property => {
              delete wallet[property];
            });

            wallet.isHardware = true;

            if (!wallet.addresses) {
              wallet.addresses = [{ address: 'invalid', confirmed: false, }];
            }

            wallet.id = wallet.addresses[0].address;
          });

          return loadedWallets;
        }

        return null;
      }),
    );
  }

  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').pipe(map(response => response.address));
  }

  verifyAddress(address: string) {
    return this.apiService.post('address/verify', { address }, {}, true)
      .pipe(map(() => true), catchError(() => of(false)));
  }

  private informDataUpdated() {
    this.walletsSubject.next(this.walletsList);
  }
}
