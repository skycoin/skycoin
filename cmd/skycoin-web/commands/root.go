// Package commands provides commands for the skycoin web interface.
package commands

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skycoin/skywire/pkg/skywire-utilities/pkg/calvin"
	"github.com/spf13/cobra"

	wasmtinygo "github.com/skycoin/skycoin/src/skycoin-lite/wasm-tinygo"
	"github.com/skycoin/skycoin/src/skycoin-web/src/gui"
)

var (
	port    int
	host    string
	nodeURL string
)

// RootCmd is the root cil command
var RootCmd = &cobra.Command{
	Use:   "skycoin-web",
	Short: "Skycoin Web Wallet",
	Long: func() (ret string) {
		ret = calvin.AsciiFont("skycoin-web")
		ret += "\nThin client web wallet for Skycoin and fibercoins."
		return ret
	}(),
	Run: func(_ *cobra.Command, _ []string) {
		serve()
	},
}

func init() {
	RootCmd.Flags().IntVarP(&port, "port", "p", 8001, "Port to serve on")
	RootCmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "Host to bind to")
	RootCmd.Flags().StringVarP(&nodeURL, "node-url", "n", "https://node.skycoin.com", "node URL")
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func serve() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Get the embedded dist directory
	distSub, err := fs.Sub(gui.DistFS, "dist")
	if err != nil {
		log.Fatalf("Failed to get dist subdirectory: %v", err)
	}

	// Serve embedded WASM files from skycoin-lite
	router.GET("/assets/scripts/skycoin-lite.wasm", func(c *gin.Context) {
		c.Header("Content-Type", "application/wasm")
		c.Data(http.StatusOK, "application/wasm", wasmtinygo.WasmFile)
	})

	router.GET("/assets/scripts/wasm_exec.js", func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript")
		c.Data(http.StatusOK, "application/javascript", wasmtinygo.WasmExecJS)
	})

	// Proxy all /api/* requests to the configured node
	router.Any("/api/*path", func(c *gin.Context) {
		// Set CORS headers first
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")

		// Handle OPTIONS preflight requests
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		// Build target URL: nodeURL + request path
		targetURL := nodeURL + c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			targetURL += "?" + c.Request.URL.RawQuery
		}

		log.Printf("[PROXY] %s %s -> %s", c.Request.Method, c.Request.URL.Path, targetURL)

		// Create proxy request
		proxyReq, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			log.Printf("[PROXY] Failed to create request: %v", err)
			c.String(http.StatusInternalServerError, "Failed to create proxy request")
			return
		}

		for name, values := range c.Request.Header {
			// Skip headers that could trigger CSRF/CORS issues
			if name == "Referer" || name == "Origin" || name == "Host" {
				continue
			}
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// Execute request to node
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			log.Printf("[PROXY] Request failed: %v", err)
			c.String(http.StatusBadGateway, "Failed to proxy request to node: %v", err)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}()

		log.Printf("[PROXY] Response: %d", resp.StatusCode)

		for name, values := range resp.Header {
			if name != "Access-Control-Allow-Origin" && name != "Access-Control-Allow-Methods" && name != "Access-Control-Allow-Headers" {
				for _, value := range values {
					c.Header(name, value)
				}
			}
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[PROXY] Failed to read response: %v", err)
			c.String(http.StatusInternalServerError, "Failed to read proxy response")
			return
		}

		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})

	// Serve static files from embedded dist directory
	router.NoRoute(func(c *gin.Context) {
		// Log all non-API requests
		if c.Request.URL.Path != "/" && c.Request.URL.Path != "/favicon.ico" {
			log.Printf("[STATIC] %s %s", c.Request.Method, c.Request.URL.Path)
		}

		// Serve from embedded filesystem
		fileServer := http.FileServer(http.FS(distSub))
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Skycoin Web Wallet starting...\n")
	fmt.Printf("Server listening on http://%s\n", addr)
	fmt.Printf("Proxying to node: %s\n", nodeURL)
	fmt.Printf("Open your browser and navigate to the address above\n")
	fmt.Printf("Press Ctrl+C to stop the server\n\n")

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
