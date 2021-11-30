// When the button was actually pressed
let _pressTime = 0;

// Delayed handle call (if any)
let _pressHandle = null;

// Delay between timer updates
const _updateDelay = 1000 / 30;
// Timer that displays the current time
let _timerLabel = null;
// Accumulated time in ms
let _accumulatedTime = 0;
// Accumulated time when the timer was last started
let _prevAccumulatedTime = 0;
// Last time at which the timer was started
let _lastStartTime = 0;
// Worker that updates the timer, if any
let _timerRunner = null;

/**
 * Setup the label that shall display the current timer
 */
function setupTimer(element) {
    _timerLabel = element;
}

/**
 * Return time with fixed drift from when the button was first pressed.
 */
let getFixedAcc = function() {
    return _accumulatedTime - (timer.now() - _pressTime);
}

/**
 * Pauses/Unpauses the timer
 *
 * @param{time} Currently accumulated time
 * @param{forceHalt} Forces the timer to pause
 */
function toggleTimer(time, forceHalt=false) {
    if (forceHalt || _timerRunner != null) {
        window.clearInterval(_timerRunner);
        _timerRunner = null;
        _accumulatedTime = time;
        setTimerText();
    }
    else if (_timerRunner == null) {
        _timerRunner = window.setInterval(updateTimer, _updateDelay);
        _prevAccumulatedTime = _accumulatedTime;
        _lastStartTime = _pressTime;
    }
}

/**
 * Update the running timer
 */
function updateTimer() {
    _accumulatedTime = _prevAccumulatedTime + (timer.now() - _lastStartTime);
    setTimerText();
    /* Update the split, if set */
    updateCurrentDiff(_accumulatedTime);
}

/**
 * Updates the timer label with the currently accumulated timer.
 */
function setTimerText() {
    if (_timerLabel) {
        _timerLabel.innerText = timeToText(_accumulatedTime);
    }
}

/**
 * Converts a given time to text.
 *
 * @param{time} The integer time to be converted to string.
 * @param{showMs} Whether the text should contain milliseconds.
 * @param{autoHideHour} Whether units should be hidden, unless greater than 0.
 */
function timeToText(time, showMs=true, autoHide=false) {
    let ms = time % 1000;
    let s = Math.floor((time / 1000) % 60);
    let min = Math.floor((time / 60000) % 60);
    let hour = Math.floor((time / 3600000) % 24);

    let txt = "";

    if (!autoHide || hour > 0) {
        txt += ("" + hour).padStart(2, "0") + ":";
    }
    if (!autoHide || hour > 0 || min > 0) {
        txt += ("" + min).padStart(2, "0") + ":";
    }
    txt += ("" + s).padStart(2, "0");
    if (showMs) {
        txt += "." + ("" + ms).padStart(3, "0");
    }

    return txt;
}

let timer = function() {
    // Maximum time allowed between multi-presses
    const _multiPress = 300;
    // Amount of presses required for a multi-press
    const _resetTimerCount = 3;
    // How long a button must be held to be considered a press
    const _pressDetectionMs = 400;

    // Time when the last fresh callback was called
    let _lastFresh = 0;
    // Amount of presses since the last single press
    let _multiPressCount = 0;
    // Whether the press has been handled
    let _handled = false;

    /**
     * Resets the timer on multipresses and (un)pause it on single presses.
     */
    let handleTimerCallback = function () {
        switch (_multiPressCount) {
        case _resetTimerCount:
            // Reset the time and pauses it
            toggleTimer(0, true);
            _prevAccumulatedTime = 0;
            setTimerText();
            reloadSplits();
            break;
        default:
            let latest = getFixedAcc();
            if (_timerRunner != null && hasMoreSplits()) {
                setCurrentSplit(latest);
            }
            /* Stop if the previous split was the last */
            if (_timerRunner == null || !hasMoreSplits()) {
                toggleTimer(latest);
            }
            break;
        }
    }

    /**
     * Callback used to trigger reseting, pausing and unpausing the timer.
     *
     * To solve issues caused by miss-pressing the button, it ignores presses
     * shorter _pressDetectionMs.
     *
     * @param{state} State of the pressed button/key.
     * @param{fresh} Whether the state just transitioned.
     * @param{ts} Timestamp of the callback, in milliseconds.
     */
    let timerCallback = function(state, fresh, ts) {
        if (state && fresh) {
            // Detect whether it's a multipress or not
            if (ts - _lastFresh > _multiPress) {
                _multiPressCount = 1;
                _pressTime = ts;
            }
            else
                _multiPressCount++;

            _lastFresh = ts;
        }
        else if (!state)
            _handled = false;
        // If the button has been held for long enough, do the callback
        else if (!_handled && ts - _lastFresh >= _pressDetectionMs) {
            handleTimerCallback();
            _handled = true;
        }
    }

    document.addEventListener('timer-onpress', function (e) {
        let state = true;
        let fresh = true;
        timerCallback(state, fresh, Math.trunc(e.timeStamp));
    });
    document.addEventListener('timer-pressed', function (e) {
        let state = true;
        let fresh = false;
        timerCallback(state, fresh, Math.trunc(e.timeStamp));
    });
    document.addEventListener('timer-onrelease', function (e) {
        let state = false;
        let fresh = true;
        timerCallback(state, fresh, Math.trunc(e.timeStamp));
    });

    return {
        /**
         * @return: The current time, since the page opened, in milliseconds.
         */
        now: function() {
            return Math.trunc((new Event('dummy')).timeStamp);
        }
    };
}();
