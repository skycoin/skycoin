module github.com/skycoin/skycoin

go 1.25.1

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/bitfield/script v0.24.1
	github.com/blang/semver v3.5.1+incompatible
	github.com/briandowns/spinner v1.23.2
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/gin-gonic/gin v1.11.0
	github.com/google/go-cmp v0.7.0
	github.com/gopherjs/gopherjs v1.17.2
	github.com/ivanpirog/coloredcobra v1.0.1
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/rs/cors v1.11.1
	github.com/shopspring/decimal v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/skycoin/encodertest v0.0.0-20190217072920-14c2e31898b9
	github.com/skycoin/hardware-wallet-daemon v0.1.1-0.20251214205958-0632d0772c3b
	github.com/skycoin/hardware-wallet-go v1.1.1-0.20251214205422-17746f44286f
	github.com/skycoin/skycoin-lite v0.0.0-20190712083345-f5a3f17e6603
	github.com/skycoin/skywire v1.3.32-0.20251008181048-e232456f8799
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.11.1
	github.com/toqueteos/webbrowser v1.2.1
	go.etcd.io/bbolt v1.4.3
	golang.org/x/term v0.38.0
)

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.2 // indirect
	github.com/bytedance/sonic/loader v0.4.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.29.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gousb v1.1.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itchyny/gojq v0.12.18 // indirect
	github.com/itchyny/timefmt-go v0.1.7 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.57.1 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/skycoin/hardware-wallet-protob v0.0.0-20250805154629-410561e1bc2f // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/arch v0.23.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/sh/v3 v3.12.0 // indirect
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
