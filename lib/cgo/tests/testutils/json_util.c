#include "json.h"
#include <string.h>
#include <math.h>

json_value* json_get_string(json_value* value, const char* key){
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

int _compareJsonValues(json_value* value1, json_value* value2, const char* ignore);

int compareJsonObjects(json_value* value1, json_value* value2, 
						const char* ignore){
	int length1 = value1->u.object.length;
	int length2 = value2->u.object.length;
	/*if( length1 != length2 )
		return 0;*/
	for (int x = 0; x < length1; x++) {
		char* name = value1->u.object.values[x].name;
		if( ignore != NULL && strcmp( ignore, name ) == 0)
			continue;
		int found = 0;
		for( int y = 0; y < length2; y++){
			if( strcmp( value2->u.object.values[y].name, name ) == 0){
				if( !_compareJsonValues( value1->u.object.values[x].value,
								value2->u.object.values[y].value, ignore )  )
					return 0;
				found = 1;
				break;
			}
		}
		if( !found )
			return 0;
	}
	return 1;
}

int compareJsonArrays(json_value* value1, json_value* value2, const char* ignore){
	int length1 = value1->u.array.length;
	int length2 = value2->u.array.length;
	if( length1 != length2 )
		return 0;
	for (int x = 0; x < length1; x++) {
		if( !_compareJsonValues(value1->u.array.values[x], 
				value2->u.array.values[x], ignore) )
			return 0;
	}
	return 1;
}

int _compareJsonValues(json_value* value1, json_value* value2, const char* ignore){
	if( value1 == NULL && value2 == NULL)
		return 1;
	if( value1 == NULL || value2 == NULL)
		return 0;
	if( value1->type != value2->type)
		return 0;
	switch (value1->type) {
    case json_null:
      return value2->type == json_null;
		case json_none:
			return 1;
		case json_object:
			return compareJsonObjects(value1, value2, ignore);
		case json_array:
			return compareJsonArrays(value1, value2, ignore);
		case json_integer:
			return value1->u.integer == value2->u.integer;
		case json_double:
			return fabs(value1->u.dbl - value2->u.dbl) < 0.000001;
		case json_string:
			return strcmp(value1->u.string.ptr, value2->u.string.ptr) == 0;
		case json_boolean:
			return value1->u.boolean == value2->u.boolean;
	}
	return 1;
}

int compareJsonValues(json_value* value1, json_value* value2){
	return _compareJsonValues(value2, value1, NULL);
}

int compareJsonValuesWithIgnoreList(json_value* value1, json_value* value2, const char* ignoreList){
	return _compareJsonValues(value2, value1, ignoreList);
}


json_value* get_json_value_not_strict(json_value* node, const char* path,
							json_type type, int allow_null){
	int n;
	const char* p = strchr(path, '/');
	if( p == NULL )
		n = strlen(path);
	else
		n = p - path;
	if( n > 0 ) {
		if( node->type == json_object){
			for (int x = 0; x < node->u.object.length; x++) {
				json_object_entry * entry = &node->u.object.values[x];
				char* name = entry->name;
				json_value* value = entry->value;
				if( strncmp( path, name, n ) == 0){
					if( p == NULL){
						if( value->type == type ||
								(allow_null && value->type == json_null))
							return value;
					}else
						return get_json_value_not_strict(
							value, p + 1, type, allow_null);
				}
			}
		} else {
			return NULL;
		}
	}
	return NULL;
}

json_value* get_json_value(json_value* node, const char* path,
							json_type type){
	return get_json_value_not_strict(node, path, type, 1);
}
