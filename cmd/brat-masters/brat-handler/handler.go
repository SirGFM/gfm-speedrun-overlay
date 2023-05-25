package brat_handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SirGFM/gfm-speedrun-overlay/logger"
	srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
)

const Prefix = "/brat-handler"

// Retrieve the path handled by `serverContext`.
func (*serverContext) Prefix() string {
	return Prefix
}

// List every other service used by this handler.
func (*serverContext) Dependencies() []string {
	return nil
}

// Map resources into themselves (as this doesn't need any fancy mapping).
func (ctx *serverContext) Map(resource []string) ([]string, error) {
	if len(resource) == 2 && resource[0] == "tmpl" {
		endpoint := resource[1]

		for _, suffix := range []string{
			"setup.go.html",
			"game.go.html",
			"twitch-iframe.go.html",
		} {
			if strings.HasSuffix(endpoint, suffix) {
				return []string{"tmpl", suffix}, nil
			}
		}
	}
	return resource, nil
}

// Implement the server interface, so serverContext may report when it changes
func (ctx *serverContext) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
	if len(urlPath) == 2 && urlPath[1] == "last-update" && req.Method == http.MethodGet {
		ms := ctx.getLastUpdate()

		r := struct{ Date int64 }{Date: ms}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		enc := json.NewEncoder(w)
		err := enc.Encode(&r)
		if err != nil {
			logger.Errorf("brat-handler: Failed to encode the response: %+v (payload: %+v)", err, r)
		}
	} else if len(urlPath) == 2 && urlPath[1] == "countdown" {
		switch req.Method {
		case http.MethodGet:
			var curTime time.Duration

			acc := ctx.countdown.Get()
			ctx.m.RLock()
			from := ctx.countdownFrom
			ctx.m.RUnlock()
			if acc < from {
				curTime = (from - acc) / time.Millisecond
			} else {
				ctx.countdown.Stop()
			}

			r := struct{ Time int64 }{Time: int64(curTime)}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			enc := json.NewEncoder(w)
			err := enc.Encode(&r)
			if err != nil {
				logger.Errorf("brat-handler: Failed to encode the response: %+v (payload: %+v)", err, r)
			}
		case http.MethodPost:
			var value uint64

			dec := json.NewDecoder(req.Body)
			err := dec.Decode(&value)
			if err != nil {
				return srv_iface.NewHttpError(err, "bra-handler", "Couldn't decode the received countdown", http.StatusBadRequest)
			}

			from := time.Millisecond * time.Duration(value)
			ctx.m.Lock()
			ctx.countdownFrom = from
			ctx.m.Unlock()

			ctx.countdown.Reset()
			ctx.countdown.Start()

			w.WriteHeader(http.StatusOK)
		default:
			res := "Invalid countdown method; want GET or POST"
			status := http.StatusMethodNotAllowed
			return srv_iface.NewHttpError(nil, "brat-handler", res, status)
		}
	} else {
		res := "Invalid path"
		status := http.StatusBadRequest
		return srv_iface.NewHttpError(nil, "brat-handler", res, status)
	}

	return nil
}
