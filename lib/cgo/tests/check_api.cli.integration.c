#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

Test(api_cli_integration, TestGenerateAddresses) {
	Config__Handle configHandle;
	App__Handle appHandle;
	GoString showConfigCommand = {"showConfig", 10 };
	GoUint32 errcode;
	
	errcode = SKY_cli_LoadConfig( &configHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_NewApp( &configHandle, &appHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
	errcode = SKY_cli_App_Run( &appHandle, showConfigCommand );
	cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");
}