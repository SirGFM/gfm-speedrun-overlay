/**
 * Handles the '/speedrun' display.
 */

const windowWidth = 1280;
const windowHeight = 720;

const _defViewTop = 8;
/** Left view width (timer, gamepad etc) */
const _leftViewWidth = 240;
/** Right view width (game area) */
const _rightViewWidth = 1040;
/** Max border, considering a 16:9 resolution (max multi == 64) */
const _16x9MaxWidth = 1024;
const _16x9MaxHeight = 576;
/** Max border, considering a 4:3 resolution (max multi == 240) */
const _4x3MaxWidth = 960;
const _4x3MaxHeight = 720;
/** Max border, considering an unrestricted resolution */
const _maxWidth = 1024;
const _maxHeight = 720;
/* Width of the label's border */
const _outlineWidth = 0.2;

/* Display loaded from the config file */
let _display;

/**
 * Retrieves a console's dimensions from its name.
 *
 * @param{name} A string with the console's name (or initials).
 * @return{Array} The dimensions for the consoles, as [width, height].
 */
let getConsoleDimension = function(name) {
    switch (name.toLowerCase()) {
    case "nes":
        return [879, 695];
    case "gb":
    case "gbc":
        return [640, 576];
    case "megadrive":
    case "md":
        return [960, 695];
    case "mastersystem":
    case "sms":
        return [879, 599];
    case "gamegear":
    case "gg":
        return [640, 599];
    case "gba":
        return [720, 480];
    case "snes":
        return [879, 672];
    default:
        throw "Invalid display mode. Check the docs!";
    }
}

/**
 * Convert the value into a valid dimension.
 *
 * @param{value} The tentative dimension.
 * @return{Int} The calculated dimension.
 */
function getValidDimension(value) {
    return makeValueEven(Math.floor(value));
}

/**
 * Retrieves the dimensions of the display as an object with attributes:
 *   - border
 *       - width
 *       - height
 *       - x
 *       - y
 *   - game
 *       - width
 *       - height
 *       - x
 *       - y
 *
 * @param{display} A 'display' object. Look at docs/config.md for more info.
 * @return An object as described above
 */
let _setDimensions = function(display) {
    let borderW = 0, borderH = 0, maxW = 0, maxH = 0, w = 0, h = 0;
    let padContent = true;

    /* 0 - Sanitize display size */
    if (display.type == "console") {
        let arr = getConsoleDimension(display.mode);
        w = arr[0];
        h = arr[1];
    }
    else if (display.type == "custom") {
        w = display.width;
        h = display.height;
    }
    else {
        throw "Invalid display object. Check the docs!";
    }

    if (Math.floor(w / 16) == Math.floor(h / 9)) {
        maxW = _16x9MaxWidth;
        maxH = _16x9MaxHeight;
    }
    else if (Math.floor(w / 4) == Math.floor(h / 3)) {
        maxW = _4x3MaxWidth;
        maxH = _4x3MaxHeight;
    }
    else {
        maxW = _maxWidth;
        maxH = _maxHeight;
    }

    /* Check whether the dimensions fit within the border */
    do {
        borderW = getBoxDimension(w, padContent);
        borderH = getBoxDimension(h, padContent);
        if (borderW > maxW || borderH > maxH) {
            w = w * 0.9;
            h = h * 0.9;
            continue;
        }
    } while (false);

    /* Correct the game view (if needed) */
    w = getValidDimension(w);
    h = getValidDimension(h);

    /* Create (or resize) the game's view */
    let view = document.getElementById('gameView');
    if (!view) {
        view = document.createElement('div');
        view.id = 'gameView';
    }

    /* Center within the window, unless it would overlap the left area */
    let borderX = (windowWidth - borderW) * 0.5;
    let borderY = _defViewTop;
    if (borderH >= maxH) {
        borderY = (windowHeight - borderH) * 0.5;
        getValidDimension(borderY);
    }
    if (borderX <= _leftViewWidth) {
        borderX = _leftViewWidth + (_rightViewWidth - borderW) * 0.5;
    }
    borderY = getValidDimension(borderY);
    borderX = getValidDimension(borderX);

    createBox(view, w, h, anchor=true, darkBG=true, hasShadow=true, padContent);
    setBoxPosition(view, borderX, borderY);
    pos = getBoxContentAbsolutePosition(view);

    return {
        border: {
            x: borderX,
            y: borderY,
            width: borderW,
            height: borderH
        },
        game: {
            x: pos.x,
            y: pos.y,
            width: w,
            height: h
        }
    };
}

/**
 * Configures 'gameView', the border, and 'gameColorKey', the actual game area,
 * based on a display object.
 *
 * @param{display} A 'display' object. Look at docs/config.md for more info.
 * @return An object as described above
 */
function configureDisplay(argsDisplay) {
    _display = _setDimensions(argsDisplay);

    /* Write the game info into the screen */
    let str = "&nbsp game position:&nbsp</br>&nbsp&nbsp{";
    str += " x: "+_display.game.x;
    str += ", y: "+_display.game.y;
    str += ",&nbsp</br>&nbsp&nbsp w: "+_display.game.width;
    str += ", h: "+_display.game.height;
    str += "&nbsp}&nbsp";
    document.getElementById("gameInfo").innerHTML = str;
}

/**
 * Calculate how many pixels are left under the border (e.g., for the title).
 */
function getUnderView() {
    let val = windowHeight - _display.border.y - _display.border.height;
    return getValidDimension(val);
}

/**
 * Calculate the actual width of the left side of the screen.
 */
function getLeftView() {
    if (_display.border.x * 0.9 > _leftViewWidth) {
        return _display.border.x * 0.9;
    }
    return _leftViewWidth;
}

function getViewInfo() {
    return {
        'x': _display.border.x,
        'y': _display.border.y,
        'width': _display.border.width,
        'height': _display.border.height
    };
}
