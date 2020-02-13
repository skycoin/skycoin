import { throwError as observableThrowError, of, Observable, SubscriptionLike, BehaviorSubject } from 'rxjs';
import { delay, timeout, mergeMap, catchError } from 'rxjs/operators';
import { Injectable, NgZone } from '@angular/core';
import { HttpClient, HttpHeaders, HttpErrorResponse } from '@angular/common/http';

import { HwWalletPinService } from './hw-wallet-pin.service';
import { HwWalletSeedWordService } from './hw-wallet-seed-word.service';
import { OperationError, HWOperationResults } from '../utils/operation-error';
import { getErrorMsg } from '../utils/errors';

/**
 * Allows to make request to the hw wallet daemon with ease.
 */
@Injectable()
export class HwWalletDaemonService {
  /**
   * URL for accessing the API.
   */
  private readonly url = 'http://127.0.0.1:9510/api/v1/';
  /**
   * Max time the service will wait for responses from the daemon.
   */
  private readonly timeoutMs = 55000;

  /**
   * Allows to know when a devices has be connected/disconnected.
   */
  private connectionEventSubject = new BehaviorSubject<boolean>(false);
  private checkHwSubscription: SubscriptionLike;
  private hwConnected = false;

  /**
   * How many times the service checked if the hw wallet is connected and it was not.
   */
  private disconnectedChecks = 0;
  /**
   * How many times the service will quickly check if the device is connected, before starting
   * to check less frequently.
   */
  private readonly maxFastDisconnectedChecks = 32;

  /**
   * Allows to know when a device has be connected/disconnected.
   */
  get connectionEvent() {
    return this.connectionEventSubject.asObservable();
  }

  constructor(
    private http: HttpClient,
    private hwWalletPinService: HwWalletPinService,
    private hwWalletSeedWordService: HwWalletSeedWordService,
    private ngZone: NgZone,
  ) { }

  /**
   * Sends a GET request to the hw wallet daemon.
   * @param route URL to send the request to. You must omit the "http://x:x/api/vx/" part.
   */
  get(route: string) {
    return this.checkResponse(this.http.get(
      this.getUrl(route),
      this.returnRequestOptions(),
    ), route.includes('/available'));
  }

  /**
   * Sends a POST request to the hw wallet daemon.
   * @param route URL to send the request to. You must omit the "http://x:x/api/vx/" part.
   * @param params Object with the key/value pairs to be sent to the daemon as params.
   */
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

  /**
   * Sends a PUT request to the hw wallet daemon.
   * @param route URL to send the request to. You must omit the "http://x:x/api/vx/" part.
   * @param params Object with the key/value pairs to be sent to the daemon as params.
   * @param sendMultipartFormData If true, the data will be sent as multipart/form-data.
   */
  put(route: string, params: any = null, sendMultipartFormData = false) {
    return this.checkResponse(this.http.put(
      this.getUrl(route),
      params,
      this.returnRequestOptions(sendMultipartFormData),
    ), false);
  }

  /**
   * Sends a DELETE request to the hw wallet daemon.
   * @param route URL to send the request to. You must omit the "http://x:x/api/vx/" part.
   */
  delete(route: string) {
    return this.checkResponse(this.http.delete(
      this.getUrl(route),
      this.returnRequestOptions(),
    ));
  }

