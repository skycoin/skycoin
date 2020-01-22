export enum OperationErrorCategories {
  GeneralApiError = 'GeneralApiError',
  HwApiError = 'HwApiError',
}

export enum OperationErrorTypes {
  NoInternet = 'NoInternet',
  Unauthorized = 'Unauthorized',
  Unknown = 'Unknown',
}

export enum HWOperationResults {
  Success = 'Success',
  FailedOrRefused = 'FailedOrRefused',
  PinMismatch = 'PinMismatch',
  WithoutSeed = 'WithoutSeed',
  WrongPin = 'WrongPin',
  IncorrectHardwareWallet = 'IncorrectHardwareWallet',
  WrongWord = 'WrongWord',
  InvalidSeed = 'InvalidSeed',
  WrongSeed = 'WrongSeed',
  UndefinedError = 'UndefinedError',
  Disconnected = 'Disconnected',
  DaemonError = 'DaemonError',
  InvalidAddress = 'InvalidAddress',
  Timeout = 'Timeout',
  NotInBootloaderMode = 'NotInBootloaderMode',
  AddressGeneratorProblem = 'AddressGeneratorProblem',
}

export class OperationError {
  category: OperationErrorCategories;
  type: OperationErrorTypes | HWOperationResults;
  originalError: any;
  originalServerErrorMsg: string;
  translatableErrorMsg: string;
}
