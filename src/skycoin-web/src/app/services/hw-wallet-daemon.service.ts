import { throwError as observableThrowError, of, Observable, SubscriptionLike, BehaviorSubject } from 'rxjs';
import { delay, timeout, mergeMap, catchError } from 'rxjs';
import { Injectable, NgZone } from '@angular/core';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';

import { HwWalletPinService } from './hw-wallet-pin.service';
import { OperationError, HWOperationResults } from '../utils/operation-error';
import { getHwErrorMsg } from '../utils/hw-errors';

@Injectable()
export class HwWalletDaemonService {
  private readonly url = 'http://127.0.0.1:9510/api/v1/';
  private readonly timeoutMs = 55000;

  private connectionEventSubject = new BehaviorSubject<boolean>(false);
  private checkHwSubscription!: SubscriptionLike;
  private hwConnected = false;

  private disconnectedChecks = 0;
  private readonly maxFastDisconnectedChecks = 32;
  private readonly updatePeriod = 10 * 1000;
  private readonly fastUpdatePeriod = 2 * 1000;

  get connectionEvent() {
    return this.connectionEventSubject.asObservable();
  }

  constructor(
    private http: HttpClient,
    private hwWalletPinService: HwWalletPinService,
    private ngZone: NgZone,
  ) { }

  get(route: string) {
    return this.checkResponse(this.http.get(
      this.getUrl(route),
      this.returnRequestOptions(),
    ), route.includes('/available'));
  }

  post(route: string, params: any = null) {
    if (!params) {
      params = {};
    }

    return this.checkResponse(this.http.post(
      this.getUrl(route),
      JSON.stringify(params),
      this.returnRequestOptions(),
    ));
  }

  put(route: string, params: any = null, sendMultipartFormData = false) {
    return this.checkResponse(this.http.put(
      this.getUrl(route),
      params,
      this.returnRequestOptions(sendMultipartFormData),
    ), false);
  }

  delete(route: string) {
    return this.checkResponse(this.http.delete(
      this.getUrl(route),
      this.returnRequestOptions(),
    ));
  }

