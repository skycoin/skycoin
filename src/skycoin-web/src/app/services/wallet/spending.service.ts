import { Injectable } from '@angular/core';
import { Observable, forkJoin, of, throwError } from 'rxjs';
import { mergeMap, map, catchError, first } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { BigNumber } from 'bignumber.js';

import { ApiService } from '../api.service';
import { CipherProvider } from '../cipher.provider';
import { Output, TransactionInput, TransactionOutput,
  Wallet, GetOutputsRequestOutput, Transaction, GetOutputsRequest } from '../../app.datatypes';
import { BaseCoin } from '../../coins/basecoin';
import { CoinService } from '../coin.service';
import { WalletService } from './wallet.service';
import { GlobalsService } from '../globals.service';
import { isEqualOrSuperiorVersion } from '../../utils/semver';
import { BlockchainService } from '../blockchain.service';
import { HwWalletService, HwInput, HwOutput } from '../hw-wallet.service';

export class Destination {
  address: string;
  coins: BigNumber;
  hours?: BigNumber;
}

export enum HoursSelectionTypes {
  Manual = 'manual',
  Auto = 'auto'
}

export class HoursSelection {
  type: HoursSelectionTypes;
  ShareFactor?: BigNumber;
}

@Injectable()
export class SpendingService {
  isInjectingTx = false;

  private currentCoin: BaseCoin;

  private get coinsMultiplier(): number {
    return this.currentCoin ? this.currentCoin.coinsMultiplier : 1000000;
  }

  constructor(
    private apiService: ApiService,
    private cipherProvider: CipherProvider,
    private translate: TranslateService,
    private walletService: WalletService,
    private globalsService: GlobalsService,
    private blockchainService: BlockchainService,
    private hwWalletService: HwWalletService,
    coinService: CoinService,
  ) {
    coinService.currentCoin.subscribe((coin) => this.currentCoin = coin);
  }

  sendBitcoin(
    wallet: Wallet,
    destinations: Destination[],
    feeRate: number): Observable<string> {
    const body = {
      wallet_id: wallet['filename'] || wallet.label,
      destinations: destinations.map(d => ({
        address: d.address,
        coins: d.coins.multipliedBy(this.currentCoin.coinsMultiplier).integerValue().toNumber(),
      })),
      fee_rate: feeRate,
    };

    return this.apiService.post('btc/send', body, { json: true });
  }

