// Helper functionalities that aren't exclusively related to web-servers.

package common

import (
    "io"
    "io/ioutil"
    "os"
)

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
