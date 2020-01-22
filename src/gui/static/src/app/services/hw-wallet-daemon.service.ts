import { throwError as observableThrowError, of, Observable, SubscriptionLike, BehaviorSubject } from 'rxjs';
import { delay, timeout, mergeMap, map, catchError } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';
import { HwWalletPinService } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import { OperationError, OperationErrorCategories, HWOperationResults } from '../utils/operation-error';
import { getErrorMsg } from '../utils/errors';

@Injectable()
export class HwWalletDaemonService {
  private readonly url = 'http://127.0.0.1:9510/api/v1';

  private checkHwSubscription: SubscriptionLike;
  private hwConnected = false;
  private connectionEventSubject = new BehaviorSubject<boolean>(false);
  private disconnectedChecks = 0;

  private readonly maxFastDisconnectedChecks = 32;

  get connectionEvent() {
    return this.connectionEventSubject.asObservable();
  }

  constructor(
    private http: HttpClient,
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

  private checkResponse(operationResponse: Observable<any>, checkingConnected = false, smallTimeout = false) {
    return operationResponse.pipe(
      timeout(smallTimeout ? 30000 : 55000),
      mergeMap((finalResponse: any) => {
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
          return this.hwWalletPinService.requestPin().pipe(mergeMap(pin => {
            if (!pin) {
              return this.put('/cancel').pipe(mergeMap(() => {
                const response = new OperationError();
                response.category = OperationErrorCategories.HwApiError;
                response.originalError = null;
                response.originalServerErrorMsg = '';
                response.type = HWOperationResults.FailedOrRefused;
                response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response);

                return observableThrowError(response);
              }));
            }

            return this.post('/intermediate/pin_matrix', {pin: pin});
          }));
        }

        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('WordRequest') !== -1) {
          return this.hwWalletSeedWordService.requestWord().pipe(mergeMap(word => {
            if (!word) {
              return this.put('/cancel').pipe(mergeMap(() => {
                const response = new OperationError();
                response.category = OperationErrorCategories.HwApiError;
                response.originalError = null;
                response.originalServerErrorMsg = '';
                response.type = HWOperationResults.FailedOrRefused;
                response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response);

                return observableThrowError(response);
              }));
            }

            return this.post('/intermediate/word', {word: word});
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
        response.category = OperationErrorCategories.HwApiError;
        response.originalError = error;

        if (error && error.name && error.name === 'TimeoutError') {
          this.put('/cancel').subscribe();

          response.originalServerErrorMsg = error.name;
          response.type = HWOperationResults.Timeout;
          response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response);

          return observableThrowError(response);
        }

        const convertedError = error as HttpErrorResponse;
        if (convertedError.status !== null && convertedError.status !== undefined) {
          if (convertedError.status === 0 || convertedError.status === 504) {
            response.originalServerErrorMsg = '';
            response.type = HWOperationResults.DaemonError;
            response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response);
          }
        }

        if (!response.originalServerErrorMsg) {
          response.originalServerErrorMsg = getErrorMsg(error);
        }

        if (!response.type) {
          response.type = this.getHardwareWalletErrorType(response.originalServerErrorMsg);
          response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response);
        }

        return observableThrowError(response);
      }));
  }

  private returnRequestOptions(sendMultipartFormData = false) {
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
        delay(wait ? (this.hwConnected || this.disconnectedChecks < this.maxFastDisconnectedChecks ? 2000 : 10000) : 0),
        mergeMap(() => this.get('/available')))
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

  getHardwareWalletErrorType(responseContent: string) {
    if (!responseContent || typeof responseContent !== 'string') {
      responseContent = '';
    }
    let result: HWOperationResults;

    if (responseContent.toUpperCase().includes('failed or refused'.toUpperCase())) {
      result = HWOperationResults.FailedOrRefused;
    } else if (responseContent.toUpperCase().includes('PIN invalid'.toUpperCase())) {
      result = HWOperationResults.WrongPin;
    } else if (responseContent.toUpperCase().includes('canceled by user'.toUpperCase())) {
      result = HWOperationResults.FailedOrRefused;
    } else if (responseContent.toUpperCase().includes('cancelled by user'.toUpperCase())) {
      result = HWOperationResults.FailedOrRefused;
    } else if (responseContent.toUpperCase().includes('Expected WordAck after Button'.toUpperCase())) {
      result = HWOperationResults.FailedOrRefused;
    } else if (responseContent.toUpperCase().includes('Wrong word retyped'.toUpperCase())) {
      result = HWOperationResults.WrongWord;
    } else if (responseContent.toUpperCase().includes('PIN mismatch'.toUpperCase())) {
      result = HWOperationResults.PinMismatch;
    } else if (responseContent.toUpperCase().includes('Mnemonic not set'.toUpperCase())) {
      result = HWOperationResults.WithoutSeed;
    } else if (responseContent.toUpperCase().includes('Mnemonic required'.toUpperCase())) {
      result = HWOperationResults.WithoutSeed;
    } else if (responseContent.toUpperCase().includes('Invalid seed, are words in correct order?'.toUpperCase())) {
      result = HWOperationResults.InvalidSeed;
    } else if (responseContent.toUpperCase().includes('The seed is valid but does not match the one in the device'.toUpperCase())) {
      result = HWOperationResults.WrongSeed;
    } else if (responseContent.toUpperCase().includes('Invalid base58 character'.toUpperCase())) {
      result = HWOperationResults.InvalidAddress;
    } else if (responseContent.toUpperCase().includes('Invalid address length'.toUpperCase())) {
      result = HWOperationResults.InvalidAddress;
    } else if (responseContent.toUpperCase().includes('LIBUSB'.toUpperCase())) {
      result = HWOperationResults.DaemonError;
    } else if (responseContent.toUpperCase().includes('hidapi'.toUpperCase())) {
      result = HWOperationResults.Disconnected;
      setTimeout(() => this.checkHw(false));
    } else if (responseContent.toUpperCase().includes('device disconnected'.toUpperCase())) {
      result = HWOperationResults.Disconnected;
      setTimeout(() => this.checkHw(false));
    } else if (responseContent.toUpperCase().includes('no device connected'.toUpperCase())) {
      result = HWOperationResults.Disconnected;
      setTimeout(() => this.checkHw(false));
    } else if (responseContent.toUpperCase().includes('MessageType_Success'.toUpperCase())) {
      result = HWOperationResults.Success;
    } else {
      result = HWOperationResults.UndefinedError;
    }

    return result;
  }

  getHardwareWalletErrorMsg(error: OperationError, genericError: string = null): string {
    let response: string;
    if (error.type) {
      if (error.type === HWOperationResults.FailedOrRefused) {
        response = 'hardware-wallet.errors.refused';
      } else if (error.type === HWOperationResults.WrongPin) {
        response = 'hardware-wallet.errors.incorrect-pin';
      } else if (error.type === HWOperationResults.IncorrectHardwareWallet) {
        response = 'hardware-wallet.errors.incorrect-wallet';
      } else if (error.type === HWOperationResults.DaemonError) {
        response = 'hardware-wallet.errors.daemon-connection';
      } else if (error.type === HWOperationResults.InvalidAddress) {
        response = 'hardware-wallet.errors.invalid-address';
      } else if (error.type === HWOperationResults.Timeout) {
        response = 'hardware-wallet.errors.timeout';
      } else if (error.type === HWOperationResults.Disconnected) {
        response = 'hardware-wallet.errors.disconnected';
      } else if (error.type === HWOperationResults.NotInBootloaderMode) {
        response = 'hardware-wallet.errors.not-in-bootloader-mode';
      } else if (error.type === HWOperationResults.PinMismatch) {
        response = 'hardware-wallet.change-pin.pin-mismatch';
      } else if (error.type === HWOperationResults.WrongWord) {
        response = 'hardware-wallet.restore-seed.error-wrong-word';
      } else if (error.type === HWOperationResults.InvalidSeed) {
        response = 'hardware-wallet.restore-seed.error-invalid-seed';
      } else if (error.type === HWOperationResults.WrongSeed) {
        response = 'hardware-wallet.restore-seed.error-wrong-seed';
      } else if (error.type === HWOperationResults.AddressGeneratorProblem) {
        response = 'hardware-wallet.errors.invalid-address-generated';
      } else {
        response = genericError ? genericError : 'hardware-wallet.errors.generic-error';
      }
    } else {
      response = genericError ? genericError : 'hardware-wallet.errors.generic-error';
    }

    return response;
  }

}
