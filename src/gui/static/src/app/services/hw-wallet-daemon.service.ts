import { Injectable, NgZone } from '@angular/core';
import { ApiService } from './api.service';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { HwWalletPinService } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import { ISubscription } from 'rxjs/Subscription';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import 'rxjs/add/operator/timeout';

@Injectable()
export class HwWalletDaemonService {

  public static readonly errorCancelled = 'Cancelled';
  public static readonly errorConnectingWithTheDaemon = 'Error connecting with the hw wallet service';
  public static readonly errorTimeout = 'The operation was canceled due to inactivity';
  private readonly url = 'http://127.0.0.1:9510/api/v1';

  private checkHwSubscription: ISubscription;
  private hwConnected = false;
  private connectionEventSubject = new BehaviorSubject<boolean>(false);
  private disconnectedChecks = 0;

  private readonly maxFastDisconnectedChecks = 32;

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
    return this.checkResponse(this.http.get(
      this.url + route,
      this.returnRequestOptions(),
    ), route.includes('/available'));
  }

  post(route: string, params = {}) {
    return this.checkResponse(this.http.post(
      this.url + route,
      JSON.stringify(params),
      this.returnRequestOptions(),
    ));
  }

  put(route: string, params: any = null, sendMultipartFormData = false, smallTimeout = false) {
    return this.checkResponse(this.http.put(
      this.url + route,
      params,
      this.returnRequestOptions(sendMultipartFormData),
    ), false, smallTimeout);
  }

  delete(route: string) {
    return this.checkResponse(this.http.delete(
      this.url + route,
      this.returnRequestOptions(),
    ));
  }

  private checkResponse(response: Observable<any>, checkingConnected = false, smallTimeout = false) {
    return response
      .timeout(smallTimeout ? 30000 : 55000)
      .flatMap((res: any) => {
        const finalResponse = res.json();

        if (finalResponse.data && finalResponse.data.length) {
          if (finalResponse.data.length === 1) {
            finalResponse.data = finalResponse.data[0];
          } else {
            finalResponse.data = finalResponse.data;
          }
        }

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

        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('ButtonRequest') !== -1) {
          return this.post('/intermediate/button');
        }

        return Observable.of(finalResponse);
      })
      .catch((error: any) => {
        if (error && error.name && error.name === 'TimeoutError') {
          this.put('/cancel').subscribe();

          return Observable.throw({_body: HwWalletDaemonService.errorTimeout });
        }

        if (error && error._body)  {
          let errorContent: string;

          if (typeof error._body === 'string')  {
            errorContent = error._body;
          } else if (error._body.error)  {
            errorContent = error._body.error;
          } else {
            try {
              errorContent = JSON.parse(error._body).error;
            } catch (e) {}
          }

          if (errorContent) {
            return this.apiService.processConnectionError({_body: errorContent}, true);
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
        .delay(wait ? (this.hwConnected || this.disconnectedChecks < this.maxFastDisconnectedChecks ? 2000 : 10000) : 0)
        .flatMap(() => this.get('/available'))
        .subscribe(
          null,
          () => this.ngZone.run(() => this.updateHwConnected(false)),
        );
    });
  }

  private updateHwConnected(connected: boolean) {
    if (connected) {
      this.disconnectedChecks = 0;
    } else {
      this.disconnectedChecks += 1;
    }

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
