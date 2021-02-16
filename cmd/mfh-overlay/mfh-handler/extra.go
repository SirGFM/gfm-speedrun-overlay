package mfh_handler

import (
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "io"
    "reflect"
)

// Lock the extra data for reading and retrieve the current data, appending
// each field to `in`.
// `ctx.unlockExtraData()` must be called (or deferred!) after a call to
// this function.
func (ctx *serverContext) getExtraData(in []customField) []customField {
    ctx.extra.m.RLock()

    for i := range ctx.extra.parsed {
        in = append(in, ctx.extra.parsed[i])
    }
    return in
}

// Unlock the extra data, locked by a previous `ctx.getExtraData()` call.
func (ctx *serverContext) unlockExtraData() {
    ctx.extra.m.RUnlock()
}

// Retrieve the extra data encoded as JSON.
func (ctx *serverContext) getExtraJSON(w io.Writer) {
    var err error

    ctx.extra.m.RLock()
    if ctx.extra.obj != nil {
        enc := json.NewEncoder(w)
        err = enc.Encode(ctx.extra.obj)
    }
    ctx.extra.m.RUnlock()

    if err != nil {
        logger.Errorf("%s/overlay-extras: Failed to encode the response: %+v", Prefix, err)
    }
}

// Store the extra data into the server, overwriting any previous values.
func (ctx *serverContext) putExtraJson(r io.Reader) error {
    var data interface{}
    var parsed []customField

    dec := json.NewDecoder(r)
    err := dec.Decode(&data)
    if err != nil {
        logger.Errorf("%s/overlay-extras: Failed to decode the extra data: %+v", Prefix, err)
        return BadJSONInput
    }

    // Read each entry from this JSON into a customField[].
    val := reflect.ValueOf(data)
    typ := val.Type()
    if typ.Kind() != reflect.Map {
        return ExtraDataNotAMap
    } else if typ.Key().Kind() != reflect.String {
        return ExtraDataNotStrKeys
    } else if typ.Elem().Kind() != reflect.Interface {
        return ExtraDataNotInterfaceMap
    }

    m := val.MapRange()
    for m.Next() {
        key := m.Key()
        if key.Type().Kind() != reflect.String {
            return ExtraDataNotStrKeys
        }
        el := m.Value()
        if el.Type().Kind() != reflect.Interface {
            return ExtraDataNotInterfaceMap
        }

        newField := customField {
            key.String(),
            el.Interface(),
        }
        parsed = append(parsed, newField)
    }

    // Save the object to the context.
    ctx.extra.m.Lock()
    ctx.extra.obj = data
    ctx.extra.parsed = parsed
    ctx.extra.m.Unlock()
    ctx.update()

    return nil
}
