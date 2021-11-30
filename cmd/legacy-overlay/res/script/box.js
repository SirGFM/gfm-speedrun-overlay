const _boxPrefix = '__auto__';
const _boxSuffix = '__box__';
const _contentPadding = 4;
const _borderSize = 6;
const _lineSize = 2;
const _shadowSize = 2;
const _lightOutline = '#ffffff';
const _darkOutline  = '#cbdbfc';
const _innerShadow  = '#3f3f74';
const _border       = '#5b6ee1';
const _background   = '#3f3f74';
const _outterShadow = '#222034';

let _boxCache = {};

let _getBoxIdById = function(id) {
    return _boxPrefix + id + _boxSuffix;
}

let _getBoxId = function(content) {
    return _getBoxIdById(content.id);
}

let _getChildId = function(content, childId) {
    return _getBoxId(content) + childId;
}

let _updateChild = function(content, childId, x, y, width, height) {
    let _child = document.getElementById(_getChildId(content, childId));

    _child.style.left = x + 'px';
    _child.style.top = y + 'px';
    _child.style.width = width + 'px';
    _child.style.height = height + 'px';

    return _child;
}

function getBoxDimension(innerDimension, padContent) {
    let _ret = innerDimension + 2 *_borderSize + 4 * _lineSize;
    if (padContent)
        return _ret + 2 * _contentPadding;
    return _ret;
}

/* Changes a box dimensions */
function setBoxDimensions(content, innerWidth, innerHeight) {
    let _padContent = _boxCache[_getBoxId(content)].padded;
    let _boxWidth = getBoxDimension(innerWidth, _padContent);
    let _boxHeight = getBoxDimension(innerHeight, _padContent);
    let _contentWidth = innerWidth;
    let _contentHeight = innerHeight;
    let _contentPos = 2 * _lineSize + _borderSize;

    if (_padContent) {
        _contentWidth += 2 * _contentPadding;
        _contentHeight += 2 * _contentPadding;
    }

    _updateChild(content,
                 'leftOutterShadow',
                 -2 * _shadowSize,
                 _shadowSize,
                 2 * _shadowSize,
                 _boxHeight);
    _updateChild(content,
                 'bottomOutterShadow',
                 -2 * _shadowSize,
                 _boxHeight,
                 _boxWidth,
                 2 * _shadowSize);
    _updateChild(content,
                 'background',
                 _contentPos,
                 _contentPos,
                 _contentWidth,
                 _contentHeight);
    _updateChild(content,
                 'leftBorder',
                 _lineSize,
                 _lineSize,
                 _borderSize,
                 _boxHeight - 2 * _lineSize);
    _updateChild(content,
                 'rightBorder',
                 _boxWidth - _lineSize - _borderSize,
                 _lineSize,
                 _borderSize,
                 _boxHeight - 2 * _lineSize);
    _updateChild(content,
                 'topBorder',
                 _lineSize,
                 _lineSize,
                 _boxWidth - 2 * _lineSize,
                 _borderSize);
    _updateChild(content,
                 'bottomBorder',
                 _lineSize,
                 _boxHeight - _lineSize - _borderSize,
                 _boxWidth - 2 * _lineSize,
                 _borderSize);
    _updateChild(content,
                 'leftInnerShadow',
                 _lineSize + _borderSize - _shadowSize,
                 2 * _lineSize + _borderSize,
                 _shadowSize,
                 _contentHeight + 2 * _lineSize);
    _updateChild(content,
                 'rightInnerShadow',
                 _boxWidth - _lineSize - _shadowSize,
                 _lineSize,
                 _shadowSize,
                 _boxHeight - 2 * _lineSize);
    _updateChild(content,
                 'topInnerShadow',
                 _lineSize,
                 _lineSize,
                 _boxWidth - 2 * _lineSize,
                 _shadowSize);
    _updateChild(content,
                 'bottomInnerShadow',
                 _lineSize + _borderSize - _shadowSize,
                 _boxHeight - _lineSize - _borderSize,
                 _contentWidth + 2 * _lineSize,
                 _shadowSize);
    _updateChild(content,
                 'leftOutterOutline',
                 0,
                 0,
                 _lineSize,
                 _boxHeight);
    _updateChild(content,
                 'rightOutterOutline',
                 _boxWidth - _lineSize,
                 0,
                 _lineSize,
                 _boxHeight);
    _updateChild(content,
                 'topOutterOutline',
                 0,
                 0,
                 _boxWidth,
                 _lineSize);
    _updateChild(content,
                 'bottomOutterOutline',
                 0,
                 _boxHeight - _lineSize,
                 _boxWidth,
                 _lineSize);
    _updateChild(content,
                 'leftInnerOutline',
                 _contentPos - _lineSize,
                 _contentPos - _lineSize,
                 _lineSize,
                 _contentHeight + 2 * _lineSize);
    _updateChild(content,
                 'rightInnerOutline',
                 _contentPos + _contentWidth,
                 _contentPos - _lineSize,
                 _lineSize,
                 _contentHeight + 2 * _lineSize);
    _updateChild(content,
                 'topInnerOutline',
                 _contentPos - _lineSize,
                 _contentPos - _lineSize,
                 _contentWidth + 2 * _lineSize,
                 _lineSize);
    _updateChild(content,
                 'bottomInnerOutline',
                 _contentPos - _lineSize,
                 _contentPos + _contentHeight,
                 _contentWidth + 2 * _lineSize,
                 _lineSize);
}

