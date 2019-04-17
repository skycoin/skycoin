import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { Http, RequestOptions, Headers } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { HwWalletPinService } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import { ISubscription } from 'rxjs/Subscription';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import 'rxjs/add/operator/timeout';

export enum ConnectionMethods {
  Get,
  Post,
  Put,
  Delete,
}

@Injectable()
export class HwWalletDaemonService {

  public static readonly errorCancelled = 'Cancelled';
  public static readonly errorConnectingWithTheDaemon = 'Error connecting with the hw wallet service';
  public static readonly errorTimeout = 'The operation was cancelled due to inactivity';
  private readonly url = 'http://127.0.0.1:9510/api/v1';

  private checkHwSubscription: ISubscription;
  private hwConnected = false;
  private connectionEventSubject = new BehaviorSubject<boolean>(false);

  get connectionEvent() {
    return this.connectionEventSubject.asObservable();
  }

  constructor(
    private http: Http,
    private apiService: ApiService,
    private hwWalletPinService: HwWalletPinService,
    private hwWalletSeedWordService: HwWalletSeedWordService,
  ) {
    this.checkHw(false);
  }

  callFunction(route: string, connectionMethod: ConnectionMethods, params = {}) {
    if (connectionMethod === ConnectionMethods.Post) {
      return this.post(route, params);
    } else if (connectionMethod === ConnectionMethods.Get) {
      return this.get(route);
    } else if (connectionMethod === ConnectionMethods.Put) {
      return this.put(route);
    } else if (connectionMethod === ConnectionMethods.Delete) {
      return this.delete(route);
    }
  }

  private get(route: string) {
    return this.checkResponse(this.http.get(
      this.url + route,
      this.returnRequestOptions(),
    ));
  }

  private post(route: string, params = {}) {
    return this.checkResponse(this.http.post(
      this.url + route,
      JSON.stringify(params),
      this.returnRequestOptions(),
    ));
  }

  private put(route: string) {
    return this.checkResponse(this.http.put(
      this.url + route,
      null,
      this.returnRequestOptions(),
    ));
  }

  private delete(route: string) {
    return this.checkResponse(this.http.delete(
      this.url + route,
      this.returnRequestOptions(),
    ));
  }

  private checkResponse(response: Observable<any>) {
    return response
      .timeout(50000)
      .flatMap((res: any) => {
        const finalResponse = res.json();

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
        if (error && error.name && error.name === 'TimeoutError') {
          this.put('/cancel').map(() => HwWalletDaemonService.errorCancelled).subscribe();

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

  private returnRequestOptions() {
    const options = new RequestOptions();
    options.headers = new Headers();
    options.headers.append('Content-Type', 'application/json');

    return options;
  }

  private checkHw(wait: boolean) {
    // Reactivate this when having a more reliable method for detecting if the device is connected.
    /*
    if (this.checkHwSubscription) {
      this.checkHwSubscription.unsubscribe();
    }

    this.checkHwSubscription = Observable.of(1)
      .delay(wait ? (this.hwConnected ? 2000 : 10000) : 0)
      .flatMap(() => this.get('/features'))
      .subscribe(
        (response: any) => this.updateHwConnected(!!response.data && !!response.data.features),
        (error: any) => this.updateHwConnected(error && error.message && typeof error.message === 'string' && error.message.indexOf('Unknown message read_tiny') !== -1),
      );
      */
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
