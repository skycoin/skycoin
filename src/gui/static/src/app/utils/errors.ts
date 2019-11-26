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

export function getHardwareWalletErrorMsg(translateService: TranslateService, error: any, genericError: string = null): string {
  if (!AppConfig.useHwWalletDaemon && !window['ipcRenderer'].sendSync('hwGetDeviceConnectedSync')) {
    if (translateService) {
    return translateService.instant('hardware-wallet.general.error-disconnected');
    } else {
      return 'hardware-wallet.general.error-disconnected';
    }
  }

  let response: string;
  if (error.result) {
    if (error.result === OperationResults.FailedOrRefused) {
      response = 'hardware-wallet.general.refused';
    } else if (error.result === OperationResults.WrongPin) {
      response = 'hardware-wallet.general.error-incorrect-pin';
    } else if (error.result === OperationResults.IncorrectHardwareWallet) {
      response = 'hardware-wallet.general.error-incorrect-wallet';
    } else if (error.result === OperationResults.DaemonError) {
      response = 'hardware-wallet.errors.daemon-connection';
    } else if (error.result === OperationResults.InvalidAddress) {
      response = 'hardware-wallet.errors.invalid-address';
    } else if (error.result === OperationResults.Timeout) {
      response = 'hardware-wallet.errors.timeout';
    } else if (error.result === OperationResults.Disconnected) {
      response = 'hardware-wallet.general.error-disconnected';
    } else if (error.result === OperationResults.NotInBootloaderMode) {
      response = 'hardware-wallet.errors.not-in-bootloader-mode';
    } else if (error.result === OperationResults.PinMismatch) {
      response = 'hardware-wallet.change-pin.pin-mismatch';
    } else if (error.result === OperationResults.WrongWord) {
      response = 'hardware-wallet.restore-seed.error-wrong-word';
    } else if (error.result === OperationResults.InvalidSeed) {
      response = 'hardware-wallet.restore-seed.error-invalid-seed';
    } else if (error.result === OperationResults.WrongSeed) {
      response = 'hardware-wallet.restore-seed.error-wrong-seed';
    } else {
      response = genericError ? genericError : 'hardware-wallet.general.generic-error';
    }
  } else {
    response = genericError ? genericError : 'hardware-wallet.general.generic-error';
  }

  if (translateService) {
    return translateService.instant(response);
  } else {
    return response;
  }
}
