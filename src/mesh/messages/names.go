package messages

import (
	"regexp"
)

const seg = "[a-z\\d_]*[a-z\\d]"
const dot = "\\."

var HostRX = regexp.MustCompile("^" + seg + "$")
var DomainRX = regexp.MustCompile("^" + seg + dot + seg + "$")
