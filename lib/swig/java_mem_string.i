%rename(SKY_secp256k1go_Field_SetHex) Java_skycoin_libjava_skycoinJNI_SKY_secp256k1go_Field_SetHex;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_secp256k1go_Field_SetHex(secp256k1go__Field* p0, char * p1){
        GoString str = { p1,strlen(p1) };
		GoUint32 result = SKY_secp256k1go_Field_SetHex(p0,  str);
		return result;
	}
}

%rename(SKY_secp256k1go_Number_SetHex) Java_skycoin_libjava_skycoinJNI_SKY_secp256k1go_Number_SetHex;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_secp256k1go_Number_SetHex(Number_Handle p0, char * p1){
        GoString str = { p1,strlen(p1) };
		GoUint32 result = SKY_secp256k1go_Number_SetHex(p0,  str);
		return result;
	}
}

%rename(SKY_base58_String2Hex) Java_skycoin_libjava_skycoinJNI_SKY_base58_String2Hex;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_base58_String2Hex(char * p1,GoSlice_ * p2){
        GoString str = { p1,strlen(p1) };
		GoUint32 result = SKY_base58_String2Hex(str,  p2);
		return result;
	}
}