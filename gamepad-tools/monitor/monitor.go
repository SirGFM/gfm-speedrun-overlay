// monitor wraps `c-wrapper` functionalities into a Golang API. This is
// mostly a low-level mapping, so prefer using `gamepad-tools` instead!
//
// To start monitoring gamepad inputs, call `New()` to retrieve a new
// `*Monitor`. Be sure call (or, ideally, defer) `Close()` when done with
// the object.
//
// To retrieve gamepad data, either call `GetLastGamepadData()` to retrieve
// data from the last gamepad that generated any event, or call
// `GamepadCount()` to retrieve the number of gamepad handles available,
// and then call `GetGamepadData()` on any of the handles. Beware that only
// a few of those handles will be associated with actual gamepads!
//
// Every error from this package is of type `ErrorCode`, which is simply a
// mapping of the C enumeration for errors.

package monitor

import (
    "log"
    "time"
)

// `error` used by this package. Should usually be converted into an actual
//`error` interface by calling `(ErrorCode).err()`, which already converts
// values that do not represent errors into `nil`.
type ErrorCode int32
const (
    // Operation successfull
    ErrOk ErrorCode = iota
    // Waiting initialization
    ErrWaiting
    // Start failed
    ErrNotStarted
    // Failed to configure the mutex
    ErrSdlStartMutex
    // Failed to start the background thread
    ErrSdlStartThread
    // Failed to initialize SDL2
    ErrSdlInit
    // Failed to enable Joystick events
    ErrSdlInitBgJoystick
    // Failed to filter SDL2's signal detection
    ErrSdlInitNoSignals
    // Failed to initialize SDL2's Event subsystem
    ErrSdlInitEvent
    // Failed to initialize SDL2's Joystick subsystem
    ErrSdlInitJoystick
    // Timed out waiting for events
    ErrSdlWait
    // Failed to synchronize the library
    ErrLock
    // Failed to unlock library's synchronization
    ErrUnlock
    // Failed to synchronize the library's mainloop
    ErrWaitLock
    // DEPRECATED: Failed to expand the list of nodes
    ErrAddJoystickEvent1
    // Failed to add a new joystick
    ErrAddJoystickEvent2
    // Failed to allocate expanded list
    ErrExpandList
    // Failed to open SDL2 Joystick device
    ErrOpenJoystick
    // Failed to allocate internal abstration for SDL2 Joysticks
    ErrAllocJoyNode
    // Failed to unlock the library's mainloop synchronization
    ErrWaitUnlock
    // Failed to retrieve the joystick index
    ErrGetJoystickID
    // Failed to expand the list of nodes
    ErrAddJoystickEvent3
)

const (
    // The Monitor interface was closed and shouldn't be used anymore
    ErrFinished ErrorCode = iota + 1000
    // Failed to call function in Windows' DLL
    ErrWinSyscall
    // Failed to retrieve the payload length
    ErrInvalidLength
    // No gamepad has been detected yet
    ErrNoGamepad
)

// `err()` converts an `ErrorCode` into an `error` interface, returning
// `nil` for values that aren't failures (e.g., `ErrOk`).
func (e ErrorCode) err() error {
    switch e {
        case ErrOk:
            return nil
        case ErrWaiting:
            return nil
        default:
            return e
    }
}

// `Error()` implements the `error` interface for `ErrorCode`.
func (e ErrorCode) Error() string {
    switch e {
        case ErrOk:
            return "Operation successfull"
        case ErrWaiting:
            return "Waiting initialization"
        case ErrNotStarted:
            return "Start failed"
        case ErrSdlStartMutex:
            return "Failed to configure the mutex"
        case ErrSdlStartThread:
            return "Failed to start the background thread"
        case ErrSdlInit:
            return "Failed to initialize SDL2"
        case ErrSdlInitBgJoystick:
            return "Failed to enable Joystick events"
        case ErrSdlInitNoSignals:
            return "Failed to filter SDL2's signal detection"
        case ErrSdlInitEvent:
            return "Failed to initialize SDL2's Event subsystem"
        case ErrSdlInitJoystick:
            return "Failed to initialize SDL2's Joystick subsystem"
        case ErrSdlWait:
            return "Timed out waiting for events"
        case ErrLock:
            return "Failed to synchronize the library"
        case ErrUnlock:
            return "Failed to unlock library's synchronization"
        case ErrWaitLock:
            return "Failed to synchronize the library's mainloop"
        case ErrAddJoystickEvent1:
            return "DEPRECATED: Failed to expand the list of nodes"
        case ErrAddJoystickEvent2:
            return "Failed to add a new joystick"
        case ErrExpandList:
            return "Failed to allocate expanded list"
        case ErrOpenJoystick:
            return "Failed to open SDL2 Joystick device"
        case ErrAllocJoyNode:
            return "Failed to allocate internal abstration for SDL2 Joysticks"
        case ErrWaitUnlock:
            return "Failed to unlock the library's mainloop synchronization"
        case ErrWinSyscall:
            return "Failed to call function in Windows' DLL"
        case ErrFinished:
            return "The Monitor interface was closed and shouldn't be used anymore"
        case ErrInvalidLength:
            return "Failed to retrieve the payload length"
        case ErrNoGamepad:
            return "No gamepad has been detected yet"
        case ErrGetJoystickID:
            return "Failed to retrieve the joystick index"
        case ErrAddJoystickEvent3:
            return "Failed to expand the list of nodes"
        default:
            return "Unknown error"
    }
}

