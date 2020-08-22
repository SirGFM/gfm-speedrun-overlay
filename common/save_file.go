// Helper functionalities that aren't exclusively related to web-servers.

package common

import (
    "io"
    "io/ioutil"
    "os"
)

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

// Interface used by `AtomicSaveFile` to write data to the file.
type writeCb (func(io.Writer) error)

// Atomically write file `filename`, by calling `writefn`
// (a `func(io.Writer) error`) on a temporary file and then renaming this
// temporary file into the desired one. The temporary file is created on
// `tmpdir`, which must already exist. Note that on some systems, `tmpdir`
// and `filename` **must be within the same partition**.
// On failure, the temporary file is properly removed and `filename`, if
// existant, is kept intact.
func AtomicSaveFile(tmpdir, filename string, writefn writeCb) error {
    // Write the contents to a temporary file, so errors do not corrupt the file
    tmp, err := ioutil.TempFile(tmpdir, "atomic_file*.json")
    if err != nil {
        return newError(err, ErrCreateTempFile)
    }
    defer func() {
        // On error, try to `Close()` and `Remove()` the file. If `Close()`
        // has already been called, it simply returns an error!
        if tmp != nil {
            tmp.Close()
            os.Remove(tmp.Name())
        }
    } ()

    err = writefn(tmp)
    if err != nil {
        return newError(err, ErrWriteTempFile)
    }

    // Rename the file, atomically replacing `filename`
    err = tmp.Close()
    if err != nil {
        return newError(err, ErrFlushTempFile)
    }
    err = os.Rename(tmp.Name(), filename)
    if err != nil {
        return newError(err, ErrAtomicReplaceFile)
    }
    tmp = nil

    return nil
}
