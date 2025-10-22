import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import 'rxjs/add/observable/fromPromise';

import { Address, TransactionInput, TransactionOutput } from '../app.datatypes';

declare var Go: any;

export interface GenerateAddressResponse {
  address: Address;
  nextSeed: string;
}

export enum InitializationResults {
  Ok = 1,
  BrowserIncompatibleWithWasm = 2,
  ErrorLoadingWasmFile = 3,
}

@Injectable()
export class CipherProvider {

  private initialized = false;

  constructor(private http: HttpClient) {}

  initialize(): Observable<InitializationResults> {
    if (!this.initialized) {
      this.initialized = true;
      if (window['WebAssembly'] && window['WebAssembly'].instantiateStreaming) {
        console.log('[WASM] Starting WASM streaming instantiation...');
        const go = new Go();
        return Observable.fromPromise(
          window['WebAssembly'].instantiateStreaming(fetch('/assets/scripts/skycoin-lite.wasm'), go.importObject)
            .then((result: any) => {
              console.log('[WASM] WASM module instantiated, running...');
              go.run(result.instance);
              console.log('[WASM] Initialization complete!');
              return InitializationResults.Ok;
            })
            .catch((err: any) => {
              console.error('[WASM] Failed to instantiate WASM module:', err);
              throw InitializationResults.ErrorLoadingWasmFile;
            })
        );
      } else if (window['WebAssembly'] && window['WebAssembly'].instantiate) {
        console.log('[WASM] Starting WASM download (fallback mode)...');
        return this.http.get('/assets/scripts/skycoin-lite.wasm', { responseType: 'arraybuffer' })
          .catch((err) => {
            console.error('[WASM] Failed to download WASM file:', err);
            return Observable.throw(InitializationResults.ErrorLoadingWasmFile);
          })
          .flatMap((response: ArrayBuffer) => {
            console.log('[WASM] WASM file downloaded, size:', response.byteLength);
            const go = new Go();
            console.log('[WASM] Instantiating WASM module...');
            return Observable.fromPromise((window['WebAssembly'].instantiate(response, go.importObject) as Promise<any>)).map(result => {
              console.log('[WASM] WASM module instantiated, running...');
              go.run(result.instance);
              console.log('[WASM] Initialization complete!');

              return InitializationResults.Ok;
            }).catch(err => {
              console.error('[WASM] Failed to instantiate WASM module:', err);
              return Observable.throw(InitializationResults.ErrorLoadingWasmFile);
            });
          });
      } else {
        console.error('[WASM] Browser does not support WebAssembly');
        return Observable.throw(InitializationResults.BrowserIncompatibleWithWasm);
      }
    }

    return null;
  }

  generateAddress(seed): Observable<GenerateAddressResponse> {
    const address = window['SkycoinCipher'].generateAddress(seed);

    if (!address.error) {
      return Observable.of(this.convertToAddress(address));
    } else {
      return Observable.throw(new Error(address.error));
    }
  }

  prepareTransaction(inputs: TransactionInput[], outputs: TransactionOutput[]): Observable<string> {
    const tx = window['SkycoinCipher'].prepareTransaction(JSON.stringify(inputs), JSON.stringify(outputs));

    if (!tx.error) {
      return Observable.of(tx);
    } else {
      return Observable.throw(new Error(tx.error));
    }
  }

  private convertToAddress(address): GenerateAddressResponse {
    return {
      nextSeed: address.nextSeed,
      address: {
        secret_key: address.secret,
        public_key: address.public,
        address: address.address
      }
    };
  }
}
