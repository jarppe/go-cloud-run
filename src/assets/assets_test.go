package assets

import (
	"github.com/labstack/echo/v4"
	asserts "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandler(t *testing.T) {
	assert := asserts.New(t)
	e := echo.New()
	assets := NewAssetsContext("/assets/", testResources)
	handler := assets.Middleware()(echo.NotFoundHandler)

	for _, testFile := range testFiles {
		req := GET(testFile.name, nil)
		resp := httptest.NewRecorder()
		c := e.NewContext(req, resp)
		if assert.NoError(handler(c)) {
			assert.Equal(http.StatusOK, resp.Code)
			assert.Equal(testFile.sha1, resp.Header().Get("ETag"))
			assert.Equal(testFile.ct, resp.Header().Get("Content-Type"))
			assert.Equal(testFile.size, resp.Header().Get("Content-Length"))
			assert.Equal(NoCache, resp.Header().Get("Cache-Control"))
		}
	}
}

func TestHandler404(t *testing.T) {
	assert := asserts.New(t)
	e := echo.New()
	assets := NewAssetsContext("/assets/", testResources)
	handler := assets.Middleware()(echo.NotFoundHandler)

	req := GET("fofo", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	if assert.NoError(handler(c)) {
		assert.Equal(http.StatusNotFound, resp.Code)
	}
}

func TestHandlerCacheHeaders(t *testing.T) {
	assert := asserts.New(t)
	e := echo.New()
	assets := NewAssetsContext("/assets/", testResources)
	handler := assets.Middleware()(echo.NotFoundHandler)

	testFile := testFiles[0]

	// Request without query param `v` should set Cache-Control to NoCache

	req := GET(testFile.name, nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	if assert.NoError(handler(c)) {
		assert.Equal(http.StatusOK, resp.Code)
		assert.Equal(NoCache, resp.Header().Get("Cache-Control"))
	}

	// With incorrect query param `v` should set Cache-Control to NoCache

	q := make(url.Values)
	q.Set("v", "1234567890")
	req = GET(testFile.name, q)
	resp = httptest.NewRecorder()
	c = e.NewContext(req, resp)

	if assert.NoError(handler(c)) {
		assert.Equal(http.StatusOK, resp.Code)
		assert.Equal(NoCache, resp.Header().Get("Cache-Control"))
	}

	// With correct query param `v` should set Cache-Control to Immutable

	q = make(url.Values)
	q.Set("v", testFile.sha1)
	req = GET(testFile.name, q)
	resp = httptest.NewRecorder()
	c = e.NewContext(req, resp)

	if assert.NoError(handler(c)) {
		assert.Equal(http.StatusOK, resp.Code)
		assert.Equal(Immutable, resp.Header().Get("Cache-Control"))
	}
}

func GET(assetName string, query url.Values) *http.Request {
	requestUrl := "/assets/" + assetName
	if query != nil {
		requestUrl += "?" + query.Encode()
	}
	return httptest.NewRequest(http.MethodGet, requestUrl, nil)
}
