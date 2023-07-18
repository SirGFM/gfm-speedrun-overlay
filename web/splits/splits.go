// `splits` store splits for games (i.e., a list of objectives within a run
// of the game). Another module should be used to time runs, as this only
// manipulates the structure of splits.
//
// The service accepts four HTTP methods: GET, POST, PUT and DELETE.
//
// ## GET
//
// GET requests should be used to retrieve/list data stored in the server.
// These requests do not expect any data in the request body.
//
// Sending a GET with the path `list` (e.g.,
// http://localhost:8080/splits/list) will retrieve the list of splits
// stored in the server. The server willww reply with a JSON in the format:
//
//     {
//         "Splits": [
//             "game 0",
//             "game 1",
//             // ...
//             "game n",
//         ]
//     }
//
// To retrieve a specific entry, the custom path `load/<split-name>` must
// be used (e.g., http://localhost:8080/splits/load/my-game). The server
// will reply with a JSON in the format:
//
//     {
//         "Name": "my-game",
//         "Entries": [
//             "entry 0",
//             "entry 1",
//             // ...
//             "entry n"
//         ]
//     }
//
// Alternatively, the public function `GetSplits()` may be used to retrieve
// a specific entry programatically.
//
// ## POST & PUT
//
// POST requests must be used to add new splits. PUT, on the other hand,
// must be used to modify existing ones. The data for both operations
// should be encoded in the payload, as following JSON:
//
//     {
//         "Name": "my-game",
//         "Entries": [
//             "entry 0",
//             "entry 1",
//             // ...
//             "entry n"
//         ]
//     }
//
// ## DELETE
//
// DELETE removes the resource from the server. This method does not
// expect any data in the request body, and the split must be specified in
// the address (e.g., http://localhost:8080/splits/my-game).

package splits

import (
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/common"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
    "path"
    "strings"
    "sync"
)

const Prefix = "/splits"

// Context for the splits service
type splitsCtx struct {
    baseDir string
    // Synchronize access to the context
    rwmut sync.RWMutex
}

// Object returned by 'list' operations
type listResp struct {
    Splits []string
}

// Object used in most operations
type splits struct {
    Name string
    Entries []string `json:",omitempty"`
}

// Build a new error
func newError(err error, res string, status int) error {
    return srv_iface.NewHttpError(err, "web"+Prefix, res, status)
}

// Retrieve the path handled by `splits`.
func (*splitsCtx) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*splitsCtx) Dependencies() []string {
    return nil
}

// Close resources associated with the `splits` (i.e, nothing)
func (*splitsCtx) Close() {
}

// Convert a name, as supplied in an URL, into a local file path.
func (ctx *splitsCtx) getFileName(name string) string {
    // QueryEscape converts space to '+',
    // while PathEscape doesn't escape every character...
    // Achieve a middle ground by splitting every space-separated fragment,
    // and then manually joining those.
    var fragments []string
    for _, fragment := range strings.Split(name, " ") {
        fragments = append(fragments, url.QueryEscape(fragment))
    }
    safeName := strings.Join(fragments, url.PathEscape(" "))

    return path.Join(ctx.baseDir, safeName+".json")
}

// Check whether or not a file exists.
// NOTE: the caller must synchronize access to the file!
func (ctx *splitsCtx) unsafeFileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false, nil
    } else if err != nil {
        return false, newError(err, "Failed to open the requested file", http.StatusInternalServerError)
    }
    return true, nil
}

// Save splits information to a file, atomically.
// NOTE: the caller must synchronize access to the file!
func (ctx *splitsCtx) unsafeSaveSplit(sp splits) error {
    writefn := func(w io.Writer) error {
        enc := json.NewEncoder(w)
        err := enc.Encode(&sp)
        if err != nil {
            return newError(err, "Failed to encode the splits file", http.StatusInternalServerError)
        }
        return nil
    }

    filename := ctx.getFileName(sp.Name)
    err := common.AtomicSaveFile(ctx.baseDir, filename, writefn)
    if err != nil {
        return newError(err, "Failed to create the splits file", http.StatusInternalServerError)
    }
    return nil
}

// Handle 'list' operations.
func (ctx *splitsCtx) listSplits() ([]string, error) {
    ctx.rwmut.RLock();
    defer ctx.rwmut.RUnlock();

    fis, err := ioutil.ReadDir(ctx.baseDir)
    if err != nil {
        return nil, newError(err, "Failed to retrieve the stored splits", http.StatusInternalServerError)
    }

    names := make([]string, 0)
    for i := range fis {
        if fis[i].IsDir() || path.Ext(fis[i].Name()) != ".json" {
            continue
        }
        name := fis[i].Name()
        name = name[:len(name) - len(".json")]
        name, err := url.PathUnescape(name)
        if err != nil {
            return nil, newError(err, "Failed to retrieve name of stored split", http.StatusInternalServerError)
        }
        names = append(names, name)
    }

    return names, nil
}

// Handle 'get' operations.
func (ctx *splitsCtx) getSplits(name string) (splits, error) {
    ctx.rwmut.RLock();
    defer ctx.rwmut.RUnlock();

    fpath := ctx.getFileName(name)
    hasFile, err := ctx.unsafeFileExists(fpath)
    if err != nil {
        return splits{}, err
    } else if !hasFile {
        return splits{}, newError(err, "Splits does not exist!", http.StatusNotFound)
    }

    fp, err := os.Open(fpath)
    if err != nil {
        return splits{}, newError(err, "Failed to retrieve the requested splits", http.StatusInternalServerError)
    }
    defer fp.Close()

    var sp splits
    dec := json.NewDecoder(fp)
    err = dec.Decode(&sp)
    if err != nil {
        err = newError(err, "Failed to decode the requested splits", http.StatusInternalServerError)
    }
    return sp, err
}

