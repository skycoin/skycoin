import { throwError as observableThrowError, of, Observable } from 'rxjs';
import { concat, delay, retryWhen, take, mergeMap, catchError, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { ApiService } from '../api.service';
import { PreviewTransaction } from '../../app.datatypes';
import { BigNumber } from 'bignumber.js';
import { HwWalletService, HwOutput, HwInput } from '../hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { StorageService, StorageType } from '../storage.service';
import { TxEncoder } from '../../utils/tx-encoder';
import { WalletBase } from './wallet-objects';
import { BalanceAndOutputsService } from './balance-and-outputs.service';

export interface TransactionDestination {
  address: string;
  coins: string;
  hours?: string;
}

export enum HoursDistributionTypes {
  Manual = 'manual',
  Auto = 'auto',
}

export interface HoursDistributionOptions {
  type: HoursDistributionTypes;
  mode?: 'share';
  share_factor?: string;
}

@Injectable()
export class SpendingService {

  constructor(
    private balanceAndOutputsService: BalanceAndOutputsService,
    private apiService: ApiService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private storageService: StorageService,
  ) { }

  createTransaction(
    wallet: WalletBase|null,
    addresses: string[]|null,
    unspents: string[]|null,
    destinations: TransactionDestination[],
    hoursSelection: HoursDistributionOptions,
    changeAddress: string|null,
    password: string|null,
    unsigned: boolean): Observable<PreviewTransaction> {

    if (unspents) {
      addresses = null;
    }

    if (wallet && wallet.isHardware && !changeAddress) {
      changeAddress = wallet.addresses[0].address;
    }

    const useV2Endpoint = !wallet || !!wallet.isHardware;

    const params = {
      hours_selection: hoursSelection,
      wallet_id: wallet ? wallet.id : null,
      password: password,
      addresses: addresses,
      unspents: unspents,
      to: destinations,
      change_address: changeAddress,
    };

    if (!useV2Endpoint) {
      params['unsigned'] = unsigned;
    }

    let response: Observable<PreviewTransaction> = this.apiService.post(
      useV2Endpoint ? 'transaction' : 'wallet/transaction',
      params,
      {
        sendDataAsJson: true,
        useV2: useV2Endpoint,
      },
    ).pipe(map(transaction => {
      const data = useV2Endpoint ? transaction.data : transaction;

      if (wallet && wallet.isHardware) {
        if (data.transaction.inputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs-outputs'));
        }
        if (data.transaction.outputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs-outputs'));
        }
      }

      return {
        ...data.transaction,
        hoursBurned: new BigNumber(data.transaction.fee),
        encoded: data.encoded_transaction,
        innerHash: data.transaction.inner_hash,
      };
    }));

    if (wallet && wallet.isHardware && !unsigned) {
      let unsignedTx: PreviewTransaction;

      response = response.pipe(mergeMap(transaction => {
        unsignedTx = transaction;

        return this.signTransaction(wallet, null, transaction);
      })).pipe(map(signedTx => {
        unsignedTx.encoded = signedTx.encoded;

        return unsignedTx;
      }));
    }

    return response;
  }

  signTransaction(
    wallet: WalletBase,
    password: string|null,
    transaction: PreviewTransaction,
    rawTransactionString = ''): Observable<PreviewTransaction> {

    if (!wallet.isHardware) {
      return this.apiService.post(
        'wallet/transaction/sign',
        {
          wallet_id: wallet ? wallet.id : null,
          password: password,
          encoded_transaction: rawTransactionString ? rawTransactionString : transaction.encoded,
        },
        {
          useV2: true,
        },
      ).pipe(map(response => {
        return {
          ...response.data.transaction,
          hoursBurned: new BigNumber(response.data.transaction.fee),
          encoded: response.data.encoded_transaction,
        };
      }));

    } else {
      if (rawTransactionString) {
        throw new Error('Raw transactions not allowed.');
      }

      const txOutputs = [];
      const txInputs = [];
      const hwOutputs: HwOutput[] = [];
      const hwInputs: HwInput[] = [];

      transaction.outputs.forEach(output => {
        txOutputs.push({
          address: output.address,
          coins: parseInt(new BigNumber(output.coins).multipliedBy(1000000).toFixed(0), 10),
          hours: parseInt(output.hours, 10),
        });

        hwOutputs.push({
          address: output.address,
          coins: new BigNumber(output.coins).toString(),
          hours: new BigNumber(output.hours).toFixed(0),
        });
      });

      if (hwOutputs.length > 1) {
        for (let i = txOutputs.length - 1; i >= 0; i--) {
          if (hwOutputs[i].address === wallet.addresses[0].address) {
            hwOutputs[i].address_index = 0;
            break;
          }
        }
      }

      const addressesMap: Map<string, number> = new Map<string, number>();
      wallet.addresses.forEach((address, i) => addressesMap.set(address.address, i));

      transaction.inputs.forEach(input => {
        txInputs.push({
          hash: input.uxid,
          secret: '',
          address: input.address,
          address_index: addressesMap.get(input.address),
          calculated_hours: parseInt(input.calculated_hours, 10),
          coins: parseInt(input.coins, 10),
        });

        hwInputs.push({
          hash: input.uxid,
          index: addressesMap.get(input.address),
        });
      });

      return this.hwWalletService.signTransaction(hwInputs, hwOutputs).pipe(mergeMap(signatures => {
        const rawTransaction = TxEncoder.encode(
          hwInputs,
          hwOutputs,
          signatures.rawResponse,
          transaction.innerHash,
        );

        return of({
          ...transaction,
          encoded: rawTransaction,
        });
      }));
    }
  }

  injectTransaction(encodedTx: string, note: string): Observable<boolean> {
    return this.apiService.post('injectTransaction', { rawtx: encodedTx }, { sendDataAsJson: true }).pipe(
      mergeMap(txId => {
        setTimeout(() => this.balanceAndOutputsService.refreshBalance(), 32);

        if (!note) {
          return of(false);
        } else {
          return this.storageService.store(StorageType.NOTES, txId, note).pipe(
            retryWhen(errors => errors.pipe(delay(1000), take(3), concat(observableThrowError(-1)))),
            catchError(err => err === -1 ? of(-1) : err),
            map(result => result === -1 ? false : true));
        }
      }));
  }
}
