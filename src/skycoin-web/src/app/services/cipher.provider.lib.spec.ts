import { testCases } from '../utils/jasmine-utils';
import { Address } from '../app.datatypes';
import { convertAsciiToHexa } from '../utils/converters';
import { GenerateAddressResponse } from './cipher.provider.js';
import { loadCipherWasm } from '../utils/wasm-test-utils';

/**
 * Checks the browser cipher against the Go cipher test vectors.
 *
 * assets/scripts/skycoin-lite.wasm derives every address this wallet shows, and
 * a break in it is silent: the UI renders and the addresses are simply wrong.
 * The fixtures are src/cipher/testsuite/testdata, the same golden files the Go
 * implementation is checked against, staged here by
 * scripts/stage-cipher-fixtures.js. Sharing them is the point — the two
 * implementations agreeing is what makes a wallet restored in the browser hold
 * the same coins as one restored by the node.
 *
 * The signature suites go through window.SkycoinCipherExtras, the verification
 * helpers src/skycoin-lite/wasm/main_wasm.go publishes alongside the wallet's
 * own entry points. Each verify* answers null when the check passes and an error
 * string when it does not.
 */

const fixturesPath = '/cipher-fixtures/';

/** A key entry of a cipher testsuite golden file. */
interface GoldenKey {
  address: string;
  public: string;
  secret: string;
  signatures: string[];
}

/** The subset of the wasm cipher this spec exercises. */
interface SkycoinCipherLib {
  generateAddress(seed: string): { address: string; public: string; secret: string; nextSeed: string };
}

/** The verification helpers, which answer null when the check passes. */
interface SkycoinCipherExtrasLib {
  addressFromPubKey(publicKey: string): string;
  addressFromSecKey(secretKey: string): string;
  verifySeckey(secretKey: string): string | null;
  verifyPubkey(publicKey: string): string | null;
  verifyPubKeySignedHash(publicKey: string, signature: string, hash: string): string | null;
  verifyAddressSignedHash(address: string, signature: string, hash: string): string | null;
  verifySignatureRecoverPubKey(signature: string, hash: string): string | null;
  pubKeyFromSig(signature: string, hash: string): string;
  signHash(hash: string, secretKey: string): string;
}

declare global {
  interface Window {
    SkycoinCipher: SkycoinCipherLib;
    SkycoinCipherExtras: SkycoinCipherExtrasLib;
  }
}

/** Addresses the cipher generates always carry both keys, unlike Address. */
type GeneratedAddress = Address & Required<Pick<Address, 'public_key' | 'secret_key'>>;

// karma-read-json resolves paths against /base/, which @angular/build:karma does
// not serve, so the fixtures are declared as assets on the test target instead.
// Reading them has to be synchronous because the test cases are generated at
// describe() time.
function readFixture(file: string): any {
  const request = new XMLHttpRequest();
  request.open('GET', file, false);
  request.send(null);
  if (request.status !== 200) {
    throw new Error(`could not read ${file}: HTTP ${request.status}. ` +
      'Run scripts/stage-cipher-fixtures.js first (npm test does it for you).');
  }

  return JSON.parse(request.responseText);
}

