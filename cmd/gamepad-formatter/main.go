package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	key_logger "github.com/SirGFM/goLogKeys/logger"
	gptools "github.com/SirGFM/gfm-speedrun-overlay/gamepad-tools"
	"github.com/SirGFM/gfm-speedrun-overlay/gamepad-tools/monitor"
	"github.com/SirGFM/gfm-speedrun-overlay/local/key-events"
	"io"
	"math"
	"os"
	"os/signal"
	"sync/atomic"
	"strconv"
	"time"
)

// How often the hotkey is checked for events each second.
const default_check_per_second = 10

// Sample gamepad made with https://kenney-assets.itch.io/input-prompts-pixel-16
func fatalf(str string, args ...interface{}) {
	panic(fmt.Sprintf(str, args...))
}

/* === Button mapper ======================================================== */

// inputType describe different input types in a gamepad, mostly for the input
// request prompt.
type inputType uint32
const (
	ButtonInput inputType = iota
	AxisInput
)

// String returns the prompt to request an input of the given type.
func (it inputType) String() string {
	switch it {
	case ButtonInput: return "Press and hold"
	case AxisInput: return "Move"
	default:
		fatalf("Invalid input type '%d'", it)
		return ""
	}
}

// input describe every input type in a gamepad handled by the application.
type input uint32
const (
	AButton input = iota
	BButton
	XButton
	YButton
	RbButton
	RtButton
	LbButton
	LtButton
	StartButton
	SelectButton
	DPadLeftButton
	DPadRightButton
	DPadUpButton
	DPadDownButton
	LeftStickHorizontalLeft
	LeftStickHorizontalRight
	LeftStickVerticalUp
	LeftStickVerticalDown
	RightStickHorizontalLeft
	RightStickHorizontalRight
	RightStickVerticalUp
	RightStickVerticalDown
	NumInputs
)

// Name returns a human-readable name for the input.
func (i input) Name() string {
	switch i {
	case AButton: return "'A' button"
	case BButton: return "'B' button"
	case XButton: return "'X' button"
	case YButton: return "'Y' button"
	case RbButton: return "'RB' button"
	case RtButton: return "'RT' button"
	case LbButton: return "'LB' button"
	case LtButton: return "'LT' button"
	case StartButton: return "'Start' button"
	case SelectButton: return "'Select' button"
	case DPadLeftButton: return "D-pad Left"
	case DPadRightButton: return "D-pad Right"
	case DPadUpButton: return "D-pad Up"
	case DPadDownButton: return "D-pad Down"
	case LeftStickHorizontalLeft: return "Left Stick Horizontally to the Left"
	case LeftStickHorizontalRight: return "Left Stick Horizontally to the Right"
	case LeftStickVerticalUp: return "Left Stick Horizontally Upward"
	case LeftStickVerticalDown: return "Left Stick Horizontally Downward"
	case RightStickHorizontalLeft: return "Right Stick Horizontally to the Left"
	case RightStickHorizontalRight: return "Right Stick Horizontally to the Right"
	case RightStickVerticalUp: return "Right Stick Horizontally Upward"
	case RightStickVerticalDown: return "Right Stick Horizontally Downward"
	default:
		fatalf("Invalid input '%d'", i)
		// Note: This return isn't reached!
		return ""
	}
}

// Name returns the name used to identify the input in a JSON object.
func (i input) JsonKey() string {
	switch i {
	case AButton: return "a"
	case BButton: return "b"
	case XButton: return "x"
	case YButton: return "y"
	case RbButton: return "rb"
	case RtButton: return "rt"
	case LbButton: return "lb"
	case LtButton: return "lt"
	case StartButton: return "start"
	case SelectButton: return "select"
	case DPadLeftButton: return "left"
	case DPadRightButton: return "right"
	case DPadUpButton: return "up"
	case DPadDownButton: return "down"
	case LeftStickHorizontalLeft: return "lstick_left"
	case LeftStickHorizontalRight: return "lstick_right"
	case LeftStickVerticalUp: return "lstick_up"
	case LeftStickVerticalDown: return "lstick_down"
	case RightStickHorizontalLeft: return "rstick_left"
	case RightStickHorizontalRight: return "rstick_right"
	case RightStickVerticalUp: return "rstick_up"
	case RightStickVerticalDown: return "rstick_down"
	default:
		fatalf("Invalid input '%d'", i)
		// Note: This return isn't reached!
		return ""
	}
}

