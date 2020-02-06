import { OperationError, OperationErrorCategories, OperationErrorTypes } from './operation-error';
import { HttpErrorResponse } from '@angular/common/http';

export function parseResponseMessage(body: string): string {
  if (typeof body === 'object') {
    if (body['_body']) {
      body = body['_body'];
    } else {
      body = body + '';
    }
  }

  if (body.indexOf('"error":') !== -1) {
    try {
      body = JSON.parse(body).error.message;
    } catch (e) { }
  }

  if (body.startsWith('400') || body.startsWith('403')) {
    const parts = body.split(' - ', 2);

    return parts.length === 2
      ? parts[1].charAt(0).toUpperCase() + parts[1].slice(1)
      : body;
  }

  return body;
}

export function processServiceError(error: any): OperationError {
  if (error.category && error.type) {
    return error;
  }

  const response = new OperationError();
  response.category = OperationErrorCategories.GeneralApiError;
  response.originalError = error;

  if (!error || typeof error === 'string') {
    response.originalServerErrorMsg = error ? error : '';
    response.translatableErrorMsg = error ? error : 'service.api.unknown-error';
    response.type = OperationErrorTypes.Unknown;

    return response;
  }

  response.originalServerErrorMsg = getErrorMsg(error);
  if (response.originalServerErrorMsg) {
    response.originalServerErrorMsg = parseResponseMessage(response.originalServerErrorMsg);
  }

  const convertedError = error as HttpErrorResponse;
  if (convertedError.status !== null && convertedError.status !== undefined) {
    if (convertedError.status === 0 || convertedError.status === 504) {
      response.type = OperationErrorTypes.NoInternet;
      response.translatableErrorMsg = 'service.api.no-internet-error';
    } else if (convertedError.status === 400 && response.originalServerErrorMsg.toUpperCase().indexOf('Invalid password'.toUpperCase()) !== -1) {
      response.type = OperationErrorTypes.Unauthorized;
      response.translatableErrorMsg = 'service.api.incorrect-password-error';
    }
  }

  if (!response.type) {
    response.type = OperationErrorTypes.Unknown;
    if (response.originalServerErrorMsg) {
      response.translatableErrorMsg = formatUnknownErrorMsg(response.originalServerErrorMsg);
    } else {
      response.translatableErrorMsg = 'service.api.unknown-error';
    }
  }

  return response;
}

function formatUnknownErrorMsg(msg: string) {
  if (!msg || msg.length === 0) {
    return msg;
  }

  msg = msg.trim();

  const firstLetter = msg.substr(0, 1);
  if (firstLetter.toUpperCase() !== firstLetter) {
    msg = firstLetter.toUpperCase() + msg.substr(1, msg.length - 1);
  }

  if (!msg.endsWith('.') && !msg.endsWith(',') && !msg.endsWith(':') && !msg.endsWith(';') && !msg.endsWith('?') && !msg.endsWith('!')) {
    msg = msg + '.';
  }

  return msg;
}

export function getErrorMsg(error: any): string {
  if (error) {
    if (typeof error['_body'] === 'string') {
      return error['_body'];
    }

    if (error.originalServerErrorMsg && typeof error.originalServerErrorMsg === 'string') {
      return error.originalServerErrorMsg;
    } else if (error.error && typeof error.error === 'string') {
      return error.error;
    } else if (error.error && error.error.error && error.error.error.message)  {
      return error.error.error.message;
    } else if (error.error && error.error.error && typeof error.error.error === 'string')  {
      return error.error.error;
    } else if (error.message) {
      return error.message;
    } else if (error._body && error._body.error)  {
      return error._body.error;
    } else {
      try {
        const errorContent = JSON.parse(error._body).error;

        return errorContent;
      } catch (e) {}
    }
  }

  return null;
}

export function redirectToErrorPage(errorCode: number) {
  window.location.assign('assets/error-alert/index.html?' + errorCode);
}
