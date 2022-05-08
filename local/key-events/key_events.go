// key_events listen to keyboard events and may both trigger configurable
// events for specific keys or send this data to a `/ram_store`-like
// endpoint, so it may later be used by another service. This module must
// run locally, instead of being serviced as an endpoint.
//
// The module may use two resources within a configurable base endpoint
// (for example, '/ram_store/keyboard`):
//
//     * base_endpoint + '/data': A 'application/json' encoded array in
// which each index represents the state of the key (0 for released, 1 for
// pressed);
//     * base_endpoint + '/map': A 'application/json' encoded array naming
// every key in the same order as they appear in base_endpoint + `/data`.
//
// This base endpoint may be left unconfigured, in which case this module
// will only trigger the registered key events.

package key_events

import (
	"bytes"
	"encoding/json"
    key_logger "github.com/SirGFM/goLogKeys/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "io"
	"net/http"
	"time"
)

// A trigger configured to run whenever a button is either pressed or
// released.
type Action interface {
	Execute(pressed bool)
}

// Configures a key logger.
type WatcherConfig struct {
	// How often should the key logger be checked each second.
	PoolPerSec int
	// Base URL where a ram_store service is running.
	BaseStoreEndpoint string
	// URL where the current data is stored.
	dataStoreEndpoint string
	// URL where the mapping between indexes and keys is stored.
	mapStoreEndpoint string
	// Map of actions executed when a given key is pressed.
	OnKeyPress map[key_logger.Key]Action
	// Map of actions executed when a given key is released.
	OnKeyRelease map[key_logger.Key]Action
}

// The local key logger, implementing the Watcher interface.
type watcher struct {
	// They key logger.
	kl key_logger.Logger
	// The logger's configuration.
	config WatcherConfig
	// Ticker that signals whenever the key logger should be verified.
	poolInterval *time.Ticker
	// Client used to communicate with the ram_store service.
	httpClient http.Client
	// Workaround channel to immediately close the key logger.
	stop chan struct{}
}

// Key logger interface.
type Watcher interface {
	io.Closer
}

// doForEachKey execute the callback for every valid key.
func doForEachKey(callback func (key_logger.Key, string)) {
	for key := key_logger.Key(0); key < key_logger.KeyCount; key++ {
		callback(key, key.String())
	}
}

// NewEventWatcher configures and starts the key logger.
func NewEventWatcher(config WatcherConfig) Watcher {
	if min, max := 0, int(time.Second); config.PoolPerSec <= min || config.PoolPerSec >= max {
		logger.Fatalf("key_events: Invalid number of checks per second! PoolPerSec must be between %d and %d", min, max)
	}

	kl, err := key_logger.GetLogger()
	if err != nil {
		logger.Fatalf("key_events: Failed to create a new key logger: %+v", err)
	}

	// Configure a defer to automatically clean up the key logger on failure.
	success := false
	defer func () {
		if !success {
			kl.Clean()
		}
	} ()

	err = kl.Setup()
	if err != nil {
		logger.Fatalf("key_events: Failed to configure the key logger: %+v", err)
	}

	err = kl.Start()
	if err != nil {
		logger.Fatalf("key_events: Failed to start the key logger: %+v", err)
	}

	newConfig := WatcherConfig {
		OnKeyPress: make(map[key_logger.Key]Action),
		OnKeyRelease: make(map[key_logger.Key]Action),
	}
	if len(config.BaseStoreEndpoint) > 0 {
		newConfig.dataStoreEndpoint = config.BaseStoreEndpoint + "/data"
		newConfig.mapStoreEndpoint = config.BaseStoreEndpoint + "/map"
	}
	for k, v := range config.OnKeyPress {
		newConfig.OnKeyPress[k] = v
	}
	for k, v := range config.OnKeyRelease {
		newConfig.OnKeyRelease[k] = v
	}

	duration := time.Second / time.Duration(config.PoolPerSec)
	w := watcher {
		kl: kl,
		config: newConfig,
		poolInterval: time.NewTicker(duration),
		stop: make(chan struct{}, 1),
	}

	// If logging to a ram_store, create a mapping of index->key_name and save
	// it in the store.
	if len(config.BaseStoreEndpoint) > 0 {
		var keyNames []string

		doForEachKey(func (key key_logger.Key, name string) {
			keyNames = append(keyNames, name)
		})

		data, err := json.Marshal(keyNames)
		if err != nil {
			logger.Errorf("key_events: Failed to encode tha key map as JSON: %+v", err)
		} else {
			w.send(data, w.config.mapStoreEndpoint)
		}
	}

	go w.run()

	logger.Infof("key_events: Started key logger to '%s' every %0.2fs",
			newConfig.dataStoreEndpoint,
			float32(duration) / float32(time.Second))

	success = true
	return &w
}

// send data to the the remote service in url.
func (w *watcher) send(data []byte, url string) {
	buf := bytes.NewBuffer(data)

	resp, err := w.httpClient.Post(url, "application/json", buf)
	if err != nil {
		logger.Errorf("key_events: Failed to send the keys to %s: %+v", err, url)
	} else if code := resp.StatusCode; code != http.StatusOK && code != http.StatusNoContent {
		logger.Errorf("key_events: Failed to send the keys to %s: %+s", resp.Status, url)
	}
}

// run the key logger, executing triggers and reporting to the configure
// ram_store.
func (w *watcher) run() {
	var keyBuf [100]key_logger.Key
	var stateBuf [100]key_logger.KeyState

	// Generate a list with states for each tracked key.
	var keyBool []int
	doForEachKey(func (key key_logger.Key, name string) {
		keyBool = append(keyBool, 0)
	})

from_outer:
	for {
		// Ticker's channel is read-only and therefore cannot be closed... Thus,
		// the logic to stop the goroutine must be implemented independently
		// from that channel.
		select {
		case <-w.stop:
			logger.Debugf("key_events: Exiting...")
			break from_outer
		case <-w.poolInterval.C:
		}

		changed := false

		// Pop keys from the key logger until it's empty.
from_inner:
		for {
			key, state, err := w.kl.PopMulti(keyBuf[:], stateBuf[:])
			if err != nil {
				logger.Errorf("key_events: Failed to read more keys: %+v", err)
			} else if len(key) == 0 {
				break from_inner
			}
			changed = true

			// Execute actions and update each logged key.
			for i := range key {
				var act Action
				var ok, pressed bool

				if state[i] == key_logger.Released {
					keyBool[key[i]] = 0

					act, ok = w.config.OnKeyRelease[key[i]]
					pressed = false
				} else {
					keyBool[key[i]] = 1

					act, ok = w.config.OnKeyPress[key[i]]
					pressed = true
				}

				if ok {
					act.Execute(pressed)
				}
			}
		}

		// If any key changed, report to the ram_store.
		if changed && len(w.config.dataStoreEndpoint) > 0 {
			data, err := json.Marshal(keyBool)
			if err != nil {
				logger.Errorf("key_events: Failed to encode tha changed keys as JSON: %+v", err)
			} else {
				w.send(data, w.config.dataStoreEndpoint)
			}
		}
	}
}

// Close stops the key logger.
func (w *watcher) Close() error {
	w.kl.Clean()
	if w.poolInterval != nil {
		w.poolInterval.Stop()
		close(w.stop)
		w.poolInterval = nil
	}

	return nil
}
