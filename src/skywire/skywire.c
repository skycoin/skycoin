

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

#define ROUTE_MAX 256
#define CONNECTION_MAX 256

struct Router {
	*uint32_t table; //routing table
	int table_idx;

	*Connection pool; //table of connections
	int pool_idx;


}

struct Route {
	uint32 src_id //packet comes from
	uint32 dst_id //packet goes to

}

void NewRouter() *Router {

	struct *Router R;
	R = (*Router) malloc(sizeof(Router));

	R.table = (*Route ) malloc(ROUTE_MAX*sizeof(*Route));
	R.table_idx = 0;

	R.pool = (*Connection) malloc(ROUTE_MAX*sizeof(Connection));
	R.pool.idx = 0;
	
}


Router *gRouter; //global router

//function pointer definition
//connectionid, data to write, amount of data to write
void (*ConnectionWriteFunc) (int, *char, int);

struct Connection {
	int id;
	ConnectionWriteFunc write_func;

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



void AddConnection(int connectionId, ConnectionWriteFunc writeFunct) {

	int idx = gRouter.connectionsIdx
}


/*
	Main Test
*/
int main() {

	printf("test \n");

	return 0;
}
