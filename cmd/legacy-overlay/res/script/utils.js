let _whiteSpaceLen = 0;

/**
 * Setup the current document so getLabelLength may be later called.
 * It may be called more than once to calculate the length of a string for
 * different fonts.
 *
 * @param{classList} List of classes used by the font, as space-separated
 *                   strings (e.g. "class1 class2").
 */
function setupGetLabelLength(classList) {
    let label = document.getElementById("label-width-calculator");

    if (!label) {
        label = document.createElement("label");
        label.id = "label-width-calculator";

        // Set some inline elements to help calculate the length
        label.style.position = "absolute";
        label.style.visibility = "hidden";
        label.style.width = "auto";
        label.style.height = "auto";
        label.style.whiteSpace = "nowrap";

        document.body.insertAdjacentElement('beforeend', label);
    }
    // Overrides any previously set classes
    label.className = classList;

    // Forces a redraw (most likely useless, but still fun :D)
    void label.offsetWidth;

    // This is ugly, but gets rid of a corner case (not being able to retrieve
    // the label length on the first call)
    getLabelLength("fix coner case");

    // Store the length of a single whitespace
    _whiteSpaceLen = getLabelLength("A A") - 2 * getLabelLength("A");
}

let _getLabelDimension = function(txt, addTrailingSpace, dimension) {
    let val = 0;
    let label = document.getElementById("label-width-calculator");

    if (!label) {
        return 0;
    }

    label.innerText = txt;
    if (addTrailingSpace) {
        val = _whiteSpaceLen;
    }
    val += label[dimension];

    label.innerText = "";

    return val;
}

/**
 * Compute the length of a given string.
 *
 * @param{txt} The string
 * @param{addTrailingSpace} Whether a single whitespace should be added at the end.
 */
function getLabelLength(txt, addTrailingSpace) {
    return _getLabelDimension(txt, addTrailingSpace, "clientWidth");
}

function getLabelHeight(txt) {
    return _getLabelDimension(txt, false, "clientHeight");
}

/**
 * Check if a given value is inside an array.
 *
 * @param{arr} The array
 * @param{val} The value
 */
function isInArray(arr, val) {
    for (var i = 0; i < arr.length; i++) {
        if (arr[i] == val) {
            return true;
        }
    }
    return false;
}

/**
 * Parse all query strings from a URL. Repeated names are
 * overwritten (instead of added to a list).
 *
 * @param{url} The url to be parsed.
 */
function parseURL(url) {
    let cmds = undefined;
    let ret = {};
    // Guaranteed to be the first hit (instead of inside some sub-string).
    let from = url.indexOf('?');

    if (from == -1) {
        return ret;
    }

    // Retrieve the argument list, removing any trailing '/'
    if (url[url.length-1] == '/') {
        cmds = url.substring(from+1, url.length-1);
    }
    else {
        cmds = url.substring(from+1, url.length);
    }

    // Parsing is quite simple: arguments are separated by '&' and name/value by '='.
    cmdList = cmds.split('&')
    for (var i = 0; i < cmdList.length; i++) {
        var name = undefined;
        var val = undefined;
        var cmd = cmdList[i];
        var split = cmd.indexOf('=');

        name = cmd.substring(0, split);
        val = cmd.substring(split+1);

        ret[name] = val;
    }

    return ret;
}

/**
 * If necessary, increases the value by 1 so it becomes even.
 *
 * @param{val} Value to be modified.
 */
function makeValueEven(val) {
    return val + (val % 2);
}

/**
 * Centers an object horizontally.
 *
 * @param{elementName} Name of the element to be centralized
 * @param{parentName} Element that should be used for centering and for the
 *                    horizontal offset. The element must no be parented to this!
 *                    If empty, the entire document is used.
 * @param{newWidth} New width of the object. If not set, defaults to offsetWidth
 * @param{forceEven} Forces every value to be even values
 * @param{maxViewWidth} Cap the view's width to this value (useful for
 *                      subsections)
 */
function centerDomHorizontal(elementName, parentName, newWidth, forceEven, maxViewWidth) {
    let left = 0;
    let viewLeft = 0;
    let viewWidth = 0;
    let element = document.getElementById(elementName);

    if (parentName) {
        let parentEl = document.getElementById(parentName);

        viewLeft = parentEl.offsetLeft;
        viewWidth = parentEl.offsetWidth;
    }
    else {
        viewLeft = document.body.offsetLeft;
        viewWidth = document.body.offsetWidth;
    }

    if (!newWidth) {
        newWidth = element.offsetWidth;
    }
    if (maxViewWidth && viewWidth > maxViewWidth) {
        viewWidth = maxViewWidth;
    }

    left = viewLeft + Math.floor((viewWidth - newWidth) * 0.5);

    if (forceEven) {
        left += left % 2;
        newWidth += newWidth % 2;
    }

    element.style.left = left+"px";
    element.style.width = newWidth+"px";
}

