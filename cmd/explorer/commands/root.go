// Package commands implements the skycoin explorer.
package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/spf13/cobra"

	"github.com/skycoin/skycoin/explorer"
)

const (
	defaultUIFilesFolder = "./dist/"
	defaultExplorerHost  = "127.0.0.1:8001"
	defaultSkycoinAddr   = "http://127.0.0.1:6420"
	// timeout for requests to the backend skycoin node
	skycoinRequestTimeout = time.Second * 30
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	serverReadTimeout  = time.Second * 10
	serverWriteTimeout = time.Second * 60
	serverIdleTimeout  = time.Second * 120
)

var (
	explorerHost  string   // URL of the web server. Override with envvar EXPLORER_HOST. Must not have scheme
	skycoinAddr   *url.URL // URL of the Skycoin node. Override with envvar SKYCOIN_ADDR. Must have scheme, e.g. http://
	uiFilesFolder string   // Path for the folder with the precompiled front-end files
	apiOnly       bool
	verify        bool
)

func init() {
	log.SetOutput(os.Stdout)

	// Defaults from environment
	explorerHost = os.Getenv("EXPLORER_HOST")
	if explorerHost == "" {
		explorerHost = defaultExplorerHost
	}

	skycoinAddrString := os.Getenv("SKYCOIN_ADDR")
	if skycoinAddrString == "" {
		skycoinAddrString = defaultSkycoinAddr
	}
	skycoinAddr = buildNodeURL(skycoinAddrString)

	RootCmd.Flags().StringVarP(&explorerHost, "server-host", "s", explorerHost, "The addr:port to bind the explorer web server to\033[0m\n\r")
	RootCmd.Flags().StringVarP(&skycoinAddrString, "node-addr", "n", skycoinAddr.String(), "The skycoin node's address\033[0m\n\r")
	RootCmd.Flags().StringVarP(&uiFilesFolder, "files-folder", "f", "", "Path for the folder with the precompiled front-end files")
	RootCmd.Flags().BoolVarP(&apiOnly, "api-only", "a", false, "Only run the API, don't serve static content")
	RootCmd.Flags().BoolVarP(&verify, "verify", "v", false, "Run init() checks and quit")

	parsedJSONAPIEndpoints = make([]ParsedJSONAPIEndpoint, len(apiEndpoints)+1)

	parsedJSONAPIEndpoints[len(apiEndpoints)] = ParsedJSONAPIEndpoint{
		APIEndpoint: docEndpoint,
	}

	for i := range apiEndpoints {
		parsedJSONAPIEndpoints[i].APIEndpoint = apiEndpoints[i]
		resp := []byte(apiEndpoints[i].ExampleResponse)
		if err := json.Unmarshal(resp, &parsedJSONAPIEndpoints[i].ParsedExampleResponse); err != nil {
			log.Println("Error parsing example response JSON:", err)
			log.Println("path:", apiEndpoints[i].ExplorerPath)
			log.Println("Example response:")
			log.Println(apiEndpoints[i].ExampleResponse)
			log.Panic(err)
		}
		if apiEndpoints[i].ExampleVerboseResponse != "" {
			resp = []byte(apiEndpoints[i].ExampleVerboseResponse)
			if err := json.Unmarshal(resp, &parsedJSONAPIEndpoints[i].ParsedExampleVerboseResponse); err != nil {
				log.Println("Error parsing example verbose response JSON:", err)
				log.Println("path:", apiEndpoints[i].ExplorerPath)
				log.Println("Example verbose response:")
				log.Println(apiEndpoints[i].ExampleVerboseResponse)
				log.Panic(err)
			}
		}
	}
	t := template.Must(template.New("docs").Parse(docTemplate))

	endpoints := []APIEndpoint{}
	endpoints = append(endpoints, apiEndpoints...)
	endpoints = append(endpoints, docEndpoint)

	b := &bytes.Buffer{}
	if err := t.Execute(b, endpoints); err != nil {
		log.Panic(err)
	}

	docTemplateBody = b.String()
}

// RootCmd is skycoin blockchain explorer
var RootCmd = &cobra.Command{
	Use: func() string {
		return strings.Split(filepath.Base(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf("%v", os.Args), "[", ""), "]", "")), " ")[0]
	}(),
	Short:                 "skycoin blockchain explorer",
	Long:                  calvin.AsciiFont("explorer"),
	SilenceUsage:          true,
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Run: func(_ *cobra.Command, _ []string) {
		// Validate and adjust values
		if uiFilesFolder != "" && uiFilesFolder[len(uiFilesFolder)-1] != '/' {
			uiFilesFolder += "/"
		}

		if verify {
			log.Println("Verified")
			return
		}

		log.Printf("Starting explorer at %s, proxying to %s\n", explorerHost, skycoinAddr.String())
		log.Println("api.html output:")
		log.Println(docTemplateBody)

		mux := http.NewServeMux()

		gzipHandle := func(path string, handler http.Handler) {
			mux.Handle(path, gziphandler.GzipHandler(handler))
		}

		if !apiOnly {
			// Needed for browsers to load JS correctly
			if err := mime.AddExtensionType(".js", "application/javascript"); err != nil {
				log.Fatalln("Unable to register js mime type.")
			}

			var fileServer http.Handler
			if uiFilesFolder != "" {
				// External GUI directory
				log.Printf("Serving GUI from local folder: %s", uiFilesFolder)
				fileServer = http.FileServer(http.Dir(uiFilesFolder))
			} else {
				// Embedded GUI
				log.Println("Serving embedded GUI")
				sub, err := fs.Sub(explorer.EmbeddedFiles, "dist")
				if err == nil {
					fileServer = http.FileServer(http.FS(sub))
				} else {
					// fallback if dist/ not present in FS
					fileServer = http.FileServer(http.FS(explorer.EmbeddedFiles))
				}
			}

			// Serve static files
			gzipHandle("/", fileServer)

			// Angular SPA fallback (/app/* → index.html)
			mux.Handle("/app/", gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if uiFilesFolder != "" {
					http.ServeFile(w, r, filepath.Join(uiFilesFolder, "index.html"))
					return
				}
				f, err := explorer.EmbeddedFiles.Open("dist/index.html")
				if err != nil {
					http.Error(w, "index.html not found", http.StatusInternalServerError)
					return
				}
				defer f.Close() //nolint:errcheck
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				if _, err := io.Copy(w, f); err != nil {
					log.Printf("Error serving index.html: %v", err)
				}
			})))

			// Backwards compatibility redirects
			gzipHandle("/blocks", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			}))
			redirectToApp := func(basePath string) {
				gzipHandle(basePath, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					block := r.URL.Path[len(basePath):]
					path := fmt.Sprintf("/app%s%s", basePath, block)
					http.Redirect(w, r, path, http.StatusMovedPermanently)
				}))
			}
			redirectToApp("/block/")
			redirectToApp("/transaction/")
			redirectToApp("/address/")

			gzipHandle("/api.html", http.HandlerFunc(htmlDocs))
		} else {
			log.Println("Running in api-only mode (no GUI served)")
		}

		// Register API endpoints
		for _, e := range apiEndpoints {
			gzipHandle(e.ExplorerPath, e)
			log.Printf("%s proxied to %s with args %v", e.ExplorerPath, e.SkycoinPath, e.QueryArgs)
		}
		gzipHandle("/api/docs", http.HandlerFunc(jsonDocs))

		log.Printf("Running skycoin explorer on http://%s", explorerHost)

		s := &http.Server{
			Addr:         explorerHost,
			Handler:      mux,
			ReadTimeout:  serverReadTimeout,
			WriteTimeout: serverWriteTimeout,
			IdleTimeout:  serverIdleTimeout,
		}

		if err := s.ListenAndServe(); err != nil {
			log.Println("Fatal:", err)
		}
	},
}

