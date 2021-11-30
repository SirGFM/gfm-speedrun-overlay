/**
 * Handle setting up the splits view, as well as configuring the times etc.
 */
{

const _minDiffWidth = 40;
const _minTimeWidth = 65;
const _minSplitWidth = _minDiffWidth + _minTimeWidth;
const _labelClass = 'm5x7 split splitLabel outlined';

/**
 * Group of elements inside 'splitsDiv', used to display the split. Each group
 * is represented by the object:
 *   - div: The div that contains every label
 *   - label: The split's label/title/name
 *   - diff: Difference between the split's time and the accumulated time
 *   - time: Target/Achieved time on the split
 *   - tgtTime: The actual (integer) target/achieved time on the split, in ms
 */
let _els = [];
/** Number of elements in use (since _els may be recycled) */
let _count = 0;
/** The current split being played */
let _cur = 0;
/** URL associated with the currently displayed splits */
let _url = "";
/** Whether the splits have already been upload */
let _hasSent = false;
/** Maximum height for the splits border */
let _maxHeight = 720;
/** Set the width of the split fields */
let _labelWidth;
let _diffWidth;
let _timeWidth;

let hideSplits = function() {
    if (_els.length == 0)
        return;
    hideBoxById('splitsDiv');
    for (let i in _els) {
        let obj = _els[i];

        obj.div.style.visibility = "hidden";
    }
    _cur = -1;
    _count = 0;
}

let loadedCallback = function(res) {
    try {
        let data = JSON.parse(res);

        if (!("entries" in data)) {
            hideSplits();
            throw "Invalid response '" + data + "'";
        }
        // TODO Setup "best" possible split?

        let splits = {};
        splits.entries = data.entries;
        initSplits(splits);

        /* Signalize window to center everything */
        let e = new Event("finishSetup");
        document.dispatchEvent(e);
    } catch (e) {
        alert(e);
        /* Throw again so it may be captured by the console */
        throw e;
    }
}

let loadNewSplits = function(url) {
    getData(url, loadedCallback);
}

let saveCurrentSplits = function() {
    if (!_url || _hasSent) {
        return;
    }

    let splits = {};
    let entries = [];

    for (let i = 0; i < _count; i++) {
        let _el = {
            label: _els[i].label.innerText,
            time: _els[i].tgtTime
        };
        entries.push(_el);
    }

    splits.entries = entries;
    let data = JSON.stringify(splits);
    sendData(_url, data);
    _hasSent = true;
}

function setupSplits(splits, maxWidth, maxHeight=720) {
    _maxHeight = maxHeight;
    if ("server" in splits) {
        _url = splits.server;
        loadNewSplits(_url);
    }
    else if ("entries" in splits) {
        initSplits(splits);
    }
    else {
        hideSplits();
    }
}

let initSplits = function(splits) {
    let div = document.getElementById('splitsDiv');
    if (!div) {
        div = document.createElement('div');
        div.id = 'splitsDiv';
        div.style.overflow = 'hidden';
        div.style.position = 'absolute';
        /** Create a box with dummy dimensions */
        createBox(div, 240, 360, anchor=true, darkBG=true, hasShadow=true, padContent=true);
        div.style.left = getBoxContentPosition(div) + 'px';
        div.style.top = getBoxContentPosition(div) + 'px';
    }

    if (!("entries" in splits)) {
        hideSplits();
        return;
    }

    /** Update the box dimension based on the number of splits */
    setupGetLabelLength(_labelClass);
    let _elHeight = getLabelHeight('TALL LETTERS I');

    let _h = _elHeight * splits.entries.length;
    while (getBoxDimension(_h, true) > _maxHeight)
        _h = getValidDimension(_h - _elHeight);
    let _w = 0;
    for (let _el in splits.entries) {
        let _elWidth = getLabelLength(splits.entries[_el]);
        if (_elWidth > _w)
            _w = _elWidth;
    }
    _w = getValidDimension(_w + _minSplitWidth * 1.2);
    while (getBoxDimension(_w, true) > _maxWidth)
        _w = getValidDimension(_w * 0.95);

    createBox(div, _w, _h, anchor=true, darkBG=true, hasShadow=true, padContent=true);
    showBox(div);
    div.style.maxHeight = _h + 'px';
    div.style.width = _w + 'px';
    div.style.height = 'auto';

    /* Calculate the width of each field */
    let _mult = 0.15;
    while (getValidDimension(_w * _mult) < _minTimeWidth)
        _mult += 0.05;
    _diffWidth = getValidDimension(_w * _mult);
    _timeWidth = getValidDimension(_w * _mult);
    _labelWidth = _w - _diffWidth - _timeWidth;

    _cur = 0;

    /* Make sure only objects that will be used are visible */
    _count = splits.entries.length;
    while (_els.length < _count) {
        let obj = {};

        obj.div = document.createElement("div");
        addLine('__auto__split_'+_els.length, '', _labelClass, maxBorderWidth=_diffWidth, isBoxed=false, fixedLen=true);
        obj.label = getLineLabel('__auto__split_'+_els.length);
        obj.diff = document.createElement("label");
        obj.diff.setAttribute("class", "m5x7 split splitDiff outlined");
        obj.diff.style.width = _diffWidth;
        obj.time = document.createElement("label");
        obj.time.setAttribute("class", "m5x7 split splitTime outlined");
        obj.time.style.width = _timeWidth;
        obj.tgtTime = 0;

        _els.push(obj);
        obj.div.append(obj.label);
        obj.div.append(obj.diff);
        obj.div.append(obj.time);
        div.append(obj.div);
    }

    for (let i = 0; i < splits.entries.length; i++) {
        let obj = _els[i];

        obj.div.style.visibility = "visible";
    }
    for (let i = splits.entries.length; i < _els.length; i++) {
        let obj = _els[i];

        obj.div.style.visibility = "hidden";
    }

    resetSplits(splits);
}

function resetSplits(splits) {
    let hideSplit = true;
    _cur = 0;

    /* Load the supplied entries into the div */
    for (let i in splits.entries) {
        let entry = splits.entries[i];
        let obj = _els[i];

        addLine('__auto__split_'+i, entry.label, _labelClass, maxBorderWidth=_diffWidth, isBoxed=false, fixedLen=true);
        if ("time" in entry) {
            obj.tgtTime = entry.time
        }
        else {
            obj.tgtTime = NaN;
        }

        updateSplit(i, obj.tgtTime, hideSplit);
        obj.div.removeAttribute("class", "highlightedBg");
    }

    /* Highlight the first split */
    _els[0].div.setAttribute("class", "highlightedBg");
    try {
        _els[0].label.scrollIntoViewIfNeeded();
    } catch (e) {
        _els[0].label.scrollIntoView();
    }
    _hasSent = false;
}

let setDiff = function(idx, time, force=false) {
    let obj = _els[idx];
    let ms = false, autoHide = true;
    let signal = "";
    let tmpTime = 0;

    if (time - obj.tgtTime >= 0) {
        tmpTime = time - obj.tgtTime;
        signal = "+";
    }
    else if (force || time - obj.tgtTime > -10000) {
        tmpTime = Math.abs(time - obj.tgtTime);
        signal = "-";
    }
    else {
        /* Split isn't too good nor has gone too bad*/
        return;
    }
    ms = (tmpTime < 60000);

    /* TODO Change color depending on diff */
    obj.diff.innerText = signal + timeToText(tmpTime, ms, autoHide);
}

let updateSplit = function(idx, time, hideSplit=false) {
    let obj = _els[idx];

    if (isNaN(time)) {
        obj.diff.innerText = "";
        obj.time.innerHTML = "NaN</br>";
        obj.tgtTime = NaN;
    }
    else {
        let autoHide = true, force = true;
        let ms = (time < 60000);

        if (hideSplit) {
            obj.diff.innerText = "";
        }
        else {
            setDiff(idx, time, force);
        }
        obj.time.innerHTML = timeToText(time, ms, autoHide) + "</br>";
        obj.tgtTime = time;
    }
}

function updateCurrentDiff(time) {
    if (!hasMoreSplits()) {
        return;
    }
    setDiff(_cur, time);
}

function setCurrentSplit(time) {
    if (!hasMoreSplits()) {
        return;
    }

    updateSplit(_cur, time);

    _els[_cur].div.removeAttribute("class", "highlightedBg");
    _cur++;
    if (_cur < _els.length) {
        _els[_cur].div.setAttribute("class", "highlightedBg");
    }

    let tgt = undefined;
    if (_cur + 1 < _els.length) {
        tgt = _els[_cur + 1]
    }
    else if (_cur < _els.length) {
        tgt = _els[_cur]
    }

    if (tgt) {
        try {
            tgt.label.scrollIntoViewIfNeeded();
        } catch (e) {
            tgt.label.scrollIntoView();
        }
    }
}

function hasMoreSplits() {
    let hasMore = (_cur >= 0 && _cur < _count);
    if (!hasMore && _cur > 0) {
        saveCurrentSplits();
    }
    return hasMore;
}

function reloadSplits() {
    if (!_url) {
        return;
    }

    loadNewSplits(_url);
    _hasSent = false;
}

function isSplitsVisible() {
    if (_els.length == 0)
        return false;
    return isBoxVisibleById('splitsDiv');
}

}
