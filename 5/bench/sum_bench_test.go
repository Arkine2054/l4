package bench

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/arkine/l4/5/internal/handler"
)

func BenchmarkSumHandler(b *testing.B) {
	body := []byte(`{"a":10,"b":20}`)
	req := httptest.NewRequest(http.MethodPost, "/sum", bytes.NewReader(body))
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		handler.SumHandler(w, req)
	}
}
