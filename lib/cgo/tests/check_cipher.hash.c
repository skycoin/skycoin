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
  Ripemd160 tmp;
  Ripemd160 r;
  Ripemd160 r2;
  unsigned char buff[257];
  GoSlice slice = { buff, 0, 257 };

  randBytes(&slice,128);
  SKY_cipher_HashRipemd160(slice,&tmp);
  randBytes(&slice,160);
  SKY_cipher_HashRipemd160(slice,&r);
  cr_assert(not(eq(u8[sizeof(Ripemd160)],tmp,r)));

  unsigned char buff1[257];
  GoSlice b = { buff1, 0, 257 };
  randBytes(&b,256);
  SKY_cipher_HashRipemd160(b,&r2);
  cr_assert(not(eq(u8[sizeof(Ripemd160)],r2,tmp)));
  freshSumRipemd160(b,&tmp);
  cr_assert(eq(u8[20],tmp,r2));
}

Test(cipher_hash,TestRipemd160Set){

  Ripemd160 h;
  unsigned char buff[101];
  GoSlice slice = { buff, 0, 101 };
  int error;

  memset(h, 0, sizeof(Ripemd160));
  randBytes(&slice,21);

  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert( error == SKY_ERROR);

  randBytes(&slice,100);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,19);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,0);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_ERROR);

  randBytes(&slice,20);
  error = SKY_cipher_Ripemd160_Set(&h,slice);
  cr_assert(error == SKY_OK);
  cr_assert(eq(u8[20], h, buff));
}

Test(cipher_hash,TestSHA256Set){

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

Test(cipher_hash,TestSHA256Hex){

  SHA256 h;
  unsigned char buff[101];
  GoSlice slice = { buff, 0, 101 };
  int error;

  memset(&h, 0, sizeof(h));
  randBytes(&slice,32);
  SKY_cipher_SHA256_Set(&h,slice);
  GoString s;

  SKY_cipher_SHA256_Hex(&h, (GoString_ *)&s);
  registerMemCleanup(&s.p);

  SHA256 h2;

  error = SKY_cipher_SHA256FromHex(s, &h2 );
  cr_assert(error == SKY_OK);
  cr_assert(eq(u8[32],h,h2));

  GoString s2;

  SKY_cipher_SHA256_Hex(&h2, (GoString_ *) &s2);
  registerMemCleanup(&s2.p);
  cr_assert(eq(type(GoString),s,s2));
}

Test(cipher_hash,TestSHA256KnownValue){


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
    GoSlice slice_input;
    GoSlice slice_output;

    slice_input.data = vals[i].input;
    slice_input.len = strlen(vals[i].input);
    slice_input.cap = strlen(vals[i].input)+1;

    SHA256 sha;

    SKY_cipher_SumSHA256(slice_input,&sha);

    GoString_ tmp_output;

    SKY_cipher_SHA256_Hex(&sha,&tmp_output);
    registerMemCleanup(&tmp_output.p);

    cr_assert(strcmp(tmp_output.p,vals[i].output)== SKY_OK);
  }
}

Test(cipher_hash,TestSumSHA256){

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

Test(cipher_hash,TestSHA256FromHex){
  unsigned int error;
  SHA256 tmp;
  // Invalid hex hash
  GoString tmp_string = {"cawcd",5};
  error = SKY_cipher_SHA256FromHex(tmp_string,&tmp);
  cr_assert(error == SKY_ERROR);
  // Truncated hex hash
  SHA256 h;
  unsigned char buff[130];
  char sbuff[300];
  GoSlice slice = { buff,0,130 };
  randBytes(&slice,128);
  SKY_cipher_SumSHA256(slice,&h);
  strnhex(h,sbuff,sizeof(h) >> 1);
  GoString s1 = { sbuff, strlen(sbuff) };
  error = SKY_cipher_SHA256FromHex(s1,&h);
  cr_assert(error == SKY_ERROR);

  // Valid hex hash
  // char sbuff1[300];
  GoString_ s2;
  // strnhex(h,sbuff1,sizeof(h));
  SKY_cipher_SHA256_Hex(&h, &s2 );
  SHA256 h2;
  error = SKY_cipher_SHA256FromHex((*((GoString *) &s2)),&h2);
  cr_assert(error == SKY_OK);
  cr_assert(eq(u8[32],h,h2));
}


Test(cipher_hash,TestDoubleSHA256){
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

Test(cipher_hash,TestAddSHA256){

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

  SKY_cipher_AddSHA256(&h,&i,&add);

  cr_assert(not(eq(u8[32],add,tmp)));
  cr_assert(not(eq(u8[32],add,h)));
  cr_assert(not(eq(u8[32],add,i)));
}

Test(cipher_hash,TestXorSHA256){

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

Test(cipher_hash,TestMerkle){
  unsigned char buff[129];
  SHA256 hashlist[5];
  GoSlice b = { buff, 0, 129 },
          hashes = { hashlist, 0, 5 };
  SHA256 h, zero, out, out1, out2, out3, out4;
  int i;

  memset(zero, 0, sizeof(zero));

  for (i = 0; i < 5; i++) {
    randBytes(&b, 128);
    SKY_cipher_SumSHA256(b, &hashlist[i]);
  }

  // Single hash input returns hash
  hashes.len = 1;
  SKY_cipher_Merkle(&hashes, &h);
  cr_assert(eq(u8[32], hashlist[0], h));

  // 2 hashes should be AddSHA256 of them
  hashes.len = 2;
  SKY_cipher_AddSHA256(&hashlist[0], &hashlist[1], &out); 
  SKY_cipher_Merkle(&hashes, &h);
  cr_assert(eq(u8[32], out, h));

  // 3 hashes should be Add(Add())
  hashes.len = 3;
  SKY_cipher_AddSHA256(&hashlist[0], &hashlist[1], &out1); 
  SKY_cipher_AddSHA256(&hashlist[2], &zero, &out2); 
  SKY_cipher_AddSHA256(&out1, &out2, &out); 
  SKY_cipher_Merkle(&hashes, &h);
  cr_assert(eq(u8[32], out, h));

  // 4 hashes should be Add(Add())
  hashes.len = 4;
  SKY_cipher_AddSHA256(&hashlist[0], &hashlist[1], &out1); 
  SKY_cipher_AddSHA256(&hashlist[2], &hashlist[3], &out2); 
  SKY_cipher_AddSHA256(&out1, &out2, &out); 
  SKY_cipher_Merkle(&hashes, &h);
  cr_assert(eq(u8[32], out, h));

  // 5 hashes
  hashes.len = 5;
  SKY_cipher_AddSHA256(&hashlist[0], &hashlist[1], &out1); 
  SKY_cipher_AddSHA256(&hashlist[2], &hashlist[3], &out2); 
  SKY_cipher_AddSHA256(&out1, &out2, &out3); 
  SKY_cipher_AddSHA256(&hashlist[4], &zero, &out1); 
  SKY_cipher_AddSHA256(&zero, &zero, &out2); 
  SKY_cipher_AddSHA256(&out1, &out2, &out4); 
  SKY_cipher_AddSHA256(&out3, &out4, &out); 
  SKY_cipher_Merkle(&hashes, &h);
  cr_assert(eq(u8[32], out, h));
}

