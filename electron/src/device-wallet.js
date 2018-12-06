const HID = require('node-hid');
const messages = require('./protob/skycoin');
const bufReceiver = require('./buffer-receiver');
const dgram = require('dgram');
const scanf = require('scanf');
const os = require('os');

let deviceType = 0;

const setDeviceType = function(devType) {
    deviceType = devType;
};

const dataBytesFromChunks = function(chunks) {
    const dataBytes = [];
    chunks.forEach((chunk, j) => {
        chunk.forEach((elt, i) => {
            dataBytes[(64 * j) + i] = elt;
        });
    });
    return dataBytes;
};

// Returns a handle to usbhid device
const getDevice = function() {
    const deviceInfo = HID.devices().find( function(d) {
        const isTeensy = d.manufacturer == "SatoshiLabs";
        return isTeensy;
    });
    if( deviceInfo ) {
        const device = new HID.HID( deviceInfo.path );
        return device;
    }
    return null;
};

// Prepares buffer containing message to device
// eslint-disable-next-line max-statements
const makeTrezorMessage = function(buffer, msgId) {
    const u8Array = new Uint8Array(buffer);
    const trezorMsg = new ArrayBuffer(10 + u8Array.byteLength - 1);
    const dv = new DataView(trezorMsg);
    // Adding the '##' at the begining of the header
    dv.setUint8(0, 35);
    dv.setUint8(1, 35);
    dv.setUint16(2, msgId);
    dv.setUint32(4, u8Array.byteLength);
    // Adding '\n' at the end of the header
    dv.setUint8(8, 10);
    const trezorMsg8 = new Uint8Array(trezorMsg);
    trezorMsg8.set(u8Array.slice(1), 9);
    let lengthToWrite = u8Array.byteLength + 9;
    const chunks = [];
    let j = 0;
    do {
        const u64pack = new Uint8Array(64);
        u64pack[0] = 63;
        u64pack.set(trezorMsg8.slice(63 * j, 63 * (j + 1)), 1);
        lengthToWrite -= 63;
        chunks[j] = u64pack;
        j += 1;
    } while (lengthToWrite > 0);
    return chunks;
};

const emulatorSend = function(client, message) {
    console.log("Sending data", message, message.length);
    const nbChunks = message.length / 64;
    for (let i = 0; i < nbChunks; i += 1) {
        client.send(
            message.slice(64 * i, 64 * (i + 1)), 0, 64, 21324, '127.0.0.1',
            function(err, bytes) {
                if (err) {
                    throw err;
                }
                console.log("Sending data", bytes);
            }
        );
    }
};

DeviceTypeEnum = {
    'EMULATOR': 1,
    'USB': 2
};

class DeviceHandler {
    constructor(devType) {
        this.deviceType = devType;
        this.devHandle = this.getDeviceHandler();
    }

    getDeviceHandler() {
        switch (this.deviceType) {
        case DeviceTypeEnum.USB:
        {
            const dev = getDevice();
            if (dev === null) {
                throw new Error("Device not connected");
            }
            return dev;
        }
        case DeviceTypeEnum.EMULATOR:
        {
            const client = dgram.createSocket('udp4');
            return client;
        }
        default:
            throw new Error("Device type not defined");
        }
    }

    read(devReadCallback) {
        const bufferReceiver = new bufReceiver.BufferReceiver();
        switch (this.deviceType) {
        case DeviceTypeEnum.USB:
            {
                const devHandle = this.devHandle;
                const devHandleCallback = function(err, data) {
                    if (err) {
                        console.error(err);
                        return;
                    }
                    bufferReceiver.receiveBuffer(data, devReadCallback);
                    if (bufferReceiver.bytesToGet > 0) {
                        console.log("Reading one more time", devHandle);
                        devHandle.read(devHandleCallback);
                    }
                };
                devHandle.read(devHandleCallback);
            }
            break;
        case DeviceTypeEnum.EMULATOR:
            this.devHandle.on('message', function(data, rinfo) {
                if (rinfo) {
                    console.log(`server got: 
                        ${data} from ${rinfo.address}:${rinfo.port}`);
                }
                bufferReceiver.receiveBuffer(data, devReadCallback);
            });
            break;
        default:
            throw new Error("Device type not defined");
        }
    }

