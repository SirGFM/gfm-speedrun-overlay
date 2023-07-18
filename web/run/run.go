// `run` track individual runs of a given game/category. It depends on a
// `splits` service for retrieving the splits (i.e., intermediate
// segments/goals) for the game/category from a `splits` service.
//
// The service accepts two HTTP methods: GET, POST.
//
// ## GET
//
// GET requests should be used to initiate a new run and to retrieve a
// run's status. When a new run is created, the service creates a unique
// token for controlling the run. After that, you may retrieve either the
// run's splits or its current timer.
//
// ### Initializing a new run
//
// To initialize a new run, first register the split in the `splits`
// service and then send a HTTP GET request to the `run` service with the
// path `new/<split-name>`. For example, to initialize a run for the
// game/category "JJAT (any%)", you would send an empty request to
// `http://localhost:8080/run/new/JJAT%20(any%25)`.
//
// The service replies the following JSON object containing the token:
//
//     {
//         "Token": "Some-random-token"
//     }
//
// ### Retrieving a run's status
//
// A run status is divided into two parts: its current time and its splits.
// Each may be accessed individually in a different path.
//
// To retrieve the run's time, send a HTTP GET request to the `run` service
// with the path `timer/<token>`. The service replies with a JSON object
// containing the current run time in milliseconds in the Time field,
// regardless of whether the timer is running or is stopped, as
// exemplified bellow:
//
//     {
//         "Time": 0
//     }
//
// A run's splits track the progress of a run through each of the segment
// in a game/category, as compared to the fastest run completition time. To
// retrieve the run's splits, send a HTTP GET request to the `run` service
// with the path `splits/<token>`. The service replies with a JSON object
// that contains information on the currently active split, the
// completition time of completed segments and the aimed time for segments
// to come, as described bellow:
//
//     {
//         "Name": "Some-game/category",
//         "Splits": [
//             {
//                 "Name": "Segment 1",
//                 "BestTime": 60000,
//                 "StartTime": 0,
//                 "EndTime": 62000,
//                 "Skipped": false,
//             },
//             {
//                 "Name": "...",
//                 "BestTime": 180000,
//                 "StartTime": 62000,
//                 "EndTime": 0,
//                 "Skipped": false,
//             },
//             {
//                 "Name": "Boss",
//                 "BestTime": 30000,
//                 "StartTime": 0,
//                 "EndTime": 0,
//                 "Skipped": false,
//             }
//         ],
//         "Best": [
//             // Same structure as Splits
//         ],
//         "Current": 1,
//     }
//
// ## POST
//
// POST requests should be used to control a previously initialized run.
// When a new run is initialized, it will be on its first segment and its
// timer will be stopped.
//
// ### Commands
//
// Command should be sent after the token. For example, to start a new
// run, the URL would be `http://localhost:8080/run/<token>/start`.
//
//   * `reset: Reset a run; discarding any unsaved progress
//   * `start`: Start the timer of the run
//   * `split`: Save the duration of the current segment and advance to
//              the next one. Finishes the run if it's on the last segment
//   * `undo`: Go back to the previous segment
//   * `skip`: Advance to the next segment, without saving the current one
//   * `pause-toggle`: Pause/continue the timer
//   * `save`: Save a completed run to a file
//
// On success, these commands reply with `StatusNoContent` and an empty
// body.
//
// ### Description
//
// The resource `start` must be used to start tracking the run. Then, the
// resource `split` should be used whenever a goal/segment is completed,
// to save the time taken on the current segment and advance the run to
// the next one. If a `split` was issued too soon, `undo` may be used to
// revert to the previous segment. If, on the other hand, a segment should
// be ignored, `skip` may be used to advance to the next split ignoring
// the current one.
//
// When the service receives a `split` after reaching the last segment, it
// automatically stops the timer. However, the run is only saved if a
// `save` is issued to the service when the run is complete. At this
// point, the same token may be used for another run, but the run must
// first be reset with a `reset` command.
//
// Lastly, it's possible to pause/continue the timer by issuing a
// `pause-toggle`.


