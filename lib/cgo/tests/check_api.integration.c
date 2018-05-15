#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define NODE_ADDRESS "SKYCOIN_NODE_HOST"
#define NODE_ADDRESS_DEFAULT "http://127.0.0.1:6420"
#define BUFFER_SIZE 1024

char* getNodeAddress(){
	GoString_ nodeAddress;
	memset(&nodeAddress, 0, sizeof(GoString_));
	GoString  nodeEnvName = {NODE_ADDRESS, strlen(NODE_ADDRESS)};
	int result = SKY_cli_Getenv(nodeEnvName, &nodeAddress);
	cr_assert(result == SKY_OK, "Couldn\'t get node address from enviroment");
	registerMemCleanup((void*)nodeAddress.p);
	if( strcmp(nodeAddress.p, "") == 0){
		return NODE_ADDRESS_DEFAULT;
	}
	return (char*)nodeAddress.p;
}

Test(api_integration, TestVersion) {
	GoString_ version;
	memset(&version, 0, sizeof(GoString_));
	
	char* pNodeAddress = getNodeAddress();
	GoString nodeAddress = {pNodeAddress, strlen(pNodeAddress)};
	Client__Handle clientHandle;
	Handle versionDataHandle;
	
	int result = SKY_api_NewClient(nodeAddress, &clientHandle);
	cr_assert(result == SKY_OK, "Couldn\'t create client");
	registerHandleClose( clientHandle );
	result = SKY_api_Client_Version( &clientHandle, &versionDataHandle );
	cr_assert(result == SKY_OK, "Couldn\'t get version");
	registerHandleClose( versionDataHandle );
	result = SKY_JsonEncode_Handle(versionDataHandle, &version);
	cr_assert(result == SKY_OK, "Couldn\'t json encode version");
	registerMemCleanup((void*)version.p);
	int versionLength = strlen(version.p);
	cr_assert(versionLength > 0, "Invalid version data");
}