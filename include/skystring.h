#ifndef LIBSKY_STRING_H
#define LIBSKY_STRING_H

#include <stdio.h>
#include <stdlib.h>
#include "libskycoin.h"

extern void randBytes(GoSlice *bytes, size_t n);

extern void strnhex(unsigned char* buf, char *str, int n);

extern void strhex(unsigned char* buf, char *str);


#endif //LIBSKY_STRING_H
