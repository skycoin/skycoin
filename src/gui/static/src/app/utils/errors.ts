import { HwWalletService, OperationResults } from '../services/hw-wallet.service';
import { TranslateService } from '@ngx-translate/core';
import { AppConfig } from '../app.config';

export function parseResponseMessage(body: string): string {
  if (typeof body === 'object') {
    if (body['_body']) {
      body = body['_body'];
    } else {
      body = body + '';
    }
  }

  if (body.indexOf('"error":') !== -1) {
    body = JSON.parse(body).error.message;
  }

  if (body.startsWith('400') || body.startsWith('403')) {
    const parts = body.split(' - ', 2);

    return parts.length === 2
      ? parts[1].charAt(0).toUpperCase() + parts[1].slice(1)
      : body;
  }

  return body;
}

export function getHardwareWalletErrorMsg(hwWalletService: HwWalletService, translateService: TranslateService, error: any): string {
  if (!AppConfig.useHwWalletDaemon && !window['ipcRenderer'].sendSync('hwGetDeviceConnectedSync')) {
    return translateService.instant('hardware-wallet.general.error-disconnected');
  }

  if (error.result) {
    if (error.result === OperationResults.FailedOrRefused) {
      return translateService.instant('hardware-wallet.general.refused');
    } else if (error.result === OperationResults.WrongPin) {
      return translateService.instant('hardware-wallet.general.error-incorrect-pin');
    } else if (error.result === OperationResults.IncorrectHardwareWallet) {
      return translateService.instant('hardware-wallet.general.error-incorrect-wallet');
    } else if (error.result === OperationResults.DaemonError) {
      return translateService.instant('hardware-wallet.errors.daemon-connection');
    } else if (error.result === OperationResults.InvalidAddress) {
      return translateService.instant('hardware-wallet.errors.invalid-address');
    } else if (error.result === OperationResults.Timeout) {
      return translateService.instant('hardware-wallet.errors.timeout');
    } else if (error.result === OperationResults.Disconnected) {
      return translateService.instant('hardware-wallet.general.error-disconnected');
    } else if (error.result === OperationResults.NotInBootloaderMode) {
      return translateService.instant('hardware-wallet.errors.not-in-bootloader-mode');
    }
  }

  return translateService.instant('hardware-wallet.general.generic-error');
}
