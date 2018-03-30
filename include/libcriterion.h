


// // TODO: Write like this cr_assert(eq(type(Address), addr1, addr2))
int cr_user_Address_eq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return 0;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return 0;
  }
  return 1;
}

char *cr_user_Address_tostr(Address *addr1)
{
  char *out;

  cr_asprintf(&out, "(Address) { .Key = %s, .Version = %llu }", addr1->Key, (unsigned long long) addr1->Version);
  return out;
}
// // TODO: Write like this cr_assert(not(eq(type(Address), addr1, addr2)))
int cr_user_Address_noteq(Address *addr1, Address *addr2){
  if(addr1->Version != addr2->Version)
    return SKY_OK;
  for (int i = 0; i < sizeof(Ripemd160); ++i) {
    if(addr1->Key[i] != addr2->Key[i])
      return SKY_OK;
  }
  return SKY_ERROR;
}

int cr_user_GoString_eq(GoString *string1, GoString *string2){

  if(  strcmp(string1->p,string2->p) != 0 )
  {
    return SKY_ERROR;
  } else {
    return SKY_OK;
  }
}

char *cr_user_GoString_tostr(GoString *string)
{
  char *out;
  cr_asprintf(&out, "(GoString) { .Data = %s, .Length = %llu }", string->p, (unsigned long long) string->n);
  return out;
}

int cr_user_GoString__eq(GoString_ *string1, GoString_ *string2){
  return cr_user_GoString_eq((GoString *)string1, (GoString *)string2);
}

char *cr_user_GoString__tostr(GoString_ *string) {
  return cr_user_GoString_tostr((GoString *)string);
}