// +build withmttcard

package main

import (
    "github.com/SirGFM/gfm-speedrun-overlay/cmd/mfh-overlay/mttcard"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
)

func AddMTTCard(s srv_iface.Server) {
    mttcCtx := mttcard.New()
    if mttcCtx != nil {
        s.AddHandler(mttcCtx)
    }
}
