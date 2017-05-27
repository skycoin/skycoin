"use strict";
var UnspentOutput = (function () {
    function UnspentOutput(uxid, time, src_block_seq, src_tx, owner_address, coins, hours, spent_block_seq, spent_tx) {
        this.uxid = uxid;
        this.time = time;
        this.src_block_seq = src_block_seq;
        this.src_tx = src_tx;
        this.owner_address = owner_address;
        this.coins = coins;
        this.hours = hours;
        this.spent_block_seq = spent_block_seq;
        this.spent_tx = spent_tx;
    }
    return UnspentOutput;
}());
exports.UnspentOutput = UnspentOutput;
var AddressBalanceResponse = (function () {
    function AddressBalanceResponse(head_outputs, outgoing_outputs, incoming_outputs) {
        this.head_outputs = head_outputs;
        this.outgoing_outputs = outgoing_outputs;
        this.incoming_outputs = incoming_outputs;
    }
    return AddressBalanceResponse;
}());
exports.AddressBalanceResponse = AddressBalanceResponse;
var HeadOutput = (function () {
    function HeadOutput(hash, src_tx, address, coins, hours) {
        this.hash = hash;
        this.src_tx = src_tx;
        this.address = address;
        this.coins = coins;
        this.hours = hours;
    }
    return HeadOutput;
}());
exports.HeadOutput = HeadOutput;
//# sourceMappingURL=UnspentOutput.js.map