import { TestBed } from '@angular/core/testing';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';

import { ApiTransportService } from './api-transport.service';

/**
 * The seam's job is to decide how a request is sent, and to make both ways look
 * the same to ApiService. The decision is the part worth pinning: a wallet with
 * no visor has to keep behaving exactly as it did.
 */
describe('ApiTransportService', () => {
  let transport: ApiTransportService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ApiTransportService,
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
      ],
    });

    transport = TestBed.inject(ApiTransportService);
    httpMock = TestBed.inject(HttpTestingController);
    delete (window as any).skywireVisor;
  });

  afterEach(() => {
    delete (window as any).skywireVisor;
  });

  describe('choosing a transport', () => {
    it('treats an http node as ordinary http', () => {
      expect(transport.isDmsgTarget('http://127.0.0.1:6420/api/v1/health')).toBe(false);
    });

    it('recognises a dmsg host', () => {
      expect(transport.isDmsgTarget('http://02a1b2c3.dmsg/api/v1/health')).toBe(true);
    });

    it('reports no visor when the page has none', () => {
      expect(transport.visorAvailable).toBe(false);
    });

    it('reports a visor only when it can actually carry a request', () => {
      (window as any).skywireVisor = {};
      expect(transport.visorAvailable).toBe(false);

      (window as any).skywireVisor = { fetchDmsg: () => Promise.resolve({ status: 200 }) };
      expect(transport.visorAvailable).toBe(true);
    });
  });

  // A visor exists on the page from the moment its wasm loads, long before
  // anyone calls boot(). fetchDmsg refuses until its dmsg client exists, so the
  // seam asks status() for the same thing rather than finding out by failing.
  describe('readiness', () => {
    it('is not ready with no visor at all', () => {
      expect(transport.visorReady).toBe(false);
    });

    it('is not ready while the visor is unbooted', () => {
      (window as any).skywireVisor = {
        fetchDmsg: () => Promise.resolve({ status: 200 }),
        status: () => ({ booted: false, dmsg: false }),
      };
      expect(transport.visorReady).toBe(false);
    });

    it('is ready once the dmsg client exists', () => {
      (window as any).skywireVisor = {
        fetchDmsg: () => Promise.resolve({ status: 200 }),
        status: () => ({ booted: true, dmsg: true }),
      };
      expect(transport.visorReady).toBe(true);
    });

    // Older builds published no status(). Refusing them would be worse than
    // letting the request report its own failure.
    it('takes a visor with no status() at its word', () => {
      (window as any).skywireVisor = { fetchDmsg: () => Promise.resolve({ status: 200 }) };
      expect(transport.visorReady).toBe(true);
    });

    it('is not ready when status() cannot be read', () => {
      (window as any).skywireVisor = {
        fetchDmsg: () => Promise.resolve({ status: 200 }),
        status: () => { throw new Error('wasm trap'); },
      };
      expect(transport.visorReady).toBe(false);
    });
  });

  // Without a visor the wallet must behave exactly as it always has, which is
  // the reason the seam defaults rather than switches.
  describe('over http', () => {
    it('sends a GET with its parameters', () => {
      let received: any = null;
      transport.request('GET', 'http://127.0.0.1:6420/api/v1/balance', null,
        { params: { addrs: 'addr1' }, headers: { 'Content-Type': 'application/x-www-form-urlencoded' } })
        .subscribe(r => received = r);

      const req = httpMock.expectOne(r => r.url === 'http://127.0.0.1:6420/api/v1/balance');
      expect(req.request.method).toBe('GET');
      expect(req.request.params.get('addrs')).toBe('addr1');
      expect(req.request.headers.get('Content-Type')).toBe('application/x-www-form-urlencoded');

      req.flush({ confirmed: { coins: 1 } });
      expect(received).toEqual({ confirmed: { coins: 1 } });
    });

    it('sends a POST body untouched', () => {
      transport.request('POST', 'http://127.0.0.1:6420/api/v1/wallet/spend', 'id=w1&dst=addr', {}).subscribe();

      const req = httpMock.expectOne(r => r.url.endsWith('/wallet/spend'));
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toBe('id=w1&dst=addr');
      req.flush({});
    });

    afterEach(() => httpMock.verify());
  });

  describe('through a visor', () => {
    it('hands the host, path and body to fetchDmsg and parses the answer', done => {
      const calls: any[] = [];
      (window as any).skywireVisor = {
        fetchDmsg: (host: string, method: string, path: string, body: string, headers: any) => {
          calls.push({ host, method, path, body, headers });

          return Promise.resolve({ status: 200, body: '{"version":{"version":"0.28.0"}}' });
        },
      };

      (window as any).skywireVisor.status = () => ({ booted: true, dmsg: true });

      transport.request('GET', 'http://02a1b2c3.dmsg/api/v1/health', null, { params: { x: '1' } })
        .subscribe(response => {
          expect(calls.length).toBe(1);
          expect(calls[0].host).toBe('02a1b2c3.dmsg');
          expect(calls[0].method).toBe('GET');
          expect(calls[0].path).toBe('/api/v1/health?x=1');
          expect(response).toEqual({ version: { version: '0.28.0' } });
          done();
        });
    });

    it('turns a failing status into an error the caller can read', done => {
      (window as any).skywireVisor = {
        fetchDmsg: () => Promise.resolve({ status: 404, body: '{"error":{"message":"not found"}}' }),
      };

      transport.request('GET', 'http://02a1b2c3.dmsg/api/v1/health', null, {}).subscribe({
        next: () => done.fail('the request should not have succeeded'),
        error: (err: any) => {
          expect(err.status).toBe(404);
          expect(err.error.error.message).toBe('not found');
          done();
        },
      });
    });

    // A dmsg address cannot be fetched by the browser, so failing loudly beats
    // an opaque network error.
    it('says the visor has not booted, which is not the same as having none', done => {
      (window as any).skywireVisor = {
        fetchDmsg: () => Promise.reject(new Error('not booted; call boot() first')),
        status: () => ({ booted: false, dmsg: false }),
      };

      transport.request('GET', 'http://02a1b2c3.dmsg/api/v1/health', null, {}).subscribe({
        next: () => done.fail('the request should not have succeeded'),
        error: (err: Error) => {
          expect(err.message).toContain('has not booted');
          done();
        },
      });
    });

    it('says so when a dmsg node is asked for with no visor present', done => {
      transport.request('GET', 'http://02a1b2c3.dmsg/api/v1/health', null, {}).subscribe({
        next: () => done.fail('the request should not have succeeded'),
        error: (err: Error) => {
          expect(err.message).toContain('skywire visor');
          done();
        },
      });
    });
  });
});
