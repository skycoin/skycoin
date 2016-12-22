import {root} from '../../util/root';
import {tryCatch} from '../../util/tryCatch';
import {errorObject} from '../../util/errorObject';
import {Observable} from '../../Observable';
import {Subscriber} from '../../Subscriber';
import {TeardownLogic} from '../../Subscription';

export interface AjaxRequest {
  url?: string;
  body?: any;
  user?: string;
  async?: boolean;
  method: string;
  headers?: Object;
  timeout?: number;
  password?: string;
  hasContent?: boolean;
  crossDomain?: boolean;
  createXHR?: () => XMLHttpRequest;
  progressSubscriber?: Subscriber<any>;
  resultSelector?: <T>(response: AjaxResponse) => T;
  responseType?: string;
}

function createXHRDefault(): XMLHttpRequest {
  let xhr = new root.XMLHttpRequest();
  if (this.crossDomain) {
    if ('withCredentials' in xhr) {
      xhr.withCredentials = true;
      return xhr;
    } else if (!!root.XDomainRequest) {
      return new root.XDomainRequest();
    } else {
      throw new Error('CORS is not supported by your browser');
    }
  } else {
    return xhr;
  }
}

export interface AjaxCreationMethod {
  <T>(urlOrRequest: string | AjaxRequest): Observable<T>;
  get<T>(url: string, resultSelector?: (response: AjaxResponse) => T, headers?: Object): Observable<T>;
  post<T>(url: string, body?: any, headers?: Object): Observable<T>;
  put<T>(url: string, body?: any, headers?: Object): Observable<T>;
  delete<T>(url: string, headers?: Object): Observable<T>;
  getJSON<T, R>(url: string, resultSelector?: (data: T) => R, headers?: Object): Observable<R>;
}

function defaultGetResultSelector<T>(response: AjaxResponse): T {
  return response.response;
}

export function ajaxGet<T>(url: string, resultSelector: (response: AjaxResponse) => T = defaultGetResultSelector, headers: Object = null) {
  return new AjaxObservable<T>({ method: 'GET', url, resultSelector, headers });
};

export function ajaxPost<T>(url: string, body?: any, headers?: Object): Observable<T> {
  return new AjaxObservable<T>({ method: 'POST', url, body, headers });
};

export function ajaxDelete<T>(url: string, headers?: Object): Observable<T> {
  return new AjaxObservable<T>({ method: 'DELETE', url, headers });
};

export function ajaxPut<T>(url: string, body?: any, headers?: Object): Observable<T> {
  return new AjaxObservable<T>({ method: 'PUT', url, body, headers });
};

export function ajaxGetJSON<T, R>(url: string, resultSelector?: (data: T) => R, headers?: Object): Observable<R> {
  const finalResultSelector = resultSelector ? (res: AjaxResponse) => resultSelector(res.response) : (res: AjaxResponse) => res.response;
  return new AjaxObservable<R>({ method: 'GET', url, responseType: 'json', resultSelector: finalResultSelector, headers });
};

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @extends {Ignored}
 * @hide true
 */
export class AjaxObservable<T> extends Observable<T> {
  /**
   * Creates an observable for an Ajax request with either a request object with
   * url, headers, etc or a string for a URL.
   *
   * @example
   * source = Rx.Observable.ajax('/products');
   * source = Rx.Observable.ajax( url: 'products', method: 'GET' });
   *
   * @param {string|Object} request Can be one of the following:
   *   A string of the URL to make the Ajax call.
   *   An object with the following properties
   *   - url: URL of the request
   *   - body: The body of the request
   *   - method: Method of the request, such as GET, POST, PUT, PATCH, DELETE
   *   - async: Whether the request is async
   *   - headers: Optional headers
   *   - crossDomain: true if a cross domain request, else false
   *   - createXHR: a function to override if you need to use an alternate
   *   XMLHttpRequest implementation.
   *   - resultSelector: a function to use to alter the output value type of
   *   the Observable. Gets {@link AjaxResponse} as an argument.
   * @return {Observable} An observable sequence containing the XMLHttpRequest.
   * @static true
   * @name ajax
   * @owner Observable
  */
  static _create_stub(): void { return null; }

  static create: AjaxCreationMethod = (() => {
    const create: any = (urlOrRequest: string | AjaxRequest) => {
      return new AjaxObservable(urlOrRequest);
    };

    create.get = ajaxGet;
    create.post = ajaxPost;
    create.delete = ajaxDelete;
    create.put = ajaxPut;
    create.getJSON = ajaxGetJSON;

    return <AjaxCreationMethod>create;
  })();

