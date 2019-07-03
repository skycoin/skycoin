import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { environment } from '../../environments/environment';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
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
    private http: Http,
    private translate: TranslateService,
  ) { }

  getTransactions(addresses: Address[]): Observable<NormalTransaction[]> {
    const formattedAddresses = addresses.map(a => a.address).join(',');

    return this.post('transactions', {addrs: formattedAddresses, verbose: true})
      .map(transactions => transactions.map(transaction => ({
        addresses: [],
        balance: new BigNumber(0),
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.txn.timestamp,
        txid: transaction.txn.txid,
        inputs: transaction.txn.inputs,
        outputs: transaction.txn.outputs,
      })));
  }

  getVersion(): Observable<Version> {
    return this.get('version');
  }

  generateSeed(entropy: number = 128): Observable<string> {
    return this.get('wallet/newSeed', { entropy }).map(response => response.seed);
  }

  getHealth() {
    return this.get('health');
  }

  getWallets(): Observable<Wallet[]> {
    return this.get('wallets')
      .map((response: GetWalletsResponseWallet[]) => {
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
      });
  }

  getWalletSeed(wallet: Wallet, password: string): Observable<string> {
    return this.post('wallet/seed', { id: wallet.filename, password })
      .map(response => response.seed);
  }

  postWalletCreate(label: string, seed: string, scan: number, password: string): Observable<Wallet> {
    const params = { label, seed, scan };

    if (password) {
      params['password'] = password;
      params['encrypt'] = true;
    }

    return this.post('wallet/create', params)
      .map(response => ({
          label: response.meta.label,
          filename: response.meta.filename,
          coins: null,
          hours: null,
          addresses: response.entries.map(entry => ({ address: entry.address, coins: null, hours: null, confirmed: true })),
          encrypted: response.meta.encrypted,
        }));
  }

  postWalletNewAddress(wallet: Wallet, num: number, password?: string): Observable<Address[]> {
    const params = new Object();
    params['id'] = wallet.filename;
    params['num'] = num;
    if (password) {
      params['password'] = password;
    }

    return this.post('wallet/newAddress', params)
      .map((response: PostWalletNewAddressResponse) => {
        const result: Address[] = [];
        response.addresses.forEach(value => {
          result.push({ address: value, coins: null, hours: null });
        });

        return result;
      });
  }

  postWalletToggleEncryption(wallet: Wallet, password: string) {
    return this.post('wallet/' + (wallet.encrypted ? 'decrypt' : 'encrypt'), { id: wallet.filename, password });
  }

  get(url, params = null, options: any = {}, useV2 = false) {
    return this.http.get(this.getUrl(url, params, useV2), this.returnRequestOptions(options))
      .map((res: any) => res.json())
      .catch((error: any) => this.processConnectionError(error));
  }

  getCsrf() {
    return this.get('csrf').map(response => response.csrf_token);
  }

  post(url, params = {}, options: any = {}, useV2 = false) {
    return this.getCsrf().first().flatMap(csrf => {
      options.csrf = csrf;

      if (useV2) {
        options.json = true;
      }

      return this.http.post(
        this.getUrl(url, null, useV2),
        options.json || useV2 ? JSON.stringify(params) : this.getQueryString(params),
        this.returnRequestOptions(options),
      )
        .map((res: any) => res.json())
        .catch((error: any) => this.processConnectionError(error));
    });
  }

  returnRequestOptions(additionalOptions) {
    const options = new RequestOptions();

    options.headers = this.getHeaders(additionalOptions);

    if (additionalOptions.csrf) {
      options.headers.append('X-CSRF-Token', additionalOptions.csrf);
    }

    return options;
  }

  private getHeaders(options) {
    const headers = new Headers();
    headers.append('Content-Type', options.json ? 'application/json' : 'application/x-www-form-urlencoded');

    return headers;
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

        return Observable.throw(error);
      }

      if (error.error && typeof error.error === 'string') {
        error['_body'] = error.error;

        return Observable.throw(error);
      } else if (error.message) {
        error['_body'] = error.message;

        return Observable.throw(error);
      }
    }
    const err = Error(this.translate.instant(connectingToHwWalletDaemon ? 'hardware-wallet.errors.daemon-connection' : 'service.api.server-error'));
    err['_body'] = err.message;

    return Observable.throw(err);
  }
}
