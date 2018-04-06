
#include <stdio.h>


#include "skycriterion.h"

#ifndef LIBSKY_TESTING_H
#define LIBSKY_TESTING_H

void * registerMemCleanup(void *p);
void fprintbuff(FILE *f, void *buff, size_t n);
extern void toGoString(GoString_ *s, GoString *r);

#endif

