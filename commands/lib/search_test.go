package lib

import (
	"strings"
	"testing"

	"github.com/arduino/arduino-cli/arduino/libraries/librariesmanager"
	rpc "github.com/arduino/arduino-cli/rpc/cc/arduino/cli/commands/v1"
	paths "github.com/arduino/go-paths-helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var customIndexPath = paths.New("testdata", "test1")
var fullIndexPath = paths.New("testdata", "full")

func TestSearchLibrary(t *testing.T) {
	lm := librariesmanager.NewLibraryManager(customIndexPath, nil)
	lm.LoadIndex()

	resp := searchLibrary(&rpc.LibrarySearchRequest{Query: "test"}, lm)
	assert := assert.New(t)
	assert.Equal(resp.GetStatus(), rpc.LibrarySearchStatus_LIBRARY_SEARCH_STATUS_SUCCESS)
	assert.Equal(len(resp.GetLibraries()), 2)
	assert.True(strings.Contains(resp.GetLibraries()[0].Name, "Test"))
	assert.True(strings.Contains(resp.GetLibraries()[1].Name, "Test"))
}

func TestSearchLibrarySimilar(t *testing.T) {
	lm := librariesmanager.NewLibraryManager(customIndexPath, nil)
	lm.LoadIndex()

	resp := searchLibrary(&rpc.LibrarySearchRequest{Query: "arduino"}, lm)
	assert := assert.New(t)
	assert.Equal(resp.GetStatus(), rpc.LibrarySearchStatus_LIBRARY_SEARCH_STATUS_SUCCESS)
	assert.Equal(len(resp.GetLibraries()), 2)
	libs := map[string]*rpc.SearchedLibrary{}
	for _, l := range resp.GetLibraries() {
		libs[l.Name] = l
	}
	assert.Contains(libs, "ArduinoTestPackage")
	assert.Contains(libs, "Arduino")
}

func TestSearchLibraryFields(t *testing.T) {
	lm := librariesmanager.NewLibraryManager(fullIndexPath, nil)
	lm.LoadIndex()

	query := func(q string) []string {
		libs := []string{}
		for _, lib := range searchLibrary(&rpc.LibrarySearchRequest{Query: q}, lm).Libraries {
			libs = append(libs, lib.Name)
		}
		return libs
	}

	res := query("SparkFun_u-blox_GNSS")
	require.Len(t, res, 3)
	require.Equal(t, "SparkFun u-blox Arduino Library", res[0])
	require.Equal(t, "SparkFun u-blox GNSS Arduino Library", res[1])
	require.Equal(t, "SparkFun u-blox SARA-R5 Arduino Library", res[2])

	res = query("SparkFun u-blox GNSS")
	require.Len(t, res, 3)
	require.Equal(t, "SparkFun u-blox Arduino Library", res[0])
	require.Equal(t, "SparkFun u-blox GNSS Arduino Library", res[1])
	require.Equal(t, "SparkFun u-blox SARA-R5 Arduino Library", res[2])

	res = query("painlessMesh")
	require.Len(t, res, 1)
	require.Equal(t, "Painless Mesh", res[0])

	res = query("cristian maglie")
	require.Len(t, res, 2)
	require.Equal(t, "Arduino_ConnectionHandler", res[0])
	require.Equal(t, "FlashStorage_SAMD", res[1])

	res = query("flashstorage")
	require.Len(t, res, 19)
	require.Equal(t, "FlashStorage", res[0])
}
