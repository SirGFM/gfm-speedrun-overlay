// A simple, customizable HTTP server. `New()` initializes a new server,
// which may then be configured by calls to `AddHandler()`. As soon as
// every base path is associated with the server, it should be started with
// a `Listen()` call. After this point, the server won't accept any other
// call!
//
// `Server.Listen()` returns a `ListeningServer`, which closes every
// handler alongside the HTTP server. `ListeningServer` waits until there's
// no pending request before closing its associated `Handler`s.

package server

import (
    "fmt"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
    "net/url"
    "path"
    "strconv"
    "strings"
    "sync"
)

// `error` used by this package.
type ErrorCode uint
const (
    // Prefix already registered in the server
    RepeatedPrefix ErrorCode = iota
    // Invalid prefix
    BadPrefix
    // Handlers must start with a leading '/'
    NonRootPrefix
    // Invalid port
    BadPort
    // Handler not registered (yet)
    InvalidHandler
)

// `Error()` implements the `error` interface for `ErrorCode`.
func (e ErrorCode) Error() string {
    switch e {
    case RepeatedPrefix:
        return "Prefix already registered in the server"
    case BadPrefix:
        return "Invalid prefix"
    case NonRootPrefix:
        return "Handlers must start with a leading '/'"
    case BadPort:
        return "Invalid port"
    case InvalidHandler:
        return "Handler not registered (yet)"
    default:
        return "Unknown error"
    }
}

type runningServer struct {
    // Go's default server that handles requests from clients.
    httpServer *http.Server
    // Default handler, used in case the URL doesn't match anything.
    defaultHandler srv_iface.Handler
    // List of accepted `Handler`s.
    handlers []srv_iface.Handler
    // Synchronize access to handlers while closing
    closing sync.RWMutex
}

// Check if a given `url` has `prefix` as its base path.
func comparePrefix(prefix, url string) bool {
    return strings.HasPrefix(url, prefix) &&
           (len(url) == len(prefix) || url[len(prefix)] == '/')
}

// ServeHTTP is called by Go's http package whenever a new HTTP request arrives
func (s *runningServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    // Normalize and strip the URL from its leading prefix (and slash)
    resUrl := path.Clean(req.URL.EscapedPath())
    if len(resUrl) > 0 && resUrl[0] == '/' {
        // NOTE: The first character must no be a '/' because of the split
        resUrl = resUrl[1:]
    } else if len(resUrl) == 1 && resUrl[0] == '.' {
        // Clean converts an empty path into a single "."
        resUrl = ""
    }

    // As part of the normalization, unescape each component individually
    var urlPath []string
    for _, p := range strings.Split(resUrl, "/") {
        cleanPath, err := url.PathUnescape(p)
        if err != nil {
            reason := fmt.Sprintf("Couldn't normalize the requested resource (%s) - %+v", resUrl, err)
            status := http.StatusInternalServerError
            srv_iface.ReplyStatus("web/server", status, reason, w)
            return
        }
        urlPath = append(urlPath, cleanPath)
    }

    logger.Debugf("New request from %+v: %s /%s", req.RemoteAddr, req.Method, resUrl)

    var err404, err error
    reason := "Couldn't find the requested resource (" + resUrl + ")"
    status := http.StatusNotFound
    err404 = srv_iface.NewHttpError(nil, "web/server", reason, status)

    // Look for, and execute, the associated prefix
    s.closing.RLock()
    defer s.closing.RUnlock()
    err = err404
    for i := range s.handlers {
        // Compare to the prefix skipping the leading '/'
        if s.handlers[i].Prefix()[1:] == urlPath[0] {
            err = s.handlers[i].Handle(w, req, urlPath)
            break
        }
    }

    // If the URL didn't match anything, try the default handler
    if err == err404 && s.defaultHandler != nil {
        // Remove the leading '/' from the prefix
        prefix := s.defaultHandler.Prefix()[1:]

        var defPath []string
        defPath = append(defPath, prefix)
        defPath = append(defPath, urlPath...)
        err = s.defaultHandler.Handle(w, req, defPath)
    }

    if err == nil {
        logger.Infof("%s: %s - OK", req.Method, resUrl)
    } else if herr, ok := err.(srv_iface.HttpError); ok {
        logger.Errorf("%s: %s - %s", req.Method, resUrl, herr.GetHttpStatus())
        logger.Debugf("%+v", herr)
        srv_iface.ReplyHttpError(herr, w)
    } else {
        // Shouldn't happend
        logger.Errorf("%s: %s - ERROR", req.Method, resUrl)
    }
}