  private request: AjaxRequest;

  constructor(urlOrRequest: string | AjaxRequest) {
    super();

    const request: AjaxRequest = {
      async: true,
      createXHR: createXHRDefault,
      crossDomain: false,
      headers: {},
      method: 'GET',
      responseType: 'json',
      timeout: 0
    };

    if (typeof urlOrRequest === 'string') {
      request.url = urlOrRequest;
    } else {
      for (const prop in urlOrRequest) {
        if (urlOrRequest.hasOwnProperty(prop)) {
          request[prop] = urlOrRequest[prop];
        }
      }
    }

    this.request = request;
  }

  protected _subscribe(subscriber: Subscriber<T>): TeardownLogic {
    return new AjaxSubscriber(subscriber, this.request);
  }
}

/**
 * We need this JSDoc comment for affecting ESDoc.
 * @ignore
 * @extends {Ignored}
 */
export class AjaxSubscriber<T> extends Subscriber<Event> {
  private xhr: XMLHttpRequest;
  private resultSelector: (response: AjaxResponse) => T;
  private done: boolean = false;

  constructor(destination: Subscriber<T>, public request: AjaxRequest) {
    super(destination);

    const headers = request.headers = request.headers || {};

    // force CORS if requested
    if (!request.crossDomain && !headers['X-Requested-With']) {
      headers['X-Requested-With'] = 'XMLHttpRequest';
    }

    // ensure content type is set
    if (!('Content-Type' in headers)) {
      headers['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8';
    }

    // properly serialize body
    request.body = this.serializeBody(request.body, request.headers['Content-Type']);

    this.resultSelector = request.resultSelector;
    this.send();
  }

  next(e: Event): void {
    this.done = true;
    const { resultSelector, xhr, request, destination } = this;
    const response = new AjaxResponse(e, xhr, request);

    if (resultSelector) {
      const result = tryCatch(resultSelector)(response);
      if (result === errorObject) {
        this.error(errorObject.e);
      } else {
        destination.next(result);
      }
    } else {
      destination.next(response);
    }
  }

  private send(): XMLHttpRequest {
    const {
      request,
      request: { user, method, url, async, password, headers, body }
    } = this;
    const createXHR = request.createXHR;
    const xhr: XMLHttpRequest = tryCatch(createXHR).call(request);

    if (<any>xhr === errorObject) {
      this.error(errorObject.e);
    } else {
      this.xhr = xhr;

      // open XHR first
      let result: any;
      if (user) {
        result = tryCatch(xhr.open).call(xhr, method, url, async, user, password);
      } else {
        result = tryCatch(xhr.open).call(xhr, method, url, async);
      }

      if (result === errorObject) {
        this.error(errorObject.e);
        return;
      }

      // timeout and responseType can be set once the XHR is open
      xhr.timeout = request.timeout;
      xhr.responseType = request.responseType;

      // set headers
      this.setHeaders(xhr, headers);

      // now set up the events
      this.setupEvents(xhr, request);

      // finally send the request
      if (body) {
        xhr.send(body);
      } else {
        xhr.send();
      }
    }
  }

  private serializeBody(body: any, contentType: string) {
    if (!body || typeof body === 'string') {
      return body;
    } else if (root.FormData && body instanceof root.FormData) {
      return body;
    }

    const splitIndex = contentType.indexOf(';');
    if (splitIndex !== -1) {
      contentType = contentType.substring(0, splitIndex);
    }

    switch (contentType) {
      case 'application/x-www-form-urlencoded':
        return Object.keys(body).map(key => `${encodeURI(key)}=${encodeURI(body[key])}`).join('&');
      case 'application/json':
        return JSON.stringify(body);
    }
  }

  private setHeaders(xhr: XMLHttpRequest, headers: Object) {
    for (let key in headers) {
      if (headers.hasOwnProperty(key)) {
        xhr.setRequestHeader(key, headers[key]);
      }
    }
  }

  private setupEvents(xhr: XMLHttpRequest, request: AjaxRequest) {
    const progressSubscriber = request.progressSubscriber;

    xhr.ontimeout = function xhrTimeout(e) {
      const {subscriber, progressSubscriber, request } = (<any>xhrTimeout);
      if (progressSubscriber) {
        progressSubscriber.error(e);
      }
      subscriber.error(new AjaxTimeoutError(this, request)); //TODO: Make betterer.
    };
    (<any>xhr.ontimeout).request = request;
    (<any>xhr.ontimeout).subscriber = this;
    (<any>xhr.ontimeout).progressSubscriber = progressSubscriber;

    if (xhr.upload && 'withCredentials' in xhr && root.XDomainRequest) {
      if (progressSubscriber) {
        xhr.onprogress = function xhrProgress(e) {
          const { progressSubscriber } = (<any>xhrProgress);
          progressSubscriber.next(e);
        };
        (<any>xhr.onprogress).progressSubscriber = progressSubscriber;
      }

      xhr.onerror = function xhrError(e) {
        const { progressSubscriber, subscriber, request } = (<any>xhrError);
        if (progressSubscriber) {
          progressSubscriber.error(e);
        }
        subscriber.error(new AjaxError('ajax error', this, request));
      };
      (<any>xhr.onerror).request = request;
      (<any>xhr.onerror).subscriber = this;
      (<any>xhr.onerror).progressSubscriber = progressSubscriber;
    }

    xhr.onreadystatechange = function xhrReadyStateChange(e) {
      const { subscriber, progressSubscriber, request } = (<any>xhrReadyStateChange);
      if (this.readyState === 4) {
        // normalize IE9 bug (http://bugs.jquery.com/ticket/1450)
        let status: number = this.status === 1223 ? 204 : this.status;
        let response: any = (this.responseType === 'text' ?  (
          this.response || this.responseText) : this.response);

        // fix status code when it is 0 (0 status is undocumented).
        // Occurs when accessing file resources or on Android 4.1 stock browser
        // while retrieving files from application cache.
        if (status === 0) {
          status = response ? 200 : 0;
        }

        if (200 <= status && status < 300) {
          if (progressSubscriber) {
            progressSubscriber.complete();
          }
          subscriber.next(e);
          subscriber.complete();
        } else {
          if (progressSubscriber) {
            progressSubscriber.error(e);
          }
          subscriber.error(new AjaxError('ajax error ' + status, this, request));
        }
      }
    };
    (<any>xhr.onreadystatechange).subscriber = this;
    (<any>xhr.onreadystatechange).progressSubscriber = progressSubscriber;
    (<any>xhr.onreadystatechange).request = request;
  }

  unsubscribe() {
    const { done, xhr } = this;
    if (!done && xhr && xhr.readyState !== 4) {
      xhr.abort();
    }
    super.unsubscribe();
  }
}

/**
 * A normalized AJAX response.
 *
 * @see {@link ajax}
 *
 * @class AjaxResponse
 */
export class AjaxResponse {
  /** @type {number} The HTTP status code */
  status: number;

