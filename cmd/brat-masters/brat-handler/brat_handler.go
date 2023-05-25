package brat_handler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/SirGFM/gfm-speedrun-overlay/logger"
	srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
	web_timer "github.com/SirGFM/gfm-speedrun-overlay/web/timer"
	"github.com/SirGFM/gfm-speedrun-overlay/web/tmpl"
)

// `srv_iface.Server`, so it may be used as a server and for templating.
type Context interface {
	tmpl.DataCRUD
	tmpl.Mapper
	srv_iface.Handler
}

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
	// Data used to customize pages, keyed by the SHA-256 of the path.
	data map[string]interface{}
	// Synchronize access to the data.
	m sync.RWMutex
	// Last time the structure was updated.
	lastUpdate timer
	// The countdown timer
	countdown web_timer.LocalTimer
	// The initial time for the countdown
	countdownFrom time.Duration
}

// Retrieve a new data server
func New() Context {
	ctx := serverContext{
		countdown: web_timer.New(),
	}
	return &ctx
}

// Clean up the container, removing all associated resources.
func (ctx *serverContext) Close() {
}

// store stores the resource into the serverContext.
func (ctx *serverContext) store(data tmpl.DataReader) error {
	var val map[string]interface{}

	dec := json.NewDecoder(data)
	err := dec.Decode(&val)
	if err != nil {
		logger.Errorf("brat-handler: Failed to decode %+v's data: %+v", data.URLPath(), err)
		return BadJSONInput
	}
	ctx.m.Lock()
	ctx.data = val
	ctx.m.Unlock()
	ctx.setLastUpdate()

	return nil
}

// Create a new resource. Identical to "Update".
func (ctx *serverContext) Create(resource []string, data tmpl.DataReader) error {
	return ctx.store(data)
}

// Update an already existing resource. Identical to "Create".
func (ctx *serverContext) Update(resource []string, data tmpl.DataReader) error {
	return ctx.store(data)
}

// Remove the resource.
func (ctx *serverContext) Delete(resource []string) error {
	ctx.m.Lock()
	ctx.data = nil
	ctx.m.Unlock()
	return nil
}

// Retrieve the data associated with a given resource.
func (ctx *serverContext) Read(resource []string) (interface{}, error) {
	tmp := make(map[string]interface{})

	ctx.m.RLock()
	for key, value := range ctx.data {
		tmp[key] = value
	}
	ctx.m.RUnlock()

	var endpoint string
	if len(resource) == 2 && resource[0] == "tmpl" {
		endpoint = resource[1]
	}

	switch endpoint {
	case "setup.go.html":
		tmp["HideTimer"] = false
	case "no-timer-setup.go.html":
		tmp["HideTimer"] = true
	case "left-twitch-iframe.go.html":
		tmp["TwitchUsername"] = tmp["LeftTwitch"]
	case "right-twitch-iframe.go.html":
		tmp["TwitchUsername"] = tmp["RightTwitch"]
	}

	return tmp, nil
}
