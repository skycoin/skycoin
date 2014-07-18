

#include "skywire.h"

#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <openssl/sha.h>

#include <stdio.h>

#include <time.h>
#include <sys/time.h>

#include <stdlib.h>     /* srand, rand */
#include <time.h>       /* time */
#include <malloc.h>

#include <stdint.h>

/*
g++ skywire.c -lssl -lcrypto
g++ skywire.c -lssl -lcrypto -lrt
g++ -O3 skywire.c -lssl -lcrypto -lrt

clang++ -O1 -g -fstanitize=address ./skywire.c -lssl -lcrypto -lrt
clang++ -O1 -g -fstanitize=address ./skywire.c -lssl -lcrypto -lrt

clang++ -O2 ./skywire.c -lssl -lcrypto -lrt
*/

/*

/*
	Steps:
	- connect (ECDH)
	- setup circuit
	- start forwarding

*/

typedef uint32_t uint32;

#define ROUTE_MAX 256
#define CONNECTION_MAX 256

struct Router {
	*uint32 table; //routing table
	int table_idx;

	*Connection pool; //table of connections
	int pool_idx;
}

#define ENFORCE_ROUTE_DST 0 //should source/dst connection on routes be enforced

struct Route {
	uint32 src_id; //packet comes from
	uint32 dst_id; //packet goes to
	
	uint32 src_con;
	uint32 dst_con;	
}

void NewRouter() *Router {

	struct *Router R;
	R = (*Router) calloc(1, sizeof(Router));

	R.table = (*Route ) calloc(1, ROUTE_MAX*sizeof(*Route));
	R.table_idx = 0;

	R.pool = (*Connection) calloc(1, ROUTE_MAX*sizeof(Connection));
	R.pool.idx = 0;

}

*Connection GetConnection(*Router router, int ConnectionId) {
	for(int i=0; i<router.pool.idx; i++) {
		if(router.pool[i].id == ConnectionId) {
			return router.pool[i];
		}
	}
	return NULL;
}

Router *gRouter; //global router

//function pointer definition
//connectionid, data to write, amount of data to write
void (*ConnectionWriteFunc) (int, *char, int);
void (*OnDisconnectCallback) (int, int);

struct Connection {
	int id;
	ConnectionWriteFunc write_func;
	DisconnectCallback on_disconnect;

	char[65] snd_pubkey; //secp256k1 public key
	char[20] rcv_seckey;

	char[20] snd_secret; //secret for sending
	char[20] rcv_secret; //secret for receiving

	int state; //0 for not inited
};



/*
	Crytography Stubs
*/

struct SHA256 {
	char[20];
};

struct SecKey {
	char[20]; //secp256k1
};

struct PubKey {
	char[65];
};

//writes 20 bytes to dst
void SumSHA256(char* dst, char* src, int len) {
	for(int i=0;i<20;i++) {dst[i] = 0;}
}

//generates 20 byte seckey
SecKey KeyGen() {
	struct SecKey seckey;
	return seckey;
}

PubKey PubFromSec(Pubkey) {
	struct PubKey pubkey;
	return pubkey;
}

/*
	Messages
*/

struct IntroMsg {
	char[65] pubkey //ephermeral public key 
	char[20] secret //ephemeral
}

/*
	Global State Setup
*/


//setup module for operation
void SETUP() {
	gRouter = NewRouter();
}

//shutdown module
void TEARDOWN() {
	free(gRouter);
}

int verifyIntroMessage(char* buf, int size) {

}
//0 on success
int AddConnection(int connectionId, 
		ConnectionWriteFunc writeFunc,
		DisconnectCallback on_disconnect,
		*char connectMsg, int size ) 
{
	int idx = gRouter.connectionsIdx
	if(idx == CONNECTION_MAX) {
		return 1; //connection max
	}
	//verify the connect message

	int ret = verifyIntroMessage(connectMsg, size);
	if ret != 0 {
		return 2; //intro message failed
	}

	//setup
	for(int i=0;i<gRouter.pool_idx;i++) {
		if gRouter.pool[idx].id = connectionid {
			return 3; //connection id already in use
		}
	}

	gRouter.pool[idx] = (*Connection) calloc(1,sizeof(Connection));
	gRouter.pool[idx].id = connectionid;
	gRouter.ConnectionWriteFunc = writeFunc;
	gRouter.DisconnectCallback = on_disconnect

	gRouter.idx++;


}

int ConnectionDataIn(int connectionId, char* data, int size) {
	//process data incoming
	struct *Connection con;
	con = GetConnection(gRouter, connectionid);
	if(con == NULL) {
		return 1; //connectionId does not exist
	}



}

/*
	Main Test
*/
int main() {

	printf("test \n");

	return 0;
}
