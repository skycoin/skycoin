#!/bin/bash
# End to end test of a joint transaction on a throwaway fiber coin chain.
# Two wallets each fund a small output, then one transaction spends both: the
# customer gets the coins back and the vendor gets coin hours. Each wallet signs
# only its own input, and the result must verify, broadcast and confirm.
# Needs go, curl and jq.

set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."
source ci-scripts/joint-chain.sh
joint_chain_setup

echo "=== create the joint transaction"
PAY=$(( CH / 2 ))
RAW=$("$SKY" cli createJointTransaction --in "$CUX" --in "$VUX" \
    --out "$C1,0.002,0" --out "$V1,0.001,$PAY")
"$SKY" cli decodeRawTransaction --wallet "$CUSTOMER" "$RAW"
[[ $(jq -r '.inputs | length' <<<"$("$SKY" cli decodeRawTransaction "$RAW")") == 2 ]] || fail "want 2 inputs"

if "$SKY" cli createJointTransaction --in "$CUX" --in "$VUX" --out "$C1,0.001,0" >/dev/null 2>&1; then
    fail "coins left over without a change address were accepted"
fi
pass "leftover coins need --change-address"

echo "=== customer signs"
OUT=$("$SKY" cli signTransaction "$CUSTOMER" "$RAW" 2>"$DATA_DIR/sign.err")
grep -q "signed 1 of 2 inputs" "$DATA_DIR/sign.err" || fail "customer: $(cat "$DATA_DIR/sign.err")"
RAW=$(jq -r '.encoded_transaction' <<<"$OUT")
if "$SKY" cli verifyTransaction "$RAW" >/dev/null 2>&1; then
    fail "a half signed transaction verified as spendable"
fi
# The web wallet reviews a half signed transaction this way
UNSIGNED_OK=$(curl -s -X POST "$HOST/api/v2/transaction/verify" -H "Content-Type: application/json" \
    -d "$(jq -n --arg tx "$RAW" '{encoded_transaction: $tx, unsigned: true}')" | jq -r '.error == null and .data.unsigned')
[[ "$UNSIGNED_OK" == true ]] || fail "verify as unsigned rejected a half signed transaction"
if "$SKY" cli signTransaction "$VENDOR" "$RAW" --sign-indexes 0 >/dev/null 2>&1; then
    fail "vendor signed the customer's input"
fi
pass "customer signed 1 of 2, not yet spendable"

echo "=== vendor signs"
OUT=$("$SKY" cli signTransaction "$VENDOR" "$RAW" 2>"$DATA_DIR/sign.err")
grep -q "signed 2 of 2 inputs" "$DATA_DIR/sign.err" || fail "vendor: $(cat "$DATA_DIR/sign.err")"
RAW=$(jq -r '.encoded_transaction' <<<"$OUT")
"$SKY" cli decodeRawTransaction --wallet "$VENDOR" "$RAW" | tee "$DATA_DIR/review.txt"
grep -q "ready to broadcast" "$DATA_DIR/review.txt" || fail "review does not say fully signed"
"$SKY" cli verifyTransaction "$RAW" | grep -q "spendable" || fail "verify"
pass "vendor signed 2 of 2, verified"

echo "=== inject"
SEQ=$(head_seq)
TXID=$(inject "$RAW")
wait_block $((SEQ + 1))
STATUS=$("$SKY" cli transaction "$TXID" | jq -r '.transaction.status.confirmed')
[[ "$STATUS" == true ]] || fail "transaction $TXID not confirmed"
GOT=$("$SKY" cli addressOutputs "$V1" | jq -r --arg tx "$TXID" '[.outputs.head_outputs[] | select(.src_tx == $tx)][0].hours')
[[ "$GOT" == "$PAY" ]] || fail "vendor output has $GOT hours, want $PAY"
BACK=$("$SKY" cli addressOutputs "$C1" | jq -r --arg tx "$TXID" '[.outputs.head_outputs[] | select(.src_tx == $tx)][0].coins')
[[ "$BACK" == "0.002000" ]] || fail "customer output has $BACK coins, want 0.002000"
pass "joint transaction $TXID confirmed, vendor got $PAY hours"

echo "ALL JOINT TRANSACTION TESTS PASSED"
