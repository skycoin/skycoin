import { throwError as observableThrowError, Observable } from 'rxjs';
import { first, map, mergeMap, catchError } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';

import { environment } from '../../environments/environment';
import { processServiceError } from '../utils/errors';
import { OperationError } from '../utils/operation-error';

// IMPORTANT: AFTER MAKING MODIFICATIONS TO THIS INTERFACE YOU MUST ALSO
// MAKE APPROPIATE CHANGES TO THE createDefaultRequestOptions
// FUNCTION INSIDE ApiService.
/**
 * Options for configuring the requests to the node API.
 */
export interface NodeApiRequestOptions {
  /**
   * If true, the request will be sent to the API v2 and not to the v1.
   */
  useV2?: boolean;
  /**
   * If true, the data will be sent to the node encoded as JSON. This only makes sense for POST
   * request while using the v1, as GET request only send the data as params on the URL and the
   * data is always sent encoded as JSON when using the API v2.
   */
  sendDataAsJson?: boolean;
}

/**
 * Allows to make request to the node api with ease. Check the node API documentation for
 * information about the API endpoints.
 */
@Injectable()
export class ApiService {
  /**
   * URL for accessing the node API.
   */
  private url = environment.nodeUrl;

  constructor(
    private http: HttpClient,
  ) { }

  /**
   * Sends a GET request to the node API.
   * @param url URL to send the request to. You must omit the "http://x:x/api/vx/" part.
   * @param params Object with the key/value pairs to be sent to the node as part of
   * the querystring.
   * @param options Request options.
   */
  get(url: string, params: any = null, options: NodeApiRequestOptions = null): Observable<any> {
    if (!options) {
      options = this.createDefaultRequestOptions();
    } else {
      options = Object.assign(this.createDefaultRequestOptions(), options);
    }

    return this.http.get(this.getUrl(url, params, options.useV2), this.returnRequestOptions(options, null)).pipe(
      catchError((error: any) => this.processConnectionError(error)));
  }

  /**
   * Sends a POST request to the node API.
   * @param url URL to send the request to. You must omit the "http://x:x/api/vx/" part.
   * @param params Object with the key/value pairs to be sent to the node as
   * x-www-form-urlencoded or JSON, as defined in the options param.
   * @param options Request options.
   */
  post(url: string, params: any = null, options: NodeApiRequestOptions = null): Observable<any> {
    if (!options) {
      options = this.createDefaultRequestOptions();
    } else {
      options = Object.assign(this.createDefaultRequestOptions(), options);
    }

    return this.getCsrf().pipe(first(), mergeMap(csrf => {
      // V2 always needs the data to be sent encoded as JSON.
      if (options.useV2) {
        options.sendDataAsJson = true;
      }

      return this.http.post(
        this.getUrl(url, null, options.useV2),
        options.sendDataAsJson ? (params ? JSON.stringify(params) : '') : this.getQueryString(params),
        this.returnRequestOptions(options, csrf),
      ).pipe(
        catchError((error: any) => this.processConnectionError(error)));
    }));
  }

  /**
   * Creates a NodeApiRequestOptions instance with the default values.
   */
  private createDefaultRequestOptions(): NodeApiRequestOptions {
    return {
      useV2: false,
      sendDataAsJson: false,
    };
  }

  /**
   * Gets a csrf token from the node, to be able to make a post request to the node API.
   */
  private getCsrf(): Observable<string> {
    return this.get('csrf').pipe(map(response => response.csrf_token));
  }

  /**
   * Returns the options object requiered by HttpClient for sending a request.
   * @param options Options that will be used for making the request.
   * @param csrfToken Csrf token to be added on a header, for being able to make
   * POST requests.
   */
  private returnRequestOptions(options: NodeApiRequestOptions, csrfToken: string): any {
    const requestOptions: any = {};

    requestOptions.headers = new HttpHeaders();
    requestOptions.headers = requestOptions.headers.append('Content-Type', options.sendDataAsJson ? 'application/json' : 'application/x-www-form-urlencoded');

    if (csrfToken) {
      requestOptions.headers = requestOptions.headers.append('X-CSRF-Token', csrfToken);
    }

    return requestOptions;
  }

  /**
   * Encodes a list of params as a query string, for being used for sending data
   * in a request.
   * @param parameters Object with the key/value pairs that will be used for
   * creating the querystring.
   */
  private getQueryString(parameters: any = null): string {
    if (!parameters) {
      return '';
    }

    return Object.keys(parameters).reduce((array, key) => {
      array.push(key + '=' + encodeURIComponent(parameters[key]));

      return array;
    }, []).join('&');
  }

  /**
   * Get the complete URL needed for making a request to the node API.
   * @param url URL to send the request to, omitting the "http://x:x/api/vx/" part.
   * @param params Object with the key/value pairs to be sent to the node as part of
   * the querystring.
   * @param useV2 If the returned URL must point to the API v2 (true) or v1 (false).
   */
  private getUrl(url: string, params: any = null, useV2 = false): string {
    if (url.startsWith('/')) {
      url = url.substr(1, url.length - 1);
    }

    return this.url + (useV2 ? 'v2/' : 'v1/') + url + '?' + this.getQueryString(params);
  }

  /**
   * Takes an error returned by the node and converts it to an instance of OperationError.
   * @param error Error obtained while triying to connect to the node API.
   */
  private processConnectionError(error: any): Observable<OperationError> {
    return observableThrowError(processServiceError(error));
  }
}
