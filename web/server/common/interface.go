// Interfaces used by the HTTP server. This separation was mainly to avoid
// possible circular inclusion issues.
//
// This package also has a few utility functions used by multiple servers.

package common

import (
    "net/http"
)

// A `http.Server` that is accepting requests in a separated Goroutine.
type ListeningServer interface {
    // Halts the `http.Server`.
    Close()
}

// Interface for handling HTTP request, on a given base path.
type Handler interface {
    // Base path for the handler (e.g., '/gamepad'). It must start with a
    // leading slash ('/')!
    Prefix() string
    // List every other service used by this handler.
    Dependencies() []string
    // Handle a single `http.Request`. `urlPath` is the sanitized path,
    // with the first member being the handler's prefix.
    Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error
    // Release resources associated with this object.
    Close()
}

// Interface for handling HTTP request, on a given base path, that needs
// to send requests to local services.
type LoopbackHandler interface {
    // Receive the server's listening port
    SetListeningPort(port int)
    // Also implements `Handler`
    Handler
}

// Public interface for configuring a `http.Server`.
type Server interface {
    // Add a new `Handler` to the list.
    AddHandler(handler Handler) error
    // Configure the default handler, selected in case the requested URL
    // does not match any other handler.
    SetDefault(prefix string) error
    // Start a new `ListeningServer`, on the requested "host:port", in a
    // separated Goroutine.
    //
    // `Listen()` SHALL NOT called more than once per `Server`!
    Listen(host string, port int) (ListeningServer, error)
}
