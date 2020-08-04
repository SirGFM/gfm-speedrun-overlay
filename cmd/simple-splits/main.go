package main

import (
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/res"
    "github.com/SirGFM/gfm-speedrun-overlay/web/server"
    "github.com/SirGFM/gfm-speedrun-overlay/web/splits"
    "log"
    "os"
    "os/signal"
    "path"
    "strings"
)

func main() {
    logger.RegisterDefault(logger.LogInfo, logger.LogDebug, os.Stdout, os.Stderr)
    log.Printf("Starting...\n")

    srv := server.New()

    env := os.Getenv("RES_DIR")
    dirs := strings.Split(env, " ")

    err := res.GetHandle(srv, dirs)
    if err != nil {
        log.Fatalf("Failed to add 'res' to the server: %+v", err)
    }

    appDir := path.Clean(path.Dir(os.Args[0]))
    splitsDir := path.Join(appDir, "splits")
    err = os.Mkdir(splitsDir, 0644)
    if err != nil && !os.IsExist(err) {
        log.Fatalf("Failed to create splits directory: %+v", err)
    }

    err = splits.GetHandle(srv, splitsDir)
    if err != nil {
        log.Fatalf("Failed to add 'splits' to the server: %+v", err)
    }

    lst, err := srv.Listen("", 8080)
    if err != nil {
        log.Fatalf("Failed to start server: %+v", err)
    }
    defer lst.Close()

    intHndlr := make(chan os.Signal, 1)
    signal.Notify(intHndlr, os.Interrupt)
    <-intHndlr
    log.Print("Exiting...")
}
