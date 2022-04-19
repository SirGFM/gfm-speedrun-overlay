// `ram_store` stores arbitrary data for later retrieval. It's quite
// similar to `tmpl`, however there's no processing of the received
// objects, and any object type is accepted.
//
// This module is quite generic, and among other things could be used to
// implement communication between arbitrary elements (e.g., a key logger
// sending the keyboard status and a page displaying the keys)
//
// The accept methods are GET, DELETE, POST and PUT, and the last two are
// handled identically. GET returns the same Content Type specified in the
// POST/PUT request, and the URL identify the resource being accessed.
//
// For example, to implement a keyboard viewer one could send the keyboard
// status as the following JSON object to `/ram_store/keyboard-viewer`:
//
//     {
//         "A": true,
//         "B": false,
//         ...
//         "SpaceBar": false,
//         "Return": false,
//     }

package ram_store

import (
    "bytes"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "io"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
    "path"
    "sync"
)

const Prefix = "/ram_store"

// Build a new error
func newError(err error, res string, status int) error {
    return srv_iface.NewHttpError(err, "web"+Prefix, res, status)
}

// Associate the received data to its content type.
type data struct {
	// The received content type for this resource
	contentType string
	// The resource's content
	body *bytes.Buffer
}

// Context for the ram_store service
type rstore struct {
    // Every received resource.
    store map[string]data
    // Synchronize access to the context.
    rwmut sync.RWMutex
}

// Retrieve the path handled by `res`
func (*rstore) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*rstore) Dependencies() []string {
    return nil
}

// Pack data required to handle a request
type request struct {
    // The requesting server
    ctx *rstore
    // The writer for the HTTP response
    w http.ResponseWriter
    // The HTTP request
    req *http.Request
    // Relative path to the requested resource
    urlPath []string
}

func (r *request) ContentType() string {
    return r.req.Header.Get("Content-Type")
}

func (r *request) URLPath() string {
    return path.Join(r.urlPath...)
}

func (r *request) Read(buf []byte) (int, error) {
    return r.req.Body.Read(buf)
}

func (r *request) send() error {
    r.ctx.rwmut.RLock()
	defer r.ctx.rwmut.RUnlock()

	data, ok := r.ctx.store[r.URLPath()]
	if !ok {
        reason := "No resource was specified"
        return newError(nil, reason, http.StatusNotFound)
	}

    r.w.Header().Set("Content-Type", data.contentType)
    r.w.WriteHeader(http.StatusOK)

	// Since this entire function is locked for reading, body won't be
	// modified until there aren't any readers. To avoid modifying the
	// reader from multiple goroutine, the returned slice is consumed
	// manually.
	body := data.body.Bytes()
	for i := 0; i < len(body); {
		n, err := r.w.Write(body[i:])
		if err != nil {
			logger.Errorf("web%s: Failed to send the response: %+v (url: %+v)", Prefix, err, r.URLPath())
			break
		}

		i += n
	}

	return nil
}

func (r *request) recv() error {
    r.ctx.rwmut.Lock()
    defer r.ctx.rwmut.Unlock()

	data, ok := r.ctx.store[r.URLPath()]
	if !ok {
		data.body = bytes.NewBuffer(nil)
	}

	data.contentType = r.ContentType()
	data.body.Reset()
	_, err := io.Copy(data.body, r.req.Body)
	if err != nil {
        reason := "Failed to store the data"
        return newError(err, reason, http.StatusInternalServerError)
	}

	r.ctx.store[r.URLPath()] = data

    r.w.WriteHeader(http.StatusNoContent)
    return nil
}

func (r *request) del() error {
    r.ctx.rwmut.Lock()
    defer r.ctx.rwmut.Unlock()

	if _, ok := r.ctx.store[r.URLPath()]; !ok {
        reason := "Resource not found in the server"
        return newError(nil, reason, http.StatusNotFound)
	}
	delete(r.ctx.store, r.URLPath())

    r.w.WriteHeader(http.StatusNoContent)
    return nil
}

func (ctx *rstore) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    r := request {
        ctx: ctx,
        w: w,
        req: req,
        urlPath: urlPath,
    }

    switch req.Method {
    case http.MethodGet:
        return r.send()
    case http.MethodPost, http.MethodPut:
        return r.recv()
    case http.MethodDelete:
        return r.del()
    default:
        return newError(nil, "Invalid method: wanted one of GET, POST, PUT or DELETE", http.StatusMethodNotAllowed)
    }
}

// Close resources associated with the `rstore`
func (ctx *rstore) Close() {
    ctx.rwmut.Lock()
    defer ctx.rwmut.Unlock()

    for c := range ctx.store {
        delete(ctx.store, c)
    }
}

// Register a `ram_store` handler in the `Server`.
func GetHandle(srv srv_iface.Server) error {
    var ctx rstore

    ctx.store = make(map[string]data)

    srv.AddHandler(&ctx)
    return nil
}
