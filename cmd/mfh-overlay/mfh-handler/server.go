package mfh_handler

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "time"
)

// The context that store page's data.
type serverContext struct {
    // Last time the structure was updated.
    lastUpdate time.Time
}

// Update the context, so watching clients may automatically refresh.
func (ctx *serverContext) update() {
    ctx.lastUpdate = time.Now()
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
