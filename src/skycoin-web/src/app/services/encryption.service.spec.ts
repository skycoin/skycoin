import { TestBed } from '@angular/core/testing';
import { EncryptionService } from './encryption.service';

describe('EncryptionService', () => {
  let service: EncryptionService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [EncryptionService]
    });
    service = TestBed.inject(EncryptionService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should encrypt and decrypt a string correctly', async () => {
    const plaintext = 'test seed phrase for wallet encryption';
    const password = 'mySecurePassword123';

    const encrypted = await service.encrypt(plaintext, password);
    expect(encrypted).toBeTruthy();
    expect(encrypted).not.toEqual(plaintext);

    const decrypted = await service.decrypt(encrypted, password);
    expect(decrypted).toEqual(plaintext);
  });

  it('should produce different ciphertexts for the same plaintext (random salt/iv)', async () => {
    const plaintext = 'same input text';
    const password = 'password';

    const encrypted1 = await service.encrypt(plaintext, password);
    const encrypted2 = await service.encrypt(plaintext, password);

    expect(encrypted1).not.toEqual(encrypted2);

    // Both should decrypt to the same value
    expect(await service.decrypt(encrypted1, password)).toEqual(plaintext);
    expect(await service.decrypt(encrypted2, password)).toEqual(plaintext);
  });

  it('should fail to decrypt with wrong password', async () => {
    const plaintext = 'secret data';
    const encrypted = await service.encrypt(plaintext, 'correctPassword');

    try {
      await service.decrypt(encrypted, 'wrongPassword');
      fail('Expected decrypt to throw with wrong password');
    } catch (e) {
      expect(e).toBeTruthy();
    }
  });

  it('should handle empty string plaintext', async () => {
    const password = 'password';
    const encrypted = await service.encrypt('', password);
    const decrypted = await service.decrypt(encrypted, password);
    expect(decrypted).toEqual('');
  });

  it('should handle unicode characters', async () => {
    const plaintext = '日本語テスト 🔑 émojis and àccénts';
    const password = 'pässwörd';

    const encrypted = await service.encrypt(plaintext, password);
    const decrypted = await service.decrypt(encrypted, password);
    expect(decrypted).toEqual(plaintext);
  });

  it('should produce valid base64 output', async () => {
    const encrypted = await service.encrypt('test', 'password');
    // Base64 should only contain valid characters
    expect(encrypted).toMatch(/^[A-Za-z0-9+/]+=*$/);
  });

  it('should produce output with correct structure (16 salt + 12 iv + ciphertext)', async () => {
    const encrypted = await service.encrypt('test', 'password');
    const packed = Uint8Array.from(atob(encrypted), c => c.charCodeAt(0));
    // salt(16) + iv(12) + ciphertext (at least 1 byte + 16 byte GCM tag)
    expect(packed.length).toBeGreaterThanOrEqual(16 + 12 + 17);
  });
});
