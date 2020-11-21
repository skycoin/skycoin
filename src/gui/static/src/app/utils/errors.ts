/**
 * This file contains functions for processing the errors and make it easier to work with them.
 */
import { HttpErrorResponse } from '@angular/common/http';

import { OperationError, OperationErrorTypes } from './operation-error';

/**
 * Prepares an error msg to be displayed on the UI.
 */
export function processErrorMsg(msg: string): string {
  if (!msg || msg.length === 0) {
    return msg;
  }

  // Check if the string is from a known error. If it is, use the correct error string.
  const knownErrorMsg = checkIfKnownErrorStrings(msg);
  if (knownErrorMsg) {
    return knownErrorMsg;
  }

  // Some times an error message could be in fact a JSON string. In those cases, the real
  // error msg is inside the "error.message" property.
  if (msg.indexOf('"error":') !== -1) {
    try {
      msg = JSON.parse(msg).error.message;
    } catch (e) { }
  }

  // Remove unnecessary error codes.
  if (msg.startsWith('400') || msg.startsWith('403')) {
    const parts = msg.split(' - ', 2);

    msg = parts.length === 2 ? parts[1] : msg;
  }

  // The msg will start with an uppercase letter and end with a period.
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

/**
 * Process an error and creates an OperationError instance from it. It can successfully
 * process various types of errors (connection errors, operation errors and more), strings
 * and even OperationError intances, so it is relatively safe to use this function to
 * process almost all errors before using them, even the ones returned by services which
 * normally return OperationError instances in case of error, just to be sure to have a standard
 * OperationError object to work with.
 * @param error Error to process.
 */
export function processServiceError(error: any): OperationError {
  // Check if the provided error is already an OperationError instance.
  if (error.type) {
    return error;
  }

  const response = new OperationError();
  response.originalError = error;

  // Check if the provided error is empty or a string.
  if (!error || typeof error === 'string') {
    response.originalServerErrorMsg = error ? error : '';
    response.translatableErrorMsg = error ? error : 'service.api.unknown-error';
    response.type = OperationErrorTypes.Unknown;

    return response;
  }

  // Extract the error msg from the provided error param.
  response.originalServerErrorMsg = getErrorMsg(error);

  // Check if the provided error is a known API error.
  const convertedError = error as HttpErrorResponse;
  if (convertedError.status !== null && convertedError.status !== undefined) {
    if (convertedError.status === 0 || convertedError.status === 504) {
      response.type = OperationErrorTypes.NoInternet;
      response.translatableErrorMsg = 'service.api.no-internet-error';
    } else if (convertedError.status === 0 || convertedError.status === 403) {
      response.type = OperationErrorTypes.ApiDisabled;
      response.translatableErrorMsg = 'service.api.api-disabled-error';
    } else if (convertedError.status === 400 && response.originalServerErrorMsg.toUpperCase().indexOf('Invalid password'.toUpperCase()) !== -1) {
      response.type = OperationErrorTypes.Unauthorized;
      response.translatableErrorMsg = 'service.api.incorrect-password-error';
    }
  }

  // Use defaults and process the error msg if needed.
  if (!response.type) {
    response.type = OperationErrorTypes.Unknown;
    if (response.originalServerErrorMsg) {
      response.translatableErrorMsg = processErrorMsg(response.originalServerErrorMsg);
    } else {
      response.translatableErrorMsg = 'service.api.unknown-error';
    }
  }

  return response;
}

/**
 * Tries to get the error msg of an error object.
 * @param error Error to process.
 * @returns The error msg, or null, if it was not possible to retrieve the error msg.
 */
export function getErrorMsg(error: any): string {
  if (error) {
    // Check different posibilities, testing a normal error object and different
    // known ubications in which the error msg could be located.
    if (typeof error['_body'] === 'string') {
      return error['_body'];
    } else if (error.originalServerErrorMsg && typeof error.originalServerErrorMsg === 'string') {
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

/**
 * Makes the browser navigate to the error page.
 * @param errorCode Error code the error page must show. Consult the code of the error page to
 * for more info about the codes.
 */
export function redirectToErrorPage(errorCode: number) {
  window.location.assign('assets/error-alert/index.html?' + errorCode);
}

/**
 * Checks if a string contains a known error msg.
 * @param errorString String to check.
 * @returns If the string is known, the translatable var for showing the error in the UI. If
 * not, null is returned.
 */
function checkIfKnownErrorStrings(errorString: string): string {
  errorString = errorString.toUpperCase();

  let translatableVar: string = null;

  if (errorString.includes('CHANGEADDRESS MUST NOT BE THE NULL ADDRESS')) {
    translatableVar = 'null-change-address-error';
  } else if (errorString.includes('TO IS REQUIRED')) {
    translatableVar = 'to-required-error';
  } else if (errorString.includes('TO.COINS MUST NOT BE ZERO')) {
    translatableVar = 'zero-coins-error';
  } else if (errorString.includes('TO.ADDRESS MUST NOT BE THE NULL ADDRESS')) {
    translatableVar = 'null-destination-error';
  } else if (errorString.includes('TO CONTAINS DUPLICATE VALUES')) {
    translatableVar = 'duplicate-destination-error';
  } else if (errorString.includes('TO.HOURS MUST BE ZERO FOR AUTO TYPE HOURS SELECTION')) {
    translatableVar = 'hours-in-automatic-mode-error';
  } else if (errorString.includes('HOURSSELECTION.MODE IS REQUIRED FOR AUTO TYPE HOURS SELECTION')) {
    translatableVar = 'hours-allocation-mode-needed-error';
  } else if (errorString.includes('INVALID HOURSSELECTION.MODE')) {
    translatableVar = 'invalid-hours-allocation-mode-error';
  } else if (errorString.includes('HOURSSELECTION.MODE CANNOT BE USED FOR MANUAL TYPE HOURS SELECTION')) {
    translatableVar = 'hours-allocation-mode-not-needed-error';
  } else if (errorString.includes('INVALID HOURSSELECTION.TYPE')) {
    translatableVar = 'invalid-hours-mode-error';
  } else if (errorString.includes('HOURSSELECTION.SHAREFACTOR MUST BE SET FOR SHARE MODE')) {
    translatableVar = 'share-factor-needed-error';
  } else if (errorString.includes('HOURSSELECTION.SHAREFACTOR CAN ONLY BE USED FOR SHARE MODE')) {
    translatableVar = 'share-factor-not-needed-error';
  } else if (errorString.includes('HOURSSELECTION.SHAREFACTOR MUST BE >= 0 AND <= 1')) {
    translatableVar = 'invalid-share-factor-error';
  } else if (errorString.includes('TRANSACTION VIOLATES HARD CONSTRAINT: DUPLICATE OUTPUT IN TRANSACTION')) {
    translatableVar = 'change-equal-to-destination-error';
  } else if (errorString.includes('FINGERPRINT CONFLICT FOR')) {
    translatableVar = 'repeated-wallet';
  }

  if (translatableVar) {
    return 'send.known-node-errors.' + translatableVar;
  }

  return null;
}
