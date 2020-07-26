// Interfaces used by the HTTP server. This separation was mainly to avoid
// possible circular inclusion issues.
//
// This package also has a few utility functions used by multiple servers.

package common

import (
    "fmt"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "net/http"
    "runtime/debug"
    "strings"
)

// Error type used throughout this repository.
type HttpError interface {
    // Whether this object is a httpError (silly, I know...).
    IsHttpError() bool
    // Retrieve the error code sent in the reply
    GetHttpStatus() string
    // Implement Go's error interface.
    error
}

// Hold the context for a given HTTP error.
type httpError struct {
    // Inner error that triggered this response.
    inner error
    // Module where the error happened.
    module string
    // Human-readable error reason, sent as the reply.
    reason string
    // Error code sent as the reply.
    httpStatus int
    // Stack trace for the error. If the error is wrapping another
    // httpError, the stack trace will be left empty.
    stackTrace []byte
    // Number of errors wrapped by this object.
    wrapCount int
}

// Create a new HTTPError.
func NewHttpError(inner error, module string, reason string, httpStatus int) HttpError {
    var stackTrace []byte

    if he, ok := inner.(*httpError); !ok {
        stackTrace = debug.Stack()
    } else if he != nil {
        errCopy := *he
        errCopy.setWrapCount(1)

        inner = &errCopy
    }

    return &httpError {
        inner: inner,
        module: module,
        reason: reason,
        httpStatus: httpStatus,
        stackTrace: stackTrace,
        wrapCount: 0,
    }
}

// Whether this object is a httpError (silly, I know...).
func (e *httpError) IsHttpError() bool {
    return true
}

// Retrieve the error code sent in the reply
func (e *httpError) GetHttpStatus() string {
    return http.StatusText(e.httpStatus)
}

// Set the wrap count for this error and every wrapped httpError.
func (e *httpError) setWrapCount(count int) {
    if he, ok := e.inner.(*httpError); ok && he != nil {
        he.setWrapCount(count + 1)
    }
    e.wrapCount = count
}

// Implement the error interface for httpError.
func (e *httpError) Error() string {
    var buf strings.Builder
    var tabsBuf strings.Builder
    var args []interface{}

    for i := e.wrapCount; i > 0; i-- {
        tabsBuf.WriteString("\t")
    }
    tabs := tabsBuf.String()

    buf.WriteString("%s%s: %s (%d)")
    args = append(args, tabs, e.module, e.reason, e.httpStatus)

    if e.inner != nil {
        buf.WriteString("\n%s%s:")
        args = append(args, tabs, "Base error")

        if he, ok := e.inner.(*httpError); ok && he != nil {
            buf.WriteString("\n%+v")
            args = append(args, he)
        } else {
            buf.WriteString(" %+v")
            args = append(args, e.inner)
        }
    }

    if e.stackTrace != nil {
        buf.WriteString("\n%s:\n\n%+v\n")
        args = append(args, "Stack trace", string(e.stackTrace))
    }

    str := buf.String()
    return fmt.Sprintf(str, args...)
}

// Send a status code response.
func ReplyStatus(mod string, status int, msg string, w http.ResponseWriter) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(status)

    for data := []byte(msg); len(data) > 0; {
        n, err := w.Write(data)
        if err != nil {
            logger.Errorf("%s: Failed to send %d: %+v", mod, err, status)
            return
        }
        data = data[n:]
    }
}

// Reply a message with a `httpError`
func ReplyHttpError(e HttpError, w http.ResponseWriter) {
    fmt.Printf("%+v", e)
    if he, ok := e.(*httpError); ok && he != nil {
        ReplyStatus(he.module, he.httpStatus, he.reason, w)
    }
}
