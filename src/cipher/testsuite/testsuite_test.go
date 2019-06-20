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
	seedFileRegex         = `^seed-\d+.golden$`
	bip32SeedFileRegex    = `^seed-bip32-\d+.golden$`
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

	seedFiles, err := traverseFiles(testdataDir, seedFileRegex)
	require.NoError(t, err)

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

func TestBip32SeedSignatures(t *testing.T) {
	fn := filepath.Join(testdataDir, inputHashesFilename)

	var inputDataJSON InputTestDataJSON
	err := file.LoadJSON(fn, &inputDataJSON)
	require.NoError(t, err)

	inputData, err := InputTestDataFromJSON(&inputDataJSON)
	require.NoError(t, err)

	seedFiles, err := traverseFiles(testdataDir, bip32SeedFileRegex)
	require.NoError(t, err)

	for _, fn := range seedFiles {
		t.Run(fn, func(t *testing.T) {
			fn = filepath.Join(testdataDir, fn)

			var seedDataJSON Bip32SeedTestDataJSON
			err := file.LoadJSON(fn, &seedDataJSON)
			require.NoError(t, err)

			seedData, err := Bip32SeedTestDataFromJSON(&seedDataJSON)
			require.NoError(t, err)

			err = ValidateBip32SeedData(seedData, inputData)
			require.NoError(t, err)
		})
	}
}

func traverseFiles(dir string, filenameTemplate string) ([]string, error) { //nolint:unparam
	files := make([]string, 0)
	if err := filepath.Walk(dir, func(_ string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(filenameTemplate, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
