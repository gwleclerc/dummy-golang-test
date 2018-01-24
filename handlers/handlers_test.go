package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gwleclerc/dummy-golang-test/handlers"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type stubCache struct {
	StoredValues map[string]string
}

// Get value from storedValues
func (c *stubCache) Get(key string) (string, error) {
	return c.StoredValues[key], nil
}

// Set value into storedValues
func (c *stubCache) Set(key, value string) error {
	c.StoredValues[key] = value
	return nil
}

// Close stub
func (c *stubCache) Close() {
	return
}

func initCtx(c echo.Context, path, key, value string) (echo.Context, *stubCache) {
	c.SetPath(path)
	c.SetParamNames(key)
	c.SetParamValues(value)
	store := stubCache{StoredValues: map[string]string{}}
	c.Set("cache", &store)
	return c, &store
}

func TestHandlerGet(t *testing.T) {
	const expectedContent = "test content"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/redis/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c, store := initCtx(c, "/redis/:key", "key", "test")
	store.StoredValues["test"] = expectedContent
	h := handlers.NewCacheHandler()

	// Assertions
	if assert.NoError(t, h.Get(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedContent, rec.Body.String())
	}
}

func TestHandlerSet(t *testing.T) {
	const savedContent = "test content"
	const expectedMsg = "Value stored successfully"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/redis/test", nil)
	req.Form = url.Values{
		"value": []string{savedContent},
	}
	req.Header = http.Header{"Content-Type": []string{"form-data"}}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c, store := initCtx(c, "/redis/:key", "key", "test")
	h := handlers.NewCacheHandler()

	// Assertions
	if assert.NoError(t, h.Set(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedMsg, rec.Body.String())
		assert.Equal(t, savedContent, store.StoredValues["test"])
	}
}
