package mfh_handler

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "github.com/SirGFM/gfm-speedrun-overlay/web/tmpl"
    "sync"
    "sync/atomic"
    "time"
)

// Errors returned by this package.
type errorCode uint
const (
    // Failed to decode the JSON input
    BadJSONInput errorCode = iota
    // Didn't find any resource associated with the given path
    ResourceNotFound
    // Couldn't create a JSON-encoded copy of the requested resource
    TmplCopyResource
    // Couldn't create a JSON-decoded copy of the requested resource
    TmplGetCopyResource
    // Can't add custom field to non-map value
    TmplResourceNotAMap
    // Can't add custom field to a map of non-string keys
    TmplResourceNotStrKeys
    // Can't add custom field to a map of non-interface elements
    TmplResourceNotInterfaceMap
    // Can't read custom fields from a non-map value
    ExtraDataNotAMap
    // Can't read custom fields from a map of non-string keys
    ExtraDataNotStrKeys
    // Can't read custom fields from a map of non-interface elements
    ExtraDataNotInterfaceMap
    // Invalid title card page
    ExtraBadTitleCardPage
    // Couldn't find the requested title card
    ExtraBadTitleCard
)

// Implement the `error` interface for `errorCode`.
func (e errorCode) Error() string {
    switch e {
    case BadJSONInput:
        return "mfh-handler: Failed to decode the JSON input"
    case ResourceNotFound:
        return "mfh-handler: Didn't find any resource associated with the given path"
    case TmplCopyResource:
        return "Couldn't create a JSON-encoded copy of the requested resource"
    case TmplGetCopyResource:
        return "Couldn't create a JSON-decoded copy of the requested resource"
    case TmplResourceNotAMap:
        return "mfh-handler: Can't add custom field to non-map value"
    case TmplResourceNotStrKeys:
        return "Can't add custom field to a map of non-string keys"
    case TmplResourceNotInterfaceMap:
        return "Can't add custom field to a map of non-interface elements"
    case ExtraDataNotAMap:
        return "Can't read custom fields from a non-map value"
    case ExtraDataNotStrKeys:
        return "Can't read custom fields from a map of non-string keys"
    case ExtraDataNotInterfaceMap:
        return "Can't read custom fields from a map of non-interface elements"
    case ExtraBadTitleCardPage:
        return "mfh-handler: Invalid title card page"
    case ExtraBadTitleCard:
        return "mfh-handler: Couldn't find the requested title card"
    default:
        return "mfh-handler: Unknown"
    }
}

// Simple timer for tracking when the service was updated.
type timer struct {
    // The stored time.
    t time.Time
    // Synchronize access to the time.Time.
    m sync.RWMutex
    // (hack to) try locking for writing without blocking.
    tryWrite uint32
}

// Describe a popup element that should be temporarily displayed.
type popup struct {
    // The ID of the HTML element to be shown and later hidden.
    Id string
    // How long, in milliseconds, until the element it hidden.
    Timeout int64
}

// List of popup elements, with synchronized access
type popupList struct {
    // The list of elements
    list []popup
    // Synchronize access to the popups list
    m sync.Mutex
}

// Extra data, configurable from the dashboard and appended to
// /tmpl/overlay.go.html.
type extraData struct {
    // The raw, decoded received object.
    obj interface{}
    // The list of elements.
    parsed []customField
    // Synchronize access to the data list.
    m sync.RWMutex
}

// The context that store page's data.
type serverContext struct {
    // Data used to customize pages, keyed by the SHA-256 of the path.
    data map[string]interface{}
    // Last time the structure was updated.
    lastUpdate timer
    // List of elements that should be temporarily displayed.
    popups popupList
    // Extra elements that may be added to the overlay.
    extra extraData
}

// Retrieve the last updated time
func (ctx *serverContext) getLastUpdate() int64 {
    ctx.lastUpdate.m.RLock()
    ms := ctx.lastUpdate.t.Unix()
    ctx.lastUpdate.m.RUnlock()
    return ms * 1000
}

// Update the context, so watching clients may automatically refresh.
func (ctx *serverContext) update() {
    // (tl;dr: This is overkill)
    //
    // Ensure that only a single writer may ever acquire the write lock
    // (using atomic.CompareAndSwapUint32 to only allow a single goroutine
    // to enter the conditional at a time) and then lock it for writing.
    if atomic.CompareAndSwapUint32(&ctx.lastUpdate.tryWrite, 0, 1) {
        ctx.lastUpdate.m.Lock()
        ctx.lastUpdate.t = time.Now()
        ctx.lastUpdate.m.Unlock()
        atomic.StoreUint32(&ctx.lastUpdate.tryWrite, 0)
    }
}

// Push a popup element into the list of elements to be temporarily shown.
func (ctx *serverContext) pushPopupElement(id string, timeout int64) {
    ctx.popups.m.Lock()
    defer ctx.popups.m.Unlock()

    p := popup {
        Id: id,
        Timeout: timeout,
    }

    ctx.popups.list = append(ctx.popups.list, p)
}

// Get the list of popup elements, emptying the list in the process.
func (ctx *serverContext) getPopupList() []popup {
    ctx.popups.m.Lock()
    defer ctx.popups.m.Unlock()

    tmp := ctx.popups.list

    ctx.popups.list = ctx.popups.list[:0]
    return tmp
}

// Clean up the container, removing all associated resources.
func (ctx *serverContext) Close() {
}

// `srv_iface.Server`, so it may be used as a server and for templating.
type Context interface {
    tmpl.DataCRUD
    tmpl.Mapper
    srv_iface.Handler
}

// Retrieve a new data server
func New() Context {
    ctx := serverContext {}
    ctx.data = make(map[string]interface{})
    return &ctx
}
