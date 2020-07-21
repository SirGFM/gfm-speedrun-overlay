package main

import (
    gptools "github.com/SirGFM/gfm-speedrun-overlay/gamepad-tools"
    "github.com/SirGFM/gfm-speedrun-overlay/gamepad-tools/monitor"
    "log"
    "os"
    "os/signal"
    "time"
)

func main() {
    cfg := monitor.Config {
        LibraryPath: "sdl-gamepad.dll",
        StartTimeout: time.Second * 5,
    }

    m := monitor.New(cfg)
    defer m.Close()

    intHndlr := make(chan os.Signal, 1)
    signal.Notify(intHndlr, os.Interrupt)

    var gp gptools.Gamepad
    var data []byte
    var err error

    hasGamepad := false
    running := true
    for running {
        select {
        case <-intHndlr:
            log.Print("Received signal! Stopping...\n")
            running = false
        case <-time.After(time.Second):
            if !hasGamepad {
                gp, data, err = gptools.GetLastGamepadData(m, data)
                if err == nil && len(gp.Name) != 0 {
                    hasGamepad = true
                }
            } else {
                data, err = gp.Update(m, data)
            }

            if err != nil {
                log.Printf("Error: %+v\n", err)
            } else {
                log.Printf("Gamepad: %+v\n", gp)
            }
        }
    }
}
