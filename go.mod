module github.com/skycoin/skycoin

go 1.26.1

require (
	github.com/0magnet/coloredcobra v1.0.2
	github.com/NYTimes/gziphandler v1.1.1
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/bitfield/script v0.25.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/briandowns/spinner v1.23.2
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/gdamore/tcell/v3 v3.4.2
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.7.0
	github.com/google/gousb v1.1.3
	github.com/google/uuid v1.6.0
	github.com/gopherjs/gopherjs v1.21.0
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/pelletier/go-toml/v2 v2.4.3
	github.com/rs/cors v1.11.1
	github.com/shopspring/decimal v1.4.0
	github.com/sirupsen/logrus v1.10.1
	github.com/skycoin/encodertest v0.0.0-20190217072920-14c2e31898b9
	github.com/skycoin/hardware-wallet-protob v0.0.0-20250805154629-410561e1bc2f
	github.com/skycoin/hardware-wallet/firmware v0.0.0-20260117000250-43c66b2bf1d2
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.12.1
	github.com/toqueteos/webbrowser v1.2.1
	go.etcd.io/bbolt v1.5.0
	golang.org/x/term v0.45.0
	golang.org/x/time v0.15.0
	modernc.org/sqlite v1.57.0
)

require (
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/gdamore/encoding v1.0.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itchyny/gojq v0.12.19 // indirect
	github.com/itchyny/timefmt-go v0.1.8 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.1 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	modernc.org/libc v1.75.5 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
	mvdan.cc/sh/v3 v3.13.1 // indirect
)

// IT IS FORBIDDEN TO USE REPLACE DIRECTIVES

// [error] The go.mod file for the module providing named packages contains one or
//	more replace directives. It must not contain directives that would cause
//	it to be interpreted differently than if it were the main module.

// Uncomment for tests with local sources

//replace github.com/skycoin/hardware-wallet-go => ../hardware-wallet-go
//replace github.com/skycoin/hardware-wallet-daemon => ../hardware-wallet-daemon
//replace github.com/skycoin/skycoin-lite => ../skycoin-lite
//replace github.com/skycoin/skywire => ../skywire
//replace github.com/skycoin/hardware-wallet/firmware => ../hardware-wallet/firmware

// Below should reflect current versions of the following deps
// To update deps to specific commit hash:
// 1) Uncomment one of the following lines and substituite version with desired commit hash:
//replace github.com/skycoin/hardware-wallet-daemon => github.com/skycoin/hardware-wallet-daemon v0.1.1-0.20251214205958-0632d0772c3b
//replace github.com/skycoin/hardware-wallet-go => github.com/skycoin/hardware-wallet-go v1.1.1-0.20251214205422-17746f44286f
//replace github.com/skycoin/skycoin-lite => github.com/skycoin/skycoin-lite v0.0.0-20190712083345-f5a3f17e6603
//replace github.com/skycoin/skywire => github.com/skycoin/skywire v1.3.32-0.20251008181048-e232456f8799
// 2) Run `go mod tidy && go mod vendor`
// 3) Copy the populated version string to the correct place in require(...) above - replacing the specified version string
// 4) Re-comment the uncommented replace directive above
// 5) Save this file.
// 6) Run `go mod tidy && go mod vendor`
