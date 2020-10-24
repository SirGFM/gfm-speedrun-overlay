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

## Extra module

Handle/entry-point: `mfh-handler/overlay-extras`

This is sort of a hack... This module was implemented before it was possible to retrieve the raw data for a dynamic page. So, to add more information to the overlay page without having to resend everything, one could simply add arbitrary data to this module.

This is mainly used to register that a player finished the goal, and optionally their completition time. These information are send in a `Player<num>Won` field, a boolean, and optionally a `Player<num>Time`, the player's completition time as a human-readable string. For example, the following would be a valid object to register that player's 1 and 4 have completed a 4-way race, while only reporting the time of 16 minutes and 10 seconds for player 1:

```json
{
    "Player1Won": true,
    "Player1Time": "16:10.000",
    "Player4Won": true
}
```

This object is automatically, and independently, added to `overlay.go.html`.

## Template module

Handle/entry-point: `tmpl`

The standard "template" module. This is quite similar to the resource module, but the retrieved resources are dynamic HTML pages, instead of static ones.

Currently, its API is documented only on the source code itself. For more information, see [web/tmpl/tmpl.go](../../../web/tmpl/tmpl.go)

The data for dynamic pages is stored in `mfh-handler`.

### Twitch stream page

Resource: `tmpl/<identifier>twitch-iframe.go.html`

One of the template's pages. This page allows configuring the twitch channel to be accessed, as well as the stream's resolution. By using different identifiers for each player on OBS (say, `pl1-twitch-iframe.go.html`, `pl2-twitch-iframe.go.html` etc), it's possible to automatically change the source without having to manually change the Browser Source.

This page must be configured with an object with a `TwitchUsername` field, the channel's name (so, for `https://www.twitch.tv/mysteryfunhouse` it would be only `mysteryfunhouse`), and two optional dimentions, on fields `Width` and `Height`. If not specified, the dimensions default to 852x480.

### Mystery Tournament Career title card page

Resource: `tmpl/pl<num>-mttcard.go.html`

Another of the template pages. This pages allows easily retrieving a racer's title card. Since this is only intended to be used when simple 1v1 racer, it only accepts two possible URLs: `tmpl/pl1-mttcard.go.html` and `tmpl/pl2-mttcard.go.html`. The server will look through its extra data for the field `Player1SRL` or `Player2SRL`, respectively, and map the field required by the requested page into an object with a single field `PlayerSRL`.

A page with both racers may be accessed on `http://localhost:8088/res/racer-cards.html`.

## Mystery Tournament Career title card

Handle/entry-point: `mttcard/<srl-user>`

Simple port of my [MT Title Card](https://github.com/SirGFM/MTTitleCard) into the interface required by this system.

This module retrieve's the racing stats for a given [SpeedRunsLive](https://www.speedrunslive.com/) user both from the site as from the "Mystery Tournament Careers" spreadsheet, organizing in a simple HTML page. The title card style may be customized by editing `res/style/mttcard/style.css`.