/**
 * Create a new box for a given element
 */
function createBox(content, innerWidth, innerHeight, anchor=true, darkBG=true, hasShadow=true, padContent=true) {
    let _id = _getBoxId(content);
    let _box = _boxCache[_id];
    let _isNew = (!_box);
    if (_isNew) {
        _box = document.createElement('div');
        _boxCache[_id] = {
            'box': _box,
            'padded': padContent
        };

        let _addChild = function(childId) {
            let _child = document.createElement('div');
            _child.id = _getChildId(content, childId);
            _child.style.position = 'absolute';
            _box.appendChild(_child);
            return _child;
        }

        let _bgColor = _background;
        if (darkBG)
            _bgColor = '#000000';

        _box.id = _id;
        _box.style.position = 'fixed';
        _addChild('leftOutterShadow').style.backgroundColor = _outterShadow;
        _addChild('bottomOutterShadow').style.backgroundColor = _outterShadow;
        _addChild('background').style.backgroundColor = _bgColor;
        if (anchor)
            _box.appendChild(content);
        _addChild('leftBorder').style.backgroundColor = _border;
        _addChild('rightBorder').style.backgroundColor = _border;
        _addChild('topBorder').style.backgroundColor = _border;
        _addChild('bottomBorder').style.backgroundColor = _border;
        _addChild('leftInnerShadow').style.backgroundColor = _innerShadow;
        _addChild('rightInnerShadow').style.backgroundColor = _innerShadow;
        _addChild('topInnerShadow').style.backgroundColor = _innerShadow;
        _addChild('bottomInnerShadow').style.backgroundColor = _innerShadow;
        _addChild('leftOutterOutline').style.backgroundColor = _darkOutline;
        _addChild('rightOutterOutline').style.backgroundColor = _lightOutline;
        _addChild('topOutterOutline').style.backgroundColor = _lightOutline;
        _addChild('bottomOutterOutline').style.backgroundColor = _darkOutline;
        _addChild('leftInnerOutline').style.backgroundColor = _darkOutline;
        _addChild('rightInnerOutline').style.backgroundColor = _lightOutline;
        _addChild('topInnerOutline').style.backgroundColor = _lightOutline;
        _addChild('bottomInnerOutline').style.backgroundColor = _darkOutline;
    }

    if (_isNew)
        document.body.insertAdjacentElement('beforeend', _box)

    setBoxDimensions(content, innerWidth, innerHeight, padContent);

    if (!hasShadow) {
        document.getElementById(_getChildId('leftOutterShadow')).style.visibility = 'hidden';
        document.getElementById(_getChildId('bottomOutterShadow')).style.visibility = 'hidden';
    }
}

/** Set a box's position */
function setBoxPosition(content, x, y) {
    return setBoxPositionById(content.id, x, y)
}

function setBoxPositionById(id, x, y) {
    let _box = _boxCache[_getBoxIdById(id)].box;
    _box.style.left = x + 'px';
    _box.style.top = y + 'px';
    return _box;
}

/** Get a box's position */
function getBoxPosition(content) {
    let _box = _boxCache[_getBoxId(content)].box;
    return {
        'x': _box.offsetLeft,
        'y': _box.offsetTop
    }
}

/** Retrieve the position of the element within a box. NOTE: The position is the same for both axis) */
function getBoxContentPosition(content) {
    let _box = _boxCache[_getBoxId(content)];
    let _pos = 2 * _lineSize + _borderSize;
    if (_box.padded)
        _pos += _contentPadding;
    return _pos;
}

/** Retrieves an object with 'x' and 'y' fields with the absolute position of the element within the box */
function getBoxContentAbsolutePosition(content) {
    let _box = _boxCache[_getBoxId(content)].box;
    let _pos = getBoxContentPosition(content);
    return {
        'x': _pos + _box.offsetLeft,
        'y': _pos + _box.offsetTop
    };
}

function hideBox(content) {
    return hideBoxById(content.id);
}

function hideBoxById(id) {
    let _box = _boxCache[_getBoxIdById(id)].box;
    _box.style.visibility = 'hidden';
}

function showBox(content) {
    showBoxById(content.id);
}

function showBoxById(id) {
    let _box = _boxCache[_getBoxIdById(id)].box;
    _box.style.visibility = 'visible';
}

function isBoxVisibleById(id) {
    let _box = _boxCache[_getBoxIdById(id)].box;
    return _box.style.visibility == 'visible';
}

function getBoxHeightById(id) {
    let _box = _boxCache[_getBoxIdById(id)].box;
    return _box.scrollHeight;
}

function getBoxWidthById(id) {
    let _box = _boxCache[_getBoxIdById(id)].box;
    return _box.scrollWidth;
}
