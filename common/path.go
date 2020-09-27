// Helper functionalities that aren't exclusively related to web-servers.

package common

import (
    "path/filepath"
    "os"
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