/**
 * Centers an object vertically.
 *
 * @param{elementName} Name of the element to be centralized
 * @param{parentName} Element that should be used for centering and for the
 *                    vertical offset. The element must no be parented to this!
 *                    If empty, the entire document is used.
 * @param{newHeight} New height of the object. If not set, defaults to offsetHeight
 * @param{forceEven} Forces every value to be even values.
 * @param{maxViewHeight} Cap the view's height to this value (useful for
 *                      subsections)
 */
function centerDomVertical(elementName, parentName, newHeight, forceEven, maxViewHeight) {
    let elTop = 0;
    let viewTop = 0;
    let viewHeight = 0;
    let element = document.getElementById(elementName);

    if (parentName) {
        let parentEl = document.getElementById(parentName);

        viewTop = parentEl.offsetTop;
        viewHeight = parentEl.offsetHeight;
    }
    else {
        viewTop = document.body.offsetTop;
        viewHeight = document.body.offsetHeight;
    }

    if (!newHeight) {
        newHeight = element.offsetHeight;
    }
    if (maxViewHeight && viewHeight > maxViewHeight) {
        viewHeight = maxViewHeight;
    }

    elTop = viewTop + Math.floor((viewHeight - newHeight) * 0.5);

    if (forceEven) {
        elTop += elTop % 2;
        newHeight += newHeight % 2;
    }

    element.style.top = elTop+"px";
    element.style.height = newHeight+"px";
}

/**
 * Position an element (presumably a line) within the parent. The object's
 * relative offset (as it's expected to have 'absolute' positioning) is
 * set from 'pos' (either top or bottom).
 *
 * @param{elementName} Name of the element to be centralized
 * @param{pos} Relative position (either 'top' or 'bottom')
 * @param{offset} Offset from the relative position
 * @param{widthDiff} How much should the width differ from its parent's.
 */
function positionDivLine(elementName, pos, offset, widthDiff) {
    let element = document.getElementById(elementName);
    let newWidth = element.parentElement.offsetWidth;

    element.style[pos] = offset+"px";
    centerDomHorizontal(elementName, element.parentElement.id, newWidth+widthDiff, true);
    element.style.left = Math.floor(widthDiff * -0.5)+"px";
}

/**
 * Position a fixed element just bellow another
 *
 * @param{elementName} Name of the element to be positioned
 * @param{parentName} Element that should be used for positioning. The element
 *                    must no be parented to this! If empty, the entire document
 *                    is used.
 * @param{offset} Distance from the parent object.
 * @param{forceEven} Forces every value to be even values.
 */
function setDomPositionBellow(elementName, parentName, offset, forceEven) {
    let elTop = 0;
    let element = document.getElementById(elementName);

    if (parentName) {
        let parentEl = document.getElementById(parentName);

        elTop = parentEl.offsetTop;
        elTop += parentEl.offsetHeight;
    }
    else {
        elTop = document.body.offsetTop;
        elTop += document.body.offsetHeight;
    }
    elTop += offset;

    if (forceEven) {
        elTop += elTop % 2;
    }

    element.style.top = elTop+"px";
}

/**
 * Position an element (presumably a line) within the parent. The object's
 * relative offset (as it's expected to have 'absolute' positioning) is
 * set from 'pos' (either top or bottom).
 *
 * @param{elementName} Name of the element to be centralized
 * @param{pos} Relative position (either 'top' or 'bottom')
 * @param{offset} Offset from the relative position
 * @param{widthDiff} How much should the width differ from its parent's.
 */
function positionDivLine(elementName, pos, offset, widthDiff) {
    let element = document.getElementById(elementName);
    let newWidth = element.parentElement.offsetWidth;

    element.style[pos] = offset+"px";
    centerDomHorizontal(elementName, element.parentElement.id, newWidth+widthDiff, true);
    element.style.left = Math.floor(widthDiff * -0.5)+"px";
}

/**
 * Retrieve an integer attribute from an element. If it's 0, the first child's
 * attribute is retrieved instead.
 *
 * @{element} The element
 * @{attrName} The attribute's name
 */
function getIntegerAttr(element, attrName) {
    if (element[attrName] == 0) {
        return element.firstElementChild[attrName];
    }
    return element[attrName];
}
