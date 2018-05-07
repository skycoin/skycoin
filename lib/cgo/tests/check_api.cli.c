
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define SKYCOIN_COIN_ENV_VAR "COIN"
#define SKYCOIN_COIN_FOO "foocoin"
#define SKYCOIN_RPC_SAMPLE "http://111.22.33.44:5555"
#define SKYCOIN_RPC_ENV_VAR "RPC_ADDR"
#define SKYCOIN_WALLET_DIR_ENV_VAR "WALLET_DIR"
#define SKYCOIN_WALLET_NAME_ENV_VAR "WALLET_NAME"
#define SKYCOIN_WALLET_DIR_SAMPLE "/home/foo/bar"
#define SKYCOIN_WALLET_NAME_SAMPLE "bar.wlt"
#define SKYCOIN_WALLET_FULL_PATH_SAMPLE "/home/foo/bar/bar.wlt"

Test(api_cli, TestLoadConfig) {
	Config__Handle configHandle;
	GoUint32 errcode;
	GoString strConfigValue;
	GoString strEnv;
	GoString strEnv2;
	GoString strEnvCoinVar = {
		SKYCOIN_COIN_ENV_VAR,
		4
	};
	
	GoString strEnvRPCVar = {
		SKYCOIN_RPC_ENV_VAR,
		8
	};
	
	GoString strEnvCoinFoo = {
		SKYCOIN_COIN_FOO,
		7
	};
	
	GoString strEnvRPCSample = {
		SKYCOIN_RPC_SAMPLE,
		24
	};
	
	GoString strWalletDirVar = {
		SKYCOIN_WALLET_DIR_ENV_VAR,
		10
	};
	
	GoString strWalletNameVar = {
		SKYCOIN_WALLET_NAME_ENV_VAR,
		11
	};
	
	GoString strWalletDirSample = {
		SKYCOIN_WALLET_DIR_SAMPLE,
		13
	};
	
	GoString strWalletNameSample = {
		SKYCOIN_WALLET_NAME_SAMPLE,
		7
	};
	
	errcode = SKY_cli_Getenv(strEnvCoinVar, &strEnv);
	cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed getting COIN");
	errcode = SKY_cli_Setenv(strEnvCoinVar, strEnvCoinFoo);
	cr_assert(errcode == SKY_OK, "SKY_cli_Setenv failed setting COIN");
	errcode = SKY_cli_LoadConfig(&configHandle);
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_Config_GetCoin(&configHandle, &strConfigValue);
	cr_assert(errcode == SKY_OK, "SKY_cli_Config_GetCoin failed");
	cr_assert(strcmp(strConfigValue.p, strEnvCoinFoo.p) == 0, "SKY_cli_LoadConfig with coin failed");
	SKY_cli_Setenv(strEnvCoinVar, strEnv); //Restore previous value
	registerMemCleanup(strConfigValue.p);
	registerMemCleanup(strEnv.p);
	SKY_handle_close((Handle)configHandle);
	
	errcode = SKY_cli_Getenv(strEnvRPCVar, &strEnv);
	cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed getting RPC_ADDR");
	errcode = SKY_cli_Setenv(strEnvRPCVar, strEnvRPCSample);
	cr_assert(errcode == SKY_OK, "SKY_cli_Setenv failed setting RPC_ADDR");
	errcode = SKY_cli_LoadConfig(&configHandle);
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_Config_GetRPCAddress(&configHandle, &strConfigValue);
	cr_assert(errcode == SKY_OK, "SKY_cli_Config_GetRPCAddress failed");
	cr_assert(strcmp(strConfigValue.p, strEnvRPCSample.p) == 0, "SKY_cli_LoadConfig with RPCAddress failed");
	SKY_cli_Setenv(strEnvCoinVar, strEnv); //Restore previous value
	registerMemCleanup(strConfigValue.p);
	registerMemCleanup(strEnv.p);
	SKY_handle_close((Handle)configHandle);
	
	//Testing Wallet Dir and Wallet Name
	
	errcode = SKY_cli_Getenv(strWalletDirVar, &strEnv);
	cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed getting WALLET_DIR");
	errcode = SKY_cli_Getenv(strWalletNameVar, &strEnv2);
	cr_assert(errcode == SKY_OK, "SKY_cli_Getenv failed getting WALLET_NAME");
	
	errcode = SKY_cli_Setenv(strWalletDirVar, strWalletDirSample);
	cr_assert(errcode == SKY_OK, "SKY_cli_Setenv failed setting WALLET_DIR");
	errcode = SKY_cli_Setenv(strWalletNameVar, strWalletNameSample);
	cr_assert(errcode == SKY_OK, "SKY_cli_Setenv failed setting WALLET_NAME");
	
	errcode = SKY_cli_LoadConfig(&configHandle);
	cr_assert(errcode == SKY_OK, "SKY_cli_LoadConfig failed");
	errcode = SKY_cli_Config_FullWalletPath(&configHandle, &strConfigValue);
	cr_assert(errcode == SKY_OK, "SKY_cli_Config_FullWalletPath failed");
	cr_assert(strcmp(strConfigValue.p, SKYCOIN_WALLET_FULL_PATH_SAMPLE) == 0, "SKY_cli_LoadConfig with Wallet Dir failed");
	SKY_cli_Setenv(strWalletDirVar, strEnv); //Restore previous value
	SKY_cli_Setenv(strWalletNameVar, strEnv2); //Restore previous value
	registerMemCleanup(strConfigValue.p);
	registerMemCleanup(strEnv.p);
	registerMemCleanup(strEnv2.p);
	SKY_handle_close((Handle)configHandle);
}