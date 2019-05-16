import * as Base58 from 'base-x';
import BigNumber from 'bignumber.js';

// Base for a class for encoding Skycoin transactions. Currently only performs internal tests.
// It must be adapted to work well with WalletService and thus be able to remove the skycoin-lite lib.
//
// To-do: - Modify the class to accept any transaction, not just the test transaction.
//        - The input transaction should be in a well defined object, not the "any" type, if possible.
//        - Create a test suite using old transaction from the blockchain, to guarantee the correct operation of the procedure.
export class TxEncoder {

  // This is the transaction in the bklock #200
  // The expected encoded transaction is:
  // dc0000000008832259284fe4965625e2ef97d1ff3b40d7832b159b3d5369fc086ebb95479f01000000820dc5d47540b0978818356b512731ae4751
  // 7c061a1b660e7fd47e8f3d6420700377fd2ec04b618504cd3dadf642111df2f6c1f6edf2a4067fd460c69e8eb07301010000006c34016037cd1762
  // 2846e71bc635914d4d8f256c147aa5a0b84a896e832294800200000000bb202804300d62db2fcfae5ee720eeb28493e3f8003ef5e9050000001ebf
  // 770000000000003be2537f8c0893fddcddc878518f38ea493d949e00ca9a3b000000001ebf770000000000
  private testTx = `{
    "length":220,
    "type":0,
    "inner_hash":"08832259284fe4965625e2ef97d1ff3b40d7832b159b3d5369fc086ebb95479f",
    "sigs":["820dc5d47540b0978818356b512731ae47517c061a1b660e7fd47e8f3d6420700377fd2ec04b618504cd3dadf642111df2f6c1f6edf2a4067fd460c69e8eb07301"],
    "inputs":[
      {
        "uxid":"6c34016037cd17622846e71bc635914d4d8f256c147aa5a0b84a896e83229480"
      }
    ],
    "outputs":[
      {
        "uxid":"6e4110a8ed6f2b8b8772516466032a99b4851de65cf9ce1b5c5673946b7408a9",
        "address":"2JJ8pgq8EDAnrzf9xxBJapE2qkYLefW4uF8",
        "coins":"25400.000000",
        "hours":"7847710"
      },
      {
        "uxid":"50e534ebc9c3f0b99461ad70b01d415eabfc046e824a5d1ba46854c913928612",
        "address":"R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
        "coins":"1000.000000",
        "hours":"7847710"
      }
    ]
  }`;

  test() {
    const transaction = JSON.parse(this.testTx);

    const buffer = new ArrayBuffer(this.encodeSizeTransaction(transaction).toNumber());
    const dataView = new DataView(buffer);
    let currentPos = 0;

    // Tx length
    dataView.setUint32(currentPos, transaction.length, true);
    currentPos += 4;

    // Tx type
    dataView.setUint8(currentPos, transaction.type);
    currentPos += 1;

    // Tx innerHash
    const innerHash = this.convertToBytes(transaction.inner_hash);
    innerHash.forEach(number => {
      dataView.setUint8(currentPos, number);
      currentPos += 1;
    });

    // Tx sigs maxlen check
    if (transaction.sigs.length > 65535) {
      throw new Error('Too many signatures.');
    }

    // Tx sigs length
    dataView.setUint32(currentPos, transaction.sigs.length, true);
    currentPos += 4;

    // Tx sigs
    (transaction.sigs as string[]).forEach(sig => {
      // Copy all bytes
      const binarySig = this.convertToBytes(sig);
      binarySig.forEach(number => {
        dataView.setUint8(currentPos, number);
        currentPos += 1;
      });
    });

    // Tx inputs maxlen check
    if (transaction.sigs.length > 65535) {
      throw new Error('Too many inputs.');
    }

    // Tx inputs length
    dataView.setUint32(currentPos, transaction.inputs.length, true);
    currentPos += 4;

    // Tx inputs
    (transaction.inputs as any[]).forEach(input => {
      // Copy all bytes
      const binaryInput = this.convertToBytes(input.uxid);
      binaryInput.forEach(number => {
        dataView.setUint8(currentPos, number);
        currentPos += 1;
      });
    });

    // Tx outputs maxlen check
    if (transaction.sigs.length > 65535) {
      throw new Error('Too many outputs.');
    }

    // Tx outputs length
    dataView.setUint32(currentPos, transaction.outputs.length, true);
    currentPos += 4;

    const decoder = Base58('123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz');

    // Tx outputs
    (transaction.outputs as any[]).forEach(output => {
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
      currentPos = this.setUint64(dataView, currentPos, new BigNumber(output.coins).multipliedBy(1000000));
      // Hours
      currentPos = this.setUint64(dataView, currentPos, new BigNumber(output.hours));
    });

    //

    alert(this.convertToHex(buffer));
    console.log(this.convertToHex(buffer));
  }

  encodeSizeTransaction(transaction: any): BigNumber {
    let size = new BigNumber(0);

    // Tx length
    size = size.plus(4);

    // Tx type
    size = size.plus(1);

    // Tx innerHash
    size = size.plus(32);

    // Tx sigs
    size = size.plus(4);
    size = size.plus((new BigNumber(65).multipliedBy(transaction.sigs.length)));

    // Tx inputs
    size = size.plus(4);
    size = size.plus((new BigNumber(32).multipliedBy(transaction.inputs.length)));

    // Tx outputs
    size = size.plus(4);
    size = size.plus((new BigNumber(37).multipliedBy(transaction.outputs.length)));

    return size;
  }

  private setUint64(dataView: DataView, currentPos: number, value: BigNumber): number {
    let hex = value.toString(16);
    if (hex.length % 2 !== 0) {
      hex = '0' + hex;
    }

    const bytes = this.convertToBytes(hex);
    for (let i = bytes.length - 1; i >= 0; i--) {
      dataView.setUint8(currentPos, bytes[i]);
      currentPos += 1;
    }

    for (let i = 0; i < 8 - bytes.length; i++) {
      dataView.setUint8(currentPos, 0);
      currentPos += 1;
    }

    return currentPos;
  }

  private convertToBytes(hexString: string): number[] {
    if (hexString.length % 2 !== 0) {
      throw new Error('Invalid hex string.');
    }

    const result: number[] = [];

    for (let i = 0; i < hexString.length; i += 2) {
      result.push(parseInt(hexString.substr(i, 2), 16));
    }

    return result;
  }

  private convertToHex(buffer: ArrayBuffer) {
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
