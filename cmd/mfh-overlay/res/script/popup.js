let popup = function() {
    /** Objects currently being shown */
    let _popups = {};

    /**
     * Hides the given component. Should be registered with setTimeout.
     *
     * @param{id} The object's identifier
     * @param{obj} The object/HTML element that will be hidden
     */
    let _hideElement = function(id, obj) {
        if (!obj.classList.contains('hidden')) {
            obj.classList.add('hidden');
        }
        delete _popups[id];
    }

    /**
     * Event handler for "/mfh-handler/popup" requests.
     *
     * Parse the response, with a single field "Elements". Elements should
     * be a list of objects, each with a string Id and the Timeout to hide
     * the element, in milliseconds.
     *
     * @param{e} The event
     */
    let _onPopup = function(e) {
        let res = JSON.parse(e);

        for (var i = 0; res && res.Elements && i < res.Elements.length; i++) {
            let el = res.Elements[i];
            let id = el.Id;
            let obj = document.getElementById(id);
            if (!obj) {
                console.log('Couldn\'t find element ${id}!');
                continue;
            }

            if (obj.classList.contains('hidden')) {
                obj.classList.remove('hidden');
            }

            /* Clear any active timeout, so the animation stay active for
             * longer. */
            if (id in _popups) {
                clearTimeout(_popups[id]);
            }
            _popups[id] = setTimeout(_hideElement, el.Timeout, id, obj);
        }
    }

    /**
     * Check if the server has issued any popup to be displayed.
     */
    let _update = function() {
        try {
            conn.getData('/mfh-handler/popup', _onPopup);
        } catch (e) {
            console.log(e);
        }
    }

    /**
     * Temporarily show the given element.
     *
     * @param{id} The element's ID
     * @param{timeout} How long until it vanishes
     */
    let _show = function(id, timeout=1000) {
        try {
            let obj = JSON.stringify({'Id': id, 'Timeout': timeout});
            conn.sendData('/mfh-handler/popup', obj, onSave=null, method='POST');
        } catch (e) {
            console.log(e);
        }
    }

    return {
        'update': _update,
        'show': _show,
    };
}();
