// Helper functionalities that aren't exclusively related to web-servers.

package common

// Errors returned by this package.
type errorCode uint
const (
    // Couldn't create the temporary file
    ErrCreateTempFile errorCode = iota
    // Couldn't write the temporary file
    ErrWriteTempFile
    // Couldn't flush the temporary file
    ErrFlushTempFile
    // Couldn't replace the file atomically
    ErrAtomicReplaceFile
)

// An error associated to its possible cause (which may be nil).
type wrappedErr struct {
    err errorCode
    cause error
}

// Implement the `error` interface for `wrappedErr`.
func (we wrappedErr) Error() string {
    var s string
    switch we.err {
    case ErrCreateTempFile:
        s = "common: Couldn't create the temporary file"
    case ErrWriteTempFile:
        s = "common: Couldn't write the temporary file"
    case ErrFlushTempFile:
        s = "common: Couldn't flush the temporary file"
    case ErrAtomicReplaceFile:
        s = "common: Couldn't replace the file atomically"
    default:
        s = "common: Unknown"
    }

    if we.cause != nil {
        return s + " - From: " + we.cause.Error()
    }
    return s
}

// Wrap a error and its possible cause.
func newError(cause error, err errorCode) error {
    return wrappedErr {
        err: err,
        cause: cause,
    }
}
