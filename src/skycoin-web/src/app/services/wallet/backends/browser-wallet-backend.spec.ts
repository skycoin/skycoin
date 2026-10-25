import { firstValueFrom } from 'rxjs';

import { Wallet } from '../../../app.datatypes';
import { EncryptionService } from '../../encryption.service';
import { BrowserWalletBackend } from './browser-wallet-backend';

describe('BrowserWalletBackend getSeed', () => {
  const encryption = new EncryptionService();
  const translate = { instant: (key: string) => key };
  const backend = new BrowserWalletBackend(null as any, translate as any, null as any, encryption);

  function wallet(fields: Partial<Wallet>): Wallet {
    return { label: 'w', addresses: [], ...fields };
  }

  it('returns a seed it holds without a password', async () => {
    expect(await firstValueFrom(backend.getSeed(wallet({ seed: 'one two three' }), null))).toBe('one two three');
  });

  it('decrypts a seed kept under the wallet password', async () => {
    const encryptedSeed = await encryption.encrypt('one two three', 'hunter2');
    expect(await firstValueFrom(backend.getSeed(wallet({ encryptedSeed }), 'hunter2'))).toBe('one two three');
    await expectAsync(firstValueFrom(backend.getSeed(wallet({ encryptedSeed }), 'wrong'))).toBeRejectedWithError('Incorrect password');
  });

  it('fails for a wallet whose seed it does not hold', async () => {
    await expectAsync(firstValueFrom(backend.getSeed(wallet({}), null))).toBeRejected();
  });
});
