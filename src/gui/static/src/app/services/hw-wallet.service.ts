import { throwError as observableThrowError, of, Observable, Subject, SubscriptionLike } from 'rxjs';
import { mergeMap, map, catchError } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { HwWalletDaemonService } from './hw-wallet-daemon.service';
import { HwWalletPinService, ChangePinStates } from './hw-wallet-pin.service';
import BigNumber from 'bignumber.js';
import { StorageService, StorageType } from './storage.service';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
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

export interface HwInput {
  hash: string;
  index: number;
}

export interface HwOutput {
  address: string;
  coins: string;
  hours: string;
  address_index?: number;
}

@Injectable()
export class HwWalletService {

  public static readonly maxLabelLength = 32;

  private readonly storageKey = 'hw-wallets';

  showOptionsWhenPossible = false;

  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();

  private savingDataSubscription: SubscriptionLike;

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
    private storageService: StorageService,
    private apiService: ApiService,
    private http: HttpClient) {

    if (this.hwWalletCompatibilityActivated) {
      hwWalletDaemonService.connectionEvent.subscribe(connected => {
        this.walletConnectedSubject.next(connected);
      });
    }
  }

  get hwWalletCompatibilityActivated(): boolean {
    return true;
  }

  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  getDeviceConnected(): Observable<boolean> {
    return this.hwWalletDaemonService.get('/available').pipe(map((response: any) => {
      return response.data;
    }));
  }

  getSavedWalletsData(): Observable<string> {
    return this.storageService.get(StorageType.CLIENT, this.storageKey).pipe(
      map(result => result.data),
      catchError((err: HttpErrorResponse) => {
        try {
          if (err.status && err.status === 404) {
            return of(null);
          }
        } catch (e) {}

        return observableThrowError(err);
      }));
  }

  saveWalletsData(walletsData: string) {
    if (this.savingDataSubscription) {
      this.savingDataSubscription.unsubscribe();
    }

    this.savingDataSubscription = this.storageService.store(StorageType.CLIENT, this.storageKey, walletsData).subscribe();
  }

  cancelLastAction(): Observable<OperationResult> {
    this.prepare();

    return this.processDaemonResponse(
      this.hwWalletDaemonService.put('/cancel', null, false, true),
    );
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
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
    }), mergeMap(response => {
      return this.verifyAddresses(response.rawResponse, 0).pipe(
        catchError(() => observableThrowError({ _body: this.translate.instant('hardware-wallet.errors.invalid-address-generated') })),
        map(() => response));
    }));
  }

  private verifyAddresses(addresses: string[], currentIndex: number): Observable<any> {
    const params = {
      address: addresses[currentIndex],
    };

    return this.apiService.post('address/verify', params, {}, true).pipe(mergeMap(() => {
      if (currentIndex !== addresses.length - 1) {
        return this.verifyAddresses(addresses, currentIndex + 1);
      } else {
        return of(0);
      }
    }));
  }

  confirmAddress(index: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
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
    }));
  }

  getFeatures(cancelPreviousOperation = true): Observable<OperationResult> {

    let cancel: Observable<any>;
    if (cancelPreviousOperation) {
      cancel = this.cancelLastAction();
    } else {
      cancel = of(0);
    }

    return cancel.pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.get('/features'),
      );
    }));
  }

  updateFirmware(downloadCompleteCallback: () => any): Observable<OperationResult> {
    this.prepare();

    return this.getFeatures(false).pipe(mergeMap(result => {
      if (!result.rawResponse.bootloader_mode) {
        const response: OperationResult = {
          result: OperationResults.NotInBootloaderMode,
          rawResponse: null,
        };

        return observableThrowError(response);
      }

      return this.http.get(AppConfig.urlForHwWalletVersionChecking, { responseType: 'text' }).pipe(
        catchError(() => observableThrowError({ _body: this.translate.instant('hardware-wallet.update-firmware.connection-error') })),
        mergeMap((res: any) => {
          let lastestFirmwareVersion: string = res.trim();
          if (lastestFirmwareVersion.toLowerCase().startsWith('v')) {
            lastestFirmwareVersion = lastestFirmwareVersion.substr(1, lastestFirmwareVersion.length - 1);
          }

          return this.http.get(AppConfig.hwWalletDownloadUrlAndPrefix + lastestFirmwareVersion + '.bin', { responseType: 'arraybuffer' }).pipe(
            catchError(() => observableThrowError({ _body: this.translate.instant('hardware-wallet.update-firmware.connection-error') })),
            mergeMap(firmware => {
              downloadCompleteCallback();
              const data = new FormData();
              data.set('file', new Blob([firmware], { type: 'application/octet-stream'}));

              return this.processDaemonResponse(
                this.hwWalletDaemonService.put('/firmware_update', data, true),
              );
            }));
        }));
    }));
  }

  changePin(changingCurrentPin: boolean): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      this.hwWalletPinService.changingPin = true;
      if (changingCurrentPin) {
        this.hwWalletPinService.changePinState = ChangePinStates.RequestingCurrentPin;
      } else {
        this.hwWalletPinService.changePinState = ChangePinStates.RequestingNewPin;
      }

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post('/configure_pin_code'),
        ['PIN changed'],
      );
    }));
  }

  removePin(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
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
    }));
  }

  generateMnemonic(wordCount: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
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
    }));
  }

  recoverMnemonic(wordCount: number, dryRun: boolean): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
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
    }));
  }

  backup(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/backup',
        ),
        ['Device backed up!'],
      );
    }));
  }

  wipe(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.delete('/wipe'),
        ['Device wiped'],
      );
    }));
  }

  changeLabel(label: string): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post('/apply_settings', {label: label}),
        ['Settings applied'],
      );
    }));
  }

  signTransaction(inputs: HwInput[], outputs: HwOutput[]): Observable<OperationResult> {
    const previewData: TxData[] = [];
    outputs.forEach(output => {
      if (output.address_index === undefined || output.address_index === null) {
        const currentOutput = new TxData();
        currentOutput.address = output.address;
        currentOutput.coins = new BigNumber(output.coins).decimalPlaces(6);
        currentOutput.hours = new BigNumber(output.hours);

        previewData.push(currentOutput);
      }
    });

    this.signTransactionDialog = this.dialog.open(this.signTransactionConfirmationComponentInternal, <MatDialogConfig> {
      width: '600px',
      data: previewData,
    });

    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {
        transaction_inputs: inputs,
        transaction_outputs: outputs,
      };

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/transaction_sign',
          params,
        ), null, true,
      ).pipe(map(response => {
        this.closeTransactionDialog();

        return response;
      }), catchError(error => {
        this.closeTransactionDialog();

        return observableThrowError(error);
      }));
    }));
  }

  checkIfCorrectHwConnected(firstAddress: string): Observable<boolean> {
    return this.getAddresses(1, 0).pipe(mergeMap(
      response => {
        if (response.rawResponse[0] !== firstAddress) {
          return observableThrowError({
            result: OperationResults.IncorrectHardwareWallet,
            rawResponse: '',
          });
        }

        return of(true);
      },
    ), catchError(error => {
      if (error.result && error.result === OperationResults.WithoutSeed) {
        return observableThrowError({
          result: OperationResults.IncorrectHardwareWallet,
          rawResponse: '',
        });
      }

      return observableThrowError(error);
    }));
  }

  private closeTransactionDialog() {
    if (this.signTransactionDialog) {
      this.signTransactionDialog.close();
      this.signTransactionDialog = null;
    }
  }

  private processDaemonResponse(daemonResponse: Observable<any>, successTexts: string[] = null, responseShouldBeArray = false) {
    return daemonResponse.pipe(catchError((error: any) => {
      return observableThrowError(this.dispatchEvent(error['_body'], false));
    }), mergeMap(result => {
      if (result !== HwWalletDaemonService.errorCancelled) {
        if (responseShouldBeArray && result.data && typeof result.data === 'string') {
          result.data = [result.data];
        }

        const response = this.dispatchEvent(
          result.data ? result.data : null,
          !successTexts ? true : typeof result.data === 'string' && successTexts.some(text => (result.data as string).includes(text)));

          if (response.result === OperationResults.Success) {
            return of(response);
          } else {
            return observableThrowError(response);
          }
      } else {
        return observableThrowError(this.dispatchEvent('canceled by user', false));
      }
    }));
  }

  private prepare() {
    this.hwWalletPinService.changingPin = false;
    this.hwWalletPinService.signingTx = false;
  }

  private dispatchEvent(rawResponse: any, success: boolean) {
    if ((!rawResponse || !rawResponse.error) && success) {
      const response: OperationResult = {
        result: OperationResults.Success,
        rawResponse: rawResponse,
      };

      return response;
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
        setTimeout(() => this.hwWalletDaemonService.checkHw(false));
      } else if (responseContent.toLocaleLowerCase().includes('device disconnected'.toLocaleLowerCase())) {
        result = OperationResults.Disconnected;
        setTimeout(() => this.hwWalletDaemonService.checkHw(false));
      } else if (responseContent.toLocaleLowerCase().includes('no device connected'.toLocaleLowerCase())) {
        result = OperationResults.Disconnected;
        setTimeout(() => this.hwWalletDaemonService.checkHw(false));
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

      return response;
    }
  }
}
