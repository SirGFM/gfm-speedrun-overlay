// `tmpl` serve extra resources, customized by a given struct. This is
// quite similar to `res`, but with a few extra steps. For instructions on
// how to write a customizable page, see the following resources:
//
// * https://golang.org/pkg/text/template/
// * https://golang.org/pkg/html/template/
//
// Differently from `res`, which only accepts GET requests, this service
// also accepts POST, PUT and DELETE, which are used to configure the page
// retrieved from a GET request. These four methods resolve in a call to
// the `DataCRUD` supplied when registering the service, which handles
// storing and retrieving custom data for pages.
//
// From the service's perspective, there's no difference between a POST and
// a PUT. The former calls `DataCRUD.Create` and the later calls
// `DataCRUD.Update`, however it's up for the implementation to
// differentiate between the two. Theoretically, POST should create a new
// entry and PUT should update and existing one.
//
// Note that the URL reported to `DataCRUD` starts on this service's
// prefix (thus, `DataReder.URLPath()[0] == tmpl`).

package tmpl

import (
    "bytes"
    "crypto/sha256"
    "github.com/SirGFM/gfm-speedrun-overlay/common"
    "html/template"
    "io"
    "io/ioutil"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "net/http"
    "os"
    "path"
    "sync"
)

const Prefix = "/tmpl"

// Build a new error
func newError(err error, res string, status int) error {
    return srv_iface.NewHttpError(err, "web"+Prefix, res, status)
}

// Functions to read data from a reader, alongside its content type.
type DataReader interface {
    // Retrieve the data type.
    ContentType() string
    // Retrieve the relative path to the requested resource.
    URLPath() []string
    // Implements Go's io.Reader, reading the request's Body.
    io.Reader
}

// Define the functions used by `tmpl` to manipulate all data needed to
// execute a template. The server properly locks itself up (for reading or
// for writing) before calling any of these functions.
type DataCRUD interface {
    // Create a new resource
    Create(resource []string, data DataReader) error
    // Retrieve the data associated with a given resource.
    Read(resource []string) (interface{}, error)
    // Update an already existing resource.
    Update(resource []string, data DataReader) error
    // Remove the resource.
    Delete(resource []string) error
    // Clean up the container, removing all associated resources.
    Close()
}

// Define a function that maps resources into other resources.
type Mapper interface {
    // Map a given resource into another resource
    Map(resource []string) ([]string, error)
}

// A template that has been parsed and cached.
type cachedPage struct {
    // Hash (SHA-256) of the associated page.
    hash []byte
    // The parsed page, ready to `Execute()` some data.
    page *template.Template
}

// Context for the tmpl service
type tmpl struct{
    // List of directories where resources may be located.
    entryPoints []string
    // Cache every accessed resource, for easier and faster access.
    pages map[string]cachedPage
    // Manages all data used by the resources.
    data DataCRUD
    // Possibly maps resources into other resources
    mapper Mapper
    // Synchronize access to the context.
    rwmut sync.RWMutex
}

// Retrieve the path handled by `res`
func (*tmpl) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*tmpl) Dependencies() []string {
    return nil
}

