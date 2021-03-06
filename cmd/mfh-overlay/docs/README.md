# MFH's General Racing Overlay

This time, let's do this right and start thing off documenting things as they are implemented!

## About the overlay

For most purposes, just launch the server and access http://localhost:8088. This page will list most things you need  in OBS, or should access in a browser.

For details on how to customize the overlay, check the [MFH's Restream Overlay](overlay.md) page.

## Dev guide

This overlay uses a few standard handlers and implement a custom one. Each handler may be accessed on the IP/port (8088, by default) where this service is running, followed by its handle (for example, `http://localhost:8088/res/`).

Each handler is described on [dev-guide](dev-guide.md).
