import { Injectable } from '@angular/core';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';

@Injectable()
export class ApiService {

  private url = 'http://127.0.0.1:6420/'; // production
  // private url = '/api/'; // test

  constructor(private http: Http) { }

  get(url, options = null) {
    return this.http.get(this.getUrl(url, options), this.getHeaders())
      .map((res: any) => res.json())
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  post(url, options = {}) {
    return this.http.post(this.getUrl(url), this.getQueryString(options), this.returnRequestOptions())
      .map((res: any) => res.json())
      .catch((error: any) => Observable.throw(error || 'Server error'));
  }

  private getHeaders() {
    const headers = new Headers();
    headers.append('Content-Type', 'application/x-www-form-urlencoded');
    return headers;
  }

  returnRequestOptions() {
    const options = new RequestOptions();

    options.headers = this.getHeaders();

    return options;
  }

  private getQueryString(parameters = null) {
    if (!parameters) {
      return '';
    }

    return Object.keys(parameters).reduce((array,key) => {
      array.push(key + '=' + encodeURIComponent(parameters[key]));
      return array;
    }, []).join('&');
  }

  private getUrl(url, options = null) {
    return this.url + url + '?' + this.getQueryString(options);
  }
}