// buildNodeURL converts a string to a *url.URL instance.
func buildNodeURL(skycoinAddrString string) *url.URL {
	origURL, err := url.Parse(skycoinAddrString)
	if err != nil {
		log.Println("SKYCOIN_ADDR must have a scheme, e.g. http://")
		log.Fatalln("Invalid SKYCOIN_ADDR", skycoinAddrString, err)
	}

	if origURL.Scheme == "" {
		log.Fatalln("SKYCOIN_ADDR must have a scheme, e.g. http://")
	}

	return &url.URL{
		Scheme: origURL.Scheme,
		Host:   origURL.Host,
	}
}

func buildSkycoinURL(path string, query url.Values) string {
	rawQuery := ""
	if query != nil {
		rawQuery = query.Encode()
	}

	u := &url.URL{
		Scheme:   skycoinAddr.Scheme,
		Host:     skycoinAddr.Host,
		Path:     path,
		RawQuery: rawQuery,
	}

	return u.String()
}

// CoinSupply coin supply struct
type CoinSupply struct {
	// Coins distributed beyond the project:
	CurrentSupply string `json:"current_supply"`
	// TotalSupply is CurrentSupply plus coins held by the distribution addresses that are spendable
	TotalSupply string `json:"total_supply"`
	// MaxSupply is the maximum number of coins to be distributed ever
	MaxSupply string `json:"max_supply"`
	// CurrentCoinHourSupply is coins hours in non distribution addresses
	CurrentCoinHourSupply string `json:"current_coinhour_supply"`
	// TotalCoinHourSupply is coin hours in all addresses including unlocked distribution addresses
	TotalCoinHourSupply string `json:"total_coinhour_supply"`
	// Distribution addresses which count towards total supply
	UnlockedAddresses []string `json:"unlocked_distribution_addresses"`
	// Distribution addresses which are locked and do not count towards total supply
	LockedAddresses []string `json:"locked_distribution_addresses"`
}

// APIEndpoint api endpoint struct
type APIEndpoint struct {
	ExplorerPath   string   `json:"explorer_path"`
	SkycoinPath    string   `json:"skycoin_path"`
	QueryArgs      []string `json:"query_args,omitempty"`
	Description    string   `json:"description"`
	ExampleRequest string   `json:"example_request"`
	// This strings will be parsed into a map[string]interface{} in order to render newlines
	ExampleResponse        string `json:"-"`
	ExampleVerboseResponse string `json:"-"`
}

