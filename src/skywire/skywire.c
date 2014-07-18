

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

#define ENFORCE_ROUTE_DST 0 //should source/dst on routes be enforced

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
};


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


//0 on success
int AddConnection(int connectionId, 
		ConnectionWriteFunc writeFunc,
		DisconnectCallback on_disconnect) 
{
	int idx = gRouter.connectionsIdx
	if(idx == CONNECTION_MAX) {
		return 1; //connection max
	}

	for(int i=0;i<gRouter.pool_idx;i++) {
		if gRouter.pool[idx].id = connectionid {
			return 2; //connection id already in use
		}
	}

	gRouter.pool[idx] = (*Connection) calloc(1,sizeof(Connection));
	gRouter.pool[idx].id = connectionid;
	gRouter.ConnectionWriteFunc = writeFunc;
	gRouter.DisconnectCallback = on_disconnect

	gRouter.idx++;
}

int void ConnectionData(int connectionId, char* data, int size) {
	//process data incoming
}

/*
	Main Test
*/
int main() {

	printf("test \n");

	return 0;
}
