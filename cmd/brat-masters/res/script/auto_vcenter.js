let auto_vcenter = function() {
	/**
	 * Run through every object with class 'vcenter'
	 * and set it's line height to the container's height and text align to center,
	 * making texts within the container to be centered.
	 *
	 * This uses the same logic as auto_resize.js to ensure that the object's height is correct.
	 */
	let set_vcenter = function(event) {
		let resizables = document.getElementsByClassName('vcenter');

		/* Otherwise, simply center everything. */
		for (let i = 0; i < resizables.length; i++) {
			/* Use the computed style to get the element's font size. */
			let el = resizables[i];

			/* Try to use the element itself, but fallback to its parent
			 * if needed. */
			let heightEl = el;
			while (heightEl.scrollHeight == 0 && heightEl.parentElement != null) {
				heightEl = heightEl.parentElement;
			}

			el.style.textAlign = 'center';
			el.style.lineHeight = heightEl.offsetHeight + 'px';
		}
	}

	return {
		'set_vcenter': set_vcenter,
	};
}();
