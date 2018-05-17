
#include <stdlib.h>
#include <time.h>
#include <stdio.h>
#include <sys/stat.h>
#include <unistd.h>
#include "json.h"
#include "skytypes.h"
#include "skytest.h"


//Define function pipe2 to avoid warning implicit declaration of function 'pipe2'
int pipe2(int pipefd[2], int flags);
//Define function SKY_handle_close to avoid including libskycoin.h
void SKY_handle_close(Handle p0);

int MEMPOOLIDX = 0;
void *MEMPOOL[1024 * 256];

int JSONPOOLIDX = 0;
json_value* JSON_POOL[128];

int HANDLEPOOLIDX = 0;
Handle HANDLE_POOL[128];

int stdout_backup;
int pipefd[2];

void * registerMemCleanup(void *p) {
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

void freeRegisteredMemCleanup(void *p){
	int i;
	for (i = 0; i < MEMPOOLIDX; i++) {
		if(MEMPOOL[i] == p){
			free(p);
			MEMPOOL[i] = NULL;
			break;
		}
	}
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

void freeRegisteredJson(void *p){
	int i;
	for (i = 0; i < JSONPOOLIDX; i++) {
		if(JSON_POOL[i] == p){
			JSON_POOL[i] = NULL;
			json_value_free( (json_value*)p );
			break;
		}
	}
}

int registerHandleClose(Handle handle){
	int i;
	for (i = 0; i < HANDLEPOOLIDX; i++) {
		if(HANDLE_POOL[i] == 0){
			HANDLE_POOL[i] = handle;
			return i;
		}
	}
	HANDLE_POOL[HANDLEPOOLIDX++] = handle;
	return HANDLEPOOLIDX - 1;
}

void closeRegisteredHandle(Handle handle){
	int i;
	for (i = 0; i < HANDLEPOOLIDX; i++) {
		if(HANDLE_POOL[i] == handle){
			HANDLE_POOL[i] = 0;
			SKY_handle_close(handle);
			break;
		}
	}
}

void cleanupMem() {
  int i;
  void **ptr;
  for (i = MEMPOOLIDX, ptr = MEMPOOL; i; --i) {
	if( *ptr )
		free(*ptr);
	ptr++;
  }
  for (i = JSONPOOLIDX, ptr = (void*)JSON_POOL; i; --i) {
	if( *ptr )
		json_value_free(*ptr);
	ptr++;
  }
  for (i = 0; i < HANDLEPOOLIDX; i++) {
	  if( HANDLE_POOL[i] )
		SKY_handle_close(HANDLE_POOL[i]);
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
