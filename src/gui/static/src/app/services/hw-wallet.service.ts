import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscriber } from 'rxjs/Subscriber';
import { Subject } from 'rxjs/Subject';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material';
import { HwPinDialogParams } from '../components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';
import { environment } from '../../environments/environment';

export enum ChangePinStates {
  RequestingCurrentPin,
  RequestingNewPin,
  ConfirmingNewPin,
}

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
}

export class OperationResult {
  result: OperationResults;
  rawResponse: any;
}

interface EventData {
  event: string;
  successTexts?: string[];
}

@Injectable()
export class HwWalletService {

  showOptionsWhenPossible = false;

  private requestSequence = 0;

  private eventsObservers = new Map<number, Subscriber<OperationResult>>();
  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();

  private signTransactionDialog: MatDialogRef<{}, any>;

  // Values to be sent to HwPinDialogComponent
  private changingPin: boolean;
  private changePinState: ChangePinStates;
  private signingTx: boolean;

  // Set on AppComponent to avoid a circular reference.
  private requestPinComponentInternal;
  set requestPinComponent(value) {
    this.requestPinComponentInternal = value;
  }
  private requestWordComponentInternal;
  set requestWordComponent(value) {
    this.requestWordComponentInternal = value;
  }
  private signTransactionConfirmationComponentInternal;
  set signTransactionConfirmationComponent(value) {
    this.signTransactionConfirmationComponentInternal = value;
  }