    // eslint-disable-next-line max-statements
    write(dataBytes) {
        switch (this.deviceType) {
        case DeviceTypeEnum.USB:
        {
            console.log("Writing a buffer of length ", dataBytes.length, "to the device");
            let j = 0;
            let lengthToWrite = dataBytes.length;
            do{
                const u64pack = dataBytes.slice(64 * j, 64 * (j + 1));
                if (os.platform() == 'win32') {
                    u64pack.unshift(0x00);
                }
                this.devHandle.write(u64pack);
                j += 1;
                lengthToWrite -= 64;
            } while (lengthToWrite > 0);
            break;
        }
        case DeviceTypeEnum.EMULATOR:
            emulatorSend(this.devHandle, Buffer.from(dataBytes));
            break;
        default:
            throw new Error("Device type not defined");
        }
    }

    reopen() {
        this.close();
        this.devHandle = this.getDeviceHandler();
    }

    close() {
        switch (this.deviceType) {
        case DeviceTypeEnum.USB:
            this.devHandle.close();
            break;
        case DeviceTypeEnum.EMULATOR:
            this.devHandle.close();
            break;
        default:
            throw new Error("Device type not defined");
        }
    }
}

const createInitializeRequest = function() {
    const msgStructure = {};
    const msg = messages.Initialize.create(msgStructure);
    const buffer = messages.Initialize.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_Initialize
    );
    return dataBytesFromChunks(chunks);
};

const createButtonAckRequest = function() {
    const msgStructure = {};
    const msg = messages.ButtonAck.create(msgStructure);
    const buffer = messages.ButtonAck.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_ButtonAck
    );
    return dataBytesFromChunks(chunks);
};

const createCancelRequest = function() {
    const msgStructure = {};
    const msg = messages.Cancel.create(msgStructure);
    const buffer = messages.Cancel.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_Cancel
    );
    return dataBytesFromChunks(chunks);
};

const createChangePinRequest = function(mnemonic) {
    const msgStructure = {
        mnemonic
    };
    const msg = messages.ChangePin.create(msgStructure);
    const buffer = messages.ChangePin.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_ChangePin
    );
    return dataBytesFromChunks(chunks);
};

const createSetMnemonicRequest = function(mnemonic) {
    const msgStructure = {
        mnemonic
    };
    const msg = messages.SetMnemonic.create(msgStructure);
    const buffer = messages.SetMnemonic.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_SetMnemonic
    );
    return dataBytesFromChunks(chunks);
};

const createGenerateMnemonicRequest = function() {
    const msgStructure = {};
    const msg = messages.GenerateMnemonic.create(msgStructure);
    const buffer = messages.GenerateMnemonic.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_GenerateMnemonic
    );
    return dataBytesFromChunks(chunks);
};

const createGetVersionRequest = function() {
    const msgStructure = {};
    const msg = messages.GetVersion.create(msgStructure);
    const buffer = messages.GetVersion.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_GetVersion
    );
    return dataBytesFromChunks(chunks);
};

const createWipeDeviceRequest = function() {
    const msgStructure = {};
    const msg = messages.WipeDevice.create(msgStructure);
    const buffer = messages.WipeDevice.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_WipeDevice
    );
    return dataBytesFromChunks(chunks);
};

const createBackupDeviceRequest = function() {
    const msgStructure = {};
    const msg = messages.BackupDevice.create(msgStructure);
    const buffer = messages.BackupDevice.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_BackupDevice
    );
    return dataBytesFromChunks(chunks);
};

const createSignMessageRequest = function(addressN, message) {
    const msgStructure = {
        addressN,
        message
    };
    const msg = messages.SkycoinSignMessage.create(msgStructure);
    const buffer = messages.SkycoinSignMessage.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_SkycoinSignMessage
    );
    return dataBytesFromChunks(chunks);
};

