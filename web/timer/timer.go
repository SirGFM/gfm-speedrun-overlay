// `timer` serve a timer that may be controlled remotelly. All messages are
// JSON encoded, and requests must be sent as POST in the format:
//
//     {
//         "Action": "some-action",
//         "Value": 0,
//     }
//
// Where:
//     - Action: A string indicating the action. Must be one of:
//         - setup: Configure the timer's initial value;
//         - start: Start the timer, from its currently accumulated value;
//         - stop: Stop the timer, and keep its value unchanged;
//         - reset: Reset the timer back to its initial value;
//         - add: Increase the time by a given amount;
//         - sub: Decrease the time by a given amount;
//     - Value: The time associated with the action, in milliseconds.
//
// Among those actions, `start`, `stop` and `reset` don't need any `Value`,
// but all others (i.e., `setup`, `add` and `sub`) must be accompanied by
// a `Value`. Those actions do not generate any response on success!
//
// To retrieve the current time, issue a GET request to the service, which
// will reply with the JSON encoded response:
//
//     {
//         "Time": 0,
//     }
//
// Where `Time` is the currently accumulated time in milliseconds.

package timer

import (
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
    "sync"
    "time"
)

const Prefix = "/timer"

// Context for the timer service
type timer struct {
    // Whether the timer was started and is running
    running bool
    // When the timer was last started
    started time.Time
    // Accumulated time since the timer started running
    acc time.Duration
    // Initial time, from which the timer will count
    init time.Duration
    // Synchronize access to the context
    rwmut sync.RWMutex
}

// Interface to control a timer locally, from the server
type LocalTimer interface {
    // Start the timer, from its currently accumulated value.
    Start()
    // Stop the timer, and keep its value unchanged.
    Stop()
    // Alternate the timer between stopped and running.
    Toggle()
    // Reset the timer back to its initial value.
    Reset()
    // Configure the timer's initial value.
    Setup(time.Duration)
    // Increase the time by a given amount.
    Add(time.Duration)
    // Decrease the time by a given amount.
    Sub(time.Duration)
    // Retrieve the current time.
    Get() time.Duration
}

// Start the timer, from its currently accumulated value, without
// synchronizing the struct.
func (t *timer) unsafeStart() {
    t.started = time.Now()
    t.running = true
}

// Start the timer, from its currently accumulated value.
func (t *timer) Start() {
    t.rwmut.Lock()
    if !t.running {
        t.unsafeStart()
    }
    t.rwmut.Unlock()
}

// Stop the timer, and keep its value unchanged, without
// synchronizing the struct.
func (t *timer) unsafeStop() {
    if t.running {
        t.acc += time.Since(t.started)
    }
    t.running = false
}

// Stop the timer, and keep its value unchanged.
func (t *timer) Stop() {
    t.rwmut.Lock()
    t.unsafeStop()
    t.rwmut.Unlock()
}

// Alternate the timer between stopped and running.
func (t *timer) Toggle() {
    t.rwmut.Lock()
    if t.running {
        t.unsafeStop()
    } else {
        t.unsafeStart()
    }
    t.rwmut.Unlock()
}

// Reset the timer back to its initial value.
func (t *timer) Reset() {
    t.rwmut.Lock()
    t.acc = 0
    if t.running {
        t.started = time.Now()
    }
    t.rwmut.Unlock()
}

// Configure the timer's initial value.
func (t *timer) Setup(val time.Duration) {
    t.rwmut.Lock()
    t.init = val
    t.rwmut.Unlock()
}

// Increase the time by a given amount.
func (t *timer) Add(val time.Duration) {
    t.rwmut.Lock()
    t.acc += val
    t.rwmut.Unlock()
}

// Decrease the time by a given amount.
func (t *timer) Sub(val time.Duration) {
    t.rwmut.Lock()
    if t.acc > val {
        t.acc -= val
    } else {
        t.acc = 0
    }
    t.rwmut.Unlock()
}

// Retrieve the current time.
func (t *timer) Get() time.Duration {
    t.rwmut.RLock()
    cur := t.init + t.acc
    if t.running {
        cur += time.Since(t.started)
    }
    t.rwmut.RUnlock()

    return cur
}

// Representation of the server's response.
type response struct {
    // The currently accumulated time.
    Time int64
}

// Representation of a client's request.
type request struct {
    // The action begin requested.
    Action string
    // The action's parameter, if any.
    Value uint64 `json:",omitempty"`
}

// Retrieve a new `srv_iface.HttpError`, possibly wrapping another `error`.
func newError(inner error, reason string, status int) srv_iface.HttpError {
    return srv_iface.NewHttpError(inner, "web"+Prefix, reason, status)
}

// Retrieve the path handled by `timer`.
func (*timer) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*timer) Dependencies() []string {
    return nil
}

// Handle GET requests, returning the current time.
func (ctx *timer) get(w http.ResponseWriter, req *http.Request) error {
    cur := ctx.Get()

    r := response {
        Time: cur.Milliseconds(),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    enc := json.NewEncoder(w)
    err := enc.Encode(&r)
    if err != nil {
        logger.Errorf("web%s: Failed to encode the responde: %+v (payload: %+v)", Prefix, err, r)
    }

    return nil
}

// Handle POST requests, configuring the timer.
func (ctx *timer) post(w http.ResponseWriter, req *http.Request) error {
    var cmd request

    dec := json.NewDecoder(req.Body)
    err := dec.Decode(&cmd)
    if err != nil {
        return newError(err, "Couldn't decode the received request", http.StatusBadRequest)
    }

    t := time.Duration(cmd.Value) * time.Millisecond
    switch cmd.Action {
    case "start":
        ctx.Start()
    case "stop":
        ctx.Stop()
    case "reset":
        ctx.Reset()
    case "setup":
        ctx.Setup(t)
    case "add":
        ctx.Add(t)
    case "sub":
        ctx.Sub(t)
    default:
        return newError(nil, "Invalid operation", http.StatusBadRequest)
    }

    w.WriteHeader(http.StatusNoContent)
    return nil
}

// Handle requests to the `timer` service, filtering and redirecting as
// necessary.
func (ctx *timer) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) != 1 {
        reason := "URL must be only " + ctx.Prefix()
        return newError(nil, reason, http.StatusBadRequest)
    } else if req.Header.Get("Content-Type") != "application/json" {
        reason := "Content-Type must be \"application/json\""
        return newError(nil, reason, http.StatusUnsupportedMediaType)
    }

    switch req.Method {
    case "GET":
        return ctx.get(w, req)
    case "POST":
        return ctx.post(w, req)
    default:
        reason := "Invalid method: wanted either GET or POST"
        return newError(nil, reason, http.StatusMethodNotAllowed)
    }
}

// Close resources associated with the `timer` (i.e, nothing)
func (*timer) Close() {
}

// Register a `timer` handler in the `Server`.
func GetHandle(srv srv_iface.Server) error {
    var ctx timer

    srv.AddHandler(&ctx)
    return nil
}

// Retrieve a new LocalTimer.
func New() LocalTimer {
    return &timer{}
}
