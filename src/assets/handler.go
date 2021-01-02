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

func (assets *Assets) Handler(content []byte, contentType string) echo.HandlerFunc {
	var gzipped bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipped)
	etagWriter := sha1.New()
	if _, err := io.Copy(io.MultiWriter(gzipWriter, etagWriter), bytes.NewReader(content)); err != nil {
		log.Fatalf("Can't gzip content: %v", err)
	}
	if err := gzipWriter.Close(); err != nil {
		log.Fatalf("Can't close gzip writer: %v", err)
	}

	data := gzipped.Bytes()
	contentLength := strconv.Itoa(len(data))
	etag := hex.EncodeToString(etagWriter.Sum(nil)[:12])

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
		header.Set(echo.HeaderContentEncoding, GZip)
		header.Set(CacheControl, cacheControl)
		header.Set(ETag, etag)
		resp.WriteHeader(http.StatusOK)
		_, err := io.Copy(resp.Writer, bytes.NewReader(data))
		return err
	}
}
