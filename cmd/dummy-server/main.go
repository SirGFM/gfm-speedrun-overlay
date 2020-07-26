package main

import (
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/server"
    "net/http"
    "os"
    "os/signal"
)

type echo struct{}

// Retrieve the path handled by `echo`
func (*echo) Prefix() string {
    return "/echo"
}

func (*echo) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)

    // This is lazy, but should be OK for a simple test...
    data := make([]byte, req.ContentLength)
    req.Body.Read(data)
    w.Write(data)
    return nil
}

// List every other service used by this handler.
func (*echo) Dependencies() []string {
    return nil
}

// Close resources associated with the `echo` (i.e, nothing)
func (*echo) Close() {
    logger.Debugf("Closing echo...")
}

type foo struct{}

// Retrieve the path handled by `foo`
func (*foo) Prefix() string {
    return "/foo"
}

func (*foo) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)

    // This is lazy, but should be OK for a simple test...
    data := make([]byte, req.ContentLength)
    req.Body.Read(data)
    w.Write([]byte("bar: "))
    w.Write(data)
    return nil
}

// List every other service used by this handler.
func (*foo) Dependencies() []string {
    return nil
}

// Close resources associated with the `foo` (i.e, nothing)
func (*foo) Close() {
    logger.Debugf("Closing foo...")
}

func main() {
    logger.RegisterDefault(logger.LogDebug, logger.LogDebug, os.Stdout, os.Stderr)

    srv := server.New()

    err := srv.AddHandler(&echo{})
    if err != nil {
        logger.Fatalf("Failed to add 'echo' to the server: %+v", err)
    }
    err = srv.AddHandler(&foo{})
    if err != nil {
        logger.Fatalf("Failed to add 'foo' to the server: %+v", err)
    }

    lst, err := srv.Listen("", 8080)
    if err != nil {
        logger.Fatalf("Failed to start server: %+v", err)
    }
    defer lst.Close()

    intHndlr := make(chan os.Signal, 1)
    signal.Notify(intHndlr, os.Interrupt)
    <-intHndlr
    logger.Debugf("Exiting...")
}
