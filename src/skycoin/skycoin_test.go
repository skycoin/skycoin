// package integration_test implements skycoin main integration tests
package integration_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/visor"
)

const (
	binaryName      = "skycoin-integration-test"
	testFixturesDir = "testdata"
)

func buildBinary(t *testing.T, version string) (string, func()) {
	binaryPath, err := filepath.Abs(binaryName)
	require.NoError(t, err)

	// Build binary file with specific version
	args := []string{"build", "-ldflags", fmt.Sprintf("-X main.Version=%s", version), "-o", binaryPath, "../../../cmd/skycoin/skycoin.go"}
	cmd := exec.Command("go", args...)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	err = cmd.Start()
	require.NoError(t, err)

	output, err := ioutil.ReadAll(stdout)
	require.NoError(t, err)

	err = cmd.Wait()
	require.NoError(t, err, "Output: %s", string(output))

	return binaryPath, func() {
		if err := os.Remove(binaryPath); err != nil {
			t.Logf("Failed to remove %s: %v", binaryPath, err)
		}
	}
}

func TestDBVerifyLogic(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows doesn't support SIGINT")
		return
	}

	cases := []struct {
		name         string
		dbFile       string
		dbVersion    string
		appVersion   string
		shouldVerify bool
		args         []string
		err          string
	}{
		{
			name:         "db no version, app version 0.25.0",
			dbFile:       "no-version.db",
			dbVersion:    "",
			appVersion:   "0.25.0",
			shouldVerify: true,
		},
		{
			name:         "db version 0.25.0, app version 0.25.0",
			dbFile:       "version-0.25.0.db",
			dbVersion:    "0.25.0",
			appVersion:   "0.25.0",
			shouldVerify: false,
		},
		{
			name:         "db version 0.25.0, app version 0.26.0x",
			dbFile:       "version-0.25.0.db",
			dbVersion:    "0.25.0",
			appVersion:   "0.26.0",
			shouldVerify: false,
		},
		{
			name:         "db version 0.25.0, app version 0.26.0, force verify",
			dbFile:       "version-0.25.0.db",
			dbVersion:    "0.25.0",
			appVersion:   "0.26.0",
			args:         []string{"-verify-db=true"},
			shouldVerify: true,
		},
		{
			name:       "db version 0.25.0, app version 0.24.0",
			dbFile:     "version-0.25.0.db",
			dbVersion:  "0.25.0",
			appVersion: "0.24.0",
			err:        "Cannot use newer DB version=0.25.0 with older software version=0.24.0",
		},
	}

	copyDBFile := func(t *testing.T, dbFile string) string {
		// Copy the database file to a temp file since it will be modified by the application
		dbf, err := os.Open(filepath.Join(testFixturesDir, dbFile))
		require.NoError(t, err)
		defer dbf.Close()

		f, err := ioutil.TempFile("", fmt.Sprintf("%s.*", dbFile))
		require.NoError(t, err)
		defer f.Close()

		_, err = io.Copy(f, dbf)
		require.NoError(t, err)

		return f.Name()
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build the binary with a specific version
			binaryPath, cleanup := buildBinary(t, tc.appVersion)
			defer cleanup()

			tmpFile := copyDBFile(t, tc.dbFile)
			defer os.Remove(tmpFile)

			// Run the binary with networking disabled
			args := append([]string{
				"-disable-networking=true",
				"-web-interface=false",
				"-download-peerlist=false",
				fmt.Sprintf("-db-path=%s", tmpFile),
			}, tc.args...)

			cmd := exec.Command(binaryPath, args...)

			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)

			scanner := bufio.NewScanner(stdout)

			err = cmd.Start()
			require.NoError(t, err)

			// Kill the process if it hasn't had an error or checked the database within a timeout,
			// so that the tests that the database is not checked can complete
			go time.AfterFunc(time.Second*3, func() {
				if tc.shouldVerify {
					_ = cmd.Process.Kill()
				} else {
					_ = cmd.Process.Signal(os.Interrupt)
				}
			})

			// Scan for an error message or for the database check marker
			didVerify := false
			foundErrMsg := false
			for scanner.Scan() {
				x := scanner.Bytes()
				fmt.Println(string(x))

				if tc.err != "" && bytes.Contains(x, []byte(tc.err)) {
					foundErrMsg = true
					break
				}

				verifyMsg := "Checking database"
				if bytes.Contains(x, []byte(verifyMsg)) {
					didVerify = true
					_ = cmd.Process.Signal(os.Interrupt)
					break
				}
			}

			err = cmd.Wait()
			if err != nil {
				require.EqualError(t, err, "exit status 1", err.Error())
				require.NotEmpty(t, tc.err, err.Error())
				require.True(t, foundErrMsg)

				// Re-open the database to check that the version was not modified
				db, err := visor.OpenDB(tmpFile, false)
				require.NoError(t, err)
				defer db.Close()

				v, err := visor.GetDBVersion(db)
				require.NoError(t, err)
				require.NotNil(t, v)

				expectVersion := semver.MustParse(tc.dbVersion)
				require.Equal(t, expectVersion, *v)

				return
			}

			require.NoError(t, err)
			require.Empty(t, tc.err)
			require.False(t, foundErrMsg)
			require.Equal(t, tc.shouldVerify, didVerify)

			// Re-open the database to check that the version was added
			db, err := visor.OpenDB(tmpFile, false)
			require.NoError(t, err)
			defer db.Close()

			v, err := visor.GetDBVersion(db)
			require.NoError(t, err)
			require.NotNil(t, v)

			appVersion := semver.MustParse(tc.appVersion)

			expectVersion := appVersion

			if tc.dbVersion != "" {
				dbVersion := semver.MustParse(tc.dbVersion)
				if appVersion.LT(dbVersion) {
					expectVersion = dbVersion
				}
			}

			require.Equal(t, expectVersion, *v)
		})
	}
}
