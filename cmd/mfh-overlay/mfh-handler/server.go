package mfh_handler

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "sync"
    "sync/atomic"
    "time"
)

// Simple timer for tracking when the service was updated.
type timer struct {
    // The stored time.
    t time.Time
    // Synchronize access to the time.Time.
    m sync.RWMutex
    // (hack to) try locking for writing without blocking.
    tryWrite uint32
}

// The context that store page's data.
type serverContext struct {
    // Last time the structure was updated.
    lastUpdate timer
}

// Retrieve the last updated time
func (ctx *serverContext) getLastUpdate() int64 {
    ctx.lastUpdate.m.RLock()
    ms := ctx.lastUpdate.t.Unix()
    ctx.lastUpdate.m.RUnlock()
    return ms * 1000
}

// Update the context, so watching clients may automatically refresh.
func (ctx *serverContext) update() {
    // (tl;dr: This is overkill)
    //
    // Ensure that only a single writer may ever acquire the write lock
    // (using atomic.CompareAndSwapUint32 to only allow a single goroutine
    // to enter the conditional at a time) and then lock it for writing.
    if atomic.CompareAndSwapUint32(&ctx.lastUpdate.tryWrite, 0, 1) {
        ctx.lastUpdate.m.Lock()
        ctx.lastUpdate.t = time.Now()
        ctx.lastUpdate.m.Unlock()
        atomic.StoreUint32(&ctx.lastUpdate.tryWrite, 0)
    }
}

// Clean up the container, removing all associated resources.
func (ctx *serverContext) Close() {
}

// `srv_iface.Server`, so it may be used as a server and for templating.
type Context interface {
    srv_iface.Handler
}

// Retrieve a new data server
func New() Context {
    ctx := serverContext {}
    return &ctx
}
