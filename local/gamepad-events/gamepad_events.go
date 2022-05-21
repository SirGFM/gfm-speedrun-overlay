package gamepad_events

import (
	"encoding/json"
	"github.com/SirGFM/gfm-speedrun-overlay/logger"
)

const pkg = "local/gamepad_events"

// Structure used to parse the state of a gamepad of a given type.
type gamepadParser struct {
	// The gamepad's name.
	Name string
	// The gamepad's GUID encoded as a hex string.
	Guid string
	// List of inputs in the gamepad.
	Inputs map[InputName]Input
}

// A gamepad monitor.
type context struct {
	// Map a hex-encoded GUID to the gamepad's gamepadParser.
	gamepadConfig map[string]gamepadParser
}

// Load a JSON gamepad into the context, so it may be used to parse a gamepad.
func (c *context) Load(data []byte) error {
	gamepad := gamepadParser {}

	// Decode the data into a generic map
	dict := make(map[string]interface{})
	err := json.Unmarshal(data, &dict)
	if err != nil {
		logger.Errorf("%s: Failed to decode the JSON: %+v", pkg, err)
		return ErrBadJson
	}

	gamepad.Inputs = make(map[InputName]Input)
	for key, val := range dict {
		var ok bool

		switch key {
		case "name":
			gamepad.Name, ok = val.(string)
			if !ok {
				return ErrJsonInvalidName
			}
		case "guid":
			gamepad.Guid, ok = val.(string)
			if !ok {
				return ErrJsonInvalidGuid
			}
		default:
			var input genericInput

			// Try parse this as an GenericInput. To do this more easily,
			// re-encode it as a JSON and then re-decode it into the proper
			// type.
			inputData, err := json.Marshal(val)
			if err != nil {
				logger.Errorf("%s: Failed to re-encode the input: %+v", pkg, err)
				return ErrJsonInvalidInput
			}

			err = json.Unmarshal(inputData, &input)
			if err != nil {
				logger.Errorf("%s: Failed to decode the input: %+v", pkg, err)
				return ErrJsonBadInput
			}

			var inputName InputName
			inputName.UnmarshalText([]byte(key))
			gamepad.Inputs[inputName] = input.GetInput()
		}
	}

	if len(gamepad.Name) == 0 {
		return ErrJsonMissingName
	} else if len(gamepad.Guid) == 0 {
		return ErrJsonMissingGuid
	} else if len(gamepad.Inputs) == 0 {
		return ErrJsonInputEmpty
	}

	// TODO: Sync
	c.gamepadConfig[gamepad.Guid] = gamepad

	return nil
}
