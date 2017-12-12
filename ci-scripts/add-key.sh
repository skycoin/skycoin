#!/usr/bin/env bash

KEY_CHAIN=login.keychain
echo "security create keychain"
security create-keychain -p travis $KEY_CHAIN
# Make the keychain the default so identities are found
echo "security default-keychain"
security default-keychain -s $KEY_CHAIN
# Unlock the keychain
echo "unlock the keychain"
security unlock-keychain -p travis $KEY_CHAIN
# Set keychain locking timeout to 3600 seconds
echo "set keychain locking timeout to 3600"
security set-keychain-settings -t 3600 -u $KEY_CHAIN

# Add certificates to keychain and allow codesign to access them
# security import ./electron/ci-scripts/certs/dist.cer -k $KEY_CHAIN -T /usr/bin/codesign
# security import ./electron/ci-scripts/certs/dev.cer -k $KEY_CHAIN -T /usr/bin/codesign
echo "import distp12"
security import $GOPATH/src/github.com/skycoin/skycoin/ci-scripts/certs/dist.p12 -k $KEY_CHAIN -P $CERT_PWD  -A /usr/bin/codesign
# security import ./scripts/certs/dev.p12 -k $KEY_CHAIN -P DEVELOPMENT_KEY_PASSWORD  -T /usr/bin/codesign
echo "set key partition list"
security set-key-partition-list -S apple-tool:,apple: -s -k travis $KEY_CHAIN

echo "list keychains: "
security list-keychains
echo " ****** "

echo "find indentities keychains: "
security find-identity -p codesigning  ~/Library/Keychains/$KEY_CHAIN
echo " ****** "
