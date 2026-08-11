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
 * The golden files also carry signatures, but src/skycoin-lite/wasm/main_wasm.go
 * publishes only generateAddress, prepareTransaction and
 * prepareTransactionWithSignatures. The signature helpers this spec used to
 * reach through window.SkycoinCipherExtras belong to the GopherJS build
 * (src/skycoin-lite/skycoin/skycoin.go, exercised by
 * src/skycoin-lite/js/tests/cipher.spec.ts), not to the wasm the wallet ships,
 * so those assertions cannot run here.
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

declare global {
  interface Window {
    SkycoinCipher: SkycoinCipherLib;
  }
}

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
    });
  });

  describe('generate address', () => {
    const addressFixtureFile = readFixture(fixturesPath + addressesFileName);
    const expectedAddresses = addressFixtureFile.keys.slice(0, testSettings.addressCount);
    let seed = convertAsciiToHexa(atob(addressFixtureFile.seed));

    testCases(expectedAddresses, (address: any) => {
      it('should generate many address correctly', () => {
        const generatedAddress = generateAddress(seed);
        seed = generatedAddress.nextSeed;

        // Compared field by field rather than as a whole: the golden files
        // also carry bitcoin_address, which the browser cipher does not derive.
        expect(generatedAddress.address.address).toEqual(address.address);
        expect(generatedAddress.address.public_key).toEqual(address.public);
        expect(generatedAddress.address.secret_key).toEqual(address.secret);
      });
    });
  });

  describe('seed signatures', () => {
    const inputHashes: string[] = readFixture(fixturesPath + inputHashesFileName).hashes;

    testCases(seedSignaturesFiles.slice(0, testSettings.seedFilesCount), (fileName: string) => {
      describe(`should pass the verification for ${fileName}`, () => {
        let seedKeys: GoldenKey[] = [];
        let actualAddresses: Address[] = [];

        beforeAll(() => {
          const signaturesFixtureFile = readFixture(fixturesPath + fileName);
          const seed = convertAsciiToHexa(atob(signaturesFixtureFile.seed));
          seedKeys = signaturesFixtureFile.keys;

          actualAddresses = generateAddresses(seed, seedKeys);
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
      });
    });
  });
});

function generateAddresses(seed: string, keys: GoldenKey[]): Address[] {
  return keys.map(() => {
    const generatedAddress = generateAddress(seed);
    seed = generatedAddress.nextSeed;

    return generatedAddress.address;
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
