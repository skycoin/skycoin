import { throwError as observableThrowError, of, Observable, concat } from 'rxjs';
import { delay, retryWhen, take, mergeMap, catchError, map } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { BigNumber } from 'bignumber.js';
import { TranslateService } from '@ngx-translate/core';

import { ApiService } from '../api.service';
import { HwWalletService, HwOutput, HwInput } from '../hw-wallet.service';
import { StorageService, StorageType } from '../storage.service';
import { TxEncoder } from '../../utils/tx-encoder';
import { WalletBase } from './wallet-objects';
import { BalanceAndOutputsService } from './balance-and-outputs.service';
import { DecodedTransaction, GeneratedTransaction } from './transaction-objects';
import { processServiceError } from '../../utils/errors';
import { OperationError } from '../../utils/operation-error';

/**
 * Defines a destination to were coins will be sent.
 */
export interface TransactionDestination {
  /**
   * Address to where the coins will be sent.
   */
  address: string;
  /**
   * How many coins to send.
   */
  coins: string;
  /**
   * How many hours to send. Only needed if the node is not supposed to calculate the
   * hours automatically.
   */
  hours?: string;
}

/**
 * Modes the node can use to distribute the hours when creating a transacton.
 */
export enum HoursDistributionTypes {
  /**
   * Every destination will have an specific amout of hours.
   */
  Manual = 'manual',
  /**
   * The node will automatically calculate how many hours to send to each output.
   */
  Auto = 'auto',
}

/**
 * Specifies how the node must distribute the hours when creating a transaction.
 */
export interface HoursDistributionOptions {
  /**
   * How the node will make the calculation.
   */
  type: HoursDistributionTypes;
  /**
   * Specific mode used if the node will automatically calculate the hours.
   */
  mode?: 'share';
  /**
   * Value used by the node to know how many hours to send and retain (is posible), if the node
   * will automatically calculate the hours.
   */
  share_factor?: string;
}

/**
 * Allows to create, prepare and send transactions.
 */
@Injectable()
export class SpendingService {

  constructor(
    private balanceAndOutputsService: BalanceAndOutputsService,
    private apiService: ApiService,
    private hwWalletService: HwWalletService,
    private translate: TranslateService,
    private storageService: StorageService,
  ) { }