// Pack data required to handle a request
type request struct {
    // The requesting server
    ctx *tmpl
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

func (r *request) URLPath() []string {
    var arr []string
    return append(arr, r.urlPath...)
}

func (r *request) Read(buf []byte) (int, error) {
    return r.req.Body.Read(buf)
}

// Writer used to check for errors while executing a HTML template.
type dummyWriter struct{}
func (*dummyWriter) Write(data []byte) (int, error) {
    return len(data), nil
}

// Serve the requested file, using the data associated with its path to
// execute the template.
func (r *request) unsafeServeFile(filePath, ctype string) error {
    // Retrieve the data associated with the resource
    data, err := r.ctx.data.Read(r.urlPath)
    if err != nil {
        res := "Couldn't retrieve the associated data"
        return newError(err, res, http.StatusInternalServerError)
    }

    // Guaranteed to be OK,from a previous check.
    cache := r.ctx.pages[filePath]

    // Check if the template would parse successfully
    dummy := &dummyWriter{}
    err = cache.page.Execute(dummy, data)
    if err != nil {
        res := "Couldn't parse and execute the template"
        return newError(err, res, http.StatusInternalServerError)
    }

    // Finally, return the parsed data
    r.w.Header().Set("Content-Type", ctype)
    r.w.WriteHeader(http.StatusOK)
    cache.page.Execute(r.w, data)
    return nil
}

// Retrieve a file in one of the valid directories.
func (ctx *tmpl) getFile(filePath string) (*os.File, string, error) {
    for _, dir := range ctx.entryPoints {
        if file, ctype, err := common.OpenFile(dir, filePath); err == nil {
            return file, ctype, nil
        }
    }

    // NOTE: The file should have always been closed if err != nil
    reason := "Couldn't find the specified resource"
    return nil, "", newError(nil, reason, http.StatusNotFound)
}

func (r *request) get() error {
    // Convert the URL to a local path
    filePath := path.Join(r.urlPath[1:]...)
    if len(filePath) == 0 {
        reason := "No resource was specified"
        return newError(nil, reason, http.StatusNotFound)
    }

    // Retrieve the file if it's in any of the listed directories
    file, ctype, err := r.ctx.getFile(filePath)
    if err != nil {
        // Try to remap the file and use that instead
        if r.ctx.mapper == nil {
            return err
        }

        newUrl, err := r.ctx.mapper.Map(r.urlPath)
        if err != nil {
            reason := "Couldn't find the specified resource"
            return newError(err, reason, http.StatusNotFound)
        }

        filePath = path.Join(newUrl[1:]...)
        file, ctype, err = r.ctx.getFile(filePath)
        if err != nil {
            return err
        }
    }

    // Check if the file has been changed since its last use
    content, err := ioutil.ReadAll(file)
    file.Close()
    if err != nil {
        res := "Couldn't read the requested file"
        return newError(err, res, http.StatusInternalServerError)
    }
    hashArr := sha256.Sum256(content)
    hash := hashArr[:]

    // XXX: Don't defer unlocking the mutex since the lock may need to be
    // upgraded to a write lock if the page must be (re)cached.
    r.ctx.rwmut.RLock()

    cache, is_valid := r.ctx.pages[filePath]
    if is_valid {
        is_valid = (bytes.Compare(cache.hash, hash) == 0)
    }

    // If the page isn't in the cache or has been modified, update it!
    if !is_valid {
        // XXX: Update the lock for writing
        r.ctx.rwmut.RUnlock()
        r.ctx.rwmut.Lock()
        defer r.ctx.rwmut.Unlock()

        cache.hash = hash
        cache.page = template.New("")
        _, err = cache.page.Parse(string(content))
        if err != nil {
            res := "Couldn't parse the requested file"
            return newError(err, res, http.StatusInternalServerError)
        }

        r.ctx.pages[filePath] = cache
        return r.unsafeServeFile(filePath, ctype)
    } else {
        defer r.ctx.rwmut.RUnlock()
        return r.unsafeServeFile(filePath, ctype)
    }
}

func (r *request) create() error {
    r.ctx.rwmut.Lock()
    defer r.ctx.rwmut.Unlock()

    err := r.ctx.data.Create(r.urlPath, r)
    if err != nil {
        reason := "Failed to create (POST) the data"
        return newError(err, reason, http.StatusInternalServerError)
    }

    r.w.WriteHeader(http.StatusNoContent)
    return nil
}

func (r *request) update() error {
    r.ctx.rwmut.Lock()
    defer r.ctx.rwmut.Unlock()

    err := r.ctx.data.Update(r.urlPath, r)
    if err != nil {
        reason := "Failed to update (PUT) the data"
        return newError(err, reason, http.StatusInternalServerError)
    }

    r.w.WriteHeader(http.StatusNoContent)
    return nil
}

func (r *request) del() error {
    r.ctx.rwmut.Lock()
    defer r.ctx.rwmut.Unlock()

    err := r.ctx.data.Delete(r.urlPath)
    if err != nil {
        reason := "Failed to delete (DELETE) the data"
        return newError(err, reason, http.StatusInternalServerError)
    }

    r.w.WriteHeader(http.StatusNoContent)
    return nil
}

func (ctx *tmpl) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    r := request {
        ctx: ctx,
        w: w,
        req: req,
        urlPath: urlPath,
    }

    switch req.Method {
    case http.MethodGet:
        return r.get()
    case http.MethodPost:
        return r.create()
    case http.MethodPut:
        return r.update()
    case http.MethodDelete:
        return r.del()
    default:
        return newError(nil, "Invalid method: wanted one of GET, POST, PUT or DELETE", http.StatusMethodNotAllowed)
    }
}

// Close resources associated with the `tmpl`
func (ctx *tmpl) Close() {
    ctx.rwmut.Lock()
    defer ctx.rwmut.Unlock()

    for p := range ctx.pages {
        delete(ctx.pages, p)
    }
    ctx.data.Close()
}

// Register a `tmpl` handler in the `Server`.
func GetHandle(srv srv_iface.Server, dirs []string, data DataCRUD, mapper Mapper) error {
    var ctx tmpl

    cwd, err := os.Getwd()
    if err != nil {
        reason := "Failed to the the current working directory"
        return newError(err, reason, http.StatusInternalServerError)
    }

    // Get a list with possible resources
    ctx.entryPoints, err = common.ListToAbsolutePath(nil, path.Join(cwd, "tmpl"))
    if err != nil {
        reason := "Failed to create the list of resource directories"
        return newError(err, reason, http.StatusInternalServerError)
    }
    ctx.entryPoints, err = common.ListToAbsolutePath(ctx.entryPoints, dirs...)
    if err != nil {
        reason := "Failed to create the list of resource directories"
        return newError(err, reason, http.StatusInternalServerError)
    }
    ctx.data = data
    ctx.mapper = mapper
    ctx.pages = make(map[string]cachedPage)

    srv.AddHandler(&ctx)
    return nil
}
