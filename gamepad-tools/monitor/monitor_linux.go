// Linux bindings for `c-wrapper`.
package monitor

// #cgo LDFLAGS: -L./c-wrapper -lsdl-gamepad
// #include <stdint.h>
// extern int32_t start();
// extern int32_t check_ready();
// extern void clean();
// extern int32_t begin_request();
// extern int32_t end_request();
// extern uint32_t get_node_size();
// extern uint32_t get_last_id();
// extern uint32_t get_num_gamepads();
// extern void get_data(uint32_t idx, uint8_t *data);
import "C"

import (
    "unsafe"
)

type linux_args struct {}

func (*linux_args) clean() {
    // Do nothing
}

func lx_start() ErrorCode {
    return ErrorCode(C.start())
}

func lx_check_ready() ErrorCode {
    return ErrorCode(C.check_ready())
}

func lx_clean() {
    C.clean()
}

func lx_begin_request() ErrorCode {
    return ErrorCode(C.begin_request())
}

func lx_end_request() ErrorCode {
    return ErrorCode(C.end_request())
}

func lx_get_node_size() uint32 {
    return uint32(C.get_node_size())
}

func lx_get_last_id() uint32 {
    return uint32(C.get_last_id())
}

func lx_get_num_gamepads() uint32 {
    return uint32(C.get_num_gamepads())
}

func lx_get_data(idx uint32 , data []byte) {
    c_idx := C.uint32_t(idx)
    c_data := (*C.uint8_t)(unsafe.Pointer(&data[0]))
    C.get_data(c_idx, c_data)
}

func get_function(conf Config) remote_func {
    return remote_func {
        start: lx_start,
        check_ready: lx_check_ready,
        clean: lx_clean,
        begin_request: lx_begin_request,
        end_request: lx_end_request,
        get_node_size: lx_get_node_size,
        get_last_id: lx_get_last_id,
        get_num_gamepads: lx_get_num_gamepads,
        get_data: lx_get_data,
        args: &linux_args{},
    }
}