// Configuration for the library initialization.
type Config struct {
    // Windows only: path to the DLL.
    LibraryPath string
    // How long should the monitor wait until the background thread
    // finishes initializing.
    StartTimeout time.Duration
}

// Interface for OS-specific details.
type extra_args interface {
    clean()
}

// Redirect calls to the Operating System's `c-wrapper` through a Golang
// interface. Also stores any extra data that should be released when done
// with the library.
type remote_func struct {
    start func() ErrorCode
    check_ready func() ErrorCode
    clean func()
    begin_request func() ErrorCode
    end_request func() ErrorCode
    get_node_size func() uint32
    get_last_id func() uint32
    get_num_gamepads func() uint32
    get_data func(idx uint32 , data []byte)
    args extra_args
}

// Accessor for the `c-wrapper` mappings.
type Monitor struct {
    funcs remote_func
    closed bool
}

// Dummy function that returns ErrFinished
func dummyErrorCode() ErrorCode {
    return ErrFinished
}

// Dummy function that returns 0
func dummyUint32() uint32 {
    return 0
}

// Dummy function that does nothing
func dummyData(idx uint32 , data []byte) {
}

// Dummy function that does nothing
func dummy() {
}

// Release every resource associated with the Monitor.
func (m *Monitor) Close() {
    if !m.closed {
        m.funcs.clean()
        m.funcs.args.clean()

        m.funcs.start = dummyErrorCode
        m.funcs.check_ready = dummyErrorCode
        m.funcs.clean = dummy
        m.funcs.begin_request = dummyErrorCode
        m.funcs.end_request = dummyErrorCode
        m.funcs.get_node_size = dummyUint32
        m.funcs.get_last_id = dummyUint32
        m.funcs.get_num_gamepads = dummyUint32
        m.funcs.get_data = dummyData
        m.funcs.args = nil

        m.closed = true
    }
}

// Retrieve the number of gamepad handles available.
func (m *Monitor)GetNumGamepads() (int, error) {
    rv := m.funcs.begin_request()
    if rv != ErrOk {
        return 0, rv
    }
    ret := int(m.funcs.get_num_gamepads())
    rv = m.funcs.end_request()
    return ret, rv.err()
}

// Retrieve the data for a given gamepad, expanding data as necessary.
// This function ignores synchronizing access to the library, since that
// should be resolved before calling it.
func (m *Monitor)unsafeGetGamepadData(idx uint, data []byte) (uint, []byte, error) {
    // Expand the supplied buffer if necessary
    l := int(m.funcs.get_node_size())
    if l == 0 {
        return idx, data, ErrInvalidLength
    }

    if len(data) == 0 {
        data = append(data, 0)
    }
    for len(data) < l {
        data = append(data, data...)
    }

    // Retrieve the data from the library
    data = data[:l]
    m.funcs.get_data(uint32(idx), data)
    return idx, data, nil
}

// Retrieve the data for a given gamepad, expanding data as necessary, and
// synchronizing the library.
func (m *Monitor)GetGamepadData(idx uint, data []byte) ([]byte, error) {
    rv := m.funcs.begin_request()
    if rv != ErrOk {
        return data, rv
    }

    _, data, err := m.unsafeGetGamepadData(idx, data)

    rv = m.funcs.end_request()
    if rv != ErrOk {
        // Overwrite any error with errors unlocking the library
        err = rv
    }
    return data, err
}

// Retrieve the data for last gamepad to receive any event, expanding data
// as necessary, and synchronizing the library.
func (m *Monitor)GetLastGamepadData(data []byte) (uint, []byte, error) {
    rv := m.funcs.begin_request()
    if rv != ErrOk {
        return 0, data, rv
    }

    idx := uint(m.funcs.get_last_id())
    idx, data, err := m.unsafeGetGamepadData(idx, data)

    rv = m.funcs.end_request()
    if rv != ErrOk {
        // Overwrite any error with errors unlocking the library
        err = rv
    }
    return idx, data, err
}

// Initialize a `Monitor` from the configuration.
//
// This panics on error!
func New(conf Config) *Monitor {
    funcs := get_function(conf)

    rv := funcs.start()
    if rv != ErrOk {
        log.Panicf("gamepad-tools/monitor: Failed to initialized the library: %+d\n", rv)
    }

    rv = ErrWaiting
    timeout := conf.StartTimeout
    for timeout > 0 && rv == ErrWaiting {
        dt := time.Millisecond * 100
        rv = funcs.check_ready()
        if rv != ErrOk {
            time.Sleep(dt)
        }
        timeout -= dt
    }
    if rv != ErrOk {
        funcs.clean()
        funcs.args.clean()
        log.Panic("gamepad-tools/monitor: Library didn't initialized in a timely manner\n")
    }

    return &Monitor {
        funcs: funcs,
        closed: false,
    }
}
