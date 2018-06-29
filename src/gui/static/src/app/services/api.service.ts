import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { environment } from '../../environments/environment';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import {
  Address, GetWalletsResponseEntry, GetWalletsResponseWallet, PostWalletNewAddressResponse, Transaction, Version,
  Wallet
} from '../app.datatypes';

@Injectable()
export class ApiService {

  private url = environment.nodeUrl;

  constructor(
    private http: Http,
  ) { }

  getExplorerAddress(address: Address): Observable<Transaction[]> {
    return this.get('explorer/address', {address: address.address})
      .map(transactions => transactions.map(transaction => ({
        addresses: [],
        balance: 0,
        block: transaction.status.block_seq,
        confirmed: transaction.status.confirmed,
        timestamp: transaction.timestamp,
        txid: transaction.txid,
        inputs: transaction.inputs,
        outputs: transaction.outputs,
      })));
  }

  getVersion(): Observable<Version> {
    return this.get('version');
  }

  getWalletNewSeed(): Observable<string> {
    return this.get('wallet/newSeed')
      .map(response => response.seed);
  }

  getWallets(): Observable<Wallet[]> {
    return this.get('wallets')
      .map((response: GetWalletsResponseWallet[]) => {
        const wallets: Wallet[] = [];
        response.forEach(wallet => {
          wallets.push(<Wallet>{
            label: wallet.meta.label,
            filename: wallet.meta.filename,
            coins: null,
            hours: null,
            addresses: wallet.entries.map((entry: GetWalletsResponseEntry) => {
              return {
                address: entry.address,
                coins: null,
                hours: null,
              };
            }),
            encrypted: wallet.meta.encrypted,
            useEmulatorWallet: wallet.meta.useEmulatorWallet,
            useHardwareWallet: wallet.meta.useHardwareWallet,
          });
        });
        return wallets;
      });
  }

  getWalletSeed(wallet: Wallet, password: string): Observable<string> {
    return this.post('wallet/seed', { id: wallet.filename, password })
      .map(response => response.seed);
  }

  postWalletCreate(label: string, seed: string, scan: number, password: string, useHardwareWallet: boolean, useEmulatorWallet: boolean): Observable<Wallet> {
    const params = { label, seed, scan };

    if (password) {
      params['password'] = password;
      params['encrypt'] = true;
    }
    if (useHardwareWallet) {
      params['useHardwareWallet'] = useHardwareWallet;
    }
    if (useEmulatorWallet) {
      params['useEmulatorWallet'] = useEmulatorWallet;
    }

    return this.post('wallet/create', params)
      .map(response => ({
          label: response.meta.label,
          filename: response.meta.filename,
          coins: null,
          hours: null,
          addresses: response.entries.map(entry => ({ address: entry.address, coins: null, hours: null })),
          encrypted: response.meta.encrypted,
          useHardwareWallet: response.meta.useHardwareWallet,
          useEmulatorWallet: response.meta.useEmulatorWallet,
        }));
  }

  postWalletNewAddress(wallet: Wallet, password?: string): Observable<Address> {
    return this.post('wallet/newAddress', { id: wallet.filename, password })
      .map((response: PostWalletNewAddressResponse) => ({ address: response.addresses[0], coins: null, hours: null }));
  }

  postWalletToggleEncryption(wallet: Wallet, password: string) {
    return this.post('wallet/' + (wallet.encrypted ? 'decrypt' : 'encrypt'), { id: wallet.filename, password });
  }

  get(url, params = null, options = {}) {
    return this.http.get(this.getUrl(url, params), this.returnRequestOptions(options))
      .map((res: any) => res.json())
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  getCsrf() {
    return this.get('csrf').map(response => response.csrf_token);
  }

  post(url, params = {}, options: any = {}) {
    return this.getCsrf().first().flatMap(csrf => {
      options.csrf = csrf;
      return this.http.post(this.getUrl(url), this.getQueryString(params), this.returnRequestOptions(options))
        .map((res: any) => res.json())
        .catch((error: any) => Observable.throw(error || 'Server error'));
    });
  }

  returnRequestOptions(additionalOptions) {
    const options = new RequestOptions();

    options.headers = this.getHeaders();

    if (additionalOptions.csrf) {
      options.headers.append('X-CSRF-Token', additionalOptions.csrf);
    }

    return options;
  }

  private getHeaders() {
    const headers = new Headers();
    headers.append('Content-Type', 'application/x-www-form-urlencoded');
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

  private getUrl(url, options = null) {
    return this.url + url + '?' + this.getQueryString(options);
  }
}
