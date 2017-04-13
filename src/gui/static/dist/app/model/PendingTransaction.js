System.register([], function(exports_1, context_1) {
    "use strict";
    var __moduleName = context_1 && context_1.id;
    var PendingTxn, Transaction, Output;
    return {
        setters:[],
        execute: function() {
            PendingTxn = (function () {
                function PendingTxn(transaction, received, checked, announced) {
                    this.transaction = transaction;
                    this.received = received;
                    this.checked = checked;
                    this.announced = announced;
                }
                return PendingTxn;
            }());
            exports_1("PendingTxn", PendingTxn);
            Transaction = (function () {
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
            exports_1("Transaction", Transaction);
            Output = (function () {
                function Output(uxid, dst, coins, hrs) {
                    this.uxid = uxid;
                    this.dst = dst;
                    this.coins = coins;
                    this.hrs = hrs;
                }
                return Output;
            }());
            exports_1("Output", Output);
        }
    }
});

//# sourceMappingURL=PendingTransaction.js.map
