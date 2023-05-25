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

	/**
	 * Try to retrieve the specified element within the parent element,
	 * throwing an error if it doesn't exist.
	 *
	 * @param{parent} The element's parent
	 * @param{id} The element's ID
	 */
	let _getChildOrThrow = function(parent, id) {
		let obj = parent.querySelector('#'+id);
		if (!obj) {
			throw('Couldn\'t find ' + id);
		}
		return obj;
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
	 * _updateTimerByValue handles updating the div by the specified time.
	 *
	 * @param{div_id} The div with the timer to be updated.
	 * @param{time_ms} The current time, in milliseconds.
	 */
	let _updateTimerByValue = function(div_id, time_ms) {
		/* Ignored anything bellow centiseconds */
		let _time = _msToHMSms(time_ms);

		let s_ld = _time.sec % 10;
		let s_hd = Math.trunc(_time.sec / 10);

		let min_ld = _time.min % 10;
		let min_hd = Math.trunc(_time.min / 10);

		parentEl = _getOrThrow(div_id);
		timer_s_ld = _getChildOrThrow(parentEl, 's-ld');
		shadow_timer_s_ld = _getChildOrThrow(parentEl, 'shadow-s-ld');
		timer_s_hd = _getChildOrThrow(parentEl, 's-hd');
		shadow_timer_s_hd = _getChildOrThrow(parentEl, 'shadow-s-hd');
		timer_min_ld = _getChildOrThrow(parentEl, 'min-ld');
		shadow_timer_min_ld = _getChildOrThrow(parentEl, 'shadow-min-ld');
		timer_min_hd = _getChildOrThrow(parentEl, 'min-hd');
		shadow_timer_min_hd = _getChildOrThrow(parentEl, 'shadow-min-hd');

		timer_s_ld.innerText = '' + s_ld;
		shadow_timer_s_ld.innerText = '' + s_ld;
		timer_s_hd.innerText = '' + s_hd;
		shadow_timer_s_hd.innerText = '' + s_hd;
		timer_min_ld.innerText = '' + min_ld;
		shadow_timer_min_ld.innerText = '' + min_ld;
		timer_min_hd.innerText = '' + min_hd;
		shadow_timer_min_hd.innerText = '' + min_hd;
	}

	/**
	 * Event handler for "/timer" requests.
	 *
	 * Parse and update the received timer, in milliseconds. The message
	 * shall be encoded as a JSON with a single "Time" field, an integer.
	 *
	 * @param{e} The event
	 */
	let _updateTimer = function(e, div_id) {
		let res = JSON.parse(e);

		_updateTimerByValue(div_id, res.Time);
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
	 * @param{url} The timer URL
	 */
	let _getTime = function(onget, onerror, url) {
		try {
			let _c = conn.newConn(url, 'GET', true);
			_c.addHeader("Content-Type", "application/json");
			_c.send(null, onget, onerror);
		} catch (e) {
			console.log(e);
		}
	}

	/**
	 * Request the current timer to the server and update the timer label.
	 */
	let _update = function(div_id) {
		_getTime(function(e) { _updateTimer(e, div_id); }, _onError, '/timer');
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
	 * @param{url} The timer URL
	 */
	let _get = function(onget, url='/timer') {
		let cb = function(e) {
			let res = JSON.parse(e);
			let _time = _msToHMSms(res.Time);
			_time.Time = res.Time;

			onget(_time);
		}
		_getTime(cb, null, url);
	}

	/**
	 * Configure the current value of the timer.
	 *
	 * @param{startingTime} The new time for the timer.
	 */
	let _set = function(startingTime) {
		_sendCmd('setup', startingTime);
		_sendCmd('reset');
	}

	return {
		'update': _update,
		'update_by_value': _updateTimerByValue,
		'start': _start,
		'stop': _stop,
		'reset': _reset,
		'set': _set,
		'get': _get,
	};
}();
