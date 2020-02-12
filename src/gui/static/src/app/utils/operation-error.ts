/**
 * This file contains the basic class and values for easily working with errors on the app.
 */

/**
 * Possible values of OperationError.type for identifying errors during general operations.
 */
export enum OperationErrorTypes {
  /**
   * There is no internet connection.
   */
  NoInternet = 'NoInternet',
  /**
   * The user is not authorized. Normally means that the password is incorrect.
   */
  Unauthorized = 'Unauthorized',
  /**
   * The error is not in the list of known errors that require special treatment. This does not
   * mean the error is rare or specially bad. Just showing the error msg should be enough.
   */
  Unknown = 'Unknown',
}

/**
 * Possible results of a hw wallet operation.
 */
export enum HWOperationResults {
  /**
   * The operation was completed successfuly.
   */
  Success = 'Success',
  /**
   * The user canceled the operation. Due to how the hw wallet daemon works, it may also mean
   * that there was an unexpected error, but only on rare cases.
   */
  FailedOrRefused = 'FailedOrRefused',
  /**
   * The user entered a different PIN when asked to repeat the new PIN to confirm it.
   */
  PinMismatch = 'PinMismatch',
  /**
   * The device does not have a seed.
   */
  WithoutSeed = 'WithoutSeed',
  /**
   * The user entered a wrong PIN.
   */
  WrongPin = 'WrongPin',
  /**
   * The currently connected device is not the one needed for completing the operation.
   */
  IncorrectHardwareWallet = 'IncorrectHardwareWallet',
  /**
   * The user entered a word which does not match the one requested by the device.
   */
  WrongWord = 'WrongWord',
  /**
   * The seed entered by the user is not a valid BIP39 seed.
   */
  InvalidSeed = 'InvalidSeed',
  /**
   * The seed entered by the user is different from the one on the device.
   */
  WrongSeed = 'WrongSeed',
  /**
   * Unknown or unexpected error.
   */
  UndefinedError = 'UndefinedError',
  /**
   * There is no device connected.
   */
  Disconnected = 'Disconnected',
  /**
   * It was not possible to connect with the hw daemon.
   */
  DaemonConnectionError = 'DaemonConnectionError',
  /**
   * An invalid address was sent to the device.
   */
  InvalidAddress = 'InvalidAddress',
  /**
   * The operation was automatically cancelled due to inactivity.
   */
  Timeout = 'Timeout',
  /**
   * The device must be in bootloader mode for completing the operation and it is not.
   */
  NotInBootloaderMode = 'NotInBootloaderMode',
  /**
   * The device generated an invalid address.
   */
  AddressGeneratorProblem = 'AddressGeneratorProblem',
}

/**
 * Base object for working with errors throughout the application.
 */
export class OperationError {
  /**
   * Specific error type. Allows to know the cause of the error.
   */
  type: OperationErrorTypes | HWOperationResults;
  /**
   * Original error object from which this OperationError instance was created.
   */
  originalError: any;
  /**
   * Original, unprocessed, error msg.
   */
  originalServerErrorMsg: string;
  /**
   * Processed error msg, which can be passed to the 'translate' pipe to display it on the UI.
   */
  translatableErrorMsg: string;
}