package run

import (
    crand "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/common"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "github.com/SirGFM/gfm-speedrun-overlay/web/timer"
    "github.com/SirGFM/gfm-speedrun-overlay/web/splits"
    "io"
    "net/http"
    "os"
    "path"
    "strconv"
    "strings"
    "sync"
    "time"
)

const Prefix = "/run"

type DurationMs struct {
    time.Duration
}

// MarshalJSON implements json.Marshaler,
// so this type may automatically convert itself to milliseconds in the JSON.
func (dur DurationMs) MarshalJSON() ([]byte, error) {
    ms := dur.Milliseconds()
    str := strconv.FormatInt(ms, 10)
    return []byte(str), nil
}

// UnmarshalJSON implements json.Unmarshaler,
// so a duration in milliseconds may be automatically decoded from the JSON.
func (dur *DurationMs) UnmarshalJSON(data []byte) error {
    var ms int64

    err := json.Unmarshal(data, &ms)
    if err != nil {
        return err
    }

    dur.Duration = time.Duration(ms) * time.Millisecond
    return nil
}

// An entry in the currently tracked run, for tracking how long that
// specific segment took and comparing it to a personal best.
type split struct {
    // The split's name.
    Name string
    // The fastest completion time for this split.
    BestTime DurationMs
    // The split's starting time, from the start of the run.
    StartTime DurationMs
    // The split's ending time, from the start of the run.
    EndTime DurationMs
    // Whether the split was skipped.
    Skipped bool
}

// A list of split entries, mainly used to encode/decode the splits for
// a run as a JSON object.
type splitList struct {
    // The list of splits
    Splits []split
}

// Local resources for a given game/category being tracked.
//
// Each game/category has its own directory in the service, within the
// run's `baseDir`, which is tracked by this indexer's `categoryDir`.
// To account for changes in the splits for a game/category, runs are
// recorded in a sub-directory named as the SHA 256 of the game/category
// splits, which is stored in `runsDir`. If needed, `runsDir` is
// automatically created whenever a new run is requested.
//
// Each directory for a different split version of a game/category has its
// own `best.json`.
type runIndexer struct {
    // Name of the game/category.
    name string
    // List of splits (as in, segments' names) in this game/category.
    splits []string
    // Directory for the game/category, where the various versions of its
    // splits are stored (in sub-directories).
    categoryDir string
    // Directory where run of the retrieved splits are stored; named after
    // the SHA 256 of the splits in use and within `categoryDir`.
    runsDir string
}

// Save `splits` to `filename` in `idx`'s `runsDir`, encoding it as a JSON.
func (idx runIndexer) _saveRun(splits []split, filename string) error {
    filePath := path.Join(idx.runsDir, filename)

    writefn := func(w io.Writer) error {
        tmp := splitList {
            Splits: splits,
        }

        enc := json.NewEncoder(w)
        err := enc.Encode(&tmp)
        if err != nil {
            return newError(err, "Couldn't save the splits", http.StatusInternalServerError)
        }
        return nil
    }
    err := common.AtomicSaveFile(idx.runsDir, filePath, writefn)
    if err != nil {
        return newError(err, "Couldn't create the splits", http.StatusInternalServerError)
    }

    return nil
}

// Save `splits` to `idx`'s `runsDir`, encoding it as a JSON, in a file
// named after the current date.
func (idx runIndexer) saveRun(splits []split) error {
    now := time.Now().UTC()
    filename := now.Format("2006-01-02T15:04:05Z.json")
    return idx._saveRun(splits, filename)
}

// Save `splits` to `idx`'s `runsDir`, encoding it as a JSON, in a file
// named `best.json`.
func (idx runIndexer) saveBestRun(splits []split) error {
    return idx._saveRun(splits, "best.json")
}

