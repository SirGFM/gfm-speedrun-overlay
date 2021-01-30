package mfh_handler

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
)

// The context that store page's data.
type serverContext struct {
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
