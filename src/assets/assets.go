package assets

import (
	"strings"
)

const (
	IfNoneMatch  = "If-None-Match"
	ETag         = "ETag"
	CacheControl = "Cache-Control"

	GZip      = "gzip"
	Identity  = "identity"
	NoCache   = "public, no-cache"
	Immutable = "public, immutable"
)

type Info struct {
	FilePath        string
	ContentType     string
	ContentLength   string
	ContentEncoding string
	ETag            string
}

type Assets struct {
	assetsUri    string
	assetsUriLen int
	assetsPath   string
	cache        map[string]*Info
}

func NewAssetsContext(assetsUri, assetsPath string) *Assets {
	if !strings.HasPrefix(assetsUri, "/") {
		assetsUri = "/" + assetsUri
	}
	if !strings.HasSuffix(assetsUri, "/") {
		assetsUri += "/"
	}

	if !strings.HasSuffix(assetsPath, "/") {
		assetsPath += "/"
	}

	return &Assets{
		assetsUri:    assetsUri,
		assetsUriLen: len(assetsUri),
		assetsPath:   assetsPath,
		cache:        make(map[string]*Info),
	}
}
