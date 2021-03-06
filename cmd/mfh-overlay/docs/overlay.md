# MFH's Restream Overlay

Documentation for editing and/or simply understanding the organization of the overlay.

Interesting links:

* `http://localhost:8088/tmpl/overlay.go.html` - overlay's entry point.
* `http://localhost:8088/res/config.html` - page used to configure the current racers and game.
* `http://localhost:8088/res/dashboard.html` - page to control the ongoing race.

# 1v1 and 4-way races

The overlay is intended to be used in restreams of 1v1 and 4-way races. Although it can seamlessly change between the 2 layouts, by marking a option in the configuration page, there are also two extra mappings for the overlay:

* `http://localhost:8088/tmpl/1v1-overlay.go.html` - overlay that forcefully uses the 1v1 layout
* `http://localhost:8088/tmpl/2v2-overlay.go.html` - overlay that forcefully uses the 4-way layout

These work by simply overriding the 4 way flag (`Layout2v2`) to true or false on the server. Every other configuration is sent as-is.

# CSS/Styling

The overlay uses 3 css files:

* `res/style/common.css`: Common definitions for both layouts
    * Page size
    * Font configuration
* `res/style/1v1.css`: Specialization for 1v1 races
* `res/style/2v2.css`: Specialization for 4-way races

# Overlay dimensions and organization

The overlay is divided into a few areas:

* Top racers information
* Bottom racers information
* Game information

There also some pre-defined space for images, logos and some other stuff, which sort of exists apart of the rest of the layout.

## Racers information

The top racers information area is used both in 1v1 races as in 4-way races, while the bottom racers information area is used exclusively during 4-way racers.

Each racer has its own `div`, which may contain (as configured):

* An audio indicator
* The racer's name
* Predictions

Both the audio indicator and the predictions are optional, but these three fields share the same `div`. Since both the audio indicator and the predictions have static, pre-defined dimensions, adding or removing these fields will change the available space for the racer's name.

In the 1v1 overlay, the racers information takes up the entire top of the overlay, while in the 4-way overlay there's an area at the top of the overlay and another at the bottom of the overlay for the racers information.

Also, if needed, the racers' name will automatically shrink to fit the available space.

## Game information

This are contains information regarding the game being raced, such as:

* Game's title
* Race's goal
* Game's platform
* "Race description"

The "race description" can be used for anything. It has an accompanying label, which defaults to a smaller font size than the description. It may be used to indicate a game's submitter, the reject reason, that this is a weekly blind race... Anything, really!

In 4-way races, this area takes up the left portion of the overlay, and is displayed in the given order. In 1v1 races, the title and the goal are in a `div` by themselves, to the left, and the platform and the description are in another `div`, to the right.

If any of these information is not given, then it won't take up any space in the overlay (i.e., it's `p` is removed).

## Dimensions

### 1v1 race

The dimensions of the layout in the 1v1 layout are:

* Top (racers) area's height: 100px (from `div.player`and `div.center-label`)
    * Note that for the player name, this height is divided between the player name and optional pronouns
* Top (racers) area's player width: 580px (from `div.player-flex`)
    * Prediction's width: 66px (from `p.predictions-flex`)
    * Audio note's width: 48px (from `p.audio-flex`)
    * Logo area between both players: 120px
* Bottom (info) area's height: 155px (from `div.info` and `div.info-box`)
    * Race info (left side)'s width: 350px (from `div.left-info`)
        * Maximum game title height: 68px (roughly 3 lines, from `p#game`)
        * Maximum goal height: 48px (roughly 2 lines, from `p#goal`)
    * Image area between both info boxes: 114px (from `div.info-padding`)
    * Right side's width: 300px (from `div.right-info`)
        * **NOTE**: The right side has a right margin of 50px, to make it the same width as the left side
        * Maximum platform height: 34px (roughly 1 lines, from `p#platform`)
        * Maximum description height: 80px (roughly 2 lines, from `p#subbed-by`)

The round, if configured, takes the bottom 35px and is centered on the entire screen.

### 4-way race

* Racer area's height (top & bottom): 26px (from `div.player`and `div.center-label`)
* Racer area's width (top & bottom): 896px (from `div.player`)
* Eace racer area's width (top & bottom): 380px (from `div.player-flex`)
    * Audio note's width: 24px (from `p.audio-flex`)
    * Round area between both players: 136px (from `p.round`)
* Info area's width: 384px (~33% of the screen, from `div#info`)
* Info area's height: 605px (from `div#info`)
    * The top 115px of the info area are reserved for images
        * If any of the images are hidden, the other will become centralized
    * Maximum game title height: 108px (roughly 3 lines, from `p#game`)
    * Maximum goal height: 66px (roughly 2 lines, from `p#goal`)
    * Maximum platform height: 34px (roughly 1 lines, from `p#platform`)
    * Maximum description height: 56px (roughly 2 lines, from `p#subbed-by`)
    * **NOTE**: From all these, the bottom 379px should be free for anything

# Timer

The timer stays on the bottom right corner on the 1v1 layout and on the bottom left corner on the 4-way layout.

For both layouts, the timer has the following properties:

* Each digit is an individual element with fixed width (from the various `span.timer-*` and `span.ms-ld`)
    * This causes the timer to stand still, regardless of the length of any given digit
* Timer's width: 192px (from `div.timer`)
* Timer's height: 42px (from `div.timer`)

# Fonts

This overlay uses two fonts:

* [Regen](https://www.dafont.com/regen.font)
* [DS-Digital](https://www.dafont.com/ds-digital.font)

Simply download these fonts and place them in `./res/font/`. Currently, the overlay expect the following modes to be available

* Regen default: `/res/font/Regen.otf`
* Regen bold: `/res/font/Regen Bold.otf`
* Regen italic: `/res/font/Regen Italic.otf`
* Regen bold italic: `/res/font/Regen Bold Italic.otf`
* DS-Digital default: `/res/font/DS-DIGIT.TTF`
    * (although this is most likely the "bold italic" version...)
* DS-Digital bold: `/res/font/DS-DIGIB.TTF`
* DS-Digital italic: `/res/font/DS-DIGII.TTF`