const createAddressGenRequest = function(addressN, startIndex) {
    const msgStructure = {
        addressN,
        startIndex
    };
    const msg = messages.SkycoinAddress.create(msgStructure);
    const buffer = messages.SkycoinAddress.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_SkycoinAddress
    );
    return dataBytesFromChunks(chunks);
};

const createCheckMessageSignatureRequest = function(address, message, signature) {
    const msgStructure = {
        address,
        message,
        signature
    };
    const msg = messages.SkycoinCheckMessageSignature.create(msgStructure);
    const buffer = messages.SkycoinCheckMessageSignature.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_SkycoinCheckMessageSignature
    );
    return dataBytesFromChunks(chunks);
};

const createFirmwareUploadRequest = function(payload, hash) {
    const msgStructure = {
        hash,
        payload
    };
    const msg = messages.FirmwareUpload.create(msgStructure);
    const buffer = messages.FirmwareUpload.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_FirmwareUpload
    );
    return dataBytesFromChunks(chunks);
};

const createFirmwareEraseRequest = function(length) {
    const msgStructure = {
        length
    };
    const msg = messages.FirmwareErase.create(msgStructure);
    const buffer = messages.FirmwareErase.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_FirmwareErase
    );
    return dataBytesFromChunks(chunks);
};

const createSendPinCodeRequest = function(pin) {
    const msgStructure = {
        pin
    };
    const msg = messages.PinMatrixAck.create(msgStructure);
    const buffer = messages.PinMatrixAck.encode(msg).finish();
    const chunks = makeTrezorMessage(
        buffer,
        messages.MessageType.MessageType_PinMatrixAck
    );
    return dataBytesFromChunks(chunks);
};

const decodeButtonRequest = function(kind) {
    if (kind != messages.MessageType.MessageType_ButtonRequest) {
        console.error("Skiping button confirmation!", messages.MessageType[kind]);
        return false;
    }
    console.log("ButtonRequest!");
    return true;
};

const decodeSuccess = function(kind, dataBuffer) {
    if (kind == messages.MessageType.MessageType_Success) {
        try {
            const answer = messages.Success.decode(dataBuffer);
            console.log(
                "Success message code",
                answer.code, "message: ",
                answer.message
                );
            return answer.message;
        } catch (e) {
            console.error("Wire format is invalid");
        }
    }
    return `decodeSuccess failed: ${kind}`;
};

const decodeFailureAndPinCode = function(kind, dataBuffer) {
    if (kind == messages.MessageType.MessageType_Failure) {
        try {
            const answer = messages.Failure.decode(dataBuffer);
            console.log(
                "Failure message code",
                answer.code, "message: ",
                answer.message
                );
            return answer.message;
        } catch (e) {
            console.error("Wire format is invalid");
        }
    }
    if (kind == messages.MessageType.MessageType_PinMatrixRequest) {
        return "Pin code required";
    }
    return "decodeFailureAndPinCode failed";
};

const decodeSignMessageAnswer = function(kind, dataBuffer) {
    let signature = "";
    decodeFailureAndPinCode(kind, dataBuffer);
    if (kind == messages.MessageType.
        MessageType_ResponseSkycoinSignMessage) {
        try {
            const answer = messages.ResponseSkycoinSignMessage.
                            decode(dataBuffer);
            signature = answer.signedMessage;
        } catch (e) {
            console.error("Wire format is invalid", e);
        }
    }
    return signature;
};

const decodeAddressGenAnswer = function(kind, dataBuffer) {
    let addresses = [];
    if (kind == messages.MessageType.MessageType_ResponseSkycoinAddress) {
        try {
            const answer = messages.ResponseSkycoinAddress.
                            decode(dataBuffer);
            console.log("Addresses", answer.addresses);
            addresses = answer.addresses;
        } catch (e) {
            console.error("Wire format is invalid", e);
        }
    } else {
        return decodeFailureAndPinCode(kind, dataBuffer);
    }
    return addresses;
};

