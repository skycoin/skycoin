/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
"use strict";

var $protobuf = require("protobufjs/minimal");

// Common aliases
var $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
var $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

/**
 * Mapping between Trezor wire identifier (uint) and a protobuf message
 * @exports MessageType
 * @enum {string}
 * @property {number} MessageType_Initialize=0 MessageType_Initialize value
 * @property {number} MessageType_Ping=1 MessageType_Ping value
 * @property {number} MessageType_Success=2 MessageType_Success value
 * @property {number} MessageType_Failure=3 MessageType_Failure value
 * @property {number} MessageType_ChangePin=4 MessageType_ChangePin value
 * @property {number} MessageType_WipeDevice=5 MessageType_WipeDevice value
 * @property {number} MessageType_FirmwareErase=6 MessageType_FirmwareErase value
 * @property {number} MessageType_FirmwareUpload=7 MessageType_FirmwareUpload value
 * @property {number} MessageType_GetEntropy=9 MessageType_GetEntropy value
 * @property {number} MessageType_Entropy=10 MessageType_Entropy value
 * @property {number} MessageType_LoadDevice=13 MessageType_LoadDevice value
 * @property {number} MessageType_ResetDevice=14 MessageType_ResetDevice value
 * @property {number} MessageType_Features=17 MessageType_Features value
 * @property {number} MessageType_PinMatrixRequest=18 MessageType_PinMatrixRequest value
 * @property {number} MessageType_PinMatrixAck=19 MessageType_PinMatrixAck value
 * @property {number} MessageType_Cancel=20 MessageType_Cancel value
 * @property {number} MessageType_ButtonRequest=26 MessageType_ButtonRequest value
 * @property {number} MessageType_ButtonAck=27 MessageType_ButtonAck value
 * @property {number} MessageType_BackupDevice=34 MessageType_BackupDevice value
 * @property {number} MessageType_EntropyRequest=35 MessageType_EntropyRequest value
 * @property {number} MessageType_EntropyAck=36 MessageType_EntropyAck value
 * @property {number} MessageType_PassphraseRequest=41 MessageType_PassphraseRequest value
 * @property {number} MessageType_PassphraseAck=42 MessageType_PassphraseAck value
 * @property {number} MessageType_PassphraseStateRequest=77 MessageType_PassphraseStateRequest value
 * @property {number} MessageType_PassphraseStateAck=78 MessageType_PassphraseStateAck value
 * @property {number} MessageType_RecoveryDevice=45 MessageType_RecoveryDevice value
 * @property {number} MessageType_WordRequest=46 MessageType_WordRequest value
 * @property {number} MessageType_WordAck=47 MessageType_WordAck value
 * @property {number} MessageType_SetMnemonic=113 MessageType_SetMnemonic value
 * @property {number} MessageType_SkycoinAddress=114 MessageType_SkycoinAddress value
 * @property {number} MessageType_SkycoinCheckMessageSignature=115 MessageType_SkycoinCheckMessageSignature value
 * @property {number} MessageType_SkycoinSignMessage=116 MessageType_SkycoinSignMessage value
 * @property {number} MessageType_ResponseSkycoinAddress=117 MessageType_ResponseSkycoinAddress value
 * @property {number} MessageType_ResponseSkycoinSignMessage=118 MessageType_ResponseSkycoinSignMessage value
 * @property {number} MessageType_GenerateMnemonic=119 MessageType_GenerateMnemonic value
 */
$root.MessageType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "MessageType_Initialize"] = 0;
    values[valuesById[1] = "MessageType_Ping"] = 1;
    values[valuesById[2] = "MessageType_Success"] = 2;
    values[valuesById[3] = "MessageType_Failure"] = 3;
    values[valuesById[4] = "MessageType_ChangePin"] = 4;
    values[valuesById[5] = "MessageType_WipeDevice"] = 5;
    values[valuesById[6] = "MessageType_FirmwareErase"] = 6;
    values[valuesById[7] = "MessageType_FirmwareUpload"] = 7;
    values[valuesById[9] = "MessageType_GetEntropy"] = 9;
    values[valuesById[10] = "MessageType_Entropy"] = 10;
    values[valuesById[13] = "MessageType_LoadDevice"] = 13;
    values[valuesById[14] = "MessageType_ResetDevice"] = 14;
    values[valuesById[17] = "MessageType_Features"] = 17;
    values[valuesById[18] = "MessageType_PinMatrixRequest"] = 18;
    values[valuesById[19] = "MessageType_PinMatrixAck"] = 19;
    values[valuesById[20] = "MessageType_Cancel"] = 20;
    values[valuesById[26] = "MessageType_ButtonRequest"] = 26;
    values[valuesById[27] = "MessageType_ButtonAck"] = 27;
    values[valuesById[34] = "MessageType_BackupDevice"] = 34;
    values[valuesById[35] = "MessageType_EntropyRequest"] = 35;
    values[valuesById[36] = "MessageType_EntropyAck"] = 36;
    values[valuesById[41] = "MessageType_PassphraseRequest"] = 41;
    values[valuesById[42] = "MessageType_PassphraseAck"] = 42;
    values[valuesById[77] = "MessageType_PassphraseStateRequest"] = 77;
    values[valuesById[78] = "MessageType_PassphraseStateAck"] = 78;
    values[valuesById[45] = "MessageType_RecoveryDevice"] = 45;
    values[valuesById[46] = "MessageType_WordRequest"] = 46;
    values[valuesById[47] = "MessageType_WordAck"] = 47;
    values[valuesById[113] = "MessageType_SetMnemonic"] = 113;
    values[valuesById[114] = "MessageType_SkycoinAddress"] = 114;
    values[valuesById[115] = "MessageType_SkycoinCheckMessageSignature"] = 115;
    values[valuesById[116] = "MessageType_SkycoinSignMessage"] = 116;
    values[valuesById[117] = "MessageType_ResponseSkycoinAddress"] = 117;
    values[valuesById[118] = "MessageType_ResponseSkycoinSignMessage"] = 118;
    values[valuesById[119] = "MessageType_GenerateMnemonic"] = 119;
    return values;
})();

$root.Initialize = (function() {

    /**
     * Properties of an Initialize.
     * @exports IInitialize
     * @interface IInitialize
     * @property {Uint8Array|null} [state] Initialize state
     */

    /**
     * Constructs a new Initialize.
     * @exports Initialize
     * @classdesc Request: Reset device to default state and ask for device details
     * @next Features
     * @implements IInitialize
     * @constructor
     * @param {IInitialize=} [properties] Properties to set
     */
    function Initialize(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Initialize state.
     * @member {Uint8Array} state
     * @memberof Initialize
     * @instance
     */
    Initialize.prototype.state = $util.newBuffer([]);

    /**
     * Creates a new Initialize instance using the specified properties.
     * @function create
     * @memberof Initialize
     * @static
     * @param {IInitialize=} [properties] Properties to set
     * @returns {Initialize} Initialize instance
     */
    Initialize.create = function create(properties) {
        return new Initialize(properties);
    };

    /**
     * Encodes the specified Initialize message. Does not implicitly {@link Initialize.verify|verify} messages.
     * @function encode
     * @memberof Initialize
     * @static
     * @param {IInitialize} message Initialize message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Initialize.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.state != null && message.hasOwnProperty("state"))
            writer.uint32(/* id 1, wireType 2 =*/10).bytes(message.state);
        return writer;
    };

    /**
     * Encodes the specified Initialize message, length delimited. Does not implicitly {@link Initialize.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Initialize
     * @static
     * @param {IInitialize} message Initialize message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Initialize.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Initialize message from the specified reader or buffer.
     * @function decode
     * @memberof Initialize
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Initialize} Initialize
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Initialize.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Initialize();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.state = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an Initialize message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Initialize
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Initialize} Initialize
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Initialize.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Initialize message.
     * @function verify
     * @memberof Initialize
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Initialize.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.state != null && message.hasOwnProperty("state"))
            if (!(message.state && typeof message.state.length === "number" || $util.isString(message.state)))
                return "state: buffer expected";
        return null;
    };

    /**
     * Creates an Initialize message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Initialize
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Initialize} Initialize
     */
    Initialize.fromObject = function fromObject(object) {
        if (object instanceof $root.Initialize)
            return object;
        var message = new $root.Initialize();
        if (object.state != null)
            if (typeof object.state === "string")
                $util.base64.decode(object.state, message.state = $util.newBuffer($util.base64.length(object.state)), 0);
            else if (object.state.length)
                message.state = object.state;
        return message;
    };

    /**
     * Creates a plain object from an Initialize message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Initialize
     * @static
     * @param {Initialize} message Initialize
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Initialize.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            if (options.bytes === String)
                object.state = "";
            else {
                object.state = [];
                if (options.bytes !== Array)
                    object.state = $util.newBuffer(object.state);
            }
        if (message.state != null && message.hasOwnProperty("state"))
            object.state = options.bytes === String ? $util.base64.encode(message.state, 0, message.state.length) : options.bytes === Array ? Array.prototype.slice.call(message.state) : message.state;
        return object;
    };

    /**
     * Converts this Initialize to JSON.
     * @function toJSON
     * @memberof Initialize
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Initialize.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Initialize;
})();

$root.GetFeatures = (function() {

    /**
     * Properties of a GetFeatures.
     * @exports IGetFeatures
     * @interface IGetFeatures
     */

    /**
     * Constructs a new GetFeatures.
     * @exports GetFeatures
     * @classdesc Request: Ask for device details (no device reset)
     * @next Features
     * @implements IGetFeatures
     * @constructor
     * @param {IGetFeatures=} [properties] Properties to set
     */
    function GetFeatures(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new GetFeatures instance using the specified properties.
     * @function create
     * @memberof GetFeatures
     * @static
     * @param {IGetFeatures=} [properties] Properties to set
     * @returns {GetFeatures} GetFeatures instance
     */
    GetFeatures.create = function create(properties) {
        return new GetFeatures(properties);
    };

    /**
     * Encodes the specified GetFeatures message. Does not implicitly {@link GetFeatures.verify|verify} messages.
     * @function encode
     * @memberof GetFeatures
     * @static
     * @param {IGetFeatures} message GetFeatures message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GetFeatures.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified GetFeatures message, length delimited. Does not implicitly {@link GetFeatures.verify|verify} messages.
     * @function encodeDelimited
     * @memberof GetFeatures
     * @static
     * @param {IGetFeatures} message GetFeatures message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GetFeatures.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a GetFeatures message from the specified reader or buffer.
     * @function decode
     * @memberof GetFeatures
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {GetFeatures} GetFeatures
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GetFeatures.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.GetFeatures();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a GetFeatures message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof GetFeatures
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {GetFeatures} GetFeatures
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GetFeatures.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a GetFeatures message.
     * @function verify
     * @memberof GetFeatures
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    GetFeatures.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a GetFeatures message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof GetFeatures
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {GetFeatures} GetFeatures
     */
    GetFeatures.fromObject = function fromObject(object) {
        if (object instanceof $root.GetFeatures)
            return object;
        return new $root.GetFeatures();
    };

    /**
     * Creates a plain object from a GetFeatures message. Also converts values to other types if specified.
     * @function toObject
     * @memberof GetFeatures
     * @static
     * @param {GetFeatures} message GetFeatures
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    GetFeatures.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this GetFeatures to JSON.
     * @function toJSON
     * @memberof GetFeatures
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    GetFeatures.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return GetFeatures;
})();

$root.Features = (function() {

    /**
     * Properties of a Features.
     * @exports IFeatures
     * @interface IFeatures
     * @property {string|null} [vendor] Features vendor
     * @property {number|null} [majorVersion] Features majorVersion
     * @property {number|null} [minorVersion] Features minorVersion
     * @property {number|null} [patchVersion] Features patchVersion
     * @property {boolean|null} [bootloaderMode] Features bootloaderMode
     * @property {string|null} [deviceId] Features deviceId
     * @property {boolean|null} [pinProtection] Features pinProtection
     * @property {boolean|null} [passphraseProtection] Features passphraseProtection
     * @property {string|null} [language] Features language
     * @property {string|null} [label] Features label
     * @property {Array.<ICoinType>|null} [coins] Features coins
     * @property {boolean|null} [initialized] Features initialized
     * @property {Uint8Array|null} [revision] Features revision
     * @property {Uint8Array|null} [bootloaderHash] Features bootloaderHash
     * @property {boolean|null} [imported] Features imported
     * @property {boolean|null} [pinCached] Features pinCached
     * @property {boolean|null} [passphraseCached] Features passphraseCached
     * @property {boolean|null} [firmwarePresent] Features firmwarePresent
     * @property {boolean|null} [needsBackup] Features needsBackup
     * @property {number|null} [flags] Features flags
     * @property {string|null} [model] Features model
     * @property {number|null} [fwMajor] Features fwMajor
     * @property {number|null} [fwMinor] Features fwMinor
     * @property {number|null} [fwPatch] Features fwPatch
     * @property {string|null} [fwVendor] Features fwVendor
     * @property {Uint8Array|null} [fwVendorKeys] Features fwVendorKeys
     * @property {boolean|null} [unfinishedBackup] Features unfinishedBackup
     */

    /**
     * Constructs a new Features.
     * @exports Features
     * @classdesc Response: Reports various information about the device
     * @prev Initialize
     * @prev GetFeatures
     * @implements IFeatures
     * @constructor
     * @param {IFeatures=} [properties] Properties to set
     */
    function Features(properties) {
        this.coins = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Features vendor.
     * @member {string} vendor
     * @memberof Features
     * @instance
     */
    Features.prototype.vendor = "";

    /**
     * Features majorVersion.
     * @member {number} majorVersion
     * @memberof Features
     * @instance
     */
    Features.prototype.majorVersion = 0;

    /**
     * Features minorVersion.
     * @member {number} minorVersion
     * @memberof Features
     * @instance
     */
    Features.prototype.minorVersion = 0;

    /**
     * Features patchVersion.
     * @member {number} patchVersion
     * @memberof Features
     * @instance
     */
    Features.prototype.patchVersion = 0;

    /**
     * Features bootloaderMode.
     * @member {boolean} bootloaderMode
     * @memberof Features
     * @instance
     */
    Features.prototype.bootloaderMode = false;

    /**
     * Features deviceId.
     * @member {string} deviceId
     * @memberof Features
     * @instance
     */
    Features.prototype.deviceId = "";

    /**
     * Features pinProtection.
     * @member {boolean} pinProtection
     * @memberof Features
     * @instance
     */
    Features.prototype.pinProtection = false;

    /**
     * Features passphraseProtection.
     * @member {boolean} passphraseProtection
     * @memberof Features
     * @instance
     */
    Features.prototype.passphraseProtection = false;

    /**
     * Features language.
     * @member {string} language
     * @memberof Features
     * @instance
     */
    Features.prototype.language = "";

    /**
     * Features label.
     * @member {string} label
     * @memberof Features
     * @instance
     */
    Features.prototype.label = "";

    /**
     * Features coins.
     * @member {Array.<ICoinType>} coins
     * @memberof Features
     * @instance
     */
    Features.prototype.coins = $util.emptyArray;

    /**
     * Features initialized.
     * @member {boolean} initialized
     * @memberof Features
     * @instance
     */
    Features.prototype.initialized = false;

    /**
     * Features revision.
     * @member {Uint8Array} revision
     * @memberof Features
     * @instance
     */
    Features.prototype.revision = $util.newBuffer([]);

    /**
     * Features bootloaderHash.
     * @member {Uint8Array} bootloaderHash
     * @memberof Features
     * @instance
     */
    Features.prototype.bootloaderHash = $util.newBuffer([]);

    /**
     * Features imported.
     * @member {boolean} imported
     * @memberof Features
     * @instance
     */
    Features.prototype.imported = false;

    /**
     * Features pinCached.
     * @member {boolean} pinCached
     * @memberof Features
     * @instance
     */
    Features.prototype.pinCached = false;

    /**
     * Features passphraseCached.
     * @member {boolean} passphraseCached
     * @memberof Features
     * @instance
     */
    Features.prototype.passphraseCached = false;

    /**
     * Features firmwarePresent.
     * @member {boolean} firmwarePresent
     * @memberof Features
     * @instance
     */
    Features.prototype.firmwarePresent = false;

    /**
     * Features needsBackup.
     * @member {boolean} needsBackup
     * @memberof Features
     * @instance
     */
    Features.prototype.needsBackup = false;

    /**
     * Features flags.
     * @member {number} flags
     * @memberof Features
     * @instance
     */
    Features.prototype.flags = 0;

    /**
     * Features model.
     * @member {string} model
     * @memberof Features
     * @instance
     */
    Features.prototype.model = "";

    /**
     * Features fwMajor.
     * @member {number} fwMajor
     * @memberof Features
     * @instance
     */
    Features.prototype.fwMajor = 0;

    /**
     * Features fwMinor.
     * @member {number} fwMinor
     * @memberof Features
     * @instance
     */
    Features.prototype.fwMinor = 0;

    /**
     * Features fwPatch.
     * @member {number} fwPatch
     * @memberof Features
     * @instance
     */
    Features.prototype.fwPatch = 0;

    /**
     * Features fwVendor.
     * @member {string} fwVendor
     * @memberof Features
     * @instance
     */
    Features.prototype.fwVendor = "";

    /**
     * Features fwVendorKeys.
     * @member {Uint8Array} fwVendorKeys
     * @memberof Features
     * @instance
     */
    Features.prototype.fwVendorKeys = $util.newBuffer([]);

    /**
     * Features unfinishedBackup.
     * @member {boolean} unfinishedBackup
     * @memberof Features
     * @instance
     */
    Features.prototype.unfinishedBackup = false;

    /**
     * Creates a new Features instance using the specified properties.
     * @function create
     * @memberof Features
     * @static
     * @param {IFeatures=} [properties] Properties to set
     * @returns {Features} Features instance
     */
    Features.create = function create(properties) {
        return new Features(properties);
    };

    /**
     * Encodes the specified Features message. Does not implicitly {@link Features.verify|verify} messages.
     * @function encode
     * @memberof Features
     * @static
     * @param {IFeatures} message Features message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Features.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.vendor != null && message.hasOwnProperty("vendor"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.vendor);
        if (message.majorVersion != null && message.hasOwnProperty("majorVersion"))
            writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.majorVersion);
        if (message.minorVersion != null && message.hasOwnProperty("minorVersion"))
            writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.minorVersion);
        if (message.patchVersion != null && message.hasOwnProperty("patchVersion"))
            writer.uint32(/* id 4, wireType 0 =*/32).uint32(message.patchVersion);
        if (message.bootloaderMode != null && message.hasOwnProperty("bootloaderMode"))
            writer.uint32(/* id 5, wireType 0 =*/40).bool(message.bootloaderMode);
        if (message.deviceId != null && message.hasOwnProperty("deviceId"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.deviceId);
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            writer.uint32(/* id 7, wireType 0 =*/56).bool(message.pinProtection);
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            writer.uint32(/* id 8, wireType 0 =*/64).bool(message.passphraseProtection);
        if (message.language != null && message.hasOwnProperty("language"))
            writer.uint32(/* id 9, wireType 2 =*/74).string(message.language);
        if (message.label != null && message.hasOwnProperty("label"))
            writer.uint32(/* id 10, wireType 2 =*/82).string(message.label);
        if (message.coins != null && message.coins.length)
            for (var i = 0; i < message.coins.length; ++i)
                $root.CoinType.encode(message.coins[i], writer.uint32(/* id 11, wireType 2 =*/90).fork()).ldelim();
        if (message.initialized != null && message.hasOwnProperty("initialized"))
            writer.uint32(/* id 12, wireType 0 =*/96).bool(message.initialized);
        if (message.revision != null && message.hasOwnProperty("revision"))
            writer.uint32(/* id 13, wireType 2 =*/106).bytes(message.revision);
        if (message.bootloaderHash != null && message.hasOwnProperty("bootloaderHash"))
            writer.uint32(/* id 14, wireType 2 =*/114).bytes(message.bootloaderHash);
        if (message.imported != null && message.hasOwnProperty("imported"))
            writer.uint32(/* id 15, wireType 0 =*/120).bool(message.imported);
        if (message.pinCached != null && message.hasOwnProperty("pinCached"))
            writer.uint32(/* id 16, wireType 0 =*/128).bool(message.pinCached);
        if (message.passphraseCached != null && message.hasOwnProperty("passphraseCached"))
            writer.uint32(/* id 17, wireType 0 =*/136).bool(message.passphraseCached);
        if (message.firmwarePresent != null && message.hasOwnProperty("firmwarePresent"))
            writer.uint32(/* id 18, wireType 0 =*/144).bool(message.firmwarePresent);
        if (message.needsBackup != null && message.hasOwnProperty("needsBackup"))
            writer.uint32(/* id 19, wireType 0 =*/152).bool(message.needsBackup);
        if (message.flags != null && message.hasOwnProperty("flags"))
            writer.uint32(/* id 20, wireType 0 =*/160).uint32(message.flags);
        if (message.model != null && message.hasOwnProperty("model"))
            writer.uint32(/* id 21, wireType 2 =*/170).string(message.model);
        if (message.fwMajor != null && message.hasOwnProperty("fwMajor"))
            writer.uint32(/* id 22, wireType 0 =*/176).uint32(message.fwMajor);
        if (message.fwMinor != null && message.hasOwnProperty("fwMinor"))
            writer.uint32(/* id 23, wireType 0 =*/184).uint32(message.fwMinor);
        if (message.fwPatch != null && message.hasOwnProperty("fwPatch"))
            writer.uint32(/* id 24, wireType 0 =*/192).uint32(message.fwPatch);
        if (message.fwVendor != null && message.hasOwnProperty("fwVendor"))
            writer.uint32(/* id 25, wireType 2 =*/202).string(message.fwVendor);
        if (message.fwVendorKeys != null && message.hasOwnProperty("fwVendorKeys"))
            writer.uint32(/* id 26, wireType 2 =*/210).bytes(message.fwVendorKeys);
        if (message.unfinishedBackup != null && message.hasOwnProperty("unfinishedBackup"))
            writer.uint32(/* id 27, wireType 0 =*/216).bool(message.unfinishedBackup);
        return writer;
    };

    /**
     * Encodes the specified Features message, length delimited. Does not implicitly {@link Features.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Features
     * @static
     * @param {IFeatures} message Features message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Features.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Features message from the specified reader or buffer.
     * @function decode
     * @memberof Features
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Features} Features
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Features.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Features();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.vendor = reader.string();
                break;
            case 2:
                message.majorVersion = reader.uint32();
                break;
            case 3:
                message.minorVersion = reader.uint32();
                break;
            case 4:
                message.patchVersion = reader.uint32();
                break;
            case 5:
                message.bootloaderMode = reader.bool();
                break;
            case 6:
                message.deviceId = reader.string();
                break;
            case 7:
                message.pinProtection = reader.bool();
                break;
            case 8:
                message.passphraseProtection = reader.bool();
                break;
            case 9:
                message.language = reader.string();
                break;
            case 10:
                message.label = reader.string();
                break;
            case 11:
                if (!(message.coins && message.coins.length))
                    message.coins = [];
                message.coins.push($root.CoinType.decode(reader, reader.uint32()));
                break;
            case 12:
                message.initialized = reader.bool();
                break;
            case 13:
                message.revision = reader.bytes();
                break;
            case 14:
                message.bootloaderHash = reader.bytes();
                break;
            case 15:
                message.imported = reader.bool();
                break;
            case 16:
                message.pinCached = reader.bool();
                break;
            case 17:
                message.passphraseCached = reader.bool();
                break;
            case 18:
                message.firmwarePresent = reader.bool();
                break;
            case 19:
                message.needsBackup = reader.bool();
                break;
            case 20:
                message.flags = reader.uint32();
                break;
            case 21:
                message.model = reader.string();
                break;
            case 22:
                message.fwMajor = reader.uint32();
                break;
            case 23:
                message.fwMinor = reader.uint32();
                break;
            case 24:
                message.fwPatch = reader.uint32();
                break;
            case 25:
                message.fwVendor = reader.string();
                break;
            case 26:
                message.fwVendorKeys = reader.bytes();
                break;
            case 27:
                message.unfinishedBackup = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Features message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Features
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Features} Features
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Features.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Features message.
     * @function verify
     * @memberof Features
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Features.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.vendor != null && message.hasOwnProperty("vendor"))
            if (!$util.isString(message.vendor))
                return "vendor: string expected";
        if (message.majorVersion != null && message.hasOwnProperty("majorVersion"))
            if (!$util.isInteger(message.majorVersion))
                return "majorVersion: integer expected";
        if (message.minorVersion != null && message.hasOwnProperty("minorVersion"))
            if (!$util.isInteger(message.minorVersion))
                return "minorVersion: integer expected";
        if (message.patchVersion != null && message.hasOwnProperty("patchVersion"))
            if (!$util.isInteger(message.patchVersion))
                return "patchVersion: integer expected";
        if (message.bootloaderMode != null && message.hasOwnProperty("bootloaderMode"))
            if (typeof message.bootloaderMode !== "boolean")
                return "bootloaderMode: boolean expected";
        if (message.deviceId != null && message.hasOwnProperty("deviceId"))
            if (!$util.isString(message.deviceId))
                return "deviceId: string expected";
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            if (typeof message.pinProtection !== "boolean")
                return "pinProtection: boolean expected";
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            if (typeof message.passphraseProtection !== "boolean")
                return "passphraseProtection: boolean expected";
        if (message.language != null && message.hasOwnProperty("language"))
            if (!$util.isString(message.language))
                return "language: string expected";
        if (message.label != null && message.hasOwnProperty("label"))
            if (!$util.isString(message.label))
                return "label: string expected";
        if (message.coins != null && message.hasOwnProperty("coins")) {
            if (!Array.isArray(message.coins))
                return "coins: array expected";
            for (var i = 0; i < message.coins.length; ++i) {
                var error = $root.CoinType.verify(message.coins[i]);
                if (error)
                    return "coins." + error;
            }
        }
        if (message.initialized != null && message.hasOwnProperty("initialized"))
            if (typeof message.initialized !== "boolean")
                return "initialized: boolean expected";
        if (message.revision != null && message.hasOwnProperty("revision"))
            if (!(message.revision && typeof message.revision.length === "number" || $util.isString(message.revision)))
                return "revision: buffer expected";
        if (message.bootloaderHash != null && message.hasOwnProperty("bootloaderHash"))
            if (!(message.bootloaderHash && typeof message.bootloaderHash.length === "number" || $util.isString(message.bootloaderHash)))
                return "bootloaderHash: buffer expected";
        if (message.imported != null && message.hasOwnProperty("imported"))
            if (typeof message.imported !== "boolean")
                return "imported: boolean expected";
        if (message.pinCached != null && message.hasOwnProperty("pinCached"))
            if (typeof message.pinCached !== "boolean")
                return "pinCached: boolean expected";
        if (message.passphraseCached != null && message.hasOwnProperty("passphraseCached"))
            if (typeof message.passphraseCached !== "boolean")
                return "passphraseCached: boolean expected";
        if (message.firmwarePresent != null && message.hasOwnProperty("firmwarePresent"))
            if (typeof message.firmwarePresent !== "boolean")
                return "firmwarePresent: boolean expected";
        if (message.needsBackup != null && message.hasOwnProperty("needsBackup"))
            if (typeof message.needsBackup !== "boolean")
                return "needsBackup: boolean expected";
        if (message.flags != null && message.hasOwnProperty("flags"))
            if (!$util.isInteger(message.flags))
                return "flags: integer expected";
        if (message.model != null && message.hasOwnProperty("model"))
            if (!$util.isString(message.model))
                return "model: string expected";
        if (message.fwMajor != null && message.hasOwnProperty("fwMajor"))
            if (!$util.isInteger(message.fwMajor))
                return "fwMajor: integer expected";
        if (message.fwMinor != null && message.hasOwnProperty("fwMinor"))
            if (!$util.isInteger(message.fwMinor))
                return "fwMinor: integer expected";
        if (message.fwPatch != null && message.hasOwnProperty("fwPatch"))
            if (!$util.isInteger(message.fwPatch))
                return "fwPatch: integer expected";
        if (message.fwVendor != null && message.hasOwnProperty("fwVendor"))
            if (!$util.isString(message.fwVendor))
                return "fwVendor: string expected";
        if (message.fwVendorKeys != null && message.hasOwnProperty("fwVendorKeys"))
            if (!(message.fwVendorKeys && typeof message.fwVendorKeys.length === "number" || $util.isString(message.fwVendorKeys)))
                return "fwVendorKeys: buffer expected";
        if (message.unfinishedBackup != null && message.hasOwnProperty("unfinishedBackup"))
            if (typeof message.unfinishedBackup !== "boolean")
                return "unfinishedBackup: boolean expected";
        return null;
    };

    /**
     * Creates a Features message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Features
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Features} Features
     */
    Features.fromObject = function fromObject(object) {
        if (object instanceof $root.Features)
            return object;
        var message = new $root.Features();
        if (object.vendor != null)
            message.vendor = String(object.vendor);
        if (object.majorVersion != null)
            message.majorVersion = object.majorVersion >>> 0;
        if (object.minorVersion != null)
            message.minorVersion = object.minorVersion >>> 0;
        if (object.patchVersion != null)
            message.patchVersion = object.patchVersion >>> 0;
        if (object.bootloaderMode != null)
            message.bootloaderMode = Boolean(object.bootloaderMode);
        if (object.deviceId != null)
            message.deviceId = String(object.deviceId);
        if (object.pinProtection != null)
            message.pinProtection = Boolean(object.pinProtection);
        if (object.passphraseProtection != null)
            message.passphraseProtection = Boolean(object.passphraseProtection);
        if (object.language != null)
            message.language = String(object.language);
        if (object.label != null)
            message.label = String(object.label);
        if (object.coins) {
            if (!Array.isArray(object.coins))
                throw TypeError(".Features.coins: array expected");
            message.coins = [];
            for (var i = 0; i < object.coins.length; ++i) {
                if (typeof object.coins[i] !== "object")
                    throw TypeError(".Features.coins: object expected");
                message.coins[i] = $root.CoinType.fromObject(object.coins[i]);
            }
        }
        if (object.initialized != null)
            message.initialized = Boolean(object.initialized);
        if (object.revision != null)
            if (typeof object.revision === "string")
                $util.base64.decode(object.revision, message.revision = $util.newBuffer($util.base64.length(object.revision)), 0);
            else if (object.revision.length)
                message.revision = object.revision;
        if (object.bootloaderHash != null)
            if (typeof object.bootloaderHash === "string")
                $util.base64.decode(object.bootloaderHash, message.bootloaderHash = $util.newBuffer($util.base64.length(object.bootloaderHash)), 0);
            else if (object.bootloaderHash.length)
                message.bootloaderHash = object.bootloaderHash;
        if (object.imported != null)
            message.imported = Boolean(object.imported);
        if (object.pinCached != null)
            message.pinCached = Boolean(object.pinCached);
        if (object.passphraseCached != null)
            message.passphraseCached = Boolean(object.passphraseCached);
        if (object.firmwarePresent != null)
            message.firmwarePresent = Boolean(object.firmwarePresent);
        if (object.needsBackup != null)
            message.needsBackup = Boolean(object.needsBackup);
        if (object.flags != null)
            message.flags = object.flags >>> 0;
        if (object.model != null)
            message.model = String(object.model);
        if (object.fwMajor != null)
            message.fwMajor = object.fwMajor >>> 0;
        if (object.fwMinor != null)
            message.fwMinor = object.fwMinor >>> 0;
        if (object.fwPatch != null)
            message.fwPatch = object.fwPatch >>> 0;
        if (object.fwVendor != null)
            message.fwVendor = String(object.fwVendor);
        if (object.fwVendorKeys != null)
            if (typeof object.fwVendorKeys === "string")
                $util.base64.decode(object.fwVendorKeys, message.fwVendorKeys = $util.newBuffer($util.base64.length(object.fwVendorKeys)), 0);
            else if (object.fwVendorKeys.length)
                message.fwVendorKeys = object.fwVendorKeys;
        if (object.unfinishedBackup != null)
            message.unfinishedBackup = Boolean(object.unfinishedBackup);
        return message;
    };

    /**
     * Creates a plain object from a Features message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Features
     * @static
     * @param {Features} message Features
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Features.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults)
            object.coins = [];
        if (options.defaults) {
            object.vendor = "";
            object.majorVersion = 0;
            object.minorVersion = 0;
            object.patchVersion = 0;
            object.bootloaderMode = false;
            object.deviceId = "";
            object.pinProtection = false;
            object.passphraseProtection = false;
            object.language = "";
            object.label = "";
            object.initialized = false;
            if (options.bytes === String)
                object.revision = "";
            else {
                object.revision = [];
                if (options.bytes !== Array)
                    object.revision = $util.newBuffer(object.revision);
            }
            if (options.bytes === String)
                object.bootloaderHash = "";
            else {
                object.bootloaderHash = [];
                if (options.bytes !== Array)
                    object.bootloaderHash = $util.newBuffer(object.bootloaderHash);
            }
            object.imported = false;
            object.pinCached = false;
            object.passphraseCached = false;
            object.firmwarePresent = false;
            object.needsBackup = false;
            object.flags = 0;
            object.model = "";
            object.fwMajor = 0;
            object.fwMinor = 0;
            object.fwPatch = 0;
            object.fwVendor = "";
            if (options.bytes === String)
                object.fwVendorKeys = "";
            else {
                object.fwVendorKeys = [];
                if (options.bytes !== Array)
                    object.fwVendorKeys = $util.newBuffer(object.fwVendorKeys);
            }
            object.unfinishedBackup = false;
        }
        if (message.vendor != null && message.hasOwnProperty("vendor"))
            object.vendor = message.vendor;
        if (message.majorVersion != null && message.hasOwnProperty("majorVersion"))
            object.majorVersion = message.majorVersion;
        if (message.minorVersion != null && message.hasOwnProperty("minorVersion"))
            object.minorVersion = message.minorVersion;
        if (message.patchVersion != null && message.hasOwnProperty("patchVersion"))
            object.patchVersion = message.patchVersion;
        if (message.bootloaderMode != null && message.hasOwnProperty("bootloaderMode"))
            object.bootloaderMode = message.bootloaderMode;
        if (message.deviceId != null && message.hasOwnProperty("deviceId"))
            object.deviceId = message.deviceId;
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            object.pinProtection = message.pinProtection;
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            object.passphraseProtection = message.passphraseProtection;
        if (message.language != null && message.hasOwnProperty("language"))
            object.language = message.language;
        if (message.label != null && message.hasOwnProperty("label"))
            object.label = message.label;
        if (message.coins && message.coins.length) {
            object.coins = [];
            for (var j = 0; j < message.coins.length; ++j)
                object.coins[j] = $root.CoinType.toObject(message.coins[j], options);
        }
        if (message.initialized != null && message.hasOwnProperty("initialized"))
            object.initialized = message.initialized;
        if (message.revision != null && message.hasOwnProperty("revision"))
            object.revision = options.bytes === String ? $util.base64.encode(message.revision, 0, message.revision.length) : options.bytes === Array ? Array.prototype.slice.call(message.revision) : message.revision;
        if (message.bootloaderHash != null && message.hasOwnProperty("bootloaderHash"))
            object.bootloaderHash = options.bytes === String ? $util.base64.encode(message.bootloaderHash, 0, message.bootloaderHash.length) : options.bytes === Array ? Array.prototype.slice.call(message.bootloaderHash) : message.bootloaderHash;
        if (message.imported != null && message.hasOwnProperty("imported"))
            object.imported = message.imported;
        if (message.pinCached != null && message.hasOwnProperty("pinCached"))
            object.pinCached = message.pinCached;
        if (message.passphraseCached != null && message.hasOwnProperty("passphraseCached"))
            object.passphraseCached = message.passphraseCached;
        if (message.firmwarePresent != null && message.hasOwnProperty("firmwarePresent"))
            object.firmwarePresent = message.firmwarePresent;
        if (message.needsBackup != null && message.hasOwnProperty("needsBackup"))
            object.needsBackup = message.needsBackup;
        if (message.flags != null && message.hasOwnProperty("flags"))
            object.flags = message.flags;
        if (message.model != null && message.hasOwnProperty("model"))
            object.model = message.model;
        if (message.fwMajor != null && message.hasOwnProperty("fwMajor"))
            object.fwMajor = message.fwMajor;
        if (message.fwMinor != null && message.hasOwnProperty("fwMinor"))
            object.fwMinor = message.fwMinor;
        if (message.fwPatch != null && message.hasOwnProperty("fwPatch"))
            object.fwPatch = message.fwPatch;
        if (message.fwVendor != null && message.hasOwnProperty("fwVendor"))
            object.fwVendor = message.fwVendor;
        if (message.fwVendorKeys != null && message.hasOwnProperty("fwVendorKeys"))
            object.fwVendorKeys = options.bytes === String ? $util.base64.encode(message.fwVendorKeys, 0, message.fwVendorKeys.length) : options.bytes === Array ? Array.prototype.slice.call(message.fwVendorKeys) : message.fwVendorKeys;
        if (message.unfinishedBackup != null && message.hasOwnProperty("unfinishedBackup"))
            object.unfinishedBackup = message.unfinishedBackup;
        return object;
    };

    /**
     * Converts this Features to JSON.
     * @function toJSON
     * @memberof Features
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Features.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Features;
})();

$root.GenerateMnemonic = (function() {

    /**
     * Properties of a GenerateMnemonic.
     * @exports IGenerateMnemonic
     * @interface IGenerateMnemonic
     */

    /**
     * Constructs a new GenerateMnemonic.
     * @exports GenerateMnemonic
     * @classdesc Request: Ask the device to generate a mnemonic and configure itself with it
     * @next Success
     * @implements IGenerateMnemonic
     * @constructor
     * @param {IGenerateMnemonic=} [properties] Properties to set
     */
    function GenerateMnemonic(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new GenerateMnemonic instance using the specified properties.
     * @function create
     * @memberof GenerateMnemonic
     * @static
     * @param {IGenerateMnemonic=} [properties] Properties to set
     * @returns {GenerateMnemonic} GenerateMnemonic instance
     */
    GenerateMnemonic.create = function create(properties) {
        return new GenerateMnemonic(properties);
    };

    /**
     * Encodes the specified GenerateMnemonic message. Does not implicitly {@link GenerateMnemonic.verify|verify} messages.
     * @function encode
     * @memberof GenerateMnemonic
     * @static
     * @param {IGenerateMnemonic} message GenerateMnemonic message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GenerateMnemonic.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified GenerateMnemonic message, length delimited. Does not implicitly {@link GenerateMnemonic.verify|verify} messages.
     * @function encodeDelimited
     * @memberof GenerateMnemonic
     * @static
     * @param {IGenerateMnemonic} message GenerateMnemonic message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GenerateMnemonic.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a GenerateMnemonic message from the specified reader or buffer.
     * @function decode
     * @memberof GenerateMnemonic
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {GenerateMnemonic} GenerateMnemonic
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GenerateMnemonic.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.GenerateMnemonic();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a GenerateMnemonic message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof GenerateMnemonic
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {GenerateMnemonic} GenerateMnemonic
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GenerateMnemonic.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a GenerateMnemonic message.
     * @function verify
     * @memberof GenerateMnemonic
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    GenerateMnemonic.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a GenerateMnemonic message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof GenerateMnemonic
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {GenerateMnemonic} GenerateMnemonic
     */
    GenerateMnemonic.fromObject = function fromObject(object) {
        if (object instanceof $root.GenerateMnemonic)
            return object;
        return new $root.GenerateMnemonic();
    };

    /**
     * Creates a plain object from a GenerateMnemonic message. Also converts values to other types if specified.
     * @function toObject
     * @memberof GenerateMnemonic
     * @static
     * @param {GenerateMnemonic} message GenerateMnemonic
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    GenerateMnemonic.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this GenerateMnemonic to JSON.
     * @function toJSON
     * @memberof GenerateMnemonic
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    GenerateMnemonic.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return GenerateMnemonic;
})();

$root.SetMnemonic = (function() {

    /**
     * Properties of a SetMnemonic.
     * @exports ISetMnemonic
     * @interface ISetMnemonic
     * @property {string} mnemonic SetMnemonic mnemonic
     */

    /**
     * Constructs a new SetMnemonic.
     * @exports SetMnemonic
     * @classdesc Request: Send a mnemonic to the device
     * @next Success
     * @implements ISetMnemonic
     * @constructor
     * @param {ISetMnemonic=} [properties] Properties to set
     */
    function SetMnemonic(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * SetMnemonic mnemonic.
     * @member {string} mnemonic
     * @memberof SetMnemonic
     * @instance
     */
    SetMnemonic.prototype.mnemonic = "";

    /**
     * Creates a new SetMnemonic instance using the specified properties.
     * @function create
     * @memberof SetMnemonic
     * @static
     * @param {ISetMnemonic=} [properties] Properties to set
     * @returns {SetMnemonic} SetMnemonic instance
     */
    SetMnemonic.create = function create(properties) {
        return new SetMnemonic(properties);
    };

    /**
     * Encodes the specified SetMnemonic message. Does not implicitly {@link SetMnemonic.verify|verify} messages.
     * @function encode
     * @memberof SetMnemonic
     * @static
     * @param {ISetMnemonic} message SetMnemonic message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SetMnemonic.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 2 =*/10).string(message.mnemonic);
        return writer;
    };

    /**
     * Encodes the specified SetMnemonic message, length delimited. Does not implicitly {@link SetMnemonic.verify|verify} messages.
     * @function encodeDelimited
     * @memberof SetMnemonic
     * @static
     * @param {ISetMnemonic} message SetMnemonic message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SetMnemonic.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a SetMnemonic message from the specified reader or buffer.
     * @function decode
     * @memberof SetMnemonic
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {SetMnemonic} SetMnemonic
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SetMnemonic.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.SetMnemonic();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.mnemonic = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("mnemonic"))
            throw $util.ProtocolError("missing required 'mnemonic'", { instance: message });
        return message;
    };

    /**
     * Decodes a SetMnemonic message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof SetMnemonic
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {SetMnemonic} SetMnemonic
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SetMnemonic.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a SetMnemonic message.
     * @function verify
     * @memberof SetMnemonic
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    SetMnemonic.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isString(message.mnemonic))
            return "mnemonic: string expected";
        return null;
    };

    /**
     * Creates a SetMnemonic message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof SetMnemonic
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {SetMnemonic} SetMnemonic
     */
    SetMnemonic.fromObject = function fromObject(object) {
        if (object instanceof $root.SetMnemonic)
            return object;
        var message = new $root.SetMnemonic();
        if (object.mnemonic != null)
            message.mnemonic = String(object.mnemonic);
        return message;
    };

    /**
     * Creates a plain object from a SetMnemonic message. Also converts values to other types if specified.
     * @function toObject
     * @memberof SetMnemonic
     * @static
     * @param {SetMnemonic} message SetMnemonic
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    SetMnemonic.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.mnemonic = "";
        if (message.mnemonic != null && message.hasOwnProperty("mnemonic"))
            object.mnemonic = message.mnemonic;
        return object;
    };

    /**
     * Converts this SetMnemonic to JSON.
     * @function toJSON
     * @memberof SetMnemonic
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    SetMnemonic.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return SetMnemonic;
})();

$root.ChangePin = (function() {

    /**
     * Properties of a ChangePin.
     * @exports IChangePin
     * @interface IChangePin
     * @property {boolean|null} [remove] ChangePin remove
     */

    /**
     * Constructs a new ChangePin.
     * @exports ChangePin
     * @classdesc Request: Starts workflow for setting/changing/removing the PIN
     * @next ButtonRequest
     * @next PinMatrixRequest
     * @implements IChangePin
     * @constructor
     * @param {IChangePin=} [properties] Properties to set
     */
    function ChangePin(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ChangePin remove.
     * @member {boolean} remove
     * @memberof ChangePin
     * @instance
     */
    ChangePin.prototype.remove = false;

    /**
     * Creates a new ChangePin instance using the specified properties.
     * @function create
     * @memberof ChangePin
     * @static
     * @param {IChangePin=} [properties] Properties to set
     * @returns {ChangePin} ChangePin instance
     */
    ChangePin.create = function create(properties) {
        return new ChangePin(properties);
    };

    /**
     * Encodes the specified ChangePin message. Does not implicitly {@link ChangePin.verify|verify} messages.
     * @function encode
     * @memberof ChangePin
     * @static
     * @param {IChangePin} message ChangePin message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ChangePin.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.remove != null && message.hasOwnProperty("remove"))
            writer.uint32(/* id 1, wireType 0 =*/8).bool(message.remove);
        return writer;
    };

    /**
     * Encodes the specified ChangePin message, length delimited. Does not implicitly {@link ChangePin.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ChangePin
     * @static
     * @param {IChangePin} message ChangePin message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ChangePin.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ChangePin message from the specified reader or buffer.
     * @function decode
     * @memberof ChangePin
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ChangePin} ChangePin
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ChangePin.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ChangePin();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.remove = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ChangePin message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ChangePin
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ChangePin} ChangePin
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ChangePin.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ChangePin message.
     * @function verify
     * @memberof ChangePin
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ChangePin.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.remove != null && message.hasOwnProperty("remove"))
            if (typeof message.remove !== "boolean")
                return "remove: boolean expected";
        return null;
    };

    /**
     * Creates a ChangePin message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ChangePin
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ChangePin} ChangePin
     */
    ChangePin.fromObject = function fromObject(object) {
        if (object instanceof $root.ChangePin)
            return object;
        var message = new $root.ChangePin();
        if (object.remove != null)
            message.remove = Boolean(object.remove);
        return message;
    };

    /**
     * Creates a plain object from a ChangePin message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ChangePin
     * @static
     * @param {ChangePin} message ChangePin
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ChangePin.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.remove = false;
        if (message.remove != null && message.hasOwnProperty("remove"))
            object.remove = message.remove;
        return object;
    };

    /**
     * Converts this ChangePin to JSON.
     * @function toJSON
     * @memberof ChangePin
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ChangePin.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ChangePin;
})();

$root.SkycoinAddress = (function() {

    /**
     * Properties of a SkycoinAddress.
     * @exports ISkycoinAddress
     * @interface ISkycoinAddress
     * @property {number} addressN SkycoinAddress addressN
     * @property {number|null} [startIndex] SkycoinAddress startIndex
     */

    /**
     * Constructs a new SkycoinAddress.
     * @exports SkycoinAddress
     * @classdesc Request: Generate a Skycoin or a Bitcoin address from a seed, device sends back the address in a Success message
     * @next Failure
     * @next ResponseSkycoinAddress
     * @implements ISkycoinAddress
     * @constructor
     * @param {ISkycoinAddress=} [properties] Properties to set
     */
    function SkycoinAddress(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * SkycoinAddress addressN.
     * @member {number} addressN
     * @memberof SkycoinAddress
     * @instance
     */
    SkycoinAddress.prototype.addressN = 0;

    /**
     * SkycoinAddress startIndex.
     * @member {number} startIndex
     * @memberof SkycoinAddress
     * @instance
     */
    SkycoinAddress.prototype.startIndex = 0;

    /**
     * Creates a new SkycoinAddress instance using the specified properties.
     * @function create
     * @memberof SkycoinAddress
     * @static
     * @param {ISkycoinAddress=} [properties] Properties to set
     * @returns {SkycoinAddress} SkycoinAddress instance
     */
    SkycoinAddress.create = function create(properties) {
        return new SkycoinAddress(properties);
    };

    /**
     * Encodes the specified SkycoinAddress message. Does not implicitly {@link SkycoinAddress.verify|verify} messages.
     * @function encode
     * @memberof SkycoinAddress
     * @static
     * @param {ISkycoinAddress} message SkycoinAddress message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SkycoinAddress.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.addressN);
        if (message.startIndex != null && message.hasOwnProperty("startIndex"))
            writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.startIndex);
        return writer;
    };

    /**
     * Encodes the specified SkycoinAddress message, length delimited. Does not implicitly {@link SkycoinAddress.verify|verify} messages.
     * @function encodeDelimited
     * @memberof SkycoinAddress
     * @static
     * @param {ISkycoinAddress} message SkycoinAddress message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SkycoinAddress.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a SkycoinAddress message from the specified reader or buffer.
     * @function decode
     * @memberof SkycoinAddress
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {SkycoinAddress} SkycoinAddress
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SkycoinAddress.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.SkycoinAddress();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.addressN = reader.uint32();
                break;
            case 2:
                message.startIndex = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("addressN"))
            throw $util.ProtocolError("missing required 'addressN'", { instance: message });
        return message;
    };

    /**
     * Decodes a SkycoinAddress message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof SkycoinAddress
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {SkycoinAddress} SkycoinAddress
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SkycoinAddress.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a SkycoinAddress message.
     * @function verify
     * @memberof SkycoinAddress
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    SkycoinAddress.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isInteger(message.addressN))
            return "addressN: integer expected";
        if (message.startIndex != null && message.hasOwnProperty("startIndex"))
            if (!$util.isInteger(message.startIndex))
                return "startIndex: integer expected";
        return null;
    };

    /**
     * Creates a SkycoinAddress message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof SkycoinAddress
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {SkycoinAddress} SkycoinAddress
     */
    SkycoinAddress.fromObject = function fromObject(object) {
        if (object instanceof $root.SkycoinAddress)
            return object;
        var message = new $root.SkycoinAddress();
        if (object.addressN != null)
            message.addressN = object.addressN >>> 0;
        if (object.startIndex != null)
            message.startIndex = object.startIndex >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a SkycoinAddress message. Also converts values to other types if specified.
     * @function toObject
     * @memberof SkycoinAddress
     * @static
     * @param {SkycoinAddress} message SkycoinAddress
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    SkycoinAddress.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.addressN = 0;
            object.startIndex = 0;
        }
        if (message.addressN != null && message.hasOwnProperty("addressN"))
            object.addressN = message.addressN;
        if (message.startIndex != null && message.hasOwnProperty("startIndex"))
            object.startIndex = message.startIndex;
        return object;
    };

    /**
     * Converts this SkycoinAddress to JSON.
     * @function toJSON
     * @memberof SkycoinAddress
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    SkycoinAddress.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return SkycoinAddress;
})();

$root.ResponseSkycoinAddress = (function() {

    /**
     * Properties of a ResponseSkycoinAddress.
     * @exports IResponseSkycoinAddress
     * @interface IResponseSkycoinAddress
     * @property {Array.<string>|null} [addresses] ResponseSkycoinAddress addresses
     */

    /**
     * Constructs a new ResponseSkycoinAddress.
     * @exports ResponseSkycoinAddress
     * @classdesc Response: Return the generated skycoin address
     * @prev SkycoinAddress
     * @implements IResponseSkycoinAddress
     * @constructor
     * @param {IResponseSkycoinAddress=} [properties] Properties to set
     */
    function ResponseSkycoinAddress(properties) {
        this.addresses = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ResponseSkycoinAddress addresses.
     * @member {Array.<string>} addresses
     * @memberof ResponseSkycoinAddress
     * @instance
     */
    ResponseSkycoinAddress.prototype.addresses = $util.emptyArray;

    /**
     * Creates a new ResponseSkycoinAddress instance using the specified properties.
     * @function create
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {IResponseSkycoinAddress=} [properties] Properties to set
     * @returns {ResponseSkycoinAddress} ResponseSkycoinAddress instance
     */
    ResponseSkycoinAddress.create = function create(properties) {
        return new ResponseSkycoinAddress(properties);
    };

    /**
     * Encodes the specified ResponseSkycoinAddress message. Does not implicitly {@link ResponseSkycoinAddress.verify|verify} messages.
     * @function encode
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {IResponseSkycoinAddress} message ResponseSkycoinAddress message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ResponseSkycoinAddress.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.addresses != null && message.addresses.length)
            for (var i = 0; i < message.addresses.length; ++i)
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.addresses[i]);
        return writer;
    };

    /**
     * Encodes the specified ResponseSkycoinAddress message, length delimited. Does not implicitly {@link ResponseSkycoinAddress.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {IResponseSkycoinAddress} message ResponseSkycoinAddress message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ResponseSkycoinAddress.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ResponseSkycoinAddress message from the specified reader or buffer.
     * @function decode
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ResponseSkycoinAddress} ResponseSkycoinAddress
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ResponseSkycoinAddress.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ResponseSkycoinAddress();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.addresses && message.addresses.length))
                    message.addresses = [];
                message.addresses.push(reader.string());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ResponseSkycoinAddress message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ResponseSkycoinAddress} ResponseSkycoinAddress
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ResponseSkycoinAddress.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ResponseSkycoinAddress message.
     * @function verify
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ResponseSkycoinAddress.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.addresses != null && message.hasOwnProperty("addresses")) {
            if (!Array.isArray(message.addresses))
                return "addresses: array expected";
            for (var i = 0; i < message.addresses.length; ++i)
                if (!$util.isString(message.addresses[i]))
                    return "addresses: string[] expected";
        }
        return null;
    };

    /**
     * Creates a ResponseSkycoinAddress message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ResponseSkycoinAddress} ResponseSkycoinAddress
     */
    ResponseSkycoinAddress.fromObject = function fromObject(object) {
        if (object instanceof $root.ResponseSkycoinAddress)
            return object;
        var message = new $root.ResponseSkycoinAddress();
        if (object.addresses) {
            if (!Array.isArray(object.addresses))
                throw TypeError(".ResponseSkycoinAddress.addresses: array expected");
            message.addresses = [];
            for (var i = 0; i < object.addresses.length; ++i)
                message.addresses[i] = String(object.addresses[i]);
        }
        return message;
    };

    /**
     * Creates a plain object from a ResponseSkycoinAddress message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ResponseSkycoinAddress
     * @static
     * @param {ResponseSkycoinAddress} message ResponseSkycoinAddress
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ResponseSkycoinAddress.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults)
            object.addresses = [];
        if (message.addresses && message.addresses.length) {
            object.addresses = [];
            for (var j = 0; j < message.addresses.length; ++j)
                object.addresses[j] = message.addresses[j];
        }
        return object;
    };

    /**
     * Converts this ResponseSkycoinAddress to JSON.
     * @function toJSON
     * @memberof ResponseSkycoinAddress
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ResponseSkycoinAddress.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ResponseSkycoinAddress;
})();

$root.SkycoinCheckMessageSignature = (function() {

    /**
     * Properties of a SkycoinCheckMessageSignature.
     * @exports ISkycoinCheckMessageSignature
     * @interface ISkycoinCheckMessageSignature
     * @property {string} address SkycoinCheckMessageSignature address
     * @property {string} message SkycoinCheckMessageSignature message
     * @property {string} signature SkycoinCheckMessageSignature signature
     */

    /**
     * Constructs a new SkycoinCheckMessageSignature.
     * @exports SkycoinCheckMessageSignature
     * @classdesc Request: Check a message signature matches the given address.
     * @next Success
     * @implements ISkycoinCheckMessageSignature
     * @constructor
     * @param {ISkycoinCheckMessageSignature=} [properties] Properties to set
     */
    function SkycoinCheckMessageSignature(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * SkycoinCheckMessageSignature address.
     * @member {string} address
     * @memberof SkycoinCheckMessageSignature
     * @instance
     */
    SkycoinCheckMessageSignature.prototype.address = "";

    /**
     * SkycoinCheckMessageSignature message.
     * @member {string} message
     * @memberof SkycoinCheckMessageSignature
     * @instance
     */
    SkycoinCheckMessageSignature.prototype.message = "";

    /**
     * SkycoinCheckMessageSignature signature.
     * @member {string} signature
     * @memberof SkycoinCheckMessageSignature
     * @instance
     */
    SkycoinCheckMessageSignature.prototype.signature = "";

    /**
     * Creates a new SkycoinCheckMessageSignature instance using the specified properties.
     * @function create
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {ISkycoinCheckMessageSignature=} [properties] Properties to set
     * @returns {SkycoinCheckMessageSignature} SkycoinCheckMessageSignature instance
     */
    SkycoinCheckMessageSignature.create = function create(properties) {
        return new SkycoinCheckMessageSignature(properties);
    };

    /**
     * Encodes the specified SkycoinCheckMessageSignature message. Does not implicitly {@link SkycoinCheckMessageSignature.verify|verify} messages.
     * @function encode
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {ISkycoinCheckMessageSignature} message SkycoinCheckMessageSignature message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SkycoinCheckMessageSignature.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 2 =*/10).string(message.address);
        writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
        writer.uint32(/* id 3, wireType 2 =*/26).string(message.signature);
        return writer;
    };

    /**
     * Encodes the specified SkycoinCheckMessageSignature message, length delimited. Does not implicitly {@link SkycoinCheckMessageSignature.verify|verify} messages.
     * @function encodeDelimited
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {ISkycoinCheckMessageSignature} message SkycoinCheckMessageSignature message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SkycoinCheckMessageSignature.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a SkycoinCheckMessageSignature message from the specified reader or buffer.
     * @function decode
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {SkycoinCheckMessageSignature} SkycoinCheckMessageSignature
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SkycoinCheckMessageSignature.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.SkycoinCheckMessageSignature();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.address = reader.string();
                break;
            case 2:
                message.message = reader.string();
                break;
            case 3:
                message.signature = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("address"))
            throw $util.ProtocolError("missing required 'address'", { instance: message });
        if (!message.hasOwnProperty("message"))
            throw $util.ProtocolError("missing required 'message'", { instance: message });
        if (!message.hasOwnProperty("signature"))
            throw $util.ProtocolError("missing required 'signature'", { instance: message });
        return message;
    };

    /**
     * Decodes a SkycoinCheckMessageSignature message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {SkycoinCheckMessageSignature} SkycoinCheckMessageSignature
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SkycoinCheckMessageSignature.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a SkycoinCheckMessageSignature message.
     * @function verify
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    SkycoinCheckMessageSignature.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isString(message.address))
            return "address: string expected";
        if (!$util.isString(message.message))
            return "message: string expected";
        if (!$util.isString(message.signature))
            return "signature: string expected";
        return null;
    };

    /**
     * Creates a SkycoinCheckMessageSignature message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {SkycoinCheckMessageSignature} SkycoinCheckMessageSignature
     */
    SkycoinCheckMessageSignature.fromObject = function fromObject(object) {
        if (object instanceof $root.SkycoinCheckMessageSignature)
            return object;
        var message = new $root.SkycoinCheckMessageSignature();
        if (object.address != null)
            message.address = String(object.address);
        if (object.message != null)
            message.message = String(object.message);
        if (object.signature != null)
            message.signature = String(object.signature);
        return message;
    };

    /**
     * Creates a plain object from a SkycoinCheckMessageSignature message. Also converts values to other types if specified.
     * @function toObject
     * @memberof SkycoinCheckMessageSignature
     * @static
     * @param {SkycoinCheckMessageSignature} message SkycoinCheckMessageSignature
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    SkycoinCheckMessageSignature.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.address = "";
            object.message = "";
            object.signature = "";
        }
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = message.message;
        if (message.signature != null && message.hasOwnProperty("signature"))
            object.signature = message.signature;
        return object;
    };

    /**
     * Converts this SkycoinCheckMessageSignature to JSON.
     * @function toJSON
     * @memberof SkycoinCheckMessageSignature
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    SkycoinCheckMessageSignature.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return SkycoinCheckMessageSignature;
})();

$root.SkycoinSignMessage = (function() {

    /**
     * Properties of a SkycoinSignMessage.
     * @exports ISkycoinSignMessage
     * @interface ISkycoinSignMessage
     * @property {number} addressN SkycoinSignMessage addressN
     * @property {string} message SkycoinSignMessage message
     */

    /**
     * Constructs a new SkycoinSignMessage.
     * @exports SkycoinSignMessage
     * @classdesc Request: Sign a message digest using the given secret key.
     * @next Failure
     * @next ResponseSkycoinSignMessage
     * @implements ISkycoinSignMessage
     * @constructor
     * @param {ISkycoinSignMessage=} [properties] Properties to set
     */
    function SkycoinSignMessage(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * SkycoinSignMessage addressN.
     * @member {number} addressN
     * @memberof SkycoinSignMessage
     * @instance
     */
    SkycoinSignMessage.prototype.addressN = 0;

    /**
     * SkycoinSignMessage message.
     * @member {string} message
     * @memberof SkycoinSignMessage
     * @instance
     */
    SkycoinSignMessage.prototype.message = "";

    /**
     * Creates a new SkycoinSignMessage instance using the specified properties.
     * @function create
     * @memberof SkycoinSignMessage
     * @static
     * @param {ISkycoinSignMessage=} [properties] Properties to set
     * @returns {SkycoinSignMessage} SkycoinSignMessage instance
     */
    SkycoinSignMessage.create = function create(properties) {
        return new SkycoinSignMessage(properties);
    };

    /**
     * Encodes the specified SkycoinSignMessage message. Does not implicitly {@link SkycoinSignMessage.verify|verify} messages.
     * @function encode
     * @memberof SkycoinSignMessage
     * @static
     * @param {ISkycoinSignMessage} message SkycoinSignMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SkycoinSignMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.addressN);
        writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
        return writer;
    };

    /**
     * Encodes the specified SkycoinSignMessage message, length delimited. Does not implicitly {@link SkycoinSignMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof SkycoinSignMessage
     * @static
     * @param {ISkycoinSignMessage} message SkycoinSignMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    SkycoinSignMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a SkycoinSignMessage message from the specified reader or buffer.
     * @function decode
     * @memberof SkycoinSignMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {SkycoinSignMessage} SkycoinSignMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SkycoinSignMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.SkycoinSignMessage();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.addressN = reader.uint32();
                break;
            case 2:
                message.message = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("addressN"))
            throw $util.ProtocolError("missing required 'addressN'", { instance: message });
        if (!message.hasOwnProperty("message"))
            throw $util.ProtocolError("missing required 'message'", { instance: message });
        return message;
    };

    /**
     * Decodes a SkycoinSignMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof SkycoinSignMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {SkycoinSignMessage} SkycoinSignMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    SkycoinSignMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a SkycoinSignMessage message.
     * @function verify
     * @memberof SkycoinSignMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    SkycoinSignMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isInteger(message.addressN))
            return "addressN: integer expected";
        if (!$util.isString(message.message))
            return "message: string expected";
        return null;
    };

    /**
     * Creates a SkycoinSignMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof SkycoinSignMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {SkycoinSignMessage} SkycoinSignMessage
     */
    SkycoinSignMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.SkycoinSignMessage)
            return object;
        var message = new $root.SkycoinSignMessage();
        if (object.addressN != null)
            message.addressN = object.addressN >>> 0;
        if (object.message != null)
            message.message = String(object.message);
        return message;
    };

    /**
     * Creates a plain object from a SkycoinSignMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof SkycoinSignMessage
     * @static
     * @param {SkycoinSignMessage} message SkycoinSignMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    SkycoinSignMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.addressN = 0;
            object.message = "";
        }
        if (message.addressN != null && message.hasOwnProperty("addressN"))
            object.addressN = message.addressN;
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = message.message;
        return object;
    };

    /**
     * Converts this SkycoinSignMessage to JSON.
     * @function toJSON
     * @memberof SkycoinSignMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    SkycoinSignMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return SkycoinSignMessage;
})();

$root.ResponseSkycoinSignMessage = (function() {

    /**
     * Properties of a ResponseSkycoinSignMessage.
     * @exports IResponseSkycoinSignMessage
     * @interface IResponseSkycoinSignMessage
     * @property {string} signedMessage ResponseSkycoinSignMessage signedMessage
     */

    /**
     * Constructs a new ResponseSkycoinSignMessage.
     * @exports ResponseSkycoinSignMessage
     * @classdesc Response: Return the generated skycoin address
     * @prev SkycoinAddress
     * @implements IResponseSkycoinSignMessage
     * @constructor
     * @param {IResponseSkycoinSignMessage=} [properties] Properties to set
     */
    function ResponseSkycoinSignMessage(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ResponseSkycoinSignMessage signedMessage.
     * @member {string} signedMessage
     * @memberof ResponseSkycoinSignMessage
     * @instance
     */
    ResponseSkycoinSignMessage.prototype.signedMessage = "";

    /**
     * Creates a new ResponseSkycoinSignMessage instance using the specified properties.
     * @function create
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {IResponseSkycoinSignMessage=} [properties] Properties to set
     * @returns {ResponseSkycoinSignMessage} ResponseSkycoinSignMessage instance
     */
    ResponseSkycoinSignMessage.create = function create(properties) {
        return new ResponseSkycoinSignMessage(properties);
    };

    /**
     * Encodes the specified ResponseSkycoinSignMessage message. Does not implicitly {@link ResponseSkycoinSignMessage.verify|verify} messages.
     * @function encode
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {IResponseSkycoinSignMessage} message ResponseSkycoinSignMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ResponseSkycoinSignMessage.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 2 =*/10).string(message.signedMessage);
        return writer;
    };

    /**
     * Encodes the specified ResponseSkycoinSignMessage message, length delimited. Does not implicitly {@link ResponseSkycoinSignMessage.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {IResponseSkycoinSignMessage} message ResponseSkycoinSignMessage message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ResponseSkycoinSignMessage.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ResponseSkycoinSignMessage message from the specified reader or buffer.
     * @function decode
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ResponseSkycoinSignMessage} ResponseSkycoinSignMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ResponseSkycoinSignMessage.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ResponseSkycoinSignMessage();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.signedMessage = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("signedMessage"))
            throw $util.ProtocolError("missing required 'signedMessage'", { instance: message });
        return message;
    };

    /**
     * Decodes a ResponseSkycoinSignMessage message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ResponseSkycoinSignMessage} ResponseSkycoinSignMessage
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ResponseSkycoinSignMessage.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ResponseSkycoinSignMessage message.
     * @function verify
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ResponseSkycoinSignMessage.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isString(message.signedMessage))
            return "signedMessage: string expected";
        return null;
    };

    /**
     * Creates a ResponseSkycoinSignMessage message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ResponseSkycoinSignMessage} ResponseSkycoinSignMessage
     */
    ResponseSkycoinSignMessage.fromObject = function fromObject(object) {
        if (object instanceof $root.ResponseSkycoinSignMessage)
            return object;
        var message = new $root.ResponseSkycoinSignMessage();
        if (object.signedMessage != null)
            message.signedMessage = String(object.signedMessage);
        return message;
    };

    /**
     * Creates a plain object from a ResponseSkycoinSignMessage message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ResponseSkycoinSignMessage
     * @static
     * @param {ResponseSkycoinSignMessage} message ResponseSkycoinSignMessage
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ResponseSkycoinSignMessage.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.signedMessage = "";
        if (message.signedMessage != null && message.hasOwnProperty("signedMessage"))
            object.signedMessage = message.signedMessage;
        return object;
    };

    /**
     * Converts this ResponseSkycoinSignMessage to JSON.
     * @function toJSON
     * @memberof ResponseSkycoinSignMessage
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ResponseSkycoinSignMessage.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ResponseSkycoinSignMessage;
})();

$root.Ping = (function() {

    /**
     * Properties of a Ping.
     * @exports IPing
     * @interface IPing
     * @property {string|null} [message] Ping message
     * @property {boolean|null} [buttonProtection] Ping buttonProtection
     * @property {boolean|null} [pinProtection] Ping pinProtection
     * @property {boolean|null} [passphraseProtection] Ping passphraseProtection
     */

    /**
     * Constructs a new Ping.
     * @exports Ping
     * @classdesc Request: Test if the device is alive, device sends back the message in Success response
     * @next Success
     * @implements IPing
     * @constructor
     * @param {IPing=} [properties] Properties to set
     */
    function Ping(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Ping message.
     * @member {string} message
     * @memberof Ping
     * @instance
     */
    Ping.prototype.message = "";

    /**
     * Ping buttonProtection.
     * @member {boolean} buttonProtection
     * @memberof Ping
     * @instance
     */
    Ping.prototype.buttonProtection = false;

    /**
     * Ping pinProtection.
     * @member {boolean} pinProtection
     * @memberof Ping
     * @instance
     */
    Ping.prototype.pinProtection = false;

    /**
     * Ping passphraseProtection.
     * @member {boolean} passphraseProtection
     * @memberof Ping
     * @instance
     */
    Ping.prototype.passphraseProtection = false;

    /**
     * Creates a new Ping instance using the specified properties.
     * @function create
     * @memberof Ping
     * @static
     * @param {IPing=} [properties] Properties to set
     * @returns {Ping} Ping instance
     */
    Ping.create = function create(properties) {
        return new Ping(properties);
    };

    /**
     * Encodes the specified Ping message. Does not implicitly {@link Ping.verify|verify} messages.
     * @function encode
     * @memberof Ping
     * @static
     * @param {IPing} message Ping message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Ping.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.message != null && message.hasOwnProperty("message"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.message);
        if (message.buttonProtection != null && message.hasOwnProperty("buttonProtection"))
            writer.uint32(/* id 2, wireType 0 =*/16).bool(message.buttonProtection);
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.pinProtection);
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            writer.uint32(/* id 4, wireType 0 =*/32).bool(message.passphraseProtection);
        return writer;
    };

    /**
     * Encodes the specified Ping message, length delimited. Does not implicitly {@link Ping.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Ping
     * @static
     * @param {IPing} message Ping message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Ping.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Ping message from the specified reader or buffer.
     * @function decode
     * @memberof Ping
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Ping} Ping
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Ping.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Ping();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.message = reader.string();
                break;
            case 2:
                message.buttonProtection = reader.bool();
                break;
            case 3:
                message.pinProtection = reader.bool();
                break;
            case 4:
                message.passphraseProtection = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Ping message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Ping
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Ping} Ping
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Ping.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Ping message.
     * @function verify
     * @memberof Ping
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Ping.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.message != null && message.hasOwnProperty("message"))
            if (!$util.isString(message.message))
                return "message: string expected";
        if (message.buttonProtection != null && message.hasOwnProperty("buttonProtection"))
            if (typeof message.buttonProtection !== "boolean")
                return "buttonProtection: boolean expected";
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            if (typeof message.pinProtection !== "boolean")
                return "pinProtection: boolean expected";
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            if (typeof message.passphraseProtection !== "boolean")
                return "passphraseProtection: boolean expected";
        return null;
    };

    /**
     * Creates a Ping message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Ping
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Ping} Ping
     */
    Ping.fromObject = function fromObject(object) {
        if (object instanceof $root.Ping)
            return object;
        var message = new $root.Ping();
        if (object.message != null)
            message.message = String(object.message);
        if (object.buttonProtection != null)
            message.buttonProtection = Boolean(object.buttonProtection);
        if (object.pinProtection != null)
            message.pinProtection = Boolean(object.pinProtection);
        if (object.passphraseProtection != null)
            message.passphraseProtection = Boolean(object.passphraseProtection);
        return message;
    };

    /**
     * Creates a plain object from a Ping message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Ping
     * @static
     * @param {Ping} message Ping
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Ping.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.message = "";
            object.buttonProtection = false;
            object.pinProtection = false;
            object.passphraseProtection = false;
        }
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = message.message;
        if (message.buttonProtection != null && message.hasOwnProperty("buttonProtection"))
            object.buttonProtection = message.buttonProtection;
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            object.pinProtection = message.pinProtection;
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            object.passphraseProtection = message.passphraseProtection;
        return object;
    };

    /**
     * Converts this Ping to JSON.
     * @function toJSON
     * @memberof Ping
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Ping.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Ping;
})();

$root.Success = (function() {

    /**
     * Properties of a Success.
     * @exports ISuccess
     * @interface ISuccess
     * @property {string|null} [message] Success message
     */

    /**
     * Constructs a new Success.
     * @exports Success
     * @classdesc Response: Success of the previous request
     * @implements ISuccess
     * @constructor
     * @param {ISuccess=} [properties] Properties to set
     */
    function Success(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Success message.
     * @member {string} message
     * @memberof Success
     * @instance
     */
    Success.prototype.message = "";

    /**
     * Creates a new Success instance using the specified properties.
     * @function create
     * @memberof Success
     * @static
     * @param {ISuccess=} [properties] Properties to set
     * @returns {Success} Success instance
     */
    Success.create = function create(properties) {
        return new Success(properties);
    };

    /**
     * Encodes the specified Success message. Does not implicitly {@link Success.verify|verify} messages.
     * @function encode
     * @memberof Success
     * @static
     * @param {ISuccess} message Success message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Success.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.message != null && message.hasOwnProperty("message"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.message);
        return writer;
    };

    /**
     * Encodes the specified Success message, length delimited. Does not implicitly {@link Success.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Success
     * @static
     * @param {ISuccess} message Success message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Success.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Success message from the specified reader or buffer.
     * @function decode
     * @memberof Success
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Success} Success
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Success.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Success();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.message = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Success message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Success
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Success} Success
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Success.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Success message.
     * @function verify
     * @memberof Success
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Success.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.message != null && message.hasOwnProperty("message"))
            if (!$util.isString(message.message))
                return "message: string expected";
        return null;
    };

    /**
     * Creates a Success message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Success
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Success} Success
     */
    Success.fromObject = function fromObject(object) {
        if (object instanceof $root.Success)
            return object;
        var message = new $root.Success();
        if (object.message != null)
            message.message = String(object.message);
        return message;
    };

    /**
     * Creates a plain object from a Success message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Success
     * @static
     * @param {Success} message Success
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Success.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.message = "";
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = message.message;
        return object;
    };

    /**
     * Converts this Success to JSON.
     * @function toJSON
     * @memberof Success
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Success.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Success;
})();

$root.Failure = (function() {

    /**
     * Properties of a Failure.
     * @exports IFailure
     * @interface IFailure
     * @property {FailureType|null} [code] Failure code
     * @property {string|null} [message] Failure message
     */

    /**
     * Constructs a new Failure.
     * @exports Failure
     * @classdesc Response: Failure of the previous request
     * @implements IFailure
     * @constructor
     * @param {IFailure=} [properties] Properties to set
     */
    function Failure(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Failure code.
     * @member {FailureType} code
     * @memberof Failure
     * @instance
     */
    Failure.prototype.code = 1;

    /**
     * Failure message.
     * @member {string} message
     * @memberof Failure
     * @instance
     */
    Failure.prototype.message = "";

    /**
     * Creates a new Failure instance using the specified properties.
     * @function create
     * @memberof Failure
     * @static
     * @param {IFailure=} [properties] Properties to set
     * @returns {Failure} Failure instance
     */
    Failure.create = function create(properties) {
        return new Failure(properties);
    };

    /**
     * Encodes the specified Failure message. Does not implicitly {@link Failure.verify|verify} messages.
     * @function encode
     * @memberof Failure
     * @static
     * @param {IFailure} message Failure message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Failure.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.code != null && message.hasOwnProperty("code"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.code);
        if (message.message != null && message.hasOwnProperty("message"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.message);
        return writer;
    };

    /**
     * Encodes the specified Failure message, length delimited. Does not implicitly {@link Failure.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Failure
     * @static
     * @param {IFailure} message Failure message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Failure.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Failure message from the specified reader or buffer.
     * @function decode
     * @memberof Failure
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Failure} Failure
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Failure.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Failure();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.code = reader.int32();
                break;
            case 2:
                message.message = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Failure message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Failure
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Failure} Failure
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Failure.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Failure message.
     * @function verify
     * @memberof Failure
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Failure.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.code != null && message.hasOwnProperty("code"))
            switch (message.code) {
            default:
                return "code: enum value expected";
            case 1:
            case 2:
            case 3:
            case 4:
            case 5:
            case 6:
            case 7:
            case 8:
            case 9:
            case 10:
            case 11:
            case 12:
            case 13:
            case 99:
                break;
            }
        if (message.message != null && message.hasOwnProperty("message"))
            if (!$util.isString(message.message))
                return "message: string expected";
        return null;
    };

    /**
     * Creates a Failure message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Failure
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Failure} Failure
     */
    Failure.fromObject = function fromObject(object) {
        if (object instanceof $root.Failure)
            return object;
        var message = new $root.Failure();
        switch (object.code) {
        case "Failure_UnexpectedMessage":
        case 1:
            message.code = 1;
            break;
        case "Failure_ButtonExpected":
        case 2:
            message.code = 2;
            break;
        case "Failure_DataError":
        case 3:
            message.code = 3;
            break;
        case "Failure_ActionCancelled":
        case 4:
            message.code = 4;
            break;
        case "Failure_PinExpected":
        case 5:
            message.code = 5;
            break;
        case "Failure_PinCancelled":
        case 6:
            message.code = 6;
            break;
        case "Failure_PinInvalid":
        case 7:
            message.code = 7;
            break;
        case "Failure_InvalidSignature":
        case 8:
            message.code = 8;
            break;
        case "Failure_ProcessError":
        case 9:
            message.code = 9;
            break;
        case "Failure_NotEnoughFunds":
        case 10:
            message.code = 10;
            break;
        case "Failure_NotInitialized":
        case 11:
            message.code = 11;
            break;
        case "Failure_PinMismatch":
        case 12:
            message.code = 12;
            break;
        case "Failure_AddressGeneration":
        case 13:
            message.code = 13;
            break;
        case "Failure_FirmwareError":
        case 99:
            message.code = 99;
            break;
        }
        if (object.message != null)
            message.message = String(object.message);
        return message;
    };

    /**
     * Creates a plain object from a Failure message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Failure
     * @static
     * @param {Failure} message Failure
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Failure.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.code = options.enums === String ? "Failure_UnexpectedMessage" : 1;
            object.message = "";
        }
        if (message.code != null && message.hasOwnProperty("code"))
            object.code = options.enums === String ? $root.FailureType[message.code] : message.code;
        if (message.message != null && message.hasOwnProperty("message"))
            object.message = message.message;
        return object;
    };

    /**
     * Converts this Failure to JSON.
     * @function toJSON
     * @memberof Failure
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Failure.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Failure;
})();

$root.ButtonRequest = (function() {

    /**
     * Properties of a ButtonRequest.
     * @exports IButtonRequest
     * @interface IButtonRequest
     * @property {ButtonRequestType|null} [code] ButtonRequest code
     * @property {string|null} [data] ButtonRequest data
     */

    /**
     * Constructs a new ButtonRequest.
     * @exports ButtonRequest
     * @classdesc Response: Device is waiting for HW button press.
     * @next ButtonAck
     * @next Cancel
     * @implements IButtonRequest
     * @constructor
     * @param {IButtonRequest=} [properties] Properties to set
     */
    function ButtonRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ButtonRequest code.
     * @member {ButtonRequestType} code
     * @memberof ButtonRequest
     * @instance
     */
    ButtonRequest.prototype.code = 1;

    /**
     * ButtonRequest data.
     * @member {string} data
     * @memberof ButtonRequest
     * @instance
     */
    ButtonRequest.prototype.data = "";

    /**
     * Creates a new ButtonRequest instance using the specified properties.
     * @function create
     * @memberof ButtonRequest
     * @static
     * @param {IButtonRequest=} [properties] Properties to set
     * @returns {ButtonRequest} ButtonRequest instance
     */
    ButtonRequest.create = function create(properties) {
        return new ButtonRequest(properties);
    };

    /**
     * Encodes the specified ButtonRequest message. Does not implicitly {@link ButtonRequest.verify|verify} messages.
     * @function encode
     * @memberof ButtonRequest
     * @static
     * @param {IButtonRequest} message ButtonRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ButtonRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.code != null && message.hasOwnProperty("code"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.code);
        if (message.data != null && message.hasOwnProperty("data"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.data);
        return writer;
    };

    /**
     * Encodes the specified ButtonRequest message, length delimited. Does not implicitly {@link ButtonRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ButtonRequest
     * @static
     * @param {IButtonRequest} message ButtonRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ButtonRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ButtonRequest message from the specified reader or buffer.
     * @function decode
     * @memberof ButtonRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ButtonRequest} ButtonRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ButtonRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ButtonRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.code = reader.int32();
                break;
            case 2:
                message.data = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ButtonRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ButtonRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ButtonRequest} ButtonRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ButtonRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ButtonRequest message.
     * @function verify
     * @memberof ButtonRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ButtonRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.code != null && message.hasOwnProperty("code"))
            switch (message.code) {
            default:
                return "code: enum value expected";
            case 1:
            case 2:
            case 3:
            case 4:
            case 5:
            case 6:
            case 7:
            case 8:
            case 9:
            case 10:
            case 11:
            case 12:
            case 13:
            case 14:
                break;
            }
        if (message.data != null && message.hasOwnProperty("data"))
            if (!$util.isString(message.data))
                return "data: string expected";
        return null;
    };

    /**
     * Creates a ButtonRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ButtonRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ButtonRequest} ButtonRequest
     */
    ButtonRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.ButtonRequest)
            return object;
        var message = new $root.ButtonRequest();
        switch (object.code) {
        case "ButtonRequest_Other":
        case 1:
            message.code = 1;
            break;
        case "ButtonRequest_FeeOverThreshold":
        case 2:
            message.code = 2;
            break;
        case "ButtonRequest_ConfirmOutput":
        case 3:
            message.code = 3;
            break;
        case "ButtonRequest_ResetDevice":
        case 4:
            message.code = 4;
            break;
        case "ButtonRequest_ConfirmWord":
        case 5:
            message.code = 5;
            break;
        case "ButtonRequest_WipeDevice":
        case 6:
            message.code = 6;
            break;
        case "ButtonRequest_ProtectCall":
        case 7:
            message.code = 7;
            break;
        case "ButtonRequest_SignTx":
        case 8:
            message.code = 8;
            break;
        case "ButtonRequest_FirmwareCheck":
        case 9:
            message.code = 9;
            break;
        case "ButtonRequest_Address":
        case 10:
            message.code = 10;
            break;
        case "ButtonRequest_PublicKey":
        case 11:
            message.code = 11;
            break;
        case "ButtonRequest_MnemonicWordCount":
        case 12:
            message.code = 12;
            break;
        case "ButtonRequest_MnemonicInput":
        case 13:
            message.code = 13;
            break;
        case "ButtonRequest_PassphraseType":
        case 14:
            message.code = 14;
            break;
        }
        if (object.data != null)
            message.data = String(object.data);
        return message;
    };

    /**
     * Creates a plain object from a ButtonRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ButtonRequest
     * @static
     * @param {ButtonRequest} message ButtonRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ButtonRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.code = options.enums === String ? "ButtonRequest_Other" : 1;
            object.data = "";
        }
        if (message.code != null && message.hasOwnProperty("code"))
            object.code = options.enums === String ? $root.ButtonRequestType[message.code] : message.code;
        if (message.data != null && message.hasOwnProperty("data"))
            object.data = message.data;
        return object;
    };

    /**
     * Converts this ButtonRequest to JSON.
     * @function toJSON
     * @memberof ButtonRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ButtonRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ButtonRequest;
})();

$root.ButtonAck = (function() {

    /**
     * Properties of a ButtonAck.
     * @exports IButtonAck
     * @interface IButtonAck
     */

    /**
     * Constructs a new ButtonAck.
     * @exports ButtonAck
     * @classdesc Request: Computer agrees to wait for HW button press
     * @prev ButtonRequest
     * @implements IButtonAck
     * @constructor
     * @param {IButtonAck=} [properties] Properties to set
     */
    function ButtonAck(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new ButtonAck instance using the specified properties.
     * @function create
     * @memberof ButtonAck
     * @static
     * @param {IButtonAck=} [properties] Properties to set
     * @returns {ButtonAck} ButtonAck instance
     */
    ButtonAck.create = function create(properties) {
        return new ButtonAck(properties);
    };

    /**
     * Encodes the specified ButtonAck message. Does not implicitly {@link ButtonAck.verify|verify} messages.
     * @function encode
     * @memberof ButtonAck
     * @static
     * @param {IButtonAck} message ButtonAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ButtonAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified ButtonAck message, length delimited. Does not implicitly {@link ButtonAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ButtonAck
     * @static
     * @param {IButtonAck} message ButtonAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ButtonAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ButtonAck message from the specified reader or buffer.
     * @function decode
     * @memberof ButtonAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ButtonAck} ButtonAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ButtonAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ButtonAck();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ButtonAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ButtonAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ButtonAck} ButtonAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ButtonAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ButtonAck message.
     * @function verify
     * @memberof ButtonAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ButtonAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a ButtonAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ButtonAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ButtonAck} ButtonAck
     */
    ButtonAck.fromObject = function fromObject(object) {
        if (object instanceof $root.ButtonAck)
            return object;
        return new $root.ButtonAck();
    };

    /**
     * Creates a plain object from a ButtonAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ButtonAck
     * @static
     * @param {ButtonAck} message ButtonAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ButtonAck.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this ButtonAck to JSON.
     * @function toJSON
     * @memberof ButtonAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ButtonAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ButtonAck;
})();

$root.PinMatrixRequest = (function() {

    /**
     * Properties of a PinMatrixRequest.
     * @exports IPinMatrixRequest
     * @interface IPinMatrixRequest
     * @property {PinMatrixRequestType|null} [type] PinMatrixRequest type
     */

    /**
     * Constructs a new PinMatrixRequest.
     * @exports PinMatrixRequest
     * @classdesc Response: Device is asking computer to show PIN matrix and awaits PIN encoded using this matrix scheme
     * @next PinMatrixAck
     * @next Cancel
     * @implements IPinMatrixRequest
     * @constructor
     * @param {IPinMatrixRequest=} [properties] Properties to set
     */
    function PinMatrixRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PinMatrixRequest type.
     * @member {PinMatrixRequestType} type
     * @memberof PinMatrixRequest
     * @instance
     */
    PinMatrixRequest.prototype.type = 1;

    /**
     * Creates a new PinMatrixRequest instance using the specified properties.
     * @function create
     * @memberof PinMatrixRequest
     * @static
     * @param {IPinMatrixRequest=} [properties] Properties to set
     * @returns {PinMatrixRequest} PinMatrixRequest instance
     */
    PinMatrixRequest.create = function create(properties) {
        return new PinMatrixRequest(properties);
    };

    /**
     * Encodes the specified PinMatrixRequest message. Does not implicitly {@link PinMatrixRequest.verify|verify} messages.
     * @function encode
     * @memberof PinMatrixRequest
     * @static
     * @param {IPinMatrixRequest} message PinMatrixRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PinMatrixRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.type);
        return writer;
    };

    /**
     * Encodes the specified PinMatrixRequest message, length delimited. Does not implicitly {@link PinMatrixRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PinMatrixRequest
     * @static
     * @param {IPinMatrixRequest} message PinMatrixRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PinMatrixRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PinMatrixRequest message from the specified reader or buffer.
     * @function decode
     * @memberof PinMatrixRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PinMatrixRequest} PinMatrixRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PinMatrixRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.PinMatrixRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.type = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PinMatrixRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PinMatrixRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PinMatrixRequest} PinMatrixRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PinMatrixRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PinMatrixRequest message.
     * @function verify
     * @memberof PinMatrixRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PinMatrixRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 1:
            case 2:
            case 3:
                break;
            }
        return null;
    };

    /**
     * Creates a PinMatrixRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PinMatrixRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PinMatrixRequest} PinMatrixRequest
     */
    PinMatrixRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.PinMatrixRequest)
            return object;
        var message = new $root.PinMatrixRequest();
        switch (object.type) {
        case "PinMatrixRequestType_Current":
        case 1:
            message.type = 1;
            break;
        case "PinMatrixRequestType_NewFirst":
        case 2:
            message.type = 2;
            break;
        case "PinMatrixRequestType_NewSecond":
        case 3:
            message.type = 3;
            break;
        }
        return message;
    };

    /**
     * Creates a plain object from a PinMatrixRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PinMatrixRequest
     * @static
     * @param {PinMatrixRequest} message PinMatrixRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PinMatrixRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.type = options.enums === String ? "PinMatrixRequestType_Current" : 1;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.PinMatrixRequestType[message.type] : message.type;
        return object;
    };

    /**
     * Converts this PinMatrixRequest to JSON.
     * @function toJSON
     * @memberof PinMatrixRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PinMatrixRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PinMatrixRequest;
})();

$root.PinMatrixAck = (function() {

    /**
     * Properties of a PinMatrixAck.
     * @exports IPinMatrixAck
     * @interface IPinMatrixAck
     * @property {string} pin PinMatrixAck pin
     */

    /**
     * Constructs a new PinMatrixAck.
     * @exports PinMatrixAck
     * @classdesc Request: Computer responds with encoded PIN
     * @prev PinMatrixRequest
     * @implements IPinMatrixAck
     * @constructor
     * @param {IPinMatrixAck=} [properties] Properties to set
     */
    function PinMatrixAck(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PinMatrixAck pin.
     * @member {string} pin
     * @memberof PinMatrixAck
     * @instance
     */
    PinMatrixAck.prototype.pin = "";

    /**
     * Creates a new PinMatrixAck instance using the specified properties.
     * @function create
     * @memberof PinMatrixAck
     * @static
     * @param {IPinMatrixAck=} [properties] Properties to set
     * @returns {PinMatrixAck} PinMatrixAck instance
     */
    PinMatrixAck.create = function create(properties) {
        return new PinMatrixAck(properties);
    };

    /**
     * Encodes the specified PinMatrixAck message. Does not implicitly {@link PinMatrixAck.verify|verify} messages.
     * @function encode
     * @memberof PinMatrixAck
     * @static
     * @param {IPinMatrixAck} message PinMatrixAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PinMatrixAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 2 =*/10).string(message.pin);
        return writer;
    };

    /**
     * Encodes the specified PinMatrixAck message, length delimited. Does not implicitly {@link PinMatrixAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PinMatrixAck
     * @static
     * @param {IPinMatrixAck} message PinMatrixAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PinMatrixAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PinMatrixAck message from the specified reader or buffer.
     * @function decode
     * @memberof PinMatrixAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PinMatrixAck} PinMatrixAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PinMatrixAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.PinMatrixAck();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.pin = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("pin"))
            throw $util.ProtocolError("missing required 'pin'", { instance: message });
        return message;
    };

    /**
     * Decodes a PinMatrixAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PinMatrixAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PinMatrixAck} PinMatrixAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PinMatrixAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PinMatrixAck message.
     * @function verify
     * @memberof PinMatrixAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PinMatrixAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isString(message.pin))
            return "pin: string expected";
        return null;
    };

    /**
     * Creates a PinMatrixAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PinMatrixAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PinMatrixAck} PinMatrixAck
     */
    PinMatrixAck.fromObject = function fromObject(object) {
        if (object instanceof $root.PinMatrixAck)
            return object;
        var message = new $root.PinMatrixAck();
        if (object.pin != null)
            message.pin = String(object.pin);
        return message;
    };

    /**
     * Creates a plain object from a PinMatrixAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PinMatrixAck
     * @static
     * @param {PinMatrixAck} message PinMatrixAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PinMatrixAck.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.pin = "";
        if (message.pin != null && message.hasOwnProperty("pin"))
            object.pin = message.pin;
        return object;
    };

    /**
     * Converts this PinMatrixAck to JSON.
     * @function toJSON
     * @memberof PinMatrixAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PinMatrixAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PinMatrixAck;
})();

$root.Cancel = (function() {

    /**
     * Properties of a Cancel.
     * @exports ICancel
     * @interface ICancel
     */

    /**
     * Constructs a new Cancel.
     * @exports Cancel
     * @classdesc Request: Abort last operation that required user interaction
     * @prev ButtonRequest
     * @prev PinMatrixRequest
     * @prev PassphraseRequest
     * @implements ICancel
     * @constructor
     * @param {ICancel=} [properties] Properties to set
     */
    function Cancel(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new Cancel instance using the specified properties.
     * @function create
     * @memberof Cancel
     * @static
     * @param {ICancel=} [properties] Properties to set
     * @returns {Cancel} Cancel instance
     */
    Cancel.create = function create(properties) {
        return new Cancel(properties);
    };

    /**
     * Encodes the specified Cancel message. Does not implicitly {@link Cancel.verify|verify} messages.
     * @function encode
     * @memberof Cancel
     * @static
     * @param {ICancel} message Cancel message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Cancel.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified Cancel message, length delimited. Does not implicitly {@link Cancel.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Cancel
     * @static
     * @param {ICancel} message Cancel message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Cancel.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a Cancel message from the specified reader or buffer.
     * @function decode
     * @memberof Cancel
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Cancel} Cancel
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Cancel.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Cancel();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a Cancel message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Cancel
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Cancel} Cancel
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Cancel.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a Cancel message.
     * @function verify
     * @memberof Cancel
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Cancel.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a Cancel message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Cancel
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Cancel} Cancel
     */
    Cancel.fromObject = function fromObject(object) {
        if (object instanceof $root.Cancel)
            return object;
        return new $root.Cancel();
    };

    /**
     * Creates a plain object from a Cancel message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Cancel
     * @static
     * @param {Cancel} message Cancel
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Cancel.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this Cancel to JSON.
     * @function toJSON
     * @memberof Cancel
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Cancel.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Cancel;
})();

$root.PassphraseRequest = (function() {

    /**
     * Properties of a PassphraseRequest.
     * @exports IPassphraseRequest
     * @interface IPassphraseRequest
     * @property {boolean|null} [onDevice] PassphraseRequest onDevice
     */

    /**
     * Constructs a new PassphraseRequest.
     * @exports PassphraseRequest
     * @classdesc Response: Device awaits encryption passphrase
     * @next PassphraseAck
     * @next Cancel
     * @implements IPassphraseRequest
     * @constructor
     * @param {IPassphraseRequest=} [properties] Properties to set
     */
    function PassphraseRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PassphraseRequest onDevice.
     * @member {boolean} onDevice
     * @memberof PassphraseRequest
     * @instance
     */
    PassphraseRequest.prototype.onDevice = false;

    /**
     * Creates a new PassphraseRequest instance using the specified properties.
     * @function create
     * @memberof PassphraseRequest
     * @static
     * @param {IPassphraseRequest=} [properties] Properties to set
     * @returns {PassphraseRequest} PassphraseRequest instance
     */
    PassphraseRequest.create = function create(properties) {
        return new PassphraseRequest(properties);
    };

    /**
     * Encodes the specified PassphraseRequest message. Does not implicitly {@link PassphraseRequest.verify|verify} messages.
     * @function encode
     * @memberof PassphraseRequest
     * @static
     * @param {IPassphraseRequest} message PassphraseRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.onDevice != null && message.hasOwnProperty("onDevice"))
            writer.uint32(/* id 1, wireType 0 =*/8).bool(message.onDevice);
        return writer;
    };

    /**
     * Encodes the specified PassphraseRequest message, length delimited. Does not implicitly {@link PassphraseRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PassphraseRequest
     * @static
     * @param {IPassphraseRequest} message PassphraseRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PassphraseRequest message from the specified reader or buffer.
     * @function decode
     * @memberof PassphraseRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PassphraseRequest} PassphraseRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.PassphraseRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.onDevice = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PassphraseRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PassphraseRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PassphraseRequest} PassphraseRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PassphraseRequest message.
     * @function verify
     * @memberof PassphraseRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PassphraseRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.onDevice != null && message.hasOwnProperty("onDevice"))
            if (typeof message.onDevice !== "boolean")
                return "onDevice: boolean expected";
        return null;
    };

    /**
     * Creates a PassphraseRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PassphraseRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PassphraseRequest} PassphraseRequest
     */
    PassphraseRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.PassphraseRequest)
            return object;
        var message = new $root.PassphraseRequest();
        if (object.onDevice != null)
            message.onDevice = Boolean(object.onDevice);
        return message;
    };

    /**
     * Creates a plain object from a PassphraseRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PassphraseRequest
     * @static
     * @param {PassphraseRequest} message PassphraseRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PassphraseRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.onDevice = false;
        if (message.onDevice != null && message.hasOwnProperty("onDevice"))
            object.onDevice = message.onDevice;
        return object;
    };

    /**
     * Converts this PassphraseRequest to JSON.
     * @function toJSON
     * @memberof PassphraseRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PassphraseRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PassphraseRequest;
})();

$root.PassphraseAck = (function() {

    /**
     * Properties of a PassphraseAck.
     * @exports IPassphraseAck
     * @interface IPassphraseAck
     * @property {string|null} [passphrase] PassphraseAck passphrase
     * @property {Uint8Array|null} [state] PassphraseAck state
     */

    /**
     * Constructs a new PassphraseAck.
     * @exports PassphraseAck
     * @classdesc Request: Send passphrase back
     * @prev PassphraseRequest
     * @next PassphraseStateRequest
     * @implements IPassphraseAck
     * @constructor
     * @param {IPassphraseAck=} [properties] Properties to set
     */
    function PassphraseAck(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PassphraseAck passphrase.
     * @member {string} passphrase
     * @memberof PassphraseAck
     * @instance
     */
    PassphraseAck.prototype.passphrase = "";

    /**
     * PassphraseAck state.
     * @member {Uint8Array} state
     * @memberof PassphraseAck
     * @instance
     */
    PassphraseAck.prototype.state = $util.newBuffer([]);

    /**
     * Creates a new PassphraseAck instance using the specified properties.
     * @function create
     * @memberof PassphraseAck
     * @static
     * @param {IPassphraseAck=} [properties] Properties to set
     * @returns {PassphraseAck} PassphraseAck instance
     */
    PassphraseAck.create = function create(properties) {
        return new PassphraseAck(properties);
    };

    /**
     * Encodes the specified PassphraseAck message. Does not implicitly {@link PassphraseAck.verify|verify} messages.
     * @function encode
     * @memberof PassphraseAck
     * @static
     * @param {IPassphraseAck} message PassphraseAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.passphrase != null && message.hasOwnProperty("passphrase"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.passphrase);
        if (message.state != null && message.hasOwnProperty("state"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.state);
        return writer;
    };

    /**
     * Encodes the specified PassphraseAck message, length delimited. Does not implicitly {@link PassphraseAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PassphraseAck
     * @static
     * @param {IPassphraseAck} message PassphraseAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PassphraseAck message from the specified reader or buffer.
     * @function decode
     * @memberof PassphraseAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PassphraseAck} PassphraseAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.PassphraseAck();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.passphrase = reader.string();
                break;
            case 2:
                message.state = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PassphraseAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PassphraseAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PassphraseAck} PassphraseAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PassphraseAck message.
     * @function verify
     * @memberof PassphraseAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PassphraseAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.passphrase != null && message.hasOwnProperty("passphrase"))
            if (!$util.isString(message.passphrase))
                return "passphrase: string expected";
        if (message.state != null && message.hasOwnProperty("state"))
            if (!(message.state && typeof message.state.length === "number" || $util.isString(message.state)))
                return "state: buffer expected";
        return null;
    };

    /**
     * Creates a PassphraseAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PassphraseAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PassphraseAck} PassphraseAck
     */
    PassphraseAck.fromObject = function fromObject(object) {
        if (object instanceof $root.PassphraseAck)
            return object;
        var message = new $root.PassphraseAck();
        if (object.passphrase != null)
            message.passphrase = String(object.passphrase);
        if (object.state != null)
            if (typeof object.state === "string")
                $util.base64.decode(object.state, message.state = $util.newBuffer($util.base64.length(object.state)), 0);
            else if (object.state.length)
                message.state = object.state;
        return message;
    };

    /**
     * Creates a plain object from a PassphraseAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PassphraseAck
     * @static
     * @param {PassphraseAck} message PassphraseAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PassphraseAck.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.passphrase = "";
            if (options.bytes === String)
                object.state = "";
            else {
                object.state = [];
                if (options.bytes !== Array)
                    object.state = $util.newBuffer(object.state);
            }
        }
        if (message.passphrase != null && message.hasOwnProperty("passphrase"))
            object.passphrase = message.passphrase;
        if (message.state != null && message.hasOwnProperty("state"))
            object.state = options.bytes === String ? $util.base64.encode(message.state, 0, message.state.length) : options.bytes === Array ? Array.prototype.slice.call(message.state) : message.state;
        return object;
    };

    /**
     * Converts this PassphraseAck to JSON.
     * @function toJSON
     * @memberof PassphraseAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PassphraseAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PassphraseAck;
})();

$root.PassphraseStateRequest = (function() {

    /**
     * Properties of a PassphraseStateRequest.
     * @exports IPassphraseStateRequest
     * @interface IPassphraseStateRequest
     * @property {Uint8Array|null} [state] PassphraseStateRequest state
     */

    /**
     * Constructs a new PassphraseStateRequest.
     * @exports PassphraseStateRequest
     * @classdesc @prev PassphraseAck
     * @next PassphraseStateAck
     * @implements IPassphraseStateRequest
     * @constructor
     * @param {IPassphraseStateRequest=} [properties] Properties to set
     */
    function PassphraseStateRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * PassphraseStateRequest state.
     * @member {Uint8Array} state
     * @memberof PassphraseStateRequest
     * @instance
     */
    PassphraseStateRequest.prototype.state = $util.newBuffer([]);

    /**
     * Creates a new PassphraseStateRequest instance using the specified properties.
     * @function create
     * @memberof PassphraseStateRequest
     * @static
     * @param {IPassphraseStateRequest=} [properties] Properties to set
     * @returns {PassphraseStateRequest} PassphraseStateRequest instance
     */
    PassphraseStateRequest.create = function create(properties) {
        return new PassphraseStateRequest(properties);
    };

    /**
     * Encodes the specified PassphraseStateRequest message. Does not implicitly {@link PassphraseStateRequest.verify|verify} messages.
     * @function encode
     * @memberof PassphraseStateRequest
     * @static
     * @param {IPassphraseStateRequest} message PassphraseStateRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseStateRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.state != null && message.hasOwnProperty("state"))
            writer.uint32(/* id 1, wireType 2 =*/10).bytes(message.state);
        return writer;
    };

    /**
     * Encodes the specified PassphraseStateRequest message, length delimited. Does not implicitly {@link PassphraseStateRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PassphraseStateRequest
     * @static
     * @param {IPassphraseStateRequest} message PassphraseStateRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseStateRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PassphraseStateRequest message from the specified reader or buffer.
     * @function decode
     * @memberof PassphraseStateRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PassphraseStateRequest} PassphraseStateRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseStateRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.PassphraseStateRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.state = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PassphraseStateRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PassphraseStateRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PassphraseStateRequest} PassphraseStateRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseStateRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PassphraseStateRequest message.
     * @function verify
     * @memberof PassphraseStateRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PassphraseStateRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.state != null && message.hasOwnProperty("state"))
            if (!(message.state && typeof message.state.length === "number" || $util.isString(message.state)))
                return "state: buffer expected";
        return null;
    };

    /**
     * Creates a PassphraseStateRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PassphraseStateRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PassphraseStateRequest} PassphraseStateRequest
     */
    PassphraseStateRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.PassphraseStateRequest)
            return object;
        var message = new $root.PassphraseStateRequest();
        if (object.state != null)
            if (typeof object.state === "string")
                $util.base64.decode(object.state, message.state = $util.newBuffer($util.base64.length(object.state)), 0);
            else if (object.state.length)
                message.state = object.state;
        return message;
    };

    /**
     * Creates a plain object from a PassphraseStateRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PassphraseStateRequest
     * @static
     * @param {PassphraseStateRequest} message PassphraseStateRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PassphraseStateRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            if (options.bytes === String)
                object.state = "";
            else {
                object.state = [];
                if (options.bytes !== Array)
                    object.state = $util.newBuffer(object.state);
            }
        if (message.state != null && message.hasOwnProperty("state"))
            object.state = options.bytes === String ? $util.base64.encode(message.state, 0, message.state.length) : options.bytes === Array ? Array.prototype.slice.call(message.state) : message.state;
        return object;
    };

    /**
     * Converts this PassphraseStateRequest to JSON.
     * @function toJSON
     * @memberof PassphraseStateRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PassphraseStateRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PassphraseStateRequest;
})();

$root.PassphraseStateAck = (function() {

    /**
     * Properties of a PassphraseStateAck.
     * @exports IPassphraseStateAck
     * @interface IPassphraseStateAck
     */

    /**
     * Constructs a new PassphraseStateAck.
     * @exports PassphraseStateAck
     * @classdesc @prev PassphraseStateRequest
     * @implements IPassphraseStateAck
     * @constructor
     * @param {IPassphraseStateAck=} [properties] Properties to set
     */
    function PassphraseStateAck(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new PassphraseStateAck instance using the specified properties.
     * @function create
     * @memberof PassphraseStateAck
     * @static
     * @param {IPassphraseStateAck=} [properties] Properties to set
     * @returns {PassphraseStateAck} PassphraseStateAck instance
     */
    PassphraseStateAck.create = function create(properties) {
        return new PassphraseStateAck(properties);
    };

    /**
     * Encodes the specified PassphraseStateAck message. Does not implicitly {@link PassphraseStateAck.verify|verify} messages.
     * @function encode
     * @memberof PassphraseStateAck
     * @static
     * @param {IPassphraseStateAck} message PassphraseStateAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseStateAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified PassphraseStateAck message, length delimited. Does not implicitly {@link PassphraseStateAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof PassphraseStateAck
     * @static
     * @param {IPassphraseStateAck} message PassphraseStateAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    PassphraseStateAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a PassphraseStateAck message from the specified reader or buffer.
     * @function decode
     * @memberof PassphraseStateAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {PassphraseStateAck} PassphraseStateAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseStateAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.PassphraseStateAck();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a PassphraseStateAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof PassphraseStateAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {PassphraseStateAck} PassphraseStateAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    PassphraseStateAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a PassphraseStateAck message.
     * @function verify
     * @memberof PassphraseStateAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    PassphraseStateAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a PassphraseStateAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof PassphraseStateAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {PassphraseStateAck} PassphraseStateAck
     */
    PassphraseStateAck.fromObject = function fromObject(object) {
        if (object instanceof $root.PassphraseStateAck)
            return object;
        return new $root.PassphraseStateAck();
    };

    /**
     * Creates a plain object from a PassphraseStateAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof PassphraseStateAck
     * @static
     * @param {PassphraseStateAck} message PassphraseStateAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    PassphraseStateAck.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this PassphraseStateAck to JSON.
     * @function toJSON
     * @memberof PassphraseStateAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    PassphraseStateAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return PassphraseStateAck;
})();

$root.GetEntropy = (function() {

    /**
     * Properties of a GetEntropy.
     * @exports IGetEntropy
     * @interface IGetEntropy
     * @property {number} size GetEntropy size
     */

    /**
     * Constructs a new GetEntropy.
     * @exports GetEntropy
     * @classdesc Request: Request a sample of random data generated by hardware RNG. May be used for testing.
     * @next ButtonRequest
     * @next Entropy
     * @next Failure
     * @implements IGetEntropy
     * @constructor
     * @param {IGetEntropy=} [properties] Properties to set
     */
    function GetEntropy(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * GetEntropy size.
     * @member {number} size
     * @memberof GetEntropy
     * @instance
     */
    GetEntropy.prototype.size = 0;

    /**
     * Creates a new GetEntropy instance using the specified properties.
     * @function create
     * @memberof GetEntropy
     * @static
     * @param {IGetEntropy=} [properties] Properties to set
     * @returns {GetEntropy} GetEntropy instance
     */
    GetEntropy.create = function create(properties) {
        return new GetEntropy(properties);
    };

    /**
     * Encodes the specified GetEntropy message. Does not implicitly {@link GetEntropy.verify|verify} messages.
     * @function encode
     * @memberof GetEntropy
     * @static
     * @param {IGetEntropy} message GetEntropy message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GetEntropy.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.size);
        return writer;
    };

    /**
     * Encodes the specified GetEntropy message, length delimited. Does not implicitly {@link GetEntropy.verify|verify} messages.
     * @function encodeDelimited
     * @memberof GetEntropy
     * @static
     * @param {IGetEntropy} message GetEntropy message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    GetEntropy.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a GetEntropy message from the specified reader or buffer.
     * @function decode
     * @memberof GetEntropy
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {GetEntropy} GetEntropy
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GetEntropy.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.GetEntropy();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.size = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("size"))
            throw $util.ProtocolError("missing required 'size'", { instance: message });
        return message;
    };

    /**
     * Decodes a GetEntropy message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof GetEntropy
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {GetEntropy} GetEntropy
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    GetEntropy.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a GetEntropy message.
     * @function verify
     * @memberof GetEntropy
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    GetEntropy.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isInteger(message.size))
            return "size: integer expected";
        return null;
    };

    /**
     * Creates a GetEntropy message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof GetEntropy
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {GetEntropy} GetEntropy
     */
    GetEntropy.fromObject = function fromObject(object) {
        if (object instanceof $root.GetEntropy)
            return object;
        var message = new $root.GetEntropy();
        if (object.size != null)
            message.size = object.size >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a GetEntropy message. Also converts values to other types if specified.
     * @function toObject
     * @memberof GetEntropy
     * @static
     * @param {GetEntropy} message GetEntropy
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    GetEntropy.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.size = 0;
        if (message.size != null && message.hasOwnProperty("size"))
            object.size = message.size;
        return object;
    };

    /**
     * Converts this GetEntropy to JSON.
     * @function toJSON
     * @memberof GetEntropy
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    GetEntropy.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return GetEntropy;
})();

$root.Entropy = (function() {

    /**
     * Properties of an Entropy.
     * @exports IEntropy
     * @interface IEntropy
     * @property {Uint8Array} entropy Entropy entropy
     */

    /**
     * Constructs a new Entropy.
     * @exports Entropy
     * @classdesc Response: Reply with random data generated by internal RNG
     * @prev GetEntropy
     * @implements IEntropy
     * @constructor
     * @param {IEntropy=} [properties] Properties to set
     */
    function Entropy(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Entropy entropy.
     * @member {Uint8Array} entropy
     * @memberof Entropy
     * @instance
     */
    Entropy.prototype.entropy = $util.newBuffer([]);

    /**
     * Creates a new Entropy instance using the specified properties.
     * @function create
     * @memberof Entropy
     * @static
     * @param {IEntropy=} [properties] Properties to set
     * @returns {Entropy} Entropy instance
     */
    Entropy.create = function create(properties) {
        return new Entropy(properties);
    };

    /**
     * Encodes the specified Entropy message. Does not implicitly {@link Entropy.verify|verify} messages.
     * @function encode
     * @memberof Entropy
     * @static
     * @param {IEntropy} message Entropy message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Entropy.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 2 =*/10).bytes(message.entropy);
        return writer;
    };

    /**
     * Encodes the specified Entropy message, length delimited. Does not implicitly {@link Entropy.verify|verify} messages.
     * @function encodeDelimited
     * @memberof Entropy
     * @static
     * @param {IEntropy} message Entropy message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    Entropy.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an Entropy message from the specified reader or buffer.
     * @function decode
     * @memberof Entropy
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {Entropy} Entropy
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Entropy.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.Entropy();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.entropy = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("entropy"))
            throw $util.ProtocolError("missing required 'entropy'", { instance: message });
        return message;
    };

    /**
     * Decodes an Entropy message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof Entropy
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {Entropy} Entropy
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    Entropy.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an Entropy message.
     * @function verify
     * @memberof Entropy
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    Entropy.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!(message.entropy && typeof message.entropy.length === "number" || $util.isString(message.entropy)))
            return "entropy: buffer expected";
        return null;
    };

    /**
     * Creates an Entropy message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof Entropy
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {Entropy} Entropy
     */
    Entropy.fromObject = function fromObject(object) {
        if (object instanceof $root.Entropy)
            return object;
        var message = new $root.Entropy();
        if (object.entropy != null)
            if (typeof object.entropy === "string")
                $util.base64.decode(object.entropy, message.entropy = $util.newBuffer($util.base64.length(object.entropy)), 0);
            else if (object.entropy.length)
                message.entropy = object.entropy;
        return message;
    };

    /**
     * Creates a plain object from an Entropy message. Also converts values to other types if specified.
     * @function toObject
     * @memberof Entropy
     * @static
     * @param {Entropy} message Entropy
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    Entropy.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            if (options.bytes === String)
                object.entropy = "";
            else {
                object.entropy = [];
                if (options.bytes !== Array)
                    object.entropy = $util.newBuffer(object.entropy);
            }
        if (message.entropy != null && message.hasOwnProperty("entropy"))
            object.entropy = options.bytes === String ? $util.base64.encode(message.entropy, 0, message.entropy.length) : options.bytes === Array ? Array.prototype.slice.call(message.entropy) : message.entropy;
        return object;
    };

    /**
     * Converts this Entropy to JSON.
     * @function toJSON
     * @memberof Entropy
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    Entropy.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return Entropy;
})();

$root.WipeDevice = (function() {

    /**
     * Properties of a WipeDevice.
     * @exports IWipeDevice
     * @interface IWipeDevice
     */

    /**
     * Constructs a new WipeDevice.
     * @exports WipeDevice
     * @classdesc Request: Request device to wipe all sensitive data and settings
     * @next ButtonRequest
     * @implements IWipeDevice
     * @constructor
     * @param {IWipeDevice=} [properties] Properties to set
     */
    function WipeDevice(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new WipeDevice instance using the specified properties.
     * @function create
     * @memberof WipeDevice
     * @static
     * @param {IWipeDevice=} [properties] Properties to set
     * @returns {WipeDevice} WipeDevice instance
     */
    WipeDevice.create = function create(properties) {
        return new WipeDevice(properties);
    };

    /**
     * Encodes the specified WipeDevice message. Does not implicitly {@link WipeDevice.verify|verify} messages.
     * @function encode
     * @memberof WipeDevice
     * @static
     * @param {IWipeDevice} message WipeDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    WipeDevice.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified WipeDevice message, length delimited. Does not implicitly {@link WipeDevice.verify|verify} messages.
     * @function encodeDelimited
     * @memberof WipeDevice
     * @static
     * @param {IWipeDevice} message WipeDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    WipeDevice.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a WipeDevice message from the specified reader or buffer.
     * @function decode
     * @memberof WipeDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {WipeDevice} WipeDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    WipeDevice.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.WipeDevice();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a WipeDevice message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof WipeDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {WipeDevice} WipeDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    WipeDevice.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a WipeDevice message.
     * @function verify
     * @memberof WipeDevice
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    WipeDevice.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a WipeDevice message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof WipeDevice
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {WipeDevice} WipeDevice
     */
    WipeDevice.fromObject = function fromObject(object) {
        if (object instanceof $root.WipeDevice)
            return object;
        return new $root.WipeDevice();
    };

    /**
     * Creates a plain object from a WipeDevice message. Also converts values to other types if specified.
     * @function toObject
     * @memberof WipeDevice
     * @static
     * @param {WipeDevice} message WipeDevice
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    WipeDevice.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this WipeDevice to JSON.
     * @function toJSON
     * @memberof WipeDevice
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    WipeDevice.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return WipeDevice;
})();

$root.LoadDevice = (function() {

    /**
     * Properties of a LoadDevice.
     * @exports ILoadDevice
     * @interface ILoadDevice
     * @property {string|null} [mnemonic] LoadDevice mnemonic
     * @property {IHDNodeType|null} [node] LoadDevice node
     * @property {string|null} [pin] LoadDevice pin
     * @property {boolean|null} [passphraseProtection] LoadDevice passphraseProtection
     * @property {string|null} [language] LoadDevice language
     * @property {string|null} [label] LoadDevice label
     * @property {boolean|null} [skipChecksum] LoadDevice skipChecksum
     * @property {number|null} [u2fCounter] LoadDevice u2fCounter
     */

    /**
     * Constructs a new LoadDevice.
     * @exports LoadDevice
     * @classdesc Request: Load seed and related internal settings from the computer
     * @next ButtonRequest
     * @next Success
     * @next Failure
     * @implements ILoadDevice
     * @constructor
     * @param {ILoadDevice=} [properties] Properties to set
     */
    function LoadDevice(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * LoadDevice mnemonic.
     * @member {string} mnemonic
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.mnemonic = "";

    /**
     * LoadDevice node.
     * @member {IHDNodeType|null|undefined} node
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.node = null;

    /**
     * LoadDevice pin.
     * @member {string} pin
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.pin = "";

    /**
     * LoadDevice passphraseProtection.
     * @member {boolean} passphraseProtection
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.passphraseProtection = false;

    /**
     * LoadDevice language.
     * @member {string} language
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.language = "english";

    /**
     * LoadDevice label.
     * @member {string} label
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.label = "";

    /**
     * LoadDevice skipChecksum.
     * @member {boolean} skipChecksum
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.skipChecksum = false;

    /**
     * LoadDevice u2fCounter.
     * @member {number} u2fCounter
     * @memberof LoadDevice
     * @instance
     */
    LoadDevice.prototype.u2fCounter = 0;

    /**
     * Creates a new LoadDevice instance using the specified properties.
     * @function create
     * @memberof LoadDevice
     * @static
     * @param {ILoadDevice=} [properties] Properties to set
     * @returns {LoadDevice} LoadDevice instance
     */
    LoadDevice.create = function create(properties) {
        return new LoadDevice(properties);
    };

    /**
     * Encodes the specified LoadDevice message. Does not implicitly {@link LoadDevice.verify|verify} messages.
     * @function encode
     * @memberof LoadDevice
     * @static
     * @param {ILoadDevice} message LoadDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    LoadDevice.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.mnemonic != null && message.hasOwnProperty("mnemonic"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.mnemonic);
        if (message.node != null && message.hasOwnProperty("node"))
            $root.HDNodeType.encode(message.node, writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.pin != null && message.hasOwnProperty("pin"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.pin);
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            writer.uint32(/* id 4, wireType 0 =*/32).bool(message.passphraseProtection);
        if (message.language != null && message.hasOwnProperty("language"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.language);
        if (message.label != null && message.hasOwnProperty("label"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.label);
        if (message.skipChecksum != null && message.hasOwnProperty("skipChecksum"))
            writer.uint32(/* id 7, wireType 0 =*/56).bool(message.skipChecksum);
        if (message.u2fCounter != null && message.hasOwnProperty("u2fCounter"))
            writer.uint32(/* id 8, wireType 0 =*/64).uint32(message.u2fCounter);
        return writer;
    };

    /**
     * Encodes the specified LoadDevice message, length delimited. Does not implicitly {@link LoadDevice.verify|verify} messages.
     * @function encodeDelimited
     * @memberof LoadDevice
     * @static
     * @param {ILoadDevice} message LoadDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    LoadDevice.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a LoadDevice message from the specified reader or buffer.
     * @function decode
     * @memberof LoadDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {LoadDevice} LoadDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    LoadDevice.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.LoadDevice();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.mnemonic = reader.string();
                break;
            case 2:
                message.node = $root.HDNodeType.decode(reader, reader.uint32());
                break;
            case 3:
                message.pin = reader.string();
                break;
            case 4:
                message.passphraseProtection = reader.bool();
                break;
            case 5:
                message.language = reader.string();
                break;
            case 6:
                message.label = reader.string();
                break;
            case 7:
                message.skipChecksum = reader.bool();
                break;
            case 8:
                message.u2fCounter = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a LoadDevice message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof LoadDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {LoadDevice} LoadDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    LoadDevice.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a LoadDevice message.
     * @function verify
     * @memberof LoadDevice
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    LoadDevice.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.mnemonic != null && message.hasOwnProperty("mnemonic"))
            if (!$util.isString(message.mnemonic))
                return "mnemonic: string expected";
        if (message.node != null && message.hasOwnProperty("node")) {
            var error = $root.HDNodeType.verify(message.node);
            if (error)
                return "node." + error;
        }
        if (message.pin != null && message.hasOwnProperty("pin"))
            if (!$util.isString(message.pin))
                return "pin: string expected";
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            if (typeof message.passphraseProtection !== "boolean")
                return "passphraseProtection: boolean expected";
        if (message.language != null && message.hasOwnProperty("language"))
            if (!$util.isString(message.language))
                return "language: string expected";
        if (message.label != null && message.hasOwnProperty("label"))
            if (!$util.isString(message.label))
                return "label: string expected";
        if (message.skipChecksum != null && message.hasOwnProperty("skipChecksum"))
            if (typeof message.skipChecksum !== "boolean")
                return "skipChecksum: boolean expected";
        if (message.u2fCounter != null && message.hasOwnProperty("u2fCounter"))
            if (!$util.isInteger(message.u2fCounter))
                return "u2fCounter: integer expected";
        return null;
    };

    /**
     * Creates a LoadDevice message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof LoadDevice
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {LoadDevice} LoadDevice
     */
    LoadDevice.fromObject = function fromObject(object) {
        if (object instanceof $root.LoadDevice)
            return object;
        var message = new $root.LoadDevice();
        if (object.mnemonic != null)
            message.mnemonic = String(object.mnemonic);
        if (object.node != null) {
            if (typeof object.node !== "object")
                throw TypeError(".LoadDevice.node: object expected");
            message.node = $root.HDNodeType.fromObject(object.node);
        }
        if (object.pin != null)
            message.pin = String(object.pin);
        if (object.passphraseProtection != null)
            message.passphraseProtection = Boolean(object.passphraseProtection);
        if (object.language != null)
            message.language = String(object.language);
        if (object.label != null)
            message.label = String(object.label);
        if (object.skipChecksum != null)
            message.skipChecksum = Boolean(object.skipChecksum);
        if (object.u2fCounter != null)
            message.u2fCounter = object.u2fCounter >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a LoadDevice message. Also converts values to other types if specified.
     * @function toObject
     * @memberof LoadDevice
     * @static
     * @param {LoadDevice} message LoadDevice
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    LoadDevice.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.mnemonic = "";
            object.node = null;
            object.pin = "";
            object.passphraseProtection = false;
            object.language = "english";
            object.label = "";
            object.skipChecksum = false;
            object.u2fCounter = 0;
        }
        if (message.mnemonic != null && message.hasOwnProperty("mnemonic"))
            object.mnemonic = message.mnemonic;
        if (message.node != null && message.hasOwnProperty("node"))
            object.node = $root.HDNodeType.toObject(message.node, options);
        if (message.pin != null && message.hasOwnProperty("pin"))
            object.pin = message.pin;
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            object.passphraseProtection = message.passphraseProtection;
        if (message.language != null && message.hasOwnProperty("language"))
            object.language = message.language;
        if (message.label != null && message.hasOwnProperty("label"))
            object.label = message.label;
        if (message.skipChecksum != null && message.hasOwnProperty("skipChecksum"))
            object.skipChecksum = message.skipChecksum;
        if (message.u2fCounter != null && message.hasOwnProperty("u2fCounter"))
            object.u2fCounter = message.u2fCounter;
        return object;
    };

    /**
     * Converts this LoadDevice to JSON.
     * @function toJSON
     * @memberof LoadDevice
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    LoadDevice.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return LoadDevice;
})();

$root.ResetDevice = (function() {

    /**
     * Properties of a ResetDevice.
     * @exports IResetDevice
     * @interface IResetDevice
     * @property {boolean|null} [displayRandom] ResetDevice displayRandom
     * @property {number|null} [strength] ResetDevice strength
     * @property {boolean|null} [passphraseProtection] ResetDevice passphraseProtection
     * @property {boolean|null} [pinProtection] ResetDevice pinProtection
     * @property {string|null} [language] ResetDevice language
     * @property {string|null} [label] ResetDevice label
     * @property {number|null} [u2fCounter] ResetDevice u2fCounter
     * @property {boolean|null} [skipBackup] ResetDevice skipBackup
     */

    /**
     * Constructs a new ResetDevice.
     * @exports ResetDevice
     * @classdesc Request: Ask device to do initialization involving user interaction
     * @next EntropyRequest
     * @next Failure
     * @implements IResetDevice
     * @constructor
     * @param {IResetDevice=} [properties] Properties to set
     */
    function ResetDevice(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * ResetDevice displayRandom.
     * @member {boolean} displayRandom
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.displayRandom = false;

    /**
     * ResetDevice strength.
     * @member {number} strength
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.strength = 256;

    /**
     * ResetDevice passphraseProtection.
     * @member {boolean} passphraseProtection
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.passphraseProtection = false;

    /**
     * ResetDevice pinProtection.
     * @member {boolean} pinProtection
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.pinProtection = false;

    /**
     * ResetDevice language.
     * @member {string} language
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.language = "english";

    /**
     * ResetDevice label.
     * @member {string} label
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.label = "";

    /**
     * ResetDevice u2fCounter.
     * @member {number} u2fCounter
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.u2fCounter = 0;

    /**
     * ResetDevice skipBackup.
     * @member {boolean} skipBackup
     * @memberof ResetDevice
     * @instance
     */
    ResetDevice.prototype.skipBackup = false;

    /**
     * Creates a new ResetDevice instance using the specified properties.
     * @function create
     * @memberof ResetDevice
     * @static
     * @param {IResetDevice=} [properties] Properties to set
     * @returns {ResetDevice} ResetDevice instance
     */
    ResetDevice.create = function create(properties) {
        return new ResetDevice(properties);
    };

    /**
     * Encodes the specified ResetDevice message. Does not implicitly {@link ResetDevice.verify|verify} messages.
     * @function encode
     * @memberof ResetDevice
     * @static
     * @param {IResetDevice} message ResetDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ResetDevice.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.displayRandom != null && message.hasOwnProperty("displayRandom"))
            writer.uint32(/* id 1, wireType 0 =*/8).bool(message.displayRandom);
        if (message.strength != null && message.hasOwnProperty("strength"))
            writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.strength);
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.passphraseProtection);
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            writer.uint32(/* id 4, wireType 0 =*/32).bool(message.pinProtection);
        if (message.language != null && message.hasOwnProperty("language"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.language);
        if (message.label != null && message.hasOwnProperty("label"))
            writer.uint32(/* id 6, wireType 2 =*/50).string(message.label);
        if (message.u2fCounter != null && message.hasOwnProperty("u2fCounter"))
            writer.uint32(/* id 7, wireType 0 =*/56).uint32(message.u2fCounter);
        if (message.skipBackup != null && message.hasOwnProperty("skipBackup"))
            writer.uint32(/* id 8, wireType 0 =*/64).bool(message.skipBackup);
        return writer;
    };

    /**
     * Encodes the specified ResetDevice message, length delimited. Does not implicitly {@link ResetDevice.verify|verify} messages.
     * @function encodeDelimited
     * @memberof ResetDevice
     * @static
     * @param {IResetDevice} message ResetDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    ResetDevice.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a ResetDevice message from the specified reader or buffer.
     * @function decode
     * @memberof ResetDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {ResetDevice} ResetDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ResetDevice.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.ResetDevice();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.displayRandom = reader.bool();
                break;
            case 2:
                message.strength = reader.uint32();
                break;
            case 3:
                message.passphraseProtection = reader.bool();
                break;
            case 4:
                message.pinProtection = reader.bool();
                break;
            case 5:
                message.language = reader.string();
                break;
            case 6:
                message.label = reader.string();
                break;
            case 7:
                message.u2fCounter = reader.uint32();
                break;
            case 8:
                message.skipBackup = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a ResetDevice message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof ResetDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {ResetDevice} ResetDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    ResetDevice.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a ResetDevice message.
     * @function verify
     * @memberof ResetDevice
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    ResetDevice.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.displayRandom != null && message.hasOwnProperty("displayRandom"))
            if (typeof message.displayRandom !== "boolean")
                return "displayRandom: boolean expected";
        if (message.strength != null && message.hasOwnProperty("strength"))
            if (!$util.isInteger(message.strength))
                return "strength: integer expected";
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            if (typeof message.passphraseProtection !== "boolean")
                return "passphraseProtection: boolean expected";
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            if (typeof message.pinProtection !== "boolean")
                return "pinProtection: boolean expected";
        if (message.language != null && message.hasOwnProperty("language"))
            if (!$util.isString(message.language))
                return "language: string expected";
        if (message.label != null && message.hasOwnProperty("label"))
            if (!$util.isString(message.label))
                return "label: string expected";
        if (message.u2fCounter != null && message.hasOwnProperty("u2fCounter"))
            if (!$util.isInteger(message.u2fCounter))
                return "u2fCounter: integer expected";
        if (message.skipBackup != null && message.hasOwnProperty("skipBackup"))
            if (typeof message.skipBackup !== "boolean")
                return "skipBackup: boolean expected";
        return null;
    };

    /**
     * Creates a ResetDevice message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof ResetDevice
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {ResetDevice} ResetDevice
     */
    ResetDevice.fromObject = function fromObject(object) {
        if (object instanceof $root.ResetDevice)
            return object;
        var message = new $root.ResetDevice();
        if (object.displayRandom != null)
            message.displayRandom = Boolean(object.displayRandom);
        if (object.strength != null)
            message.strength = object.strength >>> 0;
        if (object.passphraseProtection != null)
            message.passphraseProtection = Boolean(object.passphraseProtection);
        if (object.pinProtection != null)
            message.pinProtection = Boolean(object.pinProtection);
        if (object.language != null)
            message.language = String(object.language);
        if (object.label != null)
            message.label = String(object.label);
        if (object.u2fCounter != null)
            message.u2fCounter = object.u2fCounter >>> 0;
        if (object.skipBackup != null)
            message.skipBackup = Boolean(object.skipBackup);
        return message;
    };

    /**
     * Creates a plain object from a ResetDevice message. Also converts values to other types if specified.
     * @function toObject
     * @memberof ResetDevice
     * @static
     * @param {ResetDevice} message ResetDevice
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    ResetDevice.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.displayRandom = false;
            object.strength = 256;
            object.passphraseProtection = false;
            object.pinProtection = false;
            object.language = "english";
            object.label = "";
            object.u2fCounter = 0;
            object.skipBackup = false;
        }
        if (message.displayRandom != null && message.hasOwnProperty("displayRandom"))
            object.displayRandom = message.displayRandom;
        if (message.strength != null && message.hasOwnProperty("strength"))
            object.strength = message.strength;
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            object.passphraseProtection = message.passphraseProtection;
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            object.pinProtection = message.pinProtection;
        if (message.language != null && message.hasOwnProperty("language"))
            object.language = message.language;
        if (message.label != null && message.hasOwnProperty("label"))
            object.label = message.label;
        if (message.u2fCounter != null && message.hasOwnProperty("u2fCounter"))
            object.u2fCounter = message.u2fCounter;
        if (message.skipBackup != null && message.hasOwnProperty("skipBackup"))
            object.skipBackup = message.skipBackup;
        return object;
    };

    /**
     * Converts this ResetDevice to JSON.
     * @function toJSON
     * @memberof ResetDevice
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    ResetDevice.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return ResetDevice;
})();

$root.BackupDevice = (function() {

    /**
     * Properties of a BackupDevice.
     * @exports IBackupDevice
     * @interface IBackupDevice
     */

    /**
     * Constructs a new BackupDevice.
     * @exports BackupDevice
     * @classdesc Request: Perform backup of the device seed if not backed up using ResetDevice
     * @next ButtonRequest
     * @implements IBackupDevice
     * @constructor
     * @param {IBackupDevice=} [properties] Properties to set
     */
    function BackupDevice(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new BackupDevice instance using the specified properties.
     * @function create
     * @memberof BackupDevice
     * @static
     * @param {IBackupDevice=} [properties] Properties to set
     * @returns {BackupDevice} BackupDevice instance
     */
    BackupDevice.create = function create(properties) {
        return new BackupDevice(properties);
    };

    /**
     * Encodes the specified BackupDevice message. Does not implicitly {@link BackupDevice.verify|verify} messages.
     * @function encode
     * @memberof BackupDevice
     * @static
     * @param {IBackupDevice} message BackupDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    BackupDevice.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified BackupDevice message, length delimited. Does not implicitly {@link BackupDevice.verify|verify} messages.
     * @function encodeDelimited
     * @memberof BackupDevice
     * @static
     * @param {IBackupDevice} message BackupDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    BackupDevice.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a BackupDevice message from the specified reader or buffer.
     * @function decode
     * @memberof BackupDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {BackupDevice} BackupDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    BackupDevice.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.BackupDevice();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a BackupDevice message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof BackupDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {BackupDevice} BackupDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    BackupDevice.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a BackupDevice message.
     * @function verify
     * @memberof BackupDevice
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    BackupDevice.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates a BackupDevice message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof BackupDevice
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {BackupDevice} BackupDevice
     */
    BackupDevice.fromObject = function fromObject(object) {
        if (object instanceof $root.BackupDevice)
            return object;
        return new $root.BackupDevice();
    };

    /**
     * Creates a plain object from a BackupDevice message. Also converts values to other types if specified.
     * @function toObject
     * @memberof BackupDevice
     * @static
     * @param {BackupDevice} message BackupDevice
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    BackupDevice.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this BackupDevice to JSON.
     * @function toJSON
     * @memberof BackupDevice
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    BackupDevice.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return BackupDevice;
})();

$root.EntropyRequest = (function() {

    /**
     * Properties of an EntropyRequest.
     * @exports IEntropyRequest
     * @interface IEntropyRequest
     */

    /**
     * Constructs a new EntropyRequest.
     * @exports EntropyRequest
     * @classdesc Response: Ask for additional entropy from host computer
     * @prev ResetDevice
     * @next EntropyAck
     * @implements IEntropyRequest
     * @constructor
     * @param {IEntropyRequest=} [properties] Properties to set
     */
    function EntropyRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * Creates a new EntropyRequest instance using the specified properties.
     * @function create
     * @memberof EntropyRequest
     * @static
     * @param {IEntropyRequest=} [properties] Properties to set
     * @returns {EntropyRequest} EntropyRequest instance
     */
    EntropyRequest.create = function create(properties) {
        return new EntropyRequest(properties);
    };

    /**
     * Encodes the specified EntropyRequest message. Does not implicitly {@link EntropyRequest.verify|verify} messages.
     * @function encode
     * @memberof EntropyRequest
     * @static
     * @param {IEntropyRequest} message EntropyRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    EntropyRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        return writer;
    };

    /**
     * Encodes the specified EntropyRequest message, length delimited. Does not implicitly {@link EntropyRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof EntropyRequest
     * @static
     * @param {IEntropyRequest} message EntropyRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    EntropyRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an EntropyRequest message from the specified reader or buffer.
     * @function decode
     * @memberof EntropyRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {EntropyRequest} EntropyRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    EntropyRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.EntropyRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an EntropyRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof EntropyRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {EntropyRequest} EntropyRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    EntropyRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an EntropyRequest message.
     * @function verify
     * @memberof EntropyRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    EntropyRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        return null;
    };

    /**
     * Creates an EntropyRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof EntropyRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {EntropyRequest} EntropyRequest
     */
    EntropyRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.EntropyRequest)
            return object;
        return new $root.EntropyRequest();
    };

    /**
     * Creates a plain object from an EntropyRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof EntropyRequest
     * @static
     * @param {EntropyRequest} message EntropyRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    EntropyRequest.toObject = function toObject() {
        return {};
    };

    /**
     * Converts this EntropyRequest to JSON.
     * @function toJSON
     * @memberof EntropyRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    EntropyRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return EntropyRequest;
})();

$root.EntropyAck = (function() {

    /**
     * Properties of an EntropyAck.
     * @exports IEntropyAck
     * @interface IEntropyAck
     * @property {Uint8Array|null} [entropy] EntropyAck entropy
     */

    /**
     * Constructs a new EntropyAck.
     * @exports EntropyAck
     * @classdesc Request: Provide additional entropy for seed generation function
     * @prev EntropyRequest
     * @next ButtonRequest
     * @implements IEntropyAck
     * @constructor
     * @param {IEntropyAck=} [properties] Properties to set
     */
    function EntropyAck(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * EntropyAck entropy.
     * @member {Uint8Array} entropy
     * @memberof EntropyAck
     * @instance
     */
    EntropyAck.prototype.entropy = $util.newBuffer([]);

    /**
     * Creates a new EntropyAck instance using the specified properties.
     * @function create
     * @memberof EntropyAck
     * @static
     * @param {IEntropyAck=} [properties] Properties to set
     * @returns {EntropyAck} EntropyAck instance
     */
    EntropyAck.create = function create(properties) {
        return new EntropyAck(properties);
    };

    /**
     * Encodes the specified EntropyAck message. Does not implicitly {@link EntropyAck.verify|verify} messages.
     * @function encode
     * @memberof EntropyAck
     * @static
     * @param {IEntropyAck} message EntropyAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    EntropyAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.entropy != null && message.hasOwnProperty("entropy"))
            writer.uint32(/* id 1, wireType 2 =*/10).bytes(message.entropy);
        return writer;
    };

    /**
     * Encodes the specified EntropyAck message, length delimited. Does not implicitly {@link EntropyAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof EntropyAck
     * @static
     * @param {IEntropyAck} message EntropyAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    EntropyAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an EntropyAck message from the specified reader or buffer.
     * @function decode
     * @memberof EntropyAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {EntropyAck} EntropyAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    EntropyAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.EntropyAck();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.entropy = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an EntropyAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof EntropyAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {EntropyAck} EntropyAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    EntropyAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an EntropyAck message.
     * @function verify
     * @memberof EntropyAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    EntropyAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.entropy != null && message.hasOwnProperty("entropy"))
            if (!(message.entropy && typeof message.entropy.length === "number" || $util.isString(message.entropy)))
                return "entropy: buffer expected";
        return null;
    };

    /**
     * Creates an EntropyAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof EntropyAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {EntropyAck} EntropyAck
     */
    EntropyAck.fromObject = function fromObject(object) {
        if (object instanceof $root.EntropyAck)
            return object;
        var message = new $root.EntropyAck();
        if (object.entropy != null)
            if (typeof object.entropy === "string")
                $util.base64.decode(object.entropy, message.entropy = $util.newBuffer($util.base64.length(object.entropy)), 0);
            else if (object.entropy.length)
                message.entropy = object.entropy;
        return message;
    };

    /**
     * Creates a plain object from an EntropyAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof EntropyAck
     * @static
     * @param {EntropyAck} message EntropyAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    EntropyAck.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            if (options.bytes === String)
                object.entropy = "";
            else {
                object.entropy = [];
                if (options.bytes !== Array)
                    object.entropy = $util.newBuffer(object.entropy);
            }
        if (message.entropy != null && message.hasOwnProperty("entropy"))
            object.entropy = options.bytes === String ? $util.base64.encode(message.entropy, 0, message.entropy.length) : options.bytes === Array ? Array.prototype.slice.call(message.entropy) : message.entropy;
        return object;
    };

    /**
     * Converts this EntropyAck to JSON.
     * @function toJSON
     * @memberof EntropyAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    EntropyAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return EntropyAck;
})();

$root.RecoveryDevice = (function() {

    /**
     * Properties of a RecoveryDevice.
     * @exports IRecoveryDevice
     * @interface IRecoveryDevice
     * @property {number|null} [wordCount] RecoveryDevice wordCount
     * @property {boolean|null} [passphraseProtection] RecoveryDevice passphraseProtection
     * @property {boolean|null} [pinProtection] RecoveryDevice pinProtection
     * @property {string|null} [language] RecoveryDevice language
     * @property {string|null} [label] RecoveryDevice label
     * @property {boolean|null} [enforceWordlist] RecoveryDevice enforceWordlist
     * @property {number|null} [type] RecoveryDevice type
     * @property {boolean|null} [dryRun] RecoveryDevice dryRun
     */

    /**
     * Constructs a new RecoveryDevice.
     * @exports RecoveryDevice
     * @classdesc Request: Start recovery workflow asking user for specific words of mnemonic
     * Used to recovery device safely even on untrusted computer.
     * @next WordRequest
     * @implements IRecoveryDevice
     * @constructor
     * @param {IRecoveryDevice=} [properties] Properties to set
     */
    function RecoveryDevice(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * RecoveryDevice wordCount.
     * @member {number} wordCount
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.wordCount = 0;

    /**
     * RecoveryDevice passphraseProtection.
     * @member {boolean} passphraseProtection
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.passphraseProtection = false;

    /**
     * RecoveryDevice pinProtection.
     * @member {boolean} pinProtection
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.pinProtection = false;

    /**
     * RecoveryDevice language.
     * @member {string} language
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.language = "english";

    /**
     * RecoveryDevice label.
     * @member {string} label
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.label = "";

    /**
     * RecoveryDevice enforceWordlist.
     * @member {boolean} enforceWordlist
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.enforceWordlist = false;

    /**
     * RecoveryDevice type.
     * @member {number} type
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.type = 0;

    /**
     * RecoveryDevice dryRun.
     * @member {boolean} dryRun
     * @memberof RecoveryDevice
     * @instance
     */
    RecoveryDevice.prototype.dryRun = false;

    /**
     * Creates a new RecoveryDevice instance using the specified properties.
     * @function create
     * @memberof RecoveryDevice
     * @static
     * @param {IRecoveryDevice=} [properties] Properties to set
     * @returns {RecoveryDevice} RecoveryDevice instance
     */
    RecoveryDevice.create = function create(properties) {
        return new RecoveryDevice(properties);
    };

    /**
     * Encodes the specified RecoveryDevice message. Does not implicitly {@link RecoveryDevice.verify|verify} messages.
     * @function encode
     * @memberof RecoveryDevice
     * @static
     * @param {IRecoveryDevice} message RecoveryDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    RecoveryDevice.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.wordCount != null && message.hasOwnProperty("wordCount"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.wordCount);
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            writer.uint32(/* id 2, wireType 0 =*/16).bool(message.passphraseProtection);
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            writer.uint32(/* id 3, wireType 0 =*/24).bool(message.pinProtection);
        if (message.language != null && message.hasOwnProperty("language"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.language);
        if (message.label != null && message.hasOwnProperty("label"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.label);
        if (message.enforceWordlist != null && message.hasOwnProperty("enforceWordlist"))
            writer.uint32(/* id 6, wireType 0 =*/48).bool(message.enforceWordlist);
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 8, wireType 0 =*/64).uint32(message.type);
        if (message.dryRun != null && message.hasOwnProperty("dryRun"))
            writer.uint32(/* id 10, wireType 0 =*/80).bool(message.dryRun);
        return writer;
    };

    /**
     * Encodes the specified RecoveryDevice message, length delimited. Does not implicitly {@link RecoveryDevice.verify|verify} messages.
     * @function encodeDelimited
     * @memberof RecoveryDevice
     * @static
     * @param {IRecoveryDevice} message RecoveryDevice message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    RecoveryDevice.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a RecoveryDevice message from the specified reader or buffer.
     * @function decode
     * @memberof RecoveryDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {RecoveryDevice} RecoveryDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    RecoveryDevice.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.RecoveryDevice();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.wordCount = reader.uint32();
                break;
            case 2:
                message.passphraseProtection = reader.bool();
                break;
            case 3:
                message.pinProtection = reader.bool();
                break;
            case 4:
                message.language = reader.string();
                break;
            case 5:
                message.label = reader.string();
                break;
            case 6:
                message.enforceWordlist = reader.bool();
                break;
            case 8:
                message.type = reader.uint32();
                break;
            case 10:
                message.dryRun = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a RecoveryDevice message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof RecoveryDevice
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {RecoveryDevice} RecoveryDevice
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    RecoveryDevice.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a RecoveryDevice message.
     * @function verify
     * @memberof RecoveryDevice
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    RecoveryDevice.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.wordCount != null && message.hasOwnProperty("wordCount"))
            if (!$util.isInteger(message.wordCount))
                return "wordCount: integer expected";
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            if (typeof message.passphraseProtection !== "boolean")
                return "passphraseProtection: boolean expected";
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            if (typeof message.pinProtection !== "boolean")
                return "pinProtection: boolean expected";
        if (message.language != null && message.hasOwnProperty("language"))
            if (!$util.isString(message.language))
                return "language: string expected";
        if (message.label != null && message.hasOwnProperty("label"))
            if (!$util.isString(message.label))
                return "label: string expected";
        if (message.enforceWordlist != null && message.hasOwnProperty("enforceWordlist"))
            if (typeof message.enforceWordlist !== "boolean")
                return "enforceWordlist: boolean expected";
        if (message.type != null && message.hasOwnProperty("type"))
            if (!$util.isInteger(message.type))
                return "type: integer expected";
        if (message.dryRun != null && message.hasOwnProperty("dryRun"))
            if (typeof message.dryRun !== "boolean")
                return "dryRun: boolean expected";
        return null;
    };

    /**
     * Creates a RecoveryDevice message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof RecoveryDevice
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {RecoveryDevice} RecoveryDevice
     */
    RecoveryDevice.fromObject = function fromObject(object) {
        if (object instanceof $root.RecoveryDevice)
            return object;
        var message = new $root.RecoveryDevice();
        if (object.wordCount != null)
            message.wordCount = object.wordCount >>> 0;
        if (object.passphraseProtection != null)
            message.passphraseProtection = Boolean(object.passphraseProtection);
        if (object.pinProtection != null)
            message.pinProtection = Boolean(object.pinProtection);
        if (object.language != null)
            message.language = String(object.language);
        if (object.label != null)
            message.label = String(object.label);
        if (object.enforceWordlist != null)
            message.enforceWordlist = Boolean(object.enforceWordlist);
        if (object.type != null)
            message.type = object.type >>> 0;
        if (object.dryRun != null)
            message.dryRun = Boolean(object.dryRun);
        return message;
    };

    /**
     * Creates a plain object from a RecoveryDevice message. Also converts values to other types if specified.
     * @function toObject
     * @memberof RecoveryDevice
     * @static
     * @param {RecoveryDevice} message RecoveryDevice
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    RecoveryDevice.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.wordCount = 0;
            object.passphraseProtection = false;
            object.pinProtection = false;
            object.language = "english";
            object.label = "";
            object.enforceWordlist = false;
            object.type = 0;
            object.dryRun = false;
        }
        if (message.wordCount != null && message.hasOwnProperty("wordCount"))
            object.wordCount = message.wordCount;
        if (message.passphraseProtection != null && message.hasOwnProperty("passphraseProtection"))
            object.passphraseProtection = message.passphraseProtection;
        if (message.pinProtection != null && message.hasOwnProperty("pinProtection"))
            object.pinProtection = message.pinProtection;
        if (message.language != null && message.hasOwnProperty("language"))
            object.language = message.language;
        if (message.label != null && message.hasOwnProperty("label"))
            object.label = message.label;
        if (message.enforceWordlist != null && message.hasOwnProperty("enforceWordlist"))
            object.enforceWordlist = message.enforceWordlist;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = message.type;
        if (message.dryRun != null && message.hasOwnProperty("dryRun"))
            object.dryRun = message.dryRun;
        return object;
    };

    /**
     * Converts this RecoveryDevice to JSON.
     * @function toJSON
     * @memberof RecoveryDevice
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    RecoveryDevice.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return RecoveryDevice;
})();

$root.WordRequest = (function() {

    /**
     * Properties of a WordRequest.
     * @exports IWordRequest
     * @interface IWordRequest
     * @property {WordRequestType|null} [type] WordRequest type
     */

    /**
     * Constructs a new WordRequest.
     * @exports WordRequest
     * @classdesc Response: Device is waiting for user to enter word of the mnemonic
     * Its position is shown only on device's internal display.
     * @prev RecoveryDevice
     * @prev WordAck
     * @implements IWordRequest
     * @constructor
     * @param {IWordRequest=} [properties] Properties to set
     */
    function WordRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * WordRequest type.
     * @member {WordRequestType} type
     * @memberof WordRequest
     * @instance
     */
    WordRequest.prototype.type = 0;

    /**
     * Creates a new WordRequest instance using the specified properties.
     * @function create
     * @memberof WordRequest
     * @static
     * @param {IWordRequest=} [properties] Properties to set
     * @returns {WordRequest} WordRequest instance
     */
    WordRequest.create = function create(properties) {
        return new WordRequest(properties);
    };

    /**
     * Encodes the specified WordRequest message. Does not implicitly {@link WordRequest.verify|verify} messages.
     * @function encode
     * @memberof WordRequest
     * @static
     * @param {IWordRequest} message WordRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    WordRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.type != null && message.hasOwnProperty("type"))
            writer.uint32(/* id 1, wireType 0 =*/8).int32(message.type);
        return writer;
    };

    /**
     * Encodes the specified WordRequest message, length delimited. Does not implicitly {@link WordRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof WordRequest
     * @static
     * @param {IWordRequest} message WordRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    WordRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a WordRequest message from the specified reader or buffer.
     * @function decode
     * @memberof WordRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {WordRequest} WordRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    WordRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.WordRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.type = reader.int32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a WordRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof WordRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {WordRequest} WordRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    WordRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a WordRequest message.
     * @function verify
     * @memberof WordRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    WordRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.type != null && message.hasOwnProperty("type"))
            switch (message.type) {
            default:
                return "type: enum value expected";
            case 0:
            case 1:
            case 2:
                break;
            }
        return null;
    };

    /**
     * Creates a WordRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof WordRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {WordRequest} WordRequest
     */
    WordRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.WordRequest)
            return object;
        var message = new $root.WordRequest();
        switch (object.type) {
        case "WordRequestType_Plain":
        case 0:
            message.type = 0;
            break;
        case "WordRequestType_Matrix9":
        case 1:
            message.type = 1;
            break;
        case "WordRequestType_Matrix6":
        case 2:
            message.type = 2;
            break;
        }
        return message;
    };

    /**
     * Creates a plain object from a WordRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof WordRequest
     * @static
     * @param {WordRequest} message WordRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    WordRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.type = options.enums === String ? "WordRequestType_Plain" : 0;
        if (message.type != null && message.hasOwnProperty("type"))
            object.type = options.enums === String ? $root.WordRequestType[message.type] : message.type;
        return object;
    };

    /**
     * Converts this WordRequest to JSON.
     * @function toJSON
     * @memberof WordRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    WordRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return WordRequest;
})();

$root.WordAck = (function() {

    /**
     * Properties of a WordAck.
     * @exports IWordAck
     * @interface IWordAck
     * @property {string} word WordAck word
     */

    /**
     * Constructs a new WordAck.
     * @exports WordAck
     * @classdesc Request: Computer replies with word from the mnemonic
     * @prev WordRequest
     * @next WordRequest
     * @next Success
     * @next Failure
     * @implements IWordAck
     * @constructor
     * @param {IWordAck=} [properties] Properties to set
     */
    function WordAck(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * WordAck word.
     * @member {string} word
     * @memberof WordAck
     * @instance
     */
    WordAck.prototype.word = "";

    /**
     * Creates a new WordAck instance using the specified properties.
     * @function create
     * @memberof WordAck
     * @static
     * @param {IWordAck=} [properties] Properties to set
     * @returns {WordAck} WordAck instance
     */
    WordAck.create = function create(properties) {
        return new WordAck(properties);
    };

    /**
     * Encodes the specified WordAck message. Does not implicitly {@link WordAck.verify|verify} messages.
     * @function encode
     * @memberof WordAck
     * @static
     * @param {IWordAck} message WordAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    WordAck.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 2 =*/10).string(message.word);
        return writer;
    };

    /**
     * Encodes the specified WordAck message, length delimited. Does not implicitly {@link WordAck.verify|verify} messages.
     * @function encodeDelimited
     * @memberof WordAck
     * @static
     * @param {IWordAck} message WordAck message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    WordAck.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a WordAck message from the specified reader or buffer.
     * @function decode
     * @memberof WordAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {WordAck} WordAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    WordAck.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.WordAck();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.word = reader.string();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("word"))
            throw $util.ProtocolError("missing required 'word'", { instance: message });
        return message;
    };

    /**
     * Decodes a WordAck message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof WordAck
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {WordAck} WordAck
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    WordAck.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a WordAck message.
     * @function verify
     * @memberof WordAck
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    WordAck.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isString(message.word))
            return "word: string expected";
        return null;
    };

    /**
     * Creates a WordAck message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof WordAck
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {WordAck} WordAck
     */
    WordAck.fromObject = function fromObject(object) {
        if (object instanceof $root.WordAck)
            return object;
        var message = new $root.WordAck();
        if (object.word != null)
            message.word = String(object.word);
        return message;
    };

    /**
     * Creates a plain object from a WordAck message. Also converts values to other types if specified.
     * @function toObject
     * @memberof WordAck
     * @static
     * @param {WordAck} message WordAck
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    WordAck.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.word = "";
        if (message.word != null && message.hasOwnProperty("word"))
            object.word = message.word;
        return object;
    };

    /**
     * Converts this WordAck to JSON.
     * @function toJSON
     * @memberof WordAck
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    WordAck.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return WordAck;
})();

$root.FirmwareErase = (function() {

    /**
     * Properties of a FirmwareErase.
     * @exports IFirmwareErase
     * @interface IFirmwareErase
     * @property {number|null} [length] FirmwareErase length
     */

    /**
     * Constructs a new FirmwareErase.
     * @exports FirmwareErase
     * @classdesc Request: Ask device to erase its firmware (so it can be replaced via FirmwareUpload)
     * @start
     * @next FirmwareRequest
     * @implements IFirmwareErase
     * @constructor
     * @param {IFirmwareErase=} [properties] Properties to set
     */
    function FirmwareErase(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FirmwareErase length.
     * @member {number} length
     * @memberof FirmwareErase
     * @instance
     */
    FirmwareErase.prototype.length = 0;

    /**
     * Creates a new FirmwareErase instance using the specified properties.
     * @function create
     * @memberof FirmwareErase
     * @static
     * @param {IFirmwareErase=} [properties] Properties to set
     * @returns {FirmwareErase} FirmwareErase instance
     */
    FirmwareErase.create = function create(properties) {
        return new FirmwareErase(properties);
    };

    /**
     * Encodes the specified FirmwareErase message. Does not implicitly {@link FirmwareErase.verify|verify} messages.
     * @function encode
     * @memberof FirmwareErase
     * @static
     * @param {IFirmwareErase} message FirmwareErase message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FirmwareErase.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.length != null && message.hasOwnProperty("length"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.length);
        return writer;
    };

    /**
     * Encodes the specified FirmwareErase message, length delimited. Does not implicitly {@link FirmwareErase.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FirmwareErase
     * @static
     * @param {IFirmwareErase} message FirmwareErase message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FirmwareErase.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FirmwareErase message from the specified reader or buffer.
     * @function decode
     * @memberof FirmwareErase
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FirmwareErase} FirmwareErase
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FirmwareErase.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.FirmwareErase();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.length = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FirmwareErase message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FirmwareErase
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FirmwareErase} FirmwareErase
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FirmwareErase.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FirmwareErase message.
     * @function verify
     * @memberof FirmwareErase
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FirmwareErase.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.length != null && message.hasOwnProperty("length"))
            if (!$util.isInteger(message.length))
                return "length: integer expected";
        return null;
    };

    /**
     * Creates a FirmwareErase message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FirmwareErase
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FirmwareErase} FirmwareErase
     */
    FirmwareErase.fromObject = function fromObject(object) {
        if (object instanceof $root.FirmwareErase)
            return object;
        var message = new $root.FirmwareErase();
        if (object.length != null)
            message.length = object.length >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a FirmwareErase message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FirmwareErase
     * @static
     * @param {FirmwareErase} message FirmwareErase
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FirmwareErase.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults)
            object.length = 0;
        if (message.length != null && message.hasOwnProperty("length"))
            object.length = message.length;
        return object;
    };

    /**
     * Converts this FirmwareErase to JSON.
     * @function toJSON
     * @memberof FirmwareErase
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FirmwareErase.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FirmwareErase;
})();

$root.FirmwareRequest = (function() {

    /**
     * Properties of a FirmwareRequest.
     * @exports IFirmwareRequest
     * @interface IFirmwareRequest
     * @property {number|null} [offset] FirmwareRequest offset
     * @property {number|null} [length] FirmwareRequest length
     */

    /**
     * Constructs a new FirmwareRequest.
     * @exports FirmwareRequest
     * @classdesc Response: Ask for firmware chunk
     * @next FirmwareUpload
     * @implements IFirmwareRequest
     * @constructor
     * @param {IFirmwareRequest=} [properties] Properties to set
     */
    function FirmwareRequest(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * FirmwareRequest offset.
     * @member {number} offset
     * @memberof FirmwareRequest
     * @instance
     */
    FirmwareRequest.prototype.offset = 0;

    /**
     * FirmwareRequest length.
     * @member {number} length
     * @memberof FirmwareRequest
     * @instance
     */
    FirmwareRequest.prototype.length = 0;

    /**
     * Creates a new FirmwareRequest instance using the specified properties.
     * @function create
     * @memberof FirmwareRequest
     * @static
     * @param {IFirmwareRequest=} [properties] Properties to set
     * @returns {FirmwareRequest} FirmwareRequest instance
     */
    FirmwareRequest.create = function create(properties) {
        return new FirmwareRequest(properties);
    };

    /**
     * Encodes the specified FirmwareRequest message. Does not implicitly {@link FirmwareRequest.verify|verify} messages.
     * @function encode
     * @memberof FirmwareRequest
     * @static
     * @param {IFirmwareRequest} message FirmwareRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FirmwareRequest.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.offset != null && message.hasOwnProperty("offset"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.offset);
        if (message.length != null && message.hasOwnProperty("length"))
            writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.length);
        return writer;
    };

    /**
     * Encodes the specified FirmwareRequest message, length delimited. Does not implicitly {@link FirmwareRequest.verify|verify} messages.
     * @function encodeDelimited
     * @memberof FirmwareRequest
     * @static
     * @param {IFirmwareRequest} message FirmwareRequest message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    FirmwareRequest.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a FirmwareRequest message from the specified reader or buffer.
     * @function decode
     * @memberof FirmwareRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {FirmwareRequest} FirmwareRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FirmwareRequest.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.FirmwareRequest();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.offset = reader.uint32();
                break;
            case 2:
                message.length = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a FirmwareRequest message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof FirmwareRequest
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {FirmwareRequest} FirmwareRequest
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    FirmwareRequest.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a FirmwareRequest message.
     * @function verify
     * @memberof FirmwareRequest
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    FirmwareRequest.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.offset != null && message.hasOwnProperty("offset"))
            if (!$util.isInteger(message.offset))
                return "offset: integer expected";
        if (message.length != null && message.hasOwnProperty("length"))
            if (!$util.isInteger(message.length))
                return "length: integer expected";
        return null;
    };

    /**
     * Creates a FirmwareRequest message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof FirmwareRequest
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {FirmwareRequest} FirmwareRequest
     */
    FirmwareRequest.fromObject = function fromObject(object) {
        if (object instanceof $root.FirmwareRequest)
            return object;
        var message = new $root.FirmwareRequest();
        if (object.offset != null)
            message.offset = object.offset >>> 0;
        if (object.length != null)
            message.length = object.length >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a FirmwareRequest message. Also converts values to other types if specified.
     * @function toObject
     * @memberof FirmwareRequest
     * @static
     * @param {FirmwareRequest} message FirmwareRequest
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    FirmwareRequest.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.offset = 0;
            object.length = 0;
        }
        if (message.offset != null && message.hasOwnProperty("offset"))
            object.offset = message.offset;
        if (message.length != null && message.hasOwnProperty("length"))
            object.length = message.length;
        return object;
    };

    /**
     * Converts this FirmwareRequest to JSON.
     * @function toJSON
     * @memberof FirmwareRequest
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    FirmwareRequest.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return FirmwareRequest;
})();

/**
 * Type of failures returned by Failure message
 * @used_in Failure
 * @exports FailureType
 * @enum {string}
 * @property {number} Failure_UnexpectedMessage=1 Failure_UnexpectedMessage value
 * @property {number} Failure_ButtonExpected=2 Failure_ButtonExpected value
 * @property {number} Failure_DataError=3 Failure_DataError value
 * @property {number} Failure_ActionCancelled=4 Failure_ActionCancelled value
 * @property {number} Failure_PinExpected=5 Failure_PinExpected value
 * @property {number} Failure_PinCancelled=6 Failure_PinCancelled value
 * @property {number} Failure_PinInvalid=7 Failure_PinInvalid value
 * @property {number} Failure_InvalidSignature=8 Failure_InvalidSignature value
 * @property {number} Failure_ProcessError=9 Failure_ProcessError value
 * @property {number} Failure_NotEnoughFunds=10 Failure_NotEnoughFunds value
 * @property {number} Failure_NotInitialized=11 Failure_NotInitialized value
 * @property {number} Failure_PinMismatch=12 Failure_PinMismatch value
 * @property {number} Failure_AddressGeneration=13 Failure_AddressGeneration value
 * @property {number} Failure_FirmwareError=99 Failure_FirmwareError value
 */
$root.FailureType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[1] = "Failure_UnexpectedMessage"] = 1;
    values[valuesById[2] = "Failure_ButtonExpected"] = 2;
    values[valuesById[3] = "Failure_DataError"] = 3;
    values[valuesById[4] = "Failure_ActionCancelled"] = 4;
    values[valuesById[5] = "Failure_PinExpected"] = 5;
    values[valuesById[6] = "Failure_PinCancelled"] = 6;
    values[valuesById[7] = "Failure_PinInvalid"] = 7;
    values[valuesById[8] = "Failure_InvalidSignature"] = 8;
    values[valuesById[9] = "Failure_ProcessError"] = 9;
    values[valuesById[10] = "Failure_NotEnoughFunds"] = 10;
    values[valuesById[11] = "Failure_NotInitialized"] = 11;
    values[valuesById[12] = "Failure_PinMismatch"] = 12;
    values[valuesById[13] = "Failure_AddressGeneration"] = 13;
    values[valuesById[99] = "Failure_FirmwareError"] = 99;
    return values;
})();

/**
 * Type of script which will be used for transaction output
 * @used_in TxOutputType
 * @exports OutputScriptType
 * @enum {string}
 * @property {number} PAYTOADDRESS=0 PAYTOADDRESS value
 * @property {number} PAYTOSCRIPTHASH=1 PAYTOSCRIPTHASH value
 * @property {number} PAYTOMULTISIG=2 PAYTOMULTISIG value
 * @property {number} PAYTOOPRETURN=3 PAYTOOPRETURN value
 * @property {number} PAYTOWITNESS=4 PAYTOWITNESS value
 * @property {number} PAYTOP2SHWITNESS=5 PAYTOP2SHWITNESS value
 */
$root.OutputScriptType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "PAYTOADDRESS"] = 0;
    values[valuesById[1] = "PAYTOSCRIPTHASH"] = 1;
    values[valuesById[2] = "PAYTOMULTISIG"] = 2;
    values[valuesById[3] = "PAYTOOPRETURN"] = 3;
    values[valuesById[4] = "PAYTOWITNESS"] = 4;
    values[valuesById[5] = "PAYTOP2SHWITNESS"] = 5;
    return values;
})();

/**
 * Type of script which will be used for transaction output
 * @used_in TxInputType
 * @exports InputScriptType
 * @enum {string}
 * @property {number} SPENDADDRESS=0 SPENDADDRESS value
 * @property {number} SPENDMULTISIG=1 SPENDMULTISIG value
 * @property {number} EXTERNAL=2 EXTERNAL value
 * @property {number} SPENDWITNESS=3 SPENDWITNESS value
 * @property {number} SPENDP2SHWITNESS=4 SPENDP2SHWITNESS value
 */
$root.InputScriptType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "SPENDADDRESS"] = 0;
    values[valuesById[1] = "SPENDMULTISIG"] = 1;
    values[valuesById[2] = "EXTERNAL"] = 2;
    values[valuesById[3] = "SPENDWITNESS"] = 3;
    values[valuesById[4] = "SPENDP2SHWITNESS"] = 4;
    return values;
})();

/**
 * Type of information required by transaction signing process
 * @used_in TxRequest
 * @exports RequestType
 * @enum {string}
 * @property {number} TXINPUT=0 TXINPUT value
 * @property {number} TXOUTPUT=1 TXOUTPUT value
 * @property {number} TXMETA=2 TXMETA value
 * @property {number} TXFINISHED=3 TXFINISHED value
 * @property {number} TXEXTRADATA=4 TXEXTRADATA value
 */
$root.RequestType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "TXINPUT"] = 0;
    values[valuesById[1] = "TXOUTPUT"] = 1;
    values[valuesById[2] = "TXMETA"] = 2;
    values[valuesById[3] = "TXFINISHED"] = 3;
    values[valuesById[4] = "TXEXTRADATA"] = 4;
    return values;
})();

/**
 * Type of button request
 * @used_in ButtonRequest
 * @exports ButtonRequestType
 * @enum {string}
 * @property {number} ButtonRequest_Other=1 ButtonRequest_Other value
 * @property {number} ButtonRequest_FeeOverThreshold=2 ButtonRequest_FeeOverThreshold value
 * @property {number} ButtonRequest_ConfirmOutput=3 ButtonRequest_ConfirmOutput value
 * @property {number} ButtonRequest_ResetDevice=4 ButtonRequest_ResetDevice value
 * @property {number} ButtonRequest_ConfirmWord=5 ButtonRequest_ConfirmWord value
 * @property {number} ButtonRequest_WipeDevice=6 ButtonRequest_WipeDevice value
 * @property {number} ButtonRequest_ProtectCall=7 ButtonRequest_ProtectCall value
 * @property {number} ButtonRequest_SignTx=8 ButtonRequest_SignTx value
 * @property {number} ButtonRequest_FirmwareCheck=9 ButtonRequest_FirmwareCheck value
 * @property {number} ButtonRequest_Address=10 ButtonRequest_Address value
 * @property {number} ButtonRequest_PublicKey=11 ButtonRequest_PublicKey value
 * @property {number} ButtonRequest_MnemonicWordCount=12 ButtonRequest_MnemonicWordCount value
 * @property {number} ButtonRequest_MnemonicInput=13 ButtonRequest_MnemonicInput value
 * @property {number} ButtonRequest_PassphraseType=14 ButtonRequest_PassphraseType value
 */
$root.ButtonRequestType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[1] = "ButtonRequest_Other"] = 1;
    values[valuesById[2] = "ButtonRequest_FeeOverThreshold"] = 2;
    values[valuesById[3] = "ButtonRequest_ConfirmOutput"] = 3;
    values[valuesById[4] = "ButtonRequest_ResetDevice"] = 4;
    values[valuesById[5] = "ButtonRequest_ConfirmWord"] = 5;
    values[valuesById[6] = "ButtonRequest_WipeDevice"] = 6;
    values[valuesById[7] = "ButtonRequest_ProtectCall"] = 7;
    values[valuesById[8] = "ButtonRequest_SignTx"] = 8;
    values[valuesById[9] = "ButtonRequest_FirmwareCheck"] = 9;
    values[valuesById[10] = "ButtonRequest_Address"] = 10;
    values[valuesById[11] = "ButtonRequest_PublicKey"] = 11;
    values[valuesById[12] = "ButtonRequest_MnemonicWordCount"] = 12;
    values[valuesById[13] = "ButtonRequest_MnemonicInput"] = 13;
    values[valuesById[14] = "ButtonRequest_PassphraseType"] = 14;
    return values;
})();

/**
 * Type of PIN request
 * @used_in PinMatrixRequest
 * @exports PinMatrixRequestType
 * @enum {string}
 * @property {number} PinMatrixRequestType_Current=1 PinMatrixRequestType_Current value
 * @property {number} PinMatrixRequestType_NewFirst=2 PinMatrixRequestType_NewFirst value
 * @property {number} PinMatrixRequestType_NewSecond=3 PinMatrixRequestType_NewSecond value
 */
$root.PinMatrixRequestType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[1] = "PinMatrixRequestType_Current"] = 1;
    values[valuesById[2] = "PinMatrixRequestType_NewFirst"] = 2;
    values[valuesById[3] = "PinMatrixRequestType_NewSecond"] = 3;
    return values;
})();

/**
 * Type of recovery procedure. These should be used as bitmask, e.g.,
 * `RecoveryDeviceType_ScrambledWords | RecoveryDeviceType_Matrix`
 * listing every method supported by the host computer.
 * 
 * Note that ScrambledWords must be supported by every implementation
 * for backward compatibility; there is no way to not support it.
 * 
 * @used_in RecoveryDevice
 * @exports RecoveryDeviceType
 * @enum {string}
 * @property {number} RecoveryDeviceType_ScrambledWords=0 RecoveryDeviceType_ScrambledWords value
 * @property {number} RecoveryDeviceType_Matrix=1 RecoveryDeviceType_Matrix value
 */
$root.RecoveryDeviceType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "RecoveryDeviceType_ScrambledWords"] = 0;
    values[valuesById[1] = "RecoveryDeviceType_Matrix"] = 1;
    return values;
})();

/**
 * Type of Recovery Word request
 * @used_in WordRequest
 * @exports WordRequestType
 * @enum {string}
 * @property {number} WordRequestType_Plain=0 WordRequestType_Plain value
 * @property {number} WordRequestType_Matrix9=1 WordRequestType_Matrix9 value
 * @property {number} WordRequestType_Matrix6=2 WordRequestType_Matrix6 value
 */
$root.WordRequestType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[0] = "WordRequestType_Plain"] = 0;
    values[valuesById[1] = "WordRequestType_Matrix9"] = 1;
    values[valuesById[2] = "WordRequestType_Matrix6"] = 2;
    return values;
})();

$root.HDNodeType = (function() {

    /**
     * Properties of a HDNodeType.
     * @exports IHDNodeType
     * @interface IHDNodeType
     * @property {number} depth HDNodeType depth
     * @property {number} fingerprint HDNodeType fingerprint
     * @property {number} childNum HDNodeType childNum
     * @property {Uint8Array} chainCode HDNodeType chainCode
     * @property {Uint8Array|null} [privateKey] HDNodeType privateKey
     * @property {Uint8Array|null} [publicKey] HDNodeType publicKey
     */

    /**
     * Constructs a new HDNodeType.
     * @exports HDNodeType
     * @classdesc Structure representing BIP32 (hierarchical deterministic) node
     * Used for imports of private key into the device and exporting public key out of device
     * @used_in PublicKey
     * @used_in LoadDevice
     * @used_in DebugLinkState
     * @used_in Storage
     * @implements IHDNodeType
     * @constructor
     * @param {IHDNodeType=} [properties] Properties to set
     */
    function HDNodeType(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * HDNodeType depth.
     * @member {number} depth
     * @memberof HDNodeType
     * @instance
     */
    HDNodeType.prototype.depth = 0;

    /**
     * HDNodeType fingerprint.
     * @member {number} fingerprint
     * @memberof HDNodeType
     * @instance
     */
    HDNodeType.prototype.fingerprint = 0;

    /**
     * HDNodeType childNum.
     * @member {number} childNum
     * @memberof HDNodeType
     * @instance
     */
    HDNodeType.prototype.childNum = 0;

    /**
     * HDNodeType chainCode.
     * @member {Uint8Array} chainCode
     * @memberof HDNodeType
     * @instance
     */
    HDNodeType.prototype.chainCode = $util.newBuffer([]);

    /**
     * HDNodeType privateKey.
     * @member {Uint8Array} privateKey
     * @memberof HDNodeType
     * @instance
     */
    HDNodeType.prototype.privateKey = $util.newBuffer([]);

    /**
     * HDNodeType publicKey.
     * @member {Uint8Array} publicKey
     * @memberof HDNodeType
     * @instance
     */
    HDNodeType.prototype.publicKey = $util.newBuffer([]);

    /**
     * Creates a new HDNodeType instance using the specified properties.
     * @function create
     * @memberof HDNodeType
     * @static
     * @param {IHDNodeType=} [properties] Properties to set
     * @returns {HDNodeType} HDNodeType instance
     */
    HDNodeType.create = function create(properties) {
        return new HDNodeType(properties);
    };

    /**
     * Encodes the specified HDNodeType message. Does not implicitly {@link HDNodeType.verify|verify} messages.
     * @function encode
     * @memberof HDNodeType
     * @static
     * @param {IHDNodeType} message HDNodeType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    HDNodeType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.depth);
        writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.fingerprint);
        writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.childNum);
        writer.uint32(/* id 4, wireType 2 =*/34).bytes(message.chainCode);
        if (message.privateKey != null && message.hasOwnProperty("privateKey"))
            writer.uint32(/* id 5, wireType 2 =*/42).bytes(message.privateKey);
        if (message.publicKey != null && message.hasOwnProperty("publicKey"))
            writer.uint32(/* id 6, wireType 2 =*/50).bytes(message.publicKey);
        return writer;
    };

    /**
     * Encodes the specified HDNodeType message, length delimited. Does not implicitly {@link HDNodeType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof HDNodeType
     * @static
     * @param {IHDNodeType} message HDNodeType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    HDNodeType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a HDNodeType message from the specified reader or buffer.
     * @function decode
     * @memberof HDNodeType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {HDNodeType} HDNodeType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    HDNodeType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.HDNodeType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.depth = reader.uint32();
                break;
            case 2:
                message.fingerprint = reader.uint32();
                break;
            case 3:
                message.childNum = reader.uint32();
                break;
            case 4:
                message.chainCode = reader.bytes();
                break;
            case 5:
                message.privateKey = reader.bytes();
                break;
            case 6:
                message.publicKey = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("depth"))
            throw $util.ProtocolError("missing required 'depth'", { instance: message });
        if (!message.hasOwnProperty("fingerprint"))
            throw $util.ProtocolError("missing required 'fingerprint'", { instance: message });
        if (!message.hasOwnProperty("childNum"))
            throw $util.ProtocolError("missing required 'childNum'", { instance: message });
        if (!message.hasOwnProperty("chainCode"))
            throw $util.ProtocolError("missing required 'chainCode'", { instance: message });
        return message;
    };

    /**
     * Decodes a HDNodeType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof HDNodeType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {HDNodeType} HDNodeType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    HDNodeType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a HDNodeType message.
     * @function verify
     * @memberof HDNodeType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    HDNodeType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isInteger(message.depth))
            return "depth: integer expected";
        if (!$util.isInteger(message.fingerprint))
            return "fingerprint: integer expected";
        if (!$util.isInteger(message.childNum))
            return "childNum: integer expected";
        if (!(message.chainCode && typeof message.chainCode.length === "number" || $util.isString(message.chainCode)))
            return "chainCode: buffer expected";
        if (message.privateKey != null && message.hasOwnProperty("privateKey"))
            if (!(message.privateKey && typeof message.privateKey.length === "number" || $util.isString(message.privateKey)))
                return "privateKey: buffer expected";
        if (message.publicKey != null && message.hasOwnProperty("publicKey"))
            if (!(message.publicKey && typeof message.publicKey.length === "number" || $util.isString(message.publicKey)))
                return "publicKey: buffer expected";
        return null;
    };

    /**
     * Creates a HDNodeType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof HDNodeType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {HDNodeType} HDNodeType
     */
    HDNodeType.fromObject = function fromObject(object) {
        if (object instanceof $root.HDNodeType)
            return object;
        var message = new $root.HDNodeType();
        if (object.depth != null)
            message.depth = object.depth >>> 0;
        if (object.fingerprint != null)
            message.fingerprint = object.fingerprint >>> 0;
        if (object.childNum != null)
            message.childNum = object.childNum >>> 0;
        if (object.chainCode != null)
            if (typeof object.chainCode === "string")
                $util.base64.decode(object.chainCode, message.chainCode = $util.newBuffer($util.base64.length(object.chainCode)), 0);
            else if (object.chainCode.length)
                message.chainCode = object.chainCode;
        if (object.privateKey != null)
            if (typeof object.privateKey === "string")
                $util.base64.decode(object.privateKey, message.privateKey = $util.newBuffer($util.base64.length(object.privateKey)), 0);
            else if (object.privateKey.length)
                message.privateKey = object.privateKey;
        if (object.publicKey != null)
            if (typeof object.publicKey === "string")
                $util.base64.decode(object.publicKey, message.publicKey = $util.newBuffer($util.base64.length(object.publicKey)), 0);
            else if (object.publicKey.length)
                message.publicKey = object.publicKey;
        return message;
    };

    /**
     * Creates a plain object from a HDNodeType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof HDNodeType
     * @static
     * @param {HDNodeType} message HDNodeType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    HDNodeType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.depth = 0;
            object.fingerprint = 0;
            object.childNum = 0;
            if (options.bytes === String)
                object.chainCode = "";
            else {
                object.chainCode = [];
                if (options.bytes !== Array)
                    object.chainCode = $util.newBuffer(object.chainCode);
            }
            if (options.bytes === String)
                object.privateKey = "";
            else {
                object.privateKey = [];
                if (options.bytes !== Array)
                    object.privateKey = $util.newBuffer(object.privateKey);
            }
            if (options.bytes === String)
                object.publicKey = "";
            else {
                object.publicKey = [];
                if (options.bytes !== Array)
                    object.publicKey = $util.newBuffer(object.publicKey);
            }
        }
        if (message.depth != null && message.hasOwnProperty("depth"))
            object.depth = message.depth;
        if (message.fingerprint != null && message.hasOwnProperty("fingerprint"))
            object.fingerprint = message.fingerprint;
        if (message.childNum != null && message.hasOwnProperty("childNum"))
            object.childNum = message.childNum;
        if (message.chainCode != null && message.hasOwnProperty("chainCode"))
            object.chainCode = options.bytes === String ? $util.base64.encode(message.chainCode, 0, message.chainCode.length) : options.bytes === Array ? Array.prototype.slice.call(message.chainCode) : message.chainCode;
        if (message.privateKey != null && message.hasOwnProperty("privateKey"))
            object.privateKey = options.bytes === String ? $util.base64.encode(message.privateKey, 0, message.privateKey.length) : options.bytes === Array ? Array.prototype.slice.call(message.privateKey) : message.privateKey;
        if (message.publicKey != null && message.hasOwnProperty("publicKey"))
            object.publicKey = options.bytes === String ? $util.base64.encode(message.publicKey, 0, message.publicKey.length) : options.bytes === Array ? Array.prototype.slice.call(message.publicKey) : message.publicKey;
        return object;
    };

    /**
     * Converts this HDNodeType to JSON.
     * @function toJSON
     * @memberof HDNodeType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    HDNodeType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return HDNodeType;
})();

$root.HDNodePathType = (function() {

    /**
     * Properties of a HDNodePathType.
     * @exports IHDNodePathType
     * @interface IHDNodePathType
     * @property {IHDNodeType} node HDNodePathType node
     * @property {Array.<number>|null} [addressN] HDNodePathType addressN
     */

    /**
     * Constructs a new HDNodePathType.
     * @exports HDNodePathType
     * @classdesc Represents a HDNodePathType.
     * @implements IHDNodePathType
     * @constructor
     * @param {IHDNodePathType=} [properties] Properties to set
     */
    function HDNodePathType(properties) {
        this.addressN = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * HDNodePathType node.
     * @member {IHDNodeType} node
     * @memberof HDNodePathType
     * @instance
     */
    HDNodePathType.prototype.node = null;

    /**
     * HDNodePathType addressN.
     * @member {Array.<number>} addressN
     * @memberof HDNodePathType
     * @instance
     */
    HDNodePathType.prototype.addressN = $util.emptyArray;

    /**
     * Creates a new HDNodePathType instance using the specified properties.
     * @function create
     * @memberof HDNodePathType
     * @static
     * @param {IHDNodePathType=} [properties] Properties to set
     * @returns {HDNodePathType} HDNodePathType instance
     */
    HDNodePathType.create = function create(properties) {
        return new HDNodePathType(properties);
    };

    /**
     * Encodes the specified HDNodePathType message. Does not implicitly {@link HDNodePathType.verify|verify} messages.
     * @function encode
     * @memberof HDNodePathType
     * @static
     * @param {IHDNodePathType} message HDNodePathType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    HDNodePathType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        $root.HDNodeType.encode(message.node, writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.addressN != null && message.addressN.length)
            for (var i = 0; i < message.addressN.length; ++i)
                writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.addressN[i]);
        return writer;
    };

    /**
     * Encodes the specified HDNodePathType message, length delimited. Does not implicitly {@link HDNodePathType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof HDNodePathType
     * @static
     * @param {IHDNodePathType} message HDNodePathType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    HDNodePathType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a HDNodePathType message from the specified reader or buffer.
     * @function decode
     * @memberof HDNodePathType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {HDNodePathType} HDNodePathType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    HDNodePathType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.HDNodePathType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.node = $root.HDNodeType.decode(reader, reader.uint32());
                break;
            case 2:
                if (!(message.addressN && message.addressN.length))
                    message.addressN = [];
                if ((tag & 7) === 2) {
                    var end2 = reader.uint32() + reader.pos;
                    while (reader.pos < end2)
                        message.addressN.push(reader.uint32());
                } else
                    message.addressN.push(reader.uint32());
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("node"))
            throw $util.ProtocolError("missing required 'node'", { instance: message });
        return message;
    };

    /**
     * Decodes a HDNodePathType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof HDNodePathType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {HDNodePathType} HDNodePathType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    HDNodePathType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a HDNodePathType message.
     * @function verify
     * @memberof HDNodePathType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    HDNodePathType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        {
            var error = $root.HDNodeType.verify(message.node);
            if (error)
                return "node." + error;
        }
        if (message.addressN != null && message.hasOwnProperty("addressN")) {
            if (!Array.isArray(message.addressN))
                return "addressN: array expected";
            for (var i = 0; i < message.addressN.length; ++i)
                if (!$util.isInteger(message.addressN[i]))
                    return "addressN: integer[] expected";
        }
        return null;
    };

    /**
     * Creates a HDNodePathType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof HDNodePathType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {HDNodePathType} HDNodePathType
     */
    HDNodePathType.fromObject = function fromObject(object) {
        if (object instanceof $root.HDNodePathType)
            return object;
        var message = new $root.HDNodePathType();
        if (object.node != null) {
            if (typeof object.node !== "object")
                throw TypeError(".HDNodePathType.node: object expected");
            message.node = $root.HDNodeType.fromObject(object.node);
        }
        if (object.addressN) {
            if (!Array.isArray(object.addressN))
                throw TypeError(".HDNodePathType.addressN: array expected");
            message.addressN = [];
            for (var i = 0; i < object.addressN.length; ++i)
                message.addressN[i] = object.addressN[i] >>> 0;
        }
        return message;
    };

    /**
     * Creates a plain object from a HDNodePathType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof HDNodePathType
     * @static
     * @param {HDNodePathType} message HDNodePathType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    HDNodePathType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults)
            object.addressN = [];
        if (options.defaults)
            object.node = null;
        if (message.node != null && message.hasOwnProperty("node"))
            object.node = $root.HDNodeType.toObject(message.node, options);
        if (message.addressN && message.addressN.length) {
            object.addressN = [];
            for (var j = 0; j < message.addressN.length; ++j)
                object.addressN[j] = message.addressN[j];
        }
        return object;
    };

    /**
     * Converts this HDNodePathType to JSON.
     * @function toJSON
     * @memberof HDNodePathType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    HDNodePathType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return HDNodePathType;
})();

$root.CoinType = (function() {

    /**
     * Properties of a CoinType.
     * @exports ICoinType
     * @interface ICoinType
     * @property {string|null} [coinName] CoinType coinName
     * @property {string|null} [coinShortcut] CoinType coinShortcut
     * @property {number|null} [addressType] CoinType addressType
     * @property {number|Long|null} [maxfeeKb] CoinType maxfeeKb
     * @property {number|null} [addressTypeP2sh] CoinType addressTypeP2sh
     * @property {string|null} [signedMessageHeader] CoinType signedMessageHeader
     * @property {number|null} [xpubMagic] CoinType xpubMagic
     * @property {number|null} [xprvMagic] CoinType xprvMagic
     * @property {boolean|null} [segwit] CoinType segwit
     * @property {number|null} [forkid] CoinType forkid
     * @property {boolean|null} [forceBip143] CoinType forceBip143
     */

    /**
     * Constructs a new CoinType.
     * @exports CoinType
     * @classdesc Structure representing Coin
     * @used_in Features
     * @implements ICoinType
     * @constructor
     * @param {ICoinType=} [properties] Properties to set
     */
    function CoinType(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * CoinType coinName.
     * @member {string} coinName
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.coinName = "";

    /**
     * CoinType coinShortcut.
     * @member {string} coinShortcut
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.coinShortcut = "";

    /**
     * CoinType addressType.
     * @member {number} addressType
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.addressType = 0;

    /**
     * CoinType maxfeeKb.
     * @member {number|Long} maxfeeKb
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.maxfeeKb = $util.Long ? $util.Long.fromBits(0,0,true) : 0;

    /**
     * CoinType addressTypeP2sh.
     * @member {number} addressTypeP2sh
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.addressTypeP2sh = 5;

    /**
     * CoinType signedMessageHeader.
     * @member {string} signedMessageHeader
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.signedMessageHeader = "";

    /**
     * CoinType xpubMagic.
     * @member {number} xpubMagic
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.xpubMagic = 76067358;

    /**
     * CoinType xprvMagic.
     * @member {number} xprvMagic
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.xprvMagic = 76066276;

    /**
     * CoinType segwit.
     * @member {boolean} segwit
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.segwit = false;

    /**
     * CoinType forkid.
     * @member {number} forkid
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.forkid = 0;

    /**
     * CoinType forceBip143.
     * @member {boolean} forceBip143
     * @memberof CoinType
     * @instance
     */
    CoinType.prototype.forceBip143 = false;

    /**
     * Creates a new CoinType instance using the specified properties.
     * @function create
     * @memberof CoinType
     * @static
     * @param {ICoinType=} [properties] Properties to set
     * @returns {CoinType} CoinType instance
     */
    CoinType.create = function create(properties) {
        return new CoinType(properties);
    };

    /**
     * Encodes the specified CoinType message. Does not implicitly {@link CoinType.verify|verify} messages.
     * @function encode
     * @memberof CoinType
     * @static
     * @param {ICoinType} message CoinType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CoinType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.coinName != null && message.hasOwnProperty("coinName"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.coinName);
        if (message.coinShortcut != null && message.hasOwnProperty("coinShortcut"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.coinShortcut);
        if (message.addressType != null && message.hasOwnProperty("addressType"))
            writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.addressType);
        if (message.maxfeeKb != null && message.hasOwnProperty("maxfeeKb"))
            writer.uint32(/* id 4, wireType 0 =*/32).uint64(message.maxfeeKb);
        if (message.addressTypeP2sh != null && message.hasOwnProperty("addressTypeP2sh"))
            writer.uint32(/* id 5, wireType 0 =*/40).uint32(message.addressTypeP2sh);
        if (message.signedMessageHeader != null && message.hasOwnProperty("signedMessageHeader"))
            writer.uint32(/* id 8, wireType 2 =*/66).string(message.signedMessageHeader);
        if (message.xpubMagic != null && message.hasOwnProperty("xpubMagic"))
            writer.uint32(/* id 9, wireType 0 =*/72).uint32(message.xpubMagic);
        if (message.xprvMagic != null && message.hasOwnProperty("xprvMagic"))
            writer.uint32(/* id 10, wireType 0 =*/80).uint32(message.xprvMagic);
        if (message.segwit != null && message.hasOwnProperty("segwit"))
            writer.uint32(/* id 11, wireType 0 =*/88).bool(message.segwit);
        if (message.forkid != null && message.hasOwnProperty("forkid"))
            writer.uint32(/* id 12, wireType 0 =*/96).uint32(message.forkid);
        if (message.forceBip143 != null && message.hasOwnProperty("forceBip143"))
            writer.uint32(/* id 13, wireType 0 =*/104).bool(message.forceBip143);
        return writer;
    };

    /**
     * Encodes the specified CoinType message, length delimited. Does not implicitly {@link CoinType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof CoinType
     * @static
     * @param {ICoinType} message CoinType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    CoinType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a CoinType message from the specified reader or buffer.
     * @function decode
     * @memberof CoinType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {CoinType} CoinType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CoinType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.CoinType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.coinName = reader.string();
                break;
            case 2:
                message.coinShortcut = reader.string();
                break;
            case 3:
                message.addressType = reader.uint32();
                break;
            case 4:
                message.maxfeeKb = reader.uint64();
                break;
            case 5:
                message.addressTypeP2sh = reader.uint32();
                break;
            case 8:
                message.signedMessageHeader = reader.string();
                break;
            case 9:
                message.xpubMagic = reader.uint32();
                break;
            case 10:
                message.xprvMagic = reader.uint32();
                break;
            case 11:
                message.segwit = reader.bool();
                break;
            case 12:
                message.forkid = reader.uint32();
                break;
            case 13:
                message.forceBip143 = reader.bool();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a CoinType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof CoinType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {CoinType} CoinType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    CoinType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a CoinType message.
     * @function verify
     * @memberof CoinType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    CoinType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.coinName != null && message.hasOwnProperty("coinName"))
            if (!$util.isString(message.coinName))
                return "coinName: string expected";
        if (message.coinShortcut != null && message.hasOwnProperty("coinShortcut"))
            if (!$util.isString(message.coinShortcut))
                return "coinShortcut: string expected";
        if (message.addressType != null && message.hasOwnProperty("addressType"))
            if (!$util.isInteger(message.addressType))
                return "addressType: integer expected";
        if (message.maxfeeKb != null && message.hasOwnProperty("maxfeeKb"))
            if (!$util.isInteger(message.maxfeeKb) && !(message.maxfeeKb && $util.isInteger(message.maxfeeKb.low) && $util.isInteger(message.maxfeeKb.high)))
                return "maxfeeKb: integer|Long expected";
        if (message.addressTypeP2sh != null && message.hasOwnProperty("addressTypeP2sh"))
            if (!$util.isInteger(message.addressTypeP2sh))
                return "addressTypeP2sh: integer expected";
        if (message.signedMessageHeader != null && message.hasOwnProperty("signedMessageHeader"))
            if (!$util.isString(message.signedMessageHeader))
                return "signedMessageHeader: string expected";
        if (message.xpubMagic != null && message.hasOwnProperty("xpubMagic"))
            if (!$util.isInteger(message.xpubMagic))
                return "xpubMagic: integer expected";
        if (message.xprvMagic != null && message.hasOwnProperty("xprvMagic"))
            if (!$util.isInteger(message.xprvMagic))
                return "xprvMagic: integer expected";
        if (message.segwit != null && message.hasOwnProperty("segwit"))
            if (typeof message.segwit !== "boolean")
                return "segwit: boolean expected";
        if (message.forkid != null && message.hasOwnProperty("forkid"))
            if (!$util.isInteger(message.forkid))
                return "forkid: integer expected";
        if (message.forceBip143 != null && message.hasOwnProperty("forceBip143"))
            if (typeof message.forceBip143 !== "boolean")
                return "forceBip143: boolean expected";
        return null;
    };

    /**
     * Creates a CoinType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof CoinType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {CoinType} CoinType
     */
    CoinType.fromObject = function fromObject(object) {
        if (object instanceof $root.CoinType)
            return object;
        var message = new $root.CoinType();
        if (object.coinName != null)
            message.coinName = String(object.coinName);
        if (object.coinShortcut != null)
            message.coinShortcut = String(object.coinShortcut);
        if (object.addressType != null)
            message.addressType = object.addressType >>> 0;
        if (object.maxfeeKb != null)
            if ($util.Long)
                (message.maxfeeKb = $util.Long.fromValue(object.maxfeeKb)).unsigned = true;
            else if (typeof object.maxfeeKb === "string")
                message.maxfeeKb = parseInt(object.maxfeeKb, 10);
            else if (typeof object.maxfeeKb === "number")
                message.maxfeeKb = object.maxfeeKb;
            else if (typeof object.maxfeeKb === "object")
                message.maxfeeKb = new $util.LongBits(object.maxfeeKb.low >>> 0, object.maxfeeKb.high >>> 0).toNumber(true);
        if (object.addressTypeP2sh != null)
            message.addressTypeP2sh = object.addressTypeP2sh >>> 0;
        if (object.signedMessageHeader != null)
            message.signedMessageHeader = String(object.signedMessageHeader);
        if (object.xpubMagic != null)
            message.xpubMagic = object.xpubMagic >>> 0;
        if (object.xprvMagic != null)
            message.xprvMagic = object.xprvMagic >>> 0;
        if (object.segwit != null)
            message.segwit = Boolean(object.segwit);
        if (object.forkid != null)
            message.forkid = object.forkid >>> 0;
        if (object.forceBip143 != null)
            message.forceBip143 = Boolean(object.forceBip143);
        return message;
    };

    /**
     * Creates a plain object from a CoinType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof CoinType
     * @static
     * @param {CoinType} message CoinType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    CoinType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.coinName = "";
            object.coinShortcut = "";
            object.addressType = 0;
            if ($util.Long) {
                var long = new $util.Long(0, 0, true);
                object.maxfeeKb = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.maxfeeKb = options.longs === String ? "0" : 0;
            object.addressTypeP2sh = 5;
            object.signedMessageHeader = "";
            object.xpubMagic = 76067358;
            object.xprvMagic = 76066276;
            object.segwit = false;
            object.forkid = 0;
            object.forceBip143 = false;
        }
        if (message.coinName != null && message.hasOwnProperty("coinName"))
            object.coinName = message.coinName;
        if (message.coinShortcut != null && message.hasOwnProperty("coinShortcut"))
            object.coinShortcut = message.coinShortcut;
        if (message.addressType != null && message.hasOwnProperty("addressType"))
            object.addressType = message.addressType;
        if (message.maxfeeKb != null && message.hasOwnProperty("maxfeeKb"))
            if (typeof message.maxfeeKb === "number")
                object.maxfeeKb = options.longs === String ? String(message.maxfeeKb) : message.maxfeeKb;
            else
                object.maxfeeKb = options.longs === String ? $util.Long.prototype.toString.call(message.maxfeeKb) : options.longs === Number ? new $util.LongBits(message.maxfeeKb.low >>> 0, message.maxfeeKb.high >>> 0).toNumber(true) : message.maxfeeKb;
        if (message.addressTypeP2sh != null && message.hasOwnProperty("addressTypeP2sh"))
            object.addressTypeP2sh = message.addressTypeP2sh;
        if (message.signedMessageHeader != null && message.hasOwnProperty("signedMessageHeader"))
            object.signedMessageHeader = message.signedMessageHeader;
        if (message.xpubMagic != null && message.hasOwnProperty("xpubMagic"))
            object.xpubMagic = message.xpubMagic;
        if (message.xprvMagic != null && message.hasOwnProperty("xprvMagic"))
            object.xprvMagic = message.xprvMagic;
        if (message.segwit != null && message.hasOwnProperty("segwit"))
            object.segwit = message.segwit;
        if (message.forkid != null && message.hasOwnProperty("forkid"))
            object.forkid = message.forkid;
        if (message.forceBip143 != null && message.hasOwnProperty("forceBip143"))
            object.forceBip143 = message.forceBip143;
        return object;
    };

    /**
     * Converts this CoinType to JSON.
     * @function toJSON
     * @memberof CoinType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    CoinType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return CoinType;
})();

$root.MultisigRedeemScriptType = (function() {

    /**
     * Properties of a MultisigRedeemScriptType.
     * @exports IMultisigRedeemScriptType
     * @interface IMultisigRedeemScriptType
     * @property {Array.<IHDNodePathType>|null} [pubkeys] MultisigRedeemScriptType pubkeys
     * @property {Array.<Uint8Array>|null} [signatures] MultisigRedeemScriptType signatures
     * @property {number|null} [m] MultisigRedeemScriptType m
     */

    /**
     * Constructs a new MultisigRedeemScriptType.
     * @exports MultisigRedeemScriptType
     * @classdesc Type of redeem script used in input
     * @used_in TxInputType
     * @implements IMultisigRedeemScriptType
     * @constructor
     * @param {IMultisigRedeemScriptType=} [properties] Properties to set
     */
    function MultisigRedeemScriptType(properties) {
        this.pubkeys = [];
        this.signatures = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * MultisigRedeemScriptType pubkeys.
     * @member {Array.<IHDNodePathType>} pubkeys
     * @memberof MultisigRedeemScriptType
     * @instance
     */
    MultisigRedeemScriptType.prototype.pubkeys = $util.emptyArray;

    /**
     * MultisigRedeemScriptType signatures.
     * @member {Array.<Uint8Array>} signatures
     * @memberof MultisigRedeemScriptType
     * @instance
     */
    MultisigRedeemScriptType.prototype.signatures = $util.emptyArray;

    /**
     * MultisigRedeemScriptType m.
     * @member {number} m
     * @memberof MultisigRedeemScriptType
     * @instance
     */
    MultisigRedeemScriptType.prototype.m = 0;

    /**
     * Creates a new MultisigRedeemScriptType instance using the specified properties.
     * @function create
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {IMultisigRedeemScriptType=} [properties] Properties to set
     * @returns {MultisigRedeemScriptType} MultisigRedeemScriptType instance
     */
    MultisigRedeemScriptType.create = function create(properties) {
        return new MultisigRedeemScriptType(properties);
    };

    /**
     * Encodes the specified MultisigRedeemScriptType message. Does not implicitly {@link MultisigRedeemScriptType.verify|verify} messages.
     * @function encode
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {IMultisigRedeemScriptType} message MultisigRedeemScriptType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    MultisigRedeemScriptType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.pubkeys != null && message.pubkeys.length)
            for (var i = 0; i < message.pubkeys.length; ++i)
                $root.HDNodePathType.encode(message.pubkeys[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
        if (message.signatures != null && message.signatures.length)
            for (var i = 0; i < message.signatures.length; ++i)
                writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.signatures[i]);
        if (message.m != null && message.hasOwnProperty("m"))
            writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.m);
        return writer;
    };

    /**
     * Encodes the specified MultisigRedeemScriptType message, length delimited. Does not implicitly {@link MultisigRedeemScriptType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {IMultisigRedeemScriptType} message MultisigRedeemScriptType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    MultisigRedeemScriptType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a MultisigRedeemScriptType message from the specified reader or buffer.
     * @function decode
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {MultisigRedeemScriptType} MultisigRedeemScriptType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    MultisigRedeemScriptType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.MultisigRedeemScriptType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.pubkeys && message.pubkeys.length))
                    message.pubkeys = [];
                message.pubkeys.push($root.HDNodePathType.decode(reader, reader.uint32()));
                break;
            case 2:
                if (!(message.signatures && message.signatures.length))
                    message.signatures = [];
                message.signatures.push(reader.bytes());
                break;
            case 3:
                message.m = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a MultisigRedeemScriptType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {MultisigRedeemScriptType} MultisigRedeemScriptType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    MultisigRedeemScriptType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a MultisigRedeemScriptType message.
     * @function verify
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    MultisigRedeemScriptType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.pubkeys != null && message.hasOwnProperty("pubkeys")) {
            if (!Array.isArray(message.pubkeys))
                return "pubkeys: array expected";
            for (var i = 0; i < message.pubkeys.length; ++i) {
                var error = $root.HDNodePathType.verify(message.pubkeys[i]);
                if (error)
                    return "pubkeys." + error;
            }
        }
        if (message.signatures != null && message.hasOwnProperty("signatures")) {
            if (!Array.isArray(message.signatures))
                return "signatures: array expected";
            for (var i = 0; i < message.signatures.length; ++i)
                if (!(message.signatures[i] && typeof message.signatures[i].length === "number" || $util.isString(message.signatures[i])))
                    return "signatures: buffer[] expected";
        }
        if (message.m != null && message.hasOwnProperty("m"))
            if (!$util.isInteger(message.m))
                return "m: integer expected";
        return null;
    };

    /**
     * Creates a MultisigRedeemScriptType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {MultisigRedeemScriptType} MultisigRedeemScriptType
     */
    MultisigRedeemScriptType.fromObject = function fromObject(object) {
        if (object instanceof $root.MultisigRedeemScriptType)
            return object;
        var message = new $root.MultisigRedeemScriptType();
        if (object.pubkeys) {
            if (!Array.isArray(object.pubkeys))
                throw TypeError(".MultisigRedeemScriptType.pubkeys: array expected");
            message.pubkeys = [];
            for (var i = 0; i < object.pubkeys.length; ++i) {
                if (typeof object.pubkeys[i] !== "object")
                    throw TypeError(".MultisigRedeemScriptType.pubkeys: object expected");
                message.pubkeys[i] = $root.HDNodePathType.fromObject(object.pubkeys[i]);
            }
        }
        if (object.signatures) {
            if (!Array.isArray(object.signatures))
                throw TypeError(".MultisigRedeemScriptType.signatures: array expected");
            message.signatures = [];
            for (var i = 0; i < object.signatures.length; ++i)
                if (typeof object.signatures[i] === "string")
                    $util.base64.decode(object.signatures[i], message.signatures[i] = $util.newBuffer($util.base64.length(object.signatures[i])), 0);
                else if (object.signatures[i].length)
                    message.signatures[i] = object.signatures[i];
        }
        if (object.m != null)
            message.m = object.m >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a MultisigRedeemScriptType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof MultisigRedeemScriptType
     * @static
     * @param {MultisigRedeemScriptType} message MultisigRedeemScriptType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    MultisigRedeemScriptType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults) {
            object.pubkeys = [];
            object.signatures = [];
        }
        if (options.defaults)
            object.m = 0;
        if (message.pubkeys && message.pubkeys.length) {
            object.pubkeys = [];
            for (var j = 0; j < message.pubkeys.length; ++j)
                object.pubkeys[j] = $root.HDNodePathType.toObject(message.pubkeys[j], options);
        }
        if (message.signatures && message.signatures.length) {
            object.signatures = [];
            for (var j = 0; j < message.signatures.length; ++j)
                object.signatures[j] = options.bytes === String ? $util.base64.encode(message.signatures[j], 0, message.signatures[j].length) : options.bytes === Array ? Array.prototype.slice.call(message.signatures[j]) : message.signatures[j];
        }
        if (message.m != null && message.hasOwnProperty("m"))
            object.m = message.m;
        return object;
    };

    /**
     * Converts this MultisigRedeemScriptType to JSON.
     * @function toJSON
     * @memberof MultisigRedeemScriptType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    MultisigRedeemScriptType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return MultisigRedeemScriptType;
})();

$root.TxInputType = (function() {

    /**
     * Properties of a TxInputType.
     * @exports ITxInputType
     * @interface ITxInputType
     * @property {Array.<number>|null} [addressN] TxInputType addressN
     * @property {Uint8Array} prevHash TxInputType prevHash
     * @property {number} prevIndex TxInputType prevIndex
     * @property {Uint8Array|null} [scriptSig] TxInputType scriptSig
     * @property {number|null} [sequence] TxInputType sequence
     * @property {InputScriptType|null} [scriptType] TxInputType scriptType
     * @property {IMultisigRedeemScriptType|null} [multisig] TxInputType multisig
     * @property {number|Long|null} [amount] TxInputType amount
     * @property {number|null} [decredTree] TxInputType decredTree
     * @property {number|null} [decredScriptVersion] TxInputType decredScriptVersion
     */

    /**
     * Constructs a new TxInputType.
     * @exports TxInputType
     * @classdesc Structure representing transaction input
     * @used_in SimpleSignTx
     * @used_in TransactionType
     * @implements ITxInputType
     * @constructor
     * @param {ITxInputType=} [properties] Properties to set
     */
    function TxInputType(properties) {
        this.addressN = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TxInputType addressN.
     * @member {Array.<number>} addressN
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.addressN = $util.emptyArray;

    /**
     * TxInputType prevHash.
     * @member {Uint8Array} prevHash
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.prevHash = $util.newBuffer([]);

    /**
     * TxInputType prevIndex.
     * @member {number} prevIndex
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.prevIndex = 0;

    /**
     * TxInputType scriptSig.
     * @member {Uint8Array} scriptSig
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.scriptSig = $util.newBuffer([]);

    /**
     * TxInputType sequence.
     * @member {number} sequence
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.sequence = 4294967295;

    /**
     * TxInputType scriptType.
     * @member {InputScriptType} scriptType
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.scriptType = 0;

    /**
     * TxInputType multisig.
     * @member {IMultisigRedeemScriptType|null|undefined} multisig
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.multisig = null;

    /**
     * TxInputType amount.
     * @member {number|Long} amount
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.amount = $util.Long ? $util.Long.fromBits(0,0,true) : 0;

    /**
     * TxInputType decredTree.
     * @member {number} decredTree
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.decredTree = 0;

    /**
     * TxInputType decredScriptVersion.
     * @member {number} decredScriptVersion
     * @memberof TxInputType
     * @instance
     */
    TxInputType.prototype.decredScriptVersion = 0;

    /**
     * Creates a new TxInputType instance using the specified properties.
     * @function create
     * @memberof TxInputType
     * @static
     * @param {ITxInputType=} [properties] Properties to set
     * @returns {TxInputType} TxInputType instance
     */
    TxInputType.create = function create(properties) {
        return new TxInputType(properties);
    };

    /**
     * Encodes the specified TxInputType message. Does not implicitly {@link TxInputType.verify|verify} messages.
     * @function encode
     * @memberof TxInputType
     * @static
     * @param {ITxInputType} message TxInputType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxInputType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.addressN != null && message.addressN.length)
            for (var i = 0; i < message.addressN.length; ++i)
                writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.addressN[i]);
        writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.prevHash);
        writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.prevIndex);
        if (message.scriptSig != null && message.hasOwnProperty("scriptSig"))
            writer.uint32(/* id 4, wireType 2 =*/34).bytes(message.scriptSig);
        if (message.sequence != null && message.hasOwnProperty("sequence"))
            writer.uint32(/* id 5, wireType 0 =*/40).uint32(message.sequence);
        if (message.scriptType != null && message.hasOwnProperty("scriptType"))
            writer.uint32(/* id 6, wireType 0 =*/48).int32(message.scriptType);
        if (message.multisig != null && message.hasOwnProperty("multisig"))
            $root.MultisigRedeemScriptType.encode(message.multisig, writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
        if (message.amount != null && message.hasOwnProperty("amount"))
            writer.uint32(/* id 8, wireType 0 =*/64).uint64(message.amount);
        if (message.decredTree != null && message.hasOwnProperty("decredTree"))
            writer.uint32(/* id 9, wireType 0 =*/72).uint32(message.decredTree);
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            writer.uint32(/* id 10, wireType 0 =*/80).uint32(message.decredScriptVersion);
        return writer;
    };

    /**
     * Encodes the specified TxInputType message, length delimited. Does not implicitly {@link TxInputType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TxInputType
     * @static
     * @param {ITxInputType} message TxInputType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxInputType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TxInputType message from the specified reader or buffer.
     * @function decode
     * @memberof TxInputType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TxInputType} TxInputType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxInputType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.TxInputType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                if (!(message.addressN && message.addressN.length))
                    message.addressN = [];
                if ((tag & 7) === 2) {
                    var end2 = reader.uint32() + reader.pos;
                    while (reader.pos < end2)
                        message.addressN.push(reader.uint32());
                } else
                    message.addressN.push(reader.uint32());
                break;
            case 2:
                message.prevHash = reader.bytes();
                break;
            case 3:
                message.prevIndex = reader.uint32();
                break;
            case 4:
                message.scriptSig = reader.bytes();
                break;
            case 5:
                message.sequence = reader.uint32();
                break;
            case 6:
                message.scriptType = reader.int32();
                break;
            case 7:
                message.multisig = $root.MultisigRedeemScriptType.decode(reader, reader.uint32());
                break;
            case 8:
                message.amount = reader.uint64();
                break;
            case 9:
                message.decredTree = reader.uint32();
                break;
            case 10:
                message.decredScriptVersion = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("prevHash"))
            throw $util.ProtocolError("missing required 'prevHash'", { instance: message });
        if (!message.hasOwnProperty("prevIndex"))
            throw $util.ProtocolError("missing required 'prevIndex'", { instance: message });
        return message;
    };

    /**
     * Decodes a TxInputType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TxInputType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TxInputType} TxInputType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxInputType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TxInputType message.
     * @function verify
     * @memberof TxInputType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TxInputType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.addressN != null && message.hasOwnProperty("addressN")) {
            if (!Array.isArray(message.addressN))
                return "addressN: array expected";
            for (var i = 0; i < message.addressN.length; ++i)
                if (!$util.isInteger(message.addressN[i]))
                    return "addressN: integer[] expected";
        }
        if (!(message.prevHash && typeof message.prevHash.length === "number" || $util.isString(message.prevHash)))
            return "prevHash: buffer expected";
        if (!$util.isInteger(message.prevIndex))
            return "prevIndex: integer expected";
        if (message.scriptSig != null && message.hasOwnProperty("scriptSig"))
            if (!(message.scriptSig && typeof message.scriptSig.length === "number" || $util.isString(message.scriptSig)))
                return "scriptSig: buffer expected";
        if (message.sequence != null && message.hasOwnProperty("sequence"))
            if (!$util.isInteger(message.sequence))
                return "sequence: integer expected";
        if (message.scriptType != null && message.hasOwnProperty("scriptType"))
            switch (message.scriptType) {
            default:
                return "scriptType: enum value expected";
            case 0:
            case 1:
            case 2:
            case 3:
            case 4:
                break;
            }
        if (message.multisig != null && message.hasOwnProperty("multisig")) {
            var error = $root.MultisigRedeemScriptType.verify(message.multisig);
            if (error)
                return "multisig." + error;
        }
        if (message.amount != null && message.hasOwnProperty("amount"))
            if (!$util.isInteger(message.amount) && !(message.amount && $util.isInteger(message.amount.low) && $util.isInteger(message.amount.high)))
                return "amount: integer|Long expected";
        if (message.decredTree != null && message.hasOwnProperty("decredTree"))
            if (!$util.isInteger(message.decredTree))
                return "decredTree: integer expected";
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            if (!$util.isInteger(message.decredScriptVersion))
                return "decredScriptVersion: integer expected";
        return null;
    };

    /**
     * Creates a TxInputType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TxInputType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TxInputType} TxInputType
     */
    TxInputType.fromObject = function fromObject(object) {
        if (object instanceof $root.TxInputType)
            return object;
        var message = new $root.TxInputType();
        if (object.addressN) {
            if (!Array.isArray(object.addressN))
                throw TypeError(".TxInputType.addressN: array expected");
            message.addressN = [];
            for (var i = 0; i < object.addressN.length; ++i)
                message.addressN[i] = object.addressN[i] >>> 0;
        }
        if (object.prevHash != null)
            if (typeof object.prevHash === "string")
                $util.base64.decode(object.prevHash, message.prevHash = $util.newBuffer($util.base64.length(object.prevHash)), 0);
            else if (object.prevHash.length)
                message.prevHash = object.prevHash;
        if (object.prevIndex != null)
            message.prevIndex = object.prevIndex >>> 0;
        if (object.scriptSig != null)
            if (typeof object.scriptSig === "string")
                $util.base64.decode(object.scriptSig, message.scriptSig = $util.newBuffer($util.base64.length(object.scriptSig)), 0);
            else if (object.scriptSig.length)
                message.scriptSig = object.scriptSig;
        if (object.sequence != null)
            message.sequence = object.sequence >>> 0;
        switch (object.scriptType) {
        case "SPENDADDRESS":
        case 0:
            message.scriptType = 0;
            break;
        case "SPENDMULTISIG":
        case 1:
            message.scriptType = 1;
            break;
        case "EXTERNAL":
        case 2:
            message.scriptType = 2;
            break;
        case "SPENDWITNESS":
        case 3:
            message.scriptType = 3;
            break;
        case "SPENDP2SHWITNESS":
        case 4:
            message.scriptType = 4;
            break;
        }
        if (object.multisig != null) {
            if (typeof object.multisig !== "object")
                throw TypeError(".TxInputType.multisig: object expected");
            message.multisig = $root.MultisigRedeemScriptType.fromObject(object.multisig);
        }
        if (object.amount != null)
            if ($util.Long)
                (message.amount = $util.Long.fromValue(object.amount)).unsigned = true;
            else if (typeof object.amount === "string")
                message.amount = parseInt(object.amount, 10);
            else if (typeof object.amount === "number")
                message.amount = object.amount;
            else if (typeof object.amount === "object")
                message.amount = new $util.LongBits(object.amount.low >>> 0, object.amount.high >>> 0).toNumber(true);
        if (object.decredTree != null)
            message.decredTree = object.decredTree >>> 0;
        if (object.decredScriptVersion != null)
            message.decredScriptVersion = object.decredScriptVersion >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a TxInputType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TxInputType
     * @static
     * @param {TxInputType} message TxInputType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TxInputType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults)
            object.addressN = [];
        if (options.defaults) {
            if (options.bytes === String)
                object.prevHash = "";
            else {
                object.prevHash = [];
                if (options.bytes !== Array)
                    object.prevHash = $util.newBuffer(object.prevHash);
            }
            object.prevIndex = 0;
            if (options.bytes === String)
                object.scriptSig = "";
            else {
                object.scriptSig = [];
                if (options.bytes !== Array)
                    object.scriptSig = $util.newBuffer(object.scriptSig);
            }
            object.sequence = 4294967295;
            object.scriptType = options.enums === String ? "SPENDADDRESS" : 0;
            object.multisig = null;
            if ($util.Long) {
                var long = new $util.Long(0, 0, true);
                object.amount = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.amount = options.longs === String ? "0" : 0;
            object.decredTree = 0;
            object.decredScriptVersion = 0;
        }
        if (message.addressN && message.addressN.length) {
            object.addressN = [];
            for (var j = 0; j < message.addressN.length; ++j)
                object.addressN[j] = message.addressN[j];
        }
        if (message.prevHash != null && message.hasOwnProperty("prevHash"))
            object.prevHash = options.bytes === String ? $util.base64.encode(message.prevHash, 0, message.prevHash.length) : options.bytes === Array ? Array.prototype.slice.call(message.prevHash) : message.prevHash;
        if (message.prevIndex != null && message.hasOwnProperty("prevIndex"))
            object.prevIndex = message.prevIndex;
        if (message.scriptSig != null && message.hasOwnProperty("scriptSig"))
            object.scriptSig = options.bytes === String ? $util.base64.encode(message.scriptSig, 0, message.scriptSig.length) : options.bytes === Array ? Array.prototype.slice.call(message.scriptSig) : message.scriptSig;
        if (message.sequence != null && message.hasOwnProperty("sequence"))
            object.sequence = message.sequence;
        if (message.scriptType != null && message.hasOwnProperty("scriptType"))
            object.scriptType = options.enums === String ? $root.InputScriptType[message.scriptType] : message.scriptType;
        if (message.multisig != null && message.hasOwnProperty("multisig"))
            object.multisig = $root.MultisigRedeemScriptType.toObject(message.multisig, options);
        if (message.amount != null && message.hasOwnProperty("amount"))
            if (typeof message.amount === "number")
                object.amount = options.longs === String ? String(message.amount) : message.amount;
            else
                object.amount = options.longs === String ? $util.Long.prototype.toString.call(message.amount) : options.longs === Number ? new $util.LongBits(message.amount.low >>> 0, message.amount.high >>> 0).toNumber(true) : message.amount;
        if (message.decredTree != null && message.hasOwnProperty("decredTree"))
            object.decredTree = message.decredTree;
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            object.decredScriptVersion = message.decredScriptVersion;
        return object;
    };

    /**
     * Converts this TxInputType to JSON.
     * @function toJSON
     * @memberof TxInputType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TxInputType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TxInputType;
})();

$root.TxOutputType = (function() {

    /**
     * Properties of a TxOutputType.
     * @exports ITxOutputType
     * @interface ITxOutputType
     * @property {string|null} [address] TxOutputType address
     * @property {Array.<number>|null} [addressN] TxOutputType addressN
     * @property {number|Long} amount TxOutputType amount
     * @property {OutputScriptType} scriptType TxOutputType scriptType
     * @property {IMultisigRedeemScriptType|null} [multisig] TxOutputType multisig
     * @property {Uint8Array|null} [opReturnData] TxOutputType opReturnData
     * @property {number|null} [decredScriptVersion] TxOutputType decredScriptVersion
     */

    /**
     * Constructs a new TxOutputType.
     * @exports TxOutputType
     * @classdesc Structure representing transaction output
     * @used_in SimpleSignTx
     * @used_in TransactionType
     * @implements ITxOutputType
     * @constructor
     * @param {ITxOutputType=} [properties] Properties to set
     */
    function TxOutputType(properties) {
        this.addressN = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TxOutputType address.
     * @member {string} address
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.address = "";

    /**
     * TxOutputType addressN.
     * @member {Array.<number>} addressN
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.addressN = $util.emptyArray;

    /**
     * TxOutputType amount.
     * @member {number|Long} amount
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.amount = $util.Long ? $util.Long.fromBits(0,0,true) : 0;

    /**
     * TxOutputType scriptType.
     * @member {OutputScriptType} scriptType
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.scriptType = 0;

    /**
     * TxOutputType multisig.
     * @member {IMultisigRedeemScriptType|null|undefined} multisig
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.multisig = null;

    /**
     * TxOutputType opReturnData.
     * @member {Uint8Array} opReturnData
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.opReturnData = $util.newBuffer([]);

    /**
     * TxOutputType decredScriptVersion.
     * @member {number} decredScriptVersion
     * @memberof TxOutputType
     * @instance
     */
    TxOutputType.prototype.decredScriptVersion = 0;

    /**
     * Creates a new TxOutputType instance using the specified properties.
     * @function create
     * @memberof TxOutputType
     * @static
     * @param {ITxOutputType=} [properties] Properties to set
     * @returns {TxOutputType} TxOutputType instance
     */
    TxOutputType.create = function create(properties) {
        return new TxOutputType(properties);
    };

    /**
     * Encodes the specified TxOutputType message. Does not implicitly {@link TxOutputType.verify|verify} messages.
     * @function encode
     * @memberof TxOutputType
     * @static
     * @param {ITxOutputType} message TxOutputType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxOutputType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.address != null && message.hasOwnProperty("address"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.address);
        if (message.addressN != null && message.addressN.length)
            for (var i = 0; i < message.addressN.length; ++i)
                writer.uint32(/* id 2, wireType 0 =*/16).uint32(message.addressN[i]);
        writer.uint32(/* id 3, wireType 0 =*/24).uint64(message.amount);
        writer.uint32(/* id 4, wireType 0 =*/32).int32(message.scriptType);
        if (message.multisig != null && message.hasOwnProperty("multisig"))
            $root.MultisigRedeemScriptType.encode(message.multisig, writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        if (message.opReturnData != null && message.hasOwnProperty("opReturnData"))
            writer.uint32(/* id 6, wireType 2 =*/50).bytes(message.opReturnData);
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            writer.uint32(/* id 7, wireType 0 =*/56).uint32(message.decredScriptVersion);
        return writer;
    };

    /**
     * Encodes the specified TxOutputType message, length delimited. Does not implicitly {@link TxOutputType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TxOutputType
     * @static
     * @param {ITxOutputType} message TxOutputType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxOutputType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TxOutputType message from the specified reader or buffer.
     * @function decode
     * @memberof TxOutputType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TxOutputType} TxOutputType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxOutputType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.TxOutputType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.address = reader.string();
                break;
            case 2:
                if (!(message.addressN && message.addressN.length))
                    message.addressN = [];
                if ((tag & 7) === 2) {
                    var end2 = reader.uint32() + reader.pos;
                    while (reader.pos < end2)
                        message.addressN.push(reader.uint32());
                } else
                    message.addressN.push(reader.uint32());
                break;
            case 3:
                message.amount = reader.uint64();
                break;
            case 4:
                message.scriptType = reader.int32();
                break;
            case 5:
                message.multisig = $root.MultisigRedeemScriptType.decode(reader, reader.uint32());
                break;
            case 6:
                message.opReturnData = reader.bytes();
                break;
            case 7:
                message.decredScriptVersion = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("amount"))
            throw $util.ProtocolError("missing required 'amount'", { instance: message });
        if (!message.hasOwnProperty("scriptType"))
            throw $util.ProtocolError("missing required 'scriptType'", { instance: message });
        return message;
    };

    /**
     * Decodes a TxOutputType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TxOutputType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TxOutputType} TxOutputType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxOutputType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TxOutputType message.
     * @function verify
     * @memberof TxOutputType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TxOutputType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.address != null && message.hasOwnProperty("address"))
            if (!$util.isString(message.address))
                return "address: string expected";
        if (message.addressN != null && message.hasOwnProperty("addressN")) {
            if (!Array.isArray(message.addressN))
                return "addressN: array expected";
            for (var i = 0; i < message.addressN.length; ++i)
                if (!$util.isInteger(message.addressN[i]))
                    return "addressN: integer[] expected";
        }
        if (!$util.isInteger(message.amount) && !(message.amount && $util.isInteger(message.amount.low) && $util.isInteger(message.amount.high)))
            return "amount: integer|Long expected";
        switch (message.scriptType) {
        default:
            return "scriptType: enum value expected";
        case 0:
        case 1:
        case 2:
        case 3:
        case 4:
        case 5:
            break;
        }
        if (message.multisig != null && message.hasOwnProperty("multisig")) {
            var error = $root.MultisigRedeemScriptType.verify(message.multisig);
            if (error)
                return "multisig." + error;
        }
        if (message.opReturnData != null && message.hasOwnProperty("opReturnData"))
            if (!(message.opReturnData && typeof message.opReturnData.length === "number" || $util.isString(message.opReturnData)))
                return "opReturnData: buffer expected";
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            if (!$util.isInteger(message.decredScriptVersion))
                return "decredScriptVersion: integer expected";
        return null;
    };

    /**
     * Creates a TxOutputType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TxOutputType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TxOutputType} TxOutputType
     */
    TxOutputType.fromObject = function fromObject(object) {
        if (object instanceof $root.TxOutputType)
            return object;
        var message = new $root.TxOutputType();
        if (object.address != null)
            message.address = String(object.address);
        if (object.addressN) {
            if (!Array.isArray(object.addressN))
                throw TypeError(".TxOutputType.addressN: array expected");
            message.addressN = [];
            for (var i = 0; i < object.addressN.length; ++i)
                message.addressN[i] = object.addressN[i] >>> 0;
        }
        if (object.amount != null)
            if ($util.Long)
                (message.amount = $util.Long.fromValue(object.amount)).unsigned = true;
            else if (typeof object.amount === "string")
                message.amount = parseInt(object.amount, 10);
            else if (typeof object.amount === "number")
                message.amount = object.amount;
            else if (typeof object.amount === "object")
                message.amount = new $util.LongBits(object.amount.low >>> 0, object.amount.high >>> 0).toNumber(true);
        switch (object.scriptType) {
        case "PAYTOADDRESS":
        case 0:
            message.scriptType = 0;
            break;
        case "PAYTOSCRIPTHASH":
        case 1:
            message.scriptType = 1;
            break;
        case "PAYTOMULTISIG":
        case 2:
            message.scriptType = 2;
            break;
        case "PAYTOOPRETURN":
        case 3:
            message.scriptType = 3;
            break;
        case "PAYTOWITNESS":
        case 4:
            message.scriptType = 4;
            break;
        case "PAYTOP2SHWITNESS":
        case 5:
            message.scriptType = 5;
            break;
        }
        if (object.multisig != null) {
            if (typeof object.multisig !== "object")
                throw TypeError(".TxOutputType.multisig: object expected");
            message.multisig = $root.MultisigRedeemScriptType.fromObject(object.multisig);
        }
        if (object.opReturnData != null)
            if (typeof object.opReturnData === "string")
                $util.base64.decode(object.opReturnData, message.opReturnData = $util.newBuffer($util.base64.length(object.opReturnData)), 0);
            else if (object.opReturnData.length)
                message.opReturnData = object.opReturnData;
        if (object.decredScriptVersion != null)
            message.decredScriptVersion = object.decredScriptVersion >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a TxOutputType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TxOutputType
     * @static
     * @param {TxOutputType} message TxOutputType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TxOutputType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults)
            object.addressN = [];
        if (options.defaults) {
            object.address = "";
            if ($util.Long) {
                var long = new $util.Long(0, 0, true);
                object.amount = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.amount = options.longs === String ? "0" : 0;
            object.scriptType = options.enums === String ? "PAYTOADDRESS" : 0;
            object.multisig = null;
            if (options.bytes === String)
                object.opReturnData = "";
            else {
                object.opReturnData = [];
                if (options.bytes !== Array)
                    object.opReturnData = $util.newBuffer(object.opReturnData);
            }
            object.decredScriptVersion = 0;
        }
        if (message.address != null && message.hasOwnProperty("address"))
            object.address = message.address;
        if (message.addressN && message.addressN.length) {
            object.addressN = [];
            for (var j = 0; j < message.addressN.length; ++j)
                object.addressN[j] = message.addressN[j];
        }
        if (message.amount != null && message.hasOwnProperty("amount"))
            if (typeof message.amount === "number")
                object.amount = options.longs === String ? String(message.amount) : message.amount;
            else
                object.amount = options.longs === String ? $util.Long.prototype.toString.call(message.amount) : options.longs === Number ? new $util.LongBits(message.amount.low >>> 0, message.amount.high >>> 0).toNumber(true) : message.amount;
        if (message.scriptType != null && message.hasOwnProperty("scriptType"))
            object.scriptType = options.enums === String ? $root.OutputScriptType[message.scriptType] : message.scriptType;
        if (message.multisig != null && message.hasOwnProperty("multisig"))
            object.multisig = $root.MultisigRedeemScriptType.toObject(message.multisig, options);
        if (message.opReturnData != null && message.hasOwnProperty("opReturnData"))
            object.opReturnData = options.bytes === String ? $util.base64.encode(message.opReturnData, 0, message.opReturnData.length) : options.bytes === Array ? Array.prototype.slice.call(message.opReturnData) : message.opReturnData;
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            object.decredScriptVersion = message.decredScriptVersion;
        return object;
    };

    /**
     * Converts this TxOutputType to JSON.
     * @function toJSON
     * @memberof TxOutputType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TxOutputType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TxOutputType;
})();

$root.TxOutputBinType = (function() {

    /**
     * Properties of a TxOutputBinType.
     * @exports ITxOutputBinType
     * @interface ITxOutputBinType
     * @property {number|Long} amount TxOutputBinType amount
     * @property {Uint8Array} scriptPubkey TxOutputBinType scriptPubkey
     * @property {number|null} [decredScriptVersion] TxOutputBinType decredScriptVersion
     */

    /**
     * Constructs a new TxOutputBinType.
     * @exports TxOutputBinType
     * @classdesc Structure representing compiled transaction output
     * @used_in TransactionType
     * @implements ITxOutputBinType
     * @constructor
     * @param {ITxOutputBinType=} [properties] Properties to set
     */
    function TxOutputBinType(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TxOutputBinType amount.
     * @member {number|Long} amount
     * @memberof TxOutputBinType
     * @instance
     */
    TxOutputBinType.prototype.amount = $util.Long ? $util.Long.fromBits(0,0,true) : 0;

    /**
     * TxOutputBinType scriptPubkey.
     * @member {Uint8Array} scriptPubkey
     * @memberof TxOutputBinType
     * @instance
     */
    TxOutputBinType.prototype.scriptPubkey = $util.newBuffer([]);

    /**
     * TxOutputBinType decredScriptVersion.
     * @member {number} decredScriptVersion
     * @memberof TxOutputBinType
     * @instance
     */
    TxOutputBinType.prototype.decredScriptVersion = 0;

    /**
     * Creates a new TxOutputBinType instance using the specified properties.
     * @function create
     * @memberof TxOutputBinType
     * @static
     * @param {ITxOutputBinType=} [properties] Properties to set
     * @returns {TxOutputBinType} TxOutputBinType instance
     */
    TxOutputBinType.create = function create(properties) {
        return new TxOutputBinType(properties);
    };

    /**
     * Encodes the specified TxOutputBinType message. Does not implicitly {@link TxOutputBinType.verify|verify} messages.
     * @function encode
     * @memberof TxOutputBinType
     * @static
     * @param {ITxOutputBinType} message TxOutputBinType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxOutputBinType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        writer.uint32(/* id 1, wireType 0 =*/8).uint64(message.amount);
        writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.scriptPubkey);
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.decredScriptVersion);
        return writer;
    };

    /**
     * Encodes the specified TxOutputBinType message, length delimited. Does not implicitly {@link TxOutputBinType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TxOutputBinType
     * @static
     * @param {ITxOutputBinType} message TxOutputBinType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxOutputBinType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TxOutputBinType message from the specified reader or buffer.
     * @function decode
     * @memberof TxOutputBinType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TxOutputBinType} TxOutputBinType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxOutputBinType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.TxOutputBinType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.amount = reader.uint64();
                break;
            case 2:
                message.scriptPubkey = reader.bytes();
                break;
            case 3:
                message.decredScriptVersion = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        if (!message.hasOwnProperty("amount"))
            throw $util.ProtocolError("missing required 'amount'", { instance: message });
        if (!message.hasOwnProperty("scriptPubkey"))
            throw $util.ProtocolError("missing required 'scriptPubkey'", { instance: message });
        return message;
    };

    /**
     * Decodes a TxOutputBinType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TxOutputBinType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TxOutputBinType} TxOutputBinType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxOutputBinType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TxOutputBinType message.
     * @function verify
     * @memberof TxOutputBinType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TxOutputBinType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (!$util.isInteger(message.amount) && !(message.amount && $util.isInteger(message.amount.low) && $util.isInteger(message.amount.high)))
            return "amount: integer|Long expected";
        if (!(message.scriptPubkey && typeof message.scriptPubkey.length === "number" || $util.isString(message.scriptPubkey)))
            return "scriptPubkey: buffer expected";
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            if (!$util.isInteger(message.decredScriptVersion))
                return "decredScriptVersion: integer expected";
        return null;
    };

    /**
     * Creates a TxOutputBinType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TxOutputBinType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TxOutputBinType} TxOutputBinType
     */
    TxOutputBinType.fromObject = function fromObject(object) {
        if (object instanceof $root.TxOutputBinType)
            return object;
        var message = new $root.TxOutputBinType();
        if (object.amount != null)
            if ($util.Long)
                (message.amount = $util.Long.fromValue(object.amount)).unsigned = true;
            else if (typeof object.amount === "string")
                message.amount = parseInt(object.amount, 10);
            else if (typeof object.amount === "number")
                message.amount = object.amount;
            else if (typeof object.amount === "object")
                message.amount = new $util.LongBits(object.amount.low >>> 0, object.amount.high >>> 0).toNumber(true);
        if (object.scriptPubkey != null)
            if (typeof object.scriptPubkey === "string")
                $util.base64.decode(object.scriptPubkey, message.scriptPubkey = $util.newBuffer($util.base64.length(object.scriptPubkey)), 0);
            else if (object.scriptPubkey.length)
                message.scriptPubkey = object.scriptPubkey;
        if (object.decredScriptVersion != null)
            message.decredScriptVersion = object.decredScriptVersion >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a TxOutputBinType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TxOutputBinType
     * @static
     * @param {TxOutputBinType} message TxOutputBinType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TxOutputBinType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            if ($util.Long) {
                var long = new $util.Long(0, 0, true);
                object.amount = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
            } else
                object.amount = options.longs === String ? "0" : 0;
            if (options.bytes === String)
                object.scriptPubkey = "";
            else {
                object.scriptPubkey = [];
                if (options.bytes !== Array)
                    object.scriptPubkey = $util.newBuffer(object.scriptPubkey);
            }
            object.decredScriptVersion = 0;
        }
        if (message.amount != null && message.hasOwnProperty("amount"))
            if (typeof message.amount === "number")
                object.amount = options.longs === String ? String(message.amount) : message.amount;
            else
                object.amount = options.longs === String ? $util.Long.prototype.toString.call(message.amount) : options.longs === Number ? new $util.LongBits(message.amount.low >>> 0, message.amount.high >>> 0).toNumber(true) : message.amount;
        if (message.scriptPubkey != null && message.hasOwnProperty("scriptPubkey"))
            object.scriptPubkey = options.bytes === String ? $util.base64.encode(message.scriptPubkey, 0, message.scriptPubkey.length) : options.bytes === Array ? Array.prototype.slice.call(message.scriptPubkey) : message.scriptPubkey;
        if (message.decredScriptVersion != null && message.hasOwnProperty("decredScriptVersion"))
            object.decredScriptVersion = message.decredScriptVersion;
        return object;
    };

    /**
     * Converts this TxOutputBinType to JSON.
     * @function toJSON
     * @memberof TxOutputBinType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TxOutputBinType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TxOutputBinType;
})();

$root.TransactionType = (function() {

    /**
     * Properties of a TransactionType.
     * @exports ITransactionType
     * @interface ITransactionType
     * @property {number|null} [version] TransactionType version
     * @property {Array.<ITxInputType>|null} [inputs] TransactionType inputs
     * @property {Array.<ITxOutputBinType>|null} [binOutputs] TransactionType binOutputs
     * @property {Array.<ITxOutputType>|null} [outputs] TransactionType outputs
     * @property {number|null} [lockTime] TransactionType lockTime
     * @property {number|null} [inputsCnt] TransactionType inputsCnt
     * @property {number|null} [outputsCnt] TransactionType outputsCnt
     * @property {Uint8Array|null} [extraData] TransactionType extraData
     * @property {number|null} [extraDataLen] TransactionType extraDataLen
     * @property {number|null} [decredExpiry] TransactionType decredExpiry
     */

    /**
     * Constructs a new TransactionType.
     * @exports TransactionType
     * @classdesc Structure representing transaction
     * @used_in SimpleSignTx
     * @implements ITransactionType
     * @constructor
     * @param {ITransactionType=} [properties] Properties to set
     */
    function TransactionType(properties) {
        this.inputs = [];
        this.binOutputs = [];
        this.outputs = [];
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TransactionType version.
     * @member {number} version
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.version = 0;

    /**
     * TransactionType inputs.
     * @member {Array.<ITxInputType>} inputs
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.inputs = $util.emptyArray;

    /**
     * TransactionType binOutputs.
     * @member {Array.<ITxOutputBinType>} binOutputs
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.binOutputs = $util.emptyArray;

    /**
     * TransactionType outputs.
     * @member {Array.<ITxOutputType>} outputs
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.outputs = $util.emptyArray;

    /**
     * TransactionType lockTime.
     * @member {number} lockTime
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.lockTime = 0;

    /**
     * TransactionType inputsCnt.
     * @member {number} inputsCnt
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.inputsCnt = 0;

    /**
     * TransactionType outputsCnt.
     * @member {number} outputsCnt
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.outputsCnt = 0;

    /**
     * TransactionType extraData.
     * @member {Uint8Array} extraData
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.extraData = $util.newBuffer([]);

    /**
     * TransactionType extraDataLen.
     * @member {number} extraDataLen
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.extraDataLen = 0;

    /**
     * TransactionType decredExpiry.
     * @member {number} decredExpiry
     * @memberof TransactionType
     * @instance
     */
    TransactionType.prototype.decredExpiry = 0;

    /**
     * Creates a new TransactionType instance using the specified properties.
     * @function create
     * @memberof TransactionType
     * @static
     * @param {ITransactionType=} [properties] Properties to set
     * @returns {TransactionType} TransactionType instance
     */
    TransactionType.create = function create(properties) {
        return new TransactionType(properties);
    };

    /**
     * Encodes the specified TransactionType message. Does not implicitly {@link TransactionType.verify|verify} messages.
     * @function encode
     * @memberof TransactionType
     * @static
     * @param {ITransactionType} message TransactionType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TransactionType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.version != null && message.hasOwnProperty("version"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.version);
        if (message.inputs != null && message.inputs.length)
            for (var i = 0; i < message.inputs.length; ++i)
                $root.TxInputType.encode(message.inputs[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
        if (message.binOutputs != null && message.binOutputs.length)
            for (var i = 0; i < message.binOutputs.length; ++i)
                $root.TxOutputBinType.encode(message.binOutputs[i], writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
        if (message.lockTime != null && message.hasOwnProperty("lockTime"))
            writer.uint32(/* id 4, wireType 0 =*/32).uint32(message.lockTime);
        if (message.outputs != null && message.outputs.length)
            for (var i = 0; i < message.outputs.length; ++i)
                $root.TxOutputType.encode(message.outputs[i], writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
        if (message.inputsCnt != null && message.hasOwnProperty("inputsCnt"))
            writer.uint32(/* id 6, wireType 0 =*/48).uint32(message.inputsCnt);
        if (message.outputsCnt != null && message.hasOwnProperty("outputsCnt"))
            writer.uint32(/* id 7, wireType 0 =*/56).uint32(message.outputsCnt);
        if (message.extraData != null && message.hasOwnProperty("extraData"))
            writer.uint32(/* id 8, wireType 2 =*/66).bytes(message.extraData);
        if (message.extraDataLen != null && message.hasOwnProperty("extraDataLen"))
            writer.uint32(/* id 9, wireType 0 =*/72).uint32(message.extraDataLen);
        if (message.decredExpiry != null && message.hasOwnProperty("decredExpiry"))
            writer.uint32(/* id 10, wireType 0 =*/80).uint32(message.decredExpiry);
        return writer;
    };

    /**
     * Encodes the specified TransactionType message, length delimited. Does not implicitly {@link TransactionType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TransactionType
     * @static
     * @param {ITransactionType} message TransactionType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TransactionType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TransactionType message from the specified reader or buffer.
     * @function decode
     * @memberof TransactionType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TransactionType} TransactionType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TransactionType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.TransactionType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.version = reader.uint32();
                break;
            case 2:
                if (!(message.inputs && message.inputs.length))
                    message.inputs = [];
                message.inputs.push($root.TxInputType.decode(reader, reader.uint32()));
                break;
            case 3:
                if (!(message.binOutputs && message.binOutputs.length))
                    message.binOutputs = [];
                message.binOutputs.push($root.TxOutputBinType.decode(reader, reader.uint32()));
                break;
            case 5:
                if (!(message.outputs && message.outputs.length))
                    message.outputs = [];
                message.outputs.push($root.TxOutputType.decode(reader, reader.uint32()));
                break;
            case 4:
                message.lockTime = reader.uint32();
                break;
            case 6:
                message.inputsCnt = reader.uint32();
                break;
            case 7:
                message.outputsCnt = reader.uint32();
                break;
            case 8:
                message.extraData = reader.bytes();
                break;
            case 9:
                message.extraDataLen = reader.uint32();
                break;
            case 10:
                message.decredExpiry = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a TransactionType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TransactionType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TransactionType} TransactionType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TransactionType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TransactionType message.
     * @function verify
     * @memberof TransactionType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TransactionType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.version != null && message.hasOwnProperty("version"))
            if (!$util.isInteger(message.version))
                return "version: integer expected";
        if (message.inputs != null && message.hasOwnProperty("inputs")) {
            if (!Array.isArray(message.inputs))
                return "inputs: array expected";
            for (var i = 0; i < message.inputs.length; ++i) {
                var error = $root.TxInputType.verify(message.inputs[i]);
                if (error)
                    return "inputs." + error;
            }
        }
        if (message.binOutputs != null && message.hasOwnProperty("binOutputs")) {
            if (!Array.isArray(message.binOutputs))
                return "binOutputs: array expected";
            for (var i = 0; i < message.binOutputs.length; ++i) {
                var error = $root.TxOutputBinType.verify(message.binOutputs[i]);
                if (error)
                    return "binOutputs." + error;
            }
        }
        if (message.outputs != null && message.hasOwnProperty("outputs")) {
            if (!Array.isArray(message.outputs))
                return "outputs: array expected";
            for (var i = 0; i < message.outputs.length; ++i) {
                var error = $root.TxOutputType.verify(message.outputs[i]);
                if (error)
                    return "outputs." + error;
            }
        }
        if (message.lockTime != null && message.hasOwnProperty("lockTime"))
            if (!$util.isInteger(message.lockTime))
                return "lockTime: integer expected";
        if (message.inputsCnt != null && message.hasOwnProperty("inputsCnt"))
            if (!$util.isInteger(message.inputsCnt))
                return "inputsCnt: integer expected";
        if (message.outputsCnt != null && message.hasOwnProperty("outputsCnt"))
            if (!$util.isInteger(message.outputsCnt))
                return "outputsCnt: integer expected";
        if (message.extraData != null && message.hasOwnProperty("extraData"))
            if (!(message.extraData && typeof message.extraData.length === "number" || $util.isString(message.extraData)))
                return "extraData: buffer expected";
        if (message.extraDataLen != null && message.hasOwnProperty("extraDataLen"))
            if (!$util.isInteger(message.extraDataLen))
                return "extraDataLen: integer expected";
        if (message.decredExpiry != null && message.hasOwnProperty("decredExpiry"))
            if (!$util.isInteger(message.decredExpiry))
                return "decredExpiry: integer expected";
        return null;
    };

    /**
     * Creates a TransactionType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TransactionType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TransactionType} TransactionType
     */
    TransactionType.fromObject = function fromObject(object) {
        if (object instanceof $root.TransactionType)
            return object;
        var message = new $root.TransactionType();
        if (object.version != null)
            message.version = object.version >>> 0;
        if (object.inputs) {
            if (!Array.isArray(object.inputs))
                throw TypeError(".TransactionType.inputs: array expected");
            message.inputs = [];
            for (var i = 0; i < object.inputs.length; ++i) {
                if (typeof object.inputs[i] !== "object")
                    throw TypeError(".TransactionType.inputs: object expected");
                message.inputs[i] = $root.TxInputType.fromObject(object.inputs[i]);
            }
        }
        if (object.binOutputs) {
            if (!Array.isArray(object.binOutputs))
                throw TypeError(".TransactionType.binOutputs: array expected");
            message.binOutputs = [];
            for (var i = 0; i < object.binOutputs.length; ++i) {
                if (typeof object.binOutputs[i] !== "object")
                    throw TypeError(".TransactionType.binOutputs: object expected");
                message.binOutputs[i] = $root.TxOutputBinType.fromObject(object.binOutputs[i]);
            }
        }
        if (object.outputs) {
            if (!Array.isArray(object.outputs))
                throw TypeError(".TransactionType.outputs: array expected");
            message.outputs = [];
            for (var i = 0; i < object.outputs.length; ++i) {
                if (typeof object.outputs[i] !== "object")
                    throw TypeError(".TransactionType.outputs: object expected");
                message.outputs[i] = $root.TxOutputType.fromObject(object.outputs[i]);
            }
        }
        if (object.lockTime != null)
            message.lockTime = object.lockTime >>> 0;
        if (object.inputsCnt != null)
            message.inputsCnt = object.inputsCnt >>> 0;
        if (object.outputsCnt != null)
            message.outputsCnt = object.outputsCnt >>> 0;
        if (object.extraData != null)
            if (typeof object.extraData === "string")
                $util.base64.decode(object.extraData, message.extraData = $util.newBuffer($util.base64.length(object.extraData)), 0);
            else if (object.extraData.length)
                message.extraData = object.extraData;
        if (object.extraDataLen != null)
            message.extraDataLen = object.extraDataLen >>> 0;
        if (object.decredExpiry != null)
            message.decredExpiry = object.decredExpiry >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a TransactionType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TransactionType
     * @static
     * @param {TransactionType} message TransactionType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TransactionType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.arrays || options.defaults) {
            object.inputs = [];
            object.binOutputs = [];
            object.outputs = [];
        }
        if (options.defaults) {
            object.version = 0;
            object.lockTime = 0;
            object.inputsCnt = 0;
            object.outputsCnt = 0;
            if (options.bytes === String)
                object.extraData = "";
            else {
                object.extraData = [];
                if (options.bytes !== Array)
                    object.extraData = $util.newBuffer(object.extraData);
            }
            object.extraDataLen = 0;
            object.decredExpiry = 0;
        }
        if (message.version != null && message.hasOwnProperty("version"))
            object.version = message.version;
        if (message.inputs && message.inputs.length) {
            object.inputs = [];
            for (var j = 0; j < message.inputs.length; ++j)
                object.inputs[j] = $root.TxInputType.toObject(message.inputs[j], options);
        }
        if (message.binOutputs && message.binOutputs.length) {
            object.binOutputs = [];
            for (var j = 0; j < message.binOutputs.length; ++j)
                object.binOutputs[j] = $root.TxOutputBinType.toObject(message.binOutputs[j], options);
        }
        if (message.lockTime != null && message.hasOwnProperty("lockTime"))
            object.lockTime = message.lockTime;
        if (message.outputs && message.outputs.length) {
            object.outputs = [];
            for (var j = 0; j < message.outputs.length; ++j)
                object.outputs[j] = $root.TxOutputType.toObject(message.outputs[j], options);
        }
        if (message.inputsCnt != null && message.hasOwnProperty("inputsCnt"))
            object.inputsCnt = message.inputsCnt;
        if (message.outputsCnt != null && message.hasOwnProperty("outputsCnt"))
            object.outputsCnt = message.outputsCnt;
        if (message.extraData != null && message.hasOwnProperty("extraData"))
            object.extraData = options.bytes === String ? $util.base64.encode(message.extraData, 0, message.extraData.length) : options.bytes === Array ? Array.prototype.slice.call(message.extraData) : message.extraData;
        if (message.extraDataLen != null && message.hasOwnProperty("extraDataLen"))
            object.extraDataLen = message.extraDataLen;
        if (message.decredExpiry != null && message.hasOwnProperty("decredExpiry"))
            object.decredExpiry = message.decredExpiry;
        return object;
    };

    /**
     * Converts this TransactionType to JSON.
     * @function toJSON
     * @memberof TransactionType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TransactionType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TransactionType;
})();

$root.TxRequestDetailsType = (function() {

    /**
     * Properties of a TxRequestDetailsType.
     * @exports ITxRequestDetailsType
     * @interface ITxRequestDetailsType
     * @property {number|null} [requestIndex] TxRequestDetailsType requestIndex
     * @property {Uint8Array|null} [txHash] TxRequestDetailsType txHash
     * @property {number|null} [extraDataLen] TxRequestDetailsType extraDataLen
     * @property {number|null} [extraDataOffset] TxRequestDetailsType extraDataOffset
     */

    /**
     * Constructs a new TxRequestDetailsType.
     * @exports TxRequestDetailsType
     * @classdesc Structure representing request details
     * @used_in TxRequest
     * @implements ITxRequestDetailsType
     * @constructor
     * @param {ITxRequestDetailsType=} [properties] Properties to set
     */
    function TxRequestDetailsType(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TxRequestDetailsType requestIndex.
     * @member {number} requestIndex
     * @memberof TxRequestDetailsType
     * @instance
     */
    TxRequestDetailsType.prototype.requestIndex = 0;

    /**
     * TxRequestDetailsType txHash.
     * @member {Uint8Array} txHash
     * @memberof TxRequestDetailsType
     * @instance
     */
    TxRequestDetailsType.prototype.txHash = $util.newBuffer([]);

    /**
     * TxRequestDetailsType extraDataLen.
     * @member {number} extraDataLen
     * @memberof TxRequestDetailsType
     * @instance
     */
    TxRequestDetailsType.prototype.extraDataLen = 0;

    /**
     * TxRequestDetailsType extraDataOffset.
     * @member {number} extraDataOffset
     * @memberof TxRequestDetailsType
     * @instance
     */
    TxRequestDetailsType.prototype.extraDataOffset = 0;

    /**
     * Creates a new TxRequestDetailsType instance using the specified properties.
     * @function create
     * @memberof TxRequestDetailsType
     * @static
     * @param {ITxRequestDetailsType=} [properties] Properties to set
     * @returns {TxRequestDetailsType} TxRequestDetailsType instance
     */
    TxRequestDetailsType.create = function create(properties) {
        return new TxRequestDetailsType(properties);
    };

    /**
     * Encodes the specified TxRequestDetailsType message. Does not implicitly {@link TxRequestDetailsType.verify|verify} messages.
     * @function encode
     * @memberof TxRequestDetailsType
     * @static
     * @param {ITxRequestDetailsType} message TxRequestDetailsType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxRequestDetailsType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.requestIndex != null && message.hasOwnProperty("requestIndex"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.requestIndex);
        if (message.txHash != null && message.hasOwnProperty("txHash"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.txHash);
        if (message.extraDataLen != null && message.hasOwnProperty("extraDataLen"))
            writer.uint32(/* id 3, wireType 0 =*/24).uint32(message.extraDataLen);
        if (message.extraDataOffset != null && message.hasOwnProperty("extraDataOffset"))
            writer.uint32(/* id 4, wireType 0 =*/32).uint32(message.extraDataOffset);
        return writer;
    };

    /**
     * Encodes the specified TxRequestDetailsType message, length delimited. Does not implicitly {@link TxRequestDetailsType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TxRequestDetailsType
     * @static
     * @param {ITxRequestDetailsType} message TxRequestDetailsType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxRequestDetailsType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TxRequestDetailsType message from the specified reader or buffer.
     * @function decode
     * @memberof TxRequestDetailsType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TxRequestDetailsType} TxRequestDetailsType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxRequestDetailsType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.TxRequestDetailsType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.requestIndex = reader.uint32();
                break;
            case 2:
                message.txHash = reader.bytes();
                break;
            case 3:
                message.extraDataLen = reader.uint32();
                break;
            case 4:
                message.extraDataOffset = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a TxRequestDetailsType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TxRequestDetailsType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TxRequestDetailsType} TxRequestDetailsType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxRequestDetailsType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TxRequestDetailsType message.
     * @function verify
     * @memberof TxRequestDetailsType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TxRequestDetailsType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.requestIndex != null && message.hasOwnProperty("requestIndex"))
            if (!$util.isInteger(message.requestIndex))
                return "requestIndex: integer expected";
        if (message.txHash != null && message.hasOwnProperty("txHash"))
            if (!(message.txHash && typeof message.txHash.length === "number" || $util.isString(message.txHash)))
                return "txHash: buffer expected";
        if (message.extraDataLen != null && message.hasOwnProperty("extraDataLen"))
            if (!$util.isInteger(message.extraDataLen))
                return "extraDataLen: integer expected";
        if (message.extraDataOffset != null && message.hasOwnProperty("extraDataOffset"))
            if (!$util.isInteger(message.extraDataOffset))
                return "extraDataOffset: integer expected";
        return null;
    };

    /**
     * Creates a TxRequestDetailsType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TxRequestDetailsType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TxRequestDetailsType} TxRequestDetailsType
     */
    TxRequestDetailsType.fromObject = function fromObject(object) {
        if (object instanceof $root.TxRequestDetailsType)
            return object;
        var message = new $root.TxRequestDetailsType();
        if (object.requestIndex != null)
            message.requestIndex = object.requestIndex >>> 0;
        if (object.txHash != null)
            if (typeof object.txHash === "string")
                $util.base64.decode(object.txHash, message.txHash = $util.newBuffer($util.base64.length(object.txHash)), 0);
            else if (object.txHash.length)
                message.txHash = object.txHash;
        if (object.extraDataLen != null)
            message.extraDataLen = object.extraDataLen >>> 0;
        if (object.extraDataOffset != null)
            message.extraDataOffset = object.extraDataOffset >>> 0;
        return message;
    };

    /**
     * Creates a plain object from a TxRequestDetailsType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TxRequestDetailsType
     * @static
     * @param {TxRequestDetailsType} message TxRequestDetailsType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TxRequestDetailsType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.requestIndex = 0;
            if (options.bytes === String)
                object.txHash = "";
            else {
                object.txHash = [];
                if (options.bytes !== Array)
                    object.txHash = $util.newBuffer(object.txHash);
            }
            object.extraDataLen = 0;
            object.extraDataOffset = 0;
        }
        if (message.requestIndex != null && message.hasOwnProperty("requestIndex"))
            object.requestIndex = message.requestIndex;
        if (message.txHash != null && message.hasOwnProperty("txHash"))
            object.txHash = options.bytes === String ? $util.base64.encode(message.txHash, 0, message.txHash.length) : options.bytes === Array ? Array.prototype.slice.call(message.txHash) : message.txHash;
        if (message.extraDataLen != null && message.hasOwnProperty("extraDataLen"))
            object.extraDataLen = message.extraDataLen;
        if (message.extraDataOffset != null && message.hasOwnProperty("extraDataOffset"))
            object.extraDataOffset = message.extraDataOffset;
        return object;
    };

    /**
     * Converts this TxRequestDetailsType to JSON.
     * @function toJSON
     * @memberof TxRequestDetailsType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TxRequestDetailsType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TxRequestDetailsType;
})();

$root.TxRequestSerializedType = (function() {

    /**
     * Properties of a TxRequestSerializedType.
     * @exports ITxRequestSerializedType
     * @interface ITxRequestSerializedType
     * @property {number|null} [signatureIndex] TxRequestSerializedType signatureIndex
     * @property {Uint8Array|null} [signature] TxRequestSerializedType signature
     * @property {Uint8Array|null} [serializedTx] TxRequestSerializedType serializedTx
     */

    /**
     * Constructs a new TxRequestSerializedType.
     * @exports TxRequestSerializedType
     * @classdesc Structure representing serialized data
     * @used_in TxRequest
     * @implements ITxRequestSerializedType
     * @constructor
     * @param {ITxRequestSerializedType=} [properties] Properties to set
     */
    function TxRequestSerializedType(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * TxRequestSerializedType signatureIndex.
     * @member {number} signatureIndex
     * @memberof TxRequestSerializedType
     * @instance
     */
    TxRequestSerializedType.prototype.signatureIndex = 0;

    /**
     * TxRequestSerializedType signature.
     * @member {Uint8Array} signature
     * @memberof TxRequestSerializedType
     * @instance
     */
    TxRequestSerializedType.prototype.signature = $util.newBuffer([]);

    /**
     * TxRequestSerializedType serializedTx.
     * @member {Uint8Array} serializedTx
     * @memberof TxRequestSerializedType
     * @instance
     */
    TxRequestSerializedType.prototype.serializedTx = $util.newBuffer([]);

    /**
     * Creates a new TxRequestSerializedType instance using the specified properties.
     * @function create
     * @memberof TxRequestSerializedType
     * @static
     * @param {ITxRequestSerializedType=} [properties] Properties to set
     * @returns {TxRequestSerializedType} TxRequestSerializedType instance
     */
    TxRequestSerializedType.create = function create(properties) {
        return new TxRequestSerializedType(properties);
    };

    /**
     * Encodes the specified TxRequestSerializedType message. Does not implicitly {@link TxRequestSerializedType.verify|verify} messages.
     * @function encode
     * @memberof TxRequestSerializedType
     * @static
     * @param {ITxRequestSerializedType} message TxRequestSerializedType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxRequestSerializedType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.signatureIndex != null && message.hasOwnProperty("signatureIndex"))
            writer.uint32(/* id 1, wireType 0 =*/8).uint32(message.signatureIndex);
        if (message.signature != null && message.hasOwnProperty("signature"))
            writer.uint32(/* id 2, wireType 2 =*/18).bytes(message.signature);
        if (message.serializedTx != null && message.hasOwnProperty("serializedTx"))
            writer.uint32(/* id 3, wireType 2 =*/26).bytes(message.serializedTx);
        return writer;
    };

    /**
     * Encodes the specified TxRequestSerializedType message, length delimited. Does not implicitly {@link TxRequestSerializedType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof TxRequestSerializedType
     * @static
     * @param {ITxRequestSerializedType} message TxRequestSerializedType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    TxRequestSerializedType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes a TxRequestSerializedType message from the specified reader or buffer.
     * @function decode
     * @memberof TxRequestSerializedType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {TxRequestSerializedType} TxRequestSerializedType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxRequestSerializedType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.TxRequestSerializedType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.signatureIndex = reader.uint32();
                break;
            case 2:
                message.signature = reader.bytes();
                break;
            case 3:
                message.serializedTx = reader.bytes();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes a TxRequestSerializedType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof TxRequestSerializedType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {TxRequestSerializedType} TxRequestSerializedType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    TxRequestSerializedType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies a TxRequestSerializedType message.
     * @function verify
     * @memberof TxRequestSerializedType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    TxRequestSerializedType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.signatureIndex != null && message.hasOwnProperty("signatureIndex"))
            if (!$util.isInteger(message.signatureIndex))
                return "signatureIndex: integer expected";
        if (message.signature != null && message.hasOwnProperty("signature"))
            if (!(message.signature && typeof message.signature.length === "number" || $util.isString(message.signature)))
                return "signature: buffer expected";
        if (message.serializedTx != null && message.hasOwnProperty("serializedTx"))
            if (!(message.serializedTx && typeof message.serializedTx.length === "number" || $util.isString(message.serializedTx)))
                return "serializedTx: buffer expected";
        return null;
    };

    /**
     * Creates a TxRequestSerializedType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof TxRequestSerializedType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {TxRequestSerializedType} TxRequestSerializedType
     */
    TxRequestSerializedType.fromObject = function fromObject(object) {
        if (object instanceof $root.TxRequestSerializedType)
            return object;
        var message = new $root.TxRequestSerializedType();
        if (object.signatureIndex != null)
            message.signatureIndex = object.signatureIndex >>> 0;
        if (object.signature != null)
            if (typeof object.signature === "string")
                $util.base64.decode(object.signature, message.signature = $util.newBuffer($util.base64.length(object.signature)), 0);
            else if (object.signature.length)
                message.signature = object.signature;
        if (object.serializedTx != null)
            if (typeof object.serializedTx === "string")
                $util.base64.decode(object.serializedTx, message.serializedTx = $util.newBuffer($util.base64.length(object.serializedTx)), 0);
            else if (object.serializedTx.length)
                message.serializedTx = object.serializedTx;
        return message;
    };

    /**
     * Creates a plain object from a TxRequestSerializedType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof TxRequestSerializedType
     * @static
     * @param {TxRequestSerializedType} message TxRequestSerializedType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    TxRequestSerializedType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.signatureIndex = 0;
            if (options.bytes === String)
                object.signature = "";
            else {
                object.signature = [];
                if (options.bytes !== Array)
                    object.signature = $util.newBuffer(object.signature);
            }
            if (options.bytes === String)
                object.serializedTx = "";
            else {
                object.serializedTx = [];
                if (options.bytes !== Array)
                    object.serializedTx = $util.newBuffer(object.serializedTx);
            }
        }
        if (message.signatureIndex != null && message.hasOwnProperty("signatureIndex"))
            object.signatureIndex = message.signatureIndex;
        if (message.signature != null && message.hasOwnProperty("signature"))
            object.signature = options.bytes === String ? $util.base64.encode(message.signature, 0, message.signature.length) : options.bytes === Array ? Array.prototype.slice.call(message.signature) : message.signature;
        if (message.serializedTx != null && message.hasOwnProperty("serializedTx"))
            object.serializedTx = options.bytes === String ? $util.base64.encode(message.serializedTx, 0, message.serializedTx.length) : options.bytes === Array ? Array.prototype.slice.call(message.serializedTx) : message.serializedTx;
        return object;
    };

    /**
     * Converts this TxRequestSerializedType to JSON.
     * @function toJSON
     * @memberof TxRequestSerializedType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    TxRequestSerializedType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return TxRequestSerializedType;
})();

$root.IdentityType = (function() {

    /**
     * Properties of an IdentityType.
     * @exports IIdentityType
     * @interface IIdentityType
     * @property {string|null} [proto] IdentityType proto
     * @property {string|null} [user] IdentityType user
     * @property {string|null} [host] IdentityType host
     * @property {string|null} [port] IdentityType port
     * @property {string|null} [path] IdentityType path
     * @property {number|null} [index] IdentityType index
     */

    /**
     * Constructs a new IdentityType.
     * @exports IdentityType
     * @classdesc Structure representing identity data
     * @used_in IdentityType
     * @implements IIdentityType
     * @constructor
     * @param {IIdentityType=} [properties] Properties to set
     */
    function IdentityType(properties) {
        if (properties)
            for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                if (properties[keys[i]] != null)
                    this[keys[i]] = properties[keys[i]];
    }

    /**
     * IdentityType proto.
     * @member {string} proto
     * @memberof IdentityType
     * @instance
     */
    IdentityType.prototype.proto = "";

    /**
     * IdentityType user.
     * @member {string} user
     * @memberof IdentityType
     * @instance
     */
    IdentityType.prototype.user = "";

    /**
     * IdentityType host.
     * @member {string} host
     * @memberof IdentityType
     * @instance
     */
    IdentityType.prototype.host = "";

    /**
     * IdentityType port.
     * @member {string} port
     * @memberof IdentityType
     * @instance
     */
    IdentityType.prototype.port = "";

    /**
     * IdentityType path.
     * @member {string} path
     * @memberof IdentityType
     * @instance
     */
    IdentityType.prototype.path = "";

    /**
     * IdentityType index.
     * @member {number} index
     * @memberof IdentityType
     * @instance
     */
    IdentityType.prototype.index = 0;

    /**
     * Creates a new IdentityType instance using the specified properties.
     * @function create
     * @memberof IdentityType
     * @static
     * @param {IIdentityType=} [properties] Properties to set
     * @returns {IdentityType} IdentityType instance
     */
    IdentityType.create = function create(properties) {
        return new IdentityType(properties);
    };

    /**
     * Encodes the specified IdentityType message. Does not implicitly {@link IdentityType.verify|verify} messages.
     * @function encode
     * @memberof IdentityType
     * @static
     * @param {IIdentityType} message IdentityType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    IdentityType.encode = function encode(message, writer) {
        if (!writer)
            writer = $Writer.create();
        if (message.proto != null && message.hasOwnProperty("proto"))
            writer.uint32(/* id 1, wireType 2 =*/10).string(message.proto);
        if (message.user != null && message.hasOwnProperty("user"))
            writer.uint32(/* id 2, wireType 2 =*/18).string(message.user);
        if (message.host != null && message.hasOwnProperty("host"))
            writer.uint32(/* id 3, wireType 2 =*/26).string(message.host);
        if (message.port != null && message.hasOwnProperty("port"))
            writer.uint32(/* id 4, wireType 2 =*/34).string(message.port);
        if (message.path != null && message.hasOwnProperty("path"))
            writer.uint32(/* id 5, wireType 2 =*/42).string(message.path);
        if (message.index != null && message.hasOwnProperty("index"))
            writer.uint32(/* id 6, wireType 0 =*/48).uint32(message.index);
        return writer;
    };

    /**
     * Encodes the specified IdentityType message, length delimited. Does not implicitly {@link IdentityType.verify|verify} messages.
     * @function encodeDelimited
     * @memberof IdentityType
     * @static
     * @param {IIdentityType} message IdentityType message or plain object to encode
     * @param {$protobuf.Writer} [writer] Writer to encode to
     * @returns {$protobuf.Writer} Writer
     */
    IdentityType.encodeDelimited = function encodeDelimited(message, writer) {
        return this.encode(message, writer).ldelim();
    };

    /**
     * Decodes an IdentityType message from the specified reader or buffer.
     * @function decode
     * @memberof IdentityType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @param {number} [length] Message length if known beforehand
     * @returns {IdentityType} IdentityType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    IdentityType.decode = function decode(reader, length) {
        if (!(reader instanceof $Reader))
            reader = $Reader.create(reader);
        var end = length === undefined ? reader.len : reader.pos + length, message = new $root.IdentityType();
        while (reader.pos < end) {
            var tag = reader.uint32();
            switch (tag >>> 3) {
            case 1:
                message.proto = reader.string();
                break;
            case 2:
                message.user = reader.string();
                break;
            case 3:
                message.host = reader.string();
                break;
            case 4:
                message.port = reader.string();
                break;
            case 5:
                message.path = reader.string();
                break;
            case 6:
                message.index = reader.uint32();
                break;
            default:
                reader.skipType(tag & 7);
                break;
            }
        }
        return message;
    };

    /**
     * Decodes an IdentityType message from the specified reader or buffer, length delimited.
     * @function decodeDelimited
     * @memberof IdentityType
     * @static
     * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
     * @returns {IdentityType} IdentityType
     * @throws {Error} If the payload is not a reader or valid buffer
     * @throws {$protobuf.util.ProtocolError} If required fields are missing
     */
    IdentityType.decodeDelimited = function decodeDelimited(reader) {
        if (!(reader instanceof $Reader))
            reader = new $Reader(reader);
        return this.decode(reader, reader.uint32());
    };

    /**
     * Verifies an IdentityType message.
     * @function verify
     * @memberof IdentityType
     * @static
     * @param {Object.<string,*>} message Plain object to verify
     * @returns {string|null} `null` if valid, otherwise the reason why it is not
     */
    IdentityType.verify = function verify(message) {
        if (typeof message !== "object" || message === null)
            return "object expected";
        if (message.proto != null && message.hasOwnProperty("proto"))
            if (!$util.isString(message.proto))
                return "proto: string expected";
        if (message.user != null && message.hasOwnProperty("user"))
            if (!$util.isString(message.user))
                return "user: string expected";
        if (message.host != null && message.hasOwnProperty("host"))
            if (!$util.isString(message.host))
                return "host: string expected";
        if (message.port != null && message.hasOwnProperty("port"))
            if (!$util.isString(message.port))
                return "port: string expected";
        if (message.path != null && message.hasOwnProperty("path"))
            if (!$util.isString(message.path))
                return "path: string expected";
        if (message.index != null && message.hasOwnProperty("index"))
            if (!$util.isInteger(message.index))
                return "index: integer expected";
        return null;
    };

    /**
     * Creates an IdentityType message from a plain object. Also converts values to their respective internal types.
     * @function fromObject
     * @memberof IdentityType
     * @static
     * @param {Object.<string,*>} object Plain object
     * @returns {IdentityType} IdentityType
     */
    IdentityType.fromObject = function fromObject(object) {
        if (object instanceof $root.IdentityType)
            return object;
        var message = new $root.IdentityType();
        if (object.proto != null)
            message.proto = String(object.proto);
        if (object.user != null)
            message.user = String(object.user);
        if (object.host != null)
            message.host = String(object.host);
        if (object.port != null)
            message.port = String(object.port);
        if (object.path != null)
            message.path = String(object.path);
        if (object.index != null)
            message.index = object.index >>> 0;
        return message;
    };

    /**
     * Creates a plain object from an IdentityType message. Also converts values to other types if specified.
     * @function toObject
     * @memberof IdentityType
     * @static
     * @param {IdentityType} message IdentityType
     * @param {$protobuf.IConversionOptions} [options] Conversion options
     * @returns {Object.<string,*>} Plain object
     */
    IdentityType.toObject = function toObject(message, options) {
        if (!options)
            options = {};
        var object = {};
        if (options.defaults) {
            object.proto = "";
            object.user = "";
            object.host = "";
            object.port = "";
            object.path = "";
            object.index = 0;
        }
        if (message.proto != null && message.hasOwnProperty("proto"))
            object.proto = message.proto;
        if (message.user != null && message.hasOwnProperty("user"))
            object.user = message.user;
        if (message.host != null && message.hasOwnProperty("host"))
            object.host = message.host;
        if (message.port != null && message.hasOwnProperty("port"))
            object.port = message.port;
        if (message.path != null && message.hasOwnProperty("path"))
            object.path = message.path;
        if (message.index != null && message.hasOwnProperty("index"))
            object.index = message.index;
        return object;
    };

    /**
     * Converts this IdentityType to JSON.
     * @function toJSON
     * @memberof IdentityType
     * @instance
     * @returns {Object.<string,*>} JSON object
     */
    IdentityType.prototype.toJSON = function toJSON() {
        return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
    };

    return IdentityType;
})();

/**
 * Ask trezor to generate a skycoin address
 * @used_in SkycoinAddress｛
 * @exports SkycoinAddressType
 * @enum {string}
 * @property {number} AddressTypeSkycoin=1 AddressTypeSkycoin value
 * @property {number} AddressTypeBitcoin=2 AddressTypeBitcoin value
 */
$root.SkycoinAddressType = (function() {
    var valuesById = {}, values = Object.create(valuesById);
    values[valuesById[1] = "AddressTypeSkycoin"] = 1;
    values[valuesById[2] = "AddressTypeBitcoin"] = 2;
    return values;
})();

$root.google = (function() {

    /**
     * Namespace google.
     * @exports google
     * @namespace
     */
    var google = {};

    google.protobuf = (function() {

        /**
         * Namespace protobuf.
         * @memberof google
         * @namespace
         */
        var protobuf = {};

        protobuf.FileDescriptorSet = (function() {

            /**
             * Properties of a FileDescriptorSet.
             * @memberof google.protobuf
             * @interface IFileDescriptorSet
             * @property {Array.<google.protobuf.IFileDescriptorProto>|null} [file] FileDescriptorSet file
             */

            /**
             * Constructs a new FileDescriptorSet.
             * @memberof google.protobuf
             * @classdesc Represents a FileDescriptorSet.
             * @implements IFileDescriptorSet
             * @constructor
             * @param {google.protobuf.IFileDescriptorSet=} [properties] Properties to set
             */
            function FileDescriptorSet(properties) {
                this.file = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * FileDescriptorSet file.
             * @member {Array.<google.protobuf.IFileDescriptorProto>} file
             * @memberof google.protobuf.FileDescriptorSet
             * @instance
             */
            FileDescriptorSet.prototype.file = $util.emptyArray;

            /**
             * Creates a new FileDescriptorSet instance using the specified properties.
             * @function create
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {google.protobuf.IFileDescriptorSet=} [properties] Properties to set
             * @returns {google.protobuf.FileDescriptorSet} FileDescriptorSet instance
             */
            FileDescriptorSet.create = function create(properties) {
                return new FileDescriptorSet(properties);
            };

            /**
             * Encodes the specified FileDescriptorSet message. Does not implicitly {@link google.protobuf.FileDescriptorSet.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {google.protobuf.IFileDescriptorSet} message FileDescriptorSet message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FileDescriptorSet.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.file != null && message.file.length)
                    for (var i = 0; i < message.file.length; ++i)
                        $root.google.protobuf.FileDescriptorProto.encode(message.file[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified FileDescriptorSet message, length delimited. Does not implicitly {@link google.protobuf.FileDescriptorSet.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {google.protobuf.IFileDescriptorSet} message FileDescriptorSet message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FileDescriptorSet.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a FileDescriptorSet message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.FileDescriptorSet} FileDescriptorSet
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FileDescriptorSet.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.FileDescriptorSet();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        if (!(message.file && message.file.length))
                            message.file = [];
                        message.file.push($root.google.protobuf.FileDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a FileDescriptorSet message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.FileDescriptorSet} FileDescriptorSet
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FileDescriptorSet.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a FileDescriptorSet message.
             * @function verify
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            FileDescriptorSet.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.file != null && message.hasOwnProperty("file")) {
                    if (!Array.isArray(message.file))
                        return "file: array expected";
                    for (var i = 0; i < message.file.length; ++i) {
                        var error = $root.google.protobuf.FileDescriptorProto.verify(message.file[i]);
                        if (error)
                            return "file." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a FileDescriptorSet message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.FileDescriptorSet} FileDescriptorSet
             */
            FileDescriptorSet.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.FileDescriptorSet)
                    return object;
                var message = new $root.google.protobuf.FileDescriptorSet();
                if (object.file) {
                    if (!Array.isArray(object.file))
                        throw TypeError(".google.protobuf.FileDescriptorSet.file: array expected");
                    message.file = [];
                    for (var i = 0; i < object.file.length; ++i) {
                        if (typeof object.file[i] !== "object")
                            throw TypeError(".google.protobuf.FileDescriptorSet.file: object expected");
                        message.file[i] = $root.google.protobuf.FileDescriptorProto.fromObject(object.file[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a FileDescriptorSet message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.FileDescriptorSet
             * @static
             * @param {google.protobuf.FileDescriptorSet} message FileDescriptorSet
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            FileDescriptorSet.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.file = [];
                if (message.file && message.file.length) {
                    object.file = [];
                    for (var j = 0; j < message.file.length; ++j)
                        object.file[j] = $root.google.protobuf.FileDescriptorProto.toObject(message.file[j], options);
                }
                return object;
            };

            /**
             * Converts this FileDescriptorSet to JSON.
             * @function toJSON
             * @memberof google.protobuf.FileDescriptorSet
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            FileDescriptorSet.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return FileDescriptorSet;
        })();

        protobuf.FileDescriptorProto = (function() {

            /**
             * Properties of a FileDescriptorProto.
             * @memberof google.protobuf
             * @interface IFileDescriptorProto
             * @property {string|null} [name] FileDescriptorProto name
             * @property {string|null} ["package"] FileDescriptorProto package
             * @property {Array.<string>|null} [dependency] FileDescriptorProto dependency
             * @property {Array.<number>|null} [publicDependency] FileDescriptorProto publicDependency
             * @property {Array.<number>|null} [weakDependency] FileDescriptorProto weakDependency
             * @property {Array.<google.protobuf.IDescriptorProto>|null} [messageType] FileDescriptorProto messageType
             * @property {Array.<google.protobuf.IEnumDescriptorProto>|null} [enumType] FileDescriptorProto enumType
             * @property {Array.<google.protobuf.IServiceDescriptorProto>|null} [service] FileDescriptorProto service
             * @property {Array.<google.protobuf.IFieldDescriptorProto>|null} [extension] FileDescriptorProto extension
             * @property {google.protobuf.IFileOptions|null} [options] FileDescriptorProto options
             * @property {google.protobuf.ISourceCodeInfo|null} [sourceCodeInfo] FileDescriptorProto sourceCodeInfo
             * @property {string|null} [syntax] FileDescriptorProto syntax
             */

            /**
             * Constructs a new FileDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents a FileDescriptorProto.
             * @implements IFileDescriptorProto
             * @constructor
             * @param {google.protobuf.IFileDescriptorProto=} [properties] Properties to set
             */
            function FileDescriptorProto(properties) {
                this.dependency = [];
                this.publicDependency = [];
                this.weakDependency = [];
                this.messageType = [];
                this.enumType = [];
                this.service = [];
                this.extension = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * FileDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.name = "";

            /**
             * FileDescriptorProto package.
             * @member {string} package
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype["package"] = "";

            /**
             * FileDescriptorProto dependency.
             * @member {Array.<string>} dependency
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.dependency = $util.emptyArray;

            /**
             * FileDescriptorProto publicDependency.
             * @member {Array.<number>} publicDependency
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.publicDependency = $util.emptyArray;

            /**
             * FileDescriptorProto weakDependency.
             * @member {Array.<number>} weakDependency
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.weakDependency = $util.emptyArray;

            /**
             * FileDescriptorProto messageType.
             * @member {Array.<google.protobuf.IDescriptorProto>} messageType
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.messageType = $util.emptyArray;

            /**
             * FileDescriptorProto enumType.
             * @member {Array.<google.protobuf.IEnumDescriptorProto>} enumType
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.enumType = $util.emptyArray;

            /**
             * FileDescriptorProto service.
             * @member {Array.<google.protobuf.IServiceDescriptorProto>} service
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.service = $util.emptyArray;

            /**
             * FileDescriptorProto extension.
             * @member {Array.<google.protobuf.IFieldDescriptorProto>} extension
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.extension = $util.emptyArray;

            /**
             * FileDescriptorProto options.
             * @member {google.protobuf.IFileOptions|null|undefined} options
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.options = null;

            /**
             * FileDescriptorProto sourceCodeInfo.
             * @member {google.protobuf.ISourceCodeInfo|null|undefined} sourceCodeInfo
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.sourceCodeInfo = null;

            /**
             * FileDescriptorProto syntax.
             * @member {string} syntax
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             */
            FileDescriptorProto.prototype.syntax = "";

            /**
             * Creates a new FileDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {google.protobuf.IFileDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.FileDescriptorProto} FileDescriptorProto instance
             */
            FileDescriptorProto.create = function create(properties) {
                return new FileDescriptorProto(properties);
            };

            /**
             * Encodes the specified FileDescriptorProto message. Does not implicitly {@link google.protobuf.FileDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {google.protobuf.IFileDescriptorProto} message FileDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FileDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message["package"] != null && message.hasOwnProperty("package"))
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message["package"]);
                if (message.dependency != null && message.dependency.length)
                    for (var i = 0; i < message.dependency.length; ++i)
                        writer.uint32(/* id 3, wireType 2 =*/26).string(message.dependency[i]);
                if (message.messageType != null && message.messageType.length)
                    for (var i = 0; i < message.messageType.length; ++i)
                        $root.google.protobuf.DescriptorProto.encode(message.messageType[i], writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
                if (message.enumType != null && message.enumType.length)
                    for (var i = 0; i < message.enumType.length; ++i)
                        $root.google.protobuf.EnumDescriptorProto.encode(message.enumType[i], writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
                if (message.service != null && message.service.length)
                    for (var i = 0; i < message.service.length; ++i)
                        $root.google.protobuf.ServiceDescriptorProto.encode(message.service[i], writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
                if (message.extension != null && message.extension.length)
                    for (var i = 0; i < message.extension.length; ++i)
                        $root.google.protobuf.FieldDescriptorProto.encode(message.extension[i], writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.FileOptions.encode(message.options, writer.uint32(/* id 8, wireType 2 =*/66).fork()).ldelim();
                if (message.sourceCodeInfo != null && message.hasOwnProperty("sourceCodeInfo"))
                    $root.google.protobuf.SourceCodeInfo.encode(message.sourceCodeInfo, writer.uint32(/* id 9, wireType 2 =*/74).fork()).ldelim();
                if (message.publicDependency != null && message.publicDependency.length)
                    for (var i = 0; i < message.publicDependency.length; ++i)
                        writer.uint32(/* id 10, wireType 0 =*/80).int32(message.publicDependency[i]);
                if (message.weakDependency != null && message.weakDependency.length)
                    for (var i = 0; i < message.weakDependency.length; ++i)
                        writer.uint32(/* id 11, wireType 0 =*/88).int32(message.weakDependency[i]);
                if (message.syntax != null && message.hasOwnProperty("syntax"))
                    writer.uint32(/* id 12, wireType 2 =*/98).string(message.syntax);
                return writer;
            };

            /**
             * Encodes the specified FileDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.FileDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {google.protobuf.IFileDescriptorProto} message FileDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FileDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a FileDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.FileDescriptorProto} FileDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FileDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.FileDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 2:
                        message["package"] = reader.string();
                        break;
                    case 3:
                        if (!(message.dependency && message.dependency.length))
                            message.dependency = [];
                        message.dependency.push(reader.string());
                        break;
                    case 10:
                        if (!(message.publicDependency && message.publicDependency.length))
                            message.publicDependency = [];
                        if ((tag & 7) === 2) {
                            var end2 = reader.uint32() + reader.pos;
                            while (reader.pos < end2)
                                message.publicDependency.push(reader.int32());
                        } else
                            message.publicDependency.push(reader.int32());
                        break;
                    case 11:
                        if (!(message.weakDependency && message.weakDependency.length))
                            message.weakDependency = [];
                        if ((tag & 7) === 2) {
                            var end2 = reader.uint32() + reader.pos;
                            while (reader.pos < end2)
                                message.weakDependency.push(reader.int32());
                        } else
                            message.weakDependency.push(reader.int32());
                        break;
                    case 4:
                        if (!(message.messageType && message.messageType.length))
                            message.messageType = [];
                        message.messageType.push($root.google.protobuf.DescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 5:
                        if (!(message.enumType && message.enumType.length))
                            message.enumType = [];
                        message.enumType.push($root.google.protobuf.EnumDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 6:
                        if (!(message.service && message.service.length))
                            message.service = [];
                        message.service.push($root.google.protobuf.ServiceDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 7:
                        if (!(message.extension && message.extension.length))
                            message.extension = [];
                        message.extension.push($root.google.protobuf.FieldDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 8:
                        message.options = $root.google.protobuf.FileOptions.decode(reader, reader.uint32());
                        break;
                    case 9:
                        message.sourceCodeInfo = $root.google.protobuf.SourceCodeInfo.decode(reader, reader.uint32());
                        break;
                    case 12:
                        message.syntax = reader.string();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a FileDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.FileDescriptorProto} FileDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FileDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a FileDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            FileDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message["package"] != null && message.hasOwnProperty("package"))
                    if (!$util.isString(message["package"]))
                        return "package: string expected";
                if (message.dependency != null && message.hasOwnProperty("dependency")) {
                    if (!Array.isArray(message.dependency))
                        return "dependency: array expected";
                    for (var i = 0; i < message.dependency.length; ++i)
                        if (!$util.isString(message.dependency[i]))
                            return "dependency: string[] expected";
                }
                if (message.publicDependency != null && message.hasOwnProperty("publicDependency")) {
                    if (!Array.isArray(message.publicDependency))
                        return "publicDependency: array expected";
                    for (var i = 0; i < message.publicDependency.length; ++i)
                        if (!$util.isInteger(message.publicDependency[i]))
                            return "publicDependency: integer[] expected";
                }
                if (message.weakDependency != null && message.hasOwnProperty("weakDependency")) {
                    if (!Array.isArray(message.weakDependency))
                        return "weakDependency: array expected";
                    for (var i = 0; i < message.weakDependency.length; ++i)
                        if (!$util.isInteger(message.weakDependency[i]))
                            return "weakDependency: integer[] expected";
                }
                if (message.messageType != null && message.hasOwnProperty("messageType")) {
                    if (!Array.isArray(message.messageType))
                        return "messageType: array expected";
                    for (var i = 0; i < message.messageType.length; ++i) {
                        var error = $root.google.protobuf.DescriptorProto.verify(message.messageType[i]);
                        if (error)
                            return "messageType." + error;
                    }
                }
                if (message.enumType != null && message.hasOwnProperty("enumType")) {
                    if (!Array.isArray(message.enumType))
                        return "enumType: array expected";
                    for (var i = 0; i < message.enumType.length; ++i) {
                        var error = $root.google.protobuf.EnumDescriptorProto.verify(message.enumType[i]);
                        if (error)
                            return "enumType." + error;
                    }
                }
                if (message.service != null && message.hasOwnProperty("service")) {
                    if (!Array.isArray(message.service))
                        return "service: array expected";
                    for (var i = 0; i < message.service.length; ++i) {
                        var error = $root.google.protobuf.ServiceDescriptorProto.verify(message.service[i]);
                        if (error)
                            return "service." + error;
                    }
                }
                if (message.extension != null && message.hasOwnProperty("extension")) {
                    if (!Array.isArray(message.extension))
                        return "extension: array expected";
                    for (var i = 0; i < message.extension.length; ++i) {
                        var error = $root.google.protobuf.FieldDescriptorProto.verify(message.extension[i]);
                        if (error)
                            return "extension." + error;
                    }
                }
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.FileOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                if (message.sourceCodeInfo != null && message.hasOwnProperty("sourceCodeInfo")) {
                    var error = $root.google.protobuf.SourceCodeInfo.verify(message.sourceCodeInfo);
                    if (error)
                        return "sourceCodeInfo." + error;
                }
                if (message.syntax != null && message.hasOwnProperty("syntax"))
                    if (!$util.isString(message.syntax))
                        return "syntax: string expected";
                return null;
            };

            /**
             * Creates a FileDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.FileDescriptorProto} FileDescriptorProto
             */
            FileDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.FileDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.FileDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object["package"] != null)
                    message["package"] = String(object["package"]);
                if (object.dependency) {
                    if (!Array.isArray(object.dependency))
                        throw TypeError(".google.protobuf.FileDescriptorProto.dependency: array expected");
                    message.dependency = [];
                    for (var i = 0; i < object.dependency.length; ++i)
                        message.dependency[i] = String(object.dependency[i]);
                }
                if (object.publicDependency) {
                    if (!Array.isArray(object.publicDependency))
                        throw TypeError(".google.protobuf.FileDescriptorProto.publicDependency: array expected");
                    message.publicDependency = [];
                    for (var i = 0; i < object.publicDependency.length; ++i)
                        message.publicDependency[i] = object.publicDependency[i] | 0;
                }
                if (object.weakDependency) {
                    if (!Array.isArray(object.weakDependency))
                        throw TypeError(".google.protobuf.FileDescriptorProto.weakDependency: array expected");
                    message.weakDependency = [];
                    for (var i = 0; i < object.weakDependency.length; ++i)
                        message.weakDependency[i] = object.weakDependency[i] | 0;
                }
                if (object.messageType) {
                    if (!Array.isArray(object.messageType))
                        throw TypeError(".google.protobuf.FileDescriptorProto.messageType: array expected");
                    message.messageType = [];
                    for (var i = 0; i < object.messageType.length; ++i) {
                        if (typeof object.messageType[i] !== "object")
                            throw TypeError(".google.protobuf.FileDescriptorProto.messageType: object expected");
                        message.messageType[i] = $root.google.protobuf.DescriptorProto.fromObject(object.messageType[i]);
                    }
                }
                if (object.enumType) {
                    if (!Array.isArray(object.enumType))
                        throw TypeError(".google.protobuf.FileDescriptorProto.enumType: array expected");
                    message.enumType = [];
                    for (var i = 0; i < object.enumType.length; ++i) {
                        if (typeof object.enumType[i] !== "object")
                            throw TypeError(".google.protobuf.FileDescriptorProto.enumType: object expected");
                        message.enumType[i] = $root.google.protobuf.EnumDescriptorProto.fromObject(object.enumType[i]);
                    }
                }
                if (object.service) {
                    if (!Array.isArray(object.service))
                        throw TypeError(".google.protobuf.FileDescriptorProto.service: array expected");
                    message.service = [];
                    for (var i = 0; i < object.service.length; ++i) {
                        if (typeof object.service[i] !== "object")
                            throw TypeError(".google.protobuf.FileDescriptorProto.service: object expected");
                        message.service[i] = $root.google.protobuf.ServiceDescriptorProto.fromObject(object.service[i]);
                    }
                }
                if (object.extension) {
                    if (!Array.isArray(object.extension))
                        throw TypeError(".google.protobuf.FileDescriptorProto.extension: array expected");
                    message.extension = [];
                    for (var i = 0; i < object.extension.length; ++i) {
                        if (typeof object.extension[i] !== "object")
                            throw TypeError(".google.protobuf.FileDescriptorProto.extension: object expected");
                        message.extension[i] = $root.google.protobuf.FieldDescriptorProto.fromObject(object.extension[i]);
                    }
                }
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.FileDescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.FileOptions.fromObject(object.options);
                }
                if (object.sourceCodeInfo != null) {
                    if (typeof object.sourceCodeInfo !== "object")
                        throw TypeError(".google.protobuf.FileDescriptorProto.sourceCodeInfo: object expected");
                    message.sourceCodeInfo = $root.google.protobuf.SourceCodeInfo.fromObject(object.sourceCodeInfo);
                }
                if (object.syntax != null)
                    message.syntax = String(object.syntax);
                return message;
            };

            /**
             * Creates a plain object from a FileDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.FileDescriptorProto
             * @static
             * @param {google.protobuf.FileDescriptorProto} message FileDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            FileDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults) {
                    object.dependency = [];
                    object.messageType = [];
                    object.enumType = [];
                    object.service = [];
                    object.extension = [];
                    object.publicDependency = [];
                    object.weakDependency = [];
                }
                if (options.defaults) {
                    object.name = "";
                    object["package"] = "";
                    object.options = null;
                    object.sourceCodeInfo = null;
                    object.syntax = "";
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message["package"] != null && message.hasOwnProperty("package"))
                    object["package"] = message["package"];
                if (message.dependency && message.dependency.length) {
                    object.dependency = [];
                    for (var j = 0; j < message.dependency.length; ++j)
                        object.dependency[j] = message.dependency[j];
                }
                if (message.messageType && message.messageType.length) {
                    object.messageType = [];
                    for (var j = 0; j < message.messageType.length; ++j)
                        object.messageType[j] = $root.google.protobuf.DescriptorProto.toObject(message.messageType[j], options);
                }
                if (message.enumType && message.enumType.length) {
                    object.enumType = [];
                    for (var j = 0; j < message.enumType.length; ++j)
                        object.enumType[j] = $root.google.protobuf.EnumDescriptorProto.toObject(message.enumType[j], options);
                }
                if (message.service && message.service.length) {
                    object.service = [];
                    for (var j = 0; j < message.service.length; ++j)
                        object.service[j] = $root.google.protobuf.ServiceDescriptorProto.toObject(message.service[j], options);
                }
                if (message.extension && message.extension.length) {
                    object.extension = [];
                    for (var j = 0; j < message.extension.length; ++j)
                        object.extension[j] = $root.google.protobuf.FieldDescriptorProto.toObject(message.extension[j], options);
                }
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.FileOptions.toObject(message.options, options);
                if (message.sourceCodeInfo != null && message.hasOwnProperty("sourceCodeInfo"))
                    object.sourceCodeInfo = $root.google.protobuf.SourceCodeInfo.toObject(message.sourceCodeInfo, options);
                if (message.publicDependency && message.publicDependency.length) {
                    object.publicDependency = [];
                    for (var j = 0; j < message.publicDependency.length; ++j)
                        object.publicDependency[j] = message.publicDependency[j];
                }
                if (message.weakDependency && message.weakDependency.length) {
                    object.weakDependency = [];
                    for (var j = 0; j < message.weakDependency.length; ++j)
                        object.weakDependency[j] = message.weakDependency[j];
                }
                if (message.syntax != null && message.hasOwnProperty("syntax"))
                    object.syntax = message.syntax;
                return object;
            };

            /**
             * Converts this FileDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.FileDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            FileDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return FileDescriptorProto;
        })();

        protobuf.DescriptorProto = (function() {

            /**
             * Properties of a DescriptorProto.
             * @memberof google.protobuf
             * @interface IDescriptorProto
             * @property {string|null} [name] DescriptorProto name
             * @property {Array.<google.protobuf.IFieldDescriptorProto>|null} [field] DescriptorProto field
             * @property {Array.<google.protobuf.IFieldDescriptorProto>|null} [extension] DescriptorProto extension
             * @property {Array.<google.protobuf.IDescriptorProto>|null} [nestedType] DescriptorProto nestedType
             * @property {Array.<google.protobuf.IEnumDescriptorProto>|null} [enumType] DescriptorProto enumType
             * @property {Array.<google.protobuf.DescriptorProto.IExtensionRange>|null} [extensionRange] DescriptorProto extensionRange
             * @property {Array.<google.protobuf.IOneofDescriptorProto>|null} [oneofDecl] DescriptorProto oneofDecl
             * @property {google.protobuf.IMessageOptions|null} [options] DescriptorProto options
             */

            /**
             * Constructs a new DescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents a DescriptorProto.
             * @implements IDescriptorProto
             * @constructor
             * @param {google.protobuf.IDescriptorProto=} [properties] Properties to set
             */
            function DescriptorProto(properties) {
                this.field = [];
                this.extension = [];
                this.nestedType = [];
                this.enumType = [];
                this.extensionRange = [];
                this.oneofDecl = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * DescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.name = "";

            /**
             * DescriptorProto field.
             * @member {Array.<google.protobuf.IFieldDescriptorProto>} field
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.field = $util.emptyArray;

            /**
             * DescriptorProto extension.
             * @member {Array.<google.protobuf.IFieldDescriptorProto>} extension
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.extension = $util.emptyArray;

            /**
             * DescriptorProto nestedType.
             * @member {Array.<google.protobuf.IDescriptorProto>} nestedType
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.nestedType = $util.emptyArray;

            /**
             * DescriptorProto enumType.
             * @member {Array.<google.protobuf.IEnumDescriptorProto>} enumType
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.enumType = $util.emptyArray;

            /**
             * DescriptorProto extensionRange.
             * @member {Array.<google.protobuf.DescriptorProto.IExtensionRange>} extensionRange
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.extensionRange = $util.emptyArray;

            /**
             * DescriptorProto oneofDecl.
             * @member {Array.<google.protobuf.IOneofDescriptorProto>} oneofDecl
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.oneofDecl = $util.emptyArray;

            /**
             * DescriptorProto options.
             * @member {google.protobuf.IMessageOptions|null|undefined} options
             * @memberof google.protobuf.DescriptorProto
             * @instance
             */
            DescriptorProto.prototype.options = null;

            /**
             * Creates a new DescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {google.protobuf.IDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.DescriptorProto} DescriptorProto instance
             */
            DescriptorProto.create = function create(properties) {
                return new DescriptorProto(properties);
            };

            /**
             * Encodes the specified DescriptorProto message. Does not implicitly {@link google.protobuf.DescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {google.protobuf.IDescriptorProto} message DescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            DescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message.field != null && message.field.length)
                    for (var i = 0; i < message.field.length; ++i)
                        $root.google.protobuf.FieldDescriptorProto.encode(message.field[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                if (message.nestedType != null && message.nestedType.length)
                    for (var i = 0; i < message.nestedType.length; ++i)
                        $root.google.protobuf.DescriptorProto.encode(message.nestedType[i], writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
                if (message.enumType != null && message.enumType.length)
                    for (var i = 0; i < message.enumType.length; ++i)
                        $root.google.protobuf.EnumDescriptorProto.encode(message.enumType[i], writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
                if (message.extensionRange != null && message.extensionRange.length)
                    for (var i = 0; i < message.extensionRange.length; ++i)
                        $root.google.protobuf.DescriptorProto.ExtensionRange.encode(message.extensionRange[i], writer.uint32(/* id 5, wireType 2 =*/42).fork()).ldelim();
                if (message.extension != null && message.extension.length)
                    for (var i = 0; i < message.extension.length; ++i)
                        $root.google.protobuf.FieldDescriptorProto.encode(message.extension[i], writer.uint32(/* id 6, wireType 2 =*/50).fork()).ldelim();
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.MessageOptions.encode(message.options, writer.uint32(/* id 7, wireType 2 =*/58).fork()).ldelim();
                if (message.oneofDecl != null && message.oneofDecl.length)
                    for (var i = 0; i < message.oneofDecl.length; ++i)
                        $root.google.protobuf.OneofDescriptorProto.encode(message.oneofDecl[i], writer.uint32(/* id 8, wireType 2 =*/66).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified DescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.DescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {google.protobuf.IDescriptorProto} message DescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            DescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a DescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.DescriptorProto} DescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            DescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.DescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 2:
                        if (!(message.field && message.field.length))
                            message.field = [];
                        message.field.push($root.google.protobuf.FieldDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 6:
                        if (!(message.extension && message.extension.length))
                            message.extension = [];
                        message.extension.push($root.google.protobuf.FieldDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 3:
                        if (!(message.nestedType && message.nestedType.length))
                            message.nestedType = [];
                        message.nestedType.push($root.google.protobuf.DescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 4:
                        if (!(message.enumType && message.enumType.length))
                            message.enumType = [];
                        message.enumType.push($root.google.protobuf.EnumDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 5:
                        if (!(message.extensionRange && message.extensionRange.length))
                            message.extensionRange = [];
                        message.extensionRange.push($root.google.protobuf.DescriptorProto.ExtensionRange.decode(reader, reader.uint32()));
                        break;
                    case 8:
                        if (!(message.oneofDecl && message.oneofDecl.length))
                            message.oneofDecl = [];
                        message.oneofDecl.push($root.google.protobuf.OneofDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 7:
                        message.options = $root.google.protobuf.MessageOptions.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a DescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.DescriptorProto} DescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            DescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a DescriptorProto message.
             * @function verify
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            DescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message.field != null && message.hasOwnProperty("field")) {
                    if (!Array.isArray(message.field))
                        return "field: array expected";
                    for (var i = 0; i < message.field.length; ++i) {
                        var error = $root.google.protobuf.FieldDescriptorProto.verify(message.field[i]);
                        if (error)
                            return "field." + error;
                    }
                }
                if (message.extension != null && message.hasOwnProperty("extension")) {
                    if (!Array.isArray(message.extension))
                        return "extension: array expected";
                    for (var i = 0; i < message.extension.length; ++i) {
                        var error = $root.google.protobuf.FieldDescriptorProto.verify(message.extension[i]);
                        if (error)
                            return "extension." + error;
                    }
                }
                if (message.nestedType != null && message.hasOwnProperty("nestedType")) {
                    if (!Array.isArray(message.nestedType))
                        return "nestedType: array expected";
                    for (var i = 0; i < message.nestedType.length; ++i) {
                        var error = $root.google.protobuf.DescriptorProto.verify(message.nestedType[i]);
                        if (error)
                            return "nestedType." + error;
                    }
                }
                if (message.enumType != null && message.hasOwnProperty("enumType")) {
                    if (!Array.isArray(message.enumType))
                        return "enumType: array expected";
                    for (var i = 0; i < message.enumType.length; ++i) {
                        var error = $root.google.protobuf.EnumDescriptorProto.verify(message.enumType[i]);
                        if (error)
                            return "enumType." + error;
                    }
                }
                if (message.extensionRange != null && message.hasOwnProperty("extensionRange")) {
                    if (!Array.isArray(message.extensionRange))
                        return "extensionRange: array expected";
                    for (var i = 0; i < message.extensionRange.length; ++i) {
                        var error = $root.google.protobuf.DescriptorProto.ExtensionRange.verify(message.extensionRange[i]);
                        if (error)
                            return "extensionRange." + error;
                    }
                }
                if (message.oneofDecl != null && message.hasOwnProperty("oneofDecl")) {
                    if (!Array.isArray(message.oneofDecl))
                        return "oneofDecl: array expected";
                    for (var i = 0; i < message.oneofDecl.length; ++i) {
                        var error = $root.google.protobuf.OneofDescriptorProto.verify(message.oneofDecl[i]);
                        if (error)
                            return "oneofDecl." + error;
                    }
                }
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.MessageOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                return null;
            };

            /**
             * Creates a DescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.DescriptorProto} DescriptorProto
             */
            DescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.DescriptorProto)
                    return object;
                var message = new $root.google.protobuf.DescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object.field) {
                    if (!Array.isArray(object.field))
                        throw TypeError(".google.protobuf.DescriptorProto.field: array expected");
                    message.field = [];
                    for (var i = 0; i < object.field.length; ++i) {
                        if (typeof object.field[i] !== "object")
                            throw TypeError(".google.protobuf.DescriptorProto.field: object expected");
                        message.field[i] = $root.google.protobuf.FieldDescriptorProto.fromObject(object.field[i]);
                    }
                }
                if (object.extension) {
                    if (!Array.isArray(object.extension))
                        throw TypeError(".google.protobuf.DescriptorProto.extension: array expected");
                    message.extension = [];
                    for (var i = 0; i < object.extension.length; ++i) {
                        if (typeof object.extension[i] !== "object")
                            throw TypeError(".google.protobuf.DescriptorProto.extension: object expected");
                        message.extension[i] = $root.google.protobuf.FieldDescriptorProto.fromObject(object.extension[i]);
                    }
                }
                if (object.nestedType) {
                    if (!Array.isArray(object.nestedType))
                        throw TypeError(".google.protobuf.DescriptorProto.nestedType: array expected");
                    message.nestedType = [];
                    for (var i = 0; i < object.nestedType.length; ++i) {
                        if (typeof object.nestedType[i] !== "object")
                            throw TypeError(".google.protobuf.DescriptorProto.nestedType: object expected");
                        message.nestedType[i] = $root.google.protobuf.DescriptorProto.fromObject(object.nestedType[i]);
                    }
                }
                if (object.enumType) {
                    if (!Array.isArray(object.enumType))
                        throw TypeError(".google.protobuf.DescriptorProto.enumType: array expected");
                    message.enumType = [];
                    for (var i = 0; i < object.enumType.length; ++i) {
                        if (typeof object.enumType[i] !== "object")
                            throw TypeError(".google.protobuf.DescriptorProto.enumType: object expected");
                        message.enumType[i] = $root.google.protobuf.EnumDescriptorProto.fromObject(object.enumType[i]);
                    }
                }
                if (object.extensionRange) {
                    if (!Array.isArray(object.extensionRange))
                        throw TypeError(".google.protobuf.DescriptorProto.extensionRange: array expected");
                    message.extensionRange = [];
                    for (var i = 0; i < object.extensionRange.length; ++i) {
                        if (typeof object.extensionRange[i] !== "object")
                            throw TypeError(".google.protobuf.DescriptorProto.extensionRange: object expected");
                        message.extensionRange[i] = $root.google.protobuf.DescriptorProto.ExtensionRange.fromObject(object.extensionRange[i]);
                    }
                }
                if (object.oneofDecl) {
                    if (!Array.isArray(object.oneofDecl))
                        throw TypeError(".google.protobuf.DescriptorProto.oneofDecl: array expected");
                    message.oneofDecl = [];
                    for (var i = 0; i < object.oneofDecl.length; ++i) {
                        if (typeof object.oneofDecl[i] !== "object")
                            throw TypeError(".google.protobuf.DescriptorProto.oneofDecl: object expected");
                        message.oneofDecl[i] = $root.google.protobuf.OneofDescriptorProto.fromObject(object.oneofDecl[i]);
                    }
                }
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.DescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.MessageOptions.fromObject(object.options);
                }
                return message;
            };

            /**
             * Creates a plain object from a DescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.DescriptorProto
             * @static
             * @param {google.protobuf.DescriptorProto} message DescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            DescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults) {
                    object.field = [];
                    object.nestedType = [];
                    object.enumType = [];
                    object.extensionRange = [];
                    object.extension = [];
                    object.oneofDecl = [];
                }
                if (options.defaults) {
                    object.name = "";
                    object.options = null;
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message.field && message.field.length) {
                    object.field = [];
                    for (var j = 0; j < message.field.length; ++j)
                        object.field[j] = $root.google.protobuf.FieldDescriptorProto.toObject(message.field[j], options);
                }
                if (message.nestedType && message.nestedType.length) {
                    object.nestedType = [];
                    for (var j = 0; j < message.nestedType.length; ++j)
                        object.nestedType[j] = $root.google.protobuf.DescriptorProto.toObject(message.nestedType[j], options);
                }
                if (message.enumType && message.enumType.length) {
                    object.enumType = [];
                    for (var j = 0; j < message.enumType.length; ++j)
                        object.enumType[j] = $root.google.protobuf.EnumDescriptorProto.toObject(message.enumType[j], options);
                }
                if (message.extensionRange && message.extensionRange.length) {
                    object.extensionRange = [];
                    for (var j = 0; j < message.extensionRange.length; ++j)
                        object.extensionRange[j] = $root.google.protobuf.DescriptorProto.ExtensionRange.toObject(message.extensionRange[j], options);
                }
                if (message.extension && message.extension.length) {
                    object.extension = [];
                    for (var j = 0; j < message.extension.length; ++j)
                        object.extension[j] = $root.google.protobuf.FieldDescriptorProto.toObject(message.extension[j], options);
                }
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.MessageOptions.toObject(message.options, options);
                if (message.oneofDecl && message.oneofDecl.length) {
                    object.oneofDecl = [];
                    for (var j = 0; j < message.oneofDecl.length; ++j)
                        object.oneofDecl[j] = $root.google.protobuf.OneofDescriptorProto.toObject(message.oneofDecl[j], options);
                }
                return object;
            };

            /**
             * Converts this DescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.DescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            DescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            DescriptorProto.ExtensionRange = (function() {

                /**
                 * Properties of an ExtensionRange.
                 * @memberof google.protobuf.DescriptorProto
                 * @interface IExtensionRange
                 * @property {number|null} [start] ExtensionRange start
                 * @property {number|null} [end] ExtensionRange end
                 */

                /**
                 * Constructs a new ExtensionRange.
                 * @memberof google.protobuf.DescriptorProto
                 * @classdesc Represents an ExtensionRange.
                 * @implements IExtensionRange
                 * @constructor
                 * @param {google.protobuf.DescriptorProto.IExtensionRange=} [properties] Properties to set
                 */
                function ExtensionRange(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * ExtensionRange start.
                 * @member {number} start
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @instance
                 */
                ExtensionRange.prototype.start = 0;

                /**
                 * ExtensionRange end.
                 * @member {number} end
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @instance
                 */
                ExtensionRange.prototype.end = 0;

                /**
                 * Creates a new ExtensionRange instance using the specified properties.
                 * @function create
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {google.protobuf.DescriptorProto.IExtensionRange=} [properties] Properties to set
                 * @returns {google.protobuf.DescriptorProto.ExtensionRange} ExtensionRange instance
                 */
                ExtensionRange.create = function create(properties) {
                    return new ExtensionRange(properties);
                };

                /**
                 * Encodes the specified ExtensionRange message. Does not implicitly {@link google.protobuf.DescriptorProto.ExtensionRange.verify|verify} messages.
                 * @function encode
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {google.protobuf.DescriptorProto.IExtensionRange} message ExtensionRange message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                ExtensionRange.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.start != null && message.hasOwnProperty("start"))
                        writer.uint32(/* id 1, wireType 0 =*/8).int32(message.start);
                    if (message.end != null && message.hasOwnProperty("end"))
                        writer.uint32(/* id 2, wireType 0 =*/16).int32(message.end);
                    return writer;
                };

                /**
                 * Encodes the specified ExtensionRange message, length delimited. Does not implicitly {@link google.protobuf.DescriptorProto.ExtensionRange.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {google.protobuf.DescriptorProto.IExtensionRange} message ExtensionRange message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                ExtensionRange.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes an ExtensionRange message from the specified reader or buffer.
                 * @function decode
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {google.protobuf.DescriptorProto.ExtensionRange} ExtensionRange
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                ExtensionRange.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.DescriptorProto.ExtensionRange();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.start = reader.int32();
                            break;
                        case 2:
                            message.end = reader.int32();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes an ExtensionRange message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {google.protobuf.DescriptorProto.ExtensionRange} ExtensionRange
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                ExtensionRange.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies an ExtensionRange message.
                 * @function verify
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                ExtensionRange.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.start != null && message.hasOwnProperty("start"))
                        if (!$util.isInteger(message.start))
                            return "start: integer expected";
                    if (message.end != null && message.hasOwnProperty("end"))
                        if (!$util.isInteger(message.end))
                            return "end: integer expected";
                    return null;
                };

                /**
                 * Creates an ExtensionRange message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {google.protobuf.DescriptorProto.ExtensionRange} ExtensionRange
                 */
                ExtensionRange.fromObject = function fromObject(object) {
                    if (object instanceof $root.google.protobuf.DescriptorProto.ExtensionRange)
                        return object;
                    var message = new $root.google.protobuf.DescriptorProto.ExtensionRange();
                    if (object.start != null)
                        message.start = object.start | 0;
                    if (object.end != null)
                        message.end = object.end | 0;
                    return message;
                };

                /**
                 * Creates a plain object from an ExtensionRange message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @static
                 * @param {google.protobuf.DescriptorProto.ExtensionRange} message ExtensionRange
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                ExtensionRange.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.start = 0;
                        object.end = 0;
                    }
                    if (message.start != null && message.hasOwnProperty("start"))
                        object.start = message.start;
                    if (message.end != null && message.hasOwnProperty("end"))
                        object.end = message.end;
                    return object;
                };

                /**
                 * Converts this ExtensionRange to JSON.
                 * @function toJSON
                 * @memberof google.protobuf.DescriptorProto.ExtensionRange
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                ExtensionRange.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return ExtensionRange;
            })();

            return DescriptorProto;
        })();

        protobuf.FieldDescriptorProto = (function() {

            /**
             * Properties of a FieldDescriptorProto.
             * @memberof google.protobuf
             * @interface IFieldDescriptorProto
             * @property {string|null} [name] FieldDescriptorProto name
             * @property {number|null} [number] FieldDescriptorProto number
             * @property {google.protobuf.FieldDescriptorProto.Label|null} [label] FieldDescriptorProto label
             * @property {google.protobuf.FieldDescriptorProto.Type|null} [type] FieldDescriptorProto type
             * @property {string|null} [typeName] FieldDescriptorProto typeName
             * @property {string|null} [extendee] FieldDescriptorProto extendee
             * @property {string|null} [defaultValue] FieldDescriptorProto defaultValue
             * @property {number|null} [oneofIndex] FieldDescriptorProto oneofIndex
             * @property {google.protobuf.IFieldOptions|null} [options] FieldDescriptorProto options
             */

            /**
             * Constructs a new FieldDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents a FieldDescriptorProto.
             * @implements IFieldDescriptorProto
             * @constructor
             * @param {google.protobuf.IFieldDescriptorProto=} [properties] Properties to set
             */
            function FieldDescriptorProto(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * FieldDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.name = "";

            /**
             * FieldDescriptorProto number.
             * @member {number} number
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.number = 0;

            /**
             * FieldDescriptorProto label.
             * @member {google.protobuf.FieldDescriptorProto.Label} label
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.label = 1;

            /**
             * FieldDescriptorProto type.
             * @member {google.protobuf.FieldDescriptorProto.Type} type
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.type = 1;

            /**
             * FieldDescriptorProto typeName.
             * @member {string} typeName
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.typeName = "";

            /**
             * FieldDescriptorProto extendee.
             * @member {string} extendee
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.extendee = "";

            /**
             * FieldDescriptorProto defaultValue.
             * @member {string} defaultValue
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.defaultValue = "";

            /**
             * FieldDescriptorProto oneofIndex.
             * @member {number} oneofIndex
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.oneofIndex = 0;

            /**
             * FieldDescriptorProto options.
             * @member {google.protobuf.IFieldOptions|null|undefined} options
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             */
            FieldDescriptorProto.prototype.options = null;

            /**
             * Creates a new FieldDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {google.protobuf.IFieldDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.FieldDescriptorProto} FieldDescriptorProto instance
             */
            FieldDescriptorProto.create = function create(properties) {
                return new FieldDescriptorProto(properties);
            };

            /**
             * Encodes the specified FieldDescriptorProto message. Does not implicitly {@link google.protobuf.FieldDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {google.protobuf.IFieldDescriptorProto} message FieldDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FieldDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message.extendee != null && message.hasOwnProperty("extendee"))
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message.extendee);
                if (message.number != null && message.hasOwnProperty("number"))
                    writer.uint32(/* id 3, wireType 0 =*/24).int32(message.number);
                if (message.label != null && message.hasOwnProperty("label"))
                    writer.uint32(/* id 4, wireType 0 =*/32).int32(message.label);
                if (message.type != null && message.hasOwnProperty("type"))
                    writer.uint32(/* id 5, wireType 0 =*/40).int32(message.type);
                if (message.typeName != null && message.hasOwnProperty("typeName"))
                    writer.uint32(/* id 6, wireType 2 =*/50).string(message.typeName);
                if (message.defaultValue != null && message.hasOwnProperty("defaultValue"))
                    writer.uint32(/* id 7, wireType 2 =*/58).string(message.defaultValue);
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.FieldOptions.encode(message.options, writer.uint32(/* id 8, wireType 2 =*/66).fork()).ldelim();
                if (message.oneofIndex != null && message.hasOwnProperty("oneofIndex"))
                    writer.uint32(/* id 9, wireType 0 =*/72).int32(message.oneofIndex);
                return writer;
            };

            /**
             * Encodes the specified FieldDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.FieldDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {google.protobuf.IFieldDescriptorProto} message FieldDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FieldDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a FieldDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.FieldDescriptorProto} FieldDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FieldDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.FieldDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 3:
                        message.number = reader.int32();
                        break;
                    case 4:
                        message.label = reader.int32();
                        break;
                    case 5:
                        message.type = reader.int32();
                        break;
                    case 6:
                        message.typeName = reader.string();
                        break;
                    case 2:
                        message.extendee = reader.string();
                        break;
                    case 7:
                        message.defaultValue = reader.string();
                        break;
                    case 9:
                        message.oneofIndex = reader.int32();
                        break;
                    case 8:
                        message.options = $root.google.protobuf.FieldOptions.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a FieldDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.FieldDescriptorProto} FieldDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FieldDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a FieldDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            FieldDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message.number != null && message.hasOwnProperty("number"))
                    if (!$util.isInteger(message.number))
                        return "number: integer expected";
                if (message.label != null && message.hasOwnProperty("label"))
                    switch (message.label) {
                    default:
                        return "label: enum value expected";
                    case 1:
                    case 2:
                    case 3:
                        break;
                    }
                if (message.type != null && message.hasOwnProperty("type"))
                    switch (message.type) {
                    default:
                        return "type: enum value expected";
                    case 1:
                    case 2:
                    case 3:
                    case 4:
                    case 5:
                    case 6:
                    case 7:
                    case 8:
                    case 9:
                    case 10:
                    case 11:
                    case 12:
                    case 13:
                    case 14:
                    case 15:
                    case 16:
                    case 17:
                    case 18:
                        break;
                    }
                if (message.typeName != null && message.hasOwnProperty("typeName"))
                    if (!$util.isString(message.typeName))
                        return "typeName: string expected";
                if (message.extendee != null && message.hasOwnProperty("extendee"))
                    if (!$util.isString(message.extendee))
                        return "extendee: string expected";
                if (message.defaultValue != null && message.hasOwnProperty("defaultValue"))
                    if (!$util.isString(message.defaultValue))
                        return "defaultValue: string expected";
                if (message.oneofIndex != null && message.hasOwnProperty("oneofIndex"))
                    if (!$util.isInteger(message.oneofIndex))
                        return "oneofIndex: integer expected";
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.FieldOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                return null;
            };

            /**
             * Creates a FieldDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.FieldDescriptorProto} FieldDescriptorProto
             */
            FieldDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.FieldDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.FieldDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object.number != null)
                    message.number = object.number | 0;
                switch (object.label) {
                case "LABEL_OPTIONAL":
                case 1:
                    message.label = 1;
                    break;
                case "LABEL_REQUIRED":
                case 2:
                    message.label = 2;
                    break;
                case "LABEL_REPEATED":
                case 3:
                    message.label = 3;
                    break;
                }
                switch (object.type) {
                case "TYPE_DOUBLE":
                case 1:
                    message.type = 1;
                    break;
                case "TYPE_FLOAT":
                case 2:
                    message.type = 2;
                    break;
                case "TYPE_INT64":
                case 3:
                    message.type = 3;
                    break;
                case "TYPE_UINT64":
                case 4:
                    message.type = 4;
                    break;
                case "TYPE_INT32":
                case 5:
                    message.type = 5;
                    break;
                case "TYPE_FIXED64":
                case 6:
                    message.type = 6;
                    break;
                case "TYPE_FIXED32":
                case 7:
                    message.type = 7;
                    break;
                case "TYPE_BOOL":
                case 8:
                    message.type = 8;
                    break;
                case "TYPE_STRING":
                case 9:
                    message.type = 9;
                    break;
                case "TYPE_GROUP":
                case 10:
                    message.type = 10;
                    break;
                case "TYPE_MESSAGE":
                case 11:
                    message.type = 11;
                    break;
                case "TYPE_BYTES":
                case 12:
                    message.type = 12;
                    break;
                case "TYPE_UINT32":
                case 13:
                    message.type = 13;
                    break;
                case "TYPE_ENUM":
                case 14:
                    message.type = 14;
                    break;
                case "TYPE_SFIXED32":
                case 15:
                    message.type = 15;
                    break;
                case "TYPE_SFIXED64":
                case 16:
                    message.type = 16;
                    break;
                case "TYPE_SINT32":
                case 17:
                    message.type = 17;
                    break;
                case "TYPE_SINT64":
                case 18:
                    message.type = 18;
                    break;
                }
                if (object.typeName != null)
                    message.typeName = String(object.typeName);
                if (object.extendee != null)
                    message.extendee = String(object.extendee);
                if (object.defaultValue != null)
                    message.defaultValue = String(object.defaultValue);
                if (object.oneofIndex != null)
                    message.oneofIndex = object.oneofIndex | 0;
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.FieldDescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.FieldOptions.fromObject(object.options);
                }
                return message;
            };

            /**
             * Creates a plain object from a FieldDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.FieldDescriptorProto
             * @static
             * @param {google.protobuf.FieldDescriptorProto} message FieldDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            FieldDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.name = "";
                    object.extendee = "";
                    object.number = 0;
                    object.label = options.enums === String ? "LABEL_OPTIONAL" : 1;
                    object.type = options.enums === String ? "TYPE_DOUBLE" : 1;
                    object.typeName = "";
                    object.defaultValue = "";
                    object.options = null;
                    object.oneofIndex = 0;
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message.extendee != null && message.hasOwnProperty("extendee"))
                    object.extendee = message.extendee;
                if (message.number != null && message.hasOwnProperty("number"))
                    object.number = message.number;
                if (message.label != null && message.hasOwnProperty("label"))
                    object.label = options.enums === String ? $root.google.protobuf.FieldDescriptorProto.Label[message.label] : message.label;
                if (message.type != null && message.hasOwnProperty("type"))
                    object.type = options.enums === String ? $root.google.protobuf.FieldDescriptorProto.Type[message.type] : message.type;
                if (message.typeName != null && message.hasOwnProperty("typeName"))
                    object.typeName = message.typeName;
                if (message.defaultValue != null && message.hasOwnProperty("defaultValue"))
                    object.defaultValue = message.defaultValue;
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.FieldOptions.toObject(message.options, options);
                if (message.oneofIndex != null && message.hasOwnProperty("oneofIndex"))
                    object.oneofIndex = message.oneofIndex;
                return object;
            };

            /**
             * Converts this FieldDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.FieldDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            FieldDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * Type enum.
             * @name google.protobuf.FieldDescriptorProto.Type
             * @enum {string}
             * @property {number} TYPE_DOUBLE=1 TYPE_DOUBLE value
             * @property {number} TYPE_FLOAT=2 TYPE_FLOAT value
             * @property {number} TYPE_INT64=3 TYPE_INT64 value
             * @property {number} TYPE_UINT64=4 TYPE_UINT64 value
             * @property {number} TYPE_INT32=5 TYPE_INT32 value
             * @property {number} TYPE_FIXED64=6 TYPE_FIXED64 value
             * @property {number} TYPE_FIXED32=7 TYPE_FIXED32 value
             * @property {number} TYPE_BOOL=8 TYPE_BOOL value
             * @property {number} TYPE_STRING=9 TYPE_STRING value
             * @property {number} TYPE_GROUP=10 TYPE_GROUP value
             * @property {number} TYPE_MESSAGE=11 TYPE_MESSAGE value
             * @property {number} TYPE_BYTES=12 TYPE_BYTES value
             * @property {number} TYPE_UINT32=13 TYPE_UINT32 value
             * @property {number} TYPE_ENUM=14 TYPE_ENUM value
             * @property {number} TYPE_SFIXED32=15 TYPE_SFIXED32 value
             * @property {number} TYPE_SFIXED64=16 TYPE_SFIXED64 value
             * @property {number} TYPE_SINT32=17 TYPE_SINT32 value
             * @property {number} TYPE_SINT64=18 TYPE_SINT64 value
             */
            FieldDescriptorProto.Type = (function() {
                var valuesById = {}, values = Object.create(valuesById);
                values[valuesById[1] = "TYPE_DOUBLE"] = 1;
                values[valuesById[2] = "TYPE_FLOAT"] = 2;
                values[valuesById[3] = "TYPE_INT64"] = 3;
                values[valuesById[4] = "TYPE_UINT64"] = 4;
                values[valuesById[5] = "TYPE_INT32"] = 5;
                values[valuesById[6] = "TYPE_FIXED64"] = 6;
                values[valuesById[7] = "TYPE_FIXED32"] = 7;
                values[valuesById[8] = "TYPE_BOOL"] = 8;
                values[valuesById[9] = "TYPE_STRING"] = 9;
                values[valuesById[10] = "TYPE_GROUP"] = 10;
                values[valuesById[11] = "TYPE_MESSAGE"] = 11;
                values[valuesById[12] = "TYPE_BYTES"] = 12;
                values[valuesById[13] = "TYPE_UINT32"] = 13;
                values[valuesById[14] = "TYPE_ENUM"] = 14;
                values[valuesById[15] = "TYPE_SFIXED32"] = 15;
                values[valuesById[16] = "TYPE_SFIXED64"] = 16;
                values[valuesById[17] = "TYPE_SINT32"] = 17;
                values[valuesById[18] = "TYPE_SINT64"] = 18;
                return values;
            })();

            /**
             * Label enum.
             * @name google.protobuf.FieldDescriptorProto.Label
             * @enum {string}
             * @property {number} LABEL_OPTIONAL=1 LABEL_OPTIONAL value
             * @property {number} LABEL_REQUIRED=2 LABEL_REQUIRED value
             * @property {number} LABEL_REPEATED=3 LABEL_REPEATED value
             */
            FieldDescriptorProto.Label = (function() {
                var valuesById = {}, values = Object.create(valuesById);
                values[valuesById[1] = "LABEL_OPTIONAL"] = 1;
                values[valuesById[2] = "LABEL_REQUIRED"] = 2;
                values[valuesById[3] = "LABEL_REPEATED"] = 3;
                return values;
            })();

            return FieldDescriptorProto;
        })();

        protobuf.OneofDescriptorProto = (function() {

            /**
             * Properties of an OneofDescriptorProto.
             * @memberof google.protobuf
             * @interface IOneofDescriptorProto
             * @property {string|null} [name] OneofDescriptorProto name
             */

            /**
             * Constructs a new OneofDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents an OneofDescriptorProto.
             * @implements IOneofDescriptorProto
             * @constructor
             * @param {google.protobuf.IOneofDescriptorProto=} [properties] Properties to set
             */
            function OneofDescriptorProto(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * OneofDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.OneofDescriptorProto
             * @instance
             */
            OneofDescriptorProto.prototype.name = "";

            /**
             * Creates a new OneofDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {google.protobuf.IOneofDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.OneofDescriptorProto} OneofDescriptorProto instance
             */
            OneofDescriptorProto.create = function create(properties) {
                return new OneofDescriptorProto(properties);
            };

            /**
             * Encodes the specified OneofDescriptorProto message. Does not implicitly {@link google.protobuf.OneofDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {google.protobuf.IOneofDescriptorProto} message OneofDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            OneofDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                return writer;
            };

            /**
             * Encodes the specified OneofDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.OneofDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {google.protobuf.IOneofDescriptorProto} message OneofDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            OneofDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an OneofDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.OneofDescriptorProto} OneofDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            OneofDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.OneofDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an OneofDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.OneofDescriptorProto} OneofDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            OneofDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an OneofDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            OneofDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                return null;
            };

            /**
             * Creates an OneofDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.OneofDescriptorProto} OneofDescriptorProto
             */
            OneofDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.OneofDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.OneofDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                return message;
            };

            /**
             * Creates a plain object from an OneofDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.OneofDescriptorProto
             * @static
             * @param {google.protobuf.OneofDescriptorProto} message OneofDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            OneofDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults)
                    object.name = "";
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                return object;
            };

            /**
             * Converts this OneofDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.OneofDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            OneofDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return OneofDescriptorProto;
        })();

        protobuf.EnumDescriptorProto = (function() {

            /**
             * Properties of an EnumDescriptorProto.
             * @memberof google.protobuf
             * @interface IEnumDescriptorProto
             * @property {string|null} [name] EnumDescriptorProto name
             * @property {Array.<google.protobuf.IEnumValueDescriptorProto>|null} [value] EnumDescriptorProto value
             * @property {google.protobuf.IEnumOptions|null} [options] EnumDescriptorProto options
             */

            /**
             * Constructs a new EnumDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents an EnumDescriptorProto.
             * @implements IEnumDescriptorProto
             * @constructor
             * @param {google.protobuf.IEnumDescriptorProto=} [properties] Properties to set
             */
            function EnumDescriptorProto(properties) {
                this.value = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * EnumDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.EnumDescriptorProto
             * @instance
             */
            EnumDescriptorProto.prototype.name = "";

            /**
             * EnumDescriptorProto value.
             * @member {Array.<google.protobuf.IEnumValueDescriptorProto>} value
             * @memberof google.protobuf.EnumDescriptorProto
             * @instance
             */
            EnumDescriptorProto.prototype.value = $util.emptyArray;

            /**
             * EnumDescriptorProto options.
             * @member {google.protobuf.IEnumOptions|null|undefined} options
             * @memberof google.protobuf.EnumDescriptorProto
             * @instance
             */
            EnumDescriptorProto.prototype.options = null;

            /**
             * Creates a new EnumDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {google.protobuf.IEnumDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.EnumDescriptorProto} EnumDescriptorProto instance
             */
            EnumDescriptorProto.create = function create(properties) {
                return new EnumDescriptorProto(properties);
            };

            /**
             * Encodes the specified EnumDescriptorProto message. Does not implicitly {@link google.protobuf.EnumDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {google.protobuf.IEnumDescriptorProto} message EnumDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message.value != null && message.value.length)
                    for (var i = 0; i < message.value.length; ++i)
                        $root.google.protobuf.EnumValueDescriptorProto.encode(message.value[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.EnumOptions.encode(message.options, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified EnumDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.EnumDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {google.protobuf.IEnumDescriptorProto} message EnumDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an EnumDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.EnumDescriptorProto} EnumDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.EnumDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 2:
                        if (!(message.value && message.value.length))
                            message.value = [];
                        message.value.push($root.google.protobuf.EnumValueDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 3:
                        message.options = $root.google.protobuf.EnumOptions.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an EnumDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.EnumDescriptorProto} EnumDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an EnumDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            EnumDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message.value != null && message.hasOwnProperty("value")) {
                    if (!Array.isArray(message.value))
                        return "value: array expected";
                    for (var i = 0; i < message.value.length; ++i) {
                        var error = $root.google.protobuf.EnumValueDescriptorProto.verify(message.value[i]);
                        if (error)
                            return "value." + error;
                    }
                }
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.EnumOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                return null;
            };

            /**
             * Creates an EnumDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.EnumDescriptorProto} EnumDescriptorProto
             */
            EnumDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.EnumDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.EnumDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object.value) {
                    if (!Array.isArray(object.value))
                        throw TypeError(".google.protobuf.EnumDescriptorProto.value: array expected");
                    message.value = [];
                    for (var i = 0; i < object.value.length; ++i) {
                        if (typeof object.value[i] !== "object")
                            throw TypeError(".google.protobuf.EnumDescriptorProto.value: object expected");
                        message.value[i] = $root.google.protobuf.EnumValueDescriptorProto.fromObject(object.value[i]);
                    }
                }
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.EnumDescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.EnumOptions.fromObject(object.options);
                }
                return message;
            };

            /**
             * Creates a plain object from an EnumDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.EnumDescriptorProto
             * @static
             * @param {google.protobuf.EnumDescriptorProto} message EnumDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            EnumDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.value = [];
                if (options.defaults) {
                    object.name = "";
                    object.options = null;
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message.value && message.value.length) {
                    object.value = [];
                    for (var j = 0; j < message.value.length; ++j)
                        object.value[j] = $root.google.protobuf.EnumValueDescriptorProto.toObject(message.value[j], options);
                }
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.EnumOptions.toObject(message.options, options);
                return object;
            };

            /**
             * Converts this EnumDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.EnumDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            EnumDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return EnumDescriptorProto;
        })();

        protobuf.EnumValueDescriptorProto = (function() {

            /**
             * Properties of an EnumValueDescriptorProto.
             * @memberof google.protobuf
             * @interface IEnumValueDescriptorProto
             * @property {string|null} [name] EnumValueDescriptorProto name
             * @property {number|null} [number] EnumValueDescriptorProto number
             * @property {google.protobuf.IEnumValueOptions|null} [options] EnumValueDescriptorProto options
             */

            /**
             * Constructs a new EnumValueDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents an EnumValueDescriptorProto.
             * @implements IEnumValueDescriptorProto
             * @constructor
             * @param {google.protobuf.IEnumValueDescriptorProto=} [properties] Properties to set
             */
            function EnumValueDescriptorProto(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * EnumValueDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @instance
             */
            EnumValueDescriptorProto.prototype.name = "";

            /**
             * EnumValueDescriptorProto number.
             * @member {number} number
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @instance
             */
            EnumValueDescriptorProto.prototype.number = 0;

            /**
             * EnumValueDescriptorProto options.
             * @member {google.protobuf.IEnumValueOptions|null|undefined} options
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @instance
             */
            EnumValueDescriptorProto.prototype.options = null;

            /**
             * Creates a new EnumValueDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {google.protobuf.IEnumValueDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.EnumValueDescriptorProto} EnumValueDescriptorProto instance
             */
            EnumValueDescriptorProto.create = function create(properties) {
                return new EnumValueDescriptorProto(properties);
            };

            /**
             * Encodes the specified EnumValueDescriptorProto message. Does not implicitly {@link google.protobuf.EnumValueDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {google.protobuf.IEnumValueDescriptorProto} message EnumValueDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumValueDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message.number != null && message.hasOwnProperty("number"))
                    writer.uint32(/* id 2, wireType 0 =*/16).int32(message.number);
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.EnumValueOptions.encode(message.options, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified EnumValueDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.EnumValueDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {google.protobuf.IEnumValueDescriptorProto} message EnumValueDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumValueDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an EnumValueDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.EnumValueDescriptorProto} EnumValueDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumValueDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.EnumValueDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 2:
                        message.number = reader.int32();
                        break;
                    case 3:
                        message.options = $root.google.protobuf.EnumValueOptions.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an EnumValueDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.EnumValueDescriptorProto} EnumValueDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumValueDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an EnumValueDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            EnumValueDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message.number != null && message.hasOwnProperty("number"))
                    if (!$util.isInteger(message.number))
                        return "number: integer expected";
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.EnumValueOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                return null;
            };

            /**
             * Creates an EnumValueDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.EnumValueDescriptorProto} EnumValueDescriptorProto
             */
            EnumValueDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.EnumValueDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.EnumValueDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object.number != null)
                    message.number = object.number | 0;
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.EnumValueDescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.EnumValueOptions.fromObject(object.options);
                }
                return message;
            };

            /**
             * Creates a plain object from an EnumValueDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @static
             * @param {google.protobuf.EnumValueDescriptorProto} message EnumValueDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            EnumValueDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.name = "";
                    object.number = 0;
                    object.options = null;
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message.number != null && message.hasOwnProperty("number"))
                    object.number = message.number;
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.EnumValueOptions.toObject(message.options, options);
                return object;
            };

            /**
             * Converts this EnumValueDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.EnumValueDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            EnumValueDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return EnumValueDescriptorProto;
        })();

        protobuf.ServiceDescriptorProto = (function() {

            /**
             * Properties of a ServiceDescriptorProto.
             * @memberof google.protobuf
             * @interface IServiceDescriptorProto
             * @property {string|null} [name] ServiceDescriptorProto name
             * @property {Array.<google.protobuf.IMethodDescriptorProto>|null} [method] ServiceDescriptorProto method
             * @property {google.protobuf.IServiceOptions|null} [options] ServiceDescriptorProto options
             */

            /**
             * Constructs a new ServiceDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents a ServiceDescriptorProto.
             * @implements IServiceDescriptorProto
             * @constructor
             * @param {google.protobuf.IServiceDescriptorProto=} [properties] Properties to set
             */
            function ServiceDescriptorProto(properties) {
                this.method = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * ServiceDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.ServiceDescriptorProto
             * @instance
             */
            ServiceDescriptorProto.prototype.name = "";

            /**
             * ServiceDescriptorProto method.
             * @member {Array.<google.protobuf.IMethodDescriptorProto>} method
             * @memberof google.protobuf.ServiceDescriptorProto
             * @instance
             */
            ServiceDescriptorProto.prototype.method = $util.emptyArray;

            /**
             * ServiceDescriptorProto options.
             * @member {google.protobuf.IServiceOptions|null|undefined} options
             * @memberof google.protobuf.ServiceDescriptorProto
             * @instance
             */
            ServiceDescriptorProto.prototype.options = null;

            /**
             * Creates a new ServiceDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {google.protobuf.IServiceDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.ServiceDescriptorProto} ServiceDescriptorProto instance
             */
            ServiceDescriptorProto.create = function create(properties) {
                return new ServiceDescriptorProto(properties);
            };

            /**
             * Encodes the specified ServiceDescriptorProto message. Does not implicitly {@link google.protobuf.ServiceDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {google.protobuf.IServiceDescriptorProto} message ServiceDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ServiceDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message.method != null && message.method.length)
                    for (var i = 0; i < message.method.length; ++i)
                        $root.google.protobuf.MethodDescriptorProto.encode(message.method[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.ServiceOptions.encode(message.options, writer.uint32(/* id 3, wireType 2 =*/26).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified ServiceDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.ServiceDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {google.protobuf.IServiceDescriptorProto} message ServiceDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ServiceDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a ServiceDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.ServiceDescriptorProto} ServiceDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ServiceDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.ServiceDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 2:
                        if (!(message.method && message.method.length))
                            message.method = [];
                        message.method.push($root.google.protobuf.MethodDescriptorProto.decode(reader, reader.uint32()));
                        break;
                    case 3:
                        message.options = $root.google.protobuf.ServiceOptions.decode(reader, reader.uint32());
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a ServiceDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.ServiceDescriptorProto} ServiceDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ServiceDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a ServiceDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            ServiceDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message.method != null && message.hasOwnProperty("method")) {
                    if (!Array.isArray(message.method))
                        return "method: array expected";
                    for (var i = 0; i < message.method.length; ++i) {
                        var error = $root.google.protobuf.MethodDescriptorProto.verify(message.method[i]);
                        if (error)
                            return "method." + error;
                    }
                }
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.ServiceOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                return null;
            };

            /**
             * Creates a ServiceDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.ServiceDescriptorProto} ServiceDescriptorProto
             */
            ServiceDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.ServiceDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.ServiceDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object.method) {
                    if (!Array.isArray(object.method))
                        throw TypeError(".google.protobuf.ServiceDescriptorProto.method: array expected");
                    message.method = [];
                    for (var i = 0; i < object.method.length; ++i) {
                        if (typeof object.method[i] !== "object")
                            throw TypeError(".google.protobuf.ServiceDescriptorProto.method: object expected");
                        message.method[i] = $root.google.protobuf.MethodDescriptorProto.fromObject(object.method[i]);
                    }
                }
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.ServiceDescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.ServiceOptions.fromObject(object.options);
                }
                return message;
            };

            /**
             * Creates a plain object from a ServiceDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.ServiceDescriptorProto
             * @static
             * @param {google.protobuf.ServiceDescriptorProto} message ServiceDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            ServiceDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.method = [];
                if (options.defaults) {
                    object.name = "";
                    object.options = null;
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message.method && message.method.length) {
                    object.method = [];
                    for (var j = 0; j < message.method.length; ++j)
                        object.method[j] = $root.google.protobuf.MethodDescriptorProto.toObject(message.method[j], options);
                }
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.ServiceOptions.toObject(message.options, options);
                return object;
            };

            /**
             * Converts this ServiceDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.ServiceDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            ServiceDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return ServiceDescriptorProto;
        })();

        protobuf.MethodDescriptorProto = (function() {

            /**
             * Properties of a MethodDescriptorProto.
             * @memberof google.protobuf
             * @interface IMethodDescriptorProto
             * @property {string|null} [name] MethodDescriptorProto name
             * @property {string|null} [inputType] MethodDescriptorProto inputType
             * @property {string|null} [outputType] MethodDescriptorProto outputType
             * @property {google.protobuf.IMethodOptions|null} [options] MethodDescriptorProto options
             * @property {boolean|null} [clientStreaming] MethodDescriptorProto clientStreaming
             * @property {boolean|null} [serverStreaming] MethodDescriptorProto serverStreaming
             */

            /**
             * Constructs a new MethodDescriptorProto.
             * @memberof google.protobuf
             * @classdesc Represents a MethodDescriptorProto.
             * @implements IMethodDescriptorProto
             * @constructor
             * @param {google.protobuf.IMethodDescriptorProto=} [properties] Properties to set
             */
            function MethodDescriptorProto(properties) {
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * MethodDescriptorProto name.
             * @member {string} name
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             */
            MethodDescriptorProto.prototype.name = "";

            /**
             * MethodDescriptorProto inputType.
             * @member {string} inputType
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             */
            MethodDescriptorProto.prototype.inputType = "";

            /**
             * MethodDescriptorProto outputType.
             * @member {string} outputType
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             */
            MethodDescriptorProto.prototype.outputType = "";

            /**
             * MethodDescriptorProto options.
             * @member {google.protobuf.IMethodOptions|null|undefined} options
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             */
            MethodDescriptorProto.prototype.options = null;

            /**
             * MethodDescriptorProto clientStreaming.
             * @member {boolean} clientStreaming
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             */
            MethodDescriptorProto.prototype.clientStreaming = false;

            /**
             * MethodDescriptorProto serverStreaming.
             * @member {boolean} serverStreaming
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             */
            MethodDescriptorProto.prototype.serverStreaming = false;

            /**
             * Creates a new MethodDescriptorProto instance using the specified properties.
             * @function create
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {google.protobuf.IMethodDescriptorProto=} [properties] Properties to set
             * @returns {google.protobuf.MethodDescriptorProto} MethodDescriptorProto instance
             */
            MethodDescriptorProto.create = function create(properties) {
                return new MethodDescriptorProto(properties);
            };

            /**
             * Encodes the specified MethodDescriptorProto message. Does not implicitly {@link google.protobuf.MethodDescriptorProto.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {google.protobuf.IMethodDescriptorProto} message MethodDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            MethodDescriptorProto.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.hasOwnProperty("name"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
                if (message.inputType != null && message.hasOwnProperty("inputType"))
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message.inputType);
                if (message.outputType != null && message.hasOwnProperty("outputType"))
                    writer.uint32(/* id 3, wireType 2 =*/26).string(message.outputType);
                if (message.options != null && message.hasOwnProperty("options"))
                    $root.google.protobuf.MethodOptions.encode(message.options, writer.uint32(/* id 4, wireType 2 =*/34).fork()).ldelim();
                if (message.clientStreaming != null && message.hasOwnProperty("clientStreaming"))
                    writer.uint32(/* id 5, wireType 0 =*/40).bool(message.clientStreaming);
                if (message.serverStreaming != null && message.hasOwnProperty("serverStreaming"))
                    writer.uint32(/* id 6, wireType 0 =*/48).bool(message.serverStreaming);
                return writer;
            };

            /**
             * Encodes the specified MethodDescriptorProto message, length delimited. Does not implicitly {@link google.protobuf.MethodDescriptorProto.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {google.protobuf.IMethodDescriptorProto} message MethodDescriptorProto message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            MethodDescriptorProto.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a MethodDescriptorProto message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.MethodDescriptorProto} MethodDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            MethodDescriptorProto.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.MethodDescriptorProto();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.name = reader.string();
                        break;
                    case 2:
                        message.inputType = reader.string();
                        break;
                    case 3:
                        message.outputType = reader.string();
                        break;
                    case 4:
                        message.options = $root.google.protobuf.MethodOptions.decode(reader, reader.uint32());
                        break;
                    case 5:
                        message.clientStreaming = reader.bool();
                        break;
                    case 6:
                        message.serverStreaming = reader.bool();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a MethodDescriptorProto message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.MethodDescriptorProto} MethodDescriptorProto
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            MethodDescriptorProto.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a MethodDescriptorProto message.
             * @function verify
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            MethodDescriptorProto.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name"))
                    if (!$util.isString(message.name))
                        return "name: string expected";
                if (message.inputType != null && message.hasOwnProperty("inputType"))
                    if (!$util.isString(message.inputType))
                        return "inputType: string expected";
                if (message.outputType != null && message.hasOwnProperty("outputType"))
                    if (!$util.isString(message.outputType))
                        return "outputType: string expected";
                if (message.options != null && message.hasOwnProperty("options")) {
                    var error = $root.google.protobuf.MethodOptions.verify(message.options);
                    if (error)
                        return "options." + error;
                }
                if (message.clientStreaming != null && message.hasOwnProperty("clientStreaming"))
                    if (typeof message.clientStreaming !== "boolean")
                        return "clientStreaming: boolean expected";
                if (message.serverStreaming != null && message.hasOwnProperty("serverStreaming"))
                    if (typeof message.serverStreaming !== "boolean")
                        return "serverStreaming: boolean expected";
                return null;
            };

            /**
             * Creates a MethodDescriptorProto message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.MethodDescriptorProto} MethodDescriptorProto
             */
            MethodDescriptorProto.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.MethodDescriptorProto)
                    return object;
                var message = new $root.google.protobuf.MethodDescriptorProto();
                if (object.name != null)
                    message.name = String(object.name);
                if (object.inputType != null)
                    message.inputType = String(object.inputType);
                if (object.outputType != null)
                    message.outputType = String(object.outputType);
                if (object.options != null) {
                    if (typeof object.options !== "object")
                        throw TypeError(".google.protobuf.MethodDescriptorProto.options: object expected");
                    message.options = $root.google.protobuf.MethodOptions.fromObject(object.options);
                }
                if (object.clientStreaming != null)
                    message.clientStreaming = Boolean(object.clientStreaming);
                if (object.serverStreaming != null)
                    message.serverStreaming = Boolean(object.serverStreaming);
                return message;
            };

            /**
             * Creates a plain object from a MethodDescriptorProto message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.MethodDescriptorProto
             * @static
             * @param {google.protobuf.MethodDescriptorProto} message MethodDescriptorProto
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            MethodDescriptorProto.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.defaults) {
                    object.name = "";
                    object.inputType = "";
                    object.outputType = "";
                    object.options = null;
                    object.clientStreaming = false;
                    object.serverStreaming = false;
                }
                if (message.name != null && message.hasOwnProperty("name"))
                    object.name = message.name;
                if (message.inputType != null && message.hasOwnProperty("inputType"))
                    object.inputType = message.inputType;
                if (message.outputType != null && message.hasOwnProperty("outputType"))
                    object.outputType = message.outputType;
                if (message.options != null && message.hasOwnProperty("options"))
                    object.options = $root.google.protobuf.MethodOptions.toObject(message.options, options);
                if (message.clientStreaming != null && message.hasOwnProperty("clientStreaming"))
                    object.clientStreaming = message.clientStreaming;
                if (message.serverStreaming != null && message.hasOwnProperty("serverStreaming"))
                    object.serverStreaming = message.serverStreaming;
                return object;
            };

            /**
             * Converts this MethodDescriptorProto to JSON.
             * @function toJSON
             * @memberof google.protobuf.MethodDescriptorProto
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            MethodDescriptorProto.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return MethodDescriptorProto;
        })();

        protobuf.FileOptions = (function() {

            /**
             * Properties of a FileOptions.
             * @memberof google.protobuf
             * @interface IFileOptions
             * @property {string|null} [javaPackage] FileOptions javaPackage
             * @property {string|null} [javaOuterClassname] FileOptions javaOuterClassname
             * @property {boolean|null} [javaMultipleFiles] FileOptions javaMultipleFiles
             * @property {boolean|null} [javaGenerateEqualsAndHash] FileOptions javaGenerateEqualsAndHash
             * @property {boolean|null} [javaStringCheckUtf8] FileOptions javaStringCheckUtf8
             * @property {google.protobuf.FileOptions.OptimizeMode|null} [optimizeFor] FileOptions optimizeFor
             * @property {string|null} [goPackage] FileOptions goPackage
             * @property {boolean|null} [ccGenericServices] FileOptions ccGenericServices
             * @property {boolean|null} [javaGenericServices] FileOptions javaGenericServices
             * @property {boolean|null} [pyGenericServices] FileOptions pyGenericServices
             * @property {boolean|null} [deprecated] FileOptions deprecated
             * @property {boolean|null} [ccEnableArenas] FileOptions ccEnableArenas
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] FileOptions uninterpretedOption
             */

            /**
             * Constructs a new FileOptions.
             * @memberof google.protobuf
             * @classdesc Represents a FileOptions.
             * @implements IFileOptions
             * @constructor
             * @param {google.protobuf.IFileOptions=} [properties] Properties to set
             */
            function FileOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * FileOptions javaPackage.
             * @member {string} javaPackage
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.javaPackage = "";

            /**
             * FileOptions javaOuterClassname.
             * @member {string} javaOuterClassname
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.javaOuterClassname = "";

            /**
             * FileOptions javaMultipleFiles.
             * @member {boolean} javaMultipleFiles
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.javaMultipleFiles = false;

            /**
             * FileOptions javaGenerateEqualsAndHash.
             * @member {boolean} javaGenerateEqualsAndHash
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.javaGenerateEqualsAndHash = false;

            /**
             * FileOptions javaStringCheckUtf8.
             * @member {boolean} javaStringCheckUtf8
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.javaStringCheckUtf8 = false;

            /**
             * FileOptions optimizeFor.
             * @member {google.protobuf.FileOptions.OptimizeMode} optimizeFor
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.optimizeFor = 1;

            /**
             * FileOptions goPackage.
             * @member {string} goPackage
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.goPackage = "";

            /**
             * FileOptions ccGenericServices.
             * @member {boolean} ccGenericServices
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.ccGenericServices = false;

            /**
             * FileOptions javaGenericServices.
             * @member {boolean} javaGenericServices
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.javaGenericServices = false;

            /**
             * FileOptions pyGenericServices.
             * @member {boolean} pyGenericServices
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.pyGenericServices = false;

            /**
             * FileOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.deprecated = false;

            /**
             * FileOptions ccEnableArenas.
             * @member {boolean} ccEnableArenas
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.ccEnableArenas = false;

            /**
             * FileOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.FileOptions
             * @instance
             */
            FileOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * Creates a new FileOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {google.protobuf.IFileOptions=} [properties] Properties to set
             * @returns {google.protobuf.FileOptions} FileOptions instance
             */
            FileOptions.create = function create(properties) {
                return new FileOptions(properties);
            };

            /**
             * Encodes the specified FileOptions message. Does not implicitly {@link google.protobuf.FileOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {google.protobuf.IFileOptions} message FileOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FileOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.javaPackage != null && message.hasOwnProperty("javaPackage"))
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.javaPackage);
                if (message.javaOuterClassname != null && message.hasOwnProperty("javaOuterClassname"))
                    writer.uint32(/* id 8, wireType 2 =*/66).string(message.javaOuterClassname);
                if (message.optimizeFor != null && message.hasOwnProperty("optimizeFor"))
                    writer.uint32(/* id 9, wireType 0 =*/72).int32(message.optimizeFor);
                if (message.javaMultipleFiles != null && message.hasOwnProperty("javaMultipleFiles"))
                    writer.uint32(/* id 10, wireType 0 =*/80).bool(message.javaMultipleFiles);
                if (message.goPackage != null && message.hasOwnProperty("goPackage"))
                    writer.uint32(/* id 11, wireType 2 =*/90).string(message.goPackage);
                if (message.ccGenericServices != null && message.hasOwnProperty("ccGenericServices"))
                    writer.uint32(/* id 16, wireType 0 =*/128).bool(message.ccGenericServices);
                if (message.javaGenericServices != null && message.hasOwnProperty("javaGenericServices"))
                    writer.uint32(/* id 17, wireType 0 =*/136).bool(message.javaGenericServices);
                if (message.pyGenericServices != null && message.hasOwnProperty("pyGenericServices"))
                    writer.uint32(/* id 18, wireType 0 =*/144).bool(message.pyGenericServices);
                if (message.javaGenerateEqualsAndHash != null && message.hasOwnProperty("javaGenerateEqualsAndHash"))
                    writer.uint32(/* id 20, wireType 0 =*/160).bool(message.javaGenerateEqualsAndHash);
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 23, wireType 0 =*/184).bool(message.deprecated);
                if (message.javaStringCheckUtf8 != null && message.hasOwnProperty("javaStringCheckUtf8"))
                    writer.uint32(/* id 27, wireType 0 =*/216).bool(message.javaStringCheckUtf8);
                if (message.ccEnableArenas != null && message.hasOwnProperty("ccEnableArenas"))
                    writer.uint32(/* id 31, wireType 0 =*/248).bool(message.ccEnableArenas);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified FileOptions message, length delimited. Does not implicitly {@link google.protobuf.FileOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {google.protobuf.IFileOptions} message FileOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FileOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a FileOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.FileOptions} FileOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FileOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.FileOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.javaPackage = reader.string();
                        break;
                    case 8:
                        message.javaOuterClassname = reader.string();
                        break;
                    case 10:
                        message.javaMultipleFiles = reader.bool();
                        break;
                    case 20:
                        message.javaGenerateEqualsAndHash = reader.bool();
                        break;
                    case 27:
                        message.javaStringCheckUtf8 = reader.bool();
                        break;
                    case 9:
                        message.optimizeFor = reader.int32();
                        break;
                    case 11:
                        message.goPackage = reader.string();
                        break;
                    case 16:
                        message.ccGenericServices = reader.bool();
                        break;
                    case 17:
                        message.javaGenericServices = reader.bool();
                        break;
                    case 18:
                        message.pyGenericServices = reader.bool();
                        break;
                    case 23:
                        message.deprecated = reader.bool();
                        break;
                    case 31:
                        message.ccEnableArenas = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a FileOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.FileOptions} FileOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FileOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a FileOptions message.
             * @function verify
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            FileOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.javaPackage != null && message.hasOwnProperty("javaPackage"))
                    if (!$util.isString(message.javaPackage))
                        return "javaPackage: string expected";
                if (message.javaOuterClassname != null && message.hasOwnProperty("javaOuterClassname"))
                    if (!$util.isString(message.javaOuterClassname))
                        return "javaOuterClassname: string expected";
                if (message.javaMultipleFiles != null && message.hasOwnProperty("javaMultipleFiles"))
                    if (typeof message.javaMultipleFiles !== "boolean")
                        return "javaMultipleFiles: boolean expected";
                if (message.javaGenerateEqualsAndHash != null && message.hasOwnProperty("javaGenerateEqualsAndHash"))
                    if (typeof message.javaGenerateEqualsAndHash !== "boolean")
                        return "javaGenerateEqualsAndHash: boolean expected";
                if (message.javaStringCheckUtf8 != null && message.hasOwnProperty("javaStringCheckUtf8"))
                    if (typeof message.javaStringCheckUtf8 !== "boolean")
                        return "javaStringCheckUtf8: boolean expected";
                if (message.optimizeFor != null && message.hasOwnProperty("optimizeFor"))
                    switch (message.optimizeFor) {
                    default:
                        return "optimizeFor: enum value expected";
                    case 1:
                    case 2:
                    case 3:
                        break;
                    }
                if (message.goPackage != null && message.hasOwnProperty("goPackage"))
                    if (!$util.isString(message.goPackage))
                        return "goPackage: string expected";
                if (message.ccGenericServices != null && message.hasOwnProperty("ccGenericServices"))
                    if (typeof message.ccGenericServices !== "boolean")
                        return "ccGenericServices: boolean expected";
                if (message.javaGenericServices != null && message.hasOwnProperty("javaGenericServices"))
                    if (typeof message.javaGenericServices !== "boolean")
                        return "javaGenericServices: boolean expected";
                if (message.pyGenericServices != null && message.hasOwnProperty("pyGenericServices"))
                    if (typeof message.pyGenericServices !== "boolean")
                        return "pyGenericServices: boolean expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.ccEnableArenas != null && message.hasOwnProperty("ccEnableArenas"))
                    if (typeof message.ccEnableArenas !== "boolean")
                        return "ccEnableArenas: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a FileOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.FileOptions} FileOptions
             */
            FileOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.FileOptions)
                    return object;
                var message = new $root.google.protobuf.FileOptions();
                if (object.javaPackage != null)
                    message.javaPackage = String(object.javaPackage);
                if (object.javaOuterClassname != null)
                    message.javaOuterClassname = String(object.javaOuterClassname);
                if (object.javaMultipleFiles != null)
                    message.javaMultipleFiles = Boolean(object.javaMultipleFiles);
                if (object.javaGenerateEqualsAndHash != null)
                    message.javaGenerateEqualsAndHash = Boolean(object.javaGenerateEqualsAndHash);
                if (object.javaStringCheckUtf8 != null)
                    message.javaStringCheckUtf8 = Boolean(object.javaStringCheckUtf8);
                switch (object.optimizeFor) {
                case "SPEED":
                case 1:
                    message.optimizeFor = 1;
                    break;
                case "CODE_SIZE":
                case 2:
                    message.optimizeFor = 2;
                    break;
                case "LITE_RUNTIME":
                case 3:
                    message.optimizeFor = 3;
                    break;
                }
                if (object.goPackage != null)
                    message.goPackage = String(object.goPackage);
                if (object.ccGenericServices != null)
                    message.ccGenericServices = Boolean(object.ccGenericServices);
                if (object.javaGenericServices != null)
                    message.javaGenericServices = Boolean(object.javaGenericServices);
                if (object.pyGenericServices != null)
                    message.pyGenericServices = Boolean(object.pyGenericServices);
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.ccEnableArenas != null)
                    message.ccEnableArenas = Boolean(object.ccEnableArenas);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.FileOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.FileOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a FileOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.FileOptions
             * @static
             * @param {google.protobuf.FileOptions} message FileOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            FileOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults) {
                    object.javaPackage = "";
                    object.javaOuterClassname = "";
                    object.optimizeFor = options.enums === String ? "SPEED" : 1;
                    object.javaMultipleFiles = false;
                    object.goPackage = "";
                    object.ccGenericServices = false;
                    object.javaGenericServices = false;
                    object.pyGenericServices = false;
                    object.javaGenerateEqualsAndHash = false;
                    object.deprecated = false;
                    object.javaStringCheckUtf8 = false;
                    object.ccEnableArenas = false;
                }
                if (message.javaPackage != null && message.hasOwnProperty("javaPackage"))
                    object.javaPackage = message.javaPackage;
                if (message.javaOuterClassname != null && message.hasOwnProperty("javaOuterClassname"))
                    object.javaOuterClassname = message.javaOuterClassname;
                if (message.optimizeFor != null && message.hasOwnProperty("optimizeFor"))
                    object.optimizeFor = options.enums === String ? $root.google.protobuf.FileOptions.OptimizeMode[message.optimizeFor] : message.optimizeFor;
                if (message.javaMultipleFiles != null && message.hasOwnProperty("javaMultipleFiles"))
                    object.javaMultipleFiles = message.javaMultipleFiles;
                if (message.goPackage != null && message.hasOwnProperty("goPackage"))
                    object.goPackage = message.goPackage;
                if (message.ccGenericServices != null && message.hasOwnProperty("ccGenericServices"))
                    object.ccGenericServices = message.ccGenericServices;
                if (message.javaGenericServices != null && message.hasOwnProperty("javaGenericServices"))
                    object.javaGenericServices = message.javaGenericServices;
                if (message.pyGenericServices != null && message.hasOwnProperty("pyGenericServices"))
                    object.pyGenericServices = message.pyGenericServices;
                if (message.javaGenerateEqualsAndHash != null && message.hasOwnProperty("javaGenerateEqualsAndHash"))
                    object.javaGenerateEqualsAndHash = message.javaGenerateEqualsAndHash;
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.javaStringCheckUtf8 != null && message.hasOwnProperty("javaStringCheckUtf8"))
                    object.javaStringCheckUtf8 = message.javaStringCheckUtf8;
                if (message.ccEnableArenas != null && message.hasOwnProperty("ccEnableArenas"))
                    object.ccEnableArenas = message.ccEnableArenas;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                return object;
            };

            /**
             * Converts this FileOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.FileOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            FileOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * OptimizeMode enum.
             * @name google.protobuf.FileOptions.OptimizeMode
             * @enum {string}
             * @property {number} SPEED=1 SPEED value
             * @property {number} CODE_SIZE=2 CODE_SIZE value
             * @property {number} LITE_RUNTIME=3 LITE_RUNTIME value
             */
            FileOptions.OptimizeMode = (function() {
                var valuesById = {}, values = Object.create(valuesById);
                values[valuesById[1] = "SPEED"] = 1;
                values[valuesById[2] = "CODE_SIZE"] = 2;
                values[valuesById[3] = "LITE_RUNTIME"] = 3;
                return values;
            })();

            return FileOptions;
        })();

        protobuf.MessageOptions = (function() {

            /**
             * Properties of a MessageOptions.
             * @memberof google.protobuf
             * @interface IMessageOptions
             * @property {boolean|null} [messageSetWireFormat] MessageOptions messageSetWireFormat
             * @property {boolean|null} [noStandardDescriptorAccessor] MessageOptions noStandardDescriptorAccessor
             * @property {boolean|null} [deprecated] MessageOptions deprecated
             * @property {boolean|null} [mapEntry] MessageOptions mapEntry
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] MessageOptions uninterpretedOption
             */

            /**
             * Constructs a new MessageOptions.
             * @memberof google.protobuf
             * @classdesc Represents a MessageOptions.
             * @implements IMessageOptions
             * @constructor
             * @param {google.protobuf.IMessageOptions=} [properties] Properties to set
             */
            function MessageOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * MessageOptions messageSetWireFormat.
             * @member {boolean} messageSetWireFormat
             * @memberof google.protobuf.MessageOptions
             * @instance
             */
            MessageOptions.prototype.messageSetWireFormat = false;

            /**
             * MessageOptions noStandardDescriptorAccessor.
             * @member {boolean} noStandardDescriptorAccessor
             * @memberof google.protobuf.MessageOptions
             * @instance
             */
            MessageOptions.prototype.noStandardDescriptorAccessor = false;

            /**
             * MessageOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.MessageOptions
             * @instance
             */
            MessageOptions.prototype.deprecated = false;

            /**
             * MessageOptions mapEntry.
             * @member {boolean} mapEntry
             * @memberof google.protobuf.MessageOptions
             * @instance
             */
            MessageOptions.prototype.mapEntry = false;

            /**
             * MessageOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.MessageOptions
             * @instance
             */
            MessageOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * Creates a new MessageOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {google.protobuf.IMessageOptions=} [properties] Properties to set
             * @returns {google.protobuf.MessageOptions} MessageOptions instance
             */
            MessageOptions.create = function create(properties) {
                return new MessageOptions(properties);
            };

            /**
             * Encodes the specified MessageOptions message. Does not implicitly {@link google.protobuf.MessageOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {google.protobuf.IMessageOptions} message MessageOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            MessageOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.messageSetWireFormat != null && message.hasOwnProperty("messageSetWireFormat"))
                    writer.uint32(/* id 1, wireType 0 =*/8).bool(message.messageSetWireFormat);
                if (message.noStandardDescriptorAccessor != null && message.hasOwnProperty("noStandardDescriptorAccessor"))
                    writer.uint32(/* id 2, wireType 0 =*/16).bool(message.noStandardDescriptorAccessor);
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 3, wireType 0 =*/24).bool(message.deprecated);
                if (message.mapEntry != null && message.hasOwnProperty("mapEntry"))
                    writer.uint32(/* id 7, wireType 0 =*/56).bool(message.mapEntry);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified MessageOptions message, length delimited. Does not implicitly {@link google.protobuf.MessageOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {google.protobuf.IMessageOptions} message MessageOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            MessageOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a MessageOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.MessageOptions} MessageOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            MessageOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.MessageOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.messageSetWireFormat = reader.bool();
                        break;
                    case 2:
                        message.noStandardDescriptorAccessor = reader.bool();
                        break;
                    case 3:
                        message.deprecated = reader.bool();
                        break;
                    case 7:
                        message.mapEntry = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a MessageOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.MessageOptions} MessageOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            MessageOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a MessageOptions message.
             * @function verify
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            MessageOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.messageSetWireFormat != null && message.hasOwnProperty("messageSetWireFormat"))
                    if (typeof message.messageSetWireFormat !== "boolean")
                        return "messageSetWireFormat: boolean expected";
                if (message.noStandardDescriptorAccessor != null && message.hasOwnProperty("noStandardDescriptorAccessor"))
                    if (typeof message.noStandardDescriptorAccessor !== "boolean")
                        return "noStandardDescriptorAccessor: boolean expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.mapEntry != null && message.hasOwnProperty("mapEntry"))
                    if (typeof message.mapEntry !== "boolean")
                        return "mapEntry: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a MessageOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.MessageOptions} MessageOptions
             */
            MessageOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.MessageOptions)
                    return object;
                var message = new $root.google.protobuf.MessageOptions();
                if (object.messageSetWireFormat != null)
                    message.messageSetWireFormat = Boolean(object.messageSetWireFormat);
                if (object.noStandardDescriptorAccessor != null)
                    message.noStandardDescriptorAccessor = Boolean(object.noStandardDescriptorAccessor);
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.mapEntry != null)
                    message.mapEntry = Boolean(object.mapEntry);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.MessageOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.MessageOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a MessageOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.MessageOptions
             * @static
             * @param {google.protobuf.MessageOptions} message MessageOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            MessageOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults) {
                    object.messageSetWireFormat = false;
                    object.noStandardDescriptorAccessor = false;
                    object.deprecated = false;
                    object.mapEntry = false;
                }
                if (message.messageSetWireFormat != null && message.hasOwnProperty("messageSetWireFormat"))
                    object.messageSetWireFormat = message.messageSetWireFormat;
                if (message.noStandardDescriptorAccessor != null && message.hasOwnProperty("noStandardDescriptorAccessor"))
                    object.noStandardDescriptorAccessor = message.noStandardDescriptorAccessor;
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.mapEntry != null && message.hasOwnProperty("mapEntry"))
                    object.mapEntry = message.mapEntry;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                return object;
            };

            /**
             * Converts this MessageOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.MessageOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            MessageOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return MessageOptions;
        })();

        protobuf.FieldOptions = (function() {

            /**
             * Properties of a FieldOptions.
             * @memberof google.protobuf
             * @interface IFieldOptions
             * @property {google.protobuf.FieldOptions.CType|null} [ctype] FieldOptions ctype
             * @property {boolean|null} [packed] FieldOptions packed
             * @property {boolean|null} [lazy] FieldOptions lazy
             * @property {boolean|null} [deprecated] FieldOptions deprecated
             * @property {boolean|null} [weak] FieldOptions weak
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] FieldOptions uninterpretedOption
             */

            /**
             * Constructs a new FieldOptions.
             * @memberof google.protobuf
             * @classdesc Represents a FieldOptions.
             * @implements IFieldOptions
             * @constructor
             * @param {google.protobuf.IFieldOptions=} [properties] Properties to set
             */
            function FieldOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * FieldOptions ctype.
             * @member {google.protobuf.FieldOptions.CType} ctype
             * @memberof google.protobuf.FieldOptions
             * @instance
             */
            FieldOptions.prototype.ctype = 0;

            /**
             * FieldOptions packed.
             * @member {boolean} packed
             * @memberof google.protobuf.FieldOptions
             * @instance
             */
            FieldOptions.prototype.packed = false;

            /**
             * FieldOptions lazy.
             * @member {boolean} lazy
             * @memberof google.protobuf.FieldOptions
             * @instance
             */
            FieldOptions.prototype.lazy = false;

            /**
             * FieldOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.FieldOptions
             * @instance
             */
            FieldOptions.prototype.deprecated = false;

            /**
             * FieldOptions weak.
             * @member {boolean} weak
             * @memberof google.protobuf.FieldOptions
             * @instance
             */
            FieldOptions.prototype.weak = false;

            /**
             * FieldOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.FieldOptions
             * @instance
             */
            FieldOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * Creates a new FieldOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {google.protobuf.IFieldOptions=} [properties] Properties to set
             * @returns {google.protobuf.FieldOptions} FieldOptions instance
             */
            FieldOptions.create = function create(properties) {
                return new FieldOptions(properties);
            };

            /**
             * Encodes the specified FieldOptions message. Does not implicitly {@link google.protobuf.FieldOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {google.protobuf.IFieldOptions} message FieldOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FieldOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.ctype != null && message.hasOwnProperty("ctype"))
                    writer.uint32(/* id 1, wireType 0 =*/8).int32(message.ctype);
                if (message.packed != null && message.hasOwnProperty("packed"))
                    writer.uint32(/* id 2, wireType 0 =*/16).bool(message.packed);
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 3, wireType 0 =*/24).bool(message.deprecated);
                if (message.lazy != null && message.hasOwnProperty("lazy"))
                    writer.uint32(/* id 5, wireType 0 =*/40).bool(message.lazy);
                if (message.weak != null && message.hasOwnProperty("weak"))
                    writer.uint32(/* id 10, wireType 0 =*/80).bool(message.weak);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified FieldOptions message, length delimited. Does not implicitly {@link google.protobuf.FieldOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {google.protobuf.IFieldOptions} message FieldOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            FieldOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a FieldOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.FieldOptions} FieldOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FieldOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.FieldOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.ctype = reader.int32();
                        break;
                    case 2:
                        message.packed = reader.bool();
                        break;
                    case 5:
                        message.lazy = reader.bool();
                        break;
                    case 3:
                        message.deprecated = reader.bool();
                        break;
                    case 10:
                        message.weak = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a FieldOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.FieldOptions} FieldOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            FieldOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a FieldOptions message.
             * @function verify
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            FieldOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.ctype != null && message.hasOwnProperty("ctype"))
                    switch (message.ctype) {
                    default:
                        return "ctype: enum value expected";
                    case 0:
                    case 1:
                    case 2:
                        break;
                    }
                if (message.packed != null && message.hasOwnProperty("packed"))
                    if (typeof message.packed !== "boolean")
                        return "packed: boolean expected";
                if (message.lazy != null && message.hasOwnProperty("lazy"))
                    if (typeof message.lazy !== "boolean")
                        return "lazy: boolean expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.weak != null && message.hasOwnProperty("weak"))
                    if (typeof message.weak !== "boolean")
                        return "weak: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a FieldOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.FieldOptions} FieldOptions
             */
            FieldOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.FieldOptions)
                    return object;
                var message = new $root.google.protobuf.FieldOptions();
                switch (object.ctype) {
                case "STRING":
                case 0:
                    message.ctype = 0;
                    break;
                case "CORD":
                case 1:
                    message.ctype = 1;
                    break;
                case "STRING_PIECE":
                case 2:
                    message.ctype = 2;
                    break;
                }
                if (object.packed != null)
                    message.packed = Boolean(object.packed);
                if (object.lazy != null)
                    message.lazy = Boolean(object.lazy);
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.weak != null)
                    message.weak = Boolean(object.weak);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.FieldOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.FieldOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a FieldOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.FieldOptions
             * @static
             * @param {google.protobuf.FieldOptions} message FieldOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            FieldOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults) {
                    object.ctype = options.enums === String ? "STRING" : 0;
                    object.packed = false;
                    object.deprecated = false;
                    object.lazy = false;
                    object.weak = false;
                }
                if (message.ctype != null && message.hasOwnProperty("ctype"))
                    object.ctype = options.enums === String ? $root.google.protobuf.FieldOptions.CType[message.ctype] : message.ctype;
                if (message.packed != null && message.hasOwnProperty("packed"))
                    object.packed = message.packed;
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.lazy != null && message.hasOwnProperty("lazy"))
                    object.lazy = message.lazy;
                if (message.weak != null && message.hasOwnProperty("weak"))
                    object.weak = message.weak;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                return object;
            };

            /**
             * Converts this FieldOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.FieldOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            FieldOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            /**
             * CType enum.
             * @name google.protobuf.FieldOptions.CType
             * @enum {string}
             * @property {number} STRING=0 STRING value
             * @property {number} CORD=1 CORD value
             * @property {number} STRING_PIECE=2 STRING_PIECE value
             */
            FieldOptions.CType = (function() {
                var valuesById = {}, values = Object.create(valuesById);
                values[valuesById[0] = "STRING"] = 0;
                values[valuesById[1] = "CORD"] = 1;
                values[valuesById[2] = "STRING_PIECE"] = 2;
                return values;
            })();

            return FieldOptions;
        })();

        protobuf.EnumOptions = (function() {

            /**
             * Properties of an EnumOptions.
             * @memberof google.protobuf
             * @interface IEnumOptions
             * @property {boolean|null} [allowAlias] EnumOptions allowAlias
             * @property {boolean|null} [deprecated] EnumOptions deprecated
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] EnumOptions uninterpretedOption
             */

            /**
             * Constructs a new EnumOptions.
             * @memberof google.protobuf
             * @classdesc Represents an EnumOptions.
             * @implements IEnumOptions
             * @constructor
             * @param {google.protobuf.IEnumOptions=} [properties] Properties to set
             */
            function EnumOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * EnumOptions allowAlias.
             * @member {boolean} allowAlias
             * @memberof google.protobuf.EnumOptions
             * @instance
             */
            EnumOptions.prototype.allowAlias = false;

            /**
             * EnumOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.EnumOptions
             * @instance
             */
            EnumOptions.prototype.deprecated = false;

            /**
             * EnumOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.EnumOptions
             * @instance
             */
            EnumOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * Creates a new EnumOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {google.protobuf.IEnumOptions=} [properties] Properties to set
             * @returns {google.protobuf.EnumOptions} EnumOptions instance
             */
            EnumOptions.create = function create(properties) {
                return new EnumOptions(properties);
            };

            /**
             * Encodes the specified EnumOptions message. Does not implicitly {@link google.protobuf.EnumOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {google.protobuf.IEnumOptions} message EnumOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.allowAlias != null && message.hasOwnProperty("allowAlias"))
                    writer.uint32(/* id 2, wireType 0 =*/16).bool(message.allowAlias);
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 3, wireType 0 =*/24).bool(message.deprecated);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified EnumOptions message, length delimited. Does not implicitly {@link google.protobuf.EnumOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {google.protobuf.IEnumOptions} message EnumOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an EnumOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.EnumOptions} EnumOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.EnumOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 2:
                        message.allowAlias = reader.bool();
                        break;
                    case 3:
                        message.deprecated = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an EnumOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.EnumOptions} EnumOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an EnumOptions message.
             * @function verify
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            EnumOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.allowAlias != null && message.hasOwnProperty("allowAlias"))
                    if (typeof message.allowAlias !== "boolean")
                        return "allowAlias: boolean expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                return null;
            };

            /**
             * Creates an EnumOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.EnumOptions} EnumOptions
             */
            EnumOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.EnumOptions)
                    return object;
                var message = new $root.google.protobuf.EnumOptions();
                if (object.allowAlias != null)
                    message.allowAlias = Boolean(object.allowAlias);
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.EnumOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.EnumOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from an EnumOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.EnumOptions
             * @static
             * @param {google.protobuf.EnumOptions} message EnumOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            EnumOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults) {
                    object.allowAlias = false;
                    object.deprecated = false;
                }
                if (message.allowAlias != null && message.hasOwnProperty("allowAlias"))
                    object.allowAlias = message.allowAlias;
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                return object;
            };

            /**
             * Converts this EnumOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.EnumOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            EnumOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return EnumOptions;
        })();

        protobuf.EnumValueOptions = (function() {

            /**
             * Properties of an EnumValueOptions.
             * @memberof google.protobuf
             * @interface IEnumValueOptions
             * @property {boolean|null} [deprecated] EnumValueOptions deprecated
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] EnumValueOptions uninterpretedOption
             * @property {boolean|null} [".wireIn"] EnumValueOptions .wireIn
             * @property {boolean|null} [".wireOut"] EnumValueOptions .wireOut
             * @property {boolean|null} [".wireDebugIn"] EnumValueOptions .wireDebugIn
             * @property {boolean|null} [".wireDebugOut"] EnumValueOptions .wireDebugOut
             * @property {boolean|null} [".wireTiny"] EnumValueOptions .wireTiny
             * @property {boolean|null} [".wireBootloader"] EnumValueOptions .wireBootloader
             */

            /**
             * Constructs a new EnumValueOptions.
             * @memberof google.protobuf
             * @classdesc Represents an EnumValueOptions.
             * @implements IEnumValueOptions
             * @constructor
             * @param {google.protobuf.IEnumValueOptions=} [properties] Properties to set
             */
            function EnumValueOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * EnumValueOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype.deprecated = false;

            /**
             * EnumValueOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * EnumValueOptions .wireIn.
             * @member {boolean} .wireIn
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype[".wireIn"] = false;

            /**
             * EnumValueOptions .wireOut.
             * @member {boolean} .wireOut
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype[".wireOut"] = false;

            /**
             * EnumValueOptions .wireDebugIn.
             * @member {boolean} .wireDebugIn
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype[".wireDebugIn"] = false;

            /**
             * EnumValueOptions .wireDebugOut.
             * @member {boolean} .wireDebugOut
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype[".wireDebugOut"] = false;

            /**
             * EnumValueOptions .wireTiny.
             * @member {boolean} .wireTiny
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype[".wireTiny"] = false;

            /**
             * EnumValueOptions .wireBootloader.
             * @member {boolean} .wireBootloader
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             */
            EnumValueOptions.prototype[".wireBootloader"] = false;

            /**
             * Creates a new EnumValueOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {google.protobuf.IEnumValueOptions=} [properties] Properties to set
             * @returns {google.protobuf.EnumValueOptions} EnumValueOptions instance
             */
            EnumValueOptions.create = function create(properties) {
                return new EnumValueOptions(properties);
            };

            /**
             * Encodes the specified EnumValueOptions message. Does not implicitly {@link google.protobuf.EnumValueOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {google.protobuf.IEnumValueOptions} message EnumValueOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumValueOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 1, wireType 0 =*/8).bool(message.deprecated);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                if (message[".wireIn"] != null && message.hasOwnProperty(".wireIn"))
                    writer.uint32(/* id 50002, wireType 0 =*/400016).bool(message[".wireIn"]);
                if (message[".wireOut"] != null && message.hasOwnProperty(".wireOut"))
                    writer.uint32(/* id 50003, wireType 0 =*/400024).bool(message[".wireOut"]);
                if (message[".wireDebugIn"] != null && message.hasOwnProperty(".wireDebugIn"))
                    writer.uint32(/* id 50004, wireType 0 =*/400032).bool(message[".wireDebugIn"]);
                if (message[".wireDebugOut"] != null && message.hasOwnProperty(".wireDebugOut"))
                    writer.uint32(/* id 50005, wireType 0 =*/400040).bool(message[".wireDebugOut"]);
                if (message[".wireTiny"] != null && message.hasOwnProperty(".wireTiny"))
                    writer.uint32(/* id 50006, wireType 0 =*/400048).bool(message[".wireTiny"]);
                if (message[".wireBootloader"] != null && message.hasOwnProperty(".wireBootloader"))
                    writer.uint32(/* id 50007, wireType 0 =*/400056).bool(message[".wireBootloader"]);
                return writer;
            };

            /**
             * Encodes the specified EnumValueOptions message, length delimited. Does not implicitly {@link google.protobuf.EnumValueOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {google.protobuf.IEnumValueOptions} message EnumValueOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            EnumValueOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an EnumValueOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.EnumValueOptions} EnumValueOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumValueOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.EnumValueOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        message.deprecated = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    case 50002:
                        message[".wireIn"] = reader.bool();
                        break;
                    case 50003:
                        message[".wireOut"] = reader.bool();
                        break;
                    case 50004:
                        message[".wireDebugIn"] = reader.bool();
                        break;
                    case 50005:
                        message[".wireDebugOut"] = reader.bool();
                        break;
                    case 50006:
                        message[".wireTiny"] = reader.bool();
                        break;
                    case 50007:
                        message[".wireBootloader"] = reader.bool();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an EnumValueOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.EnumValueOptions} EnumValueOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            EnumValueOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an EnumValueOptions message.
             * @function verify
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            EnumValueOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                if (message[".wireIn"] != null && message.hasOwnProperty(".wireIn"))
                    if (typeof message[".wireIn"] !== "boolean")
                        return ".wireIn: boolean expected";
                if (message[".wireOut"] != null && message.hasOwnProperty(".wireOut"))
                    if (typeof message[".wireOut"] !== "boolean")
                        return ".wireOut: boolean expected";
                if (message[".wireDebugIn"] != null && message.hasOwnProperty(".wireDebugIn"))
                    if (typeof message[".wireDebugIn"] !== "boolean")
                        return ".wireDebugIn: boolean expected";
                if (message[".wireDebugOut"] != null && message.hasOwnProperty(".wireDebugOut"))
                    if (typeof message[".wireDebugOut"] !== "boolean")
                        return ".wireDebugOut: boolean expected";
                if (message[".wireTiny"] != null && message.hasOwnProperty(".wireTiny"))
                    if (typeof message[".wireTiny"] !== "boolean")
                        return ".wireTiny: boolean expected";
                if (message[".wireBootloader"] != null && message.hasOwnProperty(".wireBootloader"))
                    if (typeof message[".wireBootloader"] !== "boolean")
                        return ".wireBootloader: boolean expected";
                return null;
            };

            /**
             * Creates an EnumValueOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.EnumValueOptions} EnumValueOptions
             */
            EnumValueOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.EnumValueOptions)
                    return object;
                var message = new $root.google.protobuf.EnumValueOptions();
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.EnumValueOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.EnumValueOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                if (object[".wireIn"] != null)
                    message[".wireIn"] = Boolean(object[".wireIn"]);
                if (object[".wireOut"] != null)
                    message[".wireOut"] = Boolean(object[".wireOut"]);
                if (object[".wireDebugIn"] != null)
                    message[".wireDebugIn"] = Boolean(object[".wireDebugIn"]);
                if (object[".wireDebugOut"] != null)
                    message[".wireDebugOut"] = Boolean(object[".wireDebugOut"]);
                if (object[".wireTiny"] != null)
                    message[".wireTiny"] = Boolean(object[".wireTiny"]);
                if (object[".wireBootloader"] != null)
                    message[".wireBootloader"] = Boolean(object[".wireBootloader"]);
                return message;
            };

            /**
             * Creates a plain object from an EnumValueOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.EnumValueOptions
             * @static
             * @param {google.protobuf.EnumValueOptions} message EnumValueOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            EnumValueOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults) {
                    object.deprecated = false;
                    object[".wireIn"] = false;
                    object[".wireOut"] = false;
                    object[".wireDebugIn"] = false;
                    object[".wireDebugOut"] = false;
                    object[".wireTiny"] = false;
                    object[".wireBootloader"] = false;
                }
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                if (message[".wireIn"] != null && message.hasOwnProperty(".wireIn"))
                    object[".wireIn"] = message[".wireIn"];
                if (message[".wireOut"] != null && message.hasOwnProperty(".wireOut"))
                    object[".wireOut"] = message[".wireOut"];
                if (message[".wireDebugIn"] != null && message.hasOwnProperty(".wireDebugIn"))
                    object[".wireDebugIn"] = message[".wireDebugIn"];
                if (message[".wireDebugOut"] != null && message.hasOwnProperty(".wireDebugOut"))
                    object[".wireDebugOut"] = message[".wireDebugOut"];
                if (message[".wireTiny"] != null && message.hasOwnProperty(".wireTiny"))
                    object[".wireTiny"] = message[".wireTiny"];
                if (message[".wireBootloader"] != null && message.hasOwnProperty(".wireBootloader"))
                    object[".wireBootloader"] = message[".wireBootloader"];
                return object;
            };

            /**
             * Converts this EnumValueOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.EnumValueOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            EnumValueOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return EnumValueOptions;
        })();

        protobuf.ServiceOptions = (function() {

            /**
             * Properties of a ServiceOptions.
             * @memberof google.protobuf
             * @interface IServiceOptions
             * @property {boolean|null} [deprecated] ServiceOptions deprecated
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] ServiceOptions uninterpretedOption
             */

            /**
             * Constructs a new ServiceOptions.
             * @memberof google.protobuf
             * @classdesc Represents a ServiceOptions.
             * @implements IServiceOptions
             * @constructor
             * @param {google.protobuf.IServiceOptions=} [properties] Properties to set
             */
            function ServiceOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * ServiceOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.ServiceOptions
             * @instance
             */
            ServiceOptions.prototype.deprecated = false;

            /**
             * ServiceOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.ServiceOptions
             * @instance
             */
            ServiceOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * Creates a new ServiceOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {google.protobuf.IServiceOptions=} [properties] Properties to set
             * @returns {google.protobuf.ServiceOptions} ServiceOptions instance
             */
            ServiceOptions.create = function create(properties) {
                return new ServiceOptions(properties);
            };

            /**
             * Encodes the specified ServiceOptions message. Does not implicitly {@link google.protobuf.ServiceOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {google.protobuf.IServiceOptions} message ServiceOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ServiceOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 33, wireType 0 =*/264).bool(message.deprecated);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified ServiceOptions message, length delimited. Does not implicitly {@link google.protobuf.ServiceOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {google.protobuf.IServiceOptions} message ServiceOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            ServiceOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a ServiceOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.ServiceOptions} ServiceOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ServiceOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.ServiceOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 33:
                        message.deprecated = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a ServiceOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.ServiceOptions} ServiceOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            ServiceOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a ServiceOptions message.
             * @function verify
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            ServiceOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a ServiceOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.ServiceOptions} ServiceOptions
             */
            ServiceOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.ServiceOptions)
                    return object;
                var message = new $root.google.protobuf.ServiceOptions();
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.ServiceOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.ServiceOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a ServiceOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.ServiceOptions
             * @static
             * @param {google.protobuf.ServiceOptions} message ServiceOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            ServiceOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults)
                    object.deprecated = false;
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                return object;
            };

            /**
             * Converts this ServiceOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.ServiceOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            ServiceOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return ServiceOptions;
        })();

        protobuf.MethodOptions = (function() {

            /**
             * Properties of a MethodOptions.
             * @memberof google.protobuf
             * @interface IMethodOptions
             * @property {boolean|null} [deprecated] MethodOptions deprecated
             * @property {Array.<google.protobuf.IUninterpretedOption>|null} [uninterpretedOption] MethodOptions uninterpretedOption
             */

            /**
             * Constructs a new MethodOptions.
             * @memberof google.protobuf
             * @classdesc Represents a MethodOptions.
             * @implements IMethodOptions
             * @constructor
             * @param {google.protobuf.IMethodOptions=} [properties] Properties to set
             */
            function MethodOptions(properties) {
                this.uninterpretedOption = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * MethodOptions deprecated.
             * @member {boolean} deprecated
             * @memberof google.protobuf.MethodOptions
             * @instance
             */
            MethodOptions.prototype.deprecated = false;

            /**
             * MethodOptions uninterpretedOption.
             * @member {Array.<google.protobuf.IUninterpretedOption>} uninterpretedOption
             * @memberof google.protobuf.MethodOptions
             * @instance
             */
            MethodOptions.prototype.uninterpretedOption = $util.emptyArray;

            /**
             * Creates a new MethodOptions instance using the specified properties.
             * @function create
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {google.protobuf.IMethodOptions=} [properties] Properties to set
             * @returns {google.protobuf.MethodOptions} MethodOptions instance
             */
            MethodOptions.create = function create(properties) {
                return new MethodOptions(properties);
            };

            /**
             * Encodes the specified MethodOptions message. Does not implicitly {@link google.protobuf.MethodOptions.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {google.protobuf.IMethodOptions} message MethodOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            MethodOptions.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    writer.uint32(/* id 33, wireType 0 =*/264).bool(message.deprecated);
                if (message.uninterpretedOption != null && message.uninterpretedOption.length)
                    for (var i = 0; i < message.uninterpretedOption.length; ++i)
                        $root.google.protobuf.UninterpretedOption.encode(message.uninterpretedOption[i], writer.uint32(/* id 999, wireType 2 =*/7994).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified MethodOptions message, length delimited. Does not implicitly {@link google.protobuf.MethodOptions.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {google.protobuf.IMethodOptions} message MethodOptions message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            MethodOptions.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a MethodOptions message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.MethodOptions} MethodOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            MethodOptions.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.MethodOptions();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 33:
                        message.deprecated = reader.bool();
                        break;
                    case 999:
                        if (!(message.uninterpretedOption && message.uninterpretedOption.length))
                            message.uninterpretedOption = [];
                        message.uninterpretedOption.push($root.google.protobuf.UninterpretedOption.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a MethodOptions message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.MethodOptions} MethodOptions
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            MethodOptions.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a MethodOptions message.
             * @function verify
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            MethodOptions.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    if (typeof message.deprecated !== "boolean")
                        return "deprecated: boolean expected";
                if (message.uninterpretedOption != null && message.hasOwnProperty("uninterpretedOption")) {
                    if (!Array.isArray(message.uninterpretedOption))
                        return "uninterpretedOption: array expected";
                    for (var i = 0; i < message.uninterpretedOption.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.verify(message.uninterpretedOption[i]);
                        if (error)
                            return "uninterpretedOption." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a MethodOptions message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.MethodOptions} MethodOptions
             */
            MethodOptions.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.MethodOptions)
                    return object;
                var message = new $root.google.protobuf.MethodOptions();
                if (object.deprecated != null)
                    message.deprecated = Boolean(object.deprecated);
                if (object.uninterpretedOption) {
                    if (!Array.isArray(object.uninterpretedOption))
                        throw TypeError(".google.protobuf.MethodOptions.uninterpretedOption: array expected");
                    message.uninterpretedOption = [];
                    for (var i = 0; i < object.uninterpretedOption.length; ++i) {
                        if (typeof object.uninterpretedOption[i] !== "object")
                            throw TypeError(".google.protobuf.MethodOptions.uninterpretedOption: object expected");
                        message.uninterpretedOption[i] = $root.google.protobuf.UninterpretedOption.fromObject(object.uninterpretedOption[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a MethodOptions message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.MethodOptions
             * @static
             * @param {google.protobuf.MethodOptions} message MethodOptions
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            MethodOptions.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.uninterpretedOption = [];
                if (options.defaults)
                    object.deprecated = false;
                if (message.deprecated != null && message.hasOwnProperty("deprecated"))
                    object.deprecated = message.deprecated;
                if (message.uninterpretedOption && message.uninterpretedOption.length) {
                    object.uninterpretedOption = [];
                    for (var j = 0; j < message.uninterpretedOption.length; ++j)
                        object.uninterpretedOption[j] = $root.google.protobuf.UninterpretedOption.toObject(message.uninterpretedOption[j], options);
                }
                return object;
            };

            /**
             * Converts this MethodOptions to JSON.
             * @function toJSON
             * @memberof google.protobuf.MethodOptions
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            MethodOptions.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            return MethodOptions;
        })();

        protobuf.UninterpretedOption = (function() {

            /**
             * Properties of an UninterpretedOption.
             * @memberof google.protobuf
             * @interface IUninterpretedOption
             * @property {Array.<google.protobuf.UninterpretedOption.INamePart>|null} [name] UninterpretedOption name
             * @property {string|null} [identifierValue] UninterpretedOption identifierValue
             * @property {number|Long|null} [positiveIntValue] UninterpretedOption positiveIntValue
             * @property {number|Long|null} [negativeIntValue] UninterpretedOption negativeIntValue
             * @property {number|null} [doubleValue] UninterpretedOption doubleValue
             * @property {Uint8Array|null} [stringValue] UninterpretedOption stringValue
             * @property {string|null} [aggregateValue] UninterpretedOption aggregateValue
             */

            /**
             * Constructs a new UninterpretedOption.
             * @memberof google.protobuf
             * @classdesc Represents an UninterpretedOption.
             * @implements IUninterpretedOption
             * @constructor
             * @param {google.protobuf.IUninterpretedOption=} [properties] Properties to set
             */
            function UninterpretedOption(properties) {
                this.name = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * UninterpretedOption name.
             * @member {Array.<google.protobuf.UninterpretedOption.INamePart>} name
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.name = $util.emptyArray;

            /**
             * UninterpretedOption identifierValue.
             * @member {string} identifierValue
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.identifierValue = "";

            /**
             * UninterpretedOption positiveIntValue.
             * @member {number|Long} positiveIntValue
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.positiveIntValue = $util.Long ? $util.Long.fromBits(0,0,true) : 0;

            /**
             * UninterpretedOption negativeIntValue.
             * @member {number|Long} negativeIntValue
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.negativeIntValue = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

            /**
             * UninterpretedOption doubleValue.
             * @member {number} doubleValue
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.doubleValue = 0;

            /**
             * UninterpretedOption stringValue.
             * @member {Uint8Array} stringValue
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.stringValue = $util.newBuffer([]);

            /**
             * UninterpretedOption aggregateValue.
             * @member {string} aggregateValue
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             */
            UninterpretedOption.prototype.aggregateValue = "";

            /**
             * Creates a new UninterpretedOption instance using the specified properties.
             * @function create
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {google.protobuf.IUninterpretedOption=} [properties] Properties to set
             * @returns {google.protobuf.UninterpretedOption} UninterpretedOption instance
             */
            UninterpretedOption.create = function create(properties) {
                return new UninterpretedOption(properties);
            };

            /**
             * Encodes the specified UninterpretedOption message. Does not implicitly {@link google.protobuf.UninterpretedOption.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {google.protobuf.IUninterpretedOption} message UninterpretedOption message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            UninterpretedOption.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.name != null && message.name.length)
                    for (var i = 0; i < message.name.length; ++i)
                        $root.google.protobuf.UninterpretedOption.NamePart.encode(message.name[i], writer.uint32(/* id 2, wireType 2 =*/18).fork()).ldelim();
                if (message.identifierValue != null && message.hasOwnProperty("identifierValue"))
                    writer.uint32(/* id 3, wireType 2 =*/26).string(message.identifierValue);
                if (message.positiveIntValue != null && message.hasOwnProperty("positiveIntValue"))
                    writer.uint32(/* id 4, wireType 0 =*/32).uint64(message.positiveIntValue);
                if (message.negativeIntValue != null && message.hasOwnProperty("negativeIntValue"))
                    writer.uint32(/* id 5, wireType 0 =*/40).int64(message.negativeIntValue);
                if (message.doubleValue != null && message.hasOwnProperty("doubleValue"))
                    writer.uint32(/* id 6, wireType 1 =*/49).double(message.doubleValue);
                if (message.stringValue != null && message.hasOwnProperty("stringValue"))
                    writer.uint32(/* id 7, wireType 2 =*/58).bytes(message.stringValue);
                if (message.aggregateValue != null && message.hasOwnProperty("aggregateValue"))
                    writer.uint32(/* id 8, wireType 2 =*/66).string(message.aggregateValue);
                return writer;
            };

            /**
             * Encodes the specified UninterpretedOption message, length delimited. Does not implicitly {@link google.protobuf.UninterpretedOption.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {google.protobuf.IUninterpretedOption} message UninterpretedOption message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            UninterpretedOption.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes an UninterpretedOption message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.UninterpretedOption} UninterpretedOption
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            UninterpretedOption.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.UninterpretedOption();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 2:
                        if (!(message.name && message.name.length))
                            message.name = [];
                        message.name.push($root.google.protobuf.UninterpretedOption.NamePart.decode(reader, reader.uint32()));
                        break;
                    case 3:
                        message.identifierValue = reader.string();
                        break;
                    case 4:
                        message.positiveIntValue = reader.uint64();
                        break;
                    case 5:
                        message.negativeIntValue = reader.int64();
                        break;
                    case 6:
                        message.doubleValue = reader.double();
                        break;
                    case 7:
                        message.stringValue = reader.bytes();
                        break;
                    case 8:
                        message.aggregateValue = reader.string();
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes an UninterpretedOption message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.UninterpretedOption} UninterpretedOption
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            UninterpretedOption.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies an UninterpretedOption message.
             * @function verify
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            UninterpretedOption.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.name != null && message.hasOwnProperty("name")) {
                    if (!Array.isArray(message.name))
                        return "name: array expected";
                    for (var i = 0; i < message.name.length; ++i) {
                        var error = $root.google.protobuf.UninterpretedOption.NamePart.verify(message.name[i]);
                        if (error)
                            return "name." + error;
                    }
                }
                if (message.identifierValue != null && message.hasOwnProperty("identifierValue"))
                    if (!$util.isString(message.identifierValue))
                        return "identifierValue: string expected";
                if (message.positiveIntValue != null && message.hasOwnProperty("positiveIntValue"))
                    if (!$util.isInteger(message.positiveIntValue) && !(message.positiveIntValue && $util.isInteger(message.positiveIntValue.low) && $util.isInteger(message.positiveIntValue.high)))
                        return "positiveIntValue: integer|Long expected";
                if (message.negativeIntValue != null && message.hasOwnProperty("negativeIntValue"))
                    if (!$util.isInteger(message.negativeIntValue) && !(message.negativeIntValue && $util.isInteger(message.negativeIntValue.low) && $util.isInteger(message.negativeIntValue.high)))
                        return "negativeIntValue: integer|Long expected";
                if (message.doubleValue != null && message.hasOwnProperty("doubleValue"))
                    if (typeof message.doubleValue !== "number")
                        return "doubleValue: number expected";
                if (message.stringValue != null && message.hasOwnProperty("stringValue"))
                    if (!(message.stringValue && typeof message.stringValue.length === "number" || $util.isString(message.stringValue)))
                        return "stringValue: buffer expected";
                if (message.aggregateValue != null && message.hasOwnProperty("aggregateValue"))
                    if (!$util.isString(message.aggregateValue))
                        return "aggregateValue: string expected";
                return null;
            };

            /**
             * Creates an UninterpretedOption message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.UninterpretedOption} UninterpretedOption
             */
            UninterpretedOption.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.UninterpretedOption)
                    return object;
                var message = new $root.google.protobuf.UninterpretedOption();
                if (object.name) {
                    if (!Array.isArray(object.name))
                        throw TypeError(".google.protobuf.UninterpretedOption.name: array expected");
                    message.name = [];
                    for (var i = 0; i < object.name.length; ++i) {
                        if (typeof object.name[i] !== "object")
                            throw TypeError(".google.protobuf.UninterpretedOption.name: object expected");
                        message.name[i] = $root.google.protobuf.UninterpretedOption.NamePart.fromObject(object.name[i]);
                    }
                }
                if (object.identifierValue != null)
                    message.identifierValue = String(object.identifierValue);
                if (object.positiveIntValue != null)
                    if ($util.Long)
                        (message.positiveIntValue = $util.Long.fromValue(object.positiveIntValue)).unsigned = true;
                    else if (typeof object.positiveIntValue === "string")
                        message.positiveIntValue = parseInt(object.positiveIntValue, 10);
                    else if (typeof object.positiveIntValue === "number")
                        message.positiveIntValue = object.positiveIntValue;
                    else if (typeof object.positiveIntValue === "object")
                        message.positiveIntValue = new $util.LongBits(object.positiveIntValue.low >>> 0, object.positiveIntValue.high >>> 0).toNumber(true);
                if (object.negativeIntValue != null)
                    if ($util.Long)
                        (message.negativeIntValue = $util.Long.fromValue(object.negativeIntValue)).unsigned = false;
                    else if (typeof object.negativeIntValue === "string")
                        message.negativeIntValue = parseInt(object.negativeIntValue, 10);
                    else if (typeof object.negativeIntValue === "number")
                        message.negativeIntValue = object.negativeIntValue;
                    else if (typeof object.negativeIntValue === "object")
                        message.negativeIntValue = new $util.LongBits(object.negativeIntValue.low >>> 0, object.negativeIntValue.high >>> 0).toNumber();
                if (object.doubleValue != null)
                    message.doubleValue = Number(object.doubleValue);
                if (object.stringValue != null)
                    if (typeof object.stringValue === "string")
                        $util.base64.decode(object.stringValue, message.stringValue = $util.newBuffer($util.base64.length(object.stringValue)), 0);
                    else if (object.stringValue.length)
                        message.stringValue = object.stringValue;
                if (object.aggregateValue != null)
                    message.aggregateValue = String(object.aggregateValue);
                return message;
            };

            /**
             * Creates a plain object from an UninterpretedOption message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.UninterpretedOption
             * @static
             * @param {google.protobuf.UninterpretedOption} message UninterpretedOption
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            UninterpretedOption.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.name = [];
                if (options.defaults) {
                    object.identifierValue = "";
                    if ($util.Long) {
                        var long = new $util.Long(0, 0, true);
                        object.positiveIntValue = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                    } else
                        object.positiveIntValue = options.longs === String ? "0" : 0;
                    if ($util.Long) {
                        var long = new $util.Long(0, 0, false);
                        object.negativeIntValue = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                    } else
                        object.negativeIntValue = options.longs === String ? "0" : 0;
                    object.doubleValue = 0;
                    if (options.bytes === String)
                        object.stringValue = "";
                    else {
                        object.stringValue = [];
                        if (options.bytes !== Array)
                            object.stringValue = $util.newBuffer(object.stringValue);
                    }
                    object.aggregateValue = "";
                }
                if (message.name && message.name.length) {
                    object.name = [];
                    for (var j = 0; j < message.name.length; ++j)
                        object.name[j] = $root.google.protobuf.UninterpretedOption.NamePart.toObject(message.name[j], options);
                }
                if (message.identifierValue != null && message.hasOwnProperty("identifierValue"))
                    object.identifierValue = message.identifierValue;
                if (message.positiveIntValue != null && message.hasOwnProperty("positiveIntValue"))
                    if (typeof message.positiveIntValue === "number")
                        object.positiveIntValue = options.longs === String ? String(message.positiveIntValue) : message.positiveIntValue;
                    else
                        object.positiveIntValue = options.longs === String ? $util.Long.prototype.toString.call(message.positiveIntValue) : options.longs === Number ? new $util.LongBits(message.positiveIntValue.low >>> 0, message.positiveIntValue.high >>> 0).toNumber(true) : message.positiveIntValue;
                if (message.negativeIntValue != null && message.hasOwnProperty("negativeIntValue"))
                    if (typeof message.negativeIntValue === "number")
                        object.negativeIntValue = options.longs === String ? String(message.negativeIntValue) : message.negativeIntValue;
                    else
                        object.negativeIntValue = options.longs === String ? $util.Long.prototype.toString.call(message.negativeIntValue) : options.longs === Number ? new $util.LongBits(message.negativeIntValue.low >>> 0, message.negativeIntValue.high >>> 0).toNumber() : message.negativeIntValue;
                if (message.doubleValue != null && message.hasOwnProperty("doubleValue"))
                    object.doubleValue = options.json && !isFinite(message.doubleValue) ? String(message.doubleValue) : message.doubleValue;
                if (message.stringValue != null && message.hasOwnProperty("stringValue"))
                    object.stringValue = options.bytes === String ? $util.base64.encode(message.stringValue, 0, message.stringValue.length) : options.bytes === Array ? Array.prototype.slice.call(message.stringValue) : message.stringValue;
                if (message.aggregateValue != null && message.hasOwnProperty("aggregateValue"))
                    object.aggregateValue = message.aggregateValue;
                return object;
            };

            /**
             * Converts this UninterpretedOption to JSON.
             * @function toJSON
             * @memberof google.protobuf.UninterpretedOption
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            UninterpretedOption.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            UninterpretedOption.NamePart = (function() {

                /**
                 * Properties of a NamePart.
                 * @memberof google.protobuf.UninterpretedOption
                 * @interface INamePart
                 * @property {string} namePart NamePart namePart
                 * @property {boolean} isExtension NamePart isExtension
                 */

                /**
                 * Constructs a new NamePart.
                 * @memberof google.protobuf.UninterpretedOption
                 * @classdesc Represents a NamePart.
                 * @implements INamePart
                 * @constructor
                 * @param {google.protobuf.UninterpretedOption.INamePart=} [properties] Properties to set
                 */
                function NamePart(properties) {
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * NamePart namePart.
                 * @member {string} namePart
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @instance
                 */
                NamePart.prototype.namePart = "";

                /**
                 * NamePart isExtension.
                 * @member {boolean} isExtension
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @instance
                 */
                NamePart.prototype.isExtension = false;

                /**
                 * Creates a new NamePart instance using the specified properties.
                 * @function create
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {google.protobuf.UninterpretedOption.INamePart=} [properties] Properties to set
                 * @returns {google.protobuf.UninterpretedOption.NamePart} NamePart instance
                 */
                NamePart.create = function create(properties) {
                    return new NamePart(properties);
                };

                /**
                 * Encodes the specified NamePart message. Does not implicitly {@link google.protobuf.UninterpretedOption.NamePart.verify|verify} messages.
                 * @function encode
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {google.protobuf.UninterpretedOption.INamePart} message NamePart message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                NamePart.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    writer.uint32(/* id 1, wireType 2 =*/10).string(message.namePart);
                    writer.uint32(/* id 2, wireType 0 =*/16).bool(message.isExtension);
                    return writer;
                };

                /**
                 * Encodes the specified NamePart message, length delimited. Does not implicitly {@link google.protobuf.UninterpretedOption.NamePart.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {google.protobuf.UninterpretedOption.INamePart} message NamePart message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                NamePart.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a NamePart message from the specified reader or buffer.
                 * @function decode
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {google.protobuf.UninterpretedOption.NamePart} NamePart
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                NamePart.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.UninterpretedOption.NamePart();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            message.namePart = reader.string();
                            break;
                        case 2:
                            message.isExtension = reader.bool();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    if (!message.hasOwnProperty("namePart"))
                        throw $util.ProtocolError("missing required 'namePart'", { instance: message });
                    if (!message.hasOwnProperty("isExtension"))
                        throw $util.ProtocolError("missing required 'isExtension'", { instance: message });
                    return message;
                };

                /**
                 * Decodes a NamePart message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {google.protobuf.UninterpretedOption.NamePart} NamePart
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                NamePart.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a NamePart message.
                 * @function verify
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                NamePart.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (!$util.isString(message.namePart))
                        return "namePart: string expected";
                    if (typeof message.isExtension !== "boolean")
                        return "isExtension: boolean expected";
                    return null;
                };

                /**
                 * Creates a NamePart message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {google.protobuf.UninterpretedOption.NamePart} NamePart
                 */
                NamePart.fromObject = function fromObject(object) {
                    if (object instanceof $root.google.protobuf.UninterpretedOption.NamePart)
                        return object;
                    var message = new $root.google.protobuf.UninterpretedOption.NamePart();
                    if (object.namePart != null)
                        message.namePart = String(object.namePart);
                    if (object.isExtension != null)
                        message.isExtension = Boolean(object.isExtension);
                    return message;
                };

                /**
                 * Creates a plain object from a NamePart message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @static
                 * @param {google.protobuf.UninterpretedOption.NamePart} message NamePart
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                NamePart.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.defaults) {
                        object.namePart = "";
                        object.isExtension = false;
                    }
                    if (message.namePart != null && message.hasOwnProperty("namePart"))
                        object.namePart = message.namePart;
                    if (message.isExtension != null && message.hasOwnProperty("isExtension"))
                        object.isExtension = message.isExtension;
                    return object;
                };

                /**
                 * Converts this NamePart to JSON.
                 * @function toJSON
                 * @memberof google.protobuf.UninterpretedOption.NamePart
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                NamePart.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return NamePart;
            })();

            return UninterpretedOption;
        })();

        protobuf.SourceCodeInfo = (function() {

            /**
             * Properties of a SourceCodeInfo.
             * @memberof google.protobuf
             * @interface ISourceCodeInfo
             * @property {Array.<google.protobuf.SourceCodeInfo.ILocation>|null} [location] SourceCodeInfo location
             */

            /**
             * Constructs a new SourceCodeInfo.
             * @memberof google.protobuf
             * @classdesc Represents a SourceCodeInfo.
             * @implements ISourceCodeInfo
             * @constructor
             * @param {google.protobuf.ISourceCodeInfo=} [properties] Properties to set
             */
            function SourceCodeInfo(properties) {
                this.location = [];
                if (properties)
                    for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                        if (properties[keys[i]] != null)
                            this[keys[i]] = properties[keys[i]];
            }

            /**
             * SourceCodeInfo location.
             * @member {Array.<google.protobuf.SourceCodeInfo.ILocation>} location
             * @memberof google.protobuf.SourceCodeInfo
             * @instance
             */
            SourceCodeInfo.prototype.location = $util.emptyArray;

            /**
             * Creates a new SourceCodeInfo instance using the specified properties.
             * @function create
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {google.protobuf.ISourceCodeInfo=} [properties] Properties to set
             * @returns {google.protobuf.SourceCodeInfo} SourceCodeInfo instance
             */
            SourceCodeInfo.create = function create(properties) {
                return new SourceCodeInfo(properties);
            };

            /**
             * Encodes the specified SourceCodeInfo message. Does not implicitly {@link google.protobuf.SourceCodeInfo.verify|verify} messages.
             * @function encode
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {google.protobuf.ISourceCodeInfo} message SourceCodeInfo message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            SourceCodeInfo.encode = function encode(message, writer) {
                if (!writer)
                    writer = $Writer.create();
                if (message.location != null && message.location.length)
                    for (var i = 0; i < message.location.length; ++i)
                        $root.google.protobuf.SourceCodeInfo.Location.encode(message.location[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
                return writer;
            };

            /**
             * Encodes the specified SourceCodeInfo message, length delimited. Does not implicitly {@link google.protobuf.SourceCodeInfo.verify|verify} messages.
             * @function encodeDelimited
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {google.protobuf.ISourceCodeInfo} message SourceCodeInfo message or plain object to encode
             * @param {$protobuf.Writer} [writer] Writer to encode to
             * @returns {$protobuf.Writer} Writer
             */
            SourceCodeInfo.encodeDelimited = function encodeDelimited(message, writer) {
                return this.encode(message, writer).ldelim();
            };

            /**
             * Decodes a SourceCodeInfo message from the specified reader or buffer.
             * @function decode
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @param {number} [length] Message length if known beforehand
             * @returns {google.protobuf.SourceCodeInfo} SourceCodeInfo
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            SourceCodeInfo.decode = function decode(reader, length) {
                if (!(reader instanceof $Reader))
                    reader = $Reader.create(reader);
                var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.SourceCodeInfo();
                while (reader.pos < end) {
                    var tag = reader.uint32();
                    switch (tag >>> 3) {
                    case 1:
                        if (!(message.location && message.location.length))
                            message.location = [];
                        message.location.push($root.google.protobuf.SourceCodeInfo.Location.decode(reader, reader.uint32()));
                        break;
                    default:
                        reader.skipType(tag & 7);
                        break;
                    }
                }
                return message;
            };

            /**
             * Decodes a SourceCodeInfo message from the specified reader or buffer, length delimited.
             * @function decodeDelimited
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
             * @returns {google.protobuf.SourceCodeInfo} SourceCodeInfo
             * @throws {Error} If the payload is not a reader or valid buffer
             * @throws {$protobuf.util.ProtocolError} If required fields are missing
             */
            SourceCodeInfo.decodeDelimited = function decodeDelimited(reader) {
                if (!(reader instanceof $Reader))
                    reader = new $Reader(reader);
                return this.decode(reader, reader.uint32());
            };

            /**
             * Verifies a SourceCodeInfo message.
             * @function verify
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {Object.<string,*>} message Plain object to verify
             * @returns {string|null} `null` if valid, otherwise the reason why it is not
             */
            SourceCodeInfo.verify = function verify(message) {
                if (typeof message !== "object" || message === null)
                    return "object expected";
                if (message.location != null && message.hasOwnProperty("location")) {
                    if (!Array.isArray(message.location))
                        return "location: array expected";
                    for (var i = 0; i < message.location.length; ++i) {
                        var error = $root.google.protobuf.SourceCodeInfo.Location.verify(message.location[i]);
                        if (error)
                            return "location." + error;
                    }
                }
                return null;
            };

            /**
             * Creates a SourceCodeInfo message from a plain object. Also converts values to their respective internal types.
             * @function fromObject
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {Object.<string,*>} object Plain object
             * @returns {google.protobuf.SourceCodeInfo} SourceCodeInfo
             */
            SourceCodeInfo.fromObject = function fromObject(object) {
                if (object instanceof $root.google.protobuf.SourceCodeInfo)
                    return object;
                var message = new $root.google.protobuf.SourceCodeInfo();
                if (object.location) {
                    if (!Array.isArray(object.location))
                        throw TypeError(".google.protobuf.SourceCodeInfo.location: array expected");
                    message.location = [];
                    for (var i = 0; i < object.location.length; ++i) {
                        if (typeof object.location[i] !== "object")
                            throw TypeError(".google.protobuf.SourceCodeInfo.location: object expected");
                        message.location[i] = $root.google.protobuf.SourceCodeInfo.Location.fromObject(object.location[i]);
                    }
                }
                return message;
            };

            /**
             * Creates a plain object from a SourceCodeInfo message. Also converts values to other types if specified.
             * @function toObject
             * @memberof google.protobuf.SourceCodeInfo
             * @static
             * @param {google.protobuf.SourceCodeInfo} message SourceCodeInfo
             * @param {$protobuf.IConversionOptions} [options] Conversion options
             * @returns {Object.<string,*>} Plain object
             */
            SourceCodeInfo.toObject = function toObject(message, options) {
                if (!options)
                    options = {};
                var object = {};
                if (options.arrays || options.defaults)
                    object.location = [];
                if (message.location && message.location.length) {
                    object.location = [];
                    for (var j = 0; j < message.location.length; ++j)
                        object.location[j] = $root.google.protobuf.SourceCodeInfo.Location.toObject(message.location[j], options);
                }
                return object;
            };

            /**
             * Converts this SourceCodeInfo to JSON.
             * @function toJSON
             * @memberof google.protobuf.SourceCodeInfo
             * @instance
             * @returns {Object.<string,*>} JSON object
             */
            SourceCodeInfo.prototype.toJSON = function toJSON() {
                return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
            };

            SourceCodeInfo.Location = (function() {

                /**
                 * Properties of a Location.
                 * @memberof google.protobuf.SourceCodeInfo
                 * @interface ILocation
                 * @property {Array.<number>|null} [path] Location path
                 * @property {Array.<number>|null} [span] Location span
                 * @property {string|null} [leadingComments] Location leadingComments
                 * @property {string|null} [trailingComments] Location trailingComments
                 */

                /**
                 * Constructs a new Location.
                 * @memberof google.protobuf.SourceCodeInfo
                 * @classdesc Represents a Location.
                 * @implements ILocation
                 * @constructor
                 * @param {google.protobuf.SourceCodeInfo.ILocation=} [properties] Properties to set
                 */
                function Location(properties) {
                    this.path = [];
                    this.span = [];
                    if (properties)
                        for (var keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                            if (properties[keys[i]] != null)
                                this[keys[i]] = properties[keys[i]];
                }

                /**
                 * Location path.
                 * @member {Array.<number>} path
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @instance
                 */
                Location.prototype.path = $util.emptyArray;

                /**
                 * Location span.
                 * @member {Array.<number>} span
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @instance
                 */
                Location.prototype.span = $util.emptyArray;

                /**
                 * Location leadingComments.
                 * @member {string} leadingComments
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @instance
                 */
                Location.prototype.leadingComments = "";

                /**
                 * Location trailingComments.
                 * @member {string} trailingComments
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @instance
                 */
                Location.prototype.trailingComments = "";

                /**
                 * Creates a new Location instance using the specified properties.
                 * @function create
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {google.protobuf.SourceCodeInfo.ILocation=} [properties] Properties to set
                 * @returns {google.protobuf.SourceCodeInfo.Location} Location instance
                 */
                Location.create = function create(properties) {
                    return new Location(properties);
                };

                /**
                 * Encodes the specified Location message. Does not implicitly {@link google.protobuf.SourceCodeInfo.Location.verify|verify} messages.
                 * @function encode
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {google.protobuf.SourceCodeInfo.ILocation} message Location message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Location.encode = function encode(message, writer) {
                    if (!writer)
                        writer = $Writer.create();
                    if (message.path != null && message.path.length) {
                        writer.uint32(/* id 1, wireType 2 =*/10).fork();
                        for (var i = 0; i < message.path.length; ++i)
                            writer.int32(message.path[i]);
                        writer.ldelim();
                    }
                    if (message.span != null && message.span.length) {
                        writer.uint32(/* id 2, wireType 2 =*/18).fork();
                        for (var i = 0; i < message.span.length; ++i)
                            writer.int32(message.span[i]);
                        writer.ldelim();
                    }
                    if (message.leadingComments != null && message.hasOwnProperty("leadingComments"))
                        writer.uint32(/* id 3, wireType 2 =*/26).string(message.leadingComments);
                    if (message.trailingComments != null && message.hasOwnProperty("trailingComments"))
                        writer.uint32(/* id 4, wireType 2 =*/34).string(message.trailingComments);
                    return writer;
                };

                /**
                 * Encodes the specified Location message, length delimited. Does not implicitly {@link google.protobuf.SourceCodeInfo.Location.verify|verify} messages.
                 * @function encodeDelimited
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {google.protobuf.SourceCodeInfo.ILocation} message Location message or plain object to encode
                 * @param {$protobuf.Writer} [writer] Writer to encode to
                 * @returns {$protobuf.Writer} Writer
                 */
                Location.encodeDelimited = function encodeDelimited(message, writer) {
                    return this.encode(message, writer).ldelim();
                };

                /**
                 * Decodes a Location message from the specified reader or buffer.
                 * @function decode
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @param {number} [length] Message length if known beforehand
                 * @returns {google.protobuf.SourceCodeInfo.Location} Location
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Location.decode = function decode(reader, length) {
                    if (!(reader instanceof $Reader))
                        reader = $Reader.create(reader);
                    var end = length === undefined ? reader.len : reader.pos + length, message = new $root.google.protobuf.SourceCodeInfo.Location();
                    while (reader.pos < end) {
                        var tag = reader.uint32();
                        switch (tag >>> 3) {
                        case 1:
                            if (!(message.path && message.path.length))
                                message.path = [];
                            if ((tag & 7) === 2) {
                                var end2 = reader.uint32() + reader.pos;
                                while (reader.pos < end2)
                                    message.path.push(reader.int32());
                            } else
                                message.path.push(reader.int32());
                            break;
                        case 2:
                            if (!(message.span && message.span.length))
                                message.span = [];
                            if ((tag & 7) === 2) {
                                var end2 = reader.uint32() + reader.pos;
                                while (reader.pos < end2)
                                    message.span.push(reader.int32());
                            } else
                                message.span.push(reader.int32());
                            break;
                        case 3:
                            message.leadingComments = reader.string();
                            break;
                        case 4:
                            message.trailingComments = reader.string();
                            break;
                        default:
                            reader.skipType(tag & 7);
                            break;
                        }
                    }
                    return message;
                };

                /**
                 * Decodes a Location message from the specified reader or buffer, length delimited.
                 * @function decodeDelimited
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
                 * @returns {google.protobuf.SourceCodeInfo.Location} Location
                 * @throws {Error} If the payload is not a reader or valid buffer
                 * @throws {$protobuf.util.ProtocolError} If required fields are missing
                 */
                Location.decodeDelimited = function decodeDelimited(reader) {
                    if (!(reader instanceof $Reader))
                        reader = new $Reader(reader);
                    return this.decode(reader, reader.uint32());
                };

                /**
                 * Verifies a Location message.
                 * @function verify
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {Object.<string,*>} message Plain object to verify
                 * @returns {string|null} `null` if valid, otherwise the reason why it is not
                 */
                Location.verify = function verify(message) {
                    if (typeof message !== "object" || message === null)
                        return "object expected";
                    if (message.path != null && message.hasOwnProperty("path")) {
                        if (!Array.isArray(message.path))
                            return "path: array expected";
                        for (var i = 0; i < message.path.length; ++i)
                            if (!$util.isInteger(message.path[i]))
                                return "path: integer[] expected";
                    }
                    if (message.span != null && message.hasOwnProperty("span")) {
                        if (!Array.isArray(message.span))
                            return "span: array expected";
                        for (var i = 0; i < message.span.length; ++i)
                            if (!$util.isInteger(message.span[i]))
                                return "span: integer[] expected";
                    }
                    if (message.leadingComments != null && message.hasOwnProperty("leadingComments"))
                        if (!$util.isString(message.leadingComments))
                            return "leadingComments: string expected";
                    if (message.trailingComments != null && message.hasOwnProperty("trailingComments"))
                        if (!$util.isString(message.trailingComments))
                            return "trailingComments: string expected";
                    return null;
                };

                /**
                 * Creates a Location message from a plain object. Also converts values to their respective internal types.
                 * @function fromObject
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {Object.<string,*>} object Plain object
                 * @returns {google.protobuf.SourceCodeInfo.Location} Location
                 */
                Location.fromObject = function fromObject(object) {
                    if (object instanceof $root.google.protobuf.SourceCodeInfo.Location)
                        return object;
                    var message = new $root.google.protobuf.SourceCodeInfo.Location();
                    if (object.path) {
                        if (!Array.isArray(object.path))
                            throw TypeError(".google.protobuf.SourceCodeInfo.Location.path: array expected");
                        message.path = [];
                        for (var i = 0; i < object.path.length; ++i)
                            message.path[i] = object.path[i] | 0;
                    }
                    if (object.span) {
                        if (!Array.isArray(object.span))
                            throw TypeError(".google.protobuf.SourceCodeInfo.Location.span: array expected");
                        message.span = [];
                        for (var i = 0; i < object.span.length; ++i)
                            message.span[i] = object.span[i] | 0;
                    }
                    if (object.leadingComments != null)
                        message.leadingComments = String(object.leadingComments);
                    if (object.trailingComments != null)
                        message.trailingComments = String(object.trailingComments);
                    return message;
                };

                /**
                 * Creates a plain object from a Location message. Also converts values to other types if specified.
                 * @function toObject
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @static
                 * @param {google.protobuf.SourceCodeInfo.Location} message Location
                 * @param {$protobuf.IConversionOptions} [options] Conversion options
                 * @returns {Object.<string,*>} Plain object
                 */
                Location.toObject = function toObject(message, options) {
                    if (!options)
                        options = {};
                    var object = {};
                    if (options.arrays || options.defaults) {
                        object.path = [];
                        object.span = [];
                    }
                    if (options.defaults) {
                        object.leadingComments = "";
                        object.trailingComments = "";
                    }
                    if (message.path && message.path.length) {
                        object.path = [];
                        for (var j = 0; j < message.path.length; ++j)
                            object.path[j] = message.path[j];
                    }
                    if (message.span && message.span.length) {
                        object.span = [];
                        for (var j = 0; j < message.span.length; ++j)
                            object.span[j] = message.span[j];
                    }
                    if (message.leadingComments != null && message.hasOwnProperty("leadingComments"))
                        object.leadingComments = message.leadingComments;
                    if (message.trailingComments != null && message.hasOwnProperty("trailingComments"))
                        object.trailingComments = message.trailingComments;
                    return object;
                };

                /**
                 * Converts this Location to JSON.
                 * @function toJSON
                 * @memberof google.protobuf.SourceCodeInfo.Location
                 * @instance
                 * @returns {Object.<string,*>} JSON object
                 */
                Location.prototype.toJSON = function toJSON() {
                    return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
                };

                return Location;
            })();

            return SourceCodeInfo;
        })();

        return protobuf;
    })();

    return google;
})();

module.exports = $root;