// A game/category being tracked. Note that a run doesn't synchronize
// itself, so its caller must synchronize accesses to a run structure.
type run struct {
    // Name of the game/category.
    Name string
    // Splits that make up the game.
    Splits []split
    // Splits for the best time for the game/category.
    Best []split
    // Current split.
    Current int
    // Whether the timer was started.
    Started bool
    // The token used to access the run.
    token string `json:"-"`
    // Information
    idx runIndexer `json:"-"`
    // The run's timer (ignored in the JSON).
    timer timer.LocalTimer `json:"-"`
    // Last time the run was accessed (for GC purposes).
    // Must be updated atomically!
    lastUse time.Time `json:"-"`
}

// Context for the run service.
type runCtx struct {
    // Port where the server is listening to these requests.
    listeningPort int
    // Directory where records are stored.
    baseDir string
    // Currently running splits.
    tokens map[string]*run
    // Synchronize access to the context.
    rwmut sync.RWMutex
}

// Response of a GET `new`.
type getNewResponse struct {
    // The token of the generated run.
    Token string
}

// Response of a GET `timer`.
type getTimerResponse struct {
    // The currently accumulated time, in milliseconds.
    Time int64
}

// Build a new error
func newError(err error, res string, status int) error {
    return srv_iface.NewHttpError(err, "web"+Prefix, res, status)
}

// Reset the splits in a run back to the best splits.
func (t *run) resetSplits() {
    t.Splits = nil
    for i := range t.Best {
        newSplit := t.Best[i]
        t.Splits = append(t.Splits, newSplit)
    }
}

// Reset a run back to its initial state.
func (r *run) resetRun() {
    r.timer.Stop()
    r.timer.Reset()
    r.resetSplits()
    r.Current = 0
    r.Started = false
}

// Start the run.
func (r *run) start() {
    r.Started = true
    r.Splits[0].StartTime.Duration = r.timer.Get()
    r.timer.Start()
}

// Get the starting time for the current split. If the previous segment
// finished normally, then the starting time will be the previous split's
// end time. However, if the prevous split was skipped, then the starting
// time will be the previous split's starting time.
func (r run) getSplitStartingTime() time.Duration {
    if r.Current == 0 {
        return 0
    } else if last := &r.Splits[r.Current-1]; last.Skipped {
        return last.StartTime.Duration
    } else {
        return last.EndTime.Duration
    }
}

// Advance to the next split, configuring the segment about to start.
func (r *run) advanceSplits() {
    r.Current++
    if r.Current < len(r.Splits) {
        r.Splits[r.Current].StartTime.Duration = r.getSplitStartingTime()
    } else {
        r.timer.Stop()
    }
}

// Finish the current segment, updating its best time, if it were beaten.
func (r *run) finishSegment() {
    cur := &r.Splits[r.Current]
    cur.EndTime.Duration = r.timer.Get()
    cur.Skipped = false
    if dt := cur.EndTime.Duration - cur.StartTime.Duration; dt > cur.BestTime.Duration {
        cur.BestTime.Duration = dt
    }
}

// Save a run to disk, updating the best run ever if needed.
func (r *run) saveRun() error {
    if r.Current < len(r.Splits) {
        return newError(nil, "Run still hasn't finished", http.StatusBadRequest)
    }

    if last := r.Current - 1; r.Splits[last].EndTime.Duration < r.Best[last].EndTime.Duration || r.Best[last].EndTime.Duration == 0 {
        // Update the best run
        r.Best = nil
        for i := range r.Splits {
            r.Best = append(r.Best, r.Splits[i])
            r.Best[i].BestTime = r.Splits[i].BestTime
        }

        err := r.idx.saveBestRun(r.Best)
        if err != nil {
            return newError(err, "Couldn't update 'best.json'", http.StatusInternalServerError)
        }
    }

    err := r.idx.saveRun(r.Splits)
    if err != nil {
        return newError(err, "Couldn't save the new run", http.StatusInternalServerError)
    }
    return nil
}

