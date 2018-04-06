#include <stdio.h>
#include <string.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

void freshSumRipemd160(GoSlice bytes, Ripemd160* rp160){

  SKY_cipher_HashRipemd160(bytes, rp160);
}

void freshSumSHA256(GoSlice bytes, SHA256* sha256){

  SKY_cipher_SumSHA256(bytes, sha256);
}

Test(cipher,TestHashRipemd160){

  Ripemd160 rp2;
  Ripemd160 rp3;
  Ripemd160 rp1;

  unsigned char buff[129];
  GoSlice slice = { buff, 0, 129 };

  randBytes(&slice,128);

  SKY_cipher_HashRipemd160(slice,&rp1);

  unsigned char buff1[170];
  GoSlice slice1 = { buff1, 0, 170 };

  randBytes(&slice1,160);

  SKY_cipher_HashRipemd160(slice1,&rp2);

  cr_assert( strcmp((char *)rp1,(char *)rp2) !=0   );

  unsigned char buff2[257];
  GoSlice slice2;

  randBytes(&slice2,256);

  SKY_cipher_HashRipemd160(slice2,&rp3);

  cr_assert(eq(u8[32],rp1,rp2));

  freshSumRipemd160(slice2,&rp3);

  cr_assert(eq(u8[32],rp1,rp2));
}

Test(hash,TestRipemd160Set){

  Ripemd160 h;
  unsigned char buff[101];
  GoSlice slice;
  int error;

  randBytes(&slice,21);

  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert( error == SKY_OK);

  randBytes(&slice,100);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_OK);

  randBytes(&slice,19);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_OK);

  randBytes(&slice,0);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_OK);

  randBytes(&slice,20);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_OK);

  unsigned char buff1[101];
  GoSlice slice1 = { buff1, 0, 101 };

  randBytes(&slice1,20);
  error = SKY_cipher_Ripemd160_Set(&h,slice1);

  cr_assert(eq(type(GoSlice),slice1,slice));
}

Test(hash,TestSHA256Set){

  SHA256 h;
  unsigned char buff[101];
  GoSlice slice = { buff, 0, 101 };
  int error;

  randBytes(&slice,33);
  error=SKY_cipher_SHA256_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,100);
  error=SKY_cipher_SHA256_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,31);
  error=SKY_cipher_SHA256_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,0);
  error=SKY_cipher_SHA256_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,32);
  error=SKY_cipher_SHA256_Set(&h,slice);
  cr_assert(error == SKY_OK);

  cr_assert(eq(u8[32], h, slice.data));
}

Test(hash,TestSHA256Hex){

  SHA256 h;
  unsigned char buff[101];
  GoSlice slice = { buff, 0, 101 };
  int error;

  memset(&h, 0, sizeof(h));
  randBytes(&slice,32);
  SKY_cipher_SHA256_Set(&h,slice);
  GoString_ s;

  SKY_cipher_SHA256_Hex(&h,&s);
  registerMemCleanup(s.p);

  SHA256 h2;

  error = SKY_cipher_SHA256FromHex( (*((GoString*)&s)),&h2 );

  GoString_ s2;

  SKY_cipher_SHA256_Hex(&h2,&s2);
  registerMemCleanup(s2.p);

  cr_assert(eq(u8[32],h,h2));

  cr_assert(eq(type(GoString_),s,s2));

}

Test(hash,TestSHA256KnownValue){


  typedef struct 
  {
    char *input;
    char *output;
  } tmpstruct;

  tmpstruct vals[3];

  vals[0].input = "skycoin";
  vals[0].output = "5a42c0643bdb465d90bf673b99c14f5fa02db71513249d904573d2b8b63d353d";

  vals[1].input = "hello world";
  vals[1].output = "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9";

  vals[2].input = "hello world asd awd awd awdapodawpokawpod ";
  vals[2].output = "99d71f95cafe05ea2dddebc35b6083bd5af0e44850c9dc5139b4476c99950be4";

  for (int i = 0; i < 3; ++i)
  {
    // FIXME: Review
    GoSlice slice_input;
    GoSlice slice_output;

    slice_input.data = vals[i].input;
    slice_input.len = strlen(vals[i].input);
    slice_input.cap = strlen(vals[i].input)+1;

    SHA256 sha;

    SKY_cipher_SumSHA256(slice_input,&sha);

    GoString_ tmp_output;

    SKY_cipher_SHA256_Hex(&sha,&tmp_output);
    registerMemCleanup(tmp_output.p);

    cr_assert( tmp_output.p == vals[i].output );
  }

}

Test(hash,TestSumSHA256){

  unsigned char bbuff[257],
  cbuff[257];
  GoSlice b = { bbuff, 0, 257 };
  SHA256 h1;

  randBytes(&b,256);

  SKY_cipher_SumSHA256(b,&h1);

  SHA256 tmp;

  cr_assert(not(eq(u8[32],h1,tmp)));

  GoSlice c = { cbuff, 0, 257 };

  randBytes(&c,256);
  SHA256 h2;
  SKY_cipher_SumSHA256(c,&h2);


  cr_assert(not(eq(u8[32],h2,tmp)));

  SHA256 tmp_h2;

  freshSumSHA256(c,&tmp_h2);

  cr_assert(eq(u8[32],h2,tmp_h2));

}

