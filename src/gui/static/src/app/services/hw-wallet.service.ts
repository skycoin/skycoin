import { throwError as observableThrowError, of, Observable, Subject } from 'rxjs';
import { mergeMap, map, catchError } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { BigNumber } from 'bignumber.js';
import { HttpClient } from '@angular/common/http';

import { AppConfig } from '../app.config';
import { HwWalletDaemonService } from './hw-wallet-daemon.service';
import { HwWalletPinService, ChangePinStates } from './hw-wallet-pin.service';
import { ApiService } from './api.service';
import { OperationError, HWOperationResults } from '../utils/operation-error';
import { getErrorMsg } from '../utils/errors';

/**
 * Data about a transaction recipient.
 */
export class HwWalletTxRecipientData {
  address: string;
  coins: BigNumber;
  hours: BigNumber;
}

/**
 * Response when a hw wallet operation completes correctly.
 */
export class OperationResult {
  // This property is currently not very usseful, as it is always set to done, but is not
  // removed because there is chance of needing it in the future.
  /**
   * Operation result.
   */
  result: HWOperationResults;
  /**
   * Raw response returned by the hw wallet daemon.
   */
  rawResponse: any;
}

/**
 * Input of a hw wallet transaction.
 */
export interface HwInput {
  hash: string;
  index: number;
}

/**
 * Output of a hw wallet transaction.
 */
export interface HwOutput {
  address: string;
  coins: string;
  hours: string;
  /**
   * Index of the address on the hw wallet. This indicates the device that the output is
   * used for returning unused remaining coins and hours and that, because of that, it should
   * not be shown while asking the user for confirmation. It only works if the destination
   * address really is is on the device at the indicated index.
   */
  address_index?: number;
}

@Injectable()
export class HwWalletService {
  /**
   * Max number of characters the hw wallet label can have.
   */
  public static readonly maxLabelLength = 32;

  /**
   * If true, the hw wallet options modal window will be shown the next time the wallets list
   * is loaded on the UI.
   */
  showOptionsWhenPossible = false;

  /**
   * Emits every time a device connection/disconnection event is detected.
   */
  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();
  /**
   * Last modal window openned for asking the user confirmation for signing a transaction.
   */
  private signTransactionDialog: MatDialogRef<{}, any>;

  // Set on AppComponent to avoid a circular reference.
  private signTransactionConfirmationComponentInternal;
  /**
   * Sets the class of the modal window used for asking the user confirmation for signing a tx.
   */
  set signTransactionConfirmationComponent(value) {
    this.signTransactionConfirmationComponentInternal = value;
  }

  constructor(
    private dialog: MatDialog,
    private hwWalletDaemonService: HwWalletDaemonService,
    private hwWalletPinService: HwWalletPinService,
    private apiService: ApiService,
    private http: HttpClient) {

    if (this.hwWalletCompatibilityActivated) {
      hwWalletDaemonService.connectionEvent.subscribe(connected => {
        this.walletConnectedSubject.next(connected);
      });
    }
  }

  /**
   * Indicates if the hw wallet compatibility should be activated in the app or not.
   */
  get hwWalletCompatibilityActivated(): boolean {
    return true;
  }

  /**
   * Emits every time a device connection/disconnection event is detected.
   */
  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  /**
   * Detects if there is currently a hw wallet connected.
   */
  getDeviceConnected(): Observable<boolean> {
    return this.hwWalletDaemonService.get('/available').pipe(map((response: any) => {
      return response.data;
    }));
  }

  /**
   * Makes the device to cancel any currently active operation and return to the home screen.
   */
  cancelLastAction(): Observable<OperationResult> {
    this.prepare();

    return this.processDaemonResponse(
      this.hwWalletDaemonService.put('/cancel'),
    );
  }

  /**
   * Gets one or more of the addresses of the hw wallet.
   * @param addressN How many addresses to recover.
   * @param startIndex Starting index.
   */
  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    // Cancel the current pending operation, if any.
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {
        address_n: addressN,
        start_index: startIndex,
        confirm_address: false,
      };

