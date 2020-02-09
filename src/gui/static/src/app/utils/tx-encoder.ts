import BigNumber from 'bignumber.js';
import * as Base58 from 'base-x';

import { HwInput, HwOutput } from '../services/hw-wallet.service';

/**
 * Allows to encode transactions, to be able to send them to the network.
 */
export class TxEncoder {
  /**
   * Creates an encoded transaction using the Skycoin format. Check the encoded Skycoin
   * transaction reference for more info.
   * @param inputs List of all inputs.
   * @param outputs List of all outputs.
   * @param signatures List of all signatures. There must be one signature per input.
   * @param innerHash Inner hash of the transaction.
   * @param transactionType Number identifying the type of the transaction, as per the Skycoin
   * transaction format.
   */
  static encode(inputs: HwInput[], outputs: HwOutput[], signatures: string[], innerHash: string, transactionType = 0): string {
    if (inputs.length !== signatures.length) {
      throw new Error('Invalid number of signatures.');
    }

    // Calculate the size of the transaction and initialize the object used
    // for writting the byte data.
    const transactionSize = this.encodeSizeTransaction(inputs, outputs, signatures).toNumber();
    const buffer = new ArrayBuffer(transactionSize);
    const dataView = new DataView(buffer);
    let currentPos = 0;

    // Tx length
    dataView.setUint32(currentPos, transactionSize, true);
    currentPos += 4;

    // Tx type
    dataView.setUint8(currentPos, transactionType);
    currentPos += 1;

    // Tx innerHash
    const innerHashBytes = this.convertToBytes(innerHash);
    innerHashBytes.forEach(number => {
      dataView.setUint8(currentPos, number);
      currentPos += 1;
    });

    // Tx sigs maxlen check
    if (signatures.length > 65535) {
      throw new Error('Too many signatures.');
    }

    // Tx sigs length
    dataView.setUint32(currentPos, signatures.length, true);
    currentPos += 4;

    // Tx sigs
    (signatures as string[]).forEach(sig => {
      // Copy all bytes
      const binarySig = this.convertToBytes(sig);
      binarySig.forEach(number => {
        dataView.setUint8(currentPos, number);
        currentPos += 1;
      });
    });

    // Tx inputs maxlen check
    if (inputs.length > 65535) {
      throw new Error('Too many inputs.');
    }

    // Tx inputs length
    dataView.setUint32(currentPos, inputs.length, true);
    currentPos += 4;

    // Tx inputs
    inputs.forEach(input => {
      // Copy all bytes
      const binaryInput = this.convertToBytes(input.hash);
      binaryInput.forEach(number => {
        dataView.setUint8(currentPos, number);
        currentPos += 1;
      });
    });

    // Tx outputs maxlen check
    if (outputs.length > 65535) {
      throw new Error('Too many outputs.');
    }

    // Tx outputs length
    dataView.setUint32(currentPos, outputs.length, true);
    currentPos += 4;

    // Create an instance of the base58 decoder.
    const decoder = Base58('123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz');

    // Tx outputs
    outputs.forEach(output => {
      // Decode the address
      const decodedAddress = decoder.decode(output.address);

      // Address version
      dataView.setUint8(currentPos, decodedAddress[20]);
      currentPos += 1;

      // Address Key
      for (let i = 0; i < 20; i++) {
        dataView.setUint8(currentPos, decodedAddress[i]);
        currentPos += 1;
      }

      // Coins
      currentPos = this.setUint64(dataView, currentPos, new BigNumber(output.coins).multipliedBy(1000000).decimalPlaces(0));
      // Hours
      currentPos = this.setUint64(dataView, currentPos, new BigNumber(output.hours));
    });

    //

    return this.convertToHex(buffer);
  }

  /**
   * Calculates the final size, in bytes, that an encoded transaction will have.
   * @param inputs List of all inputs.
   * @param outputs List of all outputs.
   * @param signatures List of all signatures.
   */
  private static encodeSizeTransaction(inputs: HwInput[], outputs: HwOutput[], signatures: string[]): BigNumber {
    let size = new BigNumber(0);

    // Tx length
    size = size.plus(4);

    // Tx type
    size = size.plus(1);

    // Tx innerHash
    size = size.plus(32);

    // Tx sigs
    size = size.plus(4);
    size = size.plus((new BigNumber(65).multipliedBy(signatures.length)));

    // Tx inputs
    size = size.plus(4);
    size = size.plus((new BigNumber(32).multipliedBy(inputs.length)));

    // Tx outputs
    size = size.plus(4);
    size = size.plus((new BigNumber(37).multipliedBy(outputs.length)));

    return size;
  }

  /**
   * Writes an Uint64 value on a DataView.
   * @param dataView DataView in which the value will be written.
   * @param currentPos Position inside the DataView in which the value will be written.
   * @param value Value to be written.
   * @returns The position in which the next value will have to be written on the DataView.
   */
  private static setUint64(dataView: DataView, currentPos: number, value: BigNumber): number {
    let hex = value.toString(16);
    // Make sure the hex string has an even number of characters.
    if (hex.length % 2 !== 0) {
      hex = '0' + hex;
    }

    const bytes = this.convertToBytes(hex);
    for (let i = bytes.length - 1; i >= 0; i--) {
      dataView.setUint8(currentPos, bytes[i]);
      currentPos += 1;
    }

    // Add zeros to fill the remaining space.
    for (let i = 0; i < 8 - bytes.length; i++) {
      dataView.setUint8(currentPos, 0);
      currentPos += 1;
    }

    return currentPos;
  }

  /**
   * Converts a hex string to a byte array.
   * @param hexString String to convert.
   */
  private static convertToBytes(hexString: string): number[] {
    if (hexString.length % 2 !== 0) {
      throw new Error('Invalid hex string.');
    }

    const result: number[] = [];

    for (let i = 0; i < hexString.length; i += 2) {
      result.push(parseInt(hexString.substr(i, 2), 16));
    }

    return result;
  }

  /**
   * Converts an ArrayBuffer to a hex string.
   * @param buffer ArrayBuffer to convert.
   */
  private static convertToHex(buffer: ArrayBuffer) {
    let result = '';

    (new Uint8Array(buffer)).forEach((v) => {
      let val = v.toString(16);
      if (val.length === 0) {
        val = '00';
      } else if (val.length === 1) {
        val = '0' + val;
      }
      result += val;
    });

    return result;
  }
}
