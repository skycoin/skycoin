import BigNumber from 'bignumber.js';

// Base wallets
////////////////////////////////////////////////

export class WalletBase {
  label: string;
  id: string;
  addresses: AddressBase[];
  encrypted: boolean;
  isHardware: boolean;
  hasHwSecurityWarnings: boolean;
  stopShowingHwSecurityPopup: boolean;
}

export class AddressBase {
  address: string;
  confirmed: boolean; // Optional parameter for hardware wallets only
}

export function duplicateWalletBase(wallet: WalletBase, duplicateAddresses: boolean): WalletBase {
  const response = new WalletBase();
  Object.assign(response, wallet);
  removeAdditionalProperties(WalletBase, response);

  response.addresses = [];
  if (duplicateAddresses) {
    wallet.addresses.forEach(address => {
      response.addresses.push(duplicateAddressBase(address));
    });
  }

  return response;
}

function duplicateAddressBase(address: AddressBase): AddressBase {
  const response = new AddressBase();
  Object.assign(response, address);
  removeAdditionalProperties(AddressBase, response);

  return response;
}

function removeAdditionalProperties(baseClass: any, objectToClean: any) {
  const knownPropertiesMap = new Map<string, boolean>();
  Object.getOwnPropertyNames(baseClass).forEach(property => {
    knownPropertiesMap.set(property, true);
  });

  const propertiesToRemove: string[] = [];
  Object.keys(objectToClean).forEach(property => {
    if (!knownPropertiesMap.has(property)) {
      propertiesToRemove.push(property);
    }
  });

  propertiesToRemove.forEach(property => {
    delete objectToClean[property];
  });
}

// Wallets with balance
////////////////////////////////////////////////

export class WalletWithBalance extends WalletBase {
  coins = new BigNumber(0);
  hours = new BigNumber(0);
  addresses: AddressWithBalance[];
}

export class AddressWithBalance extends AddressBase {
  coins = new BigNumber(0);
  hours = new BigNumber(0);
}

export function walletWithBalanceFromBase(wallet: WalletBase): WalletWithBalance {
  const response = new WalletWithBalance();
  Object.assign(response, duplicateWalletBase(wallet, false));

  wallet.addresses.forEach(address => {
    response.addresses.push(addressWithBalanceFromBase(address));
  });

  return response;
}

function addressWithBalanceFromBase(address: AddressBase): AddressWithBalance {
  const response = new AddressWithBalance();
  Object.assign(response, duplicateAddressBase(address));

  return response;
}

// Wallets with outputs
////////////////////////////////////////////////

export class Output {
  address: string;
  coins: BigNumber;
  hash: string;
  calculated_hours: BigNumber;
}

export class WalletWithOutputs extends WalletBase {
  addresses: AddressWithOutputs[];
}

export class AddressWithOutputs extends AddressBase {
  outputs: Output[] = [];
}

export function walletWithOutputsFromBase(wallet: WalletBase): WalletWithOutputs {
  const response = new WalletWithOutputs();
  Object.assign(response, duplicateWalletBase(wallet, false));

  wallet.addresses.forEach(address => {
    response.addresses.push(addressWithOutputsFromBase(address));
  });

  return response;
}

function addressWithOutputsFromBase(address: AddressBase): AddressWithOutputs {
  const response = new AddressWithOutputs();

  return response;
}
