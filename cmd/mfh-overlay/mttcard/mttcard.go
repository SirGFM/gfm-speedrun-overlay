package mttcard

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    mttcard_cfg "github.com/SirGFM/MTTitleCard/config"
    mttcard_srv "github.com/SirGFM/MTTitleCard/page"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "net/http"
    "net/url"
    "path"
    "path/filepath"
)

const Prefix = "/mttcard"

type server struct {
    handler mttcard_srv.PageServer
}

// Retrieve the path handled by `serverContext`.
func (*server) Prefix() string {
    return Prefix
}

// List every other service used by this handler.
func (*server) Dependencies() []string {
    return nil
}

// Clean up the container, removing all associated resources.
func (*server) Close() {
}

// Implement the server interface, so serverContext may report when it changes
func (ctx *server) Handle(w http.ResponseWriter, req *http.Request, urlPath []string) error {
    var err error
    var uri string

    if len(urlPath) >= 1 {
        uri = "/" + path.Join(urlPath[1:]...)
    } else {
        uri = "/"
    }

    req.URL, err = url.ParseRequestURI(uri)
    if err != nil {
        return err
    }

    ctx.handler.ServeHTTP(w, req)
    return nil
}

// Allocs a new MT Title Card handler
func New() srv_iface.Handler {
    var err error

    config := mttcard_cfg.GetDefault()
    // The service looks for a style in {config.ServiceUri}/style.css
    config.ServiceUri = "/res/style/mttcard"

    config.CredentialFile, err = filepath.Abs("./credentials.json")
    if err != nil {
        logger.Errorf("mttcard: Couldn't open credentials file: %+v", err)
        return nil
    }

    config.TokenFile, err = filepath.Abs("./token.json")
    if err != nil {
        logger.Errorf("mttcard: Couldn't open token file: %+v", err)
        return nil
    }

    err = mttcard_cfg.LoadConfig(config)
    if err != nil {
        logger.Errorf("mttcard: Couldn't load the configuration: %+v", err)
        return nil
    }

    mttc, err := mttcard_srv.New()
    if err != nil {
        logger.Errorf("mttcard: Couldn't start the MT Title Card sub-server: %+v", err)
        return nil
    }

    ctx := server {
        handler: mttc,
    }
    return &ctx
}
