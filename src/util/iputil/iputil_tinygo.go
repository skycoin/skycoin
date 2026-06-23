//go:build tinygo

package iputil

// LocalhostIP returns the address for localhost on the machine.
//
// TinyGo's net.Interface does not implement the Addrs method used to
// enumerate interface addresses, so fall back to the standard loopback
// address.
func LocalhostIP() (string, error) {
	return "127.0.0.1", nil
}