// Create a new `run`, to be managed by the service.
func newRun(idx runIndexer, token string, best []split) *run {
    var r run

    r.token = token
    r.Name = idx.name
    r.Best = nil
    for i := range best {
        newSplit := best[i]
        r.Best = append(r.Best, newSplit)
    }
    r.resetSplits()
    r.Current = 0
    r.Started = false
    r.idx = idx
    r.timer = timer.New()
    r.lastUse = time.Now()

    return &r
}

// Retrieve the path handled by `run`.
func (*runCtx) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*runCtx) Dependencies() []string {
    return []string{splits.Prefix}
}

// Close resources associated with the `splits` (i.e, nothing)
func (ctx *runCtx) Close() {
    ctx.rwmut.Lock()
    defer ctx.rwmut.Unlock()

    for t := range ctx.tokens {
        // TODO: Clean up the run?
        run := ctx.tokens[t]
        _ = run
        delete(ctx.tokens, t)
    }
}

// Receive the server's listening port
func (ctx *runCtx) SetListeningPort(port int) {
    ctx.listeningPort = port
}

// Generate a base64-encoded, 12-bytes random token (thus, a 16-bytes
// string). Any error returned is properly wrapped into a `HttpError`.
// The token map is accessed to ensure the generated token isn't repeated
// (as unlikely as that may be).
func (ctx *runCtx) unsafeGenerateToken() (string, error) {
    var token string

    ok := true
    for ok {
        // Generate a random token, encoded in base64. Encoding from
        // binary to base64 expands the string by 4/3. So, having a
        // 12-bytes longs (~7.9*10^-28 chances of randomly guessing)
        // results in a 16 bytes long text representation.
        var tokenBytes [12]byte
        n, err := crand.Read(tokenBytes[:])
        if err != nil {
            return "", newError(err, "Failed to generate a new token", http.StatusInternalServerError)
        } else if n != len(tokenBytes) {
            return "", newError(err, "Failed to generate enough bytes for the token", http.StatusInternalServerError)
        }

        token = base64.URLEncoding.EncodeToString(tokenBytes[:])
        // Retry on the unlikely case that it repeated?
        _, ok = ctx.tokens[token]
    }

    return token, nil
}

// Retrieve a run index from a game/category, querying the associated
// `splits` service for its split names, creating the local directory
// as needed.
func (ctx *runCtx) getRunIndex(name string) (runIndexer, error) {
    var idx runIndexer

    idx.name = name

    // Remove slashs from the name
    dirName := strings.Replace(name, "/", "%2f", -1)
    dirName = strings.Replace(dirName, "\\", "%5c", -1)
    idx.categoryDir = path.Join(ctx.baseDir, dirName)

    // Retrieve the splits for the game/category (in the local server)
    entries, err := splits.GetSplits(name, "localhost", ctx.listeningPort)
    if err != nil {
        return idx, err
    }
    idx.splits = entries

    // Get the local directory for the records
    hasher := sha256.New()
    for i := range entries {
        hasher.Write([]byte(entries[i]))
    }
    dir := hasher.Sum(nil)
    idx.runsDir = path.Join(idx.categoryDir, hex.EncodeToString(dir))

    // Create a new directory if it doesn't exist yet
    ctx.rwmut.Lock()
    err = os.MkdirAll(idx.runsDir, 0750)
    ctx.rwmut.Unlock()
    if err != nil {
        reason := "Failed to create runs directory"
        return idx, newError(err, reason, http.StatusInternalServerError)
    }

    return idx, nil
}

