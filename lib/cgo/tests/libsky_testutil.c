
#include <stdlib.h>
#include <time.h>

#include "skytest.h"

int MEMPOOLIDX = 0;
void *MEMPOOL[1024 * 256];

void * registerMemCleanup(void *p) {
  MEMPOOL[MEMPOOLIDX++] = p;
  return p;
}

void cleanupMem() {
  int i;
  void **ptr;
  for (i = MEMPOOLIDX, ptr = MEMPOOL; i; --i) {
    free(*ptr++);
  }
}

void setup(void) {
  srand ((unsigned int) time (NULL));
}

void teardown(void) {
  cleanupMem();
}

// TODO: Move to libsky_io.c
void fprintbuff(FILE *f, void *buff, size_t n) {
  unsigned char *ptr = (unsigned char *) buff;
  fprintf(f, "[ ");
  for (; n; --n, ptr++) {
    fprintf(f, "%02d ", *ptr);
  }
  fprintf(f, "]");
}


void toGoString(GoString_ *s, GoString *r){
GoString * tmp = r;
  
  *tmp = (*(GoString *) s);
}