package main

import (
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/res"
    "github.com/SirGFM/gfm-speedrun-overlay/web/server"
    "os"
    "os/signal"
    "path/filepath"
)

// Change the current working directory to the binary's directory. Panics on failure.
func changeToAppDir() {
    appPath := filepath.Clean(os.Args[0])
    appPath, err := filepath.Abs(appPath)
    if err != nil {
        logger.Fatalf("Couldn't retrieve the absolute path to the application: %+v", err)
    }

    appDir := filepath.Dir(appPath)
    err = os.Chdir(appDir)
    if err != nil {
        logger.Fatalf("Couldn't cd into the application's directory: %+v", err)
    }
}

type voidCloser interface {
    Close()
}

type Config struct {
    Host string
    Port int
    DefaultPage string
}

func startServer(cfg Config) voidCloser {
    srv := server.New()

    resCfg := res.Config {
        DefaultPage: cfg.DefaultPage,
        DefaultExtension: ".html",
    }

    err := res.GetHandleFromConfig(srv, resCfg)
    if err != nil {
        logger.Fatalf("Failed to add 'res' to the server: %+v", err)
    }
    err = srv.SetDefault(res.Prefix)
    if err != nil {
        logger.Fatalf("Failed to set 'res' as the default handler: %+v", err)
    }

    lst, err := srv.Listen(cfg.Host, cfg.Port)
    if err != nil {
        logger.Fatalf("Failed to start server: %+v", err)
    }
    logger.Infof("Started running server on %s:%d", cfg.Host, cfg.Port)

    return lst
}

func main() {
    logger.RegisterDefault(logger.LogInfo, logger.LogDebug, os.Stdout, os.Stderr)

    changeToAppDir()

    cfg := Config {
        Port: 8000,
        DefaultPage: "index.html",
    }
    closer := startServer(cfg)
    defer closer.Close()

    intHndlr := make(chan os.Signal, 1)
    signal.Notify(intHndlr, os.Interrupt)
    <-intHndlr
    logger.Infof("Exiting...")
}
