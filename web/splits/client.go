// `splits` store splits for games (i.e., a list of objectives within a run
// of the game). Another module should be used to time runs, as this only
// manipulates the structure of splits.
//
// See `splits.go` for the full description.

package splits

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "path"
)

// Retrieve the list of splits for a given game/category, referenced as
// `name` in the service (hosted at `address:port`). On error, the error
// shall be properly wrapped into a `server.common.HttpError`.
func GetSplits(name, address string, port int) ([]string, error) {
    prefix := Prefix
    if prefix[0] == '/' {
        prefix = prefix[1:]
    }
    prefix = url.PathEscape(prefix)

    game := url.PathEscape(name)
    res := path.Clean(path.Join("/", prefix, "load", game))
    addr := fmt.Sprintf("http://%s:%d%s", address, port, res)

    req, err := http.NewRequest("GET", addr, nil)
    if err != nil {
        reason := "Failed to prepare the split request"
        code := http.StatusInternalServerError
        return nil, newError(err, reason, code)
    }
    req.Header.Add("Content-Type", "application/json")

    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        reason := "Failed to retrieve the requested split"
        code := http.StatusServiceUnavailable
        return nil, newError(err, reason, code)
    }
    defer resp.Body.Close()
    if code := resp.StatusCode; code != http.StatusOK {
        reason := "Bad response getting the split"
        return nil, newError(nil, reason, code)
    }

    var sp splits
    dec := json.NewDecoder(resp.Body)
    err = dec.Decode(&sp)
    if err != nil {
        reason := "Failed to decode the splits response"
        code := http.StatusInternalServerError
        return nil, newError(err, reason, code)
    }

    return sp.Entries, nil
}
