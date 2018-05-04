
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

Test(api_cli, TestLoadConfig) {
	Config__Handle configHandle;
	GoUint32 errcode;
	GoString str;
	
	setenv("COIN", "foocoin", 1);
	errcode = SKY_cli_LoadConfig(&configHandle);
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_Config_GetCoin(&configHandle, &str);
	cr_assert(errcode == SKY_OK, "SKY_cli_Config_GetCoin failed");
	SKY_handle_close((Handle)configHandle);
	unsetenv("COIN");
}