package main

import cert "github.com/skycoin/skycoin/src/util/cert"

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

//export SKY_cert_CreateCertIfNotExists
func SKY_cert_CreateCertIfNotExists(_host, _certFile, _keyFile string, _appName string) (____error_code uint32) {
	____error_code = 0
	defer func() {
		____error_code = catchApiPanic(____error_code, recover())
	}()
	host := _host
	certFile := _certFile
	keyFile := _keyFile
	appName := _appName
	____return_err := cert.CreateCertIfNotExists(host, certFile, keyFile, appName)
	____error_code = libErrorCode(____return_err)
	if ____return_err == nil {
	}
	return
}