describe('CipherProvider Lib', () => {
  const addressesFileName = 'many-addresses.golden';
  const inputHashesFileName = 'input-hashes.golden';

  const seedSignaturesFiles = [
    'seed-0000.golden', 'seed-0001.golden', 'seed-0002.golden',
    'seed-0003.golden', 'seed-0004.golden', 'seed-0005.golden',
    'seed-0006.golden', 'seed-0007.golden', 'seed-0008.golden',
    'seed-0009.golden', 'seed-0010.golden'
  ];

  const testSettings = { addressCount: 1000, seedFilesCount: 11 };

  // The suites below call the cipher directly, so none of them may depend on
  // the Initialization spec having run first.
  beforeAll(async () => {
    await loadCipherWasm();
  }, 60000);

  describe('Initialization', () => {
    it('should be initialized', () => {
      expect(window.SkycoinCipher).toBeTruthy();
      expect(window.SkycoinCipher.generateAddress).toEqual(jasmine.any(Function));
      expect(window.SkycoinCipherExtras).toBeTruthy();
      expect(window.SkycoinCipherExtras.signHash).toEqual(jasmine.any(Function));
    });
  });

  // Everything under liteclient reports failure by panicking. main_wasm.go has
  // to turn that into the {error} result the wallet is written against; when it
  // recovered without setting a return value the caller got null and crashed
  // reading .error off it.
  describe('error handling', () => {
    it('reports a seed it cannot parse as an error', () => {
      const result = window.SkycoinCipher.generateAddress('not a hex seed') as any;

      expect(result).not.toBeNull();
      expect(result.error).toBeTruthy();
    });

    it('reports transaction inputs it cannot parse as an error', () => {
      const result = (window.SkycoinCipher as any).prepareTransaction('not json', 'not json');

      expect(result).not.toBeNull();
      expect(result.error).toBeTruthy();
    });

    it('reports a missing argument as an error', () => {
      const result = (window.SkycoinCipher as any).generateAddress();

      expect(result).not.toBeNull();
      expect(result.error).toBeTruthy();
    });

    it('reports a signature it cannot verify as an error', () => {
      expect(window.SkycoinCipherExtras.verifySeckey('00')).toBeTruthy();
    });
  });

  describe('generate address', () => {
    const addressFixtureFile = readFixture(fixturesPath + addressesFileName);
    const expectedAddresses = addressFixtureFile.keys.slice(0, testSettings.addressCount);
    let seed = convertAsciiToHexa(atob(addressFixtureFile.seed));
    let generated: Address;

    testCases(expectedAddresses, (address: any) => {
      it('should generate many address correctly', () => {
        const generatedAddress = generateAddress(seed);
        seed = generatedAddress.nextSeed;
        generated = generatedAddress.address;

        // Compared field by field rather than as a whole: the golden files
        // also carry bitcoin_address, which the browser cipher does not derive.
        expect(generatedAddress.address.address).toEqual(address.address);
        expect(generatedAddress.address.public_key).toEqual(address.public);
        expect(generatedAddress.address.secret_key).toEqual(address.secret);
      });

      it('should pass the verification', () => {
        verifyAddress(generated as GeneratedAddress);
      });
    });
  });

  describe('seed signatures', () => {
    const inputHashes: string[] = readFixture(fixturesPath + inputHashesFileName).hashes;

    testCases(seedSignaturesFiles.slice(0, testSettings.seedFilesCount), (fileName: string) => {
      describe(`should pass the verification for ${fileName}`, () => {
        let seedKeys: GoldenKey[] = [];
        let actualAddresses: GeneratedAddress[] = [];
        let testData: SignatureCase[] = [];

        beforeAll(() => {
          const signaturesFixtureFile = readFixture(fixturesPath + fileName);
          const seed = convertAsciiToHexa(atob(signaturesFixtureFile.seed));
          seedKeys = signaturesFixtureFile.keys;

          actualAddresses = generateAddresses(seed, seedKeys);
          testData = getSeedTestData(inputHashes, seedKeys, actualAddresses);
        });

        it('should check number of signatures and hashes', () => {
          const result = seedKeys.some((key: GoldenKey) => key.signatures.length !== inputHashes.length);

          expect(result).toEqual(false);
        });

        it('should generate many address correctly', () => {
          expect(actualAddresses.length).toEqual(seedKeys.length);

          actualAddresses.forEach((address: Address, index: number) => {
            expect(address.address).toEqual(seedKeys[index].address);
            expect(address.public_key).toEqual(seedKeys[index].public);
            expect(address.secret_key).toEqual(seedKeys[index].secret);
          });
        });

        it('address should pass the verification', () => {
          actualAddresses.forEach(address => verifyAddress(address));
        });

        it('should verify signature correctly', () => {
          testData.forEach(data => {
            expect(window.SkycoinCipherExtras.verifyPubKeySignedHash(data.public_key, data.signature, data.hash)).toBeNull();
          });
        });

        it('should check signature correctly', () => {
          testData.forEach(data => {
            expect(window.SkycoinCipherExtras.verifyAddressSignedHash(data.address, data.signature, data.hash)).toBeNull();
          });
        });

        it('should verify signed hash correctly', () => {
          testData.forEach(data => {
            expect(window.SkycoinCipherExtras.verifySignatureRecoverPubKey(data.signature, data.hash)).toBeNull();
          });
        });

        it('should generate public key correctly', () => {
          testData.forEach(data => {
            expect(window.SkycoinCipherExtras.pubKeyFromSig(data.signature, data.hash)).toEqual(data.public_key);
          });
        });

        it('sign hash should be created', () => {
          testData.forEach(data => {
            expect(window.SkycoinCipherExtras.signHash(data.hash, data.secret_key)).toBeTruthy();
          });
        });
      });
    });
  });
});

/** One golden signature paired with the key and hash it belongs to. */
interface SignatureCase {
  signature: string;
  public_key: string;
  secret_key: string;
  address: string;
  hash: string;
}

function getSeedTestData(inputHashes: string[], seedKeys: GoldenKey[], actualAddresses: GeneratedAddress[]): SignatureCase[] {
  const data: SignatureCase[] = [];

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

function verifyAddress(address: GeneratedAddress) {
  const addressFromPubKey = window.SkycoinCipherExtras.addressFromPubKey(address.public_key);
  const addressFromSecKey = window.SkycoinCipherExtras.addressFromSecKey(address.secret_key);

  expect(addressFromPubKey).toEqual(address.address);
  expect(addressFromSecKey).toEqual(address.address);

  expect(window.SkycoinCipherExtras.verifySeckey(address.secret_key)).toBeNull();
  expect(window.SkycoinCipherExtras.verifyPubkey(address.public_key)).toBeNull();
}

function generateAddresses(seed: string, keys: GoldenKey[]): GeneratedAddress[] {
  return keys.map(() => {
    const generatedAddress = generateAddress(seed);
    seed = generatedAddress.nextSeed;

    return generatedAddress.address as GeneratedAddress;
  });
}

function generateAddress(seed: string): GenerateAddressResponse {
  const address = window.SkycoinCipher.generateAddress(seed);

  return {
    address: {
      address: address.address,
      public_key: address.public,
      secret_key: address.secret
    },
    nextSeed: address.nextSeed
  };
}
