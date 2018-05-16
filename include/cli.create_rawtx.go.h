/**
 * Structure used to specify amounts transferred in a transaction.
 */
typedef struct {
	GoString_ Addr; ///< Sender / receipient address.
	GoInt64_ Coins; ///< Amount transferred (e.g. measured in SKY)
} cli__SendAmount;
