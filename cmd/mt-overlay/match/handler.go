package match

import (
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
)

const Prefix = "/mt-server"

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
    if len(urlPath) == 2 && urlPath[1] == "last-update" && req.Method == http.MethodGet {
        ms := ctx.lastUpdate.Unix()
        ms *= 1000

        r := struct { Date int64 } { Date: ms }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        enc := json.NewEncoder(w)
        err := enc.Encode(&r)
        if err != nil {
            logger.Errorf("mt-server: Failed to encode the response: %+v (payload: %+v)", Prefix, err, r)
        }

        return nil
    } else if len(urlPath) == 2 && urlPath[1] == "tp" && req.Method == http.MethodGet {
        tp := ctx.getClearTangibleProgress()

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        enc := json.NewEncoder(w)
        err := enc.Encode(tp)
        if err != nil {
            logger.Errorf("mt-server: Failed to encode the response: %+v (payload: %+v)", Prefix, err, tp)
        }

        return nil
    } else {
        res := "Invalid path"
        status := http.StatusBadRequest
        return srv_iface.NewHttpError(nil, "mt-server", res, status)
    }
}