const devButtonRequestCallback = function(kind, callback) {
    if (decodeButtonRequest(kind)) {
        const dataBytes = createButtonAckRequest();
        const deviceHandle = new DeviceHandler(deviceType);
        const devReadCallback = function(datakind, data) {
            console.log("User hit a button, calling: ", callback);
            deviceHandle.close();
            if (callback !== null && callback !== undefined) {
                // eslint-disable-next-line callback-return
                callback(datakind, data);
            }
        };
        deviceHandle.read(devReadCallback);
        deviceHandle.write(dataBytes);
    }
};

const devUpdateFirmware = function(data, hash) {
    return new Promise((resolve, reject) => {
        const dataBytes = createInitializeRequest();
        const deviceHandle = new DeviceHandler(deviceType);
        const uploadFirmwareCallback = function(kind) {
            deviceHandle.close();
            devButtonRequestCallback(kind, (datakind) => {
                if (datakind == messages.MessageType.MessageType_Success) {
                    resolve("Update firmware operation completed");
                } else {
                    reject(new Error("Update firmware operation failed or refused"));
                }
            });
        };
        const eraseFirmwareCallback = function(eraseStatus, eraseMessage) {
            console.log(decodeSuccess(eraseStatus, eraseMessage));
            deviceHandle.reopen();
            console.log(decodeSuccess(eraseStatus, eraseMessage));
            const uploadDataBytes = createFirmwareUploadRequest(data, hash);
            deviceHandle.read(uploadFirmwareCallback);
            deviceHandle.write(uploadDataBytes);
        };
        const devReadCallback = function(kind, dataBuffer) {
            console.log(decodeSuccess(kind, dataBuffer));
            deviceHandle.reopen();
            const eraseDataBytes = createFirmwareEraseRequest(data.length);
            deviceHandle.read(eraseFirmwareCallback);
            deviceHandle.write(eraseDataBytes);
        };
        deviceHandle.read(devReadCallback);
        deviceHandle.write(dataBytes);
    });
};

const devGetVersionDevice = function() {
    return new Promise((resolve) => {
            const dataBytes = createGetVersionRequest();
            const deviceHandle = new DeviceHandler(deviceType);
            const devReadCallback = function(kind, data) {
                deviceHandle.close();
                const version = decodeSuccess(kind, data);
                if (version == "") {
                    reject(new Error("Could not get version from the device"));
                } else {
                    resolve(version);
                }
            };
            deviceHandle.read(devReadCallback);
            deviceHandle.write(dataBytes);
    });
};

const devAddressGen = function(addressN, startIndex, callback) {
    const dataBytes = createAddressGenRequest(addressN, startIndex);
    const deviceHandle = new DeviceHandler(deviceType);
    const devReadCallback = function(kind, dataBuffer) {
        deviceHandle.close();
        callback(kind, dataBuffer);
    };
    deviceHandle.read(devReadCallback);
    deviceHandle.write(dataBytes);
};

const devSendPinCodeRequest = function(pinCodeCallback, pinCodeReader) {
    const sendPinCodeRequest = function(pinCode) {
        const dataBytes = createSendPinCodeRequest(pinCode);
        const deviceHandle = new DeviceHandler(deviceType);
        deviceHandle.read((answerKind, dataBuffer) => {
            pinCodeCallback(answerKind, dataBuffer);
            deviceHandle.close();
        });
        deviceHandle.write(dataBytes);
    };
    if (pinCodeReader !== null && pinCodeReader !== undefined) {
        const pinCodePromise = pinCodeReader();
        pinCodePromise.then(
            (pinCode) => {
                sendPinCodeRequest(pinCode);
            },
            () => {
                console.log("Pin code promise rejected");
                const dataBytes = createCancelRequest();
                const deviceHandle = new DeviceHandler(deviceType);
                deviceHandle.write(dataBytes);
            }
            );
    } else {
        console.log("Please input your pin code: ");
        sendPinCodeRequest(scanf('%s'));
    }
};

