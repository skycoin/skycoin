
#include <stdlib.h>
#include <time.h>
#include <stdio.h>
#include <sys/stat.h>
#include <unistd.h>

#include "json.h"
#include "skytest.h"
#include "skytypes.h"

int MEMPOOLIDX = 0;
void *MEMPOOL[1024 * 256];

int JSONPOOLIDX = 0;
json_value* JSON_POOL[128];

void* registerMemCleanup(void* p) {
	int i;
	for (i = 0; i < MEMPOOLIDX; i++) {
		if(MEMPOOL[i] == NULL){
			MEMPOOL[i] = p;
			return p;
		}
	}
	MEMPOOL[MEMPOOLIDX++] = p;
	return p;
}

int registerJsonFree(void *p){
	int i;
	for (i = 0; i < JSONPOOLIDX; i++) {
		if(JSON_POOL[i] == NULL){
			JSON_POOL[i] = p;
			return i;
		}
	}
	JSON_POOL[JSONPOOLIDX++] = p;
	return JSONPOOLIDX-1;
}

void cleanupMem() {
	int i;
  void **ptr;

  for (i = MEMPOOLIDX, ptr = MEMPOOL; i; --i) {
    if( *ptr ) {
      free(*ptr);
      *ptr = NULL;
    }
    ptr++;
  }
  MEMPOOLIDX = 0;
  for (i = JSONPOOLIDX, ptr = (void*)JSON_POOL; i; --i) {
    if( *ptr ) {
      json_value_free(*ptr);
      *ptr = NULL;
    }
    ptr++;
  }
  JSONPOOLIDX = 0;
}

json_value* loadJsonFile(const char* filename){
  FILE *fp;
  struct stat filestatus;
  int file_size;
  char* file_contents;
  json_char* json;
  json_value* value;

  if ( stat(filename, &filestatus) != 0) {
    return NULL;
  }
  file_size = filestatus.st_size;
  file_contents = (char*)malloc(filestatus.st_size);
  if ( file_contents == NULL) {
    return NULL;
  }
  fp = fopen(filename, "rt");
  if (fp == NULL) {
    free(file_contents);
    return NULL;
  }
  if ( fread(file_contents, file_size, 1, fp) != 1 ) {
    fclose(fp);
    free(file_contents);
    return NULL;
  }
  fclose(fp);

  json = (json_char*)file_contents;
  value = json_parse(json, file_size);
  free(file_contents);
  return value;
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
