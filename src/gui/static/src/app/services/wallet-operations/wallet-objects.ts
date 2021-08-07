import BigNumber from 'bignumber.js';

import { Output } from './transaction-objects';

/**
 * This file contains the objects used to represent the wallets and addresses in the app.
 */

// Base wallets
////////////////////////////////////////////////

/**
 * Basic wallet object with the most important properties.
 */
export class WalletBase {
  // NOTE: All properties must have an initial value or there could be problems creating duplicates.

  /**
   * Name used to identify the wallet.
   */
  label = '';
  /**
   * Unique ID of the wallet. In software wallets it is the name of the file and in hw
   * wallets it is the first address.
   */
  id = '';
  /**
   * Address list.
   */
  addresses: AddressBase[] = [];
  /**
   * If the wallet is encrypted with a password. Only valid for software wallets.
   */
  encrypted = false;
  /**
   * If the wallet is temporal.
   */
  temporal = false;
  /**
   * If it is a software wallet (false) or a hw wallet (true).
   */
  isHardware = false;
  /**
   * If the last time the wallet was checked there were security warning found. Only valid for
   * hw wallets.
   */
  hasHwSecurityWarnings = false;
  /**
   * If the user asked the app to stop blocking access to some functions by showing a security
   * popup when hasHwSecurityWarnings is true. Only valid for hw wallets.
   */
  stopShowingHwSecurityPopup = false;
}

/**
 * Basic address object with the most important properties.
 */
export class AddressBase {
  // NOTE: All properties must have an initial value or there could be problems creating duplicates.

  /**
   * Address string.
   */
  address = '';
  /**
   * If the address has been confirmed by the user on the hw wallet and can be shown on the UI.
   * Only valid if the address is in a hw wallet.
   */
  confirmed = false;
}

/**
 * Creates a duplicate of a WalletBase object. If the provided wallet has properties which are not
 * part of WalletBase, those properties are removed.
 * @param wallet Object to duplicate.
 * @param duplicateAddresses If the addresses must be duplicated as instancies of AddressBase
 * (true) or if the address arrays must be returned empty (false).
 */
export function duplicateWalletBase(wallet: WalletBase, duplicateAddresses: boolean): WalletBase {
  const response = new WalletBase();
  Object.assign(response, wallet);
  removeAdditionalProperties(true, response);

  response.addresses = [];
  if (duplicateAddresses) {
    wallet.addresses.forEach(address => {
      response.addresses.push(duplicateAddressBase(address));
    });
  }

  return response;
}

/**
 * Creates a duplicate of a AddressBase object. If the provided address has properties which
 * are not part of AddressBase, those properties are removed.
 * @param address Object to duplicate.
 */
function duplicateAddressBase(address: AddressBase): AddressBase {
  const response = new AddressBase();
  Object.assign(response, address);
  removeAdditionalProperties(false, response);

  return response;
}

/**
 * Removes from an object all the properties which are not part of WalletBase or AddressBase.
 * @param useWalletBaseAsReference If true, only the properties of WalletBase will be keep; if
 * false, only the properties of AddressBase will be keep.
 * @param objectToClean Object to be cleaned.
 */
function removeAdditionalProperties(useWalletBaseAsReference: boolean, objectToClean: any) {
  const knownPropertiesMap = new Map<string, boolean>();
  const reference: Object = useWalletBaseAsReference ? new WalletBase() : new AddressBase();
  Object.keys(reference).forEach(property => {
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

/**
 * Object with the basic data of a wallet and data about its balance.
 */
export class WalletWithBalance extends WalletBase {
  coins = new BigNumber(0);
  hours = new BigNumber(0);
  addresses: AddressWithBalance[] = [];
}

/**
 * Object with the basic data of an address and data about its balance.
 */
export class AddressWithBalance extends AddressBase {
  coins = new BigNumber(0);
  hours = new BigNumber(0);
}

/**
 * Creates a new WalletWithBalance instance with copies of the values of
 * a WalletBase object.
 */
export function walletWithBalanceFromBase(wallet: WalletBase): WalletWithBalance {
  const response = new WalletWithBalance();
  Object.assign(response, duplicateWalletBase(wallet, false));

  wallet.addresses.forEach(address => {
    response.addresses.push(addressWithBalanceFromBase(address));
  });

  return response;
}

/**
 * Creates a new AddressWithBalance instance with copies of the values of
 * an AddressBase object.
 */
function addressWithBalanceFromBase(address: AddressBase): AddressWithBalance {
  const response = new AddressWithBalance();
  Object.assign(response, duplicateAddressBase(address));

  return response;
}

// Wallets with outputs
////////////////////////////////////////////////

/**
 * Object with the basic data of a wallet and data about its unspent outputs.
 */
export class WalletWithOutputs extends WalletBase {
  addresses: AddressWithOutputs[] = [];
}

/**
 * Object with the basic data of an address and data about its unspent outputs.
 */
export class AddressWithOutputs extends AddressBase {
  outputs: Output[] = [];
}

/**
 * Creates a new WalletWithOutputs instance with copies of the values of
 * a WalletBase object.
 */
export function walletWithOutputsFromBase(wallet: WalletBase): WalletWithOutputs {
  const response = new WalletWithOutputs();
  Object.assign(response, duplicateWalletBase(wallet, false));

  wallet.addresses.forEach(address => {
    response.addresses.push(addressWithOutputsFromBase(address));
  });

  return response;
}

/**
 * Creates a new AddressWithOutputs instance with copies of the values of
 * an AddressBase object.
 */
function addressWithOutputsFromBase(address: AddressBase): AddressWithOutputs {
  const response = new AddressWithOutputs();
  Object.assign(response, duplicateAddressBase(address));

  return response;
}
