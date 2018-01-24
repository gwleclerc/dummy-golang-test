package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gwleclerc/dummy-golang-test/cache/mocks"
	"github.com/gwleclerc/dummy-golang-test/handlers"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func initCtx(c echo.Context, path, key, value string) echo.Context {
	c.SetPath(path)
	c.SetParamNames(key)
	c.SetParamValues(value)
	return c
}

func TestHandlerGet(t *testing.T) {
	assert := assert.New(t)
	const expectedContent = "test content"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/redis/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c = initCtx(c, "/redis/:key", "key", "test")
	store := new(mocks.Cache)
	store.On("Get", "test").Return(expectedContent, nil).Once()
	c.Set("cache", store)
	h := handlers.NewCacheHandler()

	// Assertions
	if assert.NoError(h.Get(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Equal(expectedContent, rec.Body.String())
		assert.True(store.AssertExpectations(t))
	}
}

func TestHandlerGetError(t *testing.T) {
	assert := assert.New(t)
	const expectedError = "test error"
	const expectedRes = "Can't get value for key 'test': test error"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/redis/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c = initCtx(c, "/redis/:key", "key", "test")
	store := new(mocks.Cache)
	store.On("Get", "test").Return("", fmt.Errorf(expectedError)).Once()
	c.Set("cache", store)
	h := handlers.NewCacheHandler()

	// Assertions
	if assert.NoError(h.Get(c)) {
		assert.Equal(http.StatusBadRequest, rec.Code)
		assert.Equal(expectedRes, rec.Body.String())
		assert.True(store.AssertExpectations(t))
	}
}

func TestHandlerSet(t *testing.T) {
	assert := assert.New(t)
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
	c = initCtx(c, "/redis/:key", "key", "test")
	cache := new(mocks.Cache)
	res := ""
	cache.On("Set", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		res = args[1].(string)
	})
	c.Set("cache", cache)
	h := handlers.NewCacheHandler()

	// Assertions
	if assert.NoError(h.Set(c)) {
		assert.Equal(http.StatusOK, rec.Code)
		assert.Equal(expectedMsg, rec.Body.String())
		assert.Equal(savedContent, res)
		assert.True(cache.AssertExpectations(t))
	}
}

func TestHandlerSetError(t *testing.T) {
	assert := assert.New(t)
	const savedContent = "test content"
	const expectedError = "test error"
	const expectedRes = "Can't set value for key 'test': test error"

	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/redis/test", nil)
	req.Form = url.Values{
		"value": []string{savedContent},
	}
	req.Header = http.Header{"Content-Type": []string{"form-data"}}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c = initCtx(c, "/redis/:key", "key", "test")
	cache := new(mocks.Cache)
	cache.On("Set", mock.Anything, mock.Anything).Return(fmt.Errorf(expectedError)).Once()
	c.Set("cache", cache)
	h := handlers.NewCacheHandler()

	// Assertions
	if assert.NoError(h.Set(c)) {
		assert.Equal(http.StatusBadRequest, rec.Code)
		assert.Equal(expectedRes, rec.Body.String())
		assert.True(cache.AssertExpectations(t))
	}
}