// Halts the `http.Server`, if still running
func (s *runningServer) Close() {
    if s.httpServer != nil {
        s.httpServer.Close()
        s.httpServer = nil
    }
    // Ensure no request is being handled before closing everything
    s.closing.Lock()
    defer s.closing.Unlock()
    for len(s.handlers) > 0 {
        s.handlers[0].Close()
        s.handlers[0] = nil
        s.handlers = s.handlers[1:]
    }
}

// Tracks the handlers to be used when starting a new `ListeningServer`.
type setupServer struct {
    // Default handler, used in case the URL doesn't match anything.
    defaultHandler srv_iface.Handler
    // List of `Handler`s, with their associated prefix for an easy and
    // fast lookup.
    handlers map[string]srv_iface.Handler
    // Whether this `setupServer` is still valid, or it cannot be used
    // anymore.
    valid bool
}

// Add a new `Handler` to the list.
func (s *setupServer) AddHandler(handler srv_iface.Handler) error {
    prefix := handler.Prefix()
    if _, ok := s.handlers[prefix]; ok {
        return RepeatedPrefix
    } else if len(prefix) == 0 {
        return BadPrefix
    } else if prefix[0] != '/' {
        return NonRootPrefix
    }

    s.handlers[prefix] = handler
    return nil
}

// Configure the default handler, selected in case the requested URL
// does not match any other handler.
func (s *setupServer) SetDefault(prefix string) error {
    handler, ok := s.handlers[prefix]
    if !ok {
        return InvalidHandler
    }

    s.defaultHandler = handler
    return nil
}

// Start a new `ListeningServer`, on the requested "host:port", in a
// separated Goroutine.
func (s *setupServer) Listen(host string, port int) (srv_iface.ListeningServer, error) {
    if !s.valid {
        logger.Fatalf("Trying to listen on an invalid `Server`!")
    }

    if port <= 0 || port >= 0x10000 {
        return nil, BadPort
    }

    // Check that every dependency is met in the handlers.
    for _, h := range s.handlers {
        for _, dep := range h.Dependencies() {
            if _, ok := s.handlers[dep]; !ok {
                logger.Fatalf("Missing dependency \"%s\" for service \"%s\"!", dep, h.Prefix())
            }
        }
    }

    var srv runningServer
    srv.defaultHandler = s.defaultHandler

    // Convert `setupServer`'s maps of `Handlers` in a list, for
    // `runningServer`. Also assign the listening port, if needed.
    for p, h := range s.handlers {
        srv.handlers = append(srv.handlers, h)
        delete(s.handlers, p)

        if lh, ok := h.(srv_iface.LoopbackHandler); ok && lh != nil {
            lh.SetListeningPort(port)
        }
    }
    s.handlers = nil

    // Configure and start the `http.Server`
    sport := strconv.Itoa(port)
    srv.httpServer = &http.Server {
        Addr: host + ":" + sport,
        Handler: &srv,
    }

    go func() {
        logger.Debugf("Waiting...")
        srv.httpServer.ListenAndServe()
    } ()

    // Invalidate the `Server`, so it may not be used anymore.
    s.valid = false
    return &srv, nil
}

// Retrieve a new, empty `srv_iface.Server`.
func New() srv_iface.Server {
    return &setupServer{
        handlers: make(map[string]srv_iface.Handler),
        valid: true,
    }
}
