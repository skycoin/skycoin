import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscriber } from 'rxjs/Subscriber';
import { Subject } from 'rxjs/Subject';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';
import { MatDialog, MatDialogConfig } from '@angular/material';

export class OperationResult {
  success: boolean;
  rawResponse: any;
}

@Injectable()
export class HwWalletService {

  private eventsObservers = new Map<number, Subscriber<OperationResult>>();
  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();

  // Set on AppComponent to avoid a circular reference.
  private requestPinComponentInternal;
  set requestPinComponent(value) {
    this.requestPinComponentInternal = value;
  }

  constructor(private translate: TranslateService, private dialog: MatDialog) {
    if (window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated')) {
      window['ipcRenderer'].on('hwConnectionEvent', (event, connected) => {
        if (!connected) {
          this.eventsObservers.forEach((value, key) => {
            this.dispatchError(key, this.translate.instant('hardware-wallet.general.error-disconnected'));
          });
        }
        this.walletConnectedSubject.next(connected);
      });
      window['ipcRenderer'].on('hwPinRequested', (event) => {
        dialog.open(this.requestPinComponentInternal, <MatDialogConfig> {
          width: '450px',
        }).afterClosed().subscribe(pin => {
          if (!pin) {
            this.cancelAllOperations();
          }
          window['ipcRenderer'].send('hwSendPin', pin);
        });
      });

      window['ipcRenderer'].on('hwGetAddressesResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, true);
      });
      window['ipcRenderer'].on('hwChangePinResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, typeof result === 'string' && (result as string).includes('PIN changed'));
      });
      window['ipcRenderer'].on('hwSetMnemonicResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, typeof result === 'string' && (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwGenerateMnemonicResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, typeof result === 'string' && (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwBackupDeviceResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, typeof result === 'string' && (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwWipeResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, typeof result === 'string' && (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwSignMessageResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, true);
      });
    }
  }

  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  getDeviceSync() {
    return window['ipcRenderer'].sendSync('hwGetDeviceSync');
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwGetAddresses', requestId, addressN, startIndex);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  changePin(): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwChangePin', requestId);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  getMaxAddresses(): Observable<string[]> {
    return this.getAddressesRecursively(AppConfig.maxHardwareWalletAddresses - 1, []);
  }

  setMnemonic(mnemonic: string): Observable<OperationResult> {
    mnemonic = mnemonic.replace(/(\n|\r\n)$/, '');

    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwSetMnemonic', requestId, mnemonic);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  generateMnemonic(): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwGenerateMnemonic', requestId);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  backup(): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwBackupDevice', requestId);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  wipe(): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwWipe', requestId);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  signMessage(addressIndex: number, message: string): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwSignMessage', requestId, addressIndex, message);

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

  private createRandomID() {
    return Math.floor(Math.random() * 4000000000);
  }

  private dispatchEvent(requestId: number, rawResponse: any, success: boolean) {
    if (this.eventsObservers.has(requestId)) {
      if (!rawResponse.error) {
        this.eventsObservers.get(requestId).next({
          success: success,
          rawResponse: rawResponse,
        });
      } else {
        if (typeof rawResponse.error === 'string' && (rawResponse.error as string).includes('cancelled by user')) {
          this.eventsObservers.get(requestId).next({
            success: false,
            rawResponse: rawResponse,
          });
        } else {
          this.eventsObservers.get(requestId).error(rawResponse.error);
        }
      }
      this.eventsObservers.get(requestId).complete();
      this.eventsObservers.delete(requestId);
    }
  }

  private dispatchError(requestId: number, error: String) {
    if (this.eventsObservers.has(requestId)) {
      this.eventsObservers.get(requestId).error(error);
      this.eventsObservers.delete(requestId);
    }
  }

  private cancelAllOperations() {
    this.eventsObservers.forEach((value, key) => {
      this.dispatchEvent(key, '', false);
    });
  }

}