const devAddressGenPinCode = function(addressN, startIndex, pinCodeReader) {
    return new Promise((resolve, reject) => {
        devAddressGen(addressN, startIndex, function(kind, dataBuffer) {
            console.log("Addresses generation kindly returned", messages.MessageType[kind]);
            if (kind == messages.MessageType.MessageType_Failure) {
                reject(new Error(decodeFailureAndPinCode(kind, dataBuffer)));
            }
            if (kind == messages.MessageType.MessageType_ResponseSkycoinAddress) {
                resolve(decodeAddressGenAnswer(kind, dataBuffer));
            }
            if (kind == messages.MessageType.MessageType_PinMatrixRequest) {
                devSendPinCodeRequest(
                    (answerKind, answerBuffer) => {
                    console.log("Pin code callback got answerKind", answerKind);
                    if (answerKind == messages.MessageType.MessageType_ResponseSkycoinAddress) {
                        resolve(decodeAddressGenAnswer(answerKind, answerBuffer));
                    }
                    if (answerKind == messages.MessageType.MessageType_Failure) {
                        reject(new Error(decodeFailureAndPinCode(answerKind, answerBuffer)));
                    }
                },
                pinCodeReader
                );
            }
        });
    });
};

const devSkycoinSignMessage = function(addressN, message, callback) {
    const dataBytes = createSignMessageRequest(addressN, message);
    const deviceHandle = new DeviceHandler(deviceType);
    const devReadCallback = function(kind, dataBuffer) {
        deviceHandle.close();
        callback(kind, dataBuffer);
    };
    deviceHandle.read(devReadCallback);
    deviceHandle.write(dataBytes);
};

const devSkycoinSignMessagePinCode = function(addressN, message, pinCodeReader) {
    return new Promise((resolve, reject) => {
        devSkycoinSignMessage(addressN, message, function(kind, dataBuffer) {
            console.log("Signature generation kindly returned", messages.MessageType[kind]);
            if (kind == messages.MessageType.MessageType_Failure) {
                reject(new Error(decodeFailureAndPinCode(kind, dataBuffer)));
            }
            if (kind == messages.MessageType.MessageType_ResponseSkycoinSignMessage) {
                resolve(decodeSignMessageAnswer(kind, dataBuffer));
            }
            if (kind == messages.MessageType.MessageType_PinMatrixRequest) {
                devSendPinCodeRequest(
                    (answerKind, answerBuffer) => {
                    console.log("Pin code callback got answerKind", answerKind);
                    if (answerKind == messages.MessageType.MessageType_ResponseSkycoinSignMessage) {
                        resolve(decodeSignMessageAnswer(answerKind, answerBuffer));
                    }
                    if (answerKind == messages.MessageType.MessageType_Failure) {
                        reject(new Error(decodeFailureAndPinCode(answerKind, answerBuffer)));
                    }
                },
                pinCodeReader
                );
            }
        });
    });
};

const devCheckMessageSignature = function(address, message, signature) {
    return new Promise((resolve, reject) => {
        const dataBytes = createCheckMessageSignatureRequest(address, message, signature);
        const deviceHandle = new DeviceHandler(deviceType);
        const devReadCallback = function(kind, dataBuffer) {
            if (kind == messages.MessageType.MessageType_Success) {
                try {
                    const answer = messages.Success.
                                    decode(dataBuffer);
                    console.log("Address emiting that signature:", answer.message);
                    if (answer.message === address) {
                        resolve("Signature is correct");
                    } else {
                        reject(new Error("Wrong signature"));
                    }
                } catch (e) {
                    reject(new Error("Wire format is invalid", e));
                }
            } else {
                reject(new Error("Wrong answer kind", kind));
            }
            deviceHandle.close();
        };
        deviceHandle.read(devReadCallback);
        deviceHandle.write(dataBytes);
    });
};

