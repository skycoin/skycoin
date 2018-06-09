
#include "skystring.h"

#define ALPHANUM "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
#define ALPHANUM_LEN 62
#define SIZE_ALL -1

void randBytes(GoSlice *bytes, size_t n) {
	size_t i = 0;
	unsigned char *ptr = (unsigned char *) bytes->data;
	for (; i < n; ++i, ++ptr) {
		*ptr = ALPHANUM[rand() % ALPHANUM_LEN];
	}
	bytes->len = (GoInt) n;
}

void strnhex(unsigned char* buf, char *str, int n){
	unsigned char * pin = buf;
	const char * hex = "0123456789ABCDEF";
	char * pout = str;
	for(; n; --n){
		*pout++ = hex[(*pin>>4)&0xF];
		*pout++ = hex[(*pin++)&0xF];
	}
	*pout = 0;
}

void strnhexlower(unsigned char* buf, char *str, int n){
	unsigned char * pin = buf;
	const char * hex = "0123456789abcdef";
	char * pout = str;
	for(; n; --n){
		*pout++ = hex[(*pin>>4)&0xF];
		*pout++ = hex[(*pin++)&0xF];
	}
	*pout = 0;
}


int hexnstr(const char* hex, unsigned char* str, int n){
	const char * pin = hex;
	unsigned char * pout = str;
	unsigned char c;
	int odd = 0;
	int size = 0;
	for(; *pin && size < n; pin++){
		if(*pin >= '0' && *pin <= '9'){
			c = *pin - '0';
		} else if(*pin >= 'A' && *pin <= 'F'){
			c = 10 + (*pin - 'A');
		} else if(*pin >= 'a' && *pin <= 'f'){
			c = 10 + (*pin - 'a');
		} else {  //Invalid hex string
			return -1;
		}
		if(odd){
			*pout = (*pout << 4) | c;
			pout++;
			size++;
		} else {
			*pout = c;
		}
		odd = !odd;
	}
	if( odd )
		return -1;
	if( size < n )
		*pout = 0;
	return size;
}

int cmpGoSlice_GoSlice(GoSlice *slice1, GoSlice_ *slice2){
	return ((slice1->len == slice2->len)) && (memcmp(slice1->data,slice2->data, sizeof(GoSlice_))==0 );
}

void bin2hex(unsigned char* buf, char *str, int n){
	unsigned char * pin = buf;
	const char * hex = "0123456789ABCDEF";
	char * pout = str;
	for(; n; --n){
		*pout++ = hex[(*pin>>4)&0xF];
		*pout++ = hex[(*pin++)&0xF];
	}
	*pout = 0;
}

int string_has_suffix(const char* str, const char* suffix){
	int string_len = strlen(str);
	int suffix_len = strlen(suffix);
	if(string_len >= suffix_len){
		char* p = (char*)str + (string_len - suffix_len);
		return strcmp(p, suffix) == 0;
	}
	return 0;
}

int string_has_prefix(const char* str, const char* prefix){
	int string_len = strlen(str);
	int prefix_len = strlen(prefix);
	if(string_len >= prefix_len){
		return strncmp(str, prefix, prefix_len) == 0;
	}
	return 0;
}

extern int count_words(const char* str, int length){
	int words = 1;
	char prevChar = 0;
	for(int i = 0; i < length; i++){
		char c = str[i];
		if( c == ' ' && prevChar != ' ' ) words++;
		prevChar = c;
	}
	return words;
}
