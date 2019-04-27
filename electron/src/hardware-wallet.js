'use strict'

const { ipcMain } = require('electron');

const deviceWallet = require('hardware-wallet-js/device-wallet');

const { Observable, of } = require('rxjs');

const HID = require('node-hid');

const fs = require('fs');

const path = require('path');

// Global reference of the window object.
let win;

function setWinRef(winRef) {
  win = winRef;
}

let fullWalletsFilePath;
let walletsFilePath;
let getSavedWalletsDataSyncEvent;

function setWalletsFolderPath(folderPath) {
  fullWalletsFilePath = path.join(folderPath, 'hw.data');
  walletsFilePath = folderPath;

  if (getSavedWalletsDataSyncEvent) {
    getSavedWalletsData(getSavedWalletsDataSyncEvent);
    getSavedWalletsDataSyncEvent = undefined;
  }
}

deviceWallet.setDeviceType(deviceWallet.DeviceTypeEnum.USB);
ipcMain.on('hwCompatibilityActivated', (event) => {
  event.returnValue = true;
});

let checkHwSubscription;
let hwConnected = false;

function checkHw(wait) {
  if (checkHwSubscription) {
    checkHwSubscription.unsubscribe();
  }

  checkHwSubscription = Observable.of(1)
    .delay(wait ? (hwConnected ? 2000 : 10000) : 0)
    .subscribe(
      () => {
        const connected = HID.devices().find((d) => d.manufacturer === "SkycoinFoundation");
        if (connected && !hwConnected) {
          hwConnected = true;
          if (win) {
            win.webContents.send('hwConnectionEvent', true);
          }
        } else if (!connected && hwConnected) {
          hwConnected = false;
          if (win) {
            win.webContents.send('hwConnectionEvent', false);
          }
        }
        checkHw(true);
      }
    );
}

checkHw(false);

ipcMain.on('hwGetDeviceConnectedSync', (event) => {
  event.returnValue = HID.devices().find((d) => d.manufacturer && d.manufacturer === "SkycoinFoundation") !== undefined;
  checkHw(false);
});

ipcMain.on('hwGetSavedWalletsDataSync', (event) => {
  if (fullWalletsFilePath) {
    getSavedWalletsData(event);
  } else {
    getSavedWalletsDataSyncEvent = event;
  }
});

function getSavedWalletsData(event) {
  if (fs.existsSync(fullWalletsFilePath)) {
    event.returnValue = fs.readFileSync(fullWalletsFilePath, 'utf8');
  } else {
    event.returnValue = '';
  }
}

ipcMain.on('hwSaveWalletsDataSync', (event, data) => {
  if (!fs.existsSync(walletsFilePath)) {
    fs.mkdirSync(walletsFilePath, { recursive: true });
  }
  fs.writeFileSync(fullWalletsFilePath, data, 'utf8');
  event.returnValue = null;
});

let lastPinPromiseResolve;
let lastPinPromiseReject;

const pinEvent = function() {
  return new Promise((resolve, reject) => {
    lastPinPromiseResolve = resolve;
    lastPinPromiseReject = reject;

    console.log("Hardware wallet pin requested");
    if (win) {
      win.webContents.send('hwPinRequested');
    }
  });
};

ipcMain.on('hwSendPin', (event, pin) => {
  if (pin) {
    lastPinPromiseResolve(pin);
  } else {
    lastPinPromiseReject(new Error("Cancelled"));
  }
});

let lastPassphrasePromiseResolve;
let lastPassphrasePromiseReject;

const passphraseEvent = function() {
  return new Promise((resolve, reject) => {
    lastPassphrasePromiseResolve = resolve;
    lastPassphrasePromiseReject = reject;

    console.log("Hardware wallet passphrase requested");
    if (win) {
      win.webContents.send('hwPassphraseRequested');
    }
  });
};

ipcMain.on('hwSendPassphrase', (event, passphrase) => {
  if (passphrase) {
    lastPassphrasePromiseResolve(passphrase);
  } else {
    lastPassphrasePromiseReject(new Error("Cancelled"));
  }
});

let lastSeedWordPromiseResolve;
let lastSeedWordPromiseReject;

const requestSeedWordEvent = function() {
  return new Promise((resolve, reject) => {
    lastSeedWordPromiseResolve = resolve;
    lastSeedWordPromiseReject = reject;

    console.log("Hardware wallet seed word requested");
    if (win) {
      win.webContents.send('hwSeedWordRequested');
    }
  });
};

