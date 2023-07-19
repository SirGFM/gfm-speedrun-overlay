package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	key_events "github.com/SirGFM/gfm-speedrun-overlay/local/key-events"
	"github.com/SirGFM/gfm-speedrun-overlay/logger"
	key_logger "github.com/SirGFM/goLogKeys/logger"
)

// How long until a key press is detected.
const threshold = 300 * time.Millisecond

// Endpoint for 'run' requests.
const runEndpoint = "/run"

// Endpoint for 'timer' requests.
const timerEndpoint = "/timer"

type Hotkey struct {
	// The parent context.
	c *ctx
	// Timer action executed when the context doesn't have a token.
	timerAction string
	// Run action executed when the context has a token.
	runAction string
	// When the last event was received.
	lastPress time.Time
	// Whether the key was pressed in a previous event.
	isPressed bool
	// Whether the event has already been dispatched for the current press.
	dispatched bool
}

// Execute implements key_events.Action,
// forwarding events to the hotkey.
func (h *Hotkey) Execute(pressed bool) {
	h.c.events <- event{
		actor:   h,
		pressed: pressed,
	}
}

// sendRunAction sends an action to the /run endpoint.
func (h *Hotkey) sendRunAction(token string) {
	if h.runAction == "" {
		return
	}

	endpoint := fmt.Sprintf("%s%s/%s/%s", h.c.baseURL, runEndpoint, token, h.runAction)

	resp, err := h.c.client.Post(endpoint, "application/json", nil)
	if resp != nil {
		resp.Body.Close()
	}
	if err != nil {
		logger.Errorf("failed to send the run action: %+v", err)
	}
}

// sendTimerAction sends an action to the /timer endpoint.
func (h *Hotkey) sendTimerAction() {
	if h.timerAction == "" {
		return
	}

	var req struct {
		Action string
	}
	req.Action = h.timerAction

	data, err := json.Marshal(&req)
	if err != nil {
		logger.Errorf("failed to encode the timer action: %+v", err)
		return
	}

	body := bytes.NewBuffer(data)
	resp, err := h.c.client.Post(h.c.baseURL+timerEndpoint, "application/json", body)
	if resp != nil {
		resp.Body.Close()
	}
	if err != nil {
		logger.Errorf("failed to send the timer action: %+v", err)
	}
}

// handle manages the key press,
// dispatching the command if the key is held down for at least threshold.
func (h *Hotkey) handle(pressed bool) {
	if h.lastPress.IsZero() {
		h.lastPress = time.Now()
	}

	// Update the instant when the key was first pressed.
	if pressed && !h.isPressed {
		h.lastPress = time.Now()
	}

	// Check if the key was pressed for long enough,
	// and fire the key's action.
	dt := time.Since(h.lastPress)
	if dt >= threshold && !h.dispatched {
		token := h.c.getToken()

		if token != "" {
			h.sendRunAction(token)
		} else {
			h.sendTimerAction()
		}

		// Block re-executions if the key is being held down.
		h.dispatched = pressed
	}

	h.isPressed = pressed

	// After the command was executed while pressed,
	// block it until the button was released again.
	h.dispatched = h.dispatched && pressed
}

// StartHotkeys starts listening for hotkeys.
func StartHotkeys(baseURL, configFilename string) io.Closer {
	c := ctx{
		client:  &http.Client{},
		baseURL: baseURL,
		events:  make(chan event, 100),
	}

	go c.run()

	keyConfig := key_events.WatcherConfig{
		PoolPerSec:   10,
		OnKeyPress:   make(map[key_logger.Key]key_events.Action),
		OnKeyRelease: make(map[key_logger.Key]key_events.Action),
	}

	config := parseINI(configFilename)

	if data, ok := config["config"]; ok {
		poolRate := data["pool-rate"]
		if poolRate != "" {
			num, err := strconv.ParseInt(poolRate, 0, 32)
			if err != nil {
				logger.Fatalf("hotkeys: Failed to parse [config].pool-rate: %+v", err)
			}

			keyConfig.PoolPerSec = int(num)
			logger.Infof("hotkeys: pool rate - %d per second", num)
		}

		delete(config, "config")
	}

	for keyName, values := range config {
		key := key_logger.GetKey(keyName)
		if key == 0 {
			logger.Warnf("hotkeys: ignoring invalid key '%s'", key)
			continue
		}

		hotkey := Hotkey{
			c:           &c,
			timerAction: values["timer"],
			runAction:   values["run"],
		}

		keyConfig.OnKeyPress[key] = &hotkey
		keyConfig.OnKeyRelease[key] = &hotkey

		logger.Infof(
			"hotkeys: registered key '%s' with actions: (run: '%s') (timer: '%s')",
			keyName,
			hotkey.timerAction,
			hotkey.runAction,
		)
	}

	keyWatcher := key_events.NewEventWatcher(keyConfig)
	go c.eventHandler()

	return keyWatcher
}

func PrintKeys() {
	var key key_logger.Key

	fmt.Println("Valid keys:")
	for ; key < key_logger.KeyCount; key++ {
		keyName := strings.ToLower(key.String())
		if keyName != "" {
			fmt.Printf("\t - %s\n", keyName)
		}
	}
}
