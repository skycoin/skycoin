#ifndef LIBSKY_STRING_H
#define LIBSKY_STRING_H

#include <stdio.h>
#include <stdlib.h>
#include "libskycoin.h"
#include "skytypes.h"

extern void randBytes(GoSlice *bytes, size_t n);

extern void strnhex(unsigned char* buf, char *str, int n);

extern void strhex(unsigned char* buf, char *str);

extern void toGoString(GoString_ *s1, GoString *s2);


#endif //LIBSKY_STRING_H
