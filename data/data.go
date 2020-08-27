package data

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Fetch gets the bytes at the specified URI.
// The URI can be remote (http) or local.
// if 'limit' is greater then zero, fetch stops
// with EOF after 'limit' bytes.
func Fetch(uri string, limit int64) ([]byte, error) {
	if strings.HasPrefix(uri, "http") {
		return FetchFromURI(uri, limit)
	}

	return FetchFromFile(uri, limit)
}

// FetchFromURI fetch data (with limit) from an HTTP URL.
// if 'limit' is greater then zero, fetch stops
// with EOF after 'limit' bytes.
func FetchFromURI(uri string, limit int64) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if limit > 0 {
		return ioutil.ReadAll(io.LimitReader(res.Body, limit))
	}

	return ioutil.ReadAll(res.Body)
}

// FetchFromFile fetch data (with limit) from an file.
// if 'limit' is greater then zero, fetch stops
// with EOF after 'limit' bytes.
func FetchFromFile(filename string, limit int64) ([]byte, error) {
	fp, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	if limit > 0 {
		return ioutil.ReadAll(io.LimitReader(fp, limit))
	}

	return ioutil.ReadAll(fp)
}

// Wrap hard wrap text at the specified colBreak column.
func Wrap(text string, colBreak int) string {
	if colBreak < 1 {
		return text
	}
	text = strings.TrimSpace(text)

	var sb strings.Builder
	var i int
	for i = 0; len(text[i:]) > colBreak; i += colBreak {
		sb.WriteString(text[i : i+colBreak])
		sb.WriteString("\n")

	}
	sb.WriteString(text[i:])

	return sb.String()
}
