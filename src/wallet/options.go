package wallet

// Option represents the general options, it can be used to set optional
// parameters while creating a new wallet. Also, can be used to get
// entries service of a wallet.
type Option func(interface{})

// Bip44EntriesOptions represents the options that will be used
// by bip44 to get entries service
type Bip44EntriesOptions struct {
	Account uint32
	Change  bool
}

// Account is the option for specifying account when acquiring the entries service
func Account(index uint32) Option {
	return func(opts interface{}) {
		opts.(*Bip44EntriesOptions).Account = index
	}
}

// Change is the option for whether choosing the change chain when acquiring the entries service
func Change(change bool) Option {
	return func(opts interface{}) {
		opts.(*Bip44EntriesOptions).Change = change
	}
}
