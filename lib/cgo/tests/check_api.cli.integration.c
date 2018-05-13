#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define BUFFER_SIZE 1024


Test(api_cli_integration, TestStableShowConfig) {
	char output[1024];
	Config__Handle configHandle;
	App__Handle appHandle;
	const char* str = "boxfort-worker showConfig";
	GoString showConfigCommand = {str, strlen(str) };
	GoUint32 errcode;
	
	errcode = SKY_cli_LoadConfig( &configHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_NewApp( &configHandle, &appHandle );
	cr_assert(errcode == SKY_OK, "SKY_cli_NewApp failed");
	redirectStdOut();
	errcode = SKY_cli_App_Run( &appHandle, showConfigCommand );
	getStdOut(output, BUFFER_SIZE);
	cr_assert(errcode == SKY_OK, "SKY_cli_App_Run failed");
}