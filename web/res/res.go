// `res` serve extra resources, which are packed as static files located
// in a `res/` directory.

package res

import (
    "github.com/SirGFM/gfm-speedrun-overlay/common"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
    "os"
    "path"
    "time"
)

const Prefix = "/res"

// Build a new error
func newError(err error, res string, status int) error {
    return srv_iface.NewHttpError(err, "web"+Prefix, res, status)
}

type res struct{
    // List of directories where resources may be located
    entryPoints []string
    // Page used when the requested URL is empty
    defaultPage string
    // Default extension, if the requested URL doesn't have any **and**
    // no file was found.
    defaultExtension string
}

// Retrieve the path handled by `res`
func (*res) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*res) Dependencies() []string {
    return nil
}

// Pack data required to handle a request
type fileRequest struct {
    // The requesting server
    r *res
    // The writer for the HTTP response
    w http.ResponseWriter
    // The HTTP request
    req *http.Request
    // Relative path to the requested resource
    filePath string
}

// Server the requested file, returning whether or not it was found.
func (ctx fileRequest) serveFile() bool {
    // Server the file if it's in any of the listed directories
    for _, dir := range ctx.r.entryPoints {
        if f, ctype, err := common.OpenFile(dir, ctx.filePath); err == nil {
            defer f.Close()

            var modtime time.Time
            if fi, err := f.Stat(); err == nil {
                modtime = fi.ModTime()
            }

            ctx.w.Header().Set("Content-Type", ctype)
            http.ServeContent(ctx.w, ctx.req, "", modtime, f)
            return true
        }
    }

    return false
}

func (r *res) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    if req.Method != "GET" {
        reason := "Invalid method: wanted GET"
        return newError(nil, reason, http.StatusMethodNotAllowed)
    }

    // Convert the URL to a local path
    filePath := path.Join(urlPath[1:]...)
    if len(filePath) == 0 {
        if len(r.defaultPage) == 0 {
            reason := "No resource was specified"
            return newError(nil, reason, http.StatusNotFound)
        } else {
            filePath = r.defaultPage
        }
    }

    ctx := fileRequest {
        r: r,
        w: w,
        req: req,
        filePath: filePath,
    }

    if ctx.serveFile() {
        // File served successfully
        return nil
    }

    // If the file doesn't have an extension, try to serve the file again
    // after adding the default one.
    if len(path.Ext(filePath)) == 0 && len(r.defaultExtension) != 0 {
        ctx.filePath = filePath + r.defaultExtension
        if ctx.serveFile() {
            // File served successfully
            return nil
        }
    }

    // Couldn't find the file
    reason := "Couldn't find file '" + filePath + "'"
    return newError(nil, reason, http.StatusNotFound)
}

// Close resources associated with the `res` (i.e, nothing)
func (*res) Close() {
}

// Configure the `res` server.
type Config struct {
    // List of directories where resources may be located
    Dirs []string
    // Page used when the requested URL is empty
    DefaultPage string
    // Default extension, if the requested URL doesn't have any **and**
    // no file was found.
    DefaultExtension string
}

// Register a `res` handler in the `Server`.
func GetHandle(srv srv_iface.Server, dirs []string) error {
    cfg := Config {
        Dirs: dirs,
    }

    return GetHandleFromConfig(srv, cfg)
}

// Register a `res` handler in the `Server`. The resource is configured
// based on the supplied `cfg`.
func GetHandleFromConfig(srv srv_iface.Server, cfg Config) error {
    var r res

    cwd, err := os.Getwd()
    if err != nil {
        reason := "Failed to the the current working directory"
        return newError(err, reason, http.StatusInternalServerError)
    }

    // Get a list with possible resources
    r.entryPoints, err = common.ListToAbsolutePath(nil, path.Join(cwd, "res"))
    if err != nil {
        reason := "Failed to create the list of resource directories"
        return newError(err, reason, http.StatusInternalServerError)
    }
    r.entryPoints, err = common.ListToAbsolutePath(r.entryPoints, cfg.Dirs...)
    if err != nil {
        reason := "Failed to create the list of resource directories"
        return newError(err, reason, http.StatusInternalServerError)
    }
    r.defaultPage = cfg.DefaultPage
    r.defaultExtension = cfg.DefaultExtension

    srv.AddHandler(&r)
    return nil
}
