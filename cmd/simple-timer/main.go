package main

import (
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/res"
    "github.com/SirGFM/gfm-speedrun-overlay/web/server"
    "github.com/SirGFM/gfm-speedrun-overlay/web/timer"
    "os"
    "os/signal"
    "strings"
)

func main() {
    logger.RegisterDefault(logger.LogInfo, logger.LogDebug, os.Stdout, os.Stderr)

    srv := server.New()

    env := os.Getenv("RES_DIR")
    dirs := strings.Split(env, " ")
    cfg := res.Config {
        Dirs: dirs,
        DefaultPage: "index.html",
        DefaultExtension: ".html",
    }

    err := res.GetHandleFromConfig(srv, cfg)
    if err != nil {
        logger.Fatalf("Failed to add 'res' to the server: %+v", err)
    }
    err = srv.SetDefault(res.Prefix)
    if err != nil {
        logger.Fatalf("Failed to set 'res' as the default handler: %+v", err)
    }

    err = timer.GetHandle(srv)
    if err != nil {
        logger.Fatalf("Failed to add 'timer' to the server: %+v", err)
    }

    lst, err := srv.Listen("", 8080)
    if err != nil {
        logger.Fatalf("Failed to start server: %+v", err)
    }
    defer lst.Close()

    intHndlr := make(chan os.Signal, 1)
    signal.Notify(intHndlr, os.Interrupt)
    <-intHndlr
    logger.Debugf("Exiting...")
}
