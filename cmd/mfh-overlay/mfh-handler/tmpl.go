package mfh_handler

import (
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/tmpl"
    "reflect"
    "strings"
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

// A custom field that will be added to an existing resource.
type customField struct {
    // The field's Key.
    key string
    // The field's Value.
    val interface{}
}

// Add a list of customField to a given object, `data`, that should have
// been decoded from a JSON, and thus should be a map of string to
// interfaces.
func addCustomFields(data interface{}, fields []customField) (interface{}, error) {
    var newData interface{}

    // To avoid modifying the original data, encode the data to a JSON
    // string and decode it back to a new, independent object.
    // Yes, this is ugly, but it gets the job done. *shrug*
    encData, err := json.Marshal(data)
    if err != nil {
        logger.Errorf("mfh-handler: Couldn't encode the resource to JSON: %+v", err)
        return nil, TmplCopyResource
    }
    err = json.Unmarshal(encData, &newData)
    if err != nil {
        logger.Errorf("mfh-handler: Couldn't decode the copy of the resource: %+v", err)
        return nil, TmplGetCopyResource
    }

    // Ensure that the assumption about this being a map of string to
    // interfaces is valid.
    val := reflect.ValueOf(newData)
    typ := val.Type()
    if typ.Kind() != reflect.Map {
        return nil, TmplResourceNotAMap
    } else if typ.Key().Kind() != reflect.String {
        return nil, TmplResourceNotStrKeys
    } else if typ.Elem().Kind() != reflect.Interface {
        return nil, TmplResourceNotInterfaceMap
    }

    for i := range fields {
        key := reflect.ValueOf(fields[i].key)
        el := reflect.ValueOf(fields[i].val)
        val.SetMapIndex(key, el)
    }

    return val.Interface(), nil
}

// Retrieve the data associated with a given MT title-card.
func (ctx *serverContext) getMTTCard(resource string) (interface{}, error) {
    var exp string

    if resource == "pl1-mttcard.go.html" {
        exp = "Player1SRL"
    } else if resource == "pl2-mttcard.go.html" {
        exp = "Player2SRL"
    } else {
        logger.Errorf("mfh-handler: Invalid title card page '%s'", resource)
        return nil, ExtraBadTitleCard
    }

    fields := ctx.getExtraData(nil)
    ctx.unlockExtraData()

    // Manually look through everything stored in the extra data
    // and retrieve only the expected field, in the format that
    // '/tmpl/mttcard.go.html' needs.
    for _, f := range fields {
        if f.key == exp {
            data := struct { PlayerSRL interface{} } {
                PlayerSRL: f.val,
            }
            return &data, nil
        }
    }

    return nil, ExtraBadTitleCard
}

// Retrieve the data associated with a given resource.
func (ctx *serverContext) Read(resource []string) (interface{}, error) {
    var customFields []customField

    // Store the original resource for error reporting.
    origRes := resource

    // Retrieve any custom, hard-coded values that the resource may have,
    // also normalizing its path.
    if len(resource) == 2 && resource[0] == "tmpl" {
        if strings.HasPrefix(resource[1], "1v1-") {
            field := customField {
                "Layout2v2",
                false,
            }
            customFields  = append(customFields, field)

            lastPart := resource[1][4:]
            resource = []string{"tmpl", lastPart}
        } else if strings.HasPrefix(resource[1], "2v2-") {
            field := customField {
                "Layout2v2",
                true,
            }
            customFields  = append(customFields, field)

            lastPart := resource[1][4:]
            resource = []string{"tmpl", lastPart}
        } else if strings.HasSuffix(resource[1], "mttcard.go.html") {
            return ctx.getMTTCard(resource[1])
        }
    }

    name := resource2name(resource)
    if _, ok := ctx.data[name]; !ok {
        logger.Errorf("mfh-handler: Couldn't find the resource associated with %+v", origRes)
        return nil, ResourceNotFound
    }

    // Set the win flag (and other static resources) from the extra data.
    customFields = ctx.getExtraData(customFields)
    defer ctx.unlockExtraData()

    // Add custom fields to the retrieved data.
    data := ctx.data[name]
    if len(customFields) > 0 {
        newData, err := addCustomFields(data, customFields)
        if err != nil {
            logger.Errorf("mfh-handler: Couldn't add custom fields to the resource %+v", origRes)
            return nil, err
        }

        data = newData
    }

    return data, nil
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
    if len(resource) == 2 && resource[0] == "tmpl" {
        if strings.HasSuffix(resource[1], "twitch-iframe.go.html") {
            return []string{"tmpl", "twitch-iframe.go.html"}, nil
        } else if strings.HasSuffix(resource[1], "mttcard.go.html") {
            return []string{"tmpl", "mttcard.go.html"}, nil
        } else if strings.HasPrefix(resource[1], "1v1-") || strings.HasPrefix(resource[1], "2v2-") {
            lastPart := resource[1][4:]
            return []string{"tmpl", lastPart}, nil
        }
    }

    return resource, nil
}