  constructor(private translate: TranslateService, private dialog: MatDialog) {
    if (this.hwWalletCompatibilityActivated) {
      window['ipcRenderer'].on('hwConnectionEvent', (event, connected) => {
        if (!connected) {
          this.eventsObservers.forEach((value, key) => {
            this.dispatchError(key, OperationResults.Disconnected, this.translate.instant('hardware-wallet.general.error-disconnected'));
          });
        }
        this.walletConnectedSubject.next(connected);
      });
      window['ipcRenderer'].on('hwPinRequested', (event) => {
        dialog.open(this.requestPinComponentInternal, <MatDialogConfig> {
          width: '350px',
          autoFocus: false,
          data : <HwPinDialogParams> {
            changingPin: this.changingPin,
            changePinState: this.changePinState,
            signingTx: this.signingTx,
          },
        }).afterClosed().subscribe(pin => {
          if (!pin) {
            this.cancelAllOperations();
          } else {
            if (this.changingPin) {
              if (this.changePinState === ChangePinStates.RequestingCurrentPin) {
                this.changePinState = ChangePinStates.RequestingNewPin;
              } else if (this.changePinState === ChangePinStates.RequestingNewPin) {
                this.changePinState = ChangePinStates.ConfirmingNewPin;
              }
            }
          }
          window['ipcRenderer'].send('hwSendPin', pin);
        });
      });
      window['ipcRenderer'].on('hwSeedWordRequested', (event) => {
        dialog.open(this.requestWordComponentInternal, <MatDialogConfig> {
          width: '350px',
        }).afterClosed().subscribe(word => {
          if (!word) {
            this.cancelAllOperations();
            window['ipcRenderer'].send('hwCancelLastAction');
          }
          window['ipcRenderer'].send('hwSendSeedWord', word);
        });
      });

      window['ipcRenderer'].on('hwSignTransactionResponse', (event, requestId, result) => {
        if (this.signTransactionDialog) {
          this.signTransactionDialog.close();
          this.signTransactionDialog = null;
        }

        this.dispatchEvent(requestId, result, true);
      });

      const data: EventData[] = [
        { event: 'hwChangePinResponse', successTexts: ['PIN changed'] },
        { event: 'hwGenerateMnemonicResponse', successTexts: ['operation completed'] },
        { event: 'hwRecoverMnemonicResponse', successTexts: ['Device recovered', 'The seed is valid and matches the one in the device'] },
        { event: 'hwBackupDeviceResponse', successTexts: ['operation completed'] },
        { event: 'hwWipeResponse', successTexts: ['operation completed'] },
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

  get hwWalletCompatibilityActivated(): boolean {
    // return !environment.production && window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated');
    return false;
  }

  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  getDeviceConnectedSync() {
    return window['ipcRenderer'].sendSync('hwGetDeviceConnectedSync');
  }

  getSavedWalletsDataSync(): string {
    return window['ipcRenderer'].sendSync('hwGetSavedWalletsDataSync');
  }

  saveWalletsDataSync(walletsData: string) {
    window['ipcRenderer'].sendSync('hwSaveWalletsDataSync', walletsData);
  }

  cancelLastAction(): Observable<OperationResult> {
    const requestId = this.createRandomIdAndPrepare();
    window['ipcRenderer'].send('hwCancelLastAction', requestId);

    return this.createRequestResponse(requestId);
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwGetAddresses', requestId, addressN, startIndex, false);

      return this.createRequestResponse(requestId);
    });
  }

  confirmAddress(index: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwGetAddresses', requestId, 1, index, true);

      return this.createRequestResponse(requestId);
    });
  }

  getFeatures(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwGetFeatures', requestId);

      return this.createRequestResponse(requestId);
    });
  }

  changePin(changingCurrentPin: boolean): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      this.changingPin = true;
      if (changingCurrentPin) {
        this.changePinState = ChangePinStates.RequestingCurrentPin;
      } else {
        this.changePinState = ChangePinStates.RequestingNewPin;
      }
      window['ipcRenderer'].send('hwChangePin', requestId);

      return this.createRequestResponse(requestId);
    });
  }

  generateMnemonic(wordCount: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwGenerateMnemonic', requestId, wordCount);

      return this.createRequestResponse(requestId);
    });
  }

  recoverMnemonic(wordCount: number, dryRun: boolean): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwRecoverMnemonic', requestId, wordCount, dryRun);

      return this.createRequestResponse(requestId);
    });
  }

  backup(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwBackupDevice', requestId);

      return this.createRequestResponse(requestId);
    });
  }

  wipe(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwWipe', requestId);

      return this.createRequestResponse(requestId);
    });
  }

  signMessage(addressIndex: number, message: string): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      this.signingTx = true;
      window['ipcRenderer'].send('hwSignMessage', requestId, addressIndex, message);

      return this.createRequestResponse(requestId);
    });
  }

  signTransaction(inputs: any, outputs: any): Observable<OperationResult> {
    this.signTransactionDialog = this.dialog.open(this.signTransactionConfirmationComponentInternal, <MatDialogConfig> {
      width: '450px',
    });

    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      this.signingTx = true;
      window['ipcRenderer'].send('hwSignTransaction', requestId, inputs, outputs);

      return this.createRequestResponse(requestId);
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

  private createRequestResponse(requestId: number): Observable<OperationResult> {
    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  private createRandomIdAndPrepare() {
    this.changingPin = false;
    this.signingTx = false;

    return this.requestSequence++;
  }

  private dispatchEvent(requestId: number, rawResponse: any, success: boolean) {
    if (this.eventsObservers.has(requestId)) {
      if (!rawResponse.error && success) {
        this.eventsObservers.get(requestId).next({
          result: OperationResults.Success,
          rawResponse: rawResponse,
        });
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
        } else if (responseContent.includes('Invalid seed, are words in correct order?')) {
          result = OperationResults.InvalidSeed;
        } else if (responseContent.includes('The seed is valid but does not match the one in the device')) {
          result = OperationResults.WrongSeed;
        } else {
          result = OperationResults.UndefinedError;
        }
        this.eventsObservers.get(requestId).error({
          result: result,
          rawResponse: responseContent,
        });
      }
      this.eventsObservers.get(requestId).complete();
      this.eventsObservers.delete(requestId);
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
