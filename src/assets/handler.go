package assets

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (assets *Assets) Handler(content bytes.Buffer, contentType string) echo.HandlerFunc {
	var gzipped bytes.Buffer
	gzipper := gzip.NewWriter(&gzipped)
	etagger := sha1.New()
	if _, err := io.Copy(io.MultiWriter(gzipper, etagger), bytes.NewReader(content.Bytes())); err != nil {
		log.Fatalf("Can't gzip content: %v", err)
	}
	if err := gzipper.Close(); err != nil {
		log.Fatalf("Can't close gzip writer: %v", err)
	}

	data := gzipped.Bytes()
	contentLength := strconv.Itoa(len(data))
	etag := hex.EncodeToString(etagger.Sum(nil)[:12])

	return func(c echo.Context) error {
		req := c.Request()
		resp := c.Response()

		ifNoneMatch := req.Header.Get(IfNoneMatch)
		if ifNoneMatch == etag {
			resp.WriteHeader(http.StatusNotModified)
			return nil
		}

		cacheControl := NoCache
		if req.URL.Query().Get("v") == etag {
			cacheControl = Immutable
		}

		header := resp.Header()
		header.Set(echo.HeaderContentType, contentType)
		header.Set(echo.HeaderContentLength, contentLength)
		header.Set(echo.HeaderContentEncoding, "gzip")
		header.Set("Cache-Control", cacheControl)
		header.Set("ETag", etag)
		resp.WriteHeader(http.StatusOK)
		_, err := io.Copy(resp.Writer, bytes.NewReader(data))
		return err
	}
}
