import { readJSON } from 'karma-read-json';

import { Address, testCases, convertAsciiToHexa } from './utils'

declare var CipherExtras;
declare var Cipher;

declare var TESTING_MODE: string;

describe('CipherProvider Lib', () => {
  const fixturesPath = 'tests/test-fixtures/';
  const addressesFileName = 'many-addresses.golden';
  const inputHashesFileName = 'input-hashes.golden';

  const seedSignaturesFiles = [
    'seed-0000.golden', 'seed-0001.golden', 'seed-0002.golden',
    'seed-0003.golden', 'seed-0004.golden', 'seed-0005.golden',
    'seed-0006.golden', 'seed-0007.golden', 'seed-0008.golden',
    'seed-0009.golden', 'seed-0010.golden'
  ];

  const extensiveMode = '1';

  const testSettings = TESTING_MODE == extensiveMode
    ? { addressCount: 1000, seedFilesCount: 11 }
    : { addressCount: 30, seedFilesCount: 1 };

  describe('generate address', () => {
    const addressFixtureFile = readJSON(fixturesPath + addressesFileName);
    const expectedAddresses = addressFixtureFile.keys.slice(0, testSettings.addressCount);
    let seed = convertAsciiToHexa(atob(addressFixtureFile.seed));
    let generatedAddress;

    testCases(expectedAddresses, (address: any) => {
      it('should generate many address correctly', done => {
        generatedAddress = generateAddress(seed);
        seed = generatedAddress.next_seed;

        const convertedAddress = {
          address: generatedAddress.address,
          public: generatedAddress.public_key,
          secret: generatedAddress.secret_key
        };

        expect(convertedAddress).toEqual(address);
        done();
      });

      it('should pass the verification', done => {
        verifyAddress(generatedAddress);
        done();
      });
    });
  });

  describe('seed signatures', () => {
    const inputHashes = readJSON(fixturesPath + inputHashesFileName).hashes;

    testCases(seedSignaturesFiles.slice(0, testSettings.seedFilesCount), (fileName: string) => {
      describe(`should pass the verification for ${fileName}`, () => {
        let seedKeys;
        let actualAddresses;
        let testData: { signature: string, public_key: string, hash: string, secret_key: string, address: string }[] = [];

        beforeAll(() => {
          const signaturesFixtureFile = readJSON(fixturesPath + fileName);
          const seed = convertAsciiToHexa(atob(signaturesFixtureFile.seed));
          seedKeys = signaturesFixtureFile.keys;

          actualAddresses = generateAddresses(seed, seedKeys);
          testData = getSeedTestData(inputHashes, seedKeys, actualAddresses);
        });

        it('should check number of signatures and hashes', done => {
          const result = seedKeys.some(key => key.signatures.length !== inputHashes.length);

          expect(result).toEqual(false);
          done();
        });

        it('should generate many address correctly', done => {
          actualAddresses.forEach((address, index) => {
            expect(address.address).toEqual(seedKeys[index].address);
            expect(address.public_key).toEqual(seedKeys[index].public);
            expect(address.secret_key).toEqual(seedKeys[index].secret);
          });

          done();
        });

        it('address should pass the verification', done => {
          verifyAddresses(actualAddresses);
          done();
        });

        it(`should verify signature correctly`, done => {
          testData.forEach(data => {
            const result = CipherExtras.VerifyPubKeySignedHash(data.public_key, data.signature, data.hash);
            expect(result).toBeUndefined();
            done();
          });
        });

        it(`should check signature correctly`, done => {
          testData.forEach(data => {
            const result = CipherExtras.VerifyAddressSignedHash(data.address, data.signature, data.hash);
            expect(result).toBeUndefined();
            done();
          });
        });

        it(`should verify signed hash correctly`, done => {
          testData.forEach(data => {
            const result = CipherExtras.VerifySignatureRecoverPubKey(data.signature, data.hash);
            expect(result).toBeUndefined();
            done();
          });
        });

        it(`should generate public key correctly`, done => {
          testData.forEach(data => {
            const pubKey = CipherExtras.PubKeyFromSig(data.signature, data.hash);
            expect(pubKey).toBeTruthy();
            expect(pubKey === data.public_key).toBeTruthy();
            done();
          });
        });

        it(`sign hash should be created`, done => {
          testData.forEach(data => {
            const sig = CipherExtras.SignHash(data.hash, data.secret_key);
            expect(sig).toBeTruthy();
            done();
          });
        });
      });
    });
  });
});

function getSeedTestData(inputHashes, seedKeys, actualAddresses) {
  const data = [];

  for (let seedIndex = 0; seedIndex < seedKeys.length; seedIndex++) {
    for (let hashIndex = 0; hashIndex < inputHashes.length; hashIndex++) {
      data.push({
        signature: seedKeys[seedIndex].signatures[hashIndex],
        public_key: actualAddresses[seedIndex].public_key,
        secret_key: actualAddresses[seedIndex].secret_key,
        address: actualAddresses[seedIndex].address,
        hash: inputHashes[hashIndex]
      });
    }
  }

  return data;
}

function generateAddresses(seed: string, keys: any[]): Address[] {
  return keys.map(() => {
    const generatedAddress = generateAddress(seed);
    seed = generatedAddress.next_seed;

    return generatedAddress;
  });
}

function generateAddress(seed: string): Address {
  const address = Cipher.GenerateAddresses(seed);
  return {
    address: address.Address,
    public_key: address.Public,
    secret_key: address.Secret,
    next_seed: address.NextSeed
  };
}

function verifyAddress(address) {
  const addressFromPubKey = CipherExtras.AddressFromPubKey(address.public_key);
  const addressFromSecKey = CipherExtras.AddressFromSecKey(address.secret_key);

  expect(addressFromPubKey && addressFromSecKey && addressFromPubKey === addressFromSecKey).toBe(true);

  expect(CipherExtras.VerifySeckey(address.secret_key)).toBe(1);
  expect(CipherExtras.VerifyPubkey(address.public_key)).toBe(1);
}

function verifyAddresses(addresses) {
  addresses.forEach(address => verifyAddress(address));
}
