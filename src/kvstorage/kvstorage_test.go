package kvstorage

import (
	"testing"

	"github.com/skycoin/skycoin/src/util/file"

	"github.com/stretchr/testify/require"
)

func formTestFile(fileName string) error {
	data := map[string]string{
		"test1": "some value",
		"test2": "{\"key\":\"val\",\"key2\":2}",
	}

	return file.SaveJSON(fileName, data, 0644)
}

func TestGet(t *testing.T) {
	type expect struct {
		val string
		err error
	}

	tt := []struct {
		name     string
		fileName string
		key      string
		expect   expect
	}{
		{
			name:     "no such file",
			fileName: "nofile.txt",
			key:      "key",
			expect:   expect{},
		},
		{
			name:     "no such key",
			fileName: "./testdata/data.dat",
			key:      "no such key",
			expect: expect{
				val: "",
				err: ErrNoSuchKey,
			},
		},
		{
			name:     "simple string value",
			fileName: "./testdata/data.dat",
			key:      "test1",
			expect: expect{
				val: "some value",
				err: nil,
			},
		},
		{
			name:     "complex marshaled data",
			fileName: "./testdata/data.dat",
			key:      "test2",
			expect: expect{
				val: "{\"key\":\"val\",\"key2\":2}",
				err: nil,
			},
		},
	}

	err := formTestFile("./testdata/data.dat")
	require.NoError(t, err)

	// run separately to control the error value which might be dependent
	// on the OS, so we just check if the err is not nil in this case
	t.Run(tt[0].name, func(t *testing.T) {
		_, err := Get(tt[0].fileName, tt[0].key)
		require.Error(t, err)
	})

	for i := 1; i < len(tt); i++ {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			val, err := Get(tc.fileName, tc.key)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.val, val)
		})
	}
}

func TestGetAll(t *testing.T) {
	type expect struct {
		data map[string]string
		err  error
	}

	tt := []struct {
		name     string
		fileName string
		expect   expect
	}{
		{
			name:     "no such file",
			fileName: "nofile.txt",
			expect:   expect{},
		},
		{
			name:     "file is ok",
			fileName: "./testdata/data.dat",
			expect: expect{
				data: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
	}

	err := formTestFile("./testdata/data.dat")
	require.NoError(t, err)

	// run separately to control the error value which might be dependent
	// on the OS, so we just check if the err is not nil in this case
	t.Run(tt[0].name, func(t *testing.T) {
		_, err := GetAll(tt[0].fileName)
		require.Error(t, err)
	})

	for i := 1; i < len(tt); i++ {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			data, err := GetAll(tc.fileName)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.data, data)
		})
	}
}

func TestAdd(t *testing.T) {
	type expect struct {
		newData map[string]string
		err     error
	}

	tt := []struct {
		name     string
		fileName string
		key      string
		val      string
		expect   expect
	}{
		{
			name:     "no such file",
			fileName: "nofile.txt",
			key:      "key",
			val:      "val",
		},
		{
			name:     "file is ok, add new value",
			fileName: "./testdata/data.dat",
			key:      "new key",
			val:      "new value",
			expect: expect{
				newData: map[string]string{
					"test1":   "some value",
					"test2":   "{\"key\":\"val\",\"key2\":2}",
					"new key": "new value",
				},
			},
		},
		{
			name:     "file is ok, replace old value",
			fileName: "./testdata/data.dat",
			key:      "test1",
			val:      "oiuy",
			expect: expect{
				newData: map[string]string{
					"test1": "oiuy",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
	}

	err := formTestFile("./testdata/data.dat")
	require.NoError(t, err)

	// run separately to control the error value which might be dependent
	// on the OS, so we just check if the err is not nil in this case
	t.Run(tt[0].name, func(t *testing.T) {
		_, err := GetAll(tt[0].fileName)
		require.Error(t, err)
	})

	for i := 1; i < len(tt); i++ {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			// acquire the original data
			originalData, err := GetAll(tc.fileName)
			require.NoError(t, err)

			err = Add(tc.fileName, tc.key, tc.val)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			modifiedData, err := GetAll(tc.fileName)
			require.NoError(t, err)

			// resave the original data back to file
			err = file.SaveJSON(tc.fileName, originalData, 0644)
			require.NoError(t, err)

			require.Equal(t, tc.expect.newData, modifiedData)
		})
	}
}

func TestRemove(t *testing.T) {
	type expect struct {
		newData map[string]string
		err     error
	}

	tt := []struct {
		name     string
		fileName string
		key      string
		expect   expect
	}{
		{
			name:     "no such file",
			fileName: "nofile.txt",
			key:      "key",
		},
		{
			name:     "file is ok, no such key",
			fileName: "./testdata/data.dat",
			key:      "no key",
			expect: expect{
				newData: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
				err: ErrNoSuchKey,
			},
		},
		{
			name:     "file is ok, removed",
			fileName: "./testdata/data.dat",
			key:      "test1",
			expect: expect{
				newData: map[string]string{
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
	}

	err := formTestFile("./testdata/data.dat")
	require.NoError(t, err)

	// run separately to control the error value which might be dependent
	// on the OS, so we just check if the err is not nil in this case
	t.Run(tt[0].name, func(t *testing.T) {
		_, err := GetAll(tt[0].fileName)
		require.Error(t, err)
	})

	for i := 1; i < len(tt); i++ {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			// acquire the original data
			originalData, err := GetAll(tc.fileName)
			require.NoError(t, err)

			err = Remove(tc.fileName, tc.key)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			newData, err := GetAll(tc.fileName)
			require.NoError(t, err)

			// resave the original data back to file
			err = file.SaveJSON(tc.fileName, originalData, 0644)
			require.NoError(t, err)

			require.Equal(t, tc.expect.newData, newData)
		})
	}
}
