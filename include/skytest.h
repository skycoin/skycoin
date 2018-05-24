
#include <stdio.h>

#include "json.h"
#include "skycriterion.h"

#ifndef LIBSKY_TESTING_H
#define LIBSKY_TESTING_H

/*----------------------------------------------------------------------
 * I/O
 *----------------------------------------------------------------------
 */
void fprintbuff(FILE *f, void *buff, size_t n);
json_value* loadJsonFile(const char* filename);

/*----------------------------------------------------------------------
 * Memory handling
 *----------------------------------------------------------------------
 */
void * registerMemCleanup(void *p);
extern void toGoString(GoString_ *s, GoString *r);

/*----------------------------------------------------------------------
 * JSON helpers
 *----------------------------------------------------------------------
 */
json_value* json_get_string(json_value* value, const char* key);
int json_set_string(json_value* value, const char* new_string_value);
int registerJsonFree(void *p);
void freeRegisteredJson(void *p);

json_value* loadJsonFile(const char* filename);
int compareJsonValues(json_value* value1, json_value* value2);
json_value* get_json_value(json_value* node, const char* path,
							json_type type);
json_value* get_json_value_not_strict(json_value* node, const char* path,
							json_type type, int allow_null);

/*----------------------------------------------------------------------
 * JSON helpers
 *----------------------------------------------------------------------
 */
void setup(void);
void teardown(void);

#endif

