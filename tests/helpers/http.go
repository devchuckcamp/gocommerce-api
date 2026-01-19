package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
}

// HTTPTestContext holds context for HTTP testing
type HTTPTestContext struct {
	Router   *gin.Engine
	Recorder *httptest.ResponseRecorder
}

// NewHTTPTestContext creates a new HTTP test context
func NewHTTPTestContext() *HTTPTestContext {
	return &HTTPTestContext{
		Router:   gin.New(),
		Recorder: httptest.NewRecorder(),
	}
}

// Request performs an HTTP request and returns the recorder
func (ctx *HTTPTestContext) Request(method, path string, body interface{}) *httptest.ResponseRecorder {
	ctx.Recorder = httptest.NewRecorder()

	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	ctx.Router.ServeHTTP(ctx.Recorder, req)

	return ctx.Recorder
}

// RequestWithAuth performs an HTTP request with Authorization header
func (ctx *HTTPTestContext) RequestWithAuth(method, path, token string, body interface{}) *httptest.ResponseRecorder {
	ctx.Recorder = httptest.NewRecorder()

	var reqBody io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	ctx.Router.ServeHTTP(ctx.Recorder, req)

	return ctx.Recorder
}

// GET performs a GET request
func (ctx *HTTPTestContext) GET(path string) *httptest.ResponseRecorder {
	return ctx.Request(http.MethodGet, path, nil)
}

// POST performs a POST request
func (ctx *HTTPTestContext) POST(path string, body interface{}) *httptest.ResponseRecorder {
	return ctx.Request(http.MethodPost, path, body)
}

// PUT performs a PUT request
func (ctx *HTTPTestContext) PUT(path string, body interface{}) *httptest.ResponseRecorder {
	return ctx.Request(http.MethodPut, path, body)
}

// PATCH performs a PATCH request
func (ctx *HTTPTestContext) PATCH(path string, body interface{}) *httptest.ResponseRecorder {
	return ctx.Request(http.MethodPatch, path, body)
}

// DELETE performs a DELETE request
func (ctx *HTTPTestContext) DELETE(path string) *httptest.ResponseRecorder {
	return ctx.Request(http.MethodDelete, path, nil)
}

// ParseResponse parses the response body into the given struct
func ParseResponse(t *testing.T, recorder *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if err := json.Unmarshal(recorder.Body.Bytes(), v); err != nil {
		t.Fatalf("Failed to parse response: %v\nBody: %s", err, recorder.Body.String())
	}
}

// AssertStatus asserts the response status code
func AssertStatus(t *testing.T, recorder *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if recorder.Code != expected {
		t.Errorf("Expected status %d, got %d. Body: %s", expected, recorder.Code, recorder.Body.String())
	}
}

// AssertJSON asserts that the response body matches the expected JSON
func AssertJSON(t *testing.T, recorder *httptest.ResponseRecorder, expected map[string]interface{}) {
	t.Helper()
	var actual map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		if !ok {
			t.Errorf("Expected key %q not found in response", key)
			continue
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %q to be %v, got %v", key, expectedValue, actualValue)
		}
	}
}

// CreateTestRouter creates a new Gin router for testing
func CreateTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())
	return router
}

// MockGinContext creates a mock gin.Context for handler testing
func MockGinContext(w http.ResponseWriter, r *http.Request) (*gin.Context, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	c, _ := gin.CreateTestContext(w)
	c.Request = r
	return c, router
}
