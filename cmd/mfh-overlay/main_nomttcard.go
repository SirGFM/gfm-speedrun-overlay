// +build !withmttcard

package main

import (
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
)

func AddMTTCard(s srv_iface.Server) {
    _ = s
}