Test(hash,TestSHA256FromHex){

  // Invalid hex hash
  GoString string;
  int error;
  SHA256 sha;
  string.p = "cawcd";
  string.n = 5;
  error = SKY_cipher_SHA256FromHex(string,&sha);

  cr_assert(error != NULL);

  // 	// Truncated hex hash
  SHA256 sha1;
  unsigned char buff[130],
  sbuff[300];
  GoSlice slice = { buff, 0, 130 };

  randBytes(&slice,128);
  SKY_cipher_SumSHA256(slice,&sha1);
  int len = sizeof(sha1);

  strnhex(&sha1,&sbuff,len/2);
  GoString s1 = { sbuff, strlen(sbuff) };

  error = SKY_cipher_SHA256FromHex(s1,&sha1);
  cr_assert(error != NULL);

  // Valid hex hash

  strhex(&sha1,&sbuff);
  SHA256 h2;
  GoString s2 = { sbuff, strlen(sbuff) };

  error = SKY_cipher_SHA256FromHex(s2,&h2);

  cr_assert(eq(u8[32],sha1,h2));

  cr_assert(error == NULL);

}

Test(hash,TestDoubleSHA256){
  unsigned char bbuff[130];
  GoSlice b = { bbuff, 0, 130 };
  randBytes(&b,128);
  SHA256 h;
  SHA256 tmp;
  SKY_cipher_DoubleSHA256(b,&h);
  cr_assert(not(eq(u8[32],tmp,h)));
  freshSumSHA256(b,&tmp);
  cr_assert(not(eq(u8[32],tmp,h)));
}

Test(hash,TestAddSHA256){

  unsigned char bbuff[130];
  GoSlice b = { bbuff, 0, 130 };
  randBytes(&b,128);
  SHA256 h;
  SKY_cipher_SumSHA256(b,&h);

  unsigned char cbuff[130];
  GoSlice c = { cbuff, 0, 130 };
  randBytes(&c,64);
  SHA256 i;
  SKY_cipher_SumSHA256(c,&i);

  SHA256 add;
  SHA256 tmp;

  SKY_cipher_AddSHA256(h,i,&add);

  cr_assert(not(eq(u8[32],add,tmp)));
  cr_assert(not(eq(u8[32],add,h)));
  cr_assert(not(eq(u8[32],add,i)));
}

Test(hash,TestXorSHA256){

  unsigned char bbuff[129],
  cbuff[129];
  GoSlice b = { bbuff, 0, 129 } ;
  GoSlice c = { cbuff, 0, 129 };
  SHA256 h, i;

  randBytes(&b,128);
  SKY_cipher_SumSHA256(b,&h);
  randBytes(&c,128);
  SKY_cipher_SumSHA256(c,&i);

  SHA256 tmp_xor1;
  SHA256 tmp_xor2;
  SHA256 tmp;

  SKY_cipher_SHA256_Xor(&h,&i,&tmp_xor1);
  SKY_cipher_SHA256_Xor(&i,&h,&tmp_xor2);

  cr_assert(not(eq(u8[32],tmp_xor1,h)));
  cr_assert(not(eq(u8[32],tmp_xor1,i)));
  cr_assert(not(eq(u8[32],tmp_xor1,tmp)));
  cr_assert(eq(u8[32],tmp_xor1,tmp_xor2));

}

Test(hash,TestMerkle){

  cr_fail("Not implement");

  // GoSlice slice;
  // SHA256 h;
  // SHA256 tmp;

  // randBytes(&slice,128);

  // SKY_cipher_SumSHA256(slice,&h);

  // // Single hash input returns hash

  // assert.Equal(t, Merkle([]SHA256{h}), h)
  // h2 := SumSHA256(randBytes(t, 128))
  // // 2 hashes should be AddSHA256 of them
  // assert.Equal(t, Merkle([]SHA256{h, h2}), AddSHA256(h, h2))
  // // 3 hashes should be Add(Add())
  // h3 := SumSHA256(randBytes(t, 128))
  // out := AddSHA256(AddSHA256(h, h2), AddSHA256(h3, SHA256{}))
  // assert.Equal(t, Merkle([]SHA256{h, h2, h3}), out)
  // // 4 hashes should be Add(Add())
  // h4 := SumSHA256(randBytes(t, 128))
  // out = AddSHA256(AddSHA256(h, h2), AddSHA256(h3, h4))
  // assert.Equal(t, Merkle([]SHA256{h, h2, h3, h4}), out)
  // // 5 hashes
  // h5 := SumSHA256(randBytes(t, 128))
  // out = AddSHA256(AddSHA256(h, h2), AddSHA256(h3, h4))
  // out = AddSHA256(out, AddSHA256(AddSHA256(h5, SHA256{}),
  // 	AddSHA256(SHA256{}, SHA256{})))
  // assert.Equal(t, Merkle([]SHA256{h, h2, h3, h4, h5}), out)
}
