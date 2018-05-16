#include "skynumber.h"

int number_eq(Number* n1, Number* n2){
	if ( n1->neg != n2->neg )
		return 0;
	if ( n1->nat.len != n2->nat.len )
		return 0;
	for( int i = 0; i < n1->nat.len; i++){
		if( ((char*)n1->nat.data)[i] != ((char*)n2->nat.data)[i])
			return 0;
	}
	return 1;
}