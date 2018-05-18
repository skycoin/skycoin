package testsuite

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/file"
)

const (
	testdataDir           = "./testdata/"
	manyAddressesFilename = "many-addresses.golden"
	inputHashesFilename   = "input-hashes.golden"
	seedFileRegex         = `seed-\d+.golden`
)

func TestManyAddresses(t *testing.T) {
	fn := filepath.Join(testdataDir, manyAddressesFilename)

	var dataJSON SeedTestDataJSON
	err := file.LoadJSON(fn, &dataJSON)
	require.NoError(t, err)

	data, err := SeedTestDataFromJSON(&dataJSON)
	require.NoError(t, err)

	err = ValidateSeedData(data, nil)
	require.NoError(t, err)
}

func TestSeedSignatures(t *testing.T) {
	fn := filepath.Join(testdataDir, inputHashesFilename)

	var inputDataJSON InputTestDataJSON
	err := file.LoadJSON(fn, &inputDataJSON)
	require.NoError(t, err)

	inputData, err := InputTestDataFromJSON(&inputDataJSON)
	require.NoError(t, err)

	seedFiles := traverseFiles(testdataDir, seedFileRegex)

	for _, fn := range seedFiles {
		t.Run(fn, func(t *testing.T) {
			fn = filepath.Join(testdataDir, fn)

			var seedDataJSON SeedTestDataJSON
			err := file.LoadJSON(fn, &seedDataJSON)
			require.NoError(t, err)

			seedData, err := SeedTestDataFromJSON(&seedDataJSON)
			require.NoError(t, err)

			err = ValidateSeedData(seedData, inputData)
			require.NoError(t, err)
		})
	}
}

func traverseFiles(dir string, filenameTemplate string) []string { // nolint: unparam
	files := make([]string, 0)
	filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(filenameTemplate, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}
