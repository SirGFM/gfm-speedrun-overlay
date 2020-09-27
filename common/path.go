// Helper functionalities that aren't exclusively related to web-servers.

package common

import (
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "os"
    "path/filepath"
)

// Append the list of (possibly) local directories `dir` to the
// (optionally) already initialized `list`.
func ListToAbsolutePath(list []string, dirs ...string) ([]string, error) {
    for _, s := range dirs {
        s = filepath.Clean(s)
        s, err := filepath.Abs(s)
        if err != nil {
            return list, newError(err, ErrGetAbsPath)
        }

        fi, err := os.Stat(s)
        if err != nil {
            return list, newError(err, ErrStatFile)
        } else if !fi.IsDir() {
            return list, newError(err, ErrNotDir)
        }

        list = append(list, s)
    }

    return list, nil
}

// Open the requested file and retrieve its content type, based on the
// file's extension. If the extension can't be deduced, it defaults to
// "text/plain".
func OpenFile(dir, filePath string) (*os.File, string, error) {
    filePath = filepath.Join(dir, filePath)

    file, err := os.Open(filePath)
    if err != nil {
        // Don't do any further post-processing
        err = newError(err, ErrOpenFile)
    } else if fi, err := file.Stat(); err != nil {
        err = newError(err, ErrStatFile)
    } else if fi.IsDir() {
        err = newError(err, ErrDir)
    }

    if err != nil {
        if we, ok := err.(wrappedErr); ok && we.err != ErrOpenFile {
            file.Close()
        }
        return nil, "", err
    }

    var ctype string
    switch ext := filepath.Ext(filePath); ext {
    case ".json":
        ctype = "application/json"
    case ".css":
        ctype = "text/css"
    case ".html":
        ctype = "text/html"
    case ".js":
        ctype = "text/javascript"
    case ".bmp":
        ctype = "image/bmp"
    case ".gif":
        ctype = "image/gif"
    case ".ico":
        ctype = "image/x-icon"
    case ".jpeg", ".jpg":
        ctype = "image/jpeg"
    case ".png":
        ctype = "image/png"
    case ".svg":
        ctype = "image/svg+xml"
    default:
        logger.Warnf("common: Unknown file extension: %+s", ext)
        ctype = "text/plain"
    }

    return file, ctype, nil
}