  private checkResponse(operationResponse: Observable<any>, checkingConnected = false): Observable<any> {
    return operationResponse.pipe(
      timeout(this.timeoutMs),
      mergeMap((finalResponse: any) => {
        if (finalResponse.data && finalResponse.data.length) {
          if (finalResponse.data.length === 1) {
            finalResponse.data = finalResponse.data[0];
          }
        }

        if (checkingConnected) {
          this.ngZone.run(() => this.updateHwConnected(!!finalResponse.data));
        } else {
          this.updateHwConnected(true);
        }

        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('PinMatrixRequest') !== -1) {
          return this.hwWalletPinService.requestPin().pipe(mergeMap(pin => {
            if (!pin) {
              return this.put('/cancel').pipe(mergeMap(() => {
                const response = new OperationError();
                response.originalError = null;
                response.originalServerErrorMsg = '';
                response.type = HWOperationResults.FailedOrRefused;
                response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

                return observableThrowError(response);
              }));
            }

            return this.post('/intermediate/pin_matrix', {pin: pin});
          }));
        }

        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('ButtonRequest') !== -1) {
          return this.post('/intermediate/button');
        }

        return of(finalResponse);
      }), catchError((error: any) => {
        if ((error as OperationError).type) {
          return observableThrowError(error);
        }

        const response = new OperationError();
        response.originalError = error;

        if (error && error.name && error.name === 'TimeoutError') {
          this.put('/cancel').subscribe();

          response.originalServerErrorMsg = error.name;
          response.type = HWOperationResults.Timeout;
          response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

          return observableThrowError(response);
        }

        const convertedError = error as HttpErrorResponse;
        if (convertedError.status !== null && convertedError.status !== undefined) {
          if (convertedError.status === 0 || convertedError.status === 504) {
            response.originalServerErrorMsg = getHwErrorMsg(error);
            response.type = HWOperationResults.DaemonConnectionError;
            response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

            return observableThrowError(response);
          }
        }

        response.originalServerErrorMsg = getHwErrorMsg(error);
        response.type = this.getHardwareWalletErrorType(response.originalServerErrorMsg || '');
        response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

        return observableThrowError(response);
      }));
  }

  private returnRequestOptions(sendMultipartFormData = false): any {
    const options: any = {};
    options.headers = new HttpHeaders();
    if (!sendMultipartFormData) {
      options.headers = options.headers.append('Content-Type', 'application/json');
    }

    return options;
  }

  checkHw(wait: boolean) {
    if (this.checkHwSubscription) {
      this.checkHwSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.checkHwSubscription = of(1).pipe(
        delay(wait ? (this.hwConnected || this.disconnectedChecks < this.maxFastDisconnectedChecks ? this.fastUpdatePeriod : this.updatePeriod) : 0),
        mergeMap(() => this.get('/available')))
        .subscribe({
          next: undefined,
          error: () => this.ngZone.run(() => this.updateHwConnected(false)),
        });
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

  private getUrl(url: string): string {
    if (url.startsWith('/')) {
      url = url.substr(1, url.length - 1);
    }

    return this.url + url;
  }

  getHardwareWalletErrorType(responseContent: string): HWOperationResults {
    if (!responseContent || typeof responseContent !== 'string') {
      responseContent = '';
    }

    const upper = responseContent.toUpperCase();

    if (upper.includes('FAILED OR REFUSED')) {
      return HWOperationResults.FailedOrRefused;
    } else if (upper.includes('PIN INVALID')) {
      return HWOperationResults.WrongPin;
    } else if (upper.includes('CANCELED BY USER') || upper.includes('CANCELLED BY USER')) {
      return HWOperationResults.FailedOrRefused;
    } else if (upper.includes('EXPECTED WORDACK AFTER BUTTON')) {
      return HWOperationResults.FailedOrRefused;
    } else if (upper.includes('WRONG WORD RETYPED')) {
      return HWOperationResults.WrongWord;
    } else if (upper.includes('PIN MISMATCH')) {
      return HWOperationResults.PinMismatch;
    } else if (upper.includes('MNEMONIC NOT SET') || upper.includes('MNEMONIC REQUIRED')) {
      return HWOperationResults.WithoutSeed;
    } else if (upper.includes('INVALID SEED, ARE WORDS IN CORRECT ORDER?')) {
      return HWOperationResults.InvalidSeed;
    } else if (upper.includes('THE SEED IS VALID BUT DOES NOT MATCH')) {
      return HWOperationResults.WrongSeed;
    } else if (upper.includes('INVALID BASE58 CHARACTER') || upper.includes('INVALID ADDRESS LENGTH')) {
      return HWOperationResults.InvalidAddress;
    } else if (upper.includes('LIBUSB')) {
      return HWOperationResults.DaemonConnectionError;
    } else if (upper.includes('HIDAPI') || upper.includes('DEVICE DISCONNECTED') || upper.includes('NO DEVICE CONNECTED')) {
      setTimeout(() => this.checkHw(false));
      return HWOperationResults.Disconnected;
    } else if (upper.includes('MESSAGETYPE_SUCCESS')) {
      return HWOperationResults.Success;
    }

    return HWOperationResults.UndefinedError;
  }

  getHardwareWalletErrorMsg(errorType: HWOperationResults): string {
    switch (errorType) {
      case HWOperationResults.FailedOrRefused:
        return 'hardware-wallet.errors.refused';
      case HWOperationResults.WrongPin:
        return 'hardware-wallet.errors.incorrect-pin';
      case HWOperationResults.IncorrectHardwareWallet:
        return 'hardware-wallet.errors.incorrect-wallet';
      case HWOperationResults.DaemonConnectionError:
        return 'hardware-wallet.errors.daemon-connection';
      case HWOperationResults.InvalidAddress:
        return 'hardware-wallet.errors.invalid-address';
      case HWOperationResults.Timeout:
        return 'hardware-wallet.errors.timeout';
      case HWOperationResults.Disconnected:
        return 'hardware-wallet.errors.disconnected';
      case HWOperationResults.PinMismatch:
        return 'hardware-wallet.errors.pin-mismatch';
      case HWOperationResults.WrongWord:
        return 'hardware-wallet.errors.wrong-word';
      case HWOperationResults.InvalidSeed:
        return 'hardware-wallet.errors.invalid-seed';
      case HWOperationResults.WrongSeed:
        return 'hardware-wallet.errors.wrong-seed';
      case HWOperationResults.AddressGeneratorProblem:
        return 'hardware-wallet.errors.invalid-address-generated';
      default:
        return 'hardware-wallet.errors.generic-error';
    }
  }
}
