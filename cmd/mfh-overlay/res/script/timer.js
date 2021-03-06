let timer = function() {
    /**
     * Try to retrieve the specified element and throws an error if it
     * doesn't exist.
     *
     * @param{id} The element's ID
     */
    let _getOrThrow = function(id) {
        let obj = document.getElementById(id);
        if (!obj) {
            throw('Couldn\'t find ' + id);
        }
        return obj;
    }

    /** Labels for each timer component */
    let timer_hour_hd = null;
    let timer_hour_ld = null;
    let timer_min_colon = null;
    let timer_min_hd = null;
    let timer_min_ld = null;
    let timer_s_colon = null;
    let timer_s_hd = null;
    let timer_s_ld = null;
    let timer_ms_hd = null;
    let timer_ms_ld = null;

    if ((typeof skipTimerUpdate !== 'boolean') || !skipTimerUpdate) {
        /**
         * Register an event to load the timer after the page has been
         * completely loaded.
         */
        document.addEventListener('DOMContentLoaded', function (event) {
            timer_hour_hd = _getOrThrow('hour-hd');
            timer_hour_ld = _getOrThrow('hour-ld');
            timer_min_colon = _getOrThrow('min-colon');
            timer_min_hd = _getOrThrow('min-hd');
            timer_min_ld = _getOrThrow('min-ld');
            timer_s_colon = _getOrThrow('s-colon');
            timer_s_hd = _getOrThrow('s-hd');
            timer_s_ld = _getOrThrow('s-ld');
            timer_ms_hd = _getOrThrow('ms-hd');
            timer_ms_ld = _getOrThrow('ms-ld');
        });
    }

    /**
     * Convert a time in milliseconds into an object with a hour, minute
     * seconds and milliseconds component.
     *
     * @param{_time} The time in milliseconds.
     * @return The object with elements 'hour', 'min', 'sec' and 'ms'.
     */
    let _msToHMSms = function(_time) {
        let ret = {
             'hour': 0,
             'min': 0,
             'sec': 0,
             'ms': 0,
        };

        ret.ms = _time % 1000;
        _time = Math.trunc(_time / 1000);

        ret.sec = _time % 60;
        _time = Math.trunc(_time / 60);

        ret.min = _time % 60;
        _time = Math.trunc(_time / 60);

        ret.hour = _time % 24;

        return ret;
    }

    /**
     * Event handler for "/timer" requests.
     *
     * Parse and update the received timer, in milliseconds. The message
     * shall be encoded as a JSON with a single "Time" field, an integer.
     *
     * @param{e} The event
     */
    let _updateTimer = function(e) {
        let res = JSON.parse(e);

        /* Ignored anything bellow centiseconds */
        let _time = _msToHMSms(res.Time);
        _time.ms = Math.trunc(_time.ms / 10);

        let ms_ld = _time.ms % 10;
        let ms_hd = Math.trunc(_time.ms / 10);

        let s_ld = _time.sec % 10;
        let s_hd = Math.trunc(_time.sec / 10);

        let min_ld = _time.min % 10;
        let min_hd = Math.trunc(_time.min / 10);

        let hour_ld = _time.hour % 10;
        let hour_hd = Math.trunc(_time.hour / 10);

        timer_ms_ld.innerText = '' + ms_ld;
        timer_ms_hd.innerText = '' + ms_hd;
        timer_s_ld.innerText = '' + s_ld;
        timer_s_hd.innerText = '' + s_hd;
        timer_min_ld.innerText = '' + min_ld;
        timer_min_hd.innerText = '' + min_hd;
        timer_hour_ld.innerText = '' + hour_ld;
        timer_hour_hd.innerText = '' + hour_hd;

        if (_time.min > 0 || _time.sec > 9 || _time.hour > 0) {
            timer_s_hd.style.visibility = 'visible';
        }
        else {
            timer_s_hd.style.visibility = 'hidden';
        }

        if (_time.min > 0 || _time.hour > 0) {
            timer_s_colon.style.visibility = 'visible';
            timer_min_ld.style.visibility = 'visible';
        }
        else {
            timer_s_colon.style.visibility = 'hidden';
            timer_min_ld.style.visibility = 'hidden';
        }

        if (_time.min > 9 || _time.hour > 0) {
            timer_min_hd.style.visibility = 'visible';
        }
        else {
            timer_min_hd.style.visibility = 'hidden';
        }

        if (_time.hour > 0) {
            timer_min_colon.style.visibility = 'visible';
            timer_hour_ld.style.visibility = 'visible';
        }
        else {
            timer_min_colon.style.visibility = 'hidden';
            timer_hour_ld.style.visibility = 'hidden';
        }

        if (_time.hour > 9) {
            timer_hour_hd.style.visibility = 'visible';
        }
        else {
            timer_hour_hd.style.visibility = 'hidden';
        }
    }

    /**
     * "Handles" anything that goes wrong.
     *
     * @param{e} The event
     */
    let _onError = function(e) {
        alert(e);
    }

    /**
     * Get the current time.
     *
     * On success, the callback shall receive an object with a 'Time'
     * element, representing the currently elapsed time in milliseconds.
     *
     * @param{onget} Callback executed on success.
     * @param{onerror} Callback executed on failure.
     */
    let _getTime = function(onget, onerror) {
        try {
            let _c = conn.newConn('/timer', 'GET', true);
            _c.addHeader("Content-Type", "application/json");
            _c.send(null, onget, onerror);
        } catch (e) {
            console.log(e);
        }
    }

    /**
     * Request the current timer to the server and update the timer label.
     */
    let _update = function() {
        _getTime(_updateTimer, _onError);
    }

    /** Empty callback. */
    let _nullCb = function(e) {
    }

    /**
     * Send a command to the timer.
     *
     * The possible actions are: "setup", "start", "stop", "reset", "add"
     * and "sub". From those, "setup", "add" and "sub" take an extra
     * parameter: a time in milliseconds.
     * 
     *
     * @param{action} The action to be sent.
     * @param{data} Parameter sent alongside the action.
     */
    let _sendCmd = function(action, data=null) {
        try {
            let _obj = {'Action': action};
            if (data) {
                _obj['Value'] = data;
            }
            let obj = JSON.stringify(_obj);
            conn.sendData('/timer', obj, onSave=_nullCb, method='POST');
        } catch (e) {
            console.log(e);
        }
    }

    /** Start the timer, if it's not running yet. */
    let _start = function() {
        _sendCmd('start');
    }

    /** Stop the timer, but keep its current value. */
    let _stop = function() {
        _sendCmd('stop');
    }

    /** Stop the timer and reset it back to 0. */
    let _reset = function() {
        _sendCmd('stop');
        _sendCmd('setup', 0);
        _sendCmd('reset');
    }

    /**
     * Retrieve the current timer.
     *
     * On success, 'onget' is called with an object containing the
     * current time in milliseconds (on element 'Time') and the time
     * components on elements 'hour', 'min', 'sec' and 'ms'.
     *
     * @param{onget} Callback executed on success.
     */
    let _get = function(onget) {
        let cb = function(e) {
            let res = JSON.parse(e);
            let _time = _msToHMSms(res.Time);
            _time.Time = res.Time;

            onget(_time);
        }
        _getTime(cb, null);
    }

    /**
     * Configure an initial offset for the timer.
     *
     * If the timer is already running, this offset is applied to the
     * currently accumulated time, but does not change it. So, you keep
     * initially set the offset as 5 min, but 2 minutes after starting
     * the timer change it to 4:30 min  and everything would work (i.e.,
     * the timer would then return 6:30 min, instead of 7 min).
     *
     * @param{startingTime} The new starting time for the timer.
     */
    let _set = function(startingTime) {
        _sendCmd('setup', startingTime);
    }

    return {
        'update': _update,
        'start': _start,
        'stop': _stop,
        'reset': _reset,
        'set': _set,
        'get': _get,
    };
}();
