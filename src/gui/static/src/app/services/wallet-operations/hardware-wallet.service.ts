import { Observable, of } from 'rxjs';
import { map, catchError, mergeMap } from 'rxjs/operators';
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

import { WalletsAndAddressesService } from './wallets-and-addresses.service';
import { WalletBase } from './wallet-objects';
import { HwWalletService } from '../hw-wallet.service';
import { AppConfig } from '../../app.config';

/**
 * List with the security problems that can be detected in a hw wallet.
 */
export enum HwSecurityWarnings {
  /**
   * The user has not made a backup of the seed of the device.
   */
  NeedsBackup = 'NeedsBackup',
  /**
   * The device does not have PIN code protection activated.
   */
  NeedsPin = 'NeedsPin',
  /**
   * It was not possible to connect to the remote server to know if the device is running
   * an outdated firmware.
   */
  FirmwareVersionNotVerified = 'FirmwareVersionNotVerified',
  /**
   * The device is running an outdated firmware.
   */
  OutdatedFirmware = 'OutdatedFirmware',
}

export interface HwFeaturesResponse {
  /**
   * Features returned by the device.
   */
  features: any;
  /**
   * Security warnings found during the operation.
   */
  securityWarnings: HwSecurityWarnings[];
  /**
   * During the operation it was detected that the device is showing a label which is not
   * equal to the one known by the app, so the name shown on the app was updated to match
   * the one on the device.
   */
  walletNameUpdated: boolean;
}

/**
 * Allows to perform operations related to a hardware wallet.
 */
@Injectable()
export class HardwareWalletService {
  constructor(
    private walletsAndAddressesService: WalletsAndAddressesService,
    private hwWalletService: HwWalletService,
    private http: HttpClient,
  ) { }

  /**
   * Gets the features of the connected hw wallet and also updates the label and security warnings
   * of a WalletBase instance with the data obtained from the currently connected device.
   * @param wallet WalletBase instance to be updated. If null, the function will still work, it
   * just will not update anything.
   */
  getFeaturesAndUpdateData(wallet: WalletBase): Observable<HwFeaturesResponse> {
    if (!wallet || wallet.isHardware) {

      let lastestFirmwareVersion: string;

      // Get the number of the most recent firmware version.
      return this.http.get(AppConfig.urlForHwWalletVersionChecking, { responseType: 'text' }).pipe(
      catchError(() => of(null)),
      mergeMap((res: any) => {
        // If it was not possible to get the number of the most recent firmware version,
        // continue anyways.
        if (res) {
          lastestFirmwareVersion = res;
        } else {
          lastestFirmwareVersion = null;
        }

        // Get the features.
        return this.hwWalletService.getFeatures();
      }),
      map(result => {
        let lastestFirmwareVersionReaded = false;
        let firmwareUpdated = false;

        // Compare the version number of the firmware of the connected device and the one obtained from the server.
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

        // Create the security warnings list.
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

        // If a wallet was provided, update the label and security warnings, if needed.
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

        // Prepare the response.
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

  /**
   * Asks the user to confirm an addresses on the connected device. This allows the user
   * to confirm that the address displayed by this app is equal to the one on the device.
   * @param wallet Wallet to check.
   * @param addressIndex Index of the addresses on the provided wallet.
   */
  confirmAddress(wallet: WalletBase, addressIndex: number): Observable<void> {
    return this.hwWalletService.checkIfCorrectHwConnected(wallet.id).pipe(
      mergeMap(() => this.hwWalletService.confirmAddress(addressIndex)),
      map(() => {
        // If the user confirms the operation, update the local data.
        wallet.addresses[addressIndex].confirmed = true;
        this.walletsAndAddressesService.informValuesUpdated(wallet);
      }),
    );
  }

  /**
   * Asks the device to change the label shown on its physical scren. It also changes the
   * label/name saved on the local wallet object.
   * @param wallet Wallet to modify.
   * @param newLabel New label or name.
   */
  changeLabel(wallet: WalletBase, newLabel: string): Observable<void> {
    return this.hwWalletService.checkIfCorrectHwConnected(wallet.id).pipe(
      mergeMap(() => this.hwWalletService.changeLabel(newLabel)),
      map(() => {
        // If the user confirms the operation, update the local data.
        wallet.label = newLabel;
        this.walletsAndAddressesService.informValuesUpdated(wallet);
      }),
    );
  }
}
