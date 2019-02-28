# rm -R cmd/cxcoin

# go install ../fiber-init/cmd/fiber-init/...
# go install ./cmd/newcoin/...

if [ "$1" == "update" ]; then
    # go install -gcflags=all=-e ../fiber-init/cmd/fiber-init/...
    go install ../fiber-init/cmd/fiber-init/...
    newcoin createcoin --coin cxcoin \
    	    --template-dir "$GOPATH/src/github.com/skycoin/skycoin/template" \
    	    --config-file "fiber.toml" \
    	    --config-dir "$GOPATH/src/github.com/skycoin/skycoin/"

    go install ./cmd/cxcoin/...
fi

if [ "$1" == "distr" ]; then
    fiber-init distributecoins --coin cxcoin \
	       --template-dir "${GOPATH}/src/github.com/skycoin/skycoin/template" \
	       --config-file "fiber.toml" \
	       --config-dir "$GOPATH/src/github.com/skycoin/skycoin/" \
	       --project-root "${GOPATH}/src/github.com/skycoin/skycoin" \
	       --seckey $FIBERCOIN_GENESIS_SECKEY
fi

if [ "$1" == "skycoin" ]; then
    go build -gcflags=all=-e ./cmd/skycoin/...
fi

if [ "$1" == "master" ]; then
    cxcoin -enable-all-api-sets -block-publisher=true -blockchain-secret-key=$FIBERCOIN_GENESIS_SECKEY
fi

if [ "$1" == "node" ]; then
    PORT="$2"; ./run-cxcoin.sh \
		   -localhost-only \
		   -disable-default-peers \
		   -custom-peers-file=localhost-peers.txt \
		   -download-peerlist=false \
		   -launch-browser=false \
		   -data-dir=/tmp/$2 \
		   -web-interface-port=$(expr $2 + 420) \
		   -port=$2 \
		   -genesis-address "23v7mT1uLpViNKZHh9aww4VChxizqKsNq4E" \
		   -genesis-signature "9043250618603506aad2d4e54b748d3adff525f0b27010b86d70d3d8ea28c65518c25c379d5419f18cb9c6033a2c097e5355b5d356f56f89824df1a4f5d9cf2701" \
		   -blockchain-public-key "02583e5ebbf85522474e0f17e681e62ca37910db6b8792763af4e97663c31a7984" \
		   -max-txn-size-unconfirmed=65535 \
		   -max-txn-size-create-block=65535 \
		   -max-block-size=65535
fi

if [ "$1" == "cxProgram" ]; then
    cx --blockchain --heap-initial 100 --stack-size 100 ../cx/testing.cx
fi

if [ "$1" == "createWallet" ]; then
    # rm ~/.cxcoin/wallets/*
    ADDRESS="TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX"
    SEED="museum nothing practice weird wheel dignity economy attend mask recipe minor dress"
    LABEL="cxcoin"
    CSRF_TOKEN=$(curl -s http://127.0.0.1:6420/api/v1/csrf | jq -r '.csrf_token')
    WALLET=$(curl -s -X POST http://127.0.0.1:6420/api/v1/wallet/create \
	 -H "X-CSRF-Token: $CSRF_TOKEN" \
	 -H "Content-Type: application/x-www-form-urlencoded" \
	 -d "seed=$SEED" \
	 -d "coin=$LABEL" \
	 -d "label=$LABEL" | jq -r '.meta.filename')
    
    echo $ADDRESS
    echo $WALLET
    
    export ADDRESS
    export WALLET
fi

if [ "$1" == "wallet" ]; then
    curl http://127.0.0.1:6420/api/v1/wallet?id=$WALLET
fi

if [ "$1" == "balance" ]; then
    curl http://127.0.0.1:6420/api/v1/balance\?addrs\=$ADDRESS
fi

if [ "$1" == "txn" ]; then
    CSRF_TOKEN=$(curl -s http://127.0.0.1:6420/api/v1/csrf | jq -r '.csrf_token')
    TXN=$(curl -s -X POST http://127.0.0.1:6420/api/v1/wallet/transaction \
	 -H "X-CSRF-Token: $CSRF_TOKEN" \
	 -H 'content-type: application/json' -d '{
    "hours_selection": {
        "type": "auto",
        "mode": "share",
        "share_factor": "0.5"
    },
    "wallet": {
        "id": "2019_02_26_bd54.wlt"
    },
    "change_address": "TkyD4wD64UE6M5BkNQA17zaf7Xcg4AufwX",
    "to": [{
        "address": "2PBcLADETphmqWV7sujRZdh3UcabssgKAEB",
        "coins": "1"
    }, {
        "address": "2PBcLADETphmqWV7sujRZdh3UcabssgKAEB",
        "coins": "8.99"
    }]
}' | jq -r '.encoded_transaction')

    echo $TXN

    # curl -X POST http://127.0.0.1:6420/api/v1/injectTransaction \
    # 	 -H "X-CSRF-Token: $CSRF_TOKEN" \
    # 	 -H 'content-type: application/json' -d "{\"rawtx\": \"11010000009929ac42574d6ca3d9ac9396baac87d13b41105690fd6334ec996e4c0b5a763f010000007c3cce41ae454abdd7dc9da916cf2b28df4d8d21663e72fbcae5637bdbe04ef32ba043e51af98ae0e1baca7331bb937a8603d8029cc1153b73c3b8b5fe0721680001000000c3aac6dcd07739396d267f17a3316d43c3e622918aa937f95e8cdd21a1c78af40300000000c745a77239f02e0c6d06ded997563dc956e37a0b40420f00000000001a000000000000000000000000c745a77239f02e0c6d06ded997563dc956e37a0b302d890000000000e0000000000000000000000000427fec754e22482758ca61f781cd7f8c55e9192890a00cd4e8000000fa000000000000000000000000000000\"}"

    # curl -X POST http://127.0.0.1:6420/api/v1/injectTransaction \
    # 	 -H "X-CSRF-Token: $CSRF_TOKEN" \
    # 	 -H 'content-type: application/json' -d "{\"rawtx\": \"$TXN\"}"

    curl -X POST http://127.0.0.1:6420/api/v2/wallet/transaction/sign \
	 -H "X-CSRF-Token: $CSRF_TOKEN" \
	 -H 'content-type: application/json' -d '{
    "wallet_id": "2019_02_26_bd54.wlt",
    "password": "password",
    "encoded_transaction": "11010000009929ac42574d6ca3d9ac9396baac87d13b41105690fd6334ec996e4c0b5a763f010000007c3cce41ae454abdd7dc9da916cf2b28df4d8d21663e72fbcae5637bdbe04ef32ba043e51af98ae0e1baca7331bb937a8603d8029cc1153b73c3b8b5fe0721680001000000c3aac6dcd07739396d267f17a3316d43c3e622918aa937f95e8cdd21a1c78af40300000000c745a77239f02e0c6d06ded997563dc956e37a0b40420f00000000001a000000000000000000000000c745a77239f02e0c6d06ded997563dc956e37a0b302d890000000000e0000000000000000000000000427fec754e22482758ca61f781cd7f8c55e9192890a00cd4e8000000fa000000000000000000000000000000"
}'
fi