// Retrieve the best splits for a given indexer. If needed, it creates
// a new file in the tracked directory.
// Since this function interacts with (either reading or creating) a file
// in the run's directory, it must be synchronized by the caller!
func (ctx *runCtx) unsafeGetBestRun(idx runIndexer) ([]split, error) {
    var splits splitList

    filePath := path.Join(idx.runsDir, "best.json")
    f, err := os.Open(filePath)
    if os.IsNotExist(err) {
        // Create the first `best.json`, with empty times
        for i := range idx.splits {
            newSplit := split {
                Name: idx.splits[i],
                Skipped: true,
            }
            splits.Splits = append(splits.Splits, newSplit)
        }

        // Save it to a file
        err = idx.saveBestRun(splits.Splits)
        if err != nil {
            return nil, newError(err, "Couldn't create first 'best.json'", http.StatusInternalServerError)
        }
    } else if err == nil {
        // 'best.json' exists, so read the file
        dec := json.NewDecoder(f)
        err = dec.Decode(&splits)
        f.Close()
        if err != nil {
            return nil, newError(err, "Couldn't decode best.json", http.StatusInternalServerError)
        }
    } else if err != nil {
        // Failed to open 'best.json' at all
        return nil, newError(err, "Couldn't open best.json", http.StatusInternalServerError)
    }

    return splits.Splits, nil
}

// Handle a GET `new/<splits-path>` request, replying with a JSON-encoded
// `getNewResponse` on success.
func (ctx *runCtx) getNewRun(w http.ResponseWriter, req *http.Request, name string) error {
    idx, err := ctx.getRunIndex(name)
    if err != nil {
        return err
    }

    // Retrieve the best run. A write lock is used because we may have to
    // create a new file.
    ctx.rwmut.Lock()
    defer ctx.rwmut.Unlock()
    best, err := ctx.unsafeGetBestRun(idx)
    if err != nil {
        err = newError(err, "Failed to retrieve the best run", http.StatusInternalServerError)
        return err
    }

    token, err := ctx.unsafeGenerateToken()
    if err != nil {
        err = newError(err, "Failed to generate a token", http.StatusInternalServerError)
        return err
    }

    // Configure and save the run
    t := newRun(idx, token, best)
    ctx.tokens[token] = t

    // Reply with the token
    resp := getNewResponse {
        Token: token,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    enc := json.NewEncoder(w)
    err = enc.Encode(&resp)
    if err != nil {
        // Welp, nothing else to do... D:
        err = newError(err, "Failed to encode the responde", http.StatusInternalServerError)
        logger.Errorf("%+v", err)
    }

    return nil
}

// Retrieve a generic structure from a given token, encoding it as a JSON
// on the response. The token mutex is properly locked for reading, and the
// field is retrieved by the `getResponse` function.
func (ctx *runCtx) getGeneric(getResponse func(*run)(interface{}, error), w http.ResponseWriter, req *http.Request, token string) error {
    // Try to get the run referenced by the token
    ctx.rwmut.RLock()
    defer ctx.rwmut.RUnlock()

    r, ok := ctx.tokens[token]
    if !ok {
        return newError(nil, "Failed to find the token", http.StatusNotFound)
    }

    resp, err := getResponse(r)
    if err != nil {
        return err
    }
    enc := json.NewEncoder(w)
    err = enc.Encode(resp)
    if err != nil {
        // Welp, nothing else to do... D:
        err = newError(err, "Failed to encode the responde", http.StatusInternalServerError)
        logger.Errorf("%+v", err)
    }

    return nil
}

// Handle a GET `timer/<token>` request, replying with a JSON-encoded
// `getTimerResponse` on success.
func (ctx *runCtx) getTimer(w http.ResponseWriter, req *http.Request, token string) error {
    getResponse := func(r *run)(interface{}, error) {
        cur := r.timer.Get()
        resp := getTimerResponse {
            Time: cur.Milliseconds(),
        }
        return &resp, nil
    }

    return ctx.getGeneric(getResponse, w, req, token)
}

// Handle a GET `splits/<token>` request, replying with a JSON-encoded
// `getTimerResponse` on success.
func (ctx *runCtx) getSplits(w http.ResponseWriter, req *http.Request, token string) error {
    getResponse := func(r *run)(interface{}, error) {
        return r, nil
    }

    return ctx.getGeneric(getResponse, w, req, token)
}

// Handle GET requests.
func (ctx *runCtx) get(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) != 2 {
        return newError(nil, "Missing command/parameter", http.StatusBadRequest)
    }

    switch urlPath[0] {
    case "new":
        return ctx.getNewRun(w, req, urlPath[1])
    case "splits":
        return ctx.getSplits(w, req, urlPath[1])
    case "timer":
        return ctx.getTimer(w, req, urlPath[1])
    default:
        return newError(nil, "Invalid operation", http.StatusBadRequest)
    }

    // Shouldn't ever reach here
}

