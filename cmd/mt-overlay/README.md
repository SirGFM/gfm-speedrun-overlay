# MFH's MT Overlay

Simply double-click on mt-overlay.exe and it should start. To finish the server, press Ctrl-C or close the terminal.

## Current pages

* Homepage page: http://localhost:8088
    * Has links for the available pages
* Configuration page: http://localhost:8088/config
* Dashboard page: http://localhost:8088/dashboard
* Automatic overlay: http://localhost:8088/tmpl/index.go.html
* 1v1 overlay: http://localhost:8088/tmpl/index-1v1.go.html
* 2v2 overlay: http://localhost:8088/tmpl/index-2v2.go.html

The dashboard has buttons for displaying the winner flag and controlling the built-in timer.

## About the timer

If you choose to use the built-in timer, it runs regardless of the page.

So, **it's safe to refresh the page with the timer running!**"
However, you can't stop the the layout server...

## OBS configuration

Add a browser source to  one of:

* http://localhost:8088/tmpl/index.go.html
* http://localhost:8088/tmpl/index-1v1.go.html
* http://localhost:8088/tmpl/index-2v2.go.html

`index.go.html` won't allow switching from a  2v2 setting back to a 1v1 setting,
so maybe avoid that...

If you enable auto-update in the configuration (recommended!), you can pretty
much just ignore this source afterwards. Just be careful with cached stuff...
If the page looks wrong, try to hit the "refresh without cache" button.

Something will most likely break the layout (for example, one of the player's name may be too long).
In that case, configure the CCS in OBS to edit whatever is broken.

These are the IDs you may need to change (which should be self-explanatory), with their default values:

```
// Name for top players/players in 1v1
p#pl1-name {
    font-size: 42px;
}
p#pl2-name {
    font-size: 42XXpx;
}

// Used only in 1v1 races
p#pl1-predictions {
    font-size: 30px;
}
p#pl2-predictions {
    font-size: 30px;
}

// Used only in 4-way races
p#pl3-name {
    font-size: 20px;
}
p#pl4-name {
    font-size: 20px;
}

// Used only in 4-way races
p#top-round {
    font-size: 22px;
}
p#bottom-round {
    font-size: 22px;
}

// Used only in 1v1 races
p#round {
    font-size: 26px;
}

p#game {
    font-size: 28px;
    // 30px, in 4-way races
}
p#goal {
    font-size: 26px;
    // 28px, in 4-way races
}
p#platform {
    font-size: 28px;
    // 28px as well, in 4-way races
}
p#subbed-by {
    font-size: 28px;
    // 22px, in 4-way races
}
```
