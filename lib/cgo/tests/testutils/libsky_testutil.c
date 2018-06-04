
#include <stdlib.h>
#include <time.h>
#include <stdio.h>
#include <sys/stat.h>
#include <unistd.h>
#include "json.h"
#include "skytypes.h"
#include "skytest.h"

#define BUFFER_SIZE 1024
#define stableWalletName "integration-test.wlt"
#define STRING_SIZE 128
#define JSON_FILE_SIZE 4096
#define JSON_BIG_FILE_SIZE 102400
#define TEST_DATA_DIR "src/cli/integration/testdata/"
#define stableEncryptWalletName "integration-test-encrypted.wlt"

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

typedef struct {
	Client__Handle client;
	WalletResponse__Handle wallet;
} wallet_register;

int WALLETPOOLIDX = 0;
wallet_register WALLET_POOL[64];

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

int registerWalletClean(Client__Handle clientHandle,
						WalletResponse__Handle walletHandle){
	int i;
	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].wallet == 0 && WALLET_POOL[i].client == 0){
			WALLET_POOL[i].wallet = walletHandle;
			WALLET_POOL[i].client = clientHandle;
			return i;
		}
	}
	WALLET_POOL[WALLETPOOLIDX].wallet = walletHandle;
	WALLET_POOL[WALLETPOOLIDX].client = clientHandle;
	return WALLETPOOLIDX++;
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

void cleanupWallet(Client__Handle client, WalletResponse__Handle wallet){
	int result;
	GoString_ strWalletDir;
	GoString_ strFileName;
	memset(&strWalletDir, 0, sizeof(GoString_));
	memset(&strFileName, 0, sizeof(GoString_));


	result = SKY_api_Handle_Client_GetWalletDir(client, &strWalletDir);
	if( result != SKY_OK ){
		return;
	}
	result = SKY_api_Handle_Client_GetWalletFileName(wallet, &strFileName);
	if( result != SKY_OK ){
		free( (void*)strWalletDir.p );
		return;
	}
	char fullPath[128];
	if( strWalletDir.n + strFileName.n < 126){
		strcpy( fullPath, strWalletDir.p );
		if( fullPath[0] == 0 || fullPath[strlen(fullPath) - 1] != '/' )
			strcat(fullPath, "/");
		strcat( fullPath, strFileName.p );
		result = unlink( fullPath );
		if( strlen(fullPath) < 123 ){
			strcat( fullPath, ".bak" );
			result = unlink( fullPath );
		}
	}
	GoString str = { strFileName.p, strFileName.n };
	result = SKY_api_Client_UnloadWallet( client, str );
	GoString strFullPath = { fullPath, strlen(fullPath) };
	free( (void*)strWalletDir.p );
	free( (void*)strFileName.p );
}

void cleanRegisteredWallet(
			Client__Handle client,
			WalletResponse__Handle wallet){

	int i;
	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].wallet == wallet && WALLET_POOL[i].client == client){
			WALLET_POOL[i].wallet = 0;
			WALLET_POOL[i].client = 0;
			cleanupWallet( client, wallet );
			return;
		}
	}
}

void cleanupMem() {
	int i;

	for (i = 0; i < WALLETPOOLIDX; i++) {
		if(WALLET_POOL[i].client != 0 && WALLET_POOL[i].wallet != 0){
			cleanupWallet( WALLET_POOL[i].client, WALLET_POOL[i].wallet );
		}
	}

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

int parseBoolean(const char* str, int length){
	int result = 0;
	if(length == 1){
		result = str[0] == '1' || str[0] == 't' || str[0] == 'T';
	} else {
		result = strncmp(str, "true", length) == 0 ||
			strncmp(str, "True", length) == 0 ||
			strncmp(str, "TRUE", length) == 0;
	}
	return result;
}

void toGoString(GoString_ *s, GoString *r){
GoString * tmp = r;

  *tmp = (*(GoString *) s);
}

int useCSRF() {
  GoUint32 errcode;

  GoString strCSRFVar = {"USE_CSRF", 8};
  char buffercrsf[128];
  GoString_ crsf = {buffercrsf, 0};
  errcode = SKY_cli_Getenv(strCSRFVar, &crsf);
  cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed");
  int length = strlen(crsf.p);
  int result = 0;
  if (length == 1) {
    result = crsf.p[0] == '1' || crsf.p[0] == 't' || crsf.p[0] == 'T';
  } else {
    result = strcmp(crsf.p, "true") == 0 || strcmp(crsf.p, "True") == 0 ||
             strcmp(crsf.p, "TRUE") == 0;
  }
  free((void *)crsf.p);
  return result;
}

json_value *loadGoldenFile_Cli(const char *file) {
  char path[STRING_SIZE];
  if (strlen(TEST_DATA_DIR) + strlen(file) < STRING_SIZE) {
    strcpy(path, TEST_DATA_DIR);
    strcat(path, file);
    return loadJsonFile(path);
  }
  return NULL;
}

void createTempWalletDir(bool encrypt) {
  const char *temp = "build/libskycoin/wallet-data-dir";

  int valueMkdir = mkdir(temp, S_IRWXU);

  if (valueMkdir == -1) {
    int errr = system("rm -r build/libskycoin/wallet-data-dir/*.*");
  }

  // Copy the testdata/$stableWalletName to the temporary dir.
  char walletPath[JSON_BIG_FILE_SIZE];
  if (encrypt) {
    strcpy(walletPath, stableEncryptWalletName);
  } else {
    strcpy(walletPath, stableWalletName);
  }
  unsigned char pathnameURL[BUFFER_SIZE];
  strcpy(pathnameURL, temp);
  strcat(pathnameURL, "/");
  strcat(pathnameURL, walletPath);

  FILE *rf;
  FILE *f;
  f = fopen(pathnameURL, "wb");
  unsigned char fullUrl[BUFFER_SIZE];
  strcpy(fullUrl, TEST_DATA_DIR);
  strcat(fullUrl, walletPath);
  rf = fopen(fullUrl, "rb");
  unsigned char buff[2048];
  int readBits;
  // Copy file rf to f
  if (f && rf) {
    while ((readBits = fread(buff, 1, 2048, rf)))
      fwrite(buff, 1, readBits, f);

    fclose(rf);
    fclose(f);

    GoString WalletDir = {"WALLET_DIR", 10};
    GoString Dir = {temp, strlen(temp)};
    SKY_cli_Setenv(WalletDir, Dir);
    GoString WalletPath = {"WALLET_NAME", 11};
    GoString pathname = {walletPath, strlen(walletPath)};
    SKY_cli_Setenv(WalletPath, pathname);
  }
  strcpy(walletPath, "");
};


int getCountWord(const char *str) {
  int len = 0;
  do {
    str = strpbrk(str, " "); // find separator
    if (str)
      str += strspn(str, " "); // skip separator
    ++len;                     // increment word count
  } while (str && *str);

  return len;
}