// Type returns the inputType for a given input. This simply changes the input
// prompt when requesting user input; it doesn't affect how the input is
// detected at all.
func (i input) Type() inputType {
	switch i {
	case AButton,
		BButton,
		XButton,
		YButton,
		RbButton,
		RtButton,
		LbButton,
		LtButton,
		StartButton,
		SelectButton,
		DPadLeftButton,
		DPadRightButton,
		DPadUpButton,
		DPadDownButton: return ButtonInput
	case LeftStickHorizontalLeft,
		LeftStickHorizontalRight,
		LeftStickVerticalUp,
		LeftStickVerticalDown,
		RightStickHorizontalLeft,
		RightStickHorizontalRight,
		RightStickVerticalUp,
		RightStickVerticalDown: return AxisInput
	default:
		fatalf("Invalid input '%d'", i)
		// Note: This return isn't reached!
		return inputType(0xff)
	}
}

// GetPrompt returns the input request prompt for the given gamepad input.
func (i input) GetPrompt() string {
	return fmt.Sprintf("%s %s", i.Type().String(), i.Name())
}

// float4digits wraps a float that is encoded in JSON with 4 decimal digits.
type float4digits struct {
	value float64
}

// MarshalJSON encodes the underlying float with 4 decimal digits.
func (f float4digits) MarshalJSON() ([]byte, error) {
	str := strconv.FormatFloat(f.value, 'f', 4, 64)
	return []byte(str), nil
}

// gamepadAxisToFloat converts a 16bit value to a float in the [-1, 1] range.
func gamepadAxisToFloat(axis int16) float64 {
	if axis > 0 {
		return float64(axis) / 32767.0
	} else if axis < 0 {
		return float64(axis) / 32768.0
	}
	return 0.0
}

// gamepadInput associates a given input with the read gamepad state.
type gamepadInput struct {
	Input input
	Gamepad gptools.Gamepad
}

// GetJsonEntry compares restGp, the state of the gamepad while at rest, with
// the supplied gamepadInput to detect which physical input should be used for
// the given gamepadInput.
//
// The format of the returned value depends on the type of the physical input.
// Every input has a "type", that specifies which array must be checked in the
// Gamepad, and an "index", that specifies which index in the given array must
// be checked.
//
// Button:
//
//     Object: {
//         "type": "button",
//         "index": <some-int>
//     }
//
//     Basic button that may be either pressed (1) or released (0).
//
// Axes:
//
//     Object: {
//         "type": "axis",
//         "index": <some-int>,
//         "min": <some-float>,
//         "max": <some-float>
//     }
//
//     An axis that may variate between [-1.0, 1.0]. In the Gamepad structure,
//     the axis value is read as an int16, which must be converted to the
//     desired range. Since different axes have different resting points (e.g.,
//     shoulder triggers may rest on -1.0), the button is defined as a range.
//
// Hats:
//
//     Object: {
//         "type": "mask",
//         "index": <some-int>,
//         "mask": <some-int>
//     }
//
//     A directional input encoded as a bitmask, with each bit representing a
//     given direction for the input.
func (gi gamepadInput) GetJsonEntry(restGp *gptools.Gamepad) (string, interface{}) {
    if bytes.Compare(restGp.Guid, gi.Gamepad.Guid) != 0 {
		fatalf("GUID does not match gamepad at rest!")
	}
	if len(restGp.Balls) != len(gi.Gamepad.Balls) {
		fatalf("Balls does not match gamepad at rest!")
	}
	if len(restGp.Axes) != len(gi.Gamepad.Axes) {
		fatalf("Axes does not match gamepad at rest!")
	}
	if len(restGp.Buttons) != len(gi.Gamepad.Buttons) {
		fatalf("Buttons does not match gamepad at rest!")
	}
	if len(restGp.Hats) != len(gi.Gamepad.Hats) {
		fatalf("Hats does not match gamepad at rest!")
	}
	if len(restGp.Name) != len(gi.Gamepad.Name) {
		fatalf("Name does not match gamepad at rest!")
	}
	if len(restGp.Guid) != len(gi.Gamepad.Guid) {
		fatalf("Guid does not match gamepad at rest!")
	}

	for i := range restGp.Buttons {
		if restGp.Buttons[i] != gi.Gamepad.Buttons[i] {
			input := struct {
				Type string
				Index int
			} {
				Type: "button",
				Index: i,
			}

			return gi.Input.JsonKey(), &input
		}
	}
	for i := range restGp.Hats {
		if restGp.Hats[i] != gi.Gamepad.Hats[i] {
			input := struct {
				Type string
				Index int
				Mask int
			} {
				Type: "mask",
				Index: i,
				Mask: int(gi.Gamepad.Hats[i]),
			}

			return gi.Input.JsonKey(), &input
		}
	}
	for i := range restGp.Axes {
		rest := gamepadAxisToFloat(restGp.Axes[i])
		cur := gamepadAxisToFloat(gi.Gamepad.Axes[i])

		if cur > 0 && math.Abs(cur - rest) > 0.65 {
			input := struct {
				Type string
				Index int
				Min json.Marshaler
				Max json.Marshaler
			} {
				Type: "axis",
				Index: i,
				Min: float4digits{value: rest},
				Max: float4digits{value: cur},
			}

			return gi.Input.JsonKey(), &input
		}
	}
	for i := range restGp.Balls {
		if restGp.Balls[i].X != gi.Gamepad.Balls[i].X {
			fatalf("Got trackball.x... but TrackBall isn't currently implemented!")
		} else if restGp.Balls[i].Y != gi.Gamepad.Balls[i].Y {
			fatalf("Got trackball.y... but TrackBall isn't currently implemented!")
		}
	}

	return "", nil
}

