import { HttpErrorResponse } from '@angular/common/http';

import { OperationError, HWOperationResults } from './operation-error';

export function getHwErrorMsg(error: any): string {
  if (error) {
    if (typeof error['_body'] === 'string') {
      return error['_body'];
    } else if (error.originalServerErrorMsg && typeof error.originalServerErrorMsg === 'string') {
      return error.originalServerErrorMsg;
    } else if (error.error && typeof error.error === 'string') {
      return error.error;
    } else if (error.error && error.error.error && error.error.error.message) {
      return error.error.error.message;
    } else if (error.error && error.error.error && typeof error.error.error === 'string') {
      return error.error.error;
    } else if (error.message) {
      return error.message;
    }
  }

  return null;
}

export function processHwServiceError(error: any): OperationError {
  if (error.type) {
    return error;
  }

  const response = new OperationError();
  response.originalError = error;

  if (!error || typeof error === 'string') {
    response.originalServerErrorMsg = error ? error : '';
    response.translatableErrorMsg = error ? error : 'hardware-wallet.errors.generic-error';
    response.type = HWOperationResults.UndefinedError;

    return response;
  }

  response.originalServerErrorMsg = getHwErrorMsg(error);

  const convertedError = error as HttpErrorResponse;
  if (convertedError.status !== null && convertedError.status !== undefined) {
    if (convertedError.status === 0 || convertedError.status === 504) {
      response.type = HWOperationResults.DaemonConnectionError;
      response.translatableErrorMsg = 'hardware-wallet.errors.daemon-connection';

      return response;
    }
  }

  if (!response.type) {
    response.type = HWOperationResults.UndefinedError;
    response.translatableErrorMsg = response.originalServerErrorMsg || 'hardware-wallet.errors.generic-error';
  }

  return response;
}
