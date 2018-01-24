package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gwleclerc/dummy-golang-test/cache/mocks"
	"github.com/gwleclerc/dummy-golang-test/handlers"
	"github.com/stretchr/testify/mock"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func initCtx(c echo.Context, path, key, value string) echo.Context {
	c.SetPath(path)
	c.SetParamNames(key)
	c.SetParamValues(value)
	return c
}

func TestCacheHandler(t *testing.T) {
	Convey("Using Cache Handler", t, func() {

		var (
			engine    *echo.Echo
			recorder  *httptest.ResponseRecorder
			context   echo.Context
			cacheMock *mocks.Cache
			handler   handlers.CacheHandler
			err       error
		)

		// Init echo
		engine = echo.New()
		recorder = httptest.NewRecorder()
		handler = handlers.NewCacheHandler()

		// CacheHandler.Get
		Convey("When I try to Get 'test' key", func() {

			//Init
			var (
				request *http.Request
			)
			request = httptest.NewRequest(echo.GET, "/redis/test", nil)
			context = engine.NewContext(request, recorder)
			context = initCtx(context, "/redis/:key", "key", "test")
			cacheMock = new(mocks.Cache)
			context.Set("cache", cacheMock)

			// straight case
			Convey("with existing value in cache", func() {
				const expectedContent = "test content"
				cacheMock.On("Get", "test").Return(expectedContent, nil).Once()
				err = handler.Get(context)

				So(err, ShouldBeNil)
				So(recorder.Code, ShouldEqual, http.StatusOK)
				So(recorder.Body.String(), ShouldEqual, expectedContent)
				So(cacheMock.AssertExpectations(t), ShouldBeTrue)
			})

			// missing value case
			Convey("with missing value in cache", func() {
				const expectedError = "test error"
				const expectedBody = "Can't get value for key 'test': " + expectedError
				cacheMock.On("Get", "test").Return("", fmt.Errorf(expectedError)).Once()
				err = handler.Get(context)

				So(err, ShouldBeNil)
				So(recorder.Code, ShouldEqual, http.StatusBadRequest)
				So(recorder.Body.String(), ShouldEqual, expectedBody)
				So(cacheMock.AssertExpectations(t), ShouldBeTrue)
			})

		})

		// CacheHandler.Set
		Convey("When I try to set value for 'test' key", func() {

			//Init
			const savedContent = "test content"
			var (
				request *http.Request
			)
			request = httptest.NewRequest(echo.POST, "/redis/test", nil)
			request.Form = url.Values{
				"value": []string{savedContent},
			}
			request.Header = http.Header{"Content-Type": []string{"form-data"}}
			context = engine.NewContext(request, recorder)
			context = initCtx(context, "/redis/:key", "key", "test")
			cacheMock = new(mocks.Cache)
			context.Set("cache", cacheMock)

			// straight case
			Convey("without error", func() {
				const expectedMsg = "Value stored successfully"
				const expectedContent = "test content"
				var res string
				cacheMock.On("Set", mock.Anything, mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
					res = args[1].(string)
				})
				err = handler.Set(context)

				So(err, ShouldBeNil)
				So(recorder.Code, ShouldEqual, http.StatusOK)
				So(recorder.Body.String(), ShouldEqual, expectedMsg)
				So(res, ShouldEqual, expectedContent)
				So(cacheMock.AssertExpectations(t), ShouldBeTrue)
			})

			// cache error case
			Convey("with cache error", func() {
				const expectedError = "test error"
				const expectedRes = "Can't set value for key 'test': " + expectedError
				cacheMock.On("Set", mock.Anything, mock.Anything).Return(fmt.Errorf(expectedError)).Once()
				err = handler.Set(context)

				So(err, ShouldBeNil)
				So(recorder.Code, ShouldEqual, http.StatusBadRequest)
				So(recorder.Body.String(), ShouldEqual, expectedRes)
				So(cacheMock.AssertExpectations(t), ShouldBeTrue)
			})

		})
	})
}
