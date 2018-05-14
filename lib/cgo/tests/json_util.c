#include "json.h"
#include <string.h>

json_value* json_get_string(json_value* value, const char* key, char* str, int max_size){
	int length, x;
	if (value == NULL) {
			return NULL;
	}
	if (value->type != json_object) {
		return NULL;
	}
	length = value->u.object.length;
	for (x = 0; x < length; x++) {
		if( strcmp( value->u.object.values[x].name, key) == 0){
			if( value->u.object.values[x].value->type == json_string){
				char* p = value->u.object.values[x].value->u.string.ptr;
				int string_length = value->u.object.values[x].value->u.string.length;
				if( string_length >= max_size )
					string_length = max_size - 1;
				strncpy( str, p, string_length );
				str[string_length] = 0;
				return value->u.object.values[x].value;
			}
		}
	}
	return NULL;
}

int json_set_string(json_value* value, const char* new_string_value){
	if( value->type == json_string){
		int length = strlen(new_string_value);
		if( length > value->u.string.length ){
			value->u.string.ptr = malloc(length + 1);
		}
		strcpy( value->u.string.ptr, new_string_value );
		value->u.string.length = length;
	}
	return 0;
}