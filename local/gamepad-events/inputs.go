package gamepad_events

import (
	"fmt"
)

// InputName describe every input in a gamepad handled by the application.
//
// Use the image bellow to reference where in a controller each given button should be:
//
//     Lt             Rt
//     Lb             Rb
//
//       Select   Start
//
//   @ Left           X
//   | Stick        Y   A
//                    B
//
//   ^              @ Right
// <   > DPad       | Stick
//   v
type InputName uint32
const (
	AButton InputName = iota
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
	NumInputNames
)

// Name returns a human-readable name for the input.
func (in InputName) Name() string {
	switch in {
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
		panic(fmt.Sprintf("%s: Invalid InputName '%d'", pkg, in))
	}
}

// key returns the name used to identify the input as a JSON key.
func (in InputName) key() string {
	switch in {
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
		panic(fmt.Sprintf("%s: Invalid InputName '%d'", pkg, in))
	}
}

// MarshalText implements the encoding.TextMarshaler, so an InputName may be
// used as the key in a JSON object.
func (in InputName) MarshalText() ([]byte, error) {
	return []byte(in.key()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler, so an InputName may
// be used as the key in a JSON object.
func (in *InputName) UnmarshalText(text []byte) error {
	for str, i := string(text), InputName(0); i < NumInputNames; i++ {
		if str == i.key() {
			*in = i
			return nil
		}
	}

	panic(fmt.Sprintf("%s: Invalid InputName '%s'", pkg, string(text)))
}
