# MFH's General Racing Overlay

This is a (hopefully) easily customizable overlay, based on my sorta-automated [Mystery Tournament 15](https://github.com/SirGFM/gfm-speedrun-overlay/tree/master/cmd/mt-overlay) overlay, to be used in various races/tournaments broadcast by the [Mystery Fun House](https://www.twitch.tv/mysteryfunhouse).

This tool is in no way official (as my version of the MT overlay wasn't), but it may hopefully help restreaming matches.

## Features

* Simplified interface for customizing the broadcast skin
* Up to 4 slots for racers, with:
    * Optional slot for audio icon
    * Racer name
    * Flag indicating that the racer finished the race
        * This doesn't track the player's position!
    * TANGIBLE PROGRESS
* Dashboard for controlling the race
* Built-in timer (for those like me that don't like using LiveSplit)
* Automatically shrinking text individually for fitting labels
    * Using this race as reference: https://www.twitch.tv/videos/344872865?collection=A4bL_e-gahV_JA
* Custom GIAN2OOPA overlay, with a growing turtle in the background
* Optional Mystery Tournament career title-card integration

## Possible future features

* Built-in SRL integration (mostly for weekly races, I guess...)

## Usage

Simply double-click mfh-overlay.exe and it should start. To finish the server, press Ctrl-C or close the terminal.

The home page (by default, `http://localhost:8088`) has links for all available pages

## Compiling

(Mainly for myself, but) to cross-compile for Windows, run:

```
GOOS=windows GOARCH=amd64 go build .
```

To enable the title-card module, be sure to download all required dependencies and to build with the `withmttcard` tag:

```
go get -u google.golang.org/api/sheets/v4 & # Run this in BG as it takes forever to download
export WAIT_PID=$!
go get github.com/pkg/errors
go get golang.org/x/net/html
go get -u golang.org/x/oauth2/google
wait ${WAIT_PID}
go get github.com/SirGFM/MTTitleCard
go build -tags withmttcard .
```
