package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var (
	wasmSourceDir string
	outputGoDir   string
	outputTinyDir string
	wasmFileName  string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate WASM files for Go and TinyGo",
		Long:  `Compiles WASM from source using both Go and TinyGo, and copies their respective wasm_exec.js files.`,
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&wasmSourceDir, "source", "s", "wasm", "Directory containing WASM source code")
	rootCmd.Flags().StringVarP(&outputGoDir, "go-output", "g", "wasm-go", "Output directory for Go-compiled WASM")
	rootCmd.Flags().StringVarP(&outputTinyDir, "tinygo-output", "t", "wasm-tinygo", "Output directory for TinyGo-compiled WASM")
	rootCmd.Flags().StringVarP(&wasmFileName, "filename", "f", "skycoin-lite.wasm", "Output WASM filename")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Validate source directory exists
	if _, err := os.Stat(wasmSourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory does not exist: %s", wasmSourceDir)
	}

	// Create output directories
	os.MkdirAll(outputGoDir, 0755)
	os.MkdirAll(outputTinyDir, 0755)

	// Copy wasm_exec.js files
	copyWasmExecJS()

	// Compile WASM files
	compileWASM()
}

func copyWasmExecJS() {
	fmt.Println("Copying wasm_exec.js files...")

	// Go's wasm_exec.js
	goWasmExec := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
	if _, err := os.Stat(goWasmExec); os.IsNotExist(err) {
		// Try alternate location for Go 1.21+
		goWasmExec = filepath.Join(runtime.GOROOT(), "lib", "wasm", "wasm_exec.js")
	}

	// TinyGo's wasm_exec.js
	tinygoRoot := strings.TrimSuffix(runtime.GOROOT(), "go") + "tinygo"
	tinygoWasmExec := filepath.Join(tinygoRoot, "targets", "wasm_exec.js")

	// Copy Go wasm_exec.js
	if _, err := script.File(goWasmExec).WriteFile(filepath.Join(outputGoDir, "wasm_exec.js")); err != nil {
		log.Fatalf("Failed to copy Go wasm_exec.js: %v", err)
	}
	fmt.Printf("✓ Copied %s/wasm_exec.js\n", outputGoDir)

	// Copy TinyGo wasm_exec.js
	if _, err := script.File(tinygoWasmExec).WriteFile(filepath.Join(outputTinyDir, "wasm_exec.js")); err != nil {
		log.Fatalf("Failed to copy TinyGo wasm_exec.js: %v", err)
	}
	fmt.Printf("✓ Copied %s/wasm_exec.js\n", outputTinyDir)
}

func compileWASM() {
	s := spinner.New(spinner.CharSets[14], 25*time.Millisecond)

	// Compile with Go
	s.Suffix = " Compiling with Go..."
	s.Start()

	goOutput := filepath.Join(outputGoDir, wasmFileName)
	goBuildCmd := fmt.Sprintf(
		`bash -c 'cd %s || exit 1 ; time GOOS=js GOARCH=wasm go build -o ../%s -ldflags="-s -w" . && cd .. && du -h %s'`,
		wasmSourceDir, goOutput, goOutput,
	)
	fmt.Println("\nRunning:", goBuildCmd)

	output, err := script.Exec(goBuildCmd).String()
	s.Stop()
	if err != nil {
		log.Fatalf("Go build failed: %v\nOutput: %s", err, output)
	}
	fmt.Println(output)
	fmt.Printf("✓ Compiled %s\n", goOutput)

	// Compile with TinyGo
	s.Suffix = " Compiling with TinyGo..."
	s.Start()

	tinyOutput := filepath.Join(outputTinyDir, wasmFileName)
	tinyGoBuildCmd := fmt.Sprintf(
		`bash -c 'cd %s || exit 1 ; time GOOS=js GOARCH=wasm tinygo build -target=wasm --no-debug -o ../%s . && cd .. && du -h %s'`,
		wasmSourceDir, tinyOutput, tinyOutput,
	)
	fmt.Println("\nRunning:", tinyGoBuildCmd)

	output, err = script.Exec(tinyGoBuildCmd).String()
	s.Stop()
	if err != nil {
		log.Fatalf("TinyGo build failed: %v\nOutput: %s", err, output)
	}
	fmt.Println(output)
	fmt.Printf("✓ Compiled %s\n", tinyOutput)

	fmt.Println("\n✅ All WASM files generated successfully!")
}