// Handle 'add' operations.
func (ctx *splitsCtx) addSplits(sp splits) error {
    ctx.rwmut.Lock();
    defer ctx.rwmut.Unlock();

    fpath := ctx.getFileName(sp.Name)
    hasFile, err := ctx.unsafeFileExists(fpath)
    if err != nil {
        return err
    } else if hasFile {
        return newError(err, "Splits already exist! Must use 'PUT'!", http.StatusBadRequest)
    }

    return ctx.unsafeSaveSplit(sp)
}

// Handle 'update' operations.
func (ctx *splitsCtx) updateSplits(sp splits) error {
    ctx.rwmut.Lock();
    defer ctx.rwmut.Unlock();

    fpath := ctx.getFileName(sp.Name)
    hasFile, err := ctx.unsafeFileExists(fpath)
    if err != nil {
        return err
    } else if !hasFile {
        return newError(err, "Splits does not exist! Must use 'POST'!", http.StatusBadRequest)
    }

    err = ctx.unsafeSaveSplit(sp)
    if err != nil {
        return err
    }

    return nil
}

// Handle 'delete' operations.
func (ctx *splitsCtx) delSplits(name string) error {
    ctx.rwmut.Lock();
    defer ctx.rwmut.Unlock();

    fpath := ctx.getFileName(name)
    hasFile, err := ctx.unsafeFileExists(fpath)
    if err != nil {
        return err
    } else if !hasFile {
        return newError(err, "Splits does not exist!", http.StatusNotFound)
    }

    err = os.Remove(fpath)
    if err != nil {
        return newError(err, "Couldn't remove the splits", http.StatusInternalServerError)
    }

    return nil
}

// Handle GET requests.
func (ctx *splitsCtx) get(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) == 0 {
        return newError(nil, "Missing command", http.StatusBadRequest)
    }

    switch urlPath[0] {
    case "load":
        if len(urlPath) != 2 {
            return newError(nil, "Splits name missing or too many arguments", http.StatusBadRequest)
        }

        resp, err := ctx.getSplits(urlPath[1])
        if err != nil {
            return err
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        enc := json.NewEncoder(w)
        err = enc.Encode(&resp)
        if err != nil {
            logger.Errorf("web%s: Failed to encode the responde: %+v (payload: %+v)", Prefix, err, resp)
        }
    case "list":
        if len(urlPath) != 1 {
            return newError(nil, "Too many arguments", http.StatusBadRequest)
        }

        var resp listResp
        var err error
        resp.Splits, err = ctx.listSplits()
        if err != nil {
            return err
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        enc := json.NewEncoder(w)
        err = enc.Encode(&resp)
        if err != nil {
            logger.Errorf("web%s: Failed to encode the responde: %+v (payload: %+v)", Prefix, err, resp)
        }
    default:
        return newError(nil, "Invalid operation", http.StatusBadRequest)
    }

    return nil
}

// Handle POST request.
func (ctx *splitsCtx) post(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) != 0 {
        return newError(nil, "Too many arguments", http.StatusBadRequest)
    }

    var sp splits
    dec := json.NewDecoder(req.Body)
    err := dec.Decode(&sp)
    if err != nil {
        return newError(err, "Failed to decode the received splits", http.StatusBadRequest)
    }

    err = ctx.addSplits(sp)
    if err != nil {
        return err
    }
    w.WriteHeader(http.StatusNoContent)

    return nil
}

// Handle PUT request.
func (ctx *splitsCtx) put(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) != 0 {
        return newError(nil, "Too many arguments", http.StatusBadRequest)
    }

    var sp splits
    dec := json.NewDecoder(req.Body)
    err := dec.Decode(&sp)
    if err != nil {
        return newError(err, "Failed to decode the received splits", http.StatusBadRequest)
    }

    err = ctx.updateSplits(sp)
    if err != nil {
        return err
    }
    w.WriteHeader(http.StatusNoContent)

    return nil
}

// Handle DELETE request.
func (ctx *splitsCtx) del(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) != 1 {
        return newError(nil, "Splits name missing or too many arguments", http.StatusBadRequest)
    }

    err := ctx.delSplits(urlPath[0])
    if err != nil {
        return err
    }
    w.WriteHeader(http.StatusNoContent)

    return nil
}

// Handle requests to the `splits` service, filtering and redirecting as
// necessary.
func (ctx *splitsCtx) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if req.Header.Get("Content-Type") != "application/json" {
        reason := "Content-Type must be \"application/json\""
        return newError(nil, reason, http.StatusUnsupportedMediaType)
    }

    var err error

    urlPath = urlPath[1:]
    switch req.Method {
    case "GET":
        return ctx.get(w, req, urlPath)
    case "POST":
        return ctx.post(w, req, urlPath)
    case "PUT":
        return ctx.put(w, req, urlPath)
    case "DELETE":
        return ctx.del(w, req, urlPath)
    default:
        return newError(err, "Invalid method: wanted either GET, POST, PUT or DELETE", http.StatusMethodNotAllowed)
    }
}

// Register a `splits` handler in the `Server`.
func GetHandle(srv srv_iface.Server, baseDir string) error {
    var ctx splitsCtx

    ctx.baseDir = path.Clean(baseDir)
    srv.AddHandler(&ctx)
    return nil
}
