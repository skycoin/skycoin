
#include <stdlib.h>
#include <time.h>
#include <stdio.h>
#include <unistd.h>

#include "skytest.h"

int MEMPOOLIDX = 0;
void *MEMPOOL[1024 * 256];

int stdout_backup;
int pipefd[2];

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

void redirectStdOut(){
	stdout_backup = dup(fileno(stdout));
	pipe2(pipefd, 0);
	dup2(pipefd[1], fileno(stdout));
}

int getStdOut(char* str, unsigned int max_size){
	fflush(stdout);
	close(pipefd[1]);
	dup2(stdout_backup, fileno(stdout));
	int bytes_read = read(pipefd[0], str, max_size - 1);
	if( bytes_read > 0 && bytes_read < max_size)
		str[bytes_read] = 0;
	close(pipefd[0]);
	return bytes_read;
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