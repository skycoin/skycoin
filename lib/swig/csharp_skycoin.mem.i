%rename(SKY_cipher_SumSHA256) CSharp_skycoin_SKY_cipher_SumSHA256;
%inline {
	GoUint32 CSharp_skycoin_SKY_cipher_SumSHA256(GoSlice seed, cipher_SHA256* sha){
		GoUint32 result = SKY_cipher_SumSHA256(seed,  sha);
		return result;
	}
}

%rename(SKY_cipher_SignHash) CSharp_skycoin_SKY_cipher_SignHash;
%inline {
	GoUint32 CSharp_skycoin_SKY_cipher_SignHash(cipher_SHA256 *sha,cipher__SecKey *sec,cipher_Sig *s){
		GoUint32 result = SKY_cipher_SignHash(sha,sec,s);
		return result;
	}
}

%rename(SKY_cipher_ChkSig) CSharp_skycoin_SKY_cipher_ChkSig;
%inline {
	GoUint32 CSharp_skycoin_SKY_cipher_ChkSig(cipher__Address *a,cipher_SHA256 *sha,cipher_Sig *s){
		GoUint32 result = SKY_cipher_ChkSig(a,sha,s);
		return result;
	}
}

%rename(SKY_cipher_PubKeyFromSig) CSharp_skycoin_SKY_cipher_PubKeyFromSig;
%inline {
	GoUint32 CSharp_skycoin_SKY_cipher_PubKeyFromSig(cipher_Sig *sig,cipher_SHA256 *h,cipher_PubKey *p){
		GoUint32 result = SKY_cipher_PubKeyFromSig(sig,h,p);
		return result;
	}
}

%rename(SKY_cipher_VerifySignature) CSharp_skycoin_SKY_cipher_VerifySignature;
%inline {
	GoUint32 CSharp_skycoin_SKY_cipher_VerifySignature(cipher_PubKey *p,cipher_Sig *sig,cipher_SHA256 *h){
		GoUint32 result = SKY_cipher_VerifySignature(p,sig,h);
		return result;
	}
}

%rename(SKY_cipher_TestSecKeyHash) CSharp_skycoin_SKY_cipher_TestSecKeyHash;
%inline {
	GoUint32 CSharp_skycoin_SKY_cipher_TestSecKeyHash(cipher_SecKey *s,cipher_SHA256 *h){
		GoUint32 result = SKY_cipher_TestSecKeyHash(s,h);
		return result;
	}
}