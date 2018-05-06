package main

/*

  #include <string.h>
  #include <stdlib.h>

  #include "skytypes.h"
*/
import "C"

import (
	cli "github.com/skycoin/skycoin/src/api/cli"
)

//export SKY_cli_NewPasswordReader
func SKY_cli_NewPasswordReader(_p []byte) C.PasswordReader__Handle {
	pr := cli.NewPasswordReader(_p)
	return registerPasswordReaderHandle(pr)
}
