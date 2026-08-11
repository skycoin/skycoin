import { Injectable } from '@angular/core';
import { Observable, from, throwError, of } from 'rxjs';
import { catchError, mergeMap, map } from 'rxjs';
import { HttpClient } from '@angular/common/http';

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
      if (window['WebAssembly'] && (window['WebAssembly'] as any).instantiateStreaming) {
        console.log('[WASM] Starting WASM streaming instantiation...');
        const go = new Go();
        return from(
          window['WebAssembly'].instantiateStreaming(fetch('assets/scripts/skycoin-lite.wasm'), go.importObject)
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
      } else if (window['WebAssembly'] && (window['WebAssembly'] as any).instantiate) {
        console.log('[WASM] Starting WASM download (fallback mode)...');
        return this.http.get('assets/scripts/skycoin-lite.wasm', { responseType: 'arraybuffer' }).pipe(
          catchError((err) => {
            console.error('[WASM] Failed to download WASM file:', err);
            return throwError(() => InitializationResults.ErrorLoadingWasmFile);
          }),
          mergeMap((response: ArrayBuffer) => {
            console.log('[WASM] WASM file downloaded, size:', response.byteLength);
            const go = new Go();
            console.log('[WASM] Instantiating WASM module...');
            return from((window['WebAssembly'].instantiate(response, go.importObject) as Promise<any>)).pipe(
              map(result => {
                console.log('[WASM] WASM module instantiated, running...');
                go.run(result.instance);
                console.log('[WASM] Initialization complete!');

                return InitializationResults.Ok;
              }),
              catchError(err => {
                console.error('[WASM] Failed to instantiate WASM module:', err);
                return throwError(() => InitializationResults.ErrorLoadingWasmFile);
              }),
            );
          }),
        );
      } else {
        console.error('[WASM] Browser does not support WebAssembly');
        return throwError(() => InitializationResults.BrowserIncompatibleWithWasm);
      }
    }

    return null as any;
  }


  /**
   * Turns whatever the wasm cipher returned into a value or an error.
   *
   * Each entry point in src/skycoin-lite/wasm/main_wasm.go recovers from panics
   * without setting a return value, so a bad seed or a malformed input comes
   * back as null rather than as the {error: string} the happy path documents.
   * Reading .error off that is a TypeError, which surfaced as a crash instead of
   * the failed observable every caller is written against.
   */
  private resultOrError<T>(result: any, operation: string): Observable<T> {
    if (result === null || result === undefined) {
      return throwError(() => new Error(`the cipher could not ${operation}`));
    }
    if (result.error) {
      return throwError(() => new Error(result.error));
    }

    return of(result);
  }

  generateAddress(seed: any): Observable<GenerateAddressResponse> {
    const address = (window as any)['SkycoinCipher'].generateAddress(seed);

    return this.resultOrError<any>(address, 'generate an address')
      .pipe(map(result => this.convertToAddress(result)));
  }

  prepareTransaction(inputs: TransactionInput[], outputs: TransactionOutput[]): Observable<string> {
    const tx = (window as any)['SkycoinCipher'].prepareTransaction(JSON.stringify(inputs), JSON.stringify(outputs));

    return this.resultOrError<string>(tx, 'prepare the transaction');
  }

  prepareTransactionWithSignatures(inputs: TransactionInput[], outputs: TransactionOutput[], signatures: string[]): Observable<string> {
    const tx = (window as any)['SkycoinCipher'].prepareTransactionWithSignatures(
      JSON.stringify(inputs),
      JSON.stringify(outputs),
      JSON.stringify(signatures)
    );

    return this.resultOrError<string>(tx, 'prepare the signed transaction');
  }

  private convertToAddress(address: any): GenerateAddressResponse {
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
