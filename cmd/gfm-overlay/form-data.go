package main

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
)

// MultipartReader parses a multipart form data.
// Copied from Go 1.20's source code.
func MultipartReader(header http.Header, body io.Reader, allowMixed bool) (*multipart.Reader, error) {
	v := header.Get("Content-Type")
	if v == "" {
		return nil, errors.New("missing content type")
	}

	if body == nil {
		return nil, errors.New("missing form body")
	}

	d, params, err := mime.ParseMediaType(v)
	if err != nil || !(d == "multipart/form-data" || allowMixed && d == "multipart/mixed") {
		return nil, errors.New("invalid media type")
	}

	boundary, ok := params["boundary"]
	if !ok {
		return nil, errors.New("missing boundary")
	}

	return multipart.NewReader(body, boundary), nil
}
