import { throwError as observableThrowError, Observable } from 'rxjs';
import { first, map, mergeMap, catchError } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import {
  Address, GetWalletsResponseEntry, GetWalletsResponseWallet, NormalTransaction,
  PostWalletNewAddressResponse, Version, Wallet,
} from '../app.datatypes';

@Injectable()
export class ApiService {
  private url = environment.nodeUrl;

  constructor(
    private http: HttpClient,
    private translate: TranslateService,
  ) { }

  getTransactions(addresses: Address[]): Observable<NormalTransaction[]> {
    const formattedAddresses = addresses.map(a => a.address).join(',');

    return this.post('transactions', {addrs: formattedAddresses, verbose: true}).pipe(
      map(transactions => transactions.map(transaction => ({
        addresses: [],
        balance: new BigNumber(0),
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.txn.timestamp,
        txid: transaction.txn.txid,
        inputs: transaction.txn.inputs,
        outputs: transaction.txn.outputs,
      }))));
  }

  getVersion(): Observable<Version> {
    return this.get('version');
  }

  generateSeed(entropy: number = 128): Observable<string> {
    return this.get('wallet/newSeed', { entropy }).pipe(map(response => response.seed));
  }

  getHealth() {
    return this.get('health');
  }

  getWallets(): Observable<Wallet[]> {
    return this.get('wallets').pipe(
      map((response: GetWalletsResponseWallet[]) => {
        const wallets: Wallet[] = [];
        response.forEach(wallet => {
          const processedWallet: Wallet = {
            label: wallet.meta.label,
            filename: wallet.meta.filename,
            coins: null,
            hours: null,
            addresses: [],
            encrypted: wallet.meta.encrypted,
          };

          if (wallet.entries) {
            processedWallet.addresses = wallet.entries.map((entry: GetWalletsResponseEntry) => {
              return {
                address: entry.address,
                coins: null,
                hours: null,
                confirmed: true,
              };
            });
          }

          wallets.push(processedWallet);
        });

        return wallets;
      }));
  }

  getWalletSeed(wallet: Wallet, password: string): Observable<string> {
    return this.post('wallet/seed', { id: wallet.filename, password }).pipe(
      map(response => response.seed));
  }

  postWalletCreate(label: string, seed: string, scan: number, password: string, type: string): Observable<Wallet> {
    const params = { label, seed, scan, type };

    if (password) {
      params['password'] = password;
      params['encrypt'] = true;
    }

    return this.post('wallet/create', params).pipe(
      map(response => ({
          label: response.meta.label,
          filename: response.meta.filename,
          coins: null,
          hours: null,
          addresses: response.entries.map(entry => ({ address: entry.address, coins: null, hours: null, confirmed: true })),
          encrypted: response.meta.encrypted,
        })));
  }

  postWalletNewAddress(wallet: Wallet, num: number, password?: string): Observable<Address[]> {
    const params = new Object();
    params['id'] = wallet.filename;
    params['num'] = num;
    if (password) {
      params['password'] = password;
    }

    return this.post('wallet/newAddress', params).pipe(
      map((response: PostWalletNewAddressResponse) => {
        const result: Address[] = [];
        response.addresses.forEach(value => {
          result.push({ address: value, coins: null, hours: null });
        });

        return result;
      }));
  }

  postWalletToggleEncryption(wallet: Wallet, password: string) {
    return this.post('wallet/' + (wallet.encrypted ? 'decrypt' : 'encrypt'), { id: wallet.filename, password });
  }

  get(url, params = null, options: any = {}, useV2 = false) {
    return this.http.get(this.getUrl(url, params, useV2), this.returnRequestOptions(options)).pipe(
      map((res: any) => res as any),
      catchError((error: any) => this.processConnectionError(error)));
  }

  getCsrf() {
    return this.get('csrf').pipe(map(response => response.csrf_token));
  }

  post(url, params = {}, options: any = {}, useV2 = false) {
    return this.getCsrf().pipe(first(), mergeMap(csrf => {
      options.csrf = csrf;

      if (useV2) {
        options.json = true;
      }

      return this.http.post(
        this.getUrl(url, null, useV2),
        options.json || useV2 ? JSON.stringify(params) : this.getQueryString(params),
        this.returnRequestOptions(options),
      ).pipe(
        map((res: any) => res as any),
        catchError((error: any) => this.processConnectionError(error)));
    }));
  }

  private returnRequestOptions(options) {
    const requestOptions: any = {};

    requestOptions.headers = new HttpHeaders();
    requestOptions.headers = requestOptions.headers.append('Content-Type', options.json ? 'application/json' : 'application/x-www-form-urlencoded');

    if (options.csrf) {
      requestOptions.headers = requestOptions.headers.append('X-CSRF-Token', options.csrf);
    }

    return requestOptions;
  }

  private getQueryString(parameters = null) {
    if (!parameters) {
      return '';
    }

    return Object.keys(parameters).reduce((array, key) => {
      array.push(key + '=' + encodeURIComponent(parameters[key]));

      return array;
    }, []).join('&');
  }

  private getUrl(url, options = null, useV2 = false) {
    if ((url as string).startsWith('/')) {
      url = (url as string).substr(1, (url as string).length - 1);
    }

    return this.url + (useV2 ? 'v2/' : 'v1/') + url + '?' + this.getQueryString(options);
  }

  processConnectionError(error: any, connectingToHwWalletDaemon = false): Observable<void> {
    if (error) {
      if (typeof error['_body'] === 'string') {

        return observableThrowError(error);
      }

      if (error.error && typeof error.error === 'string') {
        error['_body'] = error.error;

        return observableThrowError(error);
      } else if (error.error && error.error.error && error.error.error.message)  {
        error['_body'] = error.error.error.message;

        return observableThrowError(error);
      } else if (error.error && error.error.error && typeof error.error.error === 'string')  {
        error['_body'] = error.error.error;

        return observableThrowError(error);
      } else if (error.message) {
        error['_body'] = error.message;

        return observableThrowError(error);
      }
    }
    const err = Error(this.translate.instant(connectingToHwWalletDaemon ? 'hardware-wallet.errors.daemon-connection' : 'service.api.server-error'));
    err['_body'] = err.message;

    return observableThrowError(err);
  }
}