/* === Hotkey implementation ================================================ */

// Type for setting a value when a key is pressed.
type action struct{
	// isSet report whether the action has been set. Using channels would
	// cause issues if multiple events were generated without first
	// consuming every generated event.
	//
	// Therefore, isSet is accessed atomically.
	isSet int32
}

// Execute set the actionSet atomically.
func (a *action) Execute(pressed bool) {
	atomic.StoreInt32(&a.isSet, 1)
}

// Check atomically check if the action is set, clearing it afterwards.
func (a *action) Check() bool {
	return atomic.CompareAndSwapInt32(&a.isSet, 1, 0)
}

// Wait until the action was set, clearing the event.
func (a *action) Wait(checksPerSecond int) {
	delay := time.Second / time.Duration(checksPerSecond)
	for !a.Check() {
		time.Sleep(delay)
	}
}

/* === Application ========================================================== */

// isNoGamepad checks if an error indicates that there's no connected gamepad.
func isNoGamepad(err error) bool {
	merr, ok := err.(monitor.ErrorCode)
	return ok && merr == monitor.ErrNoGamepad
}

// isInvalidLength checks if an error indicates that there's no message.
func isInvalidLength(err error) bool {
	merr, ok := err.(monitor.ErrorCode)
	return ok && merr == monitor.ErrInvalidLength
}

// cloneGamepad makes a deep copy of the gamepad.
func cloneGamepad(gp *gptools.Gamepad) gptools.Gamepad {
	var newGp gptools.Gamepad

	data, err := json.Marshal(gp)
	if err != nil {
		fatalf("Couldn't clone the gamepad data: %+v", err)
	}

	err = json.Unmarshal(data, &newGp)
	if err != nil {
		fatalf("Couldn't clone the gamepad data: %+v", err)
	}

	return newGp
}

// run executes the main application.
func run(stop chan struct{}, m *monitor.Monitor, a *action, output io.Writer, gamepad string) {
	var gp gptools.Gamepad
	var data []byte
	var err error

	fmt.Print("Please connect a gamepad ")
	for {
		gp, data, err = gptools.GetLastGamepadData(m, data)
		if err == nil {
			break
		} else if isNoGamepad(err) || isInvalidLength(err) {
			fmt.Print(".")
			time.Sleep(time.Second / 3)
		} else {
			fatalf("Failed to detect any gamepads: %+v\n", err)
		}
	}

	// Clear any previously generated events
	a.Check()

	fmt.Printf(`
Name: %s
GUID: %s

!!! NOTE !!!

Some gamepads may report their axis incorrectly until they are initially
moved. For example, some gamepads may report a value of MAX_NEGATIVE_INT
for the should triggers on rest, but until the first movement those axis
are reported as 0.

So, to avoid issues, move every axis in the controller and press Enter.`,
		gp.Name,
		hex.EncodeToString(gp.Guid))
	a.Wait(default_check_per_second)

	data, err = gp.Update(m, data)
	if err != nil {
		fatalf("Couldn't read the gamepad at rest: %+v", err)
	}
	restGp := cloneGamepad(&gp)

	fmt.Printf("\n%s\nRefer to the image above for instructions of which button should be pressed.\n",
			gamepad)
	fmt.Println(`(Source: https://kenney-assets.itch.io/input-prompts-pixel-16)
Notes:
  - LT/RT are the shoulder triggers (button in the back);
  - The 'right arrow' is the start/pause button;
  - The 'left arrow' is the select button.
  - To skip any button, press enter without pressing anything on the gamepad`)

	var buttonList []gamepadInput
	for i := input(0); i < NumInputs; i++ {
		fmt.Printf("%s on the gamepad then press Enter ...\n", i.GetPrompt())

		a.Wait(default_check_per_second)
		data, err = gp.Update(m, data)
		if err != nil {
			fatalf("Couldn't read the gamepad for %s: %+v", input.Name, err)
		}

		button := gamepadInput {
			Input: i,
			Gamepad: cloneGamepad(&gp),
		}
		buttonList = append(buttonList, button)
	}

	// Encode the retrieved inputs to JSON.
	dict := make(map[string]interface{})
	for _, button := range buttonList {
		key, val := button.GetJsonEntry(&restGp)
		if val != nil {
			dict[key] = val
		}
	}
	if nameLen := len(gp.Name); nameLen > 0 && gp.Name[nameLen-1] == '\000' {
		gp.Name = gp.Name[:nameLen-1]
	}
	dict["name"] = gp.Name
	dict["guid"] = hex.EncodeToString(gp.Guid)

	enc := json.NewEncoder(output)
	err = enc.Encode(dict)
	if err != nil {
		fatalf("Failed to encode the gamepad data: %+v")
	}

	close(stop)
}

