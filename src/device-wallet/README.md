# Go lang api for skycoin hardware wallet

## Install protoc

    sudo apt-get install protobuf-compiler golang-goprotobuf-dev
    protoc -I../tiny-firmware/vendor/nanopb/generator/proto/ -I ./protob  --go_out=./protob protob/messages.proto protob/types.proto

## Generate protobuf files

Only once each time the messages change:

    protoc -I ./protob  --go_out=./protob protob/messages.proto protob/types.proto
    # will generate messages.pb.go and types.pb.go files

## Imported files

The folders "wire, usbhid, and usb" have been copy pasted from trezor: trezor project https://github.com/trezor/trezord-go/blob/master/ 

master HEAD at the time of the copy paste: 4527402f7004dfe677225315a0dd4d4b1b74be49