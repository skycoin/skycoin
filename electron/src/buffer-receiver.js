const messages = require('./protob/skycoin');

class BufferReceiver {
    constructor() {
        this.msgIndex = 0;
        this.msgSize = undefined;
        this.bytesToGet = undefined;
        this.kind = undefined;
        this.dataBuffer = undefined;
    }

    parseHeader(data) {
        const dv8 = new Uint8Array(data);
        this.kind = new Uint16Array(dv8.slice(4, 5))[0];
        this.msgSize = new Uint32Array(dv8.slice(8, 11))[0];
        this.dataBuffer = new Uint8Array(64 * Math.ceil(this.msgSize / 64));
        this.dataBuffer.set(dv8.slice(9));
        this.bytesToGet = this.msgSize + 9 - 64;

        console.log(
            "Received header",
            " msg this.kind: ", messages.MessageType[this.kind],
            " size: ", this.msgSize,
            "buffer lenght: ", this.dataBuffer.byteLength,
            "\nRemaining bytesToGet:", this.bytesToGet
            );
    }

    // eslint-disable-next-line max-statements
    receiveBuffer(data, callback) {

        if (this.bytesToGet === undefined) {
            this.parseHeader(data);

            if (this.bytesToGet > 0) {
                return;
            }
            callback(this.kind, this.dataBuffer.slice(0, this.msgSize));
            return;
        }

        this.dataBuffer.set(data.slice(1), (63 * this.msgIndex) + 55);
        this.msgIndex += 1;
        this.bytesToGet -= 64;

        console.log(
            "Received data", " msg kind: ",
            messages.MessageType[this.kind],
            " size: ", this.msgSize, "buffer lenght: ",
            this.dataBuffer.byteLength
            );

        if (this.bytesToGet > 0) {
            return;
        }
        if (callback) {
            // eslint-disable-next-line callback-return
            callback(this.kind, this.dataBuffer.slice(0, this.msgSize));
        }
    }
}

module.exports = {
    BufferReceiver
};
