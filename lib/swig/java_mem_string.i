%rename(SKY_secp256k1go_Field_SetHex) Java_skycoin_libjava_skycoinJNI_SKY_secp256k1go_Field_SetHex;
%inline {
	GoUint32 Java_skycoin_libjava_skycoinJNI_SKY_secp256k1go_Field_SetHex(secp256k1go__Field* p0, char * p1){
        GoString str = { p1,strlen(p1) };
		GoUint32 result = SKY_secp256k1go_Field_SetHex(p0,  str);
		return result;
	}
}