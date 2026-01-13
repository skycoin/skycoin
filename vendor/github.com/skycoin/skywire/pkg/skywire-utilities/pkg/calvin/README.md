# calvin
convert text to Calvin S ascii font (https://patorjk.com/software/taag/#p=display&f=Calvin%20S&t=Type%20Something%20)


example:

```
$ echo "Hello, World!" | go run cmd/calvin/calvin.go
╦ ╦ ┌─┐┬  ┬  ┌─┐   ╦ ╦ ┌─┐┬─┐┬  ┌┬┐┬    
╠═╣ ├┤ │  │  │ │   ║║║ │ │├┬┘│   │││    
╩ ╩ └─┘┴─┘┴─┘└─┘┘  ╚╩╝ └─┘┴└─┴─┘─┴┘o    

```

library usage example

```
package main

import (
	"github.com/skycoin/skywire/pkg/skywire-utilities/calvin"
)

func main() {
	println(calvin.AsciiFont("Hello, World!"))
}

```
