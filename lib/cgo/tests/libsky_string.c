
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
    for(; *pin && n; --n){
        *pout++ = hex[(*pin>>4)&0xF];
        *pout++ = hex[(*pin++)&0xF];
    }
    *pout = 0;
}

void strhex(unsigned char* buf, char *str){
  strnhex(buf, str, SIZE_ALL);
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
	if( size < n )
		*pout = 0;
	return size;
}

int cmpGoSlice_GoSlice(GoSlice *slice1, GoSlice_ *slice2){

return (slice1->len == slice2->len) &&
  (strcmp( (unsigned char *) slice1->data, (unsigned char *) slice2->data) == 0);}

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
