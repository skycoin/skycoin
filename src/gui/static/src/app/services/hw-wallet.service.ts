// NOTE: some  code for using the hw wallet js library was left here only for precaution and should be deleted soon.

import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscriber } from 'rxjs/Subscriber';
import { Subject } from 'rxjs/Subject';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { HwWalletDaemonService } from './hw-wallet-daemon.service';
import { HwWalletPinService, ChangePinStates } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import BigNumber from 'bignumber.js';
import { StorageService, StorageType } from './storage.service';
import { ISubscription } from 'rxjs/Subscription';
import { HttpClient } from '@angular/common/http';
import { ApiService } from './api.service';

export enum OperationResults {
  Success,
  FailedOrRefused,
  PinMismatch,
  WithoutSeed,
  WrongPin,
  IncorrectHardwareWallet,
  WrongWord,
  InvalidSeed,
  WrongSeed,
  UndefinedError,
  Disconnected,
  DaemonError,
  InvalidAddress,
  Timeout,
  NotInBootloaderMode,
}

export class TxData {
  address: string;
  coins: BigNumber;
  hours: BigNumber;
}

export class OperationResult {
  result: OperationResults;
  rawResponse: any;
}

export interface Input {
  hashIn: string;
  index: number;
}

export interface Output {
  address: string;
  coin: number;
  hour: number;
  address_index?: number;
}

interface EventData {
  event: string;
  successTexts?: string[];
}

@Injectable()
export class HwWalletService {

  public static readonly maxLabelLength = 32;

  private readonly storageKey = 'hw-wallets';

  showOptionsWhenPossible = false;

  private requestSequence = 0;

  private eventsObservers = new Map<number, Subscriber<OperationResult>>();
  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();

  private savingDataSubscription: ISubscription;

  private signTransactionDialog: MatDialogRef<{}, any>;

  // Set on AppComponent to avoid a circular reference.
  private signTransactionConfirmationComponentInternal;
  set signTransactionConfirmationComponent(value) {
    this.signTransactionConfirmationComponentInternal = value;
  }

  constructor(
    private translate: TranslateService,
    private dialog: MatDialog,
    private hwWalletDaemonService: HwWalletDaemonService,
    private hwWalletPinService: HwWalletPinService,
    private hwWalletSeedWordService: HwWalletSeedWordService,
    private storageService: StorageService,
    private apiService: ApiService,
    private http: HttpClient) {

    if (this.hwWalletCompatibilityActivated) {
      if (!AppConfig.useHwWalletDaemon) {
        window['ipcRenderer'].on('hwConnectionEvent', (event, connected) => {
          if (!connected) {
            this.eventsObservers.forEach((value, key) => {
              this.dispatchError(key, OperationResults.Disconnected, this.translate.instant('hardware-wallet.general.error-disconnected'));
            });
          }
          this.walletConnectedSubject.next(connected);
        });
      } else {
        hwWalletDaemonService.connectionEvent.subscribe(connected => {
          this.walletConnectedSubject.next(connected);
        });
      }

      if (!AppConfig.useHwWalletDaemon) {
        window['ipcRenderer'].on('hwPinRequested', (event) => {
          this.hwWalletPinService.requestPin().subscribe(pin => {
            if (!pin) {
              this.cancelAllOperations();
              window['ipcRenderer'].send('hwCancelPin');
            } else {
              window['ipcRenderer'].send('hwSendPin', pin);
            }
          });
        });
        window['ipcRenderer'].on('hwSeedWordRequested', (event) => {
          this.hwWalletSeedWordService.requestWord().subscribe(word => {
            if (!word) {
              this.cancelAllOperations();
              window['ipcRenderer'].send('hwCancelLastAction');
            }
            window['ipcRenderer'].send('hwSendSeedWord', word);
          });
        });

        window['ipcRenderer'].on('hwSignTransactionResponse', (event, requestId, result) => {
          this.closeTransactionDialog();
          this.dispatchEvent(requestId, result, true);
        });

        const data: EventData[] = [
          { event: 'hwChangePinResponse', successTexts: ['PIN changed'] },
          { event: 'hwGenerateMnemonicResponse', successTexts: ['operation completed'] },
          { event: 'hwRecoverMnemonicResponse', successTexts: ['Device recovered', 'The seed is valid and matches the one in the device'] },
          { event: 'hwBackupDeviceResponse', successTexts: ['operation completed'] },
          { event: 'hwWipeResponse', successTexts: ['operation completed'] },
          { event: 'hwChangeLabelResponse', successTexts: ['Settings applied'] },
          { event: 'hwCancelLastActionResponse' },
          { event: 'hwGetAddressesResponse' },
          { event: 'hwGetFeaturesResponse' },
          { event: 'hwSignMessageResponse' },
        ];

        data.forEach(item => {
          window['ipcRenderer'].on(item.event, (event, requestId, result) => {
            const success = item.successTexts
              ? typeof result === 'string' && item.successTexts.some(text => (result as string).includes(text))
              : true;

            this.dispatchEvent(requestId, result, success);
          });
        });
      }
    }
  }

