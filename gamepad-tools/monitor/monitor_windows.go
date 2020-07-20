// Windows bindings for `c-wrapper`.
package monitor

import (
    "golang.org/x/sys/windows"
    "log"
    "syscall"
    "unsafe"
)

type win_args struct {
    lib *windows.LazyDLL
    fnlib_start *windows.LazyProc
    fnlib_check_ready *windows.LazyProc
    fnlib_clean *windows.LazyProc
    fnlib_begin_request *windows.LazyProc
    fnlib_end_request *windows.LazyProc
    fnlib_get_node_size *windows.LazyProc
    fnlib_get_last_id *windows.LazyProc
    fnlib_get_num_gamepads *windows.LazyProc
    fnlib_get_data *windows.LazyProc
}

func (wa *win_args) clean() {
    windows.FreeLibrary(windows.Handle(wa.lib.Handle()))
}

func logErr(proc string, err error) {
    werr := err.(syscall.Errno)
    if werr != 0 {
        log.Printf("gamepad-tools/monitor: Syscall '%s' failed with err %d: %+v\n", proc, werr, err)
    }
}

func checkErr(err error) error {
    werr := err.(syscall.Errno)
    if werr == 0 {
        return nil
    }
    return err
}

func mergeErr(r1 uintptr, err error) uint32 {
    if checkErr(err) != nil {
        return uint32(ErrWinSyscall)
    }
    return uint32(r1)
}

func (wa *win_args)_start() ErrorCode {
    r1, _, err := wa.fnlib_start.Call()
    return ErrorCode(mergeErr(r1, err))
}

func (wa *win_args)_check_ready() ErrorCode {
    r1, _, err := wa.fnlib_check_ready.Call()
    return ErrorCode(mergeErr(r1, err))
}

func (wa *win_args)_clean() {
    _, _, err := wa.fnlib_clean.Call()
    logErr("clean", err)
}

func (wa *win_args)_begin_request() ErrorCode {
    r1, _, err := wa.fnlib_begin_request.Call()
    return ErrorCode(mergeErr(r1, err))
}

func (wa *win_args)_end_request() ErrorCode {
    r1, _, err := wa.fnlib_end_request.Call()
    return ErrorCode(mergeErr(r1, err))
}

func (wa *win_args)_get_node_size() uint32 {
    r1, _, err := wa.fnlib_get_node_size.Call()
    return mergeErr(r1, err)
}

func (wa *win_args)_get_last_id() uint32 {
    r1, _, err := wa.fnlib_get_last_id.Call()
    return mergeErr(r1, err)
}

func (wa *win_args)_get_num_gamepads() uint32 {
    r1, _, err := wa.fnlib_get_num_gamepads.Call()
    return mergeErr(r1, err)
}

func (wa *win_args)_get_data(idx uint32 , data []byte) {
    c_idx := uintptr(idx)
    c_data := uintptr(unsafe.Pointer(&data[0]))

    _, _, err := wa.fnlib_get_data.Call(c_idx, c_data)
    logErr("get_data", err)
}

func (wa *win_args) get_proc(proc string) *windows.LazyProc {
    f := wa.lib.NewProc(proc)
    err := f.Find()
    if err != nil {
        log.Panicf("gamepad-tools/monitor: Failed to get function '%s': %+v\n", proc, err)
    }
    return f
}

func get_function(conf Config) remote_func {
    args := &win_args{}

    args.lib = windows.NewLazyDLL(conf.LibraryPath)
    if args.lib == nil {
        log.Panic("gamepad-tools/monitor: Failed to open the monitor DLL\n")
    }
    err := args.lib.Load()
    if err != nil {
        log.Panicf("gamepad-tools/monitor: Failed to load the monitor DLL: %+v\n", err)
    }

    args.fnlib_start = args.get_proc("start")
    args.fnlib_check_ready = args.get_proc("check_ready")
    args.fnlib_clean = args.get_proc("clean")
    args.fnlib_begin_request = args.get_proc("begin_request")
    args.fnlib_end_request = args.get_proc("end_request")
    args.fnlib_get_node_size = args.get_proc("get_node_size")
    args.fnlib_get_last_id = args.get_proc("get_last_id")
    args.fnlib_get_num_gamepads = args.get_proc("get_num_gamepads")
    args.fnlib_get_data = args.get_proc("get_data")

    return remote_func {
        start: args._start,
        check_ready: args._check_ready,
        clean: args._clean,
        begin_request: args._begin_request,
        end_request: args._end_request,
        get_node_size: args._get_node_size,
        get_last_id: args._get_last_id,
        get_num_gamepads: args._get_num_gamepads,
        get_data: args._get_data,
        args: args,
    }
}
