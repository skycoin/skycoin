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
  DaemonConnectionError = 'DaemonConnectionError',
  InvalidAddress = 'InvalidAddress',
  Timeout = 'Timeout',
  NotInBootloaderMode = 'NotInBootloaderMode',
  AddressGeneratorProblem = 'AddressGeneratorProblem',
}

export class OperationError {
  type!: HWOperationResults;
  originalError: any;
  originalServerErrorMsg!: string | null;
  translatableErrorMsg!: string;
}
