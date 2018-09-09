/* Handle not as pointer is input. */
%typemap(in) Handle {
	SWIG_AsVal_long($input, (long*)&$1);
} 
%inline %{
	typedef GoInt64_ Handle;
/**
 * Memory handle for internal object retrieving password to read
 * encrypted wallets.
 */
typedef Handle PasswordReader__Handle;

/**
 * Memory handle to perform Skycoin RPC API calls
 * encrypted wallets.
 */
typedef Handle WebRpcClient__Handle;

/**
 * Memory handle providing access to wallet data
 */
typedef Handle Wallet__Handle;

/**
 * Memory handle Options Handle
*/
typedef Handle Options__Handle;

/**
 * Memory handle to access to Skycoin CLI configuration
 */
typedef Handle Config__Handle;
/**
 * Memory handle to access to coin.Transaction
 */
typedef Handle Transaction__Handle;

/**
 * Memory handle to access to coin.Transactions
 */
typedef Handle Transactions__Handle;

/**
 * Memory handle to access to api.CreatedTransaction
 */
typedef Handle CreatedTransaction__Handle;

/**
 * Memory handle to access to api.CreatedTransactionOutput
 */
typedef Handle CreatedTransactionOutput__Handle;

/**
 * Memory handle to access to api.CreatedTransactionInput
 */
typedef Handle CreatedTransactionInput__Handle;

/**
 * Memory handle to access to api.CreateTransactionResponse
 */
typedef Handle CreateTransactionResponse__Handle;

/**
 * Memory handle to access to coin.Block
 */
typedef Handle Block__Handle;

/**
 * Memory handle to access to coin.SignedBlock
 */
typedef Handle SignedBlock__Handle;

/**
 * Memory handle to access to coin.BlockBody
 */
typedef Handle BlockBody__Handle;

/**
 * Memory handle to access to cli.BalanceResult
 */

typedef Handle BalanceResult_Handle;


/**
 * Memory handle to access to api.SpendResult
 */

typedef Handle SpendResult_Handle;

/**
 * Memory handle to access to coin.Transactions
 */

typedef Handle TransactionResult_Handle;

/**
 * Memory handle to access to coin.SortableTransactions
 */

typedef Handle SortableTransactionResult_Handle;

/**
 * Memory handle to access to wallet.Notes
 */


/**
 * Memory handle to access to wallet.ReadableNotes
 */

typedef Handle WalletReadableNotes_Handle;

/**
 * Memory handle to access to webrpc.OutputsResult
 */

typedef Handle OutputsResult_Handle;

/**
 * Memory handle to access to webrpc.StatusResult
 */

typedef Handle StatusResult_Handle;

/**
 * Memory handle to access to coin.AddressUxOuts
 */

typedef Handle AddressUxOuts_Handle;

/**
 * Memory handle to access to visor.BuildInfo (BuildInfo)
 */

typedef Handle BuildInfo_Handle;

/**
 * Memory handle for hash (ripemd160.digest)
 */

typedef Handle Hash_Handle;

/**
* Handle for Number type
*/

typedef Handle Number_Handle;

/**
* Handle for Signature type
*/

typedef Handle Signature_Handle;

%}

/**
*
* typemaps for Handles
*
**/

/* Handle reference typemap. */
%typemap(in, numinputs=0) Handle* (Handle temp) {
	$1 = &temp;
}



%apply Handle { Wallet__Handle, Options__Handle, ReadableEntry__Handle, ReadableWallet__Handle, WebRpcClient__Handle,
	WalletResponse__Handle, Client__Handle, Strings__Handle, Wallets__Handle, Config__Handle, App__Handle, Context__Handle,
	GoStringMap, PasswordReader__Handle_,
	Transaction__Handle, Transactions__Handle, CreatedTransaction__Handle,
	CreatedTransactionOutput__Handle, CreatedTransactionInput__Handle, CreateTransactionResponse__Handle,
	Block__Handle, SignedBlock__Handle, BlockBody__Handle, BuildInfo_Handle, Number_Handle, Signature_Handle, ReadableOutputSet__Handle
	}

/*GoString* parameter as reference */
%typemap(in, numinputs=0) GoString* (GoString temp) {
	temp.p = NULL;
	temp.n = 0;
	$1 = ($1_type)&temp;
}

/**
* Import library
**/
%include "typemaps.i"
%include cpointer.i

%typemap(in) (cipher_PubKey*) (cipher_PubKey temp) {
	//Typemap in *cipher_PubKey
	$1 = &temp; 
}

%typemap(in) (GoUint8_ (*)[33]) (cipher_PubKey temp) {
	//Typemap in *GoUint8_ (*)[33] is cipher_PubKey
	$1 = (cipher_PubKey*)&temp; 
}

%typemap(freearg) (GoUint8_ (*)[33]) {
	//Typemap freearg *GoUint8_ (*)[33] is cipher_PubKey
	free($1);
}


%typemap(freearg) (cipher_PubKey*) {
	//Typemap freearg *cipher_PubKey
	free($1);
}

%typemap(freearg) (GoUint8_ * [33]) {
	//Typemap freearg *cipher_PubKey
}
/**
*%typemap(ctype, in="(GoUint8_ (*) [33])") (GoUint8_ (*) [33]) "cipher_PubKey*"
*%typemap(in,in="(GoUint8_ (*) [33])") (GoUint8_ (*) [33]) "$1 = (cipher_PubKey*)$arg;"
*%typemap(in,args="(GoUint8_ (*) [33])$1") (GoUint8_ (*) [33]) "$1 = (cipher_PubKey*)$arg;"
*%typemap(cstype, in="SWIGTYPE_p_a_33__GoUint8_") (GoUint8_ (*) [33]) "cipher_PubKey"
*
*%typemap(csvarin  , in="SWIGTYPE_p_a_33__GoUint8_") (GoUint8_ (*) [33]) "cipher_PubKey"
*%typemap(imtype, in="SWIGTYPE_p_a_33__GoUint8_") (GoUint8_ (*) [33]) "cipher_PubKey"
**/