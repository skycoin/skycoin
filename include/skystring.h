#ifndef LIBSKY_STRING_H
#define LIBSKY_STRING_H

#include <stdio.h>
#include <stdlib.h>
#include "libskycoin.h"

extern void randBytes(GoSlice *bytes, size_t n);

extern void strnhex(unsigned char* buf, char *str, int n);

extern void strnhexlower(unsigned char* buf, char *str, int n);

extern int hexnstr(const char* hex, unsigned char* str, int n);

extern void bin2hex(unsigned char* buf, char *str, int n);

extern int cmpGoSlice_GoSlice(GoSlice *slice1, GoSlice_ *slice2);

extern void bin2hex(unsigned char* buf, char *str, int n);

extern int string_has_suffix(const char* str, const char* suffix);

extern int string_has_prefix(const char* str, const char* prefix);

extern int count_words(const char* str, int length);

#endif //LIBSKY_STRING_H
