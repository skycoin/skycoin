import { Injectable } from '@angular/core';

@Injectable()
export class EncryptionService {

  private textEncoder = new TextEncoder();
  private textDecoder = new TextDecoder();

  async encrypt(plaintext: string, password: string): Promise<string> {
    const salt = crypto.getRandomValues(new Uint8Array(16));
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const key = await this.deriveKey(password, salt);

    const ciphertext = await crypto.subtle.encrypt(
      { name: 'AES-GCM', iv: iv as BufferSource },
      key,
      this.textEncoder.encode(plaintext) as BufferSource
    );

    // Pack salt + iv + ciphertext into a single base64 string
    const packed = new Uint8Array(salt.length + iv.length + ciphertext.byteLength);
    packed.set(salt, 0);
    packed.set(iv, salt.length);
    packed.set(new Uint8Array(ciphertext), salt.length + iv.length);

    return btoa(String.fromCharCode(...packed));
  }

  async decrypt(encoded: string, password: string): Promise<string> {
    const packed = Uint8Array.from(atob(encoded), c => c.charCodeAt(0));

    const salt = packed.slice(0, 16);
    const iv = packed.slice(16, 28);
    const ciphertext = packed.slice(28);

    const key = await this.deriveKey(password, salt);

    const plaintext = await crypto.subtle.decrypt(
      { name: 'AES-GCM', iv: iv as BufferSource },
      key,
      ciphertext as BufferSource
    );

    return this.textDecoder.decode(plaintext);
  }

  private async deriveKey(password: string, salt: Uint8Array): Promise<CryptoKey> {
    const keyMaterial = await crypto.subtle.importKey(
      'raw',
      this.textEncoder.encode(password) as BufferSource,
      'PBKDF2',
      false,
      ['deriveKey']
    );

    return crypto.subtle.deriveKey(
      { name: 'PBKDF2', salt: salt as BufferSource, iterations: 100000, hash: 'SHA-256' },
      keyMaterial,
      { name: 'AES-GCM', length: 256 },
      false,
      ['encrypt', 'decrypt']
    );
  }
}
