/**
 * Credits to Alpha_5 for the original implementation, using jQuery.
 */
let auto_resize = function() {
	/** How many times the callback was called without resizing anything. */
	let _noRetryCount = 0;

	/**
	 * Run through every object with class 'resize' and shrink its font
	 * until it fits its area. However, the object must have a 'maxHeight'
	 * properly set!
	 *
	 * Additionally, if auto_vcenter.js was loaded,
	 * center everything vertically after this is done.
	 *
	 * For some reason, some browsers (well... Chrome(ium), really)
	 * eventually stop calculating the elements' 'scrollHeight' and
	 * 'offsetHeight'. So, this function is scheduled to be re-executed
	 * until it has ran once as many times as there are resizeable element
	 * and without resizing anything.
	 *
	 * Yes, this sounds dumb. No, I don't care anymore.
	 */
	let do_resize = function(event) {
		let resizables = document.getElementsByClassName('resize');

		_noRetryCount++;
		for (let i = 0; i < resizables.length; i++) {
			/* Use the computed style to get the element's font size. */
			let el = resizables[i];
			let css = window.getComputedStyle(el);
			/* Extract the size as a number and the extension */
			let old = parseInt(css.fontSize);
			let ext = css.fontSize.substr((old+'').length);

			/* Try to use the element itself, but fallback to its parent
			 * if needed. */
			let heightEl = el;
			while (heightEl.scrollHeight == 0 && heightEl.parentElement != null) {
				heightEl = heightEl.parentElement;
			}

			/* Set the element's line height temporarily,
			 * so this may be taken into account while re-sizing it,
			 * Storing the element's original line height and alignment.
			 */
			origTextAlign = el.style.textAlign;
			origLineHeight = el.style.lineHeight;
			el.style.textAlign = 'center';
			el.style.lineHeight = heightEl.offsetHeight + ext;

			/* Simply decrease the font slightly until it fits */
			while (heightEl.scrollHeight > heightEl.offsetHeight) {
				let _new = old - 1;

				el.style.fontSize = _new + ext;
				old = _new;

				/* Reset the counter as at least one element was resized */
				_noRetryCount = 0;
			}

			/* Recover the element's line height and alignment. */
			el.style.textAlign = origTextAlign;
			el.style.lineHeight = origLineHeight;
		}

		if (_noRetryCount <= resizables.length) {
			requestAnimationFrame(do_resize);
		}
		else if (typeof(auto_vcenter) !== 'undefined') {
			auto_vcenter.set_vcenter();
		}
	}

	/**
	 * Register an event to resize every element with class 'resize'.
	 */
	document.addEventListener('DOMContentLoaded', do_resize);

	/* This object runs on boot and that's it! */
	return {};
}();
