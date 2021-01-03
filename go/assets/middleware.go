package assets

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strings"
)

func (assets *Assets) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			path := req.URL.Path
			method := req.Method

			if (method != http.MethodGet && method != http.MethodHead) || !strings.HasPrefix(path, assets.assetsUri) {
				return next(c)
			}

			info, err := assets.Get(path)
			if err != nil {
				return err
			}

			resp := c.Response()

			if info == nil {
				return c.String(http.StatusNotFound, fmt.Sprintf("Can't find asset %q", path))
			}

			ifNoneMatch := req.Header.Get(IfNoneMatch)
			if ifNoneMatch == info.ETag {
				resp.WriteHeader(http.StatusNotModified)
				return nil
			}

			cacheControl := NoCache
			if req.URL.Query().Get("v") == info.ETag {
				cacheControl = Immutable
			}

			header := resp.Header()
			header.Set(echo.HeaderContentType, info.ContentType)
			header.Set(echo.HeaderContentLength, info.ContentLength)
			header.Set(echo.HeaderContentEncoding, info.ContentEncoding)
			header.Set(ETag, info.ETag)
			header.Set(CacheControl, cacheControl)
			resp.WriteHeader(http.StatusOK)

			if method == http.MethodHead {
				return nil
			}

			file, err := os.Open(info.FilePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(resp.Writer, file)

			return err
		}
	}
}
