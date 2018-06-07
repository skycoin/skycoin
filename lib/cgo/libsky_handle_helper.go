package main

/*

  #include <string.h>
  #include <stdlib.h>


  #include "skytypes.h"
*/
import "C"

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"unsafe"

	api "github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

//export SKY_JsonEncode_Handle
func SKY_JsonEncode_Handle(handle C.Handle, json_string *C.GoString_) uint32 {
	obj, ok := lookupHandle(handle)
	if ok {
		jsonBytes, err := json.Marshal(obj)
		if err == nil {
			copyString(string(jsonBytes), json_string)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Progress_GetCurrent
func SKY_Handle_Progress_GetCurrent(handle C.Handle, current *uint64) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*daemon.BlockchainProgress); isOK {
			*current = obj.Current
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Block_GetHeadSeq
func SKY_Handle_Block_GetHeadSeq(handle C.Handle, seq *uint64) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*visor.ReadableBlock); isOK {
			*seq = obj.Head.BkSeq
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Block_GetHeadHash
func SKY_Handle_Block_GetHeadHash(handle C.Handle, hash *C.GoString_) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*visor.ReadableBlock); isOK {
			copyString(obj.Head.BlockHash, hash)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Block_GetPreviousBlockHash
func SKY_Handle_Block_GetPreviousBlockHash(handle C.Handle, hash *C.GoString_) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*visor.ReadableBlock); isOK {
			copyString(obj.Head.PreviousBlockHash, hash)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Blocks_GetAt
func SKY_Handle_Blocks_GetAt(handle C.Handle,
	index uint64, blockHandle *C.Handle) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*visor.ReadableBlocks); isOK {
			*blockHandle = registerHandle(&obj.Blocks[index])
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Blocks_GetCount
func SKY_Handle_Blocks_GetCount(handle C.Handle,
	count *uint64) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*visor.ReadableBlocks); isOK {
			*count = uint64(len(obj.Blocks))
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Connections_GetCount
func SKY_Handle_Connections_GetCount(handle C.Handle,
	count *uint64) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).(*api.Connections); isOK {
			*count = uint64(len(obj.Connections))
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Strings_GetCount
func SKY_Handle_Strings_GetCount(handle C.Strings__Handle,
	count *uint32) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).([]string); isOK {
			*count = uint32(len(obj))
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Strings_Sort
func SKY_Handle_Strings_Sort(handle C.Strings__Handle) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).([]string); isOK {
			sort.Strings(obj)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_Handle_Strings_GetAt
func SKY_Handle_Strings_GetAt(handle C.Strings__Handle,
	index int,
	str *C.GoString_) uint32 {
	obj, ok := lookupHandle(C.Handle(handle))
	if ok {
		if obj, isOK := (obj).([]string); isOK {
			copyString(obj[index], str)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_api_Handle_Client_GetWalletDir
func SKY_api_Handle_Client_GetWalletDir(handle C.Client__Handle,
	walletDir *C.GoString_) uint32 {
	client, ok := lookupClientHandle(handle)
	if ok {
		wf, err := client.WalletFolderName()
		if err == nil {
			copyString(wf.Address, walletDir)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_api_Handle_Client_GetWalletFileName
func SKY_api_Handle_Client_GetWalletFileName(handle C.WalletResponse__Handle,
	walletFileName *C.GoString_) uint32 {
	w, ok := lookupWalletResponseHandle(handle)
	if ok {
		copyString(w.Meta.Filename, walletFileName)
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_Client_GetWalletLabel
func SKY_api_Handle_Client_GetWalletLabel(handle C.WalletResponse__Handle,
	walletLabel *C.GoString_) uint32 {
	w, ok := lookupWalletResponseHandle(handle)
	if ok {
		copyString(w.Meta.Label, walletLabel)
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_Client_GetWalletFullPath
func SKY_api_Handle_Client_GetWalletFullPath(
	clientHandle C.Client__Handle,
	walletHandle C.WalletResponse__Handle,
	fullPath *C.GoString_) uint32 {
	client, ok := lookupClientHandle(clientHandle)
	if ok {
		wf, err := client.WalletFolderName()
		if err == nil {
			w, okw := lookupWalletResponseHandle(walletHandle)
			if okw {
				walletPath := filepath.Join(wf.Address, w.Meta.Filename)
				copyString(walletPath, fullPath)
				return SKY_OK
			}
		}
	}
	return SKY_ERROR
}

//export SKY_api_Handle_GetWalletMeta
func SKY_api_Handle_GetWalletMeta(handle C.Wallet__Handle,
	gomap *C.GoStringMap_) uint32 {
	w, ok := lookupWalletHandle(handle)
	if ok {
		copyToStringMap(w.Meta, gomap)
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_GetWalletEntriesCount
func SKY_api_Handle_GetWalletEntriesCount(handle C.Wallet__Handle,
	count *uint32) uint32 {
	w, ok := lookupWalletHandle(handle)
	if ok {
		*count = uint32(len(w.Entries))
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_Client_GetWalletResponseEntriesCount
func SKY_api_Handle_Client_GetWalletResponseEntriesCount(
	handle C.WalletResponse__Handle,
	count *uint32) uint32 {
	w, ok := lookupWalletResponseHandle(handle)
	if ok {
		*count = uint32(len(w.Entries))
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_WalletGetEntry
func SKY_api_Handle_WalletGetEntry(handle C.Wallet__Handle,
	index uint32,
	address *C.cipher__Address,
	pubkey *C.cipher__PubKey) uint32 {
	w, ok := lookupWalletHandle(handle)
	if ok {
		if index < uint32(len(w.Entries)) {
			*address = *(*C.cipher__Address)(unsafe.Pointer(&w.Entries[index].Address))
			*pubkey = *(*C.cipher__PubKey)(unsafe.Pointer(&w.Entries[index].Public))
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_api_Handle_WalletResponseGetEntry
func SKY_api_Handle_WalletResponseGetEntry(handle C.WalletResponse__Handle,
	index uint32,
	address *C.GoString_,
	pubkey *C.GoString_) uint32 {
	w, ok := lookupWalletResponseHandle(handle)
	if ok {
		if index < uint32(len(w.Entries)) {
			copyString(w.Entries[index].Address, address)
			copyString(w.Entries[index].Public, pubkey)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_api_Handle_WalletResponseIsEncrypted
func SKY_api_Handle_WalletResponseIsEncrypted(
	handle C.WalletResponse__Handle,
	isEncrypted *bool) uint32 {
	w, ok := lookupWalletResponseHandle(handle)
	if ok {
		*isEncrypted = w.Meta.Encrypted
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_WalletResponseGetCryptoType
func SKY_api_Handle_WalletResponseGetCryptoType(
	handle C.WalletResponse__Handle,
	cryptoType *C.GoString_) uint32 {
	w, ok := lookupWalletResponseHandle(handle)
	if ok {
		copyString(w.Meta.CryptoType, cryptoType)
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_WalletsResponseGetCount
func SKY_api_Handle_WalletsResponseGetCount(
	handle C.Wallets__Handle,
	count *uint32) uint32 {
	w, ok := lookupWalletsHandle(handle)
	if ok {
		*count = uint32(len(w))
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_WalletsResponseGetAt
func SKY_api_Handle_WalletsResponseGetAt(
	walletsHandle C.Wallets__Handle,
	index uint32,
	walletHandle *C.WalletResponse__Handle) uint32 {
	w, ok := lookupWalletsHandle(walletsHandle)
	if ok {
		if index < uint32(len(w)) {
			*walletHandle = registerWalletResponseHandle(w[index])
		}
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_GetWalletFolderAddress
func SKY_api_Handle_GetWalletFolderAddress(
	folderHandle C.Handle,
	address *C.GoString_) uint32 {
	obj, ok := lookupHandle(folderHandle)
	if ok {
		if obj, isOK := (obj).(*api.WalletFolder); isOK {
			copyString(obj.Address, address)
			return SKY_OK
		}
	}
	return SKY_ERROR
}

//export SKY_api_Handle_GetWalletSeed
func SKY_api_Handle_GetWalletSeed(handle C.Wallet__Handle,
	seed *C.GoString_) uint32 {
	w, ok := lookupWalletHandle(handle)
	if ok {
		copyString(w.Meta["seed"], seed)
		return SKY_OK
	}
	return SKY_ERROR
}

//export SKY_api_Handle_GetWalletLastSeed
func SKY_api_Handle_GetWalletLastSeed(handle C.Wallet__Handle,
	lastSeed *C.GoString_) uint32 {
	w, ok := lookupWalletHandle(handle)
	if ok {
		copyString(w.Meta["lastSeed"], lastSeed)
		return SKY_OK
	}
	return SKY_ERROR
}
