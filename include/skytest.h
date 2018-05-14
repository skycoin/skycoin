
#include <stdio.h>
#include "json.h"

#include "skycriterion.h"

#ifndef LIBSKY_TESTING_H
#define LIBSKY_TESTING_H

void * registerMemCleanup(void *p);
void fprintbuff(FILE *f, void *buff, size_t n);
void redirectStdOut();
int getStdOut(char* str, unsigned int max_size);
json_value* json_get_string(json_value* value, const char* key);
int json_set_string(json_value* value, const char* new_string_value);
void registerJsonFree(void *p);
json_value* loadJsonFile(const char* filename);
int compareJsonValues(json_value* value1, json_value* value2);
extern void toGoString(GoString_ *s, GoString *r);

#endif

