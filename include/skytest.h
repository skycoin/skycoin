
#include <stdio.h>
#include "json.h"
#include "skytypes.h"
#include "skycriterion.h"

#ifndef LIBSKY_TESTING_H
#define LIBSKY_TESTING_H

void * registerMemCleanup(void *p);

void fprintbuff(FILE *f, void *buff, size_t n);

void redirectStdOut();

int getStdOut(char* str, unsigned int max_size);

json_value* json_get_string(json_value* value, const char* key);

int json_set_string(json_value* value, const char* new_string_value);

int registerJsonFree(void *p);

void freeRegisteredJson(void *p);

int registerHandleClose(Handle handle);

void closeRegisteredHandle(Handle handle);

void freeRegisteredMemCleanup(void *p);

int registerWalletClean(Client__Handle clientHandle,
						WalletResponse__Handle walletHandle);
						
void cleanRegisteredWallet(
			Client__Handle client, 
			WalletResponse__Handle wallet);

json_value* loadJsonFile(const char* filename);

int compareJsonValues(json_value* value1, json_value* value2);

json_value* get_json_value(json_value* node, const char* path,
							json_type type);
							
json_value* get_json_value_not_strict(json_value* node, const char* path,
							json_type type, int allow_null);
														
extern void toGoString(GoString_ *s, GoString *r);

#endif