      // Recover the addresses.
      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/generate_addresses',
          params,
        ), null, true,
      );
    }), mergeMap(response => {
      // Check if the device returned valid addresses and create an appropiate error if nedded.
      return this.verifyAddresses(response.rawResponse, 0).pipe(
        catchError(err => {
          const resp = new OperationError();
          resp.originalError = err;
          resp.type = HWOperationResults.AddressGeneratorProblem;
          resp.translatableErrorMsg = this.hwWalletDaemonService.getHardwareWalletErrorMsg(resp.type);
          resp.originalServerErrorMsg = '';

          return observableThrowError(resp);
        }), map(() => response));
    }));
  }

  /**
   * Uses the node to check if all the addresses on an array are valid. Uses concurrency.
   * @param addresses List with the addresses to check.
   * @param currentIndex Index of the address that will be checked on this pass, as this
   * function uses concurrency. This value must normally be 0 when calling this function
   * from outside.
   * @returns No useful value is returned, but the observable will fail if there is an error
   * in any of the addresses.
   */
  private verifyAddresses(addresses: string[], currentIndex: number): Observable<any> {
    const params = {
      address: addresses[currentIndex],
    };

    return this.apiService.post('address/verify', params, {useV2: true}).pipe(mergeMap(() => {
      if (currentIndex !== addresses.length - 1) {
        return this.verifyAddresses(addresses, currentIndex + 1);
      } else {
        return of(0);
      }
    }));
  }

  /**
   * Makes the device ask the user to confirm an address. This allows to check if the software and
   * hardware wallets are showing the same data.
   * @param index Index of the address to confirm.
   */
  confirmAddress(index: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {
        address_n: 1,
        start_index: index,
        confirm_address: true,
      };

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/generate_addresses',
          params,
        ), null, true,
      );
    }));
  }

  /**
   * Gets the features of the connected hw wallet.
   * @param cancelPreviousOperation If true, the function will cancel any current pending
   * operation. Trying to cancel the pending operations may cause problems if the device
   * is in bootloader mode.
   */
  getFeatures(cancelPreviousOperation = true): Observable<OperationResult> {

    let cancel: Observable<any>;
    if (cancelPreviousOperation) {
      cancel = this.cancelLastAction();
    } else {
      cancel = of(0);
    }

    return cancel.pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.get('/features'),
      );
    }));
  }

  /**
   * Downloads the lastest firmware version and installs it on the connected hw wallet.
   * For the operation to work the device must be in bootloader mode.
   * @param downloadCompleteCallback Function called just after the firmware has been
   * downloaded and is going to be sent to the device.
   */
  updateFirmware(downloadCompleteCallback: () => any): Observable<OperationResult> {
    this.prepare();

    // Check if the device is in bootloader mode.
    return this.getFeatures(false).pipe(mergeMap(result => {
      if (!result.rawResponse.bootloader_mode) {
        const resp = new OperationError();
        resp.originalError = result;
        resp.type = HWOperationResults.NotInBootloaderMode;
        resp.translatableErrorMsg = this.hwWalletDaemonService.getHardwareWalletErrorMsg(resp.type);
        resp.originalServerErrorMsg = '';

        return observableThrowError(resp);
      }

      // Get the version number of the lastest firmware.
      return this.http.get(AppConfig.urlForHwWalletVersionChecking, { responseType: 'text' }).pipe(
        catchError(() => {
          return observableThrowError('hardware-wallet.update-firmware.connection-error');
        }),
        mergeMap((res: any) => {
          let lastestFirmwareVersion: string = res.trim();
          if (lastestFirmwareVersion.toLowerCase().startsWith('v')) {
            lastestFirmwareVersion = lastestFirmwareVersion.substr(1, lastestFirmwareVersion.length - 1);
          }

          // Download the lastest firmware.
          return this.http.get(AppConfig.hwWalletDownloadUrlAndPrefix + lastestFirmwareVersion + '.bin', { responseType: 'arraybuffer' }).pipe(
            catchError(() => {
              return observableThrowError('hardware-wallet.update-firmware.connection-error');
            }),
            mergeMap(firmware => {
              downloadCompleteCallback();
              const data = new FormData();
              data.set('file', new Blob([firmware], { type: 'application/octet-stream'}));

              return this.hwWalletDaemonService.put('/firmware_update', data, true);
            }));
        }));
    }));
  }

  /**
   * Sets or changes the PIN of the device.
   * @param changingCurrentPin false if the function was called for setting a PIN in a device
   * which does not have one, true otherwise.
   */
  changePin(changingCurrentPin: boolean): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      // Configure the modal window which will be used to ask for the PIN code.
      this.hwWalletPinService.changingPin = true;
      if (changingCurrentPin) {
        this.hwWalletPinService.changePinState = ChangePinStates.RequestingCurrentPin;
      } else {
        this.hwWalletPinService.changePinState = ChangePinStates.RequestingNewPin;
      }

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post('/configure_pin_code'),
        ['PIN changed'],
      );
    }));
  }

  /**
   * Removes the PIN code protection from the connected device.
   */
  removePin(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {};
      params['remove_pin'] = true;

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/configure_pin_code',
          params,
        ),
        ['PIN removed'],
      );
    }));
  }

  /**
   * Makes a seedless device configure itself with a new random seed.
   * @param wordCount How many words the new seed must have. Must be 12 or 24.
   */
  generateMnemonic(wordCount: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {};
      params['word_count'] = wordCount;
      params['use_passphrase'] = false;

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/generate_mnemonic',
          params,
        ),
        ['Mnemonic successfully configured'],
      );
    }));
  }

  /**
   * Configures a seedless device with a seed entered by the user or checks if a seed entered
   * by the user is equal to the one on the device. The procedure started by the function takes
   * care of showing the UI needed for the user to enter the seed.
   * @param wordCount How many words the seed has.
   * @param dryRun If false, the function will be used to configure a seedless device with the
   * seed provided by the user. If true, the function will just check if the seed provided
   * by the user is equal to the one on the device.
   */
  recoverMnemonic(wordCount: number, dryRun: boolean): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {};
      params['word_count'] = wordCount;
      params['use_passphrase'] = false;
      params['dry_run'] = dryRun;

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/recovery',
          params,
        ),
        ['Device recovered', 'The seed is valid and matches the one in the device'],
      );
    }));
  }

  /**
   * Makes the device show the words of its seed, so the user can back it up. This function
   * works only if the user has not completed a backup before.
   */
  backup(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/backup',
        ),
        ['Device backed up!'],
      );
    }));
  }

  /**
   * Wipes the connected device, deleting the seed, label, PIN, etc.
   */
  wipe(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.delete('/wipe'),
        ['Device wiped'],
      );
    }));
  }

  /**
   * Changes the label the device displays on the home screen.
   * @param label New label to show.
   */
  changeLabel(label: string): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post('/apply_settings', {label: label}),
        ['Settings applied'],
      );
    }));
  }

  /**
   * Makes the connected device create the signatures for a transaction.
   * @param inputs Transaction inputs.
   * @param outputs Transaction outputs.
   */
  signTransaction(inputs: HwInput[], outputs: HwOutput[]): Observable<OperationResult> {
    // Show the confirmation dialog.
    const previewData: HwWalletTxRecipientData[] = [];
    outputs.forEach(output => {
      if (output.address_index === undefined || output.address_index === null) {
        const currentOutput = new HwWalletTxRecipientData();
        currentOutput.address = output.address;
        currentOutput.coins = new BigNumber(output.coins).decimalPlaces(6);
        currentOutput.hours = new BigNumber(output.hours);

        previewData.push(currentOutput);
      }
    });
    this.signTransactionDialog = this.dialog.open(this.signTransactionConfirmationComponentInternal, <MatDialogConfig> {
      width: '600px',
      data: previewData,
    });

    // Make the device ask for confirmation and create the signatures.
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {
        transaction_inputs: inputs,
        transaction_outputs: outputs,
      };

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/transaction_sign',
          params,
        ), null, true,
      ).pipe(map(response => {
        this.closeTransactionDialog();

        return response;
      }), catchError(error => {
        this.closeTransactionDialog();

        return observableThrowError(error);
      }));
    }));
  }

  /**
   * Checks if the first address of the connected hw wallet is equal to the provided address.
   * This allows to check if the connected hw wallet is the one this app needs for an operation.
   * @param firstAddress Address the device should have at index 0.
   * @returns An observable which will fail if the connected device is not the expected one.
   */
  checkIfCorrectHwConnected(firstAddress: string): Observable<any> {
    return this.getAddresses(1, 0).pipe(mergeMap(
      response => {
        if (response.rawResponse[0] !== firstAddress) {
          const resp = new OperationError();
          resp.originalError = response;
          resp.type = HWOperationResults.IncorrectHardwareWallet;
          resp.translatableErrorMsg = this.hwWalletDaemonService.getHardwareWalletErrorMsg(resp.type);
          resp.originalServerErrorMsg = '';

          return observableThrowError(resp);
        }

        return of(true);
      },
    ), catchError(error => {
      const convertedError = error as OperationError;
      if (convertedError.type && convertedError.type === HWOperationResults.WithoutSeed) {
        const resp = new OperationError();
        resp.originalError = error;
        resp.type = HWOperationResults.IncorrectHardwareWallet;
        resp.translatableErrorMsg = this.hwWalletDaemonService.getHardwareWalletErrorMsg(resp.type);
        resp.originalServerErrorMsg = '';

        return observableThrowError(resp);
      }

      return observableThrowError(error);
    }));
  }

  /**
   * Closes the modal window used for asking the user to confirm a transaction.
   */
  private closeTransactionDialog() {
    if (this.signTransactionDialog) {
      this.signTransactionDialog.close();
      this.signTransactionDialog = null;
    }
  }

  /**
   * Gets the observable of an operation made with the hw wallet daemon service and
   * adds to it the steps neded to process the response and errors, to get a properly
   * formated response.
   * @param daemonResponse Observable which will emit the response obtained from the daemon.
   * @param successTexts Texts which are known to be part of the daemon response when the
   * operation finishes correctly. If provided, the daemon response will have to contain
   * any of the text on the array for the operation to be considered successful.
   * @param responseShouldBeArray True if the daemon is expected to return the response
   * as an array.
   * @returns daemonResponse, but with extra steps for making all appropiate operations with
   * the daemon response before emiting it to the subscription.
   */
  private processDaemonResponse(daemonResponse: Observable<any>, successTexts: string[] = null, responseShouldBeArray = false): Observable<any> {
    return daemonResponse.pipe(catchError((error: any) => {
      // Process the error to get it in an appropiate format.
      return observableThrowError(this.buildResponseObject(error, false));
    }), mergeMap(result => {
      // If the response was expected to be an array but was a single value, add it to an array.
      if (responseShouldBeArray && result.data && typeof result.data === 'string') {
        result.data = [result.data];
      }

      // Process the response to get it in an appropiate format.
      const response = this.buildResponseObject(
        result.data ? result.data : null,
        !successTexts ? true : typeof result.data === 'string' && successTexts.some(text => (result.data as string).includes(text)),
      );

      if (response.result === HWOperationResults.Success) {
        return of(response);
      } else {
        return observableThrowError(response);
      }
    }));
  }

  /**
   * Makes all initial preparations needed before sending a request to the hw wallet daemon.
   */
  private prepare() {
    this.hwWalletPinService.changingPin = false;
    this.hwWalletPinService.signingTx = false;
  }

  /**
   * Process a response obtained from the daemon and creates an appropriate object to return to
   * the code which originally requested the operation.
   * @param rawResponse Response obtained from the daemon.
   * @param success If the operation was completed successfuly (true) or ended in an error (false).
   */
  private buildResponseObject(rawResponse: any, success: boolean): any {
    // If the daemon did not respond with an error.
    if ((!rawResponse || !rawResponse.error) && success) {
      const response: OperationResult = {
        result: HWOperationResults.Success,
        rawResponse: rawResponse,
      };

      return response;

      // If the daemon did respond with an error.
    } else {
      // If the response is already an OperationError instance, return it.
      if ((rawResponse as OperationError).type) {
        return rawResponse;
      }

      // Create an appropiate OperationError instance.
      const response = new OperationError();
      response.originalError = rawResponse;
      response.originalServerErrorMsg = getErrorMsg(rawResponse);

      if (!response.originalServerErrorMsg) {
        response.originalServerErrorMsg = rawResponse + '';
      }

      response.type = this.hwWalletDaemonService.getHardwareWalletErrorType(response.originalServerErrorMsg);
      response.translatableErrorMsg = this.hwWalletDaemonService.getHardwareWalletErrorMsg(response.type);

      return response;
    }
  }
}
