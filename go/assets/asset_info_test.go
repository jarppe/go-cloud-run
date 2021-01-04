package assets

import (
	asserts "github.com/stretchr/testify/assert"
	"testing"
)

var testResources = "../test-resources"

var testFiles = []struct {
	name string
	size string
	sha1 string
	ct   string
}{
	{"foo.html", "149", "a2dc2892242e8e821e9e5bde", "text/html; charset=utf-8"},
	{"bar.html", "127", "7a97a9a0c11ea35238a7698e", "text/html; charset=utf-8"},
	{"bin", "10240", "057692a7632bab6b85612dee", "application/octet-stream"},
}

func TestGet(t *testing.T) {
	assert := asserts.New(t)
	assets := NewAssetsContext("/assets/", testResources)
	for _, testFile := range testFiles {
		info, err := assets.Get("/assets/" + testFile.name)
		if assert.NoError(err) && assert.NotNil(t, info) {
			assert.Equal(testFile.size, info.ContentLength)
			assert.Equal(testFile.sha1, info.ETag)
		}
	}
}

func TestGetNonExistentFile(t *testing.T) {
	assert := asserts.New(t)
	assets := NewAssetsContext("/assets/", testResources)
	info, err := assets.Get("/assets/" + "fofo")
	assert.Nil(info)
	assert.Nil(err)
}
