import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import { catchError, map, mergeMap, timeout } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { CoinService } from './coin.service';
import { BaseCoin } from '../coins/basecoin';
import { parseResponseMessage } from '../utils/errors';
import { ApiRequestOptions, ApiTransportService } from './api-transport.service';

@Injectable()
export class ApiService {

  /** How long a request waits for the node. */
  requestTimeoutMs = 120000;

  private url!: string;

  constructor(private transport: ApiTransportService,
              private translate: TranslateService,
              private coinService: CoinService) {
    this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        if (!coin) {
          return;
        }
        const customUrl = (coinService.customNodeUrls as any)[coin.id.toString()];
        this.url = customUrl ? customUrl : coin.nodeUrl;
        if (this.url.endsWith('/')) {
          this.url = this.url.substring(0, this.url.length - 1);
        }
        this.url += '/api/';
      });
  }

  get(url: any, params: any = null, options = {}): Observable<any> {
    return this.request('GET', this.getUrl(url), null, this.getRequestOptions(options, params)).pipe(
      catchError((error: any) => this.getErrorMessage(error)));
  }

  delete(url: any, params = null, options = {}): Observable<any> {
    return this.withCsrf(options).pipe(
      mergeMap(opts => this.request('DELETE', this.getUrl(url), null, this.getRequestOptions(opts, params))),
      catchError((error: any) => this.getErrorMessage(error)));
  }

  post(url: any, body = {}, options: any = {}, useV2 = false): Observable<any> {
    if (useV2) {
      options.json = true;
    }

    return this.withCsrf(options).pipe(
      mergeMap(opts => this.request(
        'POST',
        this.getUrl(url, useV2),
        opts.json ? JSON.stringify(body) : this.getQueryString(body),
        this.getRequestOptions(opts),
      )),
      catchError((error: any) => this.getErrorMessage(error)));
  }

  /** A request that fails when the node has not answered within requestTimeoutMs. */
  private request(method: string, url: string, body: string | null, options: ApiRequestOptions): Observable<any> {
    return this.transport.request(method, url, body, options).pipe(
      timeout({
        first: this.requestTimeoutMs,
        with: () => throwError(() => new Error(this.translate.instant('service.api.timeout'))),
      }));
  }

  /**
   * Adds the node's CSRF token to the options of a request that changes state. A node
   * checks it on every POST and DELETE; a node with CSRF off, or a server in front of
   * one, answers without a token and the request goes as it is.
   */
  private withCsrf(options: any): Observable<any> {
    if (options.csrf) {
      return of(options);
    }

    return this.request('GET', this.getUrl('csrf'), null, this.getRequestOptions({})).pipe(
      map((response: any) => response && response.csrf_token ? { ...options, csrf: response.csrf_token } : options),
      catchError(() => of(options)),
    );
  }

  private getQueryString(parameters: any = null) {
    if (!parameters) {
      return '';
    }

    return Object.keys(parameters).reduce((array, key) => {
      array.push(key + '=' + encodeURIComponent(parameters[key]));

      return array;
    }, [] as string[]).join('&');
  }

  /**
   * Builds the parts of a request that do not depend on how it is sent. The
   * transport turns these into HttpParams and HttpHeaders, or into the
   * arguments the visor's fetch takes.
   */
  private getRequestOptions(additionalOptions: any, parameters: any = null): ApiRequestOptions {
    const headers: { [key: string]: string } = {
      'Content-Type': additionalOptions.json ? 'application/json' : 'application/x-www-form-urlencoded',
    };

    if (additionalOptions.csrf) {
      headers['X-CSRF-Token'] = additionalOptions.csrf;
    }

    const params: { [key: string]: string } = {};
    if (parameters) {
      Object.keys(parameters).forEach((key: string) => params[key] = parameters[key]);
    }

    return { params: params, headers: headers };
  }

  private getUrl(url: string, useV2 = false): string {
    if (url.startsWith('/')) {
      url = url.substring(1);
    }

    return this.url + (useV2 ? 'v2/' : 'v1/') + url;
  }

  private getErrorMessage(error: any): Observable<string> {
    if (error) {
      if (typeof error['_body'] === 'string') {
        return throwError(() => new Error(parseResponseMessage(error)));
      }

      if (error.error && typeof error.error === 'string') {
        return throwError(() => new Error(parseResponseMessage(error.error.trim())));
      } else if (error.error && error.error.error && error.error.error.message) {
        return throwError(() => new Error(parseResponseMessage(error.error.error.message.trim())));
      } else if (error.message) {
        return throwError(() => new Error(parseResponseMessage(error.message.trim())));
      }
    }

    return this.translate.get('service.api.server-error').pipe(
      mergeMap(message => throwError(() => new Error(message))));
  }
}
