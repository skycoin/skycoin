// package skycoin implements skycoin main integration tests
package skycoin

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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/visor"
)

const (
	testFixturesDir   = "testdata"
	checkpointVersion = "0.25.0"
)

var (
	// import path for github.com/skycoin/skycoin/cmd/skycoin
	ldflagsNameCmd string

	// import path for github.com/skycoin/skycoin/src/skycoin
	ldflagsNameSkyLib string
)

func TestMain(m *testing.M) {
	coin := getCoinName()
	output, err := exec.Command("go", "list", fmt.Sprintf("../../cmd/%s", coin)).CombinedOutput() //nolint:gosec
	if err != nil {
		fmt.Fprintf(os.Stderr, "go list failed: %s", output)
		os.Exit(1)
	}

	ldflagsNameCmd = strings.TrimSpace(string(output))

	output, err = exec.Command("go", "list", ".").CombinedOutput() //nolint:gosec
	if err != nil {
		fmt.Fprintf(os.Stderr, "go list failed: %s", output)
		os.Exit(1)
	}

	ldflagsNameSkyLib = strings.TrimSpace(string(output))

	ret := m.Run()
	os.Exit(ret)
}

func getCoinName() string {
	coin := os.Getenv("COIN")
	if coin == "" {
		coin = "skycoin"
	}
	return coin
}

func versionUpgradeWaitTimeout(t *testing.T) time.Duration {
	x := os.Getenv("VERSION_UPGRADE_TEST_WAIT_TIMEOUT")
	if x == "" {
		return time.Second * 5
	}

	d, err := time.ParseDuration(x)
	require.NoError(t, err)
	return d
}

func buildBinary(t *testing.T, version string) (string, func()) {
	coin := getCoinName()

	binaryName := fmt.Sprintf("%s-skycoin-pkg-%s.test", coin, version)
	binaryPath, err := filepath.Abs(binaryName)
	require.NoError(t, err)

	// coverpkgName will be like github.com/skycoin/skycoin
	coverpkgName := filepath.Dir(filepath.Dir(ldflagsNameCmd))

	// Build binary file with specific app version and db checkpoint version
	args := []string{
		"test", "-c",
		"-ldflags", fmt.Sprintf("-X %s.Version=%s -X %s.DBVerifyCheckpointVersion=%s", ldflagsNameCmd, version, ldflagsNameSkyLib, checkpointVersion),
		"-tags", "testrunmain",
		"-o", binaryPath,
		fmt.Sprintf("-coverpkg=%s/...", coverpkgName),
		fmt.Sprintf("../../cmd/%s/", coin),
	}
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
			name:         "db version 0.24.1, app version 0.24.1",
			dbFile:       "version-0.24.1.db",
			dbVersion:    "0.24.1",
			appVersion:   "0.24.1",
			shouldVerify: false,
		},
		{
			name:         "db version 0.25.0, app version 0.25.0",
			dbFile:       "version-0.25.0.db",
			dbVersion:    "0.25.0",
			appVersion:   "0.25.0",
			shouldVerify: false,
		},
		{
			name:         "db version 0.25.0, app version 0.26.0",
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
			name:         "db version 0.24.1, app version 0.26.0",
			dbFile:       "version-0.24.1.db",
			dbVersion:    "0.24.1",
			appVersion:   "0.26.0",
			shouldVerify: true,
		},
		{
			name:       "db version 0.25.0, app version 0.24.1",
			dbFile:     "version-0.25.0.db",
			dbVersion:  "0.25.0",
			appVersion: "0.24.1",
			err:        "Cannot use newer DB version=0.25.0 with older software version=0.24.1",
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

	err := os.MkdirAll("../../coverage", 0750)
	require.NoError(t, err)

	// Cache for prebuilt binaries (reduces test time by not recompiling needlessly)
	appCache := make(map[string]string)
	var appCacheLock sync.Mutex
	var cleanups []func()
	defer func() {
		for _, f := range cleanups {
			f()
		}
	}()

	for i, tc := range cases {
		coverageFile := fmt.Sprintf("../../coverage/db-verify-logic-%d.coverage.out", i)
		t.Run(tc.name, func(t *testing.T) {
			// Build the binary with a specific version
			binaryPath := func() string {
				appCacheLock.Lock()
				defer appCacheLock.Unlock()
				binaryPath := appCache[tc.appVersion]
				if binaryPath == "" {
					var cleanup func()
					binaryPath, cleanup = buildBinary(t, tc.appVersion)
					appCache[tc.appVersion] = binaryPath
					cleanups = append(cleanups, cleanup)
				}
				return binaryPath
			}()

			tmpFile := copyDBFile(t, tc.dbFile)
			defer os.Remove(tmpFile)

			// Run the binary with networking disabled
			args := append([]string{
				"-disable-networking=true",
				"-web-interface=false",
				"-download-peerlist=false",
				fmt.Sprintf("-db-path=%s", tmpFile),
				"-test.run", "^TestRunMain$",
				fmt.Sprintf("-test.coverprofile=%s", coverageFile),
			}, tc.args...)

			cmd := exec.Command(binaryPath, args...)

			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)

			scanner := bufio.NewScanner(stdout)

			err = cmd.Start()
			require.NoError(t, err)

			// Kill the process if it hasn't had an error or checked the database within a timeout,
			// so that the tests that test that the database is not checked can complete
			go time.AfterFunc(versionUpgradeWaitTimeout(t), func() {
				if tc.shouldVerify {
					cmd.Process.Kill() //nolint:errcheck
				} else {
					cmd.Process.Signal(os.Interrupt) //nolint:errcheck
				}
			})

			// Scan for an error message or for the database check marker
			didVerify := false
			foundErrMsg := false
			for scanner.Scan() {
				x := scanner.Bytes()

				if tc.err != "" && bytes.Contains(x, []byte(tc.err)) {
					foundErrMsg = true
					break
				}

				verifyMsg := "Checking database"
				if bytes.Contains(x, []byte(verifyMsg)) {
					didVerify = true
					cmd.Process.Signal(os.Interrupt) //nolint:errcheck
					break
				}
			}

			err = cmd.Wait()
			if err != nil {
				require.EqualError(t, err, "exit status 1", err.Error())
				require.NotEmpty(t, tc.err, "unexpected error: %v", err.Error())
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
