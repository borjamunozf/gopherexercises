package urlshortener

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	t    *testing.T
	url  string
	want string
}

func TestMapHandler(t *testing.T) {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	mux := DefaultMux()
	mapHandler := MapHandler(pathsToUrls, mux)

	okcases := []TestCase{
		{t: t, url: "/urlshort-godoc", want: "https://godoc.org/github.com/gophercises/urlshort"},
		{t: t, url: "/yaml-godoc", want: "https://godoc.org/gopkg.in/yaml.v2"},
	}

	for _, tt := range okcases {
		tt.t.Run("Test map handler correct for existent paths", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			resp := httptest.NewRecorder()

			mapHandler(resp, req)
			got, err := resp.Result().Location()

			if err != nil {
				t.Fatal("Could not read location")
			}

			if got.String() != tt.want {
				t.Errorf("got %s url redirected but expected %q", got.String(), tt.want)
			}
		})
	}

	t.Run("Test map handler returns fallback for nonexistent path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/quefacemos", nil)
		resp := httptest.NewRecorder()

		mapHandler(resp, req)

		assert.Equal(t, resp.Code, http.StatusOK)
		assert.Equal(t, resp.Body.String(), "Hello, world\n")
	})

}

func TestYmlHandler(t *testing.T) {
	const testYaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
}

func TestParseYaml(t *testing.T) {

	const testYaml = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	want := []UrlPath{
		{Path: "/urlshort", Url: "https://github.com/gophercises/urlshort"},
		{Path: "/urlshort-final", Url: "https://github.com/gophercises/urlshort/tree/solution"},
	}

	data, err := parseYaml([]byte(testYaml))
	if err != nil {
		t.Errorf("could parse yaml data %v", err)
	}

	if !reflect.DeepEqual(data, want) {
		t.Errorf("got %v, want %v", data, want)
	}

}
