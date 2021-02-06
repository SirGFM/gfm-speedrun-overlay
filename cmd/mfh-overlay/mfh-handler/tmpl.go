package mfh_handler

import (
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/tmpl"
)

// Converts a resource's path into a unique name, from its SHA-256. To
// convert it into a readable, slightly packed string, encode it as base64.
func resource2name(resource []string) string {
    var data []byte
    dgst := sha256.New()

    for i := range resource {
        // Go is dumb and says hash.Hash.Write "never returns an error",
        // so let's use it to our advantage...
        dgst.Write([]byte(resource[i]))
    }

    data = dgst.Sum(data)
    return base64.URLEncoding.EncodeToString(data)
}

// store the resource 'data' into the server at 'name'.
func (ctx *serverContext) store(name string, data tmpl.DataReader) error {
    var val interface{}

    dec := json.NewDecoder(data)
    err := dec.Decode(&val)
    if err != nil {
        logger.Errorf("mfh-handler: Failed to decode %+v's data: %+v", data.URLPath(), err)
        return BadJSONInput
    }
    ctx.data[name] = val
    ctx.update()

    return nil
}

// Create a new resource. Identical to "Update".
func (ctx *serverContext) Create(resource []string, data tmpl.DataReader) error {
    return ctx.store(resource2name(resource), data)
}

// Retrieve the data associated with a given resource.
func (ctx *serverContext) Read(resource []string) (interface{}, error) {
    name := resource2name(resource)
    if _, ok := ctx.data[name]; !ok {
        logger.Errorf("mfh-handler: Couldn't find the resource associated with %+v", resource)
        return nil, ResourceNotFound
    }

    return ctx.data[name], nil
}

// Update an already existing resource. Identical to "Create".
func (ctx *serverContext) Update(resource []string, data tmpl.DataReader) error {
    return ctx.store(resource2name(resource), data)
}

// Remove the resource.
func (ctx *serverContext) Delete(resource []string) error {
    name := resource2name(resource)
    if _, ok := ctx.data[name]; !ok {
        logger.Errorf("mfh-handler: No resource associated with %+v", resource)
        return ResourceNotFound
    }

    delete(ctx.data, name)
    return nil
}

// Map resources into themselves (as this doesn't need any fancy mapping).
func (ctx *serverContext) Map(resource []string) ([]string, error) {
    return resource, nil
}
