"use strict";
var Block = (function () {
    function Block(header, body) {
        this.header = header;
        this.body = body;
    }
    return Block;
}());
exports.Block = Block;
var BlockHeader = (function () {
    function BlockHeader(seq, block_hash, previous_block_hash, timestamp, fee, version, tx_body_hash) {
        this.seq = seq;
        this.block_hash = block_hash;
        this.previous_block_hash = previous_block_hash;
        this.timestamp = timestamp;
        this.fee = fee;
        this.version = version;
        this.tx_body_hash = tx_body_hash;
    }
    return BlockHeader;
}());
exports.BlockHeader = BlockHeader;
var BlockBody = (function () {
    function BlockBody(txns) {
        this.txns = txns;
    }
    return BlockBody;
}());
exports.BlockBody = BlockBody;
var Output = (function () {
    function Output(uxid, dst, coins, hrs) {
        this.uxid = uxid;
        this.dst = dst;
        this.coins = coins;
        this.hrs = hrs;
    }
    return Output;
}());
exports.Output = Output;
var BlockChainMetaDataHead = (function () {
    function BlockChainMetaDataHead(seq, block_hash, previous_block_hash, timestamp, fee, version, tx_body_hash) {
        this.seq = seq;
        this.block_hash = block_hash;
        this.previous_block_hash = previous_block_hash;
        this.timestamp = timestamp;
        this.fee = fee;
        this.version = version;
        this.tx_body_hash = tx_body_hash;
    }
    return BlockChainMetaDataHead;
}());
exports.BlockChainMetaDataHead = BlockChainMetaDataHead;
var BlockChainMetaData = (function () {
    function BlockChainMetaData(head, unspents, unconfirmed) {
        this.head = head;
        this.unspents = unspents;
        this.unconfirmed = unconfirmed;
    }
    return BlockChainMetaData;
}());
exports.BlockChainMetaData = BlockChainMetaData;
var Transaction = (function () {
    function Transaction(length, type, txid, inner_hash, sigs, inputs, outputs) {
        this.length = length;
        this.type = type;
        this.txid = txid;
        this.inner_hash = inner_hash;
        this.sigs = sigs;
        this.inputs = inputs;
        this.outputs = outputs;
    }
    return Transaction;
}());
exports.Transaction = Transaction;
var CoinSupply = (function () {
    function CoinSupply(coinSupply, coinCap) {
        this.coinSupply = coinSupply;
        this.coinCap = coinCap;
    }
    return CoinSupply;
}());
exports.CoinSupply = CoinSupply;
var BlockResponse = (function () {
    function BlockResponse(blocks) {
        this.blocks = blocks;
    }
    return BlockResponse;
}());
exports.BlockResponse = BlockResponse;
//# sourceMappingURL=block.js.map