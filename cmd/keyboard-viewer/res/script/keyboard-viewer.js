let keyboard_viewer = function() {
	/** URL from where the keyboard state is retrieved. */
	let _url = '';
	/** Dictionary of keys to DOM. */
	let _keys = null;
	/** Map of key index in the service to key name. */
	let _map = null;
	/** Interval's ID. */
	let pollThID = null;

	/**
	 * Send a GET HTTP message.
	 *
     * @param{url} HTTP address of the resource.
	 * @param{okCb} Callback if the operation is successfull.
	 * @param{errCb} Callback if the operation fails.
	 */
	let _get = function(url, okCb, errCb) {
        /* Open the connection. */
        let xhr = new XMLHttpRequest();
        xhr.open('GET', url, true);

		if (okCb) {
			let cb = function (e) {
				let htmlResponse = e.target;
				if (htmlResponse.status == 200) {
					okCb(htmlResponse.response);
				}
				else if (errCb) {
					errCb(e);
				}
			};
			xhr.addEventListener('loadend', cb);
		}
		if (errCb) {
			xhr.addEventListener('error', errCb);
		}

        /** Function used to ignore callbacks. */
        let _ignore = function(e) {}

		/* Ignore other events */
		xhr.addEventListener('loadstart', _ignore);
		xhr.addEventListener('load', _ignore);
		xhr.addEventListener('progress', _ignore);

		xhr.send(null);
	}

	/**
	 * Update the keyboard viewer based on the received messages.
	 */
	let _updateKeyboard = function(e) {
		let state = JSON.parse(e);

		for (idx in state) {
			try {
				let key = _map[idx];

				if (state[idx] == 1) {
					_keys[key].img.style.visibility = 'visible';
				}
				else {
					_keys[key].img.style.visibility = 'hidden';
				}
			} catch (e) {}
		}
	}

	/**
	 * Callback called on an interval to update the keyboard.
	 */
	let _pollKeyboard = function() {
		_get(_url, _updateKeyboard, null);
	};

	/**
	 * Start keyboard viewer.
	 *
	 * On start, get the map of key index to key name. The keys are
	 * updated on a timer, based on the requested FPS.
	 *
	 * The keys dictionary must map the key name into the key's DOM, which
	 * must be on the object's 'img' attribute. E.g.:
	 *
	 *     {
	 *         'Backspace': {
	 *             // ...
	 *             'img': // The DOM for this key
	 *         },
	 *         // ...
	 *         'Period': {
	 *             // ...
	 *             'img': // The DOM for this key
	 *         }
	 *     }
	 *
     * @param{baseUrl} Base address of 'ram_store' server.
     * @param{fps} How many times per second should the keyboard be updated.
     * @param{keys} Dictionary of keys.
	 */
	let _start_interval = function(baseUrl, fps, keys) {
		_url = baseUrl + '/data';
		_keys = keys;

		if (_map == null) {
			let cb = function (e) {
                _map = JSON.parse(e);
				console.log(_map);
			}

			_get(baseUrl + '/map', cb, null);
		}

		if (pollThID == null) {
			pollThID = setInterval(_pollKeyboard, 1000 / fps);
		}
	}

	/**
	 * Stop the keyboard viewer.
	 */
	let _stop_interval = function() {
		if (pollThID != null) {
			clearInterval(pollThID);
		}
		pollThID = null;
	}

	return {
		'get': _get,
		'start_watching': _start_interval,
		'stop_watching': _stop_interval,
	};
}();
