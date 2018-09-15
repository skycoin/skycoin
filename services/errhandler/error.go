package errhandler

// SkyErr provides detailed error information occured in sky blockhain components
type SkyErr struct {
	// Code is taken fron the code constants below. Feel free to add your own constants
	Code int `json:"code"`
	// Description is a public error description which may be exposed to external service/user
	Description string `json:"description"`
	// Error field used for logging the error into the metrics/monitoring/logging system. It's intended for internal usage only,
	// not for publice exposure
	Error error `json:"-"`
}

const (
	// RPCInvalidRequest invalid RPC request
	RPCInvalidRequest = -32600
	// RPCMethodNotFound method not found
	RPCMethodNotFound = -32601
	// RPCInvalidParams invalid call params
	RPCInvalidParams = -32602
	// RPCInternalError server internal error
	RPCInternalError = -32603
	// RPCParceError parce message error
	RPCParceError = -32700

	// RPCTypeError unexpected type was passed as parameter
	RPCTypeError = -3
	// RPCInvalidAddressOrKey invalid address or key
	RPCInvalidAddressOrKey = -5
	// RPCInvalidParameter invalid, missing or duplicate parameter
	RPCInvalidParameter = -8
	// RPCDatabaceError database error
	RPCDatabaceError = -20
	// RPCDeserializationError error parsing or validating structure in raw format
	RPCDeserializationError = -22
	// RPCTransactionError error during transaction or block submission
	RPCTransactionError = -25
	// RPCTransactionRejected transaction or block was rejected by network rules
	RPCTransactionRejected = -26
	// RPCTransactionAlreadyInChain transaction already in chain
	RPCTransactionAlreadyInChain = -27
	// RPCInWarmup client still warming up
	RPCInWarmup = -28

	// RPCClientNotConnected skycoin is not connected
	RPCClientNotConnected = -9
	// RPCClientInInitialDownload Still downloading initial blocks
	RPCClientInInitialDownload = -10
	// RPCClientNodeAlreadyAdded Node is already added
	RPCClientNodeAlreadyAdded = -23
	// RPCClientNodeNotAdded node has not been added before
	RPCClientNodeNotAdded = -24
	// RPCClientNodeNotConnected node to disconnect not found in connected nodes
	RPCClientNodeNotConnected = -29
	// RPCClientInvalidIPOrSubnet invalid IP/Subnet
	RPCClientInvalidIPOrSubnet = -30

	//Wallet errors
	// RPCWalletError unspecified problem with wallet
	RPCWalletError = -4
	// RPCWalletInsufficientFunds not enough funds in account
	RPCWalletInsufficientFunds = -6 //! Not enough funds in wallet or account
	// RPCWalletInvalidAccountName invalid account name
	RPCWalletInvalidAccountName = -11
	// RPCWalletKeypoolRanOut keypool ran out, call keypoolrefill first
	RPCWalletKeypoolRanOut = -12
	// RPCWalletUnlockNeeded enter the wallet pass
	RPCWalletUnlockNeeded = -13
	// RPCWalletPassphraseIncorrect the wallet credentials entered was incorrect
	RPCWalletPassphraseIncorrect = -14
	// RPCWalletWrongEncState command given in wrong wallet encryption state
	RPCWalletWrongEncState = -15
	// RPCWalletEncryptionFailed failed to encrypt the wallet
	RPCWalletEncryptionFailed = -16
	// RPCWalletAlreadyUnlocked wallet is already unlocked
	RPCWalletAlreadyUnlocked = -17
)