  /**
   * Makes the node to create a transaction, but does not send it to the network.
   * @param wallet Wallet from which the coins will be send. If null is provided, you will have to
   * provide a list of addresses or unspent outputs from were the coins will be sent and the function
   * will return an unsigned transaction.
   * @param addresses Optional list of addresses from were the coins will be sent. All addresses should
   * be from the provided wallet (if any). If an unspent outputs list is provided, this param is ignored.
   * @param unspents Optional list of unspent outputs from were the coins will be sent. All outputs
   * should be from the provided wallet (if any).
   * @param hoursDistributionOptions Object indicating how the hours will be distributed.
   * @param destinations Array with indications about hows many coins will be sent and where.
   * @param changeAddress Optional custom address where the remaining coins and hours will be sent. If not
   * provided, the node will select one automatically.
   * @param password Wallet password, if the wallet is encrypted.
   * @param unsigned If the transaction must be signed or not. When using a hw wallet the transaction will
   * have to be signed by the device, so it will have to be connected. If no wallet param was provided, this
   * param is ignored and the transaction will be unsigned.
   * @returns The generated transaction, without the note.
   */
  createTransaction(
    wallet: WalletBase|null,
    addresses: string[]|null,
    unspents: string[]|null,
    destinations: TransactionDestination[],
    hoursDistributionOptions: HoursDistributionOptions,
    changeAddress: string|null,
    password: string|null,
    unsigned: boolean): Observable<GeneratedTransaction> {

    // Create a string indicating where the coins come from.
    let senderString = '';
    if (wallet) {
      senderString = wallet.label;
    } else {
      if (addresses) {
        senderString = addresses.join(', ');
      } else if (unspents) {
        senderString = unspents.join(', ');
      }
    }

    // Ignore the source addresses if specific source outputs were provided.
    if (unspents) {
      addresses = null;
    }

    if (wallet && wallet.isHardware && !changeAddress) {
      // Use the first address of the hw wallet as return address.
      changeAddress = wallet.addresses[0].address;
    }

    const useUnsignedTxEndpoint = !wallet || !!wallet.isHardware;

    const params = {
      hours_selection: hoursDistributionOptions,
      wallet_id: !useUnsignedTxEndpoint ? wallet.id : null,
      password: password,
      addresses: addresses,
      unspents: unspents,
      to: destinations,
      change_address: changeAddress,
    };
    if (!useUnsignedTxEndpoint) {
      params['unsigned'] = unsigned;
    }


    // Make the node create the transaction by using the appropiate URL and sending the
    // previously defined params.
    let response: Observable<GeneratedTransaction> = this.apiService.post(
      useUnsignedTxEndpoint ? 'transaction' : 'wallet/transaction',
      params,
      {
        sendDataAsJson: true,
        useV2: useUnsignedTxEndpoint,
      },
    ).pipe(map(transaction => {
      const data = useUnsignedTxEndpoint ? transaction.data : transaction;

      // Return an error if using a hw wallet and the transaction has too many inputs or outputs.
      if (wallet && wallet.isHardware) {
        if (data.transaction.inputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs-outputs'));
        }
        if (data.transaction.outputs.length > 8) {
          throw new Error(this.translate.instant('hardware-wallet.errors.too-many-inputs-outputs'));
        }
      }

      // Calculate how many coins and hours are being sent.
      let amountToSend = new BigNumber(0);
      destinations.map(destination => amountToSend = amountToSend.plus(destination.coins));

      let hoursToSend = new BigNumber(0);
      data.transaction.outputs
        .filter(o => destinations.map(dest => dest.address).find(addr => addr === o.address))
        .map(o => hoursToSend = hoursToSend.plus(new BigNumber(o.hours)));

      // Process the node response and create a known object.
      const tx: GeneratedTransaction = {
        inputs: (data.transaction.inputs as any[]).map(input => {
          return {
            hash: input.uxid,
            address: input.address,
            coins: new BigNumber(input.coins),
            hours: new BigNumber(input.calculated_hours),
          };
        }),
        outputs: (data.transaction.outputs as any[]).map(output => {
          return {
            hash: output.uxid,
            address: output.address,
            coins: new BigNumber(output.coins),
            hours: new BigNumber(output.hours),
          };
        }),
        coinsToSend: amountToSend,
        hoursToSend: hoursToSend,
        hoursBurned: new BigNumber(data.transaction.fee),
        from: senderString,
        to: destinations.map(destination => destination.address).join(', '),
        wallet: wallet,
        encoded: data.encoded_transaction,
        innerHash: data.transaction.inner_hash,
      };

      return tx;
    }));

    // If required, append to the response the steps needed for signing the transaction with the hw wallet.
    if (wallet && wallet.isHardware && !unsigned) {
      let unsignedTx: GeneratedTransaction;

      response = response.pipe(mergeMap(transaction => {
        unsignedTx = transaction;

        return this.signTransaction(wallet, null, transaction);
      })).pipe(map(encodedSignedTx => {
        unsignedTx.encoded = encodedSignedTx;

        return unsignedTx;
      }));
    }

    return response;
  }

