
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
