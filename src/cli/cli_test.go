package cli

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/testutil"
)

func Example() {
	// In cmd/cli/cli.go:
	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cli, err := NewCLI(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Run("set COIN", func(t *testing.T) {
		val := "foocoin"
		os.Setenv("COIN", val)
		defer os.Unsetenv("COIN")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.Coin, val)
	})

	t.Run("set RPC_ADDR", func(t *testing.T) {
		val := "http://111.22.33.44:5555"
		os.Setenv("RPC_ADDR", val)
		defer os.Unsetenv("RPC_ADDR")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.RPCAddress, val)
	})

	t.Run("set RPC_ADDR invalid", func(t *testing.T) {
		val := "111.22.33.44:5555"
		os.Setenv("RPC_ADDR", val)
		defer os.Unsetenv("RPC_ADDR")

		_, err := LoadConfig()
		testutil.RequireError(t, err, "RPC_ADDR must be in scheme://host format")
	})

	t.Run("set DATA_DIR", func(t *testing.T) {
		val := "/home/foo/"
		os.Setenv("DATA_DIR", val)
		defer os.Unsetenv("DATA_DIR")

		cfg, err := LoadConfig()
		require.NoError(t, err)
		require.Equal(t, cfg.DataDir, val)
	})
}
