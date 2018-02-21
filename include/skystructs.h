
#ifndef SKYSTRUCTS_H
#define SKYSTRUCTS_H

typedef unsigned char Ripemd160[20];

typedef struct {
	unsigned char Version;
	Ripemd160 Key;
} Address;

typedef unsigned char Checksum[4];

typedef unsigned char PubKey[33];
typedef unsigned char SecKey[32];
typedef unsigned char Checksum[4];

#endif

