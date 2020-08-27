package data

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchFromURI(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from scrawl!")
	}))
	defer ts.Close()

	data, err := FetchFromURI(ts.URL, 1024*10)
	if err != nil {
		t.Error(err)
	}

	want := "Hello from scrawl!"
	if got := string(data); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFetchFromFile(t *testing.T) {
	data, err := FetchFromFile("../testdata/spritesheet.json", 10)
	if err != nil {
		t.Error(err)
	}

	want := `{   "fram`
	if got := flatten(string(data)); got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

// remove tabs and newlines and spaces
func flatten(s string) string {
	return strings.Replace((strings.Replace(s, "\n", "", -1)), "\t", "", -1)
}
