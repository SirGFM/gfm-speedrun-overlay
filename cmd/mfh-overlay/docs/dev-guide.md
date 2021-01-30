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
