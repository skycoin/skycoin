//nolint
// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package consensus

//"github.com/skycoin/skycoin/src/cipher"

////////////////////////////////////////////////////////////////////////////////
type ConnectionManagerInterface interface {
	SendBlockToAllMySubscriber(blockPtr *BlockBase)

	Print() // For debugging

	// IMPORTANT: When connection manager (i.e. an implementation of
	// this interface) receives a message with 'BlockBase', the
	// manager should call
	//
	//    ConsensusParticipant.OnBlockHeaderArrived(blockPtr *BlockBase)
	//
	// function. This is not currently enforced, but is required for the
	// consensus to work properly.
}

////////////////////////////////////////////////////////////////////////////////
