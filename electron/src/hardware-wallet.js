'use strict'

const { ipcMain } = require('electron');

const deviceWallet = require('hardware-wallet-js/device-wallet');

const { Observable, of } = require('rxjs');

const HID = require('node-hid');

// Global reference of the window object.
let win;

function setWinRef(winRef) {
  win = winRef;
}

// Detect if the code is running with the "hw" arg. The "hw" arg is added when running npm
// start. If this is true, the UI will show the hardware wallet options.
let hw = process.argv.find(arg => arg === 'hw') ? true : false;

deviceWallet.setDeviceType(deviceWallet.DeviceTypeEnum.USB);
ipcMain.on('hwCompatibilityActivated', (event) => {
  event.returnValue = hw;
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
    lastPinPromiseReject(new Error("Cancelled"))
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
    lastSeedWordPromiseReject(new Error("Cancelled"))
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

ipcMain.on('hwGetAddresses', (event, requestId, addressN, startIndex) => {
  const promise = deviceWallet.devAddressGen(addressN, startIndex, false, pinEvent);
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

ipcMain.on('hwGenerateMnemonic', (event, requestId) => {
  const promise = deviceWallet.devGenerateMnemonic(false);
  promise.then(
    result => { console.log("Generate mnemonic promise resolved", result); event.sender.send('hwGenerateMnemonicResponse', requestId, result); },
    error => { console.log("Generate mnemonic promise errored: ", error); event.sender.send('hwGenerateMnemonicResponse', requestId, { error: error.toString() }); }
  );
});

ipcMain.on('hwRecoverMnemonic', (event, requestId) => {
  const promise = deviceWallet.devRecoveryDevice(false, requestSeedWordEvent);
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
  const promise = deviceWallet.devSkycoinSignMessage(addressIndex, message, pinEvent);
  promise.then(
    result => { console.log("Signature promise resolved", result); event.sender.send('hwSignMessageResponse', requestId, result); },
    error => { console.log("Signature promise errored: ", error); event.sender.send('hwSignMessageResponse', requestId, { error: error.toString() }); }
  );
});

module.exports = {
  setWinRef
}