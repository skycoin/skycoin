/* Handle not as pointer is input. */
%typemap(in) Handle {
	$input =  (long*)&$1;
} 
%typemap(in) Handle* {
	$input =  (long*)&$1;
} 

%apply Handle { Wallet__Handle, Options__Handle, ReadableEntry__Handle, ReadableWallet__Handle, WebRpcClient__Handle,
	WalletResponse__Handle, Client__Handle, Strings__Handle, Wallets__Handle, Config__Handle, App__Handle, Context__Handle,
	GoStringMap, PasswordReader__Handle_,
	Transaction__Handle, Transactions__Handle, CreatedTransaction__Handle,
	CreatedTransactionOutput__Handle, CreatedTransactionInput__Handle, CreateTransactionResponse__Handle,
	Block__Handle, SignedBlock__Handle, BlockBody__Handle, BuildInfo_Handle, Number_Handle, Signature_Handle, ReadableOutputSet__Handle
	}

%apply Handle* { Wallet__Handle*, Options__Handle*, ReadableEntry__Handle*, ReadableWallet__Handle*, WebRpcClient__Handle*,
	WalletResponse__Handle*, Client__Handle*, Strings__Handle*, Wallets__Handle*, Config__Handle*, App__Handle*, Context__Handle*,
	GoStringMap*, PasswordReader__Handle_*,
	Transaction__Handle*, Transactions__Handle*, CreatedTransaction__Handle*,
	CreatedTransactionOutput__Handle*, CreatedTransactionInput__Handle*, CreateTransactionResponse__Handle*,
	Block__Handle*, SignedBlock__Handle*, BlockBody__Handle*, BuildInfo_Handle*, Number_Handle*, Signature_Handle*, ReadableOutputSet__Handle*
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
// Pubkey
%typemap(ctype,pre="cipher_PubKey tmp$csinput = new_cipher_PubKeyp();") (GoUint8_ (*) [33])  "cipher__PubKey*"
%typemap(cstype,pre="var tmp$csinput = cipher_PubKey.getCPtr ($csinput);") (GoUint8_ (*) [33])  "cipher_PubKey"
%typemap(csin,pre="var tmp$csinput = cipher_PubKey.getCPtr ($csinput);") (GoUint8_ (*) [33])  "tmp$csinput"


// Seckey
%typemap(ctype,pre="cipher_SecKey tmp$csinput = new_cipher_SecKeyp();") (GoUint8_ (*) [32])  "cipher_SecKey*"
%typemap(cstype,pre=" var tmp$csinput = cipher_SecKey.getCPtr ($csinput);") (GoUint8_ (*) [32])  "cipher_SecKey"
%typemap(csin,pre="var tmp$csinput = cipher_SecKey.getCPtr ($csinput);") (GoUint8_ (*) [32])  "tmp$csinput"

// Sig
%typemap(ctype,pre="cipher_Sig tmp$csinput = new cipher_Sig();") (GoUint8_ (*) [65])  "cipher_Sig*"
%typemap(cstype,pre=" var tmp$csinput = cipher_Sig.getCPtr ($csinput);") (GoUint8_ (*) [65])  "cipher_Sig"
%typemap(csin,pre="var tmp$csinput = cipher_Sig.getCPtr ($csinput);") (GoUint8_ (*) [65])  "tmp$csinput"

// cipher__Ripemd160
%typemap(ctype,pre="cipher__Ripemd160 tmp$csinput = new_cipher_Ripemd160p();") (GoUint8_ (*) [20])  "cipher_Ripemd160*"
%typemap(cstype,pre=" var tmp$csinput = cipher_Ripemd160.getCPtr ($csinput);") (GoUint8_ (*) [20])  "cipher_Ripemd160"
%typemap(csin,pre="var tmp$csinput = cipher_Ripemd160.getCPtr ($csinput);") (GoUint8_ (*) [20])  "tmp$csinput"



// GoString
%typemap(cstype,pre=" var tmp$csinput = $csinput;") GoString "string"
%typemap(csin,pre="var tmp$csinput = $csinput;") GoString  "tmp$csinput"
%typemap(imtype,pre="var tmp$csinput  = $csinput;") GoString  "string"
%typemap(ctype) GoString  "char*"
%typemap(in) GoString  "$1.p=$input;$1.n=strlen($input);"

%typemap(ctype,pre="GoString_ tmp$csinput = new_GoStringp_();") GoString_*  "GoString*"
%typemap(cstype,pre=" var tmp$csinput = _GoString_.getCPtr ($csinput);") GoString_*  "_GoString_"
%typemap(csin,pre="var tmp$csinput = _GoString_.getCPtr ($csinput);") GoString_*  "tmp$csinput"

// // GoSlice
%typemap(ctype) GoSlice_*  "GoSlice_ *"
%typemap(cstype,pre=" var tmp$csinput = GoSlice.getCPtr ($csinput);") GoSlice_*  "GoSlice"
%typemap(csin) GoSlice_*  "GoSlice.getCPtr ($csinput)"

%apply int {GoInt_};
%apply char {byte,GoUint8_};





