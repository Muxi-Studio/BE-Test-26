package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"blog-backend/internal/store"
)

func serve(h *Handler, method, target string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	h.Register(mux)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestListTodosKeepsTitle(t *testing.T) {
	h := NewHandler(store.NewMemory())
	h.store.Add("写题解")

	rec := serve(h, "GET", "/todos")
	if rec.Code != http.StatusOK {
		t.Fatalf("状态码 = %d", rec.Code)
	}

	var body struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("条数 = %d", len(body.Data))
	}
	if _, ok := body.Data[0]["title"]; !ok {
		t.Fatalf("标题字段丢失: %v", body.Data[0])
	}
}

func TestErrorDoesNotLeakDetail(t *testing.T) {
	rec := httptest.NewRecorder()
	writeError(rec, errors.New("dial tcp 10.0.0.1:5432: connection refused"))

	body := rec.Body.String()
	if strings.Contains(body, "connection refused") {
		t.Fatalf("内部错误泄漏到响应: %s", body)
	}
	if !strings.Contains(body, "internal error") {
		t.Fatalf("错误信息不符: %s", body)
	}
}
