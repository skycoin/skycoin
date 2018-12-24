import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscriber } from 'rxjs/Subscriber';
import { Subject } from 'rxjs/Subject';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';
import { MatDialog, MatDialogConfig } from '@angular/material';
import { HwPinDialogParams } from '../components/layout/hardware-wallet/hw-pin-dialog/hw-pin-dialog.component';

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
  UndefinedError,
  Disconnected,
}

export class OperationResult {
  result: OperationResults;
  rawResponse: any;
}

interface EventData {
  event: string;
  successText?: string;
}

@Injectable()
export class HwWalletService {

  showOptionsWhenPossible = false;

  private requestSequence = 0;

  private eventsObservers = new Map<number, Subscriber<OperationResult>>();
  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();

  // Values to be sent to HwPinDialogComponent
  private changingPin: boolean;
  private changePinState: ChangePinStates;
  private signingTx: boolean;
  private currentSignature: number;
  private totalSignatures: number;

  // Set on AppComponent to avoid a circular reference.
  private requestPinComponentInternal;
  set requestPinComponent(value) {
    this.requestPinComponentInternal = value;
  }
  private requestWordComponentInternal;
  set requestWordComponent(value) {
    this.requestWordComponentInternal = value;
  }

  constructor(private translate: TranslateService, dialog: MatDialog) {
    if (window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated')) {
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
            currentSignature: this.currentSignature,
            totalSignatures: this.totalSignatures,
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

      const data: EventData[] = [
        { event: 'hwChangePinResponse', successText: 'PIN changed' },
        { event: 'hwGenerateMnemonicResponse', successText: 'operation completed' },
        { event: 'hwRecoverMnemonicResponse', successText: 'Device recovered' },
        { event: 'hwBackupDeviceResponse', successText: 'operation completed' },
        { event: 'hwWipeResponse', successText: 'operation completed' },
        { event: 'hwCancelLastActionResponse' },
        { event: 'hwGetAddressesResponse' },
        { event: 'hwGetFeaturesResponse' },
        { event: 'hwSignMessageResponse' },
      ];

      data.forEach(item => {
        window['ipcRenderer'].on(item.event, (event, requestId, result) => {
          const success = item.successText
            ? typeof result === 'string' && (result as string).includes(item.successText)
            : true;

          this.dispatchEvent(requestId, result, success);
        });
      });
    }
  }

  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  getDeviceConnectedSync() {
    return window['ipcRenderer'].sendSync('hwGetDeviceConnectedSync');
  }

  cancelLastAction(): Observable<OperationResult> {
    const requestId = this.createRandomIdAndPrepare();
    window['ipcRenderer'].send('hwCancelLastAction', requestId);

    return this.createRequestResponse(requestId);
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwGetAddresses', requestId, addressN, startIndex);

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

  getMaxAddresses(): Observable<string[]> {
    return this.getAddressesRecursively(AppConfig.maxHardwareWalletAddresses - 1, []);
  }

  generateMnemonic(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwGenerateMnemonic', requestId);

      return this.createRequestResponse(requestId);
    });
  }

  recoverMnemonic(): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      window['ipcRenderer'].send('hwRecoverMnemonic', requestId);

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

  signMessage(addressIndex: number, message: string, currentSignature: number, totalSignatures: number): Observable<OperationResult> {
    return this.cancelLastAction().flatMap(() => {
      const requestId = this.createRandomIdAndPrepare();
      this.signingTx = true;
      this.currentSignature = currentSignature;
      this.totalSignatures = totalSignatures;
      window['ipcRenderer'].send('hwSignMessage', requestId, addressIndex, message);

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

  private getAddressesRecursively(index: number, addresses: string[]): Observable<string[]> {
    let chain: Observable<any>;
    if (index > 0) {
      chain = this.getAddressesRecursively(index - 1, addresses).first();
    } else {
      chain = Observable.of(1);
    }

    chain = chain.flatMap(() => {
      return this.getAddresses(1, index)
      .map(response => {
        addresses.push(response.rawResponse[0]);

        return addresses;
      });
    });

    return chain;
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
