import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { Subscriber } from 'rxjs/Subscriber';

@Injectable()
export class HwWalletService {

  private eventsObservers = new Map<number, Subscriber<{}>>();

  constructor() {
    if (window['isElectron'] && window['ipcRenderer'].sendSync('hwCompatibilityActivated')) {
      window['ipcRenderer'].on('hwGetAddressesResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result);
      });
      window['ipcRenderer'].on('hwSetMnemonicResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result);
      });
      window['ipcRenderer'].on('hwWipeResponse', (event, requestId, result) => {
        this.dispatchEvent(requestId, result);
      });

      window['ipcRenderer'].on('hwConnectionEvent', (event, params) => {
        alert(JSON.stringify(params));
      });
    }
  }

  getDevice() {
    return window['ipcRenderer'].sendSync('hwGetDevice');
  }

  getAddresses(addressN, startIndex): Observable<any> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwGetAddresses', requestId, addressN, startIndex);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  setMnemonic(mnemonic): Observable<any> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwSetMnemonic', requestId, mnemonic);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  wipe(): Observable<any> {
    const requestId = this.createRandomID();
    window['ipcRenderer'].send('hwWipe', requestId);

    return new Observable(observer => {
      this.eventsObservers.set(requestId, observer);
    });
  }

  private createRandomID() {
    return Math.floor(Math.random() * 4000000000);
  }

  private dispatchEvent(requestId, result) {
    if (this.eventsObservers.has(requestId)) {
      if (!result.error) {
        this.eventsObservers.get(requestId).next(result);
      } else {
        this.eventsObservers.get(requestId).error(result.error);
      }
      this.eventsObservers.get(requestId).complete();
      this.eventsObservers.delete(requestId);
    }
  }

}
