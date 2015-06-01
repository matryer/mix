package mix

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Handler is a http.Handler that mixes files.
type Handler struct {
	files []string
	err   error

	// Header allows you to set common response headers that will be
	// send with requests handled by this Handler.
	// Use ClearHeaders() to reset default headers.
	Header http.Header
}

var _ http.Handler = (*Handler)(nil)

// New makes a new mix handler with the specified files or
// patterns.
// By default, the following HTTP headers will be included:
//     X-Mix-Patterns - comma separated list of patterns
//     X-Mix-Files - comma separated list of matching files
func New(patterns ...string) *Handler {
	files, err := Glob(patterns...)
	h := (&Handler{
		files: files,
		err:   err,
	}).ClearHeaders()
	h.Header.Set("X-Mix-Patterns", strings.Join(patterns, ", "))
	h.Header.Set("X-Mix-Files", strings.Join(files, ", "))
	return h
}

// ClearHeaders clears the X-Mix* headers.
func (h *Handler) ClearHeaders() *Handler {
	h.Header = make(http.Header)
}

// ServeFiles serves all specified files.
// Content-Type (if not set) will be inferred from the extension in the
// request.
// Uses http.ServeContent to serve the content.
func ServeFiles(w http.ResponseWriter, r *http.Request, files ...string) {

	var recentMod time.Time
	var buf bytes.Buffer
	for _, f := range files {

		stat, err := os.Stat(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// keep track of latest modtime
		if stat.ModTime().After(recentMod) {
			recentMod = stat.ModTime()
		}

		file, err := os.Open(f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(&buf, file)
		file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write linefeed
		if _, err := buf.WriteRune('\n'); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	http.ServeContent(w, r, path.Base(r.URL.Path), recentMod, sizable(buf))

}

// ServeHTTP serves the request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// set headers
	for k, vs := range h.Header {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	// return error if something went wrong
	if h.err != nil {
		http.Error(w, h.err.Error(), http.StatusInternalServerError)
		return
	}

	ServeFiles(w, r, h.files...)

}

// Glob takes a range of patterns and produces a unique list
// of matching files.
// Files are added in pattern and alphabetical order.
// Like filepath.Glob, but you can pass in many patterns.
func Glob(patterns ...string) ([]string, error) {
	seen := make(map[string]struct{})
	var files []string
	for _, g := range patterns {
		matches, err := filepath.Glob(g)
		if err != nil {
			return nil, err
		}
		for _, match := range matches {
			match = filepath.Clean(match)
			if _, alreadySeen := seen[match]; !alreadySeen {
				files = append(files, match)
				seen[match] = struct{}{}
			}
		}
	}
	return files, nil
}

// sizableBuffer is a wrapper around a bytes.Buffer that allows
// http.ServeContent to get the content length.
// Buffers can't normally seek, so this just simulates the behaviour
// and returns buf.Len() when os.SEEK_END is requested.
type sizableBuffer struct {
	buf bytes.Buffer
}

var _ io.ReadSeeker = (*sizableBuffer)(nil)

func (s *sizableBuffer) Seek(offset int64, whence int) (int64, error) {
	if whence == os.SEEK_END {
		return int64(s.buf.Len()), nil
	}
	return 0, nil
}

func (s *sizableBuffer) Read(p []byte) (int, error) {
	return s.buf.Read(p)
}

func sizable(buf bytes.Buffer) *sizableBuffer {
	return &sizableBuffer{buf: buf}
}
