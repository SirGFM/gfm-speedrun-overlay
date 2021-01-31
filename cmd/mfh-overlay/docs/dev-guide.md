# Dev guide

This overlay uses a few standard handlers and implement a custom one. Each handler may be accessed on the IP/port (8088, by default) where this service is running, followed by its handle (for example, `http://localhost:8088/res/`).

## Resource handler

Handle/entry-point: `res`

The standard resource module, which allows fetching arbitrary files within a specific directory (such as static pages and images).

This is the service's default handler, and as such may be accessed without specifying the handle. Also, the service defines the page `index.html` as its default page. So, accessing these three URLs are equivalent for this module:

* `http://localhost:8088/res/index.html`
* `http://localhost:8088/index.html`
* `http://localhost:8088/`

## Timer handler

Handle/entry-point: `timer`

The standard timer module. Defines a timer that runs within the service and may be accessed and controlled through an API.

Currently, its API is documented only on the source code itself. For more information, see [web/timer/timer.go](../../../web/timer/timer.go)

## MFH handler

Handle/entry-point: `mfh-handler`

Custom handler that shall implement all the custom operations required by this service.

### Last update sub-handler

Handle/entry-point: `mfh-handler/last-update`

Returns an object with the last date, in milliseconds since the Unix epoch, when anything on `mfh-handler` was updated. This information may be used to automatically reload a page if anything changed.

The date is represented by the field `Date`, an integer, of the returned object. For example, the date 2021-05-29T23:21:57+00:00 would be returned as:

```json
{
    "Date": 1622330517000
}
```

To use this, simply include [auto\_reload.js)](../res/script/auto_reload.js) in the desired page and call `auto_reload.update()` periodically.

### Popup sub-handler

Handle/entry-point: `mfh-handler/last-update`

Manages a list of IDs and their timeouts, in milliseconds. These may be used to report to a page that a dynamic element (controlled by a script) should be temporarily displayed.

A single element is represented by the field `Id`, a string, and `Timeout`, an integer. Regardless of whether the server returns multiple elements or a single one, it always returns these elements in the array `Elements`. For example, a response to show the IDs `TangibleProgress-p1` for 2.5s and `BigBoints-p2` for 5s would be represented by the following object:

```json
{
    "Elements": [
        {
            "Id": "TangibleProgress-p1",
            "Timeout": 2500
        },
        {
            "Id": "BigBoints-p2",
            "Timeout": 3000
        }
    ]
}
```

To use this, simply include [popup.js)](../res/script/popup.js) in the desired page, define a CSS class `hidden` for setting an object's visibility to invisible, call `popup.update()` periodically (to check if any object should be shown) and call `popup.show()` to report to the server that an object should be shown.
