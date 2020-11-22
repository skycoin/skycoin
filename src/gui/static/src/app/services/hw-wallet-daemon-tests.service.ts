// This is a version of HwWalletDaemonService that uses a queue to ensure that
// only one operation is carried out at a time. It may be good for testing in
// some circunstances as it allow to see the exact request order without having
// to worry about racing conditions, but  to use it in production it needs more
// testing, to be sure that to be sure that the queue is not going to get stuck
// because of some error. This testing has to be done after solving the problems
// that make the firmware respond extremely slow after some calls.

/*
import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { HwWalletPinService } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import { ISubscription } from 'rxjs/Subscription';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import 'rxjs/add/operator/timeout';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class HwWalletDaemonService {

  public static readonly errorCancelled = 'Cancelled';
  public static readonly errorConnectingWithTheDaemon = 'Error connecting with the hw wallet service';
  public static readonly errorTimeout = 'The operation was canceled due to inactivity';
  private readonly url = 'http://127.0.0.1:9510/api/v1';

  private checkHwSubscription: ISubscription;
  private hwConnected = false;
  private connectionEventSubject = new BehaviorSubject<boolean>(false);
  private waiting = false;
  private busy = false;
  private queueDelay = 32;
  private requestQueue: {
    operation: Observable<any>
    subject: Subject<any>;
  }[] = [];

  get connectionEvent() {
    return this.connectionEventSubject.asObservable();
  }

  constructor(
    private http: Http,
    private apiService: ApiService,
    private hwWalletPinService: HwWalletPinService,
    private hwWalletSeedWordService: HwWalletSeedWordService,
    private ngZone: NgZone,
  ) { }

  get(route: string) {
    const trigger = new Subject<any>();
    this.requestQueue.push({
      subject: trigger,
      operation: trigger.flatMap(() => this.checkResponse(this.http.get(
      this.url + route,
      this.returnRequestOptions(),
        ), route.includes('/available'))),
    });

    setTimeout(() => this.tryToRunNextRequest(), this.queueDelay);

    return this.requestQueue[this.requestQueue.length - 1].operation;
  }

  post(route: string, params = {}) {
    const trigger = new Subject<any>();
    this.requestQueue.push({
      subject: trigger,
      operation: trigger.flatMap(() => this.checkResponse(this.http.post(
      this.url + route,
      JSON.stringify(params),
      this.returnRequestOptions(),
        ))),
    });

    setTimeout(() => this.tryToRunNextRequest(), this.queueDelay);

    return this.requestQueue[this.requestQueue.length - 1].operation;
  }

  put(route: string, params: any = null, sendMultipartFormData = false, smallTimeout = false) {
    const trigger = new Subject<any>();
    this.requestQueue.push({
      subject: trigger,
      operation: trigger.flatMap(() => this.checkResponse(this.http.put(
      this.url + route,
      params,
      this.returnRequestOptions(sendMultipartFormData),
        ), false, smallTimeout)),
    });

    setTimeout(() => this.tryToRunNextRequest(), this.queueDelay);

    return this.requestQueue[this.requestQueue.length - 1].operation;
  }

  delete(route: string) {
    const trigger = new Subject<any>();
    this.requestQueue.push({
      subject: trigger,
      operation: trigger.flatMap(() => this.checkResponse(this.http.delete(
      this.url + route,
      this.returnRequestOptions(),
        ))),
    });

    setTimeout(() => this.tryToRunNextRequest(), this.queueDelay);

    return this.requestQueue[this.requestQueue.length - 1].operation;
  }

  private tryToRunNextRequest() {
    if (!this.busy) {
      if (this.requestQueue.length > 0) {
        this.runNextRequest();
      }
    }
  }

  private runNextRequest() {
    this.prepareToStart();
    this.requestQueue[0].subject.next(1);
    this.requestQueue[0].subject.complete();
  }

  private prepareToStart() {
    if (this.busy) {
      throw new Error('The service is busy.');
    }
    this.busy = true;
  }

  private checkResponse(response: Observable<any>, checkingConnected = false, smallTimeout = false) {
    return response
      .timeout(smallTimeout ? 30000 : 50000)
      .flatMap((res: any) => {
        if (!this.waiting) {
          this.waiting = true;
          setTimeout(() => {
            this.waiting = false;
            this.busy = false;
            this.requestQueue.shift();
            this.tryToRunNextRequest();
          }, this.queueDelay);
        }

        const finalResponse = res.json();

        if (checkingConnected) {
          this.ngZone.run(() => this.updateHwConnected(!!finalResponse.data));
        } else {
          this.updateHwConnected(true);
        }

        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('PinMatrixRequest') !== -1) {
          return this.hwWalletPinService.requestPin().flatMap(pin => {
            if (!pin) {
              return this.put('/cancel').map(() => HwWalletDaemonService.errorCancelled);
            }

            return this.post('/intermediate/pin_matrix', {pin: pin});
          });
        }

        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('WordRequest') !== -1) {
          return this.hwWalletSeedWordService.requestWord().flatMap(word => {
            if (!word) {
              return this.put('/cancel').map(() => HwWalletDaemonService.errorCancelled);
            }

            return this.post('/intermediate/word', {word: word});
          });
        }

        return Observable.of(finalResponse);
      })
      .catch((error: any) => {
        if (!this.waiting) {
          this.waiting = true;
          setTimeout(() => {
            this.waiting = false;
            this.busy = false;
            this.requestQueue.shift();
            this.tryToRunNextRequest();
          }, this.queueDelay);
        }

        if (error && error.name && error.name === 'TimeoutError') {
          this.put('/cancel').subscribe();

          return Observable.throw({_body: HwWalletDaemonService.errorTimeout });
        }

        if (error && error._body)  {
          let errorContent: string;

          if (error._body.error)  {
            errorContent = error._body.error;
          } else {
            try {
              errorContent = JSON.parse(error._body).error;
            } catch (e) {}
          }

          if (errorContent) {
            return this.apiService.processConnectionError(error, true);
          }
        }

        return Observable.throw({_body: HwWalletDaemonService.errorConnectingWithTheDaemon });
      });
  }

  private returnRequestOptions(sendMultipartFormData = false) {
    const options = new RequestOptions();
    options.headers = new Headers();
    if (!sendMultipartFormData) {
      options.headers.append('Content-Type', 'application/json');
    }

    return options;
  }

  checkHw(wait: boolean) {
    if (this.checkHwSubscription) {
      this.checkHwSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.checkHwSubscription = Observable.of(1)
        .delay(wait ? (this.hwConnected ? 2000 : 10000) : 0)
        .flatMap(() => this.get('/available'))
        .subscribe(
          null,
          () => this.ngZone.run(() => this.updateHwConnected(false)),
        );
    });
  }

  private updateHwConnected(connected: boolean) {
    if (connected && !this.hwConnected) {
      this.hwConnected = true;
      this.connectionEventSubject.next(this.hwConnected);
    } else if (!connected && this.hwConnected) {
      this.hwConnected = false;
      this.connectionEventSubject.next(this.hwConnected);
    }
    this.checkHw(true);
  }

}
*/
