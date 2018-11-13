import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscriber } from 'rxjs/Subscriber';

export class OperationResult {
  success: boolean;
  rawResponse: any;
}

@Injectable()
export class HwWalletService {

  private eventsObservers = new Map<number, Subscriber<OperationResult>>();

  constructor() {
    if (window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated')) {
      window['ipcRenderer'].on('hwGetAddressesResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, true);
      });
      window['ipcRenderer'].on('hwSetMnemonicResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwGenerateMnemonicResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwWipeResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, (result as string).includes('operation completed'));
      });
      window['ipcRenderer'].on('hwSignMessageResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result, true);
      });
    }
  }

  getDevice() {
    return window['ipcRenderer'].sendSync('hwGetDevice');
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwGetAddresses', requestId, addressN, startIndex);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  setMnemonic(mnemonic: string): Observable<OperationResult> {
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
        this.eventsObservers.get(requestId).error(rawResponse.error);
      }
      this.eventsObservers.get(requestId).complete();
      this.eventsObservers.delete(requestId);
    }
  }

}
