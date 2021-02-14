let conn = function() {
    /**
     * Open a new connection that may be configured as needed.
     *
     * Configuration of the connection should be done by making successive
     * 'addHeader()' calls (which may be made in sequence, like
     * 'getter.addHeader(...).addHeader(...);').
     *
     * Then, 'send()' must be called to actually execute the request.
     *
     * @param{url} HTTP address of the connecting server.
     * @param{mode} HTTP Method for the connection (e.g., PUT or GET).
     * @param{async} HTTP Method for the connection (e.g., PUT or GET).
     */
    let _newConn = function(url, mode, async) {
        /* Open the connection. */
        let xhr = new XMLHttpRequest();
        xhr.open(mode, url, async);
        xhr.setRequestHeader("Access-Control-Allow-Origin", "*");

        let okCode = 200;
        if (mode != 'GET') {
            okCode = 204;
        }

        /**
         * Add a new header to the connection.
         *
         * @param{key} The header's key (e.g., 'Content-Type').
         * @param{value} The header's value (e.g., 'application/json').
         * @return The connection object itself.
         */
        let _addHeader = function(key, value) {
            xhr.setRequestHeader(key, value);
            return this;
        }

        /** Function used to ignore callbacks. */
        let _ignore = function(e) {}

        /**
         * Send the HTTP message.
         *
         * @param{obj} Data to be sent in the message.
         * @param{okCb} Callback if the operation is successfull.
         * @param{errCb} Callback if the operation fails.
         */
        let _send = function(obj, okCb, errCb) {
            if (okCb) {
                let cb = function (e) {
                    let htmlResponse = e.target;
                    if (htmlResponse.status == okCode) {
                        okCb(htmlResponse.response);
                    }
                    else {
                        errCb(e);
                    }
                };
                xhr.addEventListener("loadend", cb);
            }
            if (errCb) {
                xhr.addEventListener("error", errCb);
            }

            /* Ignore other events */
            xhr.addEventListener("loadstart", _ignore);
            xhr.addEventListener("load", _ignore);
            xhr.addEventListener("progress", _ignore);

            xhr.send(obj);
        }

        return {
            "addHeader": _addHeader,
            "send": _send,
        };
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
     * Send a request with the specified mode.
     *
     * @param{url} HTTP address of the connecting server.
     * @param{mode} HTTP Method for the connection (e.g., PUT or GET).
     * @param{obj} The object to be sent.
     * @param{okCb} Callback executed when the request is fully sent/received.
     *                onLoad must receive a single argument (the response).
     * @return A newly set up XMLHttpRequest.
     */
    let _sendReq = function(url, mode, obj, okCb) {
        let _c = _newConn(url, mode, true);

        if (obj) {
            _c.addHeader("Content-Type", "application/json");
        }
        _c.send(obj, okCb, _onError);
    }

    /**
     * Delete a resource from the server.
     *
     * @param{url} HTTP address of the resource server.
     * @param{onDel} Callback executed when the request is fully received.
     *               onDel must receive a single argument (the response).
     */
    function _delData(url, onDel) {
        _sendReq(url, "DELETE", null, onDel);
    }

    /**
     * Update a resource in the server.
     *
     * @param{url} HTTP address of the resource server.
     * @param{obj} The object to be sent.
     * @param{onPut} Callback executed when the request is fully received.
     *               onPut must receive a single argument (the response).
     */
    function _updateData(url, obj, onPut) {
        _sendReq(url, "PUT", obj, onPut);
    }

    /**
     * Retrieve a JSON object from the server.
     *
     * @param{url} HTTP address of the connecting server.
     * @param{onLoad} Callback executed when the request is fully received.
     *                onLoad must receive a single argument (the response).
     */
    function _getData(url, onLoad) {
        _sendReq(url, "GET", null, onLoad);
    }

    /**
     * Dummy callback for reporting that data was saved.
     */
    let _onSave = function (e) {
        alert("Object saved successfully");
    }

    /**
     * Save a JSON object on the server.
     *
     * @param{url} HTTP address of the connecting server.
     * @param{obj} The object to be sent.
     */
    function _sendData(url, obj, onSave=null, method='PUT') {
        if (!onSave) {
            onSave = _onSave;
        }
        _sendReq(url, method, obj, onSave);
    }

    return {
        "getData": _getData,
        "sendData": _sendData,
        "delData": _delData,
        "updateData": _updateData,
        "newConn": _newConn,
    };
}();