func main() {
	var pressedReturn action
	var libPath string
	var outputPath string
	var output io.WriteCloser
	var invert bool

	flag.Usage = func() {
		fmt.Print(`Generate a JSON description of the connected gamepad.

NOTE: This applicaton expects only a single gamepad to be connected.

Each detected input shall be encoded as one of the objects bellow, which
describes how this input may be checked on a
github.com/SirGFM/gfm-speedrun-overlay/gamepad-tools Gamepad object. The format
of the returned value depends on the type of the physical input. Every input
has a "type", that specifies which array must be checked in the Gamepad, and an
"index", that specifies which index in the given array must be checked.

Button:

    Object: {
        "type": "button",
        "index": <some-int>
    }

    Basic button that may be either pressed (1) or released (0).

Axes:

    Object: {
        "type": "axis",
        "index": <some-int>,
        "min": <some-float>,
        "max": <some-float>
    }

    An axis that may variate between [-1.0, 1.0]. In the Gamepad structure,
    the axis value is read as an int16, which must be converted to the
    desired range. Since different axes have different resting points (e.g.,
    shoulder triggers may rest on -1.0), the button is defined as a range.

Hats:

    Object: {
        "type": "mask",
        "index": <some-int>,
        "mask": <some-int>
    }

    A directional input encoded as a bitmask, with each bit representing a
    given direction for the input.

Usage:
`)
		flag.PrintDefaults()
	}
	flag.StringVar(&libPath, "library-path", "sdl-gamepad.dll", "(Windows only) Path to the gamepad monitoring library.")
	flag.StringVar(&outputPath, "output", "-", "File where the generated JSON should be written to. Use \"-\" for stdout.")
	flag.BoolVar(&invert, "invert", false, "Whether the colors of the printed gamepad example should be inverted.")
	flag.Parse()

	if outputPath == "-" {
		output = os.Stdout
	} else {
		// Avoid shadowing output by declaring error manually.
		var err error

		output, err = os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fatalf("Couldn't open the output file: %+v", err)
		}

		defer output.Close()
	}

	keyCfg := key_events.WatcherConfig {
		PoolPerSec: 10,
		OnKeyRelease: map[key_logger.Key]key_events.Action {
			key_logger.Return: &pressedReturn,
		},
	}
	keyWatcher := key_events.NewEventWatcher(keyCfg)
	defer keyWatcher.Close()

	cfg := monitor.Config {
		LibraryPath: libPath,
		StartTimeout: time.Second * 5,
	}

	m := monitor.New(cfg)
	defer m.Close()

	intHndlr := make(chan os.Signal, 1)
	signal.Notify(intHndlr, os.Interrupt)

	gamepad := gamepad_unicode
	if invert {
		gamepad = inverted_gamepad_unicode
	}

	stop := make(chan struct{}, 1)
	go run(stop, m, &pressedReturn, output, gamepad)

	select {
	case <-intHndlr:
		fmt.Print("Received signal! Stopping...\n")
	case <-stop:
		// Application finished successfully!
	}
}