  get hwWalletCompatibilityActivated(): boolean {
    if (!AppConfig.useHwWalletDaemon) {
      return window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated');
    } else {
      return true;
    }
  }

  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  getDeviceConnected(): Observable<boolean> {
    if (!AppConfig.useHwWalletDaemon) {
      return Observable.of(window['ipcRenderer'].sendSync('hwGetDeviceConnectedSync'));
    } else {
      return this.hwWalletDaemonService.get('/available').map(response => {
        return response.data;
      });
    }
  }

  getSavedWalletsData(): Observable<string> {
    return this.storageService.get(StorageType.CLIENT, this.storageKey)
      .map(result => result.data)
      .catch(err => {
        try {
          if (err['_body']) {
            const errorBody = JSON.parse(err['_body']);
            if (errorBody && errorBody.error && errorBody.error.code === 404) {
              return Observable.of(null);
            }
          }
        } catch (e) {}

        return Observable.throw(err);
      });
  }

  saveWalletsData(walletsData: string) {
    if (this.savingDataSubscription) {
      this.savingDataSubscription.unsubscribe();
    }

    this.savingDataSubscription = this.storageService.store(StorageType.CLIENT, this.storageKey, walletsData).subscribe();
  }

  cancelLastAction(): Observable<OperationResult> {
    if (!AppConfig.useHwWalletDaemon) {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwCancelLastAction', requestId);

      return this.createRequestResponse(requestId);
    } else {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.put('/cancel', null, false, true),
      );
    }
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwGetAddresses', requestId, addressN, startIndex, false);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        const params = {
          address_n: addressN,
          start_index: startIndex,
          confirm_address: false,
        };

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/generate_addresses',
            params,
          ), null, true,
        );
      }
    }).flatMap(response => {
      return this.verifyAddresses(response.rawResponse, 0)
        .catch(() => Observable.throw({ _body: this.translate.instant('hardware-wallet.errors.invalid-address-generated') }))
        .map(() => response);
    });
  }

  private verifyAddresses(addresses: string[], currentIndex: number): Observable<any> {
    const params = {
      address: addresses[currentIndex],
    };

    return this.apiService.post('address/verify', params, {}, true).flatMap(() => {
      if (currentIndex !== addresses.length - 1) {
        return this.verifyAddresses(addresses, currentIndex + 1);
      } else {
        return Observable.of(0);
      }
    });
  }

  confirmAddress(index: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwGetAddresses', requestId, 1, index, true);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        const params = {
          address_n: 1,
          start_index: index,
          confirm_address: true,
        };

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/generate_addresses',
            params,
          ), null, true,
        );
      }
    });
  }

  getFeatures(cancelPreviousOperation = true): Observable<OperationResult> {

    let cancel: Observable<any>;
    if (cancelPreviousOperation) {
      cancel = this.cancelLastAction();
    } else {
      cancel = Observable.of(0);
    }

    return cancel.flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwGetFeatures', requestId);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        return this.processDaemonResponse(
          this.hwWalletDaemonService.get('/features'),
        );
      }
    });
  }

  updateFirmware(downloadCompleteCallback: () => any): Observable<OperationResult> {
    if (!AppConfig.useHwWalletDaemon) {
      // Unimplemented.
      return null;
    } else {
      this.prepare();

      return this.getFeatures(false).flatMap(result => {
        if (!result.rawResponse.bootloader_mode) {
          const response: OperationResult = {
            result: OperationResults.NotInBootloaderMode,
            rawResponse: null,
          };

          return Observable.throw(response);
        }

        return this.http.get(AppConfig.urlForHwWalletVersionChecking, { responseType: 'text' })
          .catch(() => Observable.throw({ _body: this.translate.instant('hardware-wallet.update-firmware.connection-error') }))
          .flatMap((res: any) => {
            let lastestFirmwareVersion: string = res.trim();
            if (lastestFirmwareVersion.toLowerCase().startsWith('v')) {
              lastestFirmwareVersion = lastestFirmwareVersion.substr(1, lastestFirmwareVersion.length - 1);
            }

            return this.http.get(AppConfig.hwWalletDownloadUrlAndPrefix + lastestFirmwareVersion + '.bin', { responseType: 'arraybuffer' })
              .catch(() => Observable.throw({ _body: this.translate.instant('hardware-wallet.update-firmware.connection-error') }))
              .flatMap(firmware => {
                downloadCompleteCallback();
                const data = new FormData();
                data.set('file', new Blob([firmware], { type: 'application/octet-stream'}));

                return this.processDaemonResponse(
                  this.hwWalletDaemonService.put('/firmware_update', data, true),
                );
              });
          });
      });
    }
  }

  changePin(changingCurrentPin: boolean): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      let requestId = 0;
      if (!AppConfig.useHwWalletDaemon) {
        requestId = this.createRandomIdAndPrepare();
      } else {
        this.prepare();
      }

      this.hwWalletPinService.changingPin = true;
      if (changingCurrentPin) {
        this.hwWalletPinService.changePinState = ChangePinStates.RequestingCurrentPin;
      } else {
        this.hwWalletPinService.changePinState = ChangePinStates.RequestingNewPin;
      }

      if (!AppConfig.useHwWalletDaemon) {
        window['ipcRenderer'].send('hwChangePin', requestId);

        return this.createRequestResponse(requestId);
      } else {
        return this.processDaemonResponse(
          this.hwWalletDaemonService.post('/configure_pin_code'),
          ['PIN changed'],
        );
      }
    });
  }

  removePin(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
       if (!AppConfig.useHwWalletDaemon) {
        // Unimplemented.
        return null;
      } else {
        this.prepare();

        const params = {};
        params['remove_pin'] = true;

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/configure_pin_code',
            params,
          ),
          ['PIN removed'],
        );
      }
    });
  }

  generateMnemonic(wordCount: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
       if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwGenerateMnemonic', requestId, wordCount);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        const params = {};
        params['word_count'] = wordCount;
        params['use_passphrase'] = false;

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/generate_mnemonic',
            params,
          ),
          ['Mnemonic successfully configured'],
        );
      }
    });
  }

  recoverMnemonic(wordCount: number, dryRun: boolean): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwRecoverMnemonic', requestId, wordCount, dryRun);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        const params = {};
        params['word_count'] = wordCount;
        params['use_passphrase'] = false;
        params['dry_run'] = dryRun;

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/recovery',
            params,
          ),
          ['Device recovered', 'The seed is valid and matches the one in the device'],
        );
      }
    });
  }

  backup(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwBackupDevice', requestId);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/backup',
          ),
          ['Device backed up!'],
        );
      }
    });
  }

  wipe(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwWipe', requestId);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        return this.processDaemonResponse(
          this.hwWalletDaemonService.delete('/wipe'),
          ['Device wiped'],
        );
      }
    });
  }

  changeLabel(label: string): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        window['ipcRenderer'].send('hwChangeLabel', requestId, label);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post('/apply_settings', {label: label}),
          ['Settings applied'],
        );
      }
    });
  }

  signTransaction(inputs: Input[], outputs: Output[]): Observable<OperationResult> {
    const previewData: TxData[] = [];
    outputs.forEach(output => {
      if (output.address_index === undefined || output.address_index === null) {
        const currentOutput = new TxData();
        currentOutput.address = output.address;
        currentOutput.coins = new BigNumber(output.coin).dividedBy(1000000);
        currentOutput.hours = new BigNumber(output.hour);

        previewData.push(currentOutput);
      }
    });

    this.signTransactionDialog = this.dialog.open(this.signTransactionConfirmationComponentInternal, <MatDialogConfig> {
      width: '560px',
      data: previewData,
    });

    return this.cancelLastAction().flatMap(() => {
      if (!AppConfig.useHwWalletDaemon) {
        const requestId = this.createRandomIdAndPrepare();
        this.hwWalletPinService.signingTx = true;
        window['ipcRenderer'].send('hwSignTransaction', requestId, inputs, outputs);

        return this.createRequestResponse(requestId);
      } else {
        this.prepare();

        const params = {
          transaction_inputs: (inputs as any[]).map(val => {
            return {
              index: val.index,
              hash: val.hashIn,
            };
          }),
          transaction_outputs : (outputs as any[]).map(val => {
            return {
              address_index: val.address_index,
              address: val.address,
              coins: new BigNumber(val.coin).dividedBy(1000000).toFixed(6),
              hours: val.hour.toString(),
            };
          }),
        };

        return this.processDaemonResponse(
          this.hwWalletDaemonService.post(
            '/transaction_sign',
            params,
          ), null, true,
        ).map(response => {
          this.closeTransactionDialog();

          return response;
        }).catch(error => {
          this.closeTransactionDialog();

          return Observable.throw(error);
        });
      }
    });
  }

  checkIfCorrectHwConnected(firstAddress: string): Observable<boolean> {
    return this.getAddresses(1, 0).flatMap(
      response => {
        if (response.rawResponse[0] !== firstAddress) {
          return Observable.throw({
            result: OperationResults.IncorrectHardwareWallet,
            rawResponse: '',
          });
        }

        return Observable.of(true);
      },
    ).catch(error => {
      if (error.result && error.result === OperationResults.WithoutSeed) {
        return Observable.throw({
          result: OperationResults.IncorrectHardwareWallet,
          rawResponse: '',
        });
      }

      return Observable.throw(error);
    });
  }

  private closeTransactionDialog() {
    if (this.signTransactionDialog) {
      this.signTransactionDialog.close();
      this.signTransactionDialog = null;
    }
  }

  private processDaemonResponse(daemonResponse: Observable<any>, successTexts: string[] = null, responseShouldBeArray = false) {
    return daemonResponse.catch((error: any) => {
      return Observable.throw(this.dispatchEvent(0, error['_body'], false, true));
    }).flatMap(result => {
      if (result !== HwWalletDaemonService.errorCancelled) {
        if (responseShouldBeArray && result.data && typeof result.data === 'string') {
          result.data = [result.data];
        }

        const response = this.dispatchEvent(0,
          result.data ? result.data : null,
          !successTexts ? true : typeof result.data === 'string' && successTexts.some(text => (result.data as string).includes(text)),
          true);

          if (response.result === OperationResults.Success) {
            return Observable.of(response);
          } else {
            return Observable.throw(response);
          }
      } else {
        return Observable.throw(this.dispatchEvent(0, 'canceled by user', false, true));
      }
    });
  }

  private createRequestResponse(requestId: number): Observable<OperationResult> {
    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  private createRandomIdAndPrepare(): number {
    this.prepare();

    return this.requestSequence++;
  }

  private prepare() {
    this.hwWalletPinService.changingPin = false;
    this.hwWalletPinService.signingTx = false;
  }

  private dispatchEvent(requestId: number, rawResponse: any, success: boolean, justReturnTheEvent = false) {
    if (this.eventsObservers.has(requestId) || justReturnTheEvent) {
      if ((!rawResponse || !rawResponse.error) && success) {
        const response: OperationResult = {
          result: OperationResults.Success,
          rawResponse: rawResponse,
        };

        if (justReturnTheEvent) {
          return response;
        } else {
          this.eventsObservers.get(requestId).next(response);
        }
      } else {
        let responseContent: string = rawResponse.error ? rawResponse.error : rawResponse;
        if (typeof responseContent !== 'string') {
          responseContent = '';
        }
        let result: OperationResults;

        if (responseContent.includes('failed or refused')) {
          result = OperationResults.FailedOrRefused;
        } else if (responseContent.includes('PIN invalid')) {
          result = OperationResults.WrongPin;
        } else if (responseContent.includes('canceled by user')) {
          result = OperationResults.FailedOrRefused;
        } else if (responseContent.includes('cancelled by user')) {
          result = OperationResults.FailedOrRefused;
        } else if (responseContent.includes('Expected WordAck after Button')) {
          result = OperationResults.FailedOrRefused;
        } else if (responseContent.includes('Wrong word retyped')) {
          result = OperationResults.WrongWord;
        } else if (responseContent.includes('PIN mismatch')) {
          result = OperationResults.PinMismatch;
        } else if (responseContent.includes('Mnemonic not set')) {
          result = OperationResults.WithoutSeed;
        } else if (responseContent.includes('Mnemonic required')) {
          result = OperationResults.WithoutSeed;
        } else if (responseContent.includes('Invalid seed, are words in correct order?')) {
          result = OperationResults.InvalidSeed;
        } else if (responseContent.includes('The seed is valid but does not match the one in the device')) {
          result = OperationResults.WrongSeed;
        } else if (responseContent.includes('Invalid base58 character')) {
          result = OperationResults.InvalidAddress;
        } else if (responseContent.includes('Invalid address length')) {
          result = OperationResults.InvalidAddress;
        } else if (responseContent.toLocaleLowerCase().includes('LIBUSB'.toLocaleLowerCase())) {
          result = OperationResults.DaemonError;
        } else if (responseContent.toLocaleLowerCase().includes('hidapi'.toLocaleLowerCase())) {
          result = OperationResults.Disconnected;
          if (AppConfig.useHwWalletDaemon) {
            setTimeout(() => this.hwWalletDaemonService.checkHw(false));
          }
        } else if (responseContent.toLocaleLowerCase().includes('device disconnected'.toLocaleLowerCase())) {
          result = OperationResults.Disconnected;
          if (AppConfig.useHwWalletDaemon) {
            setTimeout(() => this.hwWalletDaemonService.checkHw(false));
          }
        } else if (responseContent.toLocaleLowerCase().includes('no device connected'.toLocaleLowerCase())) {
          result = OperationResults.Disconnected;
          if (AppConfig.useHwWalletDaemon) {
            setTimeout(() => this.hwWalletDaemonService.checkHw(false));
          }
        } else if (responseContent.includes(HwWalletDaemonService.errorConnectingWithTheDaemon)) {
          result = OperationResults.DaemonError;
        } else if (responseContent.includes(HwWalletDaemonService.errorTimeout)) {
          result = OperationResults.Timeout;
        } else if (responseContent.includes('MessageType_Success')) {
          result = OperationResults.Success;
        } else {
          result = OperationResults.UndefinedError;
        }

        const response: OperationResult = {
          result: result,
          rawResponse: responseContent,
        };

        if (justReturnTheEvent) {
          return response;
        } else {
          this.eventsObservers.get(requestId).error(response);
        }
      }
      if (!justReturnTheEvent) {
        this.eventsObservers.get(requestId).complete();
        this.eventsObservers.delete(requestId);
      }
    }
  }

  private dispatchError(requestId: number, result: OperationResults, error: String) {
    if (this.eventsObservers.has(requestId)) {
      this.eventsObservers.get(requestId).error({
        result: result,
        rawResponse: error,
      });
      this.eventsObservers.delete(requestId);
    }
  }

  private cancelAllOperations() {
    this.eventsObservers.forEach((value, key) => {
      this.dispatchEvent(key, 'failed or refused', false);
    });
  }

}