  /**
   * Checks and process the responses returned by the daemon. This allows to automatically
   * update the connection status, getting the PIN code and more before returning the final
   * response to the original caller.
   * @param operationResponse Observable which will get the response from the daemon.
   * @param checkingConnected true if the connection to the daemon was made specifically to
   * check if a hw wallet is currently connected.
   * @returns operationResponse, but with extra steps for making all appropiate operations with
   * the daemon response before emiting it to the subscription.
   */
  private checkResponse(operationResponse: Observable<any>, checkingConnected = false): Observable<any> {
    return operationResponse.pipe(
      // This allows to control the timeout errors, instead of having the browser
      // killing the connection at will in an unpredictable way.
      timeout(this.timeoutMs),
      mergeMap((finalResponse: any) => {
        // The daemon may return single value responses as single value arrays. This extracts
        // the response from the array.
        if (finalResponse.data && finalResponse.data.length) {
          if (finalResponse.data.length === 1) {
            finalResponse.data = finalResponse.data[0];
          }
        }

        if (checkingConnected) {
          // Update the connection state with the obtained response.
          this.ngZone.run(() => this.updateHwConnected(!!finalResponse.data));
        } else {
          this.updateHwConnected(true);
        }

        // If the daemon requested the PIN code, ask the user to enter it and send it to
        // the daemon. If the user does not enter the PIN, cancel the operation.
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

        // If the daemon requested a seed word, ask the user to enter it and send it to
        // the daemon. If the user does not enter the word, cancel the operation.
        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('WordRequest') !== -1) {
          return this.hwWalletSeedWordService.requestWord().pipe(mergeMap(word => {
            if (!word) {
              return this.put('/cancel').pipe(mergeMap(() => {
                const response = new OperationError();
                response.originalError = null;
                response.originalServerErrorMsg = '';
                response.type = HWOperationResults.FailedOrRefused;
                response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

                return observableThrowError(response);
              }));
            }

            return this.post('/intermediate/word', {word: word});
          }));
        }

        // This allows the operation to continue. It is an intermediate step when the user has
        // to press a button, which allow to reset the timeout counter, as a new http connection
        // is made.
        if (typeof finalResponse.data === 'string' && (finalResponse.data as string).indexOf('ButtonRequest') !== -1) {
          return this.post('/intermediate/button');
        }

        return of(finalResponse);
      }), catchError((error: any) => {
        // If the error is already an OperationError intance, no more processing is needed.
        // This is needed because in a long operation this part may be called several times
        // and only the first one is needed for processing the error.
        if ((error as OperationError).type) {
          return observableThrowError(error);
        }

        const response = new OperationError();
        response.originalError = error;

        // Process timeouts.
        if (error && error.name && error.name === 'TimeoutError') {
          this.put('/cancel').subscribe();

          response.originalServerErrorMsg = error.name;
          response.type = HWOperationResults.Timeout;
          response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

          return observableThrowError(response);
        }

        // Process connection errors with the daemon.
        const convertedError = error as HttpErrorResponse;
        if (convertedError.status !== null && convertedError.status !== undefined) {
          if (convertedError.status === 0 || convertedError.status === 504) {
            response.originalServerErrorMsg = getErrorMsg(error);
            response.type = HWOperationResults.DaemonConnectionError;
            response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

            return observableThrowError(response);
          }
        }

        // Process the error to get its details.
        response.originalServerErrorMsg = getErrorMsg(error);
        response.type = this.getHardwareWalletErrorType(response.originalServerErrorMsg);
        response.translatableErrorMsg = this.getHardwareWalletErrorMsg(response.type);

        return observableThrowError(response);
      }));
  }

  /**
   * Returns the options object requiered by HttpClient for sending a request.
   * @param sendMultipartFormData If true, the data will be sent as multipart/form-data.
   */
  private returnRequestOptions(sendMultipartFormData = false) {
    const options: any = {};
    options.headers = new HttpHeaders();
    if (!sendMultipartFormData) {
      options.headers = options.headers.append('Content-Type', 'application/json');
    }

    return options;
  }

  /**
   * Checks if there is a hw wallet connected.
   * @param wait false if the check must be made immediately. True if the normal delay must be
   * used.
   */
  private checkHw(wait: boolean) {
    if (this.checkHwSubscription) {
      this.checkHwSubscription.unsubscribe();
    }

    this.ngZone.runOutsideAngular(() => {
      this.checkHwSubscription = of(1).pipe(
        // The delay will be small for a limited number of times, just to catch the cases in
        // which the user connects/disconnects the device quickly
        delay(wait ? (this.hwConnected || this.disconnectedChecks < this.maxFastDisconnectedChecks ? 2000 : 10000) : 0),
        mergeMap(() => this.get('/available')))
        .subscribe(
          // After the response is obtained, the procedure in charge of processing all the
          // responses obtained from the daemon automatically updates the connection status.
          null,
          () => this.ngZone.run(() => this.updateHwConnected(false)),
        );
    });
  }

  /**
   * Receives the the current connection status of the hw wallet and, if needed, dispatch events
   * indicating if the device has been connected or disconnected. It also schedules the periodical
   * automatic connection checking procedure to continue running after an appropiate delay.
   * @param connected If the device is currently connected or not.
   */
  private updateHwConnected(connected: boolean) {
    if (connected) {
      this.disconnectedChecks = 0;
    } else {
      // Keep track of how many checks have been made in which the device has been disconnected.
      this.disconnectedChecks += 1;
    }

    // Update the state if needed.
    if (connected && !this.hwConnected) {
      this.hwConnected = true;
      this.connectionEventSubject.next(this.hwConnected);
    } else if (!connected && this.hwConnected) {
      this.hwConnected = false;
      this.connectionEventSubject.next(this.hwConnected);
    }

    // Make the automatic connection checking run periodically.
    this.checkHw(true);
  }

  /**
   * Get the complete URL needed for making a request to the daemon API.
   * @param url URL to send the request to, omitting the "http://x:x/api/vx/" part.
   */
  private getUrl(url: string): string {
    if (url.startsWith('/')) {
      url = url.substr(1, url.length - 1);
    }

    return this.url + url;
  }

  /**
   * Analyzes an error message to detect to which type it corresponds.
   * @param responseContent Error message to analyze.
   */
  getHardwareWalletErrorType(responseContent: string): HWOperationResults {
    if (!responseContent || typeof responseContent !== 'string') {
      responseContent = '';
    }
    let result: HWOperationResults;

    // Changes in the responses returned by the daemon may affect this.
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
      result = HWOperationResults.DaemonConnectionError;
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

  /**
   * Gets the translatable user-understandable message that corresponds to an error type.
   * @param errorType Error type to check.
   * @returns The error message. If the provided type does not have a corresponding error
   * message, a generic error message is returned.
   */
  getHardwareWalletErrorMsg(errorType: HWOperationResults): string {
    let response: string;
    if (errorType) {
      if (errorType === HWOperationResults.FailedOrRefused) {
        response = 'hardware-wallet.errors.refused';
      } else if (errorType === HWOperationResults.WrongPin) {
        response = 'hardware-wallet.errors.incorrect-pin';
      } else if (errorType === HWOperationResults.IncorrectHardwareWallet) {
        response = 'hardware-wallet.errors.incorrect-wallet';
      } else if (errorType === HWOperationResults.DaemonConnectionError) {
        response = 'hardware-wallet.errors.daemon-connection';
      } else if (errorType === HWOperationResults.InvalidAddress) {
        response = 'hardware-wallet.errors.invalid-address';
      } else if (errorType === HWOperationResults.Timeout) {
        response = 'hardware-wallet.errors.timeout';
      } else if (errorType === HWOperationResults.Disconnected) {
        response = 'hardware-wallet.errors.disconnected';
      } else if (errorType === HWOperationResults.NotInBootloaderMode) {
        response = 'hardware-wallet.errors.not-in-bootloader-mode';
      } else if (errorType === HWOperationResults.PinMismatch) {
        response = 'hardware-wallet.change-pin.pin-mismatch';
      } else if (errorType === HWOperationResults.WrongWord) {
        response = 'hardware-wallet.restore-seed.error-wrong-word';
      } else if (errorType === HWOperationResults.InvalidSeed) {
        response = 'hardware-wallet.restore-seed.error-invalid-seed';
      } else if (errorType === HWOperationResults.WrongSeed) {
        response = 'hardware-wallet.restore-seed.error-wrong-seed';
      } else if (errorType === HWOperationResults.AddressGeneratorProblem) {
        response = 'hardware-wallet.errors.invalid-address-generated';
      } else {
        response = 'hardware-wallet.errors.generic-error';
      }
    } else {
      response = 'hardware-wallet.errors.generic-error';
    }

    return response;
  }
}
