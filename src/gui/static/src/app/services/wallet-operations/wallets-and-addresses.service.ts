import { of, Subject, Observable, ReplaySubject, BehaviorSubject, throwError as observableThrowError, Subscription, throwError } from 'rxjs';
import { mergeMap, map, catchError } from 'rxjs/operators';
import { TranslateService } from '@ngx-translate/core';
import { Injectable } from '@angular/core';

import { ApiService } from '../api.service';
import { HwWalletService } from '../hw-wallet.service';
import { AppConfig } from '../../app.config';
import { WalletBase, AddressBase, duplicateWalletBase } from './wallet-objects';
import { processServiceError, redirectToErrorPage } from '../../utils/errors';
import { StorageService, StorageType } from '../storage.service';
import { OperationError, OperationErrorTypes } from '../../utils/operation-error';

/**
 * Manages the list with the wallets and its addresses. It works like a CRUD for the wallet list, so
 * it does not contain functions for specific things, like changing the label of a wallet.
 */
@Injectable()
export class WalletsAndAddressesService {
  /**
   * Key used for saving the hw wallet list in persistent storage.
   */
  private readonly hwWalletsDataStorageKey = 'hw-wallets';

  // Wallet list and the subject used for informing when the list has been modified.
  private walletsList: WalletBase[];
  private walletsSubject: Subject<WalletBase[]> = new ReplaySubject<WalletBase[]>(1);
  // Indicates if the initial load failed.
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

  /**
   * Allows to know if there was an error while trying to load the wallets.
   */
  get errorDuringinitialLoad(): Observable<boolean> {
    return this.initialLoadFailed.asObservable();
  }

  /**
   * Gets the wallet list. It emits every time the wallet list is updated. Please note that if
   * any value of the returned wallets is modified, the changes must be notified by calling the
   * informValuesUpdated function or the behavior will be indeterminate.
   */
  get allWallets(): Observable<WalletBase[]> {
    return this.walletsSubject.asObservable();
  }

