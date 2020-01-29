import { Observable, of } from 'rxjs';
import { map, catchError, mergeMap, first } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletBase } from './wallet-objects';
import { HwWalletService } from '../hw-wallet.service';
import { HttpClient } from '@angular/common/http';
import { AppConfig } from '../../app.config';

export enum HwSecurityWarnings {
  NeedsBackup,
  NeedsPin,
  FirmwareVersionNotVerified,
  OutdatedFirmware,
}

export interface HwFeaturesResponse {
  features: any;
  securityWarnings: HwSecurityWarnings[];
  walletNameUpdated: boolean;
}

@Injectable()
export class HardwareWalletService {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private hwWalletService: HwWalletService,
    private http: HttpClient,
  ) { }

  getFeaturesAndUpdateData(wallet: WalletBase): Observable<HwFeaturesResponse> {
    if (!wallet || wallet.isHardware) {

      let lastestFirmwareVersion: string;

      return this.http.get(AppConfig.urlForHwWalletVersionChecking, { responseType: 'text' }).pipe(
      catchError(() => of(null)),
      mergeMap((res: any) => {
        if (res) {
          lastestFirmwareVersion = res;
        } else {
          lastestFirmwareVersion = null;
        }

        return this.hwWalletService.getFeatures();
      }),
      map(result => {
        let lastestFirmwareVersionReaded = false;
        let firmwareUpdated = false;

        if (lastestFirmwareVersion) {
          lastestFirmwareVersion = lastestFirmwareVersion.trim();
          const versionParts = lastestFirmwareVersion.split('.');

          if (versionParts.length === 3) {
            lastestFirmwareVersionReaded = true;

            const numVersionParts = versionParts.map(value => Number.parseInt(value.replace(/\D/g, ''), 10));

            const devMajorVersion = result.rawResponse.fw_major;
            const devMinorVersion = result.rawResponse.fw_minor;
            const devPatchVersion = result.rawResponse.fw_patch;

            if (devMajorVersion > numVersionParts[0]) {
              firmwareUpdated = true;
            } else {
              if (devMajorVersion === numVersionParts[0]) {
                if (devMinorVersion > numVersionParts[1]) {
                  firmwareUpdated = true;
                } else {
                  if (devMinorVersion === numVersionParts[1] && devPatchVersion >= numVersionParts[2]) {
                    firmwareUpdated = true;
                  }
                }
              }
            }
          }
        }

        const warnings: HwSecurityWarnings[] = [];
        let hasHwSecurityWarnings = false;

        if (result.rawResponse.needs_backup) {
          warnings.push(HwSecurityWarnings.NeedsBackup);
          hasHwSecurityWarnings = true;
        }
        if (!result.rawResponse.pin_protection) {
          warnings.push(HwSecurityWarnings.NeedsPin);
          hasHwSecurityWarnings = true;
        }

        if (!lastestFirmwareVersionReaded) {
          warnings.push(HwSecurityWarnings.FirmwareVersionNotVerified);
        } else {
          if (!firmwareUpdated) {
            warnings.push(HwSecurityWarnings.OutdatedFirmware);
            hasHwSecurityWarnings = true;
          }
        }

        let walletNameUpdated = false;

        if (wallet) {
          let changesMade = false;

          const deviceLabel = result.rawResponse.label ? result.rawResponse.label : (result.rawResponse.deviceId ? result.rawResponse.deviceId : result.rawResponse.device_id);
          if (wallet.label !== deviceLabel) {
            wallet.label = deviceLabel;
            walletNameUpdated = true;
            changesMade = true;
          }
          if (wallet.hasHwSecurityWarnings !== hasHwSecurityWarnings) {
            wallet.hasHwSecurityWarnings = hasHwSecurityWarnings;
            changesMade = true;
          }

          if (changesMade) {
            this.walletsAndAddressesService.informValuesUpdated(wallet);
          }
        }

        const response = {
          features: result.rawResponse,
          securityWarnings: warnings,
          walletNameUpdated: walletNameUpdated,
        };

        return response;
      }));
    } else {
      return null;
    }
  }

  setAddressConfirmed(wallet: WalletBase, addressIndex: number) {
    wallet.addresses[addressIndex].confirmed = true;
    this.walletsAndAddressesService.informValuesUpdated(wallet);
  }
}
