package mfh_handler

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
)

const Prefix = "/mfh-handler"

// Retrieve the path handled by `serverContext`.
func (*serverContext) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*serverContext) Dependencies() []string {
    return nil
}

// Implement the server interface, so serverContext may report when it changes
func (ctx *serverContext) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    res := "Invalid path"
    status := http.StatusBadRequest
    return srv_iface.NewHttpError(nil, "mfh-handler", res, status)
}