  /**
   * Adds one or more addresses to a wallet.
   * @param wallet Wallet to which the addresses will be added.
   * @param num Number of addresses to create.
   * @param password Wallet password, if the wallet is encrypted.
   * @returns An array with the newly created addresses.
   */
  addAddressesToWallet(wallet: WalletBase, num: number, password?: string): Observable<AddressBase[]> {
    if (!wallet.isHardware) {
      const params = new Object();
      params['id'] = wallet.id;
      params['num'] = num;
      if (password) {
        params['password'] = password;
      }

      // Add the addresses on the backend.
      return this.apiService.post('wallet/newAddress', params).pipe(map((response: any) => {
        // Find the affected wallet on the local list and add the addresses to it.
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
      // Generate the new addresses on the device.
      return this.hwWalletService.getAddresses(num, wallet.addresses.length).pipe(map(response => {
        // Find the affected wallet on the local list and add the addresses to it.
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

  /**
   * Scans the addreses of a wallet, to find if there is an addreeses with transactions which is
   * not on the addresses list. If that happens, the last address with at least one transaction
   * and all the addresses that precede it in the deterministic generation order are added to
   * the wallet. NOTE: it does not work with hardware wallets.
   * @param wallet Wallet to scan.
   * @param password Wallet password, if the wallet is encrypted.
   * @returns true if new addresses were added to the wallet, false otherwise.
   */
  scanAddresses(wallet: WalletBase, password?: string): Observable<boolean> {
    if (!wallet.isHardware) {
      const params = new Object();
      params['id'] = wallet.id;
      if (password) {
        params['password'] = password;
      }

      // Request the backend to scan the addresses.
      return this.apiService.post('wallet/scan', params).pipe(map((response: any) => {
        // Find the affected wallet on the local list and add the addresses to it.
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

  /**
   * This function must be called when any value of a wallet is changed, to ensure the wallet
   * list is updated and inform all the subscribers of the wallet list that there was a change.
   * @param wallet Object with all the properties of the wallet. Its ID must coincide with the
   * ID of one of the wallets of the wallet list or nothing will happen. Note that this object
   * is not directly saved on the wallet list, so you must always call this function after
   * making changes to a wallet.
   */
  informValuesUpdated(wallet: WalletBase) {
    const affectedWalletIndex = this.walletsList.findIndex(w => w.id === wallet.id);
    if (affectedWalletIndex === -1) {
      return;
    }

    // Create a duplicate of the provided wallet and save it on the wallet list.
    const newWallet = duplicateWalletBase(wallet, true);
    this.walletsList[affectedWalletIndex] = newWallet;

    // Save if needed and inform the changes.
    if (wallet.isHardware) {
      this.saveHardwareWalletsAndInformUpdate();
    } else {
      this.informDataUpdated();
    }
  }

  /**
   * Adds a new wallet to the node and adds it to the wallets list.
   * @param label Name given by the user to the wallet.
   * @param seed Wallet seed.
   * @param password Wallet password, if it will be encrypted, null otherwise.
   * @returns The returned observable returns nothing, but it can fail in case of error.
   */
  createSoftwareWallet(label: string, seed: string, password?: string): Observable<void> {
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

    // Ask the node to create the wallet and return the data of the newly created wallet.
    return this.apiService.post('wallet/create', params).pipe(map(response => {
      const wallet: WalletBase = {
        label: response.meta.label,
        // The filename is saved as the id of the wallet.
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

  /**
   * Adds a new hardware wallet to the wallets list, with the data of the currently connected device.
   * @returns The newly creatd wallet.
   */
  createHardwareWallet(): Observable<WalletBase> {
    let addresses: string[];
    let lastAddressWithTx = 0;
    const addressesMap: Map<string, boolean> = new Map<string, boolean>();
    const addressesWithTxMap: Map<string, boolean> = new Map<string, boolean>();

    // Ask the device to return the as many addresses as set on AppConfig.maxHardwareWalletAddresses.
    return this.hwWalletService.getAddresses(AppConfig.maxHardwareWalletAddresses, 0).pipe(mergeMap(response => {
      // Save all addresses in a map.
      addresses = response.rawResponse;
      addresses.forEach(address => {
        addressesMap.set(address, true);
      });

      // Throw an error if any wallet has any of the addresses as ID, as hw wallets use the first
      // address of the device as ID.
      let walletAlreadyExists = false;
      this.walletsList.forEach(wallet => {
        if (addressesMap.has(wallet.id)) {
          walletAlreadyExists = true;
        }
      });
      if (walletAlreadyExists) {
        return observableThrowError(processServiceError('The wallet already exists'));
      }

      // Request the transaction history of all addresses.
      const addressesString = addresses.join(',');

      return this.apiService.post('transactions', { addrs: addressesString });
    }), map(response => {
      // Get the index of the last address of the list with transaction.
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

      // Use the first address as ID.
      newWallet.id = newWallet.addresses[0].address;

      // Add the wallet just after the las hw wallet of the wallet list.
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

  /**
   * Removes a hw wallet from the wallet list.
   * @param walletId Id of the wallet to be removed. If the ID is not on the wallet list or
   * is not for a hw wallet, nothing happens.
   */
  deleteHardwareWallet(walletId: string) {
    const index = this.walletsList.findIndex(w => w.id === walletId);
    if (index === -1 || !this.walletsList[index].isHardware) {
      return;
    }

    this.walletsList.splice(index, 1);
    this.saveHardwareWalletsAndInformUpdate();
  }

  /**
   * Saves on persistent storage the data of all the hw wallets on the wallet list. It overwrites
   * the previously saved data. I also calls informDataUpdated().
   */
  private saveHardwareWalletsAndInformUpdate() {
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

    // Cancel any previous saving operation.
    if (this.savingHwWalletDataSubscription) {
      this.savingHwWalletDataSubscription.unsubscribe();
    }

    // The data is saved as a JSON string.
    this.savingHwWalletDataSubscription =
      this.storageService.store(StorageType.CLIENT, this.hwWalletsDataStorageKey, JSON.stringify(hardwareWallets))
        .subscribe({
          next: null,
          error: () => redirectToErrorPage(3),
        });

    this.informDataUpdated();
  }

  /**
   * Helper function for creating a WalletBase object for a hw wallet.
   */
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

  /**
   * Gets the saved wallets data and populates de wallet list with it.
   */
  private loadWallets(): void {
    // Get all the software wallets managed by the node.
    const softwareWallets: WalletBase[] = [];
    this.apiService.get('wallets').pipe(mergeMap((response: any[]) => {
      // Save the software wallets list.
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

      // Get the hardware wallets.
      if (this.hwWalletService.hwWalletCompatibilityActivated) {
        return this.loadHardwareWallets();
      }

      return of([]);
    })).subscribe((hardwareWallets: WalletBase[]) => {
      // Hardware wallets are put first on the wallet list.
      this.walletsList = hardwareWallets.concat(softwareWallets);

      this.informDataUpdated();
    }, () => this.initialLoadFailed.next(true));
  }

  /**
   * Load all the hw wallets saved on persistent storage.
   * @returns The list of hw wallets.
   */
  private loadHardwareWallets(): Observable<WalletBase[]> {

    return this.storageService.get(StorageType.CLIENT, this.hwWalletsDataStorageKey).pipe(
      map(result => result.data),
      catchError((err: OperationError) => {
        // If the node returned a 404 error, it means there are no hw wallets saved, not
        // that there was an error during the operations, so we can continue.
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

          // Prepare to remove all unexpected properties, which could have been saved in a
          // previous version of the app.
          const knownPropertiesMap = new Map<string, boolean>();
          const referenceObject = new WalletBase();
          Object.keys(referenceObject).forEach(property => {
            knownPropertiesMap.set(property, true);
          });

          loadedWallets.forEach(wallet => {
            // Remove all unexpected properties.
            const propertiesToRemove: string[] = [];
            Object.keys(wallet).forEach(property => {
              if (!knownPropertiesMap.has(property)) {
                propertiesToRemove.push(property);
              }
            });
            propertiesToRemove.forEach(property => {
              delete wallet[property];
            });

            // The wallet must be identified as a hw wallet and have at least one address.
            // This is just a precaution.
            wallet.isHardware = true;
            if (!wallet.addresses) {
              wallet.addresses = [{ address: 'invalid', confirmed: false, }];
            }

            wallet.id = wallet.addresses[0].address;
          });

          return loadedWallets;
        }

        return [];
      }),
    );
  }

  /**
   * Gets the path of the folder were the node saves the data of the software wallets.
   */
  folder(): Observable<string> {
    return this.apiService.get('wallets/folderName').pipe(map(response => response.address));
  }

  /**
   * Checks if a string is a valid address.
   * @param address String to check.
   * @returns True if the address is valid or false otherwise.
   */
  verifyAddress(address: string): Observable<boolean> {
    return this.apiService.post('address/verify', { address }, {useV2: true}).pipe(
      map(() => true),
      catchError((err: OperationError) => {
        err = processServiceError(err);

        // Return false in case of error, but not if the error was for a connection problem.
        if (err.type !== OperationErrorTypes.NoInternet) {
          return of(false);
        } else {
          return throwError(err);
        }
      }),
    );
  }

  /**
   * Makes walletsSubject emit, to inform that the wallet list has been updated.
   */
  private informDataUpdated() {
    this.walletsSubject.next(this.walletsList);
  }
}
