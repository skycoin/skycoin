import { throwError as observableThrowError, of, Observable, Subject } from 'rxjs';
import { mergeMap, map, catchError } from 'rxjs';
import { Injectable } from '@angular/core';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { BigNumber } from 'bignumber.js';

import { HwWalletDaemonService } from './hw-wallet-daemon.service';
import { HwWalletPinService, ChangePinStates } from './hw-wallet-pin.service';
import { ApiService } from './api.service';
import { OperationError, HWOperationResults } from '../utils/operation-error';
import { getHwErrorMsg } from '../utils/hw-errors';

export class HwWalletTxRecipientData {
  address!: string;
  coins!: BigNumber;
  hours!: BigNumber;
}

export class OperationResult {
  result!: HWOperationResults;
  rawResponse: any;
}

export interface HwInput {
  hash: string;
  index: number;
}

export interface HwOutput {
  address: string;
  coins: string;
  hours: string;
  address_index?: number;
}

@Injectable()
export class HwWalletService {
  public static readonly maxLabelLength = 32;

  showOptionsWhenPossible = false;

  private walletConnectedSubject: Subject<boolean> = new Subject<boolean>();
  private signTransactionDialog!: MatDialogRef<{}, any> | null;

  private signTransactionConfirmationComponentInternal: any;
  set signTransactionConfirmationComponent(value: any) {
    this.signTransactionConfirmationComponentInternal = value;
  }

  constructor(
    private dialog: MatDialog,
    private hwWalletDaemonService: HwWalletDaemonService,
    private hwWalletPinService: HwWalletPinService,
    private apiService: ApiService,
  ) {
    hwWalletDaemonService.connectionEvent.subscribe(connected => {
      this.walletConnectedSubject.next(connected);
    });
  }

  get walletConnectedAsyncEvent(): Observable<boolean> {
    return this.walletConnectedSubject.asObservable();
  }

  getDeviceConnected(): Observable<boolean> {
    return this.hwWalletDaemonService.get('/available').pipe(map((response: any) => {
      return response.data;
    }));
  }

  cancelLastAction(): Observable<OperationResult> {
    this.prepare();

    return this.processDaemonResponse(
      this.hwWalletDaemonService.put('/cancel'),
    );
  }

  getAddresses(addressN: number, startIndex: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      const params = {
        address_n: addressN,
        start_index: startIndex,
        confirm_address: false,
      };

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/generate_addresses',
          params,
        ), null, true,
      );
    }), mergeMap(response => {
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

  private verifyAddresses(addresses: string[], currentIndex: number): Observable<any> {
    const params = {
      address: addresses[currentIndex],
    };

    return this.apiService.post('address/verify', params, {}, true).pipe(mergeMap(() => {
      if (currentIndex !== addresses.length - 1) {
        return this.verifyAddresses(addresses, currentIndex + 1);
      } else {
        return of(0);
      }
    }));
  }

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

  changePin(changingCurrentPin: boolean): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

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

  removePin(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/configure_pin_code',
          {remove_pin: true},
        ),
        ['PIN removed'],
      );
    }));
  }

  generateMnemonic(wordCount: number): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post(
          '/generate_mnemonic',
          {word_count: wordCount, use_passphrase: false},
        ),
        ['Mnemonic successfully configured'],
      );
    }));
  }

  backup(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post('/backup'),
        ['Device backed up!'],
      );
    }));
  }

  wipe(): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.delete('/wipe'),
        ['Device wiped'],
      );
    }));
  }

  changeLabel(label: string): Observable<OperationResult> {
    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      return this.processDaemonResponse(
        this.hwWalletDaemonService.post('/apply_settings', {label: label}),
        ['Settings applied'],
      );
    }));
  }

  signTransaction(inputs: HwInput[], outputs: HwOutput[]): Observable<OperationResult> {
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

    if (this.signTransactionConfirmationComponentInternal) {
      this.signTransactionDialog = this.signTransactionConfirmationComponentInternal.openDialog(this.dialog, previewData);
    }

    return this.cancelLastAction().pipe(mergeMap(() => {
      this.prepare();

      this.hwWalletPinService.signingTx = true;

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

  private closeTransactionDialog() {
    if (this.signTransactionDialog) {
      this.signTransactionDialog.close();
      this.signTransactionDialog = null;
    }
  }

  private processDaemonResponse(daemonResponse: Observable<any>, successTexts: string[] | null = null, responseShouldBeArray = false): Observable<any> {
    return daemonResponse.pipe(catchError((error: any) => {
      return observableThrowError(this.buildResponseObject(error, false));
    }), mergeMap(result => {
      if (responseShouldBeArray && result.data && typeof result.data === 'string') {
        result.data = [result.data];
      }

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

  private prepare() {
    this.hwWalletPinService.changingPin = false;
    this.hwWalletPinService.signingTx = false;
  }

  private buildResponseObject(rawResponse: any, success: boolean): any {
    if ((!rawResponse || !rawResponse.error) && success) {
      const response: OperationResult = {
        result: HWOperationResults.Success,
        rawResponse: rawResponse,
      };

      return response;
    } else {
      if ((rawResponse as OperationError).type) {
        return rawResponse;
      }

      const response = new OperationError();
      response.originalError = rawResponse;
      response.originalServerErrorMsg = getHwErrorMsg(rawResponse);

      if (!response.originalServerErrorMsg) {
        response.originalServerErrorMsg = rawResponse + '';
      }

      response.type = this.hwWalletDaemonService.getHardwareWalletErrorType(response.originalServerErrorMsg);
      response.translatableErrorMsg = this.hwWalletDaemonService.getHardwareWalletErrorMsg(response.type);

      return response;
    }
  }
}
