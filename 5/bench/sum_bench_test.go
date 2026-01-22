package bench

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/arkine/l4/5/internal/handler"
)

func BenchmarkSumHandler(b *testing.B) {
	body := []byte(`{"a":1,"b":2}`)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(
			http.MethodPost,
			"/sum",
			bytes.NewReader(body),
		)
		w := httptest.NewRecorder()

		handler.SumHandler(w, req)
	}
}
