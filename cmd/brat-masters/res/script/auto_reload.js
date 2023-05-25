let auto_reload = function() {
	/** Time, in milliseconds when the page finished loading. */
	let boot_time = 0;

	/**
	 * Register an event to retrieve when the page finished loading.
	 */
	document.addEventListener('DOMContentLoaded', function (event) {
		boot_time = Date.now();
	});

	/**
	 * Event handler for "/mfh-handler/last-update" requests.
	 *
	 * Parse the response, a JSON with a single "Date" field, and reload
	 * the entire page if the parsed date is newer than the last boot date.
	 *
	 * @param{e} The event
	 */
	let _didUpdate = function(e) {
		try {
			let obj = JSON.parse(e);

			if (obj.Date > boot_time) {
				last_update = obj.Date;
				console.log('updated!');

				location.reload();
			}
		} catch (e) {
			console.log(e);
		}
	}

	/**
	 * Reload the current page if the server has been updated since this
	 * was loaded.
	 */
	let _update = function() {
		try {
			conn.getData('/brat-handler/last-update', _didUpdate);
		} catch (e) {
			console.log(e);
		}
	}

	return {
		'update': _update,
	};
}();
