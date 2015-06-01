package mix_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cheekybits/is"
	"github.com/matryer/mix"
)

func TestMixHandler(t *testing.T) {
	is := is.New(t)

	h := mix.New("./test/one.js", "./test/*.js")

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/assets/all.js", nil)
	is.NoErr(err)

	h.ServeHTTP(w, r)

	is.Equal(w.Code, http.StatusOK)
	is.True(strings.Contains(w.Body.String(), "one\n"))
	is.True(strings.Contains(w.Body.String(), "two\n"))
	is.True(strings.Contains(w.Body.String(), "three\n"))
	is.Equal(w.HeaderMap.Get("Content-Type"), "application/javascript")
	is.Equal(w.HeaderMap.Get("Content-Length"), "14")

}
