package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/SirGFM/gfm-speedrun-overlay/logger"
)

// Endpoint where the configurations are stored.
const configEndpoint = "/ram_store/config"

// Name of the token within the configuration.
const tokenName = "run-token"

// How long between the checks of the server's configuration.
const refreshRate = 5 * time.Second

type event struct {
	// The key that received the event.
	actor *Hotkey
	// Whether the key was pressed or released.
	pressed bool
}

type ctx struct {
	// The HTTP client that should be used by the entire application.
	client *http.Client
	// The service's base address.
	baseURL string
	// The run's token (if any).
	token string
	// Synchronizes access to token.
	tokenMutex sync.RWMutex
	// Channel used to receive and handle events in hotkeys.
	events chan event
}

// fetchToken fetches the configured token, if any.
func (c *ctx) fetchToken() (string, error) {
	// Fetch the config (which should be a multipart/form-data.
	resp, err := c.client.Get(c.baseURL + configEndpoint)
	if err != nil {
		logger.Errorf("hotkeys: failed to fetch the configurations: %+v", err)
		return "", err
	} else if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// Parse the data.
	reader, err := MultipartReader(resp.Header, resp.Body, false)
	if err != nil {
		logger.Errorf("hotkeys: failed to parse the configurations: %+v", err)
		return "", err
	}

	// Try to find the token within the configurations.
	for {
		next, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Errorf("hotkeys: failed to parse a configuration: %+v", err)
			return "", err
		}
		defer next.Close()

		if next.FormName() != tokenName {
			continue
		}

		var buf bytes.Buffer
		_, err = buf.ReadFrom(next)
		if err != nil {
			logger.Errorf("hotkeys: failed to parse the token: %+v", err)
			return "", err
		}

		return buf.String(), nil
	}

	// If the token isn't found (or is empty), report this as a success.
	return "", nil
}

// getToken retrieves the currently configured token, if any.
func (c *ctx) getToken() string {
	c.tokenMutex.RLock()
	token := c.token
	c.tokenMutex.RUnlock()
	return token
}

// swapToken swaps the current token by the supplied one.
func (c *ctx) swapToken(newToken string) {
	// Since this shouldn't ever get called from more than one goroutine,
	// obtain the write lock only if necessary,
	// ensure that the other goroutines won't get stalled unless the value is being updated.
	c.tokenMutex.RLock()
	doSwap := (newToken != c.token)
	c.tokenMutex.RUnlock()

	if doSwap {
		c.tokenMutex.Lock()
		c.token = newToken
		c.tokenMutex.Unlock()
	}
}

// run updates the ctx with the current run/mode.
func (c *ctx) run() {
	tick := time.NewTicker(refreshRate)
	defer tick.Stop()

	for _ = range tick.C {
		token, err := c.fetchToken()
		if err != nil {
			// The error has already been logged, simply continue.
			continue
		}

		c.swapToken(token)
	}
}

// eventHandler handle hotkey events received from the key logger.
func (c *ctx) eventHandler() {
	for event := range c.events {
		event.actor.handle(event.pressed)
	}
}