const devWipeDevice = function() {
    return new Promise((resolve) => {
            const dataBytes = createWipeDeviceRequest();
            const deviceHandle = new DeviceHandler(deviceType);
            const devReadCallback = function(kind) {
                deviceHandle.close();
                devButtonRequestCallback(kind, (datakind) => {
                    if (datakind == messages.MessageType.MessageType_Success) {
                        resolve("Wipe Device operation completed");
                    } else {
                        resolve("Wipe Device operation failed or refused");
                    }
                });
            };
            deviceHandle.read(devReadCallback);
            deviceHandle.write(dataBytes);
    });
};

const devBackupDevice = function() {
    return new Promise((resolve) => {
            const dataBytes = createBackupDeviceRequest();
            const deviceHandle = new DeviceHandler(deviceType);
            const buttonAckLoop = function(kind) {
                if (kind != messages.MessageType.MessageType_ButtonRequest) {
                    deviceHandle.close();
                    if (kind == messages.MessageType.MessageType_Success) {
                        resolve("Backup Device operation completed");
                    } else {
                        resolve("Backup Device operation failed or refused");
                    }
                    return;
                }
                const buttonAckBytes = createButtonAckRequest();
                deviceHandle.read(buttonAckLoop);
                deviceHandle.write(buttonAckBytes);
            };
            deviceHandle.read(buttonAckLoop);
            deviceHandle.write(dataBytes);
    });
};

const devSetMnemonic = function(mnemonic) {
    return new Promise((resolve) => {
        const dataBytes = createSetMnemonicRequest(mnemonic);
        const deviceHandle = new DeviceHandler(deviceType);
        const devReadCallback = function(kind) {
            deviceHandle.close();
            devButtonRequestCallback(kind, (datakind) => {
                if (datakind == messages.MessageType.MessageType_Success) {
                    resolve("Set Mnemonic operation completed");
                } else {
                    resolve("Set Mnemonic operation failed or refused");
                }
            });
        };
        deviceHandle.read(devReadCallback);
        deviceHandle.write(dataBytes);
    });
};

const devGenerateMnemonic = function() {
    return new Promise((resolve) => {
        const dataBytes = createGenerateMnemonicRequest();
        const deviceHandle = new DeviceHandler(deviceType);
        const devReadCallback = function(kind) {
            deviceHandle.close();
            devButtonRequestCallback(kind, (datakind) => {
                if (datakind == messages.MessageType.MessageType_Success) {
                    resolve("Generate Mnemonic operation completed");
                } else {
                    resolve("Generate Mnemonic operation failed or refused");
                }
            });
        };
        deviceHandle.read(devReadCallback);
        deviceHandle.write(dataBytes);
    });
};

const devChangePin = function(pinCodeReader) {
    return new Promise((resolve, reject) => {
        const dataBytes = createChangePinRequest();
        const deviceHandle = new DeviceHandler(deviceType);
        const pinCodeMatrixCallback = function(datakind, dataBuffer) {
            console.log("pinCodeMatrixCallback kind:", datakind, messages.MessageType[datakind]);
            if (datakind == messages.MessageType.MessageType_PinMatrixRequest) {
                devSendPinCodeRequest(pinCodeMatrixCallback, pinCodeReader);
            }
            if (datakind == messages.MessageType.MessageType_Failure) {
                reject(new Error(decodeFailureAndPinCode(datakind, dataBuffer)));
            }
            if (datakind == messages.MessageType.MessageType_Success) {
                resolve(decodeSuccess(datakind, dataBuffer));
            }
        };
        const devReadCallback = function(kind) {
            deviceHandle.close();
            devButtonRequestCallback(kind, pinCodeMatrixCallback);
        };
        deviceHandle.read(devReadCallback);
        deviceHandle.write(dataBytes);
    });
};

module.exports = {
    DeviceTypeEnum,
    devAddressGen,
    devAddressGenPinCode,
    devBackupDevice,
    devChangePin,
    devCheckMessageSignature,
    devGenerateMnemonic,
    devGetVersionDevice,
    devSetMnemonic,
    devSkycoinSignMessagePinCode,
    devUpdateFirmware,
    devWipeDevice,
    getDevice,
    makeTrezorMessage,
    setDeviceType
};
