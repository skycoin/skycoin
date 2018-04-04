
#include <stdio.h>

#include "skytypes.h"

#ifndef LIBSKY_TESTING_H
#define LIBSKY_TESTING_H

void * registerMemCleanup(void *p);
void randBytes(GoSlice_ *bytes, size_t n);
void strnhex(unsigned char* buf, char *str, int n);
void strhex(unsigned char* buf, char *str);
void fprintbuff(FILE *f, void *buff, size_t n);

#endif

