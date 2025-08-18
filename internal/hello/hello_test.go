package hello_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marcosartorato/myapp/internal/hello"

	"github.com/stretchr/testify/assert"
)

// Basic handler test: returns 200 and "Hello, World!\n"
func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()
	hello.HelloHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "unexpected status code")

	body, _ := io.ReadAll(rr.Body)
	assert.Equal(t, "Hello, World!\n", string(body), "unexpected response body")
}