ipcMain.on('hwSendSeedWord', (event, word) => {
  if (word) {
    lastSeedWordPromiseResolve(word);
  } else {
    lastSeedWordPromiseReject(new Error("Cancelled"));
  }
});

ipcMain.on('hwCancelLastAction', (event, requestId) => {
  const promise = deviceWallet.devCancelRequest();
  promise.then(
    result => { console.log("Cancel promise resolved", result); event.sender.send('hwCancelLastActionResponse', requestId, ''); },
    error => { console.log("Cancel promise errored: ", error); event.sender.send('hwCancelLastActionResponse', requestId, ''); }
  );
});

ipcMain.on('hwGetFeatures', (event, requestId) => {
  const promise = deviceWallet.devGetFeatures();
  promise.then(
    result => { console.log("Features promise resolved", result); event.sender.send('hwGetFeaturesResponse', requestId, result); },
    error => { console.log("Features promise errored: ", error); event.sender.send('hwGetFeaturesResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwGetAddresses', (event, requestId, addressN, startIndex, confirm) => {
  const promise = deviceWallet.devAddressGen(addressN, startIndex, confirm, pinEvent, passphraseEvent);
  promise.then(
    addresses => { console.log("Addresses promise resolved", addresses); event.sender.send('hwGetAddressesResponse', requestId, addresses); },
    error => { console.log("Addresses promise errored: ", error); event.sender.send('hwGetAddressesResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwChangePin', (event, requestId) => {
  const promise = deviceWallet.devChangePin(pinEvent);
  promise.then(
    result => { console.log("Change pin promise resolved", result); event.sender.send('hwChangePinResponse', requestId, result); },
    error => { console.log("Change pin promise errored: ", error); event.sender.send('hwChangePinResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwGenerateMnemonic', (event, requestId, wordCount) => {
  const promise = deviceWallet.devGenerateMnemonic(wordCount, false);
  promise.then(
    result => { console.log("Generate mnemonic promise resolved", result); event.sender.send('hwGenerateMnemonicResponse', requestId, result); },
    error => { console.log("Generate mnemonic promise errored: ", error); event.sender.send('hwGenerateMnemonicResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwRecoverMnemonic', (event, requestId, wordCount, dryRun) => {
  const promise = deviceWallet.devRecoveryDevice(wordCount, false, requestSeedWordEvent, dryRun);
  promise.then(
    result => { console.log("Recover mnemonic promise resolved", result); event.sender.send('hwRecoverMnemonicResponse', requestId, result); },
    error => { console.log("Recover mnemonic promise errored: ", error); event.sender.send('hwRecoverMnemonicResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwBackupDevice', (event, requestId) => {
  const promise = deviceWallet.devBackupDevice(pinEvent);
  promise.then(
    result => { console.log("Backup device promise resolved", result); event.sender.send('hwBackupDeviceResponse', requestId, result); },
    error => { console.log("Backup device promise errored: ", error); event.sender.send('hwBackupDeviceResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwWipe', (event, requestId) => {
  const promise = deviceWallet.devWipeDevice();
  promise.then(
    result => { console.log("Wipe promise resolved", result); event.sender.send('hwWipeResponse', requestId, result); },
    error => { console.log("Wipe promise errored: ", error); event.sender.send('hwWipeResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwSignMessage', (event, requestId, addressIndex, message) => {
  const promise = deviceWallet.devSkycoinSignMessage(addressIndex, message, pinEvent, passphraseEvent);
  promise.then(
    result => { console.log("Signature promise resolved", result); event.sender.send('hwSignMessageResponse', requestId, result); },
    error => { console.log("Signature promise errored: ", error); event.sender.send('hwSignMessageResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwSignTransaction', (event, requestId, inputs, outputs) => {
  const promise = deviceWallet.devSkycoinTransactionSign(inputs, outputs, pinEvent, passphraseEvent);
  promise.then(
    result => { console.log("Sign transaction promise resolved", result); event.sender.send('hwSignTransactionResponse', requestId, result); },
    error => { console.log("Sign transaction promise errored: ", error); event.sender.send('hwSignTransactionResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwApplySettings', (event, requestId, usePassphrase) => {
  const promise = deviceWallet.devApplySettings(usePassphrase, pinEvent);
  promise.then(
    result => { console.log("Apply settings promise resolved", result); event.sender.send('hwApplySettingsResponse', requestId, result); },
    error => { console.log("Apply settings promise errored: ", error); event.sender.send('hwApplySettingsResponse', requestId, { error: error.toString() }); }
  );
});

module.exports = {
  setWinRef,
  setWalletsFolderPath
}
