package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// Channel used to signal that the key state changed.
	press chan bool
	// The parent context.
	c *ctx
	// Timer action executed when the context doesn't have a token.
	timerAction string
	// Run action executed when the context has a token.
	runAction string
}

// Execute implements key_events.Action,
// forwarding events to the hotkey.
func (h *Hotkey) Execute(pressed bool) {
	h.press <- pressed
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

// run watches for events,
// and executes the hotkey if the key is held down for at least threshold.
func (h *Hotkey) run() {
	lastPress := time.Now()
	wasPressed := false
	executed := false

	for pressed := range h.press {
		// Update the instant when the key was first pressed.
		if pressed && !wasPressed {
			lastPress = time.Now()
		}

		// Check if the key was pressed for long enough,
		// and fire the key's action.
		dt := time.Since(lastPress)
		if dt >= threshold && !executed {
			token := h.c.getToken()

			if token != "" {
				h.sendRunAction(token)
			} else {
				h.sendTimerAction()
			}

			// Block re-executions if the key is being held down.
			executed = pressed
		}

		wasPressed = pressed

		// After the command was executed while pressed,
		// block it until the button was released again.
		executed = executed && pressed
	}
}

// StartHotkeys starts listening for hotkeys.
func StartHotkeys(baseURL string) io.Closer {
	c := ctx{
		client:  &http.Client{},
		baseURL: baseURL,
	}

	go c.run()

	esc := Hotkey{
		press:       make(chan bool, 10),
		c:           &c,
		timerAction: "reset",
		runAction:   "reset",
	}
	go esc.run()

	space := Hotkey{
		press:       make(chan bool, 10),
		c:           &c,
		timerAction: "stop",
		runAction:   "split",
	}
	go space.run()

	backspace := Hotkey{
		press:       make(chan bool, 10),
		c:           &c,
		timerAction: "",
		runAction:   "undo",
	}
	go backspace.run()

	enter := Hotkey{
		press:       make(chan bool, 10),
		c:           &c,
		timerAction: "start",
		runAction:   "start",
	}
	go enter.run()

	sKey := Hotkey{
		press:       make(chan bool, 10),
		c:           &c,
		timerAction: "",
		runAction:   "save",
	}
	go sKey.run()

	cfg := key_events.WatcherConfig{
		PoolPerSec: 20,
		OnKeyPress: map[key_logger.Key]key_events.Action{
			key_logger.Backspace: &backspace,
			key_logger.Return:    &enter,
			key_logger.Space:     &space,
			key_logger.Esc:       &esc,
			key_logger.S:         &sKey,
		},
		OnKeyRelease: map[key_logger.Key]key_events.Action{
			key_logger.Backspace: &backspace,
			key_logger.Return:    &enter,
			key_logger.Space:     &space,
			key_logger.Esc:       &esc,
			key_logger.S:         &sKey,
		},
	}

	keyWatcher := key_events.NewEventWatcher(cfg)
	return keyWatcher
}
