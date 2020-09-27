package match

import (
    "encoding/json"
    "github.com/SirGFM/gfm-speedrun-overlay/logger"
    "github.com/SirGFM/gfm-speedrun-overlay/web/tmpl"
    srv_iface "github.com/SirGFM/gfm-speedrun-overlay/web/server/common"
    "strconv"
    "sync/atomic"
    "time"
)

// Errors returned by this package.
type errorCode uint
const (
    // Featured not implemented (yet?)
    NotImplemented errorCode = iota
    // Expected JSON input
    WantJSONData
    // Failed to decode the JSON input
    BadJSONInput
    // UPDATE does not match any resource
    InvalidUpdateURL
    // GET does not match any resource
    InvalidGetURL
    // Invalid player name for the match's winner
    InvalidWinner
    // GET couldn't map the resource into a real resource
    InvalidMapURL
    // Invalid player name for the tangible progress reaction
    InvalidTPPlayer
)

// Implement the `error` interface for `errorCode`.
func (e errorCode) Error() string {
    switch e {
    case NotImplemented:
        return "mt-server: Featured not implemented (yet?)"
    case WantJSONData:
        return "mt-server: Expected JSON input"
    case BadJSONInput:
        return "mt-server: Failed to decode the JSON input"
    case InvalidUpdateURL:
        return "mt-server: UPDATE does not match any resource"
    case InvalidGetURL:
        return "mt-server: GET does not match any resource"
    case InvalidWinner:
        return "mt-server: Invalid player name for the match's winner"
    case InvalidMapURL:
        return "mt-server: GET couldn't map the resource into a real resource"
    case InvalidTPPlayer:
        return "mt-server: Invalid player name for the tangible progress reaction"
    default:
        return "mt-server: Unknown"
    }
}

// Describes player
type player struct {
    // The player's name
    Name string
    // Whether the player has won
    Won bool
    // Predictions for this player's victory
    Predictions int
}

// Describes a on-going match
type match struct {
    // The left player
    Player1 player
    // The right player
    Player2 player
    // Predictions toward player 1 winning
    Predictions string
    // The challonge round for this match
    Round string
}

// Set the winner for the given match
func (m *match) setWinner(url string) error {
    var err error

    if url == "player1" {
        m.Player1.Won = true
        m.Player2.Won = false
    } else if url == "player2" {
        m.Player1.Won = false
        m.Player2.Won = true
    } else if url == "none" {
        m.Player1.Won = false
        m.Player2.Won = false
    } else {
        err = InvalidWinner
    }

    return err
}

// Container for the on-going race
type mtServerData struct {
    // The only 2-way match, or the top match.
    Top match
    // The bottom 2-way match, in a 4-way race.
    Bottom match
    // The game being played.
    Game string
    // The match's goal (e.g., Beat the game)
    Goal string
    // The game's platform (e.g., NES).
    Platform string
    // Who submitted this game.
    SubbedBy string
    // Whether this is a 4-way race (i.e., whether bottom should be used).
    Is4Way bool
    // Whether the page should auto-update
    AutoUpdate bool
    // Whether the built-in timer should be used
    UseTimer bool
    // Whether the tangible progress button may be used
    UseTangibleProgress bool
}

// The context that store page's data.
type serverContext struct {
    // The on-going race
    data mtServerData
    // Last time the structure was updated.
    lastUpdate time.Time `json:"-"`
    // Whether the tangible progress should be shown for any player
    tangibleProgress [4]uint32 `json:"-"`
}

// Set the flag that shows the "tangible progress" reaction
func (ctx *serverContext) setTangibleProgress(url string) error {
    var addr *uint32

    if url == "player1" {
        addr = &(ctx.tangibleProgress[0])
    } else if url == "player2" {
        addr = &(ctx.tangibleProgress[1])
    } else if url == "player3" {
        addr = &(ctx.tangibleProgress[2])
    } else if url == "player4" {
        addr = &(ctx.tangibleProgress[3])
    } else {
        return InvalidWinner
    }

    atomic.StoreUint32(addr, 1)
    return nil
}

// Get the flag that shows the "tangible progress" reaction, clearing it
// afterwards.
func (ctx *serverContext) getClearTangibleProgress() interface{} {
    var tp [len(ctx.tangibleProgress)]bool

    for i := range ctx.tangibleProgress {
        addr := &(ctx.tangibleProgress[i])
        tp[i] = (0 != atomic.SwapUint32(addr, 0))
    }

    s := struct { TangibleProgress [len(ctx.tangibleProgress)]bool } {
        TangibleProgress: tp,
    }
    return &s
}