// Handle POST request.
func (ctx *runCtx) post(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) < 2 {
        return newError(nil, "Missing command (expected \"<url>/<token>/<command>\"", http.StatusBadRequest)
    }

    // Try to get the run referenced by the token
    token := urlPath[0]
    ctx.rwmut.Lock()
    defer ctx.rwmut.Unlock()
    r, ok := ctx.tokens[token]
    if !ok {
        return newError(nil, "Failed to find the token", http.StatusNotFound)
    }

    // Ensure the operation would be valid
    switch urlPath[1] {
    case "start":
        if r.Started {
            return newError(nil, "Run was already started", http.StatusBadRequest)
        }
    case "split",
        "undo",
        "skip",
        "pause-toggle":
        if !r.Started {
            return newError(nil, "Run was already started", http.StatusBadRequest)
        }
    case "reset":
        // Works in both states
    case "save":
        if !r.Started || r.Current != len(r.Splits) {
            return newError(nil, "Run must have finished before it may be saved", http.StatusBadRequest)
        }
    default:
        return newError(nil, "Invalid operation", http.StatusBadRequest)
    }

    switch urlPath[1] {
    case "reset":
        r.resetRun()
    case "start":
        r.start()
    case "split":
        r.finishSegment()
        r.advanceSplits()
    case "undo":
        if r.Current > 0 {
            r.Current--
            // Make sure the timer is restarted if it had stopped
            r.timer.Start()
        }
        if r.Current < len(r.Splits) {
            // Recover the previous end/best time.
            r.Splits[r.Current].EndTime = r.Best[r.Current].EndTime
            r.Splits[r.Current].BestTime = r.Best[r.Current].BestTime
        }
    case "skip":
        if r.Current < len(r.Splits) {
            r.Splits[r.Current].Skipped = true
        }
        r.advanceSplits()
    case "pause-toggle":
        r.timer.Toggle()
    case "save":
        err := r.saveRun()
        if err != nil {
            return err
        }
    default:
        // Shouldn't happen
        return newError(nil, "Invalid operation", http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
    return nil
}

// Handle requests to the `run` service, filtering and
// redirecting as necessary.
func (ctx *runCtx) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if req.Header.Get("Content-Type") != "application/json" {
        reason := "Content-Type must be \"application/json\""
        return newError(nil, reason, http.StatusUnsupportedMediaType)
    }
    urlPath = urlPath[1:]

    switch req.Method {
    case "GET":
        return ctx.get(w, req, urlPath)
    case "POST":
        return ctx.post(w, req, urlPath)
    default:
        return newError(nil, "Invalid method: wanted either GET or POST", http.StatusMethodNotAllowed)
    }
}

// Register a `run` handler in the `Server`.
func GetHandle(srv srv_iface.Server, baseDir string) error {
    var ctx runCtx

    ctx.baseDir = path.Clean(baseDir)
    ctx.tokens = make(map[string]*run)
    // NOTE: ctx.listeningPort is configured by the server, by calling
    // `SetListeningPort()` in the context.

    srv.AddHandler(&ctx)
    return nil
}
