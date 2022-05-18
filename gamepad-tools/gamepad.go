// gamepad consumes data retrieved by `monitor` into the Go struct
// `Gamepad`.
//
// After initializing a `*monitor.Monitor`, the options to parse gamepads
// are:
//
//   - `ListGamepads()`: retrieve a list of `Gamepad`s, alongside their
//           names formatted with their handles;
//   - `GetLastGamepadData()`: retrieve the `Gamepad` for the last gamepad
//           to generate any event;
//   - `GetGamepadData()`: retrieve the `Gamepad` for a specific gamepad.
//           Prefer using `ListGamepads()` to retrieve valid `Handles()`,
//           and then call this for one of the handles in the list(if using
//           the `Gamepad` from `ListGamepads()` isn't an option).
//
// The caller is responsible for managing the `*monitor.Monitor` used by
// this module, as well as ensuring that the same monitor is all supplied
// to update `Gamepad`s.

package gamepad_tools

import (
    "fmt"
    "github.com/SirGFM/gfm-speedrun-overlay/gamepad-tools/monitor"
)

// Tracks "ball axes" (a 2D trackpad, maybe?). See
// http://wiki.libsdl.org/SDL_JoystickGetBall?highlight=%28%5CbCategoryJoystick%5Cb%29%7C%28CategoryEnum%29%7C%28CategoryStruct%29
type GamepadBall struct {
    X, Y int32
}

// A gamepad representation.
type Gamepad struct {
    // The internal handle, within the monitor library.
    idx uint
    // List of ball axes in the Joystick.
    Balls []GamepadBall
    // List of regular axes in the Joystick.
    Axes []int16
    // List of buttons in the Joystick.
    Buttons []uint8
    // List of hats in the Joystick.
    Hats []uint8
    // Joystick's name.
    Name string
    // Device's GUID
    Guid []uint8
}

// Check whether the given gamepad has any field set.
func (gp *Gamepad)checkEmpty() error {
    if len(gp.Balls) == 0 && len(gp.Axes) == 0 && len(gp.Buttons) == 0 &&
            len(gp.Hats) == 0 && len(gp.Name) == 0 && len(gp.Guid) == 0 {
        return monitor.ErrNoGamepad
    }
    return nil
}

// Parse `data` into the `Gamepad`.
func (gp *Gamepad)parse(data []byte) {
    numBalls := data[0]
    numAxes := data[1]
    numButtons := data[2]
    numHats := data[3]

    dec := hostOrder()
    nameLen := dec.Uint32(data[4:8])
    guidLen := dec.Uint32(data[8:12])

    data = data[12:]

    gp.Balls = gp.Balls[:0]
    for ; numBalls > 0; numBalls-- {
        gpb := GamepadBall{}
        gpb.X = int32(dec.Uint32(data[:4]))
        gpb.Y = int32(dec.Uint32(data[4:8]))
        data = data[8:]

        gp.Balls = append(gp.Balls, gpb)
    }

    gp.Axes = gp.Axes[:0]
    for ; numAxes > 0; numAxes-- {
        val := int16(dec.Uint16(data[:2]))
        data = data[2:]

        gp.Axes = append(gp.Axes, val)
    }

    gp.Buttons = gp.Buttons[:0]
    if numButtons > 0 {
        gp.Buttons = append(gp.Buttons, data[:numButtons]...)
        data = data[numButtons:]
    }

    gp.Hats = gp.Hats[:0]
    if numHats > 0 {
        gp.Hats = append(gp.Hats, data[:numHats]...)
        data = data[numHats:]
    }

    if nameLen > 0 {
        gp.Name = string(data[:nameLen])
        data = data[nameLen:]
    } else {
        gp.Name = ""
    }

    gp.Guid = gp.Guid[:0]
    if guidLen > 0 {
        gp.Guid = append(gp.Guid, data[:guidLen]...)
    }
}

// Create a new `Gamepad` from the supplied data.
func newGamepad(idx uint, data []byte) Gamepad {
    gp := Gamepad{
        idx: idx,
    }
    gp.parse(data)

    return gp
}

// Update the given gamepad, using and expanding data as necessary.
//
// Any error is of type `monitor.ErrorCode`.
func (gp *Gamepad)Update(m *monitor.Monitor, data []byte) ([]byte, error) {
    data, err := m.GetGamepadData(gp.idx, data)
    if err == nil {
        gp.parse(data)
        err = gp.checkEmpty()
    }
    return data, err
}

// Check for errors getting gamepad data from the Monitor and return either
// a parsed `Gamepad` or the `error`.
func wrapNewGamepad(idx uint, data []byte, err error) (Gamepad, []byte, error) {
    var gp Gamepad
    if err != nil {
        return gp, data, err
    }
    gp = newGamepad(idx, data)
    return gp, data, gp.checkEmpty()
}

// Retrieve the data for a given gamepad, using and expanding data as
// necessary, and parse it into a `Gamepad`.
//
// Any error is of type `monitor.ErrorCode`.
func GetGamepadData(m *monitor.Monitor, idx uint, data []byte) (Gamepad, []byte, error) {
    data, err := m.GetGamepadData(idx, data)
    return wrapNewGamepad(idx, data, err)
}

// Retrieve the data for last gamepad to receive any event, using and
// expanding data as necessary, and parse it into a `Gamepad`.
//
// Any error is of type `monitor.ErrorCode`. In particular, if no gamepad
// was found, this function return `monitor.ErrNoGamepad`.
func GetLastGamepadData(m *monitor.Monitor, data []byte) (Gamepad, []byte, error) {
    idx, data, err := m.GetLastGamepadData(data)
    return wrapNewGamepad(idx, data, err)
}

// List every gamepad currently connected and monitored.
func ListGamepads(m *monitor.Monitor) ([]string, []Gamepad, error) {
    inum, err := m.GetNumGamepads()
    if err != nil {
        return nil, nil, err
    }

    var data []byte
    var nameList []string
    var gpList []Gamepad

    num := uint(inum)
    for i := uint(0); i < num; i++ {
        var gp Gamepad

        gp, data, err = GetGamepadData(m, i, data)
        if err == monitor.ErrNoGamepad {
            continue
        } else if err != nil {
            return nil, nil, err
        }

        name := fmt.Sprintf("%d - %s", i, gp.Name)
        nameList = append(nameList, name)
        gpList = append(gpList, gp)
    }

    return nameList, gpList, nil
}