// Clone the context, filling up some placeholder text
func (ctx *serverContext) getData() mtServerData {
    var data mtServerData

    data = ctx.data

    // Fill the empty values with some pre-defined values.
    if len(data.Top.Player1.Name) == 0 {
        data.Top.Player1.Name = "data.Top.Player1.Name"
    }
    if len(data.Top.Player2.Name) == 0 {
        data.Top.Player2.Name = "data.Top.Player2.Name"
    }
    if len(data.Top.Round) == 0 {
        data.Top.Round = "data.Top.Round"
    }
    if len(data.Bottom.Player1.Name) == 0 {
        data.Bottom.Player1.Name = "data.Bottom.Player1.Name"
    }
    if len(data.Bottom.Player2.Name) == 0 {
        data.Bottom.Player2.Name = "data.Bottom.Player2.Name"
    }
    if len(data.Bottom.Round) == 0 {
        data.Bottom.Round = "data.Bottom.Round"
    }
    if len(data.Game) == 0 {
        data.Game = "data.Game"
    }
    if len(data.Goal) == 0 {
        data.Goal = "data.Goal"
    }
    if len(data.Platform) == 0 {
        data.Platform = "data.Platform"
    }
    if len(data.SubbedBy) == 0 {
        data.SubbedBy = "data.SubbedBy"
    }

    // Set player predictions manually
    topPred, err := strconv.Atoi(data.Top.Predictions)
    if err == nil {
        data.Top.Player1.Predictions = topPred
        data.Top.Player2.Predictions = 100 - topPred
    } else {
        logger.Errorf("mt-server: Failed to Top's predictions: %+v", err)
        data.Top.Player1.Predictions = 50
        data.Top.Player2.Predictions = 50
    }

    botPred, err := strconv.Atoi(data.Bottom.Predictions)
    if err == nil {
        data.Bottom.Player1.Predictions = botPred
        data.Bottom.Player2.Predictions = 100 - botPred
    } else {
        logger.Errorf("mt-server: Failed to Bottom's predictions: %+v", err)
        data.Bottom.Player1.Predictions = 50
        data.Bottom.Player2.Predictions = 50
    }

    return data
}

// Create a new resource
func (ctx *serverContext) Create(resource []string, data tmpl.DataReader) error {
    return NotImplemented
}

// Retrieve the data associated with a given resource.
func (ctx *serverContext) Read(url []string) (interface{}, error) {
    if len(url) < 2 {
        return nil, InvalidGetURL
    }

    switch url[1] {
    case "index-1v1.go.html":
        clone := ctx.getData()
        clone.Is4Way = false
        return clone, nil
    case "index-2v2.go.html":
        clone := ctx.getData()
        clone.Is4Way = true
        return clone, nil
    case "index.go.html":
        return ctx.getData(), nil
    default:
        return nil, InvalidGetURL
    }
}

// Update the server's data
func (ctx *serverContext) updateData(data tmpl.DataReader) error {
    var got mtServerData

    dec := json.NewDecoder(data)
    err := dec.Decode(&got)
    if err != nil {
        logger.Errorf("mt-server: Failed to decode UPDATE data: %+v", err)
        return BadJSONInput
    }

    ctx.data = got
    return nil
}

// Update an already existing resource.
func (ctx *serverContext) Update(resource []string, data tmpl.DataReader) error {
    if data.ContentType() != "application/json" {
        return WantJSONData
    } else if len(resource) < 2 {
        return InvalidUpdateURL
    }

    var err error

    switch resource[1] {
    case "index-1v1.go.html",
         "index-2v2.go.html",
         "index.go.html":
        err = ctx.updateData(data)
    case "top-win":
        if len(resource) < 3 {
            return InvalidUpdateURL
        }
        err = ctx.data.Top.setWinner(resource[2])
    case "bottom-win":
        if len(resource) < 3 {
            return InvalidUpdateURL
        }
        err = ctx.data.Bottom.setWinner(resource[2])
    case "tp":
        if len(resource) < 3 {
            return InvalidUpdateURL
        }
        // XXX: Exit early to avoid setting the last update time
        return ctx.setTangibleProgress(resource[2])
    default:
        return InvalidUpdateURL
    }


    if err == nil {
        ctx.lastUpdate = time.Now()
    }
    return err
}

// Remove the resource.
func (ctx *serverContext) Delete(resource []string) error {
    return NotImplemented
}

// Clean up the container, removing all associated resources.
func (ctx *serverContext) Close() {
}

// Map a given resource into another resource
func (ctx *serverContext) Map(url []string) ([]string, error) {
    if len(url) < 2 {
        return nil, InvalidMapURL
    }

    switch url[1] {
    case "index-1v1.go.html",
         "index-2v2.go.html":
        return []string{url[0], "index.go.html"}, nil
    default:
        return nil, InvalidMapURL
    }
}

// Ensure serverContext implements both `tmpl.DataCRUD` and
// `srv_iface.Server`, so it may be used as a server and for templating.
type Context interface {
    tmpl.DataCRUD
    tmpl.Mapper
    srv_iface.Handler
}

// Retrieve a new data server
func New() Context {
    ctx := serverContext {}
    return &ctx
}
