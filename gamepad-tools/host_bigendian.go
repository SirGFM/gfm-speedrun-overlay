// +build *be

// gamepad consumes data retrieved by `monitor` into the Go struct
// `Gamepad`.
//
// Since data is retrieved from `monitor/c-wrapper` in host order,
// `host_bigendian.go` implements a single function, `hostOrder()`, that
// return `binary.BigEndian`.
//
// This should only be compiled in big-endian architectures.

package gamepad_tools

import (
    "encoding/binary"
)

// Hack to retrieve the host's `binary.ByteOrder`.
func hostOrder() binary.ByteOrder {
	return binary.BigEndian
}
