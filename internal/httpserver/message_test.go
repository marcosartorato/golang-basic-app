package httpserver_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/marcosartorato/myapp/internal/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestMessageHandlerRepeatType(t *testing.T) {
	// Skip the test if the env var is set
	if os.Getenv("SKIP_HANDLER_TEST") == "true" {
		t.Skip("skipping test about handlers")
	}

	// Define support struct
	type repeatRequest struct {
		Type string `json:"type"`
		Msg  string `json:"msg"`
	}
	type repeatResponse struct {
		Type string `json:"type"`
		Msg  string `json:"msg"`
	}

	// Define test cases
	tests := map[string]struct {
		reqBody    interface{}
		wantStatus int
		wantMsg    string
	}{
		"repeat type: normal message": {
			reqBody:    repeatRequest{Type: "repeat", Msg: "hello"},
			wantStatus: http.StatusOK,
			wantMsg:    "hello",
		},
		"repeat type: empty message": {
			reqBody:    repeatRequest{Type: "repeat", Msg: ""},
			wantStatus: http.StatusOK,
			wantMsg:    "",
		},
		"repeat type: large message": {
			reqBody:    repeatRequest{Type: "repeat", Msg: strings.Repeat("x", 10_000)},
			wantStatus: http.StatusOK,
			wantMsg:    strings.Repeat("x", 10_000),
		},
		"repeat type: extra field ignored": {
			reqBody:    `{"type":"repeat","msg":"hey","extra":"ignored"}`,
			wantStatus: http.StatusOK,
			wantMsg:    "hey",
		},
		"repeat type: invalid JSON": {
			reqBody:    `{"type":"repeat","msg":"oops"`, // missing closing brace
			wantStatus: http.StatusBadRequest,
		},
	}

	// Run tests
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var bodyBytes []byte
			switch v := tc.reqBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				var err error
				bodyBytes, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal req: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/message", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			httpserver.MessageHandler(rec, req)
			res := rec.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			assert.Equal(t, tc.wantStatus, res.StatusCode, "unexpected status code")
			if res.StatusCode == http.StatusOK {
				var got repeatResponse
				if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}
				assert.Equal(t, "repeat", got.Type, "unexpected type in the response body")
				assert.Equal(t, tc.wantMsg, got.Msg, "unexpected msg field in the response body")
			}
		})
	}
}

func TestMessageHandlerTimeType(t *testing.T) {
	// Skip the test if the env var is set
	if os.Getenv("SKIP_HANDLER_TEST") == "true" {
		t.Skip("skipping test about handlers")
	}

	req := httptest.NewRequest(http.MethodPost, "/api/message", strings.NewReader(`{"type":"time"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	httpserver.MessageHandler(rec, req)

	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// Check status code
	assert.Equal(t, res.StatusCode, http.StatusOK, "unexpected status code")

	// Check body
	type timeResp struct {
		Type string `json:"type"`
		Time string `json:"time"`
	}
	got := timeResp{}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal time: %v", err)
	}
	assert.Equal(t, "time", got.Type, "unexpected type in the response body")
	if _, err := time.Parse(time.RFC3339, got.Time); err != nil {
		t.Fatalf("time not RFC3339: %q (%v)", got.Time, err)
	}
}

func TestMessageHandlerUnknownType(t *testing.T) {
	// Skip the test if the env var is set
	if os.Getenv("SKIP_HANDLER_TEST") == "true" {
		t.Skip("skipping test about handlers")
	}

	/*
		Unknown type
	*/
	req := httptest.NewRequest(http.MethodPost, "/api/message", strings.NewReader(`{"type":"nope"}`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	httpserver.MessageHandler(rec, req)

	res := rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// Check status code
	assert.Equal(t, res.StatusCode, http.StatusBadRequest, "unexpected status code")

	// Check body
	if !strings.Contains(rec.Body.String(), "unknown type") {
		t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), "unknown type")
	}

	/*
		Invalid JSON
	*/
	req = httptest.NewRequest(http.MethodPost, "/api/message", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	httpserver.MessageHandler(rec, req)

	res = rec.Result()
	defer func() {
		_ = res.Body.Close()
	}()

	// Check status code
	assert.Equal(t, res.StatusCode, http.StatusBadRequest, "unexpected status code")

	// Check body
	if !strings.Contains(rec.Body.String(), "invalid JSON") {
		t.Fatalf("body = %q, want it to contain %q", rec.Body.String(), "invalid JSON")
	}
}
