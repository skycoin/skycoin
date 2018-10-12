package main

import (
	"reflect"

	cert "github.com/skycoin/skycoin/src/util/certutil"
)

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_certutil_NewTLSCertPair
func SKY_certutil_NewTLSCertPair(organization string, validUntil string, extraHosts []string, _cert *C.GoSlice_, _key *C.GoSlice_) (____error_code uint32) {
	____time_validUntil, ____return_err := parseTimeValue(validUntil)
	if ____return_err == nil {
		cert, key, ____return_err := cert.NewTLSCertPair(organization, ____time_validUntil, extraHosts)
		if ____return_err == nil {
			copyToGoSlice(reflect.ValueOf(cert), _cert)
			copyToGoSlice(reflect.ValueOf(key), _key)
		}
	}
	____error_code = libErrorCode(____return_err)
	return
}