  /**
   * Signs an unsigned transaction.
   * @param wallet Wallet which will be used to sign the transaction.
   * @param password Wallet password, if the provided walled is an encrypted software wallet.
   * @param transaction Transaction to sign.
   * @param rawTransactionString Encoded transaction to sign. If provided, the value of the
   * transaction param is ignored. Only valid if using a software wallet.
   * @returns The encoded signed transaction.
   */
  signTransaction(
    wallet: WalletBase,
    password: string|null,
    transaction: GeneratedTransaction,
    rawTransactionString = ''): Observable<string> {

    // Code for signing a software wallet. The node is responsible for making the operation.
    if (!wallet.isHardware) {
      return this.apiService.post(
        'wallet/transaction/sign',
        {
          wallet_id: wallet.id,
          password: password,
          encoded_transaction: rawTransactionString ? rawTransactionString : transaction.encoded,
        },
        {
          useV2: true,
        },
      ).pipe(map(response => response.data.encoded_transaction));

    // Code for signing a hardware wallet.
    } else {
      if (rawTransactionString) {
        throw new Error('Raw transactions not allowed.');
      }

      const hwOutputs: HwOutput[] = [];
      const hwInputs: HwInput[] = [];

      const addressesMap: Map<string, number> = new Map<string, number>();
      wallet.addresses.forEach((address, i) => addressesMap.set(address.address, i));

      // Convert all inputs and outputs to the format used by the hw wallet.
      transaction.outputs.forEach(output => {
        hwOutputs.push({
          address: output.address,
          coins: new BigNumber(output.coins).toString(),
          hours: new BigNumber(output.hours).toFixed(0),
        });
      });
      transaction.inputs.forEach(input => {
        hwInputs.push({
          hash: input.hash,
          index: addressesMap.get(input.address),
        });
      });

      if (hwOutputs.length > 1) {
        // Try to find the return address assuming that it is the first address of the device and that
        // it should be at the end of the outputs list.
        for (let i = hwOutputs.length - 1; i >= 0; i--) {
          if (hwOutputs[i].address === wallet.addresses[0].address) {
            // This makes de device consider the output as the one used for returning the remaining coins.
            hwOutputs[i].address_index = 0;
            break;
          }
        }
      }

      // Make the device sign the transaction.
      return this.hwWalletService.signTransaction(hwInputs, hwOutputs).pipe(map(signatures => {
        const rawTransaction = TxEncoder.encode(
          hwInputs,
          hwOutputs,
          signatures.rawResponse,
          transaction.innerHash,
        );

        return rawTransaction;
      }));
    }
  }

  /**
   * Creates a DecodedTransaction instance from a raw transaction string.
   * @param rawTransactionString Raw transaction string.
   * @param usigned If the raw transaction is unsigned or not.
   */
  decodeTransaction(rawTransactionString: string, usigned: boolean): Observable<DecodedTransaction> {
    // Make the call to the back-end.
    return this.apiService.post(
      'transaction/verify',
      {
        unsigned: usigned,
        encoded_transaction: rawTransactionString,
      },
      {
        useV2: true,
      },
    ).pipe(catchError((err: OperationError) => {
      // If the node returned an error, but also the transaction, get the transaction and continue normally.
      err = processServiceError(err);
      if (err && err.originalError && err.originalError.status === 422 && err.originalError.error && err.originalError.error.data) {
        return of(err.originalError.error.data);
      }

      return observableThrowError(err);
    }),
    map(r => {
      if (r.data) {
        r = r.data;
      }

      // Process the node response and create a known object.
      let inputsInformationObtained = true;
      const tx: DecodedTransaction = {
        inputs: (r.transaction.inputs as any[]).map(input => {
          if (!input.txid) {
            inputsInformationObtained = false;
          }

          return {
            hash: input.uxid,
            address: input.address,
            coins: new BigNumber(input.coins),
            hours: new BigNumber(input.calculated_hours),
          };
        }),
        outputs: (r.transaction.outputs as any[]).map(output => {
          return {
            hash: output.uxid,
            address: output.address,
            coins: new BigNumber(output.coins),
            hours: new BigNumber(output.hours),
          };
        }),
        id: r.transaction.txid,
        hoursBurned: null,
        inputsInformationObtained: inputsInformationObtained,
      };

      return tx;
    }));
  }

  /**
   * Sends a signed transaction to the network, to efectivelly send the coins.
   * @param encodedTx Transaction to send.
   * @param note Optional local note for the transaction.
   * @returns If the note was saved or not.
   */
  injectTransaction(encodedTx: string, note: string|null): Observable<boolean> {
    // Send the transaction.
    return this.apiService.post('injectTransaction', { rawtx: encodedTx }, { sendDataAsJson: true }).pipe(
      mergeMap(txId => {
        // Refresh the balance after a small delay.
        setTimeout(() => this.balanceAndOutputsService.refreshBalance(), 32);

        if (!note) {
          return of(false);
        } else {
          // Save the note. Retry 3 times if an error is found.
          return this.storageService.store(StorageType.NOTES, txId, note).pipe(
            retryWhen(errors => concat(errors.pipe(delay(1000), take(3)), observableThrowError(-1))),
            catchError(err => err === -1 ? of(-1) : err),
            map(result => result === -1 ? false : true));
        }
      }));
  }
}