  /** @type {string|ArrayBuffer|Document|object|any} The response data */
  response: any;

  /** @type {string} The raw responseText */
  responseText: string;

  /** @type {string} The responseType (e.g. 'json', 'arraybuffer', or 'xml') */
  responseType: string;

  constructor(public originalEvent: Event, public xhr: XMLHttpRequest, public request: AjaxRequest) {
    this.status = xhr.status;
    this.responseType = xhr.responseType || request.responseType;

    switch (this.responseType) {
      case 'json':
        if ('response' in xhr) {
          //IE does not support json as responseType, parse it internally
          this.response = xhr.responseType ? xhr.response : JSON.parse(xhr.response || xhr.responseText || '');
        } else {
          this.response = JSON.parse(xhr.responseText || '');
        }
        break;
      case 'xml':
        this.response = xhr.responseXML;
        break;
      case 'text':
      default:
        this.response = ('response' in xhr) ? xhr.response : xhr.responseText;
        break;
    }
  }
}

/**
 * A normalized AJAX error.
 *
 * @see {@link ajax}
 *
 * @class AjaxError
 */
export class AjaxError extends Error {
  /** @type {XMLHttpRequest} The XHR instance associated with the error */
  xhr: XMLHttpRequest;

  /** @type {AjaxRequest} The AjaxRequest associated with the error */
  request: AjaxRequest;

  /** @type {number} The HTTP status code */
  status: number;

  constructor(message: string, xhr: XMLHttpRequest, request: AjaxRequest) {
    super(message);
    this.message = message;
    this.xhr = xhr;
    this.request = request;
    this.status = xhr.status;
  }
}

/**
 * @see {@link ajax}
 *
 * @class AjaxTimeoutError
 */
export class AjaxTimeoutError extends AjaxError {
  constructor(xhr: XMLHttpRequest, request: AjaxRequest) {
    super('ajax timeout', xhr, request);
  }
}