  createTransaction(
    wallet: Wallet,
    addresses: string[]|null,
    unspents: string[]|null,
    destinations: Destination[],
    hoursSelection: HoursSelection,
    changeAddress: string|null): Observable<Transaction> {

    return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
      if (isEqualOrSuperiorVersion(version, '0.26.0')) {
        if (unspents) {
          addresses = null;
        }

        if (!unspents && !addresses) {
          addresses = [];
          wallet.addresses.forEach(address => {
            addresses.push(address.address);
          });
        }

        const processedDestinations = [];
        destinations.forEach(destination => {
          processedDestinations.push({
            address: destination.address,
            coins: destination.coins.toString(),
          });

          if (destination.hours) {
            processedDestinations[processedDestinations.length - 1]['hours'] = destination.hours.toString();
          }
        });

        const processedHours = {
          type: hoursSelection.type,
        };

        if (hoursSelection.type === HoursSelectionTypes.Auto) {
          processedHours['mode'] = 'share';
          processedHours['share_factor'] = hoursSelection.ShareFactor;
        }

        const params = {
          hours_selection: processedHours,
          addresses: addresses,
          unspents: unspents,
          to: processedDestinations,
          change_address: changeAddress,
        };

        const response: Observable<Transaction> = this.apiService.post(
          'transaction',
          params,
          {
            json: true,
          },
          true,
        ).pipe(mergeMap(transaction => {
          const data = transaction.data;

          let hoursSent = new BigNumber('0');
          data.transaction.outputs
            .filter(o => destinations.find(dest => dest.address === o.address))
            .map(o => hoursSent = hoursSent.plus(new BigNumber(o.hours)));

          const txInputs: TransactionInput[] = [];
          data.transaction.inputs.forEach(input => {
            txInputs.push({
              hash: input.uxid,
              secret: wallet.isHardware ? '' : (wallet.addresses.find(a => a.address === input.address) || {}).secret_key,
              address: input.address,
              calculated_hours: input.calculated_hours,
              coins: input.coins,
            });
          });

          const txOutputs: TransactionOutput[] = [];
          data.transaction.outputs.forEach(output => {
            txOutputs.push({
              address: output.address,
              coins: new BigNumber(output.coins).toNumber(),
              hours: new BigNumber(output.hours).toNumber(),
            });
          });

          if (wallet.isHardware) {
            return this.signWithHardwareWallet(wallet, txInputs, txOutputs, hoursSent, new BigNumber(data.transaction.fee));
          }

          return this.generateRawTransaction(txInputs, txOutputs).pipe(
            map((rawTransaction: string) => {
              return {
                inputs: txInputs,
                outputs: txOutputs,
                hoursSent: hoursSent,
                hoursBurned: new BigNumber(data.transaction.fee),
                encoded: rawTransaction
              };
            }));
        }));

        return response;

      } else {

        // Legacy code for 0.25.1 and previous versions
        const unburnedHoursRatio = new BigNumber(1).minus(new BigNumber(1).dividedBy(this.blockchainService.burnRate));

        return this.getOutputs(wallet, addresses, unspents).pipe(
          mergeMap((outputs: Output[]) => {
            // Calculate how many coins should be sent.
            let amount = new BigNumber(0);
            destinations.map(destination => amount = amount.plus(new BigNumber(destination.coins)));

            // If the user manually set the number of hours to send, this calculates the minimum number of hours the
            // inputs must have to satisfy the requirement and the fee.
            let minimumHours = new BigNumber(0);
            if (hoursSelection.type === HoursSelectionTypes.Manual) {
              destinations.map(destination => minimumHours = minimumHours.plus(new BigNumber(destination.hours)));
            }
            minimumHours = minimumHours.multipliedBy(new BigNumber(1).dividedBy(unburnedHoursRatio));

            // Obtain the outputs that will be used as inputs.
            const minRequiredOutputs =  this.getMinRequiredOutputs(amount, minimumHours, outputs, false);
            // Create the transaction.
            let tx = this.buildTransaction(
              unburnedHoursRatio,
              minRequiredOutputs,
              wallet,
              destinations,
              hoursSelection,
              changeAddress,
              amount,
              minimumHours
            );

            // Calculate how many coins are available in the inputs.
            let availableCoins = new BigNumber(0);
            tx.inputs.map(input => availableCoins = availableCoins.plus(new BigNumber(input.coins)));

            // If the inputs have the exact number of coins to send and adding an extra input causes less
            // hours to be lost, it is added.
            if (availableCoins.isEqualTo(amount)) {
              const minRequiredOutputsPlusOne =  this.getMinRequiredOutputs(amount, minimumHours, outputs, true);
              const txExtra = this.buildTransaction(
                unburnedHoursRatio,
                minRequiredOutputsPlusOne,
                wallet, destinations,
                hoursSelection,
                changeAddress,
                amount,
                minimumHours
              );

              const hoursLossedInTx = tx.hoursBurned.plus(tx.hoursSent);
              const hoursLossedInTxExtra = txExtra.hoursBurned.plus(txExtra.hoursSent);

              if (hoursLossedInTx.isGreaterThan(hoursLossedInTxExtra)) {
                tx = txExtra;
              }
            }

            return this.generateRawTransaction(tx.inputs, tx.outputs).pipe(
              map((rawTransaction: string) => {
                return {
                  inputs: tx.inputs,
                  outputs: tx.outputs,
                  hoursSent: tx.hoursSent,
                  hoursBurned: tx.hoursBurned,
                  encoded: rawTransaction
                };
              }));
          }));
      }
    }));
  }

  injectTransaction(encodedTransaction: string): Observable<string> {
    this.isInjectingTx = true;
    return this.postTransaction(encodedTransaction).pipe(
      map(response => {
        this.isInjectingTx = false;
        return response;
      }),
      catchError(err => {
        this.isInjectingTx = false;
        return throwError(() => err);
      }),
    );
  }

  outputsWithWallets(): Observable<Wallet[]> {
    return forkJoin([
      this.walletService.currentWallets.pipe(first()),
      this.getAddressesAsString().pipe(mergeMap(addresses => addresses ? this.getRequestOutputs(addresses) : of([])), first()),
    ]).pipe(map(([wallets, outputs]: [Wallet[], GetOutputsRequestOutput[]]) => {
      return wallets.map(wallet => {
        wallet.addresses = wallet.addresses.map(address => {
          address.outputs = outputs.filter(output => output.address === address.address);

          return address;
        });
        return wallet;
      });
    }));
  }

  getWalletUnspentOutputs(wallet: Wallet): Observable<Output[]> {
    return this.getOutputs(wallet, null, null);
  }

  private signWithHardwareWallet(
    wallet: Wallet,
    txInputs: TransactionInput[],
    txOutputs: TransactionOutput[],
    hoursSent: BigNumber,
    hoursBurned: BigNumber): Observable<Transaction> {

    // Convert inputs to HW daemon format
    const hwInputs: HwInput[] = txInputs.map(input => ({
      hash: input.hash,
      index: wallet.addresses.findIndex(a => a.address === input.address),
    }));

    // Convert outputs to HW daemon format
    const hwOutputs: HwOutput[] = txOutputs.map(output => {
      const hwOutput: HwOutput = {
        address: output.address,
        coins: new BigNumber(output.coins).multipliedBy(this.coinsMultiplier).toFixed(0),
        hours: output.hours.toString(),
      };

      // Mark change outputs with address_index so HW device doesn't ask for confirmation
      const addressIndex = wallet.addresses.findIndex(a => a.address === output.address);
      if (addressIndex !== -1) {
        hwOutput.address_index = addressIndex;
      }

      return hwOutput;
    });

    // First verify the correct HW is connected, then sign
    return this.hwWalletService.checkIfCorrectHwConnected(wallet.addresses[0].address).pipe(
      mergeMap(() => this.hwWalletService.signTransaction(hwInputs, hwOutputs)),
      mergeMap(signResult => {
        const signatures: string[] = Array.isArray(signResult.rawResponse) ? signResult.rawResponse : [signResult.rawResponse];

        const convertedOutputs: TransactionOutput[] = txOutputs.map(output => ({
          ...output,
          coins: parseInt(new BigNumber(output.coins).multipliedBy(this.coinsMultiplier).toFixed(0), 10)
        }));

        return this.cipherProvider.prepareTransactionWithSignatures(txInputs, convertedOutputs, signatures).pipe(
          map((rawTransaction: string) => ({
            inputs: txInputs,
            outputs: txOutputs,
            hoursSent: hoursSent,
            hoursBurned: hoursBurned,
            encoded: rawTransaction,
          }))
        );
      })
    );
  }

  private buildTransaction(
    unburnedHoursRatio: BigNumber,
    minRequiredOutputs: Output[],
    wallet: Wallet,
    destinations: Destination[],
    hoursSelection: HoursSelection,
    changeAddress: string|null,
    amount: BigNumber,
    minimumHours: BigNumber): Transaction {

    let totalCoins = new BigNumber('0');
    minRequiredOutputs.map(output => totalCoins = totalCoins.plus(output.coins));
    let totalHours = new BigNumber('0');
    minRequiredOutputs.map(output => totalHours = totalHours.plus(output.calculated_hours));

    if (totalCoins.isLessThan(amount) || totalHours.isLessThan(minimumHours)) {
      throw new Error(
        `${this.translate.instant('service.wallet.not-enough-hours1')} ${ this.currentCoin.hoursName } ${ this.translate.instant('service.wallet.not-enough-hours2') }`
      );
    }

    const sendableHours = totalHours.multipliedBy(unburnedHoursRatio).decimalPlaces(0, BigNumber.ROUND_FLOOR);
    let hoursToSend: BigNumber;
    let hoursToReturn = new BigNumber(0);

    const txOutputs: TransactionOutput[] = [];
    const txInputs: TransactionInput[] = [];

    const changeCoins = totalCoins.minus(amount).decimalPlaces(6);

    if (hoursSelection.type === HoursSelectionTypes.Auto) {
      if (changeCoins.isGreaterThan(0)) {
        hoursToSend = sendableHours.multipliedBy(hoursSelection.ShareFactor).decimalPlaces(0, BigNumber.ROUND_FLOOR);
      } else {
        hoursToSend = sendableHours;
      }
    } else {
      hoursToSend = new BigNumber(0);
      destinations.map(destination => hoursToSend = hoursToSend.plus(new BigNumber(destination.hours)));
    }

    let sendedHours = new BigNumber(0);
    destinations.forEach(dest => {
      const output: TransactionOutput = {
        address: dest.address,
        coins: dest.coins.toNumber(),
        hours: 0
      };

      let hours: BigNumber;
      if (hoursSelection.type === HoursSelectionTypes.Auto) {
        hours = hoursToSend.multipliedBy(dest.coins.dividedBy(amount)).decimalPlaces(0, BigNumber.ROUND_FLOOR);
      } else {
        hours = dest.hours;
      }
      sendedHours = sendedHours.plus(hours);

      output.hours = hours.toNumber();
      txOutputs.push(output);
    });

    if (hoursSelection.type === HoursSelectionTypes.Auto) {
      if (sendedHours.isLessThan(hoursToSend)) {
        txOutputs[0].hours = new BigNumber(txOutputs[0].hours).plus(hoursToSend.minus(sendedHours)).toNumber();
      }
    }

    if (changeCoins.isGreaterThan(0)) {
      hoursToReturn = sendableHours.minus(hoursToSend);
      changeAddress = changeAddress ? changeAddress : wallet.addresses[0].address;

      txOutputs.push({
        address: changeAddress,
        coins: changeCoins.toNumber(),
        hours: hoursToReturn.toNumber()
      });
    }

    if (changeCoins.isGreaterThan(0) && destinations.some(destination => destination.address === changeAddress)) {
      hoursToSend = sendableHours;
      hoursToReturn = new BigNumber(0);
    }

    minRequiredOutputs.forEach(input => {
      txInputs.push({
        hash: input.hash,
        secret: wallet.addresses.find(a => a.address === input.address).secret_key,
        address: input.address,
        calculated_hours: input.calculated_hours.toNumber(),
        coins: input.coins.toNumber()
      });
    });

    return {
      inputs: txInputs,
      outputs: txOutputs,
      hoursSent: hoursToSend,
      hoursBurned: totalHours.minus(hoursToSend.plus(hoursToReturn))
    };
  }

  private getAddressesAsString(): Observable<string> {
    return this.walletService.currentWallets.pipe(map(wallets => wallets.map(wallet => {
      return wallet.addresses.map(address => address.address).join(',');
    }).join(',')));
  }

  private getOutputs(wallet: Wallet, addresses: string[]|null, unspents: string[]|null): Observable<Output[]> {
    if (!wallet) {
      return of([]);
    } else {
      return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
        const requestedAddresses = addresses ? addresses.join(',') : wallet.addresses.map(a => a.address).join(',');

        let outputsRequest: Observable<any>;
        if (isEqualOrSuperiorVersion(version, '0.25.0')) {
          outputsRequest = this.apiService.post('outputs', { addrs: requestedAddresses });
        } else {
          outputsRequest = this.apiService.get('outputs', { addrs: requestedAddresses });
        }

        let unspentsMap: Map<string, boolean>;
        let addressesMap: Map<string, boolean>;
        if (unspents) {
          unspentsMap = new Map<string, boolean>();
          unspents.map(unspent => unspentsMap.set(unspent, true));
        } else if (addresses) {
          addressesMap = new Map<string, boolean>();
          addresses.map(address => addressesMap.set(address, true));
        }

        return outputsRequest.pipe(map((response: GetOutputsRequest) => {
          const outputs: Output[] = [];
          response.head_outputs.forEach(output => {
            let addOutput = false;
            if (unspents && unspentsMap.has(output.hash)) {
              addOutput = true;
            } else if (!unspents && addresses && addressesMap.has(output.address)) {
              addOutput = true;
            } else if (!unspents && !addresses) {
              addOutput = true;
            }

            if (addOutput) {
              outputs.push({
                address: output.address,
                coins: new BigNumber(output.coins),
                hash: output.hash,
                calculated_hours: new BigNumber(output.calculated_hours)
              });
            }
          });

          return outputs;
        }));
      }));
    }
  }

  private getRequestOutputs(addresses): Observable<GetOutputsRequestOutput[]> {
    if (!addresses) {
      return of([]);
    } else {
      return this.globalsService.getValidNodeVersion().pipe(mergeMap(version => {
        let outputsRequest: Observable<any>;
        if (isEqualOrSuperiorVersion(version, '0.25.0')) {
          outputsRequest = this.apiService.post('outputs', { addrs: addresses });
        } else {
          outputsRequest = this.apiService.get('outputs', { addrs: addresses });
        }

        return outputsRequest.pipe(map(response => response.head_outputs as GetOutputsRequestOutput[]));
      }));
    }
  }

  private postTransaction(rawTransaction: string): Observable<string> {
    return this.apiService.post('injectTransaction', { rawtx: rawTransaction }, { json: true });
  }

  private generateRawTransaction(txInputs: TransactionInput[], txOutputs: TransactionOutput[]): Observable<string> {
    const convertedOutputs: TransactionOutput[] = txOutputs.map(output => {
      return {
        ...output,
        coins: parseInt(new BigNumber(output.coins).multipliedBy(this.coinsMultiplier).toFixed(0), 10)
      };
    });

    return this.cipherProvider.prepareTransaction(txInputs, convertedOutputs);
  }

  private sortOutputs(outputs: Output[], highestToLowest: boolean) {
    outputs.sort((a, b) => {
      if (b.coins.isGreaterThan(a.coins)) {
        return highestToLowest ? 1 : -1;
      } else if (b.coins.isLessThan(a.coins)) {
        return highestToLowest ? -1 : 1;
      } else {
        if (b.calculated_hours.isGreaterThan(a.calculated_hours)) {
          return -1;
        } else if (b.calculated_hours.isLessThan(a.calculated_hours)) {
          return 1;
        } else {
          return 0;
        }
      }
    });
  }

  private getMinRequiredOutputs(transactionAmount: BigNumber, transactionHoursAmount: BigNumber, outputs: Output[], addOneMore: boolean): Output[] {

    // Split the outputs into those with and without hours
    const outputsWithHours: Output[] = [];
    const outputsWitouthHours: Output[] = [];
    outputs.forEach(output => {
      if (output.calculated_hours.isGreaterThan(0)) {
        outputsWithHours.push(output);
      } else {
        outputsWitouthHours.push(output);
      }
    });

    // Abort if there are no outputs with non-zero coinhours.
    if (outputsWithHours.length === 0) {
      return [];
    }

    const sortedOutputs: Output[] = [];
    const minRequiredOutputs: Output[] = [];
    let sumCoins: BigNumber = new BigNumber(0);
    let sumHours: BigNumber = new BigNumber(0);

    // Sort the outputs with hours by coins, from highest to lowest. If two items have the same amount of
    // coins, the one with the least hours is placed first.
    this.sortOutputs(outputsWithHours, true);
    // Add the first nonzero output.
    sortedOutputs.push(outputsWithHours[0]);
    outputsWithHours.splice(0, 1);

    // Sort the outputs without hours by coins, from lowest to highest, and add then to the list.
    this.sortOutputs(outputsWitouthHours, false);
    outputsWitouthHours.map(out => sortedOutputs.push(out));

    // Sort the outputs with hours by coins, from lowest to highest, and add then to the list.
    this.sortOutputs(outputsWithHours, false);
    outputsWithHours.map(out => sortedOutputs.push(out));

    let extraAdded = false;
    // Add outputs until having the necessary amount of coins and hours.
    sortedOutputs.forEach(output => {
      if (sumCoins.isLessThan(transactionAmount) || sumHours.isLessThan(transactionHoursAmount)) {
        minRequiredOutputs.push(output);
        sumCoins = sumCoins.plus(output.coins);
        sumHours = sumCoins.plus(output.calculated_hours);
      } else if (addOneMore === !extraAdded) {
        extraAdded = true;
        minRequiredOutputs.push(output);
      }
    });

    return minRequiredOutputs;
  }
}
