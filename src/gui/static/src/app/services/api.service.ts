import { throwError as observableThrowError, Observable } from 'rxjs';
import { first, map, mergeMap, catchError } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { BigNumber } from 'bignumber.js';

import { NormalTransaction, Version } from '../app.datatypes';
import { processServiceError } from '../utils/errors';
import { AddressBase } from './wallet-operations/wallet-objects';

@Injectable()
export class ApiService {
  private url = environment.nodeUrl;

  constructor(
    private http: HttpClient,
  ) { }

  getTransactions(addresses: AddressBase[]): Observable<NormalTransaction[]> {
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

  private processConnectionError(error: any): Observable<void> {
    return observableThrowError(processServiceError(error));
  }
}
