package mfh_handler

import (
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
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

// Handle popup requests in this server
func (ctx *serverContext) handlePopup(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    switch req.Method {
    case http.MethodGet:
        popups := ctx.getPopupList()
        r := struct { Elements []popup } { Elements: popups }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        enc := json.NewEncoder(w)
        err := enc.Encode(&r)
        if err != nil {
            logger.Errorf("mfh-handler: Failed to encode the response: %+v (payload: %+v)", Prefix, err, r)
        }
        break;
    case http.MethodPost:
        var got popup

        dec := json.NewDecoder(req.Body)
        err := dec.Decode(&got)
        if err != nil {
            logger.Errorf("mfh-handler: Failed to decode popup data: %+v", err)
            return BadJSONInput
        }

        ctx.pushPopupElement(got.Id, got.Timeout)
        w.WriteHeader(http.StatusNoContent)
        break;
    default:
        res := "Invalid method for /mfh-handler/popup"
        status := http.StatusMethodNotAllowed
        return srv_iface.NewHttpError(nil, "mfh-handler", res, status)
    }

    return nil
}

// Handle overlay-extras requests in this server
func (ctx *serverContext) handleOverlayExtras(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    switch req.Method {
    case http.MethodGet:
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        ctx.getExtraJSON(w)
        break;
    case http.MethodPost:
        err := ctx.putExtraJson(req.Body)
        if err != nil {
            return err
        }
        w.WriteHeader(http.StatusNoContent)
        break;
    default:
        res := "Invalid method for /mfh-handler/overlay-extra"
        status := http.StatusMethodNotAllowed
        return srv_iface.NewHttpError(nil, "mfh-handler", res, status)
    }

    return nil
}

// Implement the server interface, so serverContext may report when it changes
func (ctx *serverContext) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if len(urlPath) == 2 && urlPath[1] == "last-update" && req.Method == http.MethodGet {
        ms := ctx.getLastUpdate()

        r := struct { Date int64 } { Date: ms }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        enc := json.NewEncoder(w)
        err := enc.Encode(&r)
        if err != nil {
            logger.Errorf("mfh-handler: Failed to encode the response: %+v (payload: %+v)", Prefix, err, r)
        }

        return nil
    } else if len(urlPath) == 2 && urlPath[1] == "popup" {
        return ctx.handlePopup(w, req, urlPath)
    } else if len(urlPath) == 2 && urlPath[1] == "overlay-extras" {
        return ctx.handleOverlayExtras(w, req, urlPath)
    } else {
        res := "Invalid path"
        status := http.StatusBadRequest
        return srv_iface.NewHttpError(nil, "mfh-handler", res, status)
    }
}
