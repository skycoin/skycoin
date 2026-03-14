import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, mergeMap } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';

import { CoinService } from './coin.service';
import { BaseCoin } from '../coins/basecoin';
import { parseResponseMessage } from '../utils/errors';

@Injectable()
export class ApiService {

  private url: string;

  constructor(private http: HttpClient,
              private translate: TranslateService,
              private coinService: CoinService) {
    this.coinService.currentCoin
      .subscribe((coin: BaseCoin) => {
        if (!coin) {
          return;
        }
        const customUrl = coinService.customNodeUrls[coin.id.toString()];
        this.url = customUrl ? customUrl : coin.nodeUrl;
        if (this.url.endsWith('/')) {
          this.url = this.url.substring(0, this.url.length - 1);
        }
        this.url += '/api/';
      });
  }

  get(url, params = null, options = {}): Observable<any> {
    return this.http.get(this.getUrl(url), this.getRequestOptions(options, params)).pipe(
      catchError((error: any) => this.getErrorMessage(error)));
  }

  delete(url, params = null, options = {}): Observable<any> {
    return this.http.delete(this.getUrl(url), this.getRequestOptions(options, params)).pipe(
      catchError((error: any) => this.getErrorMessage(error)));
  }

  post(url, body = {}, options: any = {}, useV2 = false): Observable<any> {
    if (useV2) {
      options.json = true;
    }

    return this.http.post(
      this.getUrl(url, useV2),
      options.json ? JSON.stringify(body) : this.getQueryString(body),
      this.getRequestOptions(options)
    ).pipe(catchError((error: any) => this.getErrorMessage(error)));
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

  private getRequestOptions(additionalOptions: any, parameters: any = null): any {
    const options: any = {};
    options.params = this.getQueryStringParams(parameters);
    options.headers = new HttpHeaders();

    options.headers = options.headers.append('Content-Type', additionalOptions.json ? 'application/json' : 'application/x-www-form-urlencoded');

    if (additionalOptions.csrf) {
      options.headers = options.headers.append('X-CSRF-Token', additionalOptions.csrf);
    }

    return options;
  }

  private getQueryStringParams(parameters: any): HttpParams {
    let params = new HttpParams();

    if (parameters) {
      Object.keys(parameters).forEach((key: string) => params = params.set(key, parameters[key]));
    }

    return params;
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