func (s APIEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var query url.Values
	if s.QueryArgs != nil {
		query = url.Values{}
		for _, s := range s.QueryArgs {
			query.Add(s, r.URL.Query().Get(s))
		}
	}

	skycoinURL := buildSkycoinURL(s.SkycoinPath, query)

	log.Printf("Proxying request %s to skycoin node %s with timeout %v", r.URL.String(), skycoinURL, skycoinRequestTimeout)

	c := &http.Client{
		Timeout: skycoinRequestTimeout,
	}

	resp, err := c.Get(skycoinURL)
	if err != nil {
		msg := "Request to skycoin node failed"
		log.Println("ERROR:", msg, skycoinURL)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing body: ", err)
		}
	}()

	w.WriteHeader(resp.StatusCode)

	if s.ExplorerPath == "/api/coinmarketcap" {
		w.Header().Set("Content-Type", "text/plain")
		var cs CoinSupply
		if err := json.NewDecoder(resp.Body).Decode(&cs); err != nil {
			msg := "Decode CoinSupply result failed"
			log.Println("ERROR:", msg, skycoinURL, err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintf(w, "%s", cs.CurrentSupply); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if n, err := io.Copy(w, resp.Body); err != nil {
		msg := "Copying response from skycoin node to client failed"
		if n != 0 {
			msg += fmt.Sprintf(", after %d bytes were written", n)
		}
		log.Println("ERROR:", msg, skycoinURL)

		// An error response can only be written if the ResponseWriter has not been written to
		if n == 0 {
			http.Error(w, msg, http.StatusInternalServerError)
		}

		return
	}
}

var apiEndpoints = []APIEndpoint{
	{
		ExplorerPath:   "/api/coinSupply",
		SkycoinPath:    "/api/v1/coinSupply",
		Description:    "Returns metadata about the coin distribution.",
		ExampleRequest: "/api/coinSupply",
		ExampleResponse: `{
    "current_supply": "7187500.000000",
    "total_supply": "25000000.000000",
    "max_supply": "100000000.000000",
    "current_coinhour_supply": "23499025077",
    "total_coinhour_supply": "93679828577",
    "unlocked_distribution_addresses": [
        "R6aHqKWSQfvpdo2fGSrq4F1RYXkBWR9HHJ",
        "2EYM4WFHe4Dgz6kjAdUkM6Etep7ruz2ia6h",
        "25aGyzypSA3T9K6rgPUv1ouR13efNPtWP5m"
    ],
    "locked_distribution_addresses": [
        "gpqsFSuMCZmsjPc6Rtgy1FmLx424tH86My",
        "2EUF3GPEUmfocnUc1w6YPtqXVCy3UZA4rAq",
        "TtAaxB3qGz5zEAhhiGkBY9VPV7cekhvRYS"
    ]
}`,
	},
	{
		ExplorerPath:    "/api/coinmarketcap",
		SkycoinPath:     "/api/v1/coinSupply",
		Description:     "Returns circulating supply coin number.",
		ExampleRequest:  "/api/coinmarketcap",
		ExampleResponse: "7187500.000000",
	},

	{
		ExplorerPath:   "/api/blockchain/metadata",
		SkycoinPath:    "/api/v1/blockchain/metadata",
		Description:    "Returns blockchain metadata.",
		ExampleRequest: "/api/blockchain/metadata",
		ExampleResponse: `{
    "head": {
        "seq": 1893,
        "block_hash": "e20d5832b3f9bea4da58e149e4805b4e4a962ea7c5ce3cd9f31c6d7fc72e3300",
        "previous_block_hash": "e99fe31adf3a15aab77eb81d7aaa5477a96d2ce7f74ccb3d4cdce4da2eb01cb2",
        "timestamp": 1499405825,
        "fee": 5033788,
        "version": 0,
        "tx_body_hash": "c297eb14a9e68ec5501aa886e5bb720a58fe6466be633a8264f61eee9580a2c3"
    },
    "unspents": 783,
    "unconfirmed": 0
}`,
	},

	{
		ExplorerPath:   "/api/block",
		SkycoinPath:    "/api/v1/block",
		QueryArgs:      []string{"hash", "seq", "verbose"},
		Description:    "Returns information about a block, given a hash or sequence number. Assign 1 to the \"verbose\" argument to get more data in the response.",
		ExampleRequest: "/api/block?hash=e20d5832b3f9bea4da58e149e4805b4e4a962ea7c5ce3cd9f31c6d7fc72e3300",
		ExampleResponse: `{
    "header": {
        "seq": 1893,
        "block_hash": "e20d5832b3f9bea4da58e149e4805b4e4a962ea7c5ce3cd9f31c6d7fc72e3300",
        "previous_block_hash": "e99fe31adf3a15aab77eb81d7aaa5477a96d2ce7f74ccb3d4cdce4da2eb01cb2",
        "timestamp": 1499405825,
        "fee": 5033788,
        "version": 0,
        "tx_body_hash": "c297eb14a9e68ec5501aa886e5bb720a58fe6466be633a8264f61eee9580a2c3"
    },
    "body": {
        "txns": [
            {
                "length": 414,
                "type": 0,
                "txid": "c297eb14a9e68ec5501aa886e5bb720a58fe6466be633a8264f61eee9580a2c3",
                "inner_hash": "5fcc1649794894f2c79411a832f799ba12e0528ff530d7068abaa03c10e451cf",
                "sigs": [
                    "b951fd0c1528df88c87eb90cb1ecbc3ba2b6332ace16c2f1cc731976c0cebfb10ecb1c20335374f8cf0b832364a523c3e16f32c3240ed4eccfac1803caf8815100",
                    "6a06e57d130f6e780eecbe0c2626eed7724a2443438258598af287bf8fc1b87f041d04a7082550bd2055a08f0849419200fdac27c018d5cebf84e8bba1c4f61201",
                    "8cc9ae6ae5be81456fc8e2d54b6868c45b415556c689c8c4329cd30b671a18254885a1aef8a3ac9ab76a8dc83e08607516fdef291a003935ae4507775ae53c7800"
                ],
                "inputs": [
                    "0922a7b41d1b76b6b56bfad0d5d0f6845517bbd895c660eab0ebe3899b5f63c4",
                    "d73cf1f1d04a1d493fe3480a00e48187f9201bb64828fe0c638f17c0c88bb3d9",
                    "16dd81af869743599fe60108c22d7ee1fcbf1a7f460fffd3a015fbb3f721c36d"
                ],
                "outputs": [
                    {
                        "uxid": "8a941208d3f2d2c4a32438e05645fb64dba3b4b7d83c48d52f51bc1eb9a4117a",
                        "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                        "coins": "2361.000000",
                        "hours": 1006716
                    },
                    {
                        "uxid": "a70d1f0f488066a327acd0d5ea77b87d62b3b061d3db8361c90194a6520ab29f",
                        "dst": "SeDoYN6SNaTiAZFHwArnFwQmcyz7ZvJm17",
                        "coins": "51.000000",
                        "hours": 2013433
                    }
                ]
            }
        ]
    },
    "size": 414
}`,
		ExampleVerboseResponse: `{
    "header": {
        "seq": 1893,
        "block_hash": "e20d5832b3f9bea4da58e149e4805b4e4a962ea7c5ce3cd9f31c6d7fc72e3300",
        "previous_block_hash": "e99fe31adf3a15aab77eb81d7aaa5477a96d2ce7f74ccb3d4cdce4da2eb01cb2",
        "timestamp": 1499405825,
        "fee": 5033788,
        "version": 0,
        "tx_body_hash": "c297eb14a9e68ec5501aa886e5bb720a58fe6466be633a8264f61eee9580a2c3"
    },
    "body": {
        "txns": [
            {
                "length": 414,
                "type": 0,
                "txid": "c297eb14a9e68ec5501aa886e5bb720a58fe6466be633a8264f61eee9580a2c3",
                "inner_hash": "5fcc1649794894f2c79411a832f799ba12e0528ff530d7068abaa03c10e451cf",
                "fee": 5033788,
                "sigs": [
                    "b951fd0c1528df88c87eb90cb1ecbc3ba2b6332ace16c2f1cc731976c0cebfb10ecb1c20335374f8cf0b832364a523c3e16f32c3240ed4eccfac1803caf8815100",
                    "6a06e57d130f6e780eecbe0c2626eed7724a2443438258598af287bf8fc1b87f041d04a7082550bd2055a08f0849419200fdac27c018d5cebf84e8bba1c4f61201",
                    "8cc9ae6ae5be81456fc8e2d54b6868c45b415556c689c8c4329cd30b671a18254885a1aef8a3ac9ab76a8dc83e08607516fdef291a003935ae4507775ae53c7800"
                ],
                "inputs": [
                    {
                        "uxid": "0922a7b41d1b76b6b56bfad0d5d0f6845517bbd895c660eab0ebe3899b5f63c4",
                        "owner": "2Q2VViWhgBzz6c8GkXQyDVFdQUWcBcDon4L",
                        "coins": "7.000000",
                        "hours": 851106,
                        "calculated_hours": 851311
                    },
                    {
                        "uxid": "d73cf1f1d04a1d493fe3480a00e48187f9201bb64828fe0c638f17c0c88bb3d9",
                        "owner": "YPhukwVyLsPGX1FAPQa2ktr5XnSLqyGbr5",
                        "coins": "5.000000",
                        "hours": 6402335,
                        "calculated_hours": 6402335
                    },
                    {
                        "uxid": "16dd81af869743599fe60108c22d7ee1fcbf1a7f460fffd3a015fbb3f721c36d",
                        "owner": "YPhukwVyLsPGX1FAPQa2ktr5XnSLqyGbr5",
                        "coins": "2400.000000",
                        "hours": 800291,
                        "calculated_hours": 800291
                    }
                ],
                "outputs": [
                    {
                        "uxid": "8a941208d3f2d2c4a32438e05645fb64dba3b4b7d83c48d52f51bc1eb9a4117a",
                        "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                        "coins": "2361.000000",
                        "hours": 1006716
                    },
                    {
                        "uxid": "a70d1f0f488066a327acd0d5ea77b87d62b3b061d3db8361c90194a6520ab29f",
                        "dst": "SeDoYN6SNaTiAZFHwArnFwQmcyz7ZvJm17",
                        "coins": "51.000000",
                        "hours": 2013433
                    }
                ]
            }
        ]
    },
    "size": 414
}`,
	},

	{
		ExplorerPath:   "/api/blocks",
		SkycoinPath:    "/api/v1/blocks",
		QueryArgs:      []string{"start", "end", "verbose"},
		Description:    "Returns information about a range of blocks, given a start and end block sequence number. The range of blocks will include both the start and end sequence numbers. Assign 1 to the \"verbose\" argument to get more data in the response.",
		ExampleRequest: "https://explorer.skycoin.com/api/blocks?start=1891&end=1892",
		ExampleResponse: `{
    "blocks": [
        {
            "header": {
                "seq": 1891,
                "block_hash": "015265e82b4a28eee364327bf0d3b51fed54c2067318029974c262478a09bf17",
                "previous_block_hash": "71bf5abc05a4344597cf7df25249203df12e73b0a2e37dab386eb292a6cc5faa",
                "timestamp": 1499402935,
                "fee": 38414015,
                "version": 0,
                "tx_body_hash": "ac9759e9ffc3450d9af92e205f909375131df64e8f158a26e5e1e15334f45c8d"
            },
            "body": {
                "txns": [
                    {
                        "length": 802,
                        "type": 0,
                        "txid": "ac9759e9ffc3450d9af92e205f909375131df64e8f158a26e5e1e15334f45c8d",
                        "inner_hash": "cd36c8a00ff1fa81d94747b2d64107657fa40fa43bb0e5afc70d16661eca4286",
                        "sigs": [
                            "d34f0a9f0b82028e44702c7fc694e35aeb2777630efaaae32b9854b6c59c9ad747bdb91111864a050f024efd325a744cbbb7c1dde55a776c8f0d4856ac07fa5801",
                            "aa2f3ec98866bafda37f6e0b52be6fd1201131ce6b7edf1988b762edf14c35880b26a3a3fe9369a0f9eedecbe958abae83e54d572440903f5ff5b11cb9eb5e9200",
                            "746e171c22c99939d2d7e4b615f98de74f1c6ae152cb3d1e17e4d5e25d6f52b444242080a5a3d8f36de3befb8ec915321a7cb282eafc661c0857de753a6b40e100",
                            "e4f3db5afa209cbceac814d0f5c6fff8b57a3cc5c15063cd0b390067705de0a87313488d5eabd369855a87879a8d163fce08faeacd9d6538528d29c7cdbdaedf00",
                            "038263cdaf43a544b4ef66b2c910b9dd807ecd55e2fbf513c776b7accb1ff1fb5f591fbba22af58cf399d743955ab6cc7e3e641ecc675253306b2d0b6570672f00",
                            "215d1527af817f9c0613cacc3c2d07873ff563ab113023fc83ea836e3ac827f240df8a3d390d8e867cdda149d0db8f86bf435fa354dc4f3ac8c3a9b39297d7ba00",
                            "a2d7a50176ac1ca9417ca35f9848621bc3e0aeacb0cce713e9f93d82db4988ac5bca5bf02610c64f1580b198aadf49b04a14542466f7dca49d73b6c70f2577a201"
                        ],
                        "inputs": [
                            "e72d8ba4ce2d3b37aeb71df2e3bed80ee07204b3fa633f56cbce7bca836bd39c",
                            "6a349ba12c5d2827de6c24773d3dd8f6572e86adba4c8954a6d6e68df9e165e2",
                            "bb9a579003de101bee76d83d4f8796b97b34ba55555c9729a14919648c3dd7bc",
                            "6232b9b399cd86c7c0b42f9780728dea59f47565cd451010eda90e002079cdfe",
                            "d30ded9d1e3b0067d667a2a04fdacff895099ca7371d20d36d6e6565a4180e10",
                            "f248f84f019e573ba913827b9e55ef57cdd21d622b7000ff4d89d303ca2e7c20",
                            "0e64b71334792a5e6170fd526bab7f0abb3b5a1ef0912d8b273f59c6e3ff0d75"
                        ],
                        "outputs": [
                            {
                                "uxid": "ce1075dd609622b0c28d4106f58943d7f0cae6baffbe9f048bb5bc5f23706b93",
                                "dst": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "4900.000000",
                                "hours": 6402335
                            },
                            {
                                "uxid": "d73cf1f1d04a1d493fe3480a00e48187f9201bb64828fe0c638f17c0c88bb3d9",
                                "dst": "YPhukwVyLsPGX1FAPQa2ktr5XnSLqyGbr5",
                                "coins": "5.000000",
                                "hours": 6402335
                            }
                        ]
                    }
                ]
            },
            "size": 802
        },
        {
            "header": {
                "seq": 1892,
                "block_hash": "e99fe31adf3a15aab77eb81d7aaa5477a96d2ce7f74ccb3d4cdce4da2eb01cb2",
                "previous_block_hash": "015265e82b4a28eee364327bf0d3b51fed54c2067318029974c262478a09bf17",
                "timestamp": 1499403255,
                "fee": 4801753,
                "version": 0,
                "tx_body_hash": "b3c6f0f87c5282ff7ff5e6d637c2581e6a56826a76ec3dd221d02786881e3d14"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "b3c6f0f87c5282ff7ff5e6d637c2581e6a56826a76ec3dd221d02786881e3d14",
                        "inner_hash": "3dde2cd3b6562e229e99e1a0c3d87be89ab42bb773343d1471a9627cd004fc67",
                        "sigs": [
                            "22557281582f7daaf07fc79cd1c33438c06788c8561cf2fee96d81d46ffe827f44d57e0cadeb80662ea5b720c4062d5b72827dc72ea9db7afefbfe644aefea2d00"
                        ],
                        "inputs": [
                            "ce1075dd609622b0c28d4106f58943d7f0cae6baffbe9f048bb5bc5f23706b93"
                        ],
                        "outputs": [
                            {
                                "uxid": "c8b8eac053a5640bae40144cbc3dda02746071e3c7d00a4b5dfd06d28f928ec4",
                                "dst": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "2500.000000",
                                "hours": 800291
                            },
                            {
                                "uxid": "16dd81af869743599fe60108c22d7ee1fcbf1a7f460fffd3a015fbb3f721c36d",
                                "dst": "YPhukwVyLsPGX1FAPQa2ktr5XnSLqyGbr5",
                                "coins": "2400.000000",
                                "hours": 800291
                            }
                        ]
                    }
                ]
            },
            "size": 220
        }
    ]
}`,
		ExampleVerboseResponse: `{
    "blocks": [
        {
            "header": {
                "seq": 1891,
                "block_hash": "015265e82b4a28eee364327bf0d3b51fed54c2067318029974c262478a09bf17",
                "previous_block_hash": "71bf5abc05a4344597cf7df25249203df12e73b0a2e37dab386eb292a6cc5faa",
                "timestamp": 1499402935,
                "fee": 38414015,
                "version": 0,
                "tx_body_hash": "ac9759e9ffc3450d9af92e205f909375131df64e8f158a26e5e1e15334f45c8d"
            },
            "body": {
                "txns": [
                    {
                        "length": 802,
                        "type": 0,
                        "txid": "ac9759e9ffc3450d9af92e205f909375131df64e8f158a26e5e1e15334f45c8d",
                        "inner_hash": "cd36c8a00ff1fa81d94747b2d64107657fa40fa43bb0e5afc70d16661eca4286",
                        "fee": 38414015,
                        "sigs": [
                            "d34f0a9f0b82028e44702c7fc694e35aeb2777630efaaae32b9854b6c59c9ad747bdb91111864a050f024efd325a744cbbb7c1dde55a776c8f0d4856ac07fa5801",
                            "aa2f3ec98866bafda37f6e0b52be6fd1201131ce6b7edf1988b762edf14c35880b26a3a3fe9369a0f9eedecbe958abae83e54d572440903f5ff5b11cb9eb5e9200",
                            "746e171c22c99939d2d7e4b615f98de74f1c6ae152cb3d1e17e4d5e25d6f52b444242080a5a3d8f36de3befb8ec915321a7cb282eafc661c0857de753a6b40e100",
                            "e4f3db5afa209cbceac814d0f5c6fff8b57a3cc5c15063cd0b390067705de0a87313488d5eabd369855a87879a8d163fce08faeacd9d6538528d29c7cdbdaedf00",
                            "038263cdaf43a544b4ef66b2c910b9dd807ecd55e2fbf513c776b7accb1ff1fb5f591fbba22af58cf399d743955ab6cc7e3e641ecc675253306b2d0b6570672f00",
                            "215d1527af817f9c0613cacc3c2d07873ff563ab113023fc83ea836e3ac827f240df8a3d390d8e867cdda149d0db8f86bf435fa354dc4f3ac8c3a9b39297d7ba00",
                            "a2d7a50176ac1ca9417ca35f9848621bc3e0aeacb0cce713e9f93d82db4988ac5bca5bf02610c64f1580b198aadf49b04a14542466f7dca49d73b6c70f2577a201"
                        ],
                        "inputs": [
                            {
                                "uxid": "e72d8ba4ce2d3b37aeb71df2e3bed80ee07204b3fa633f56cbce7bca836bd39c",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "3.000000",
                                "hours": 0,
                                "calculated_hours": 58814
                            },
                            {
                                "uxid": "6a349ba12c5d2827de6c24773d3dd8f6572e86adba4c8954a6d6e68df9e165e2",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "2400.000000",
                                "hours": 348891,
                                "calculated_hours": 44390964
                            },
                            {
                                "uxid": "bb9a579003de101bee76d83d4f8796b97b34ba55555c9729a14919648c3dd7bc",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "2.000000",
                                "hours": 2442307,
                                "calculated_hours": 2442329
                            },
                            {
                                "uxid": "6232b9b399cd86c7c0b42f9780728dea59f47565cd451010eda90e002079cdfe",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "100.000000",
                                "hours": 2437264,
                                "calculated_hours": 2438357
                            },
                            {
                                "uxid": "d30ded9d1e3b0067d667a2a04fdacff895099ca7371d20d36d6e6565a4180e10",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "500.000000",
                                "hours": 1252147,
                                "calculated_hours": 1257601
                            },
                            {
                                "uxid": "f248f84f019e573ba913827b9e55ef57cdd21d622b7000ff4d89d303ca2e7c20",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "900.000000",
                                "hours": 304661,
                                "calculated_hours": 314456
                            },
                            {
                                "uxid": "0e64b71334792a5e6170fd526bab7f0abb3b5a1ef0912d8b273f59c6e3ff0d75",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "1000.000000",
                                "hours": 305303,
                                "calculated_hours": 316164
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "ce1075dd609622b0c28d4106f58943d7f0cae6baffbe9f048bb5bc5f23706b93",
                                "dst": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "4900.000000",
                                "hours": 6402335
                            },
                            {
                                "uxid": "d73cf1f1d04a1d493fe3480a00e48187f9201bb64828fe0c638f17c0c88bb3d9",
                                "dst": "YPhukwVyLsPGX1FAPQa2ktr5XnSLqyGbr5",
                                "coins": "5.000000",
                                "hours": 6402335
                            }
                        ]
                    }
                ]
            },
            "size": 802
        },
        {
            "header": {
                "seq": 1892,
                "block_hash": "e99fe31adf3a15aab77eb81d7aaa5477a96d2ce7f74ccb3d4cdce4da2eb01cb2",
                "previous_block_hash": "015265e82b4a28eee364327bf0d3b51fed54c2067318029974c262478a09bf17",
                "timestamp": 1499403255,
                "fee": 4801753,
                "version": 0,
                "tx_body_hash": "b3c6f0f87c5282ff7ff5e6d637c2581e6a56826a76ec3dd221d02786881e3d14"
            },
            "body": {
                "txns": [
                    {
                        "length": 220,
                        "type": 0,
                        "txid": "b3c6f0f87c5282ff7ff5e6d637c2581e6a56826a76ec3dd221d02786881e3d14",
                        "inner_hash": "3dde2cd3b6562e229e99e1a0c3d87be89ab42bb773343d1471a9627cd004fc67",
                        "fee": 4801753,
                        "sigs": [
                            "22557281582f7daaf07fc79cd1c33438c06788c8561cf2fee96d81d46ffe827f44d57e0cadeb80662ea5b720c4062d5b72827dc72ea9db7afefbfe644aefea2d00"
                        ],
                        "inputs": [
                            {
                                "uxid": "ce1075dd609622b0c28d4106f58943d7f0cae6baffbe9f048bb5bc5f23706b93",
                                "owner": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "4900.000000",
                                "hours": 6402335,
                                "calculated_hours": 6402335
                            }
                        ],
                        "outputs": [
                            {
                                "uxid": "c8b8eac053a5640bae40144cbc3dda02746071e3c7d00a4b5dfd06d28f928ec4",
                                "dst": "PRXLNyB64cqaiG4pCoFZZ8Tuv7LWYPpa7m",
                                "coins": "2500.000000",
                                "hours": 800291
                            },
                            {
                                "uxid": "16dd81af869743599fe60108c22d7ee1fcbf1a7f460fffd3a015fbb3f721c36d",
                                "dst": "YPhukwVyLsPGX1FAPQa2ktr5XnSLqyGbr5",
                                "coins": "2400.000000",
                                "hours": 800291
                            }
                        ]
                    }
                ]
            },
            "size": 220
        }
    ]
}`,
	},

	{
		ExplorerPath:   "/api/currentBalance",
		SkycoinPath:    "/api/v1/outputs",
		QueryArgs:      []string{"addrs"},
		Description:    "Returns head outputs for a list of comma-separated addresses.  If no addresses are specified, returns all head outputs.",
		ExampleRequest: "/api/currentBalance?addrs=SeDoYN6SNaTiAZFHwArnFwQmcyz7ZvJm17,iqi5BpPhEqt35SaeMLKA94XnzBG57hToNi",
		ExampleResponse: `{
    "head": {
      "seq": 65796,
      "block_hash": "3c869dfad2fdea444fe53f888c20ead67c5b0fccd8da34c3d7da580bc8a6d23c",
      "previous_block_hash": "8dce9985e05f3ac648a0a17dd60bda2737395d646b1cb42f3039eccde2a6ce7a",
      "timestamp": 1540423494,
      "fee": 6892,
      "version": 0,
      "tx_body_hash": "bd9ea7068c96ca065f511635cddcc0fb7e13bdf0b5889dd52889292ca9e2a116",
      "ux_hash": "572660420c0b463d00a1d87f320c9390456bbc93bc8be1e16808a3116f2152bc"
    },
    "head_outputs": [
        {
            "hash": "fa8161308dee3accc99a35be1fb7921dff4d24a6fc804e98d7aae7aae99d0d0d",
            "time": 1540423494,
            "block_seq": 65796,
            "src_tx": "b125abb61f5d6ec0f44422e234007b07ab276923e9533023bdd58d51e0a0f9b7",
            "address": "iqi5BpPhEqt35SaeMLKA94XnzBG57hToNi",
            "coins": "50",
            "hours": 248,
            "calculated_hours": 3445
        },
        {
            "hash": "38926afbb00c2f50d293de866bc44713eaa14c18286b24796819fbc190efcbce",
            "time": 1540423494,
            "block_seq": 50000,
            "src_tx": "8591837d905894142119923de1447ba855dcf8f34ba451970e83a2bbfea8eeca",
            "address": "iqi5BpPhEqt35SaeMLKA94XnzBG57hToNi",
            "coins": "375",
            "hours": 9,
            "calculated_hours": 3445
        },
        {
            "hash": "a70d1f0f488066a327acd0d5ea77b87d62b3b061d3db8361c90194a6520ab29f",
            "time": 1540423494,
            "block_seq": 20000,
            "src_tx": "c297eb14a9e68ec5501aa886e5bb720a58fe6466be633a8264f61eee9580a2c3",
            "address": "SeDoYN6SNaTiAZFHwArnFwQmcyz7ZvJm17",
            "coins": "51",
            "hours": 2013433,
            "calculated_hours": 3000000
        }
    ],
    "outgoing_outputs": [],
    "incoming_outputs": []
}`,
	},

	{
		ExplorerPath:   "/api/transaction",
		SkycoinPath:    "/api/v1/transaction",
		QueryArgs:      []string{"txid", "verbose"},
		Description:    "Returns transaction metadata. Assign 1 to the \"verbose\" argument to get more data in the response.",
		ExampleRequest: "/api/transaction?txid=edd2de176948cbd27cdd8cba7ca4b5afb8dbf8174dd6de95f73ce359affe4f05",
		ExampleResponse: `{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 147,
        "block_seq": 1747,
        "unknown": false
    },
    "time": 0,
    "txn": {
        "length": 220,
        "type": 0,
        "txid": "edd2de176948cbd27cdd8cba7ca4b5afb8dbf8174dd6de95f73ce359affe4f05",
        "inner_hash": "1a0b3794925960bc4a61fc1deebbdb48879f2272cdf7247e62d97b3f9df6ba72",
        "timestamp": 1498734315,
        "sigs": [
            "642a1443755bc949114b2c1a02cdf70d996a96c4971428fca39515b4143bb04b5a1922faeb45f7201fe10e40f162eeed61802fc55345ae7b9fcd12cc1a36f7b101"
        ],
        "inputs": [
            "8aab2420737e710430e95823e4853b5b9f36ab8afb679a6963564be34bfa87a4"
        ],
        "outputs": [
            {
                "uxid": "a72bebed7cf5fab28a4f57e525de5c22a34b3bd2a3bca0074e6658d618969ab0",
                "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                "coins": "1513",
                "hours": 1518
            },
            {
                "uxid": "51a9abaf13e3b404902698e6841738e7e054b59127733c60a5edd3a67c596f6d",
                "dst": "2dbpXFWsbTRRACGMcvinhKcWyXybMbv5yx8",
                "coins": "24",
                "hours": 3037
            }
        ]
    }
}`,
		ExampleVerboseResponse: `{
    "status": {
        "confirmed": true,
        "unconfirmed": false,
        "height": 55518,
        "block_seq": 1747
    },
    "time": 1498734315,
    "txn": {
        "timestamp": 1498734315,
        "length": 220,
        "type": 0,
        "txid": "edd2de176948cbd27cdd8cba7ca4b5afb8dbf8174dd6de95f73ce359affe4f05",
        "inner_hash": "1a0b3794925960bc4a61fc1deebbdb48879f2272cdf7247e62d97b3f9df6ba72",
        "fee": 16308,
        "sigs": [
            "642a1443755bc949114b2c1a02cdf70d996a96c4971428fca39515b4143bb04b5a1922faeb45f7201fe10e40f162eeed61802fc55345ae7b9fcd12cc1a36f7b101"
        ],
        "inputs": [
            {
                "uxid": "8aab2420737e710430e95823e4853b5b9f36ab8afb679a6963564be34bfa87a4",
                "owner": "2b3974pDXkdHR6VGQvVZc1iZBzrJtNfAHPY",
                "coins": "1537.000000",
                "hours": 12150,
                "calculated_hours": 20863
            }
        ],
        "outputs": [
            {
                "uxid": "a72bebed7cf5fab28a4f57e525de5c22a34b3bd2a3bca0074e6658d618969ab0",
                "dst": "2GgFvqoyk9RjwVzj8tqfcXVXB4orBwoc9qv",
                "coins": "1513.000000",
                "hours": 1518
            },
            {
                "uxid": "51a9abaf13e3b404902698e6841738e7e054b59127733c60a5edd3a67c596f6d",
                "dst": "2dbpXFWsbTRRACGMcvinhKcWyXybMbv5yx8",
                "coins": "24.000000",
                "hours": 3037
            }
        ]
    }
}`,
	},

	{
		ExplorerPath:   "/api/uxout",
		SkycoinPath:    "/api/v1/uxout",
		QueryArgs:      []string{"uxid"},
		Description:    "Returns unspent output metadata.  Note: coins are measured in drops, which are coins * 1000000.",
		ExampleRequest: "/api/uxout?uxid=374fac152e3a920975151438ad87125c0222c4e78cae2c49183cde674cb07bd7",
		ExampleResponse: `{
    "uxid": "374fac152e3a920975151438ad87125c0222c4e78cae2c49183cde674cb07bd7",
    "time": 1492186087,
    "src_block_seq": 865,
    "src_tx": "23b942eb86270f3d8bd42d2059ffbe0b1607f277870b4272b2b864fb7dfa7457",
    "owner_address": "3fEbhoHzsUwdUr8nSoofXhU9nZRDGgycCu",
    "coins": 100000000,
    "hours": 37364678,
    "spent_block_seq": 866,
    "spent_tx": "1ffa69520c15a8eb89cf68bb1a7ef9bdfcecbc490f8b70ed01206f480fa6b8ea"
}`,
	},

	{
		ExplorerPath:   "/api/pendingTxs",
		SkycoinPath:    "/api/v1/pendingTxs",
		QueryArgs:      []string{"verbose"},
		Description:    "Returns the unconfirmed transactions in the pool. Assign 1 to the \"verbose\" argument to get more data in the response.",
		ExampleRequest: "/api/pendingTxs",
		ExampleResponse: `[
    {
        "transaction": {
            "length": 317,
            "type": 0,
            "txid": "89578005d8730fe1789288ee7dea036160a9bd43234fb673baa6abd91289a48b",
            "inner_hash": "cac977eee019832245724aa643ceff451b9d8b24612b2f6a58177c79e8a4c26f",
            "sigs": [
                "3f084a0c750731dd985d3137200f9b5fc3de06069e62edea0cdd3a91d88e56b95aff5104a3e797ab4d6d417861af0c343efb0fff2e5ba9e7cf88ab714e10f38101",
                "e9a8aa8860d189daf0b1dbfd2a4cc309fc0c7250fa81113aa7258f9603d19727793c1b7533131605db64752aeb9c1f4465198bb1d8dd597213d6406a0a81ed3701"
            ],
            "inputs": [
                "bb89d4ed40d0e6e3a82c12e70b01a4bc240d2cd4f252cfac88235abe61bd3ad0",
                "170d6fd7be1d722a1969cb3f7d45cdf4d978129c3433915dbaf098d4f075bbfc"
            ],
            "outputs": [
                {
                    "uxid": "ec9cf2f6052bab24ec57847c72cfb377c06958a9e04a077d07b6dd5bf23ec106",
                    "dst": "nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
                    "coins": "60.000000",
                    "hours": 2458
                },
                {
                    "uxid": "be40210601829ba8653bac1d6ecc4049955d97fb490a48c310fd912280422bd9",
                    "dst": "2iVtHS5ye99Km5PonsB42No3pQRGEURmxyc",
                    "coins": "1.000000",
                    "hours": 2458
                }
            ]
        },
        "received": "2017-05-09T10:11:57.14303834+02:00",
        "checked": "2017-05-09T10:19:58.801315452+02:00",
        "announced": "0001-01-01T00:00:00Z",
        "is_valid": true
    }
]`,
		ExampleVerboseResponse: `[
    {
        "transaction": {
            "length": 220,
            "type": 0,
            "txid": "d455564dcf1fb666c3846cf579ff33e21c203e2923938c6563fe7fcb8573ba44",
            "inner_hash": "4e73155db8ed04a3bd2b953218efcc9122ebfbf4c55f08f50d1563e48eacf71d",
            "fee": 12855964,
            "sigs": [
                "17330c256a50e2117ddccf51f1980fc14380f0f9476432196ade3043668759847b97e1b209961458745684d9239541f79d9ca9255582864d30a540017ab84f2b01"
            ],
            "inputs": [
                {
                    "uxid": "27e7bc48ceca4d47e806a87100a8a98592b7618702e1cd479bf4c190462a6d09",
                    "owner": "23MjQipM9YsPKkYiuaBmf6m7fD54wrzHxpd",
                    "coins": "7815.000000",
                    "hours": 279089,
                    "calculated_hours": 13101146
                }
            ],
            "outputs": [
                {
                    "uxid": "4b4ebf62acbaece798d0dfc92fcea85768a2874dad8a9b8eb5454288deae468c",
                    "dst": "23MjQipM9YsPKkYiuaBmf6m7fD54wrzHxpd",
                    "coins": "586.000000",
                    "hours": 122591
                },
                {
                    "uxid": "781cfb134d5fdad48f3c937dfcfc66b169a305adc8abdfe92a0ec94c564913f2",
                    "dst": "2ehrG4VKLRuvBNWYz3U7tS75QWvzyWR89Dg",
                    "coins": "7229.000000",
                    "hours": 122591
                }
            ]
        },
        "received": "2018-06-20T14:14:52.415702671+08:00",
        "checked": "2018-08-26T19:47:45.328131142+08:00",
        "announced": "2018-08-26T19:51:47.356083569+08:00",
        "is_valid": true
    }
]`,
	},
	{
		ExplorerPath:   "/api/richlist",
		SkycoinPath:    "/api/v1/richlist",
		QueryArgs:      []string{"n", "include-distribution"},
		Description:    "Returns top N richer with unspect outputs, If no n are specified, returns 20.",
		ExampleRequest: "/api/richlist?n=4&include-distribution=false",
		ExampleResponse: `{
    "richlist": [
        {
            "address": "zMDywYdGEDtTSvWnCyc3qsYHWwj9ogws74",
            "coins": "1000000.000000",
            "locked": false
        },
        {
            "address": "z6CJZfYLvmd41GRVE8HASjRcy5hqbpHZvE",
            "coins": "1000000.000000",
            "locked": false
        },
        {
            "address": "wyQVmno9aBJZmQ99nDSLoYWwp7YDJCWsrH",
            "coins": "1000000.000000",
            "locked": false
        },
        {
            "address": "tBaeg9zE2sgmw5ZQENaPPYd6jfwpVpGTzS",
            "coins": "1000000.000000",
            "locked": false
        }
    ]
}`,
	},
	{
		ExplorerPath:   "/api/addresscount",
		SkycoinPath:    "/api/v1/addresscount",
		Description:    "Returns count number of unique addresses with unspent outputs",
		ExampleRequest: "/api/addresscount",
		ExampleResponse: `{
    "count": 10103
}`,
	},
	{
		ExplorerPath:   "/api/transactions",
		SkycoinPath:    "/api/v1/transactions",
		QueryArgs:      []string{"addrs", "confirmed", "verbose"},
		Description:    "Returns transactions for a list of comma-separated addresses. If no addresses are specified, returns all transactions. Assign 1 to the \"verbose\" argument to get more data in the response.",
		ExampleRequest: "/api/transactions?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,6dkVxyKFbFKg9Vdg6HPg1UANLByYRqkrdY&confirmed=1",
		ExampleResponse: `[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 10492,
            "block_seq": 1177,
            "unknown": false
        },
        "time": 1494275011,
        "txn": {
            "length": 317,
            "type": 0,
            "txid": "b09cd3a8baef6a449848f50a1b97943006ca92747d4e485d0647a3ea74550eca",
            "inner_hash": "2cb370051c92521a04ba5357e229d8ffa90d9d1741ea223b44dd60a1483ee0e5",
            "timestamp": 1494275011,
            "sigs": [
                "a55155ca15f73f0762f79c15917949a936658cff668647daf82a174eed95703a02622881f9cf6c7495536676f931b2d91d389a9e7b034232b3a1519c8da6fb8800",
                "cc7d7cbd6f31adabd9bde2c0deaa9277c0f3cf807a4ec97e11872817091dc3705841a6adb74acb625ee20ab6d3525350b8663566003276073d94c3bfe22fe48e01"
            ],
            "inputs": [
                "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd"
            ],
            "outputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454
                },
                {
                    "uxid": "a1268e9bd2033b49b44afa765d20876467254f51e5515626780467267a65c563",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "1.000000",
                    "hours": 7454
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 10491,
            "block_seq": 1178,
            "unknown": false
        },
        "time": 1494275231,
        "txn": {
            "length": 183,
            "type": 0,
            "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
            "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
            "timestamp": 1494275231,
            "sigs": [
                "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
            ],
            "inputs": [
                "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191"
            ],
            "outputs": [
                {
                    "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "8.000000",
                    "hours": 931
                }
            ]
        }
    }
]`,
		ExampleVerboseResponse: `[
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 56088,
            "block_seq": 1177
        },
        "time": 1494275011,
        "txn": {
            "timestamp": 1494275011,
            "length": 317,
            "type": 0,
            "txid": "b09cd3a8baef6a449848f50a1b97943006ca92747d4e485d0647a3ea74550eca",
            "inner_hash": "2cb370051c92521a04ba5357e229d8ffa90d9d1741ea223b44dd60a1483ee0e5",
            "fee": 44726,
            "sigs": [
                "a55155ca15f73f0762f79c15917949a936658cff668647daf82a174eed95703a02622881f9cf6c7495536676f931b2d91d389a9e7b034232b3a1519c8da6fb8800",
                "cc7d7cbd6f31adabd9bde2c0deaa9277c0f3cf807a4ec97e11872817091dc3705841a6adb74acb625ee20ab6d3525350b8663566003276073d94c3bfe22fe48e01"
            ],
            "inputs": [
                {
                    "uxid": "4f4b0078a9cd19b3395e54b3f42af6adc997f77f04e0ca54016c67c4f2384e3c",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "1.000000",
                    "hours": 52836,
                    "calculated_hours": 52857
                },
                {
                    "uxid": "36f4871646b6564b2f1ab72bd768a67579a1e0242bc68bcbcf1779bc75b3dddd",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 6604,
                    "calculated_hours": 6777
                }
            ],
            "outputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "dst": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454
                },
                {
                    "uxid": "a1268e9bd2033b49b44afa765d20876467254f51e5515626780467267a65c563",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "1.000000",
                    "hours": 7454
                }
            ]
        }
    },
    {
        "status": {
            "confirmed": true,
            "unconfirmed": false,
            "height": 56087,
            "block_seq": 1178
        },
        "time": 1494275231,
        "txn": {
            "timestamp": 1494275231,
            "length": 183,
            "type": 0,
            "txid": "a6446654829a4a844add9f181949d12f8291fdd2c0fcb22200361e90e814e2d3",
            "inner_hash": "075f255d42ddd2fb228fe488b8b468526810db7a144aeed1fd091e3fd404626e",
            "fee": 6523,
            "sigs": [
                "9b6fae9a70a42464dda089c943fafbf7bae8b8402e6bf4e4077553206eebc2ed4f7630bb1bd92505131cca5bf8bd82a44477ef53058e1995411bdbf1f5dfad1f00"
            ],
            "inputs": [
                {
                    "uxid": "5287f390628909dd8c25fad0feb37859c0c1ddcf90da0c040c837c89fefd9191",
                    "owner": "2K6NuLBBapWndAssUtkxKfCtyjDQDHrEhhT",
                    "coins": "8.000000",
                    "hours": 7454,
                    "calculated_hours": 7454
                }
            ],
            "outputs": [
                {
                    "uxid": "70fa9dfb887f9ef55beb4e960f60e4703c56f98201acecf2cad729f5d7e84690",
                    "dst": "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
                    "coins": "8.000000",
                    "hours": 931
                }
            ]
        }
    }
]`,
	},
	{
		ExplorerPath:   "/api/paginatedTransactions",
		SkycoinPath:    "/api/v2/transactions",
		QueryArgs:      []string{"addrs", "confirmed", "verbose", "page", "limit", "sort"},
		Description:    "This API is almost the same as the transactions version, except that it would not return all transactions by default and has pagination supported. If there are unconfirmed transactions, they will be appended after the confirmed transactions.",
		ExampleRequest: "/api/transactions?page=1024&limit=2",
		ExampleResponse: `{
    "data": {
        "page_info": {
            "total_pages": 66530,
            "page_size": 2,
            "current_page": 1024
        },
        "txns": [
            {
                "status": {
                    "confirmed": true,
                    "unconfirmed": false,
                    "height": 128216,
                    "block_seq": 2016
                },
                "time": 1500130512,
                "txn": {
                    "timestamp": 1500130512,
                    "length": 414,
                    "type": 0,
                    "txid": "f0a3c01325f3e8f09255d49b490c804b929d668fcb70ea814e1a9868b608cfdb",
                    "inner_hash": "85a298977f5fa338b7a73359c51b83787130b4f3db4a8425a1c54e45e317499d",
                    "sigs": [
                        "25333b9a283691cb189e1d2ade7dd6eeb6a275be820ff031af9b877b56330f1546a875a528bab2e559236141a644f2248a19ee5fcc86b2271f9dc60fb296f3f701",
                        "e21fdae15af052f9b842bc062ab8a2ed42baf61fe11c60255555c0fc86b99abc659269dd907472091d392d31b3c1ad24e11176ee6a9e27da1fc57e2d8ddbd04d00",
                        "5e4aa1cfca62e0a0aac1c646c3917a96bb6d1c7b8cde2e255d01730eb9d436b446cd5a09dfec097d28f5e7038a05e7172e7d5ddfe4558b1f9e3c25367051ff4f00"
                    ],
                    "inputs": [
                        "5d83e6df94ca78079c8689e700dcabdab2de959fe9f803b36fec34b47b07d025",
                        "ba1ba491090065d943ce3990b62c5d94f363bbdf37043032d79046af3687ef4c",
                        "cdce197632464ee9c46d48cb21c959772b8bf2aa04239399353988b937b6e149"
                    ],
                    "outputs": [
                        {
                            "uxid": "d19549c470bb6d217bb8095df9ef14346ee8f86730208a4247420307fadbb0f0",
                            "dst": "WSJoAtC4XcjAxTHAFLKU6MNthhpSDX7i1z",
                            "coins": "3908.000000",
                            "hours": 1070530
                        },
                        {
                            "uxid": "1742af80ec06a3ef2123a371c6f5e82c275d881e7444f8a921818bc98032fff4",
                            "dst": "2f9JhZJ147v9D4KxnJwbj8i5iNxqeKL3xNh",
                            "coins": "50.000000",
                            "hours": 1070530
                        }
                    ]
                }
            },
            {
                "status": {
                    "confirmed": true,
                    "unconfirmed": false,
                    "height": 128215,
                    "block_seq": 2017
                },
                "time": 1500130612,
                "txn": {
                    "timestamp": 1500130612,
                    "length": 414,
                    "type": 0,
                    "txid": "7dc9ae6524abe9108fdc744f210b94274a9c9fdd3da16eaea1aa88037792c27d",
                    "inner_hash": "752142fcd1cc4b9bad972611a9e64108d91b7642e1eb6b65ac92360a0c9c6bfc",
                    "sigs": [
                        "07973d43d4782ee96af70fa0ee4c73f667b035ade8770d55524b91e6d762a73b7bb6e24358c929609dd91c0140e51fa4f55952b45638bae699e522f4009c8f0c01",
                        "6d8cd2ebabcb511b1772546c898fa456bfeffefcba69aea0f6c285cda38014c015bb664f66ce0d840c71f66403782e5b6b9fd2688a212eb2e3d275aebee5856b00",
                        "4844a246eff8b59f177f9a4a43815b1fa8a3168c18d63b50c75773ee51c0b9c047ce59431871ba9fdb7df8827c9a2b175424f22cc6cf02a784b805b50574cc7000"
                    ],
                    "inputs": [
                        "8754b0d917f6690d5e88dd0950bdcb8e09d96ffb14b76964da923f5dc8969e0c",
                        "8fe7df28494563a5b47abbe737c095da461235a2529dda2a1119a19965293c8b",
                        "421ec170519fab890e3410af5ba4cf33f71fa57d786d5d39f71b7a96ed898094"
                    ],
                    "outputs": [
                        {
                            "uxid": "fceb40fd9e8895c050fc165d861ebcbe87789eeb89809879a662fdac854bf84e",
                            "dst": "2bfYafFtdkCRNcCyuDvsATV66GvBR9xfvjy",
                            "coins": "37659.000000",
                            "hours": 93972
                        },
                        {
                            "uxid": "eb8c8677da0200be7a405f2e3497db9beaa6288734d85acc5488d573ce2b8399",
                            "dst": "fXZv5X2NXhWYShoE8jazbh5UWYVCFgUXdW",
                            "coins": "25000.000000",
                            "hours": 187945
                        }
                    ]
                }
            }
        ]
    }
}`,
	},
	{
		ExplorerPath:   "/api/balance",
		SkycoinPath:    "/api/v1/balance",
		QueryArgs:      []string{"addrs"},
		Description:    "Returns the combined balance of a list of comma-separated addresses.",
		ExampleRequest: "/api/balance?addrs=7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD,nu7eSpT6hr5P21uzw7bnbxm83B6ywSjHdq",
		ExampleResponse: `{
      "confirmed": {
          "coins": 70000000,
          "hours": 28052
      },
      "predicted": {
          "coins": 9000000,
          "hours": 8385
      }
}`,
	},
	{
		ExplorerPath:   "/api/health",
		SkycoinPath:    "/api/v1/health",
		Description:    "Returns information about the current state of the node.",
		ExampleRequest: "/api/health",
		ExampleResponse: `{
      "blockchain": {
          "head": {
              "seq": 58894,
              "block_hash": "3961bea8c4ab45d658ae42effd4caf36b81709dc52a5708fdd4c8eb1b199a1f6",
              "previous_block_hash": "8eca94e7597b87c8587286b66a6b409f6b4bf288a381a56d7fde3594e319c38a",
              "timestamp": 1537581604,
              "fee": 485194,
              "version": 0,
              "tx_body_hash": "c03c0dd28841d5aa87ce4e692ec8adde923799146ec5504e17ac0c95036362dd",
              "ux_hash": "f7d30ecb49f132283862ad58f691e8747894c9fc241cb3a864fc15bd3e2c83d3"
          },
          "unspents": 38171,
          "unconfirmed": 1,
          "time_since_last_block": "4m46s"
      },
      "version": {
          "version": "0.25.0",
          "commit": "8798b5ee43c7ce43b9b75d57a1a6cd2c1295cd1e",
          "branch": "develop"
      },
      "coin": "skycoin",
      "user_agent": "skycoin:0.25.0",
      "open_connections": 8,
      "outgoing_connections": 5,
      "incoming_connections": 3,
      "uptime": "6m30.629057248s",
      "csrf_enabled": true,
      "csp_enabled": true,
      "wallet_api_enabled": true,
      "gui_enabled": true,
      "user_verify_transaction": {
          "burn_factor": 10,
          "max_transaction_size": 32768,
          "max_decimals": 3
      },
      "unconfirmed_verify_transaction": {
          "burn_factor": 10,
          "max_transaction_size": 32768,
          "max_decimals": 3
      },
      "started_at": 1542443907,
      "fiber": {
          "name": "skycoin",
          "display_name": "Skycoin",
          "ticker": "SKY",
          "coin_hours_display_name": "Coin Hours",
          "coin_hours_ticker": "SCH",
          "explorer_url": "https://explorer.skycoin.com"
      }
}`,
	},
	{
		ExplorerPath:   "/api/blockchain/progress",
		SkycoinPath:    "/api/v1/blockchain/progress",
		Description:    "Gets the blockchain progress.",
		ExampleRequest: "/api/blockchain/progress",
		ExampleResponse: `{
    "current": 2760,
    "highest": 2760,
    "peers": [
        {
            "address": "35.157.164.126:6000",
            "height": 2760
        },
        {
            "address": "63.142.253.76:6000",
            "height": 2760
        }
    ]
}`,
	},
}

var docEndpoint = APIEndpoint{
	ExplorerPath:   "/api/docs",
	Description:    "Returns this documentation as JSON",
	ExampleRequest: "/api/docs",
}

// ParsedJSONAPIEndpoint parse the ExampleResponse string into generic interface, for human-readable
// formatting when rendered in the browser as JSON
type ParsedJSONAPIEndpoint struct {
	APIEndpoint
	ParsedExampleResponse        interface{} `json:"example_response,omitempty"`
	ParsedExampleVerboseResponse interface{} `json:"example_vervose_response,omitempty"`
}

var parsedJSONAPIEndpoints []ParsedJSONAPIEndpoint

func jsonDocs(w http.ResponseWriter, _ *http.Request) {
	wrapper := struct {
		Endpoints []ParsedJSONAPIEndpoint `json:"endpoints"`
	}{
		Endpoints: parsedJSONAPIEndpoints,
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	if err := enc.Encode(wrapper); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

const docTemplate string = `<html><head>
<title>Skycoin Explorer API Documentation</title>
<style type="text/css">
code { white-space: pre; }
pre { background: #F7FAFB; }
code.inline { border-radius: 3px; padding: 0.2em; background-color: #F7FAFB; font-size: 1.3em; }
.example { margin-left: 2em; }
</style>
</head><body><div id="main">

<h1>Skycoin Explorer API Documentation</h1>

<div>
<p>
<p>The Skycoin Explorer API proxies a subset of a Skycoin node's API.</p>
<p>All endpoints start with /api</p>
<p>Further information about an endpoint can be found at the Skycoin repo.</p>
<p>Skycoin Github:<a href="https://github.com/skycoin/skycoin">https://github.com/skycoin/skycoin</a>.</p>
<p>Skycoin Explorer Github: <a href="https://github.com/skycoin/skycoin-explorer">https://github.com/skycoin/skycoin-explorer</a></p>
</p>
</div>

<hr />

{{- range . -}}
<div class="endpoint">

<h2><p><code class="inline">{{ .ExplorerPath }}</code></p></h2>

{{- if .QueryArgs -}}
<p><ul>
{{- range .QueryArgs -}}
    <li><code class="inline">{{ . }}</code></li>
{{- end -}}
</ul></p>
{{- end -}}


{{ if .Description }}
<p>{{ .Description }}</p>
{{ end }}

{{ if .ExampleRequest }}
<p>Example request:</p>
<p class="example"><code class="inline">{{ .ExampleRequest }}</code></p>
{{ end }}

{{ if .ExampleResponse }}
<p>Example response:</p>
<p class="example"><pre class="example"><code>{{ .ExampleResponse }}</code></pre></p>
{{ end }}

{{ if .ExampleVerboseResponse }}
<p>Example verbose response:</p>
<p class="example"><pre class="example"><code>{{ .ExampleVerboseResponse }}</code></pre></p>
{{ end }}

{{ if .SkycoinPath }}
<p>Internal skycoin node path:</p>
<p class="example"><code class="inline">{{ .SkycoinPath }}</code></p>
{{ end }}

<div>
<hr />
{{- end -}}

</div></body></html>`

var docTemplateBody string

func htmlDocs(w http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprintf(w, "%s", docTemplateBody); err != nil {
		log.Printf("Error writing HTML docs: %v", err)
	}
}
