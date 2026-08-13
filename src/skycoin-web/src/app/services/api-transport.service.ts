import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { Observable, from, throwError } from 'rxjs';

/**
 * How a request to the node is actually sent.
 *
 * The wallet has always reached the node over same-origin HTTP, which is why
 * cmd/skycoin-web proxies /api: a browser cannot call a node directly, because
 * of CORS and because the node checks Host, Origin and Referer.
 *
 * A skywire visor running on the same page can make the request itself, outside
 * the browser's fetch stack, so those constraints do not apply. It publishes
 * skywireVisor.fetchDmsg for exactly this — its own comment names "form-encoded
 * skycoin API POSTs (balance, txns)" as the case it handles.
 *
 * Which one to use is not a setting. A dmsg address cannot be fetched by the
 * browser at all, and an http(s) one does not need the visor, so the target
 * decides. That keeps a wallet with no visor behaving exactly as it does now.
 */

/** Options a caller can attach to a request, independent of how it is sent. */
export interface ApiRequestOptions {
  params?: { [key: string]: string };
  headers?: { [key: string]: string };
}

/** The visor's surface, as much of it as this file uses. */
interface SkywireVisor {
  fetchDmsg(
    pkHost: string,
    method: string,
    path: string,
    body?: string | null,
    headers?: { [key: string]: string },
  ): Promise<{ status: number; headers?: { [key: string]: string }; body?: string }>;

  /**
   * Reports what the visor has running. `dmsg` is the field that matters here:
   * it is true once the dmsg client exists, which is the exact condition
   * fetchDmsg tests before doing anything — it answers "not booted; call boot()
   * first" otherwise. `dmsg_connected` is a stronger statement, that the client
   * also has a session.
   */
  status?(): { booted?: boolean; dmsg?: boolean; dmsg_connected?: boolean };
}

@Injectable()
export class ApiTransportService {
  constructor(private http: HttpClient) {}

  /**
   * Reports whether a visor is on the page.
   *
   * Detected through the global rather than by importing anything from skywire:
   * this project does not depend on it, and the visor is a drop-in precisely
   * because the contract is the object it publishes.
   *
   * Being present is not the same as being able to carry a request — see
   * visorReady. A visor exists on the page from the moment its wasm loads, well
   * before anyone has called boot().
   */
  get visorAvailable(): boolean {
    const visor = (window as any).skywireVisor as SkywireVisor | undefined;

    return !!visor && typeof visor.fetchDmsg === 'function';
  }

  /**
   * Reports whether the visor has booted far enough to carry a request.
   *
   * fetchDmsg refuses with "not booted; call boot() first" until its dmsg
   * client exists, so this asks status() for the same thing rather than finding
   * out by failing. A visor that publishes no status() is taken at its word:
   * older builds had no way to ask, and refusing them would be worse than
   * letting the request report its own failure.
   */
  get visorReady(): boolean {
    const visor = (window as any).skywireVisor as SkywireVisor | undefined;
    if (!visor || typeof visor.fetchDmsg !== 'function') {
      return false;
    }
    if (typeof visor.status !== 'function') {
      return true;
    }

    try {
      return !!visor.status().dmsg;
    } catch {
      // A visor whose status cannot be read is not one to send a request to.
      return false;
    }
  }

  /**
   * Reports whether a url names a host only the visor can reach.
   *
   * A dmsg address is either <pk>.dmsg or one of the resolver aliases the visor
   * understands. Anything else is ordinary http(s).
   */
  isDmsgTarget(url: string): boolean {
    const host = this.hostOf(url);

    return host.endsWith('.dmsg');
  }

  request(method: string, url: string, body: string | null, options: ApiRequestOptions = {}): Observable<any> {
    if (this.isDmsgTarget(url)) {
      if (!this.visorAvailable) {
        return throwError(() => new Error(
          `${this.hostOf(url)} can only be reached through a skywire visor, and none is running on this page`));
      }
      if (!this.visorReady) {
        // Distinguished from the above because the two need different things
        // from whoever sees them: one is a missing visor, the other a visor
        // that nobody has booted yet.
        return throwError(() => new Error(
          `the skywire visor on this page has not booted yet, so ${this.hostOf(url)} cannot be reached`));
      }

      return this.requestThroughVisor(method, url, body, options);
    }

    return this.requestOverHttp(method, url, body, options);
  }

  private requestOverHttp(method: string, url: string, body: string | null, options: ApiRequestOptions): Observable<any> {
    const request = {
      params: this.toParams(options.params),
      headers: this.toHeaders(options.headers),
    };

    if (method === 'GET') {
      return this.http.get(url, request);
    }
    if (method === 'DELETE') {
      return this.http.delete(url, request);
    }

    return this.http.post(url, body, request);
  }

  private requestThroughVisor(method: string, url: string, body: string | null, options: ApiRequestOptions): Observable<any> {
    const visor = (window as any).skywireVisor as SkywireVisor;
    const target = new URL(url);
    const query = this.queryString(options.params);
    const path = target.pathname + (query ? '?' + query : '');

    return from(visor.fetchDmsg(target.host, method, path, body, options.headers).then(response => {
      if (response.status < 200 || response.status >= 300) {
        // Shaped like an HttpErrorResponse so ApiService reads it the same way
        // whichever transport produced it.
        throw { status: response.status, error: this.parse(response.body) };
      }

      return this.parse(response.body);
    }));
  }

  /** The node answers JSON; anything else is handed back untouched. */
  private parse(body: string | undefined): any {
    if (!body) {
      return null;
    }

    try {
      return JSON.parse(body);
    } catch {
      return body;
    }
  }

  private hostOf(url: string): string {
    try {
      return new URL(url).hostname;
    } catch {
      return '';
    }
  }

  private queryString(params: { [key: string]: string } | undefined): string {
    if (!params) {
      return '';
    }

    return Object.keys(params)
      .map(key => key + '=' + encodeURIComponent(params[key]))
      .join('&');
  }

  private toParams(params: { [key: string]: string } | undefined): HttpParams {
    let result = new HttpParams();
    if (params) {
      Object.keys(params).forEach(key => result = result.set(key, params[key]));
    }

    return result;
  }

  private toHeaders(headers: { [key: string]: string } | undefined): HttpHeaders {
    let result = new HttpHeaders();
    if (headers) {
      Object.keys(headers).forEach(key => result = result.append(key, headers[key]));
    }

    return result;
  }
}
