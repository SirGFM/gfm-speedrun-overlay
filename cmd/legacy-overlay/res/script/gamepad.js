let gamepad = function() {
    const _fps = 30;
    const xbox_id = 'Xbox 360 Controller';
    const xbox_id_linux = 'Vendor: 045e Product: 028e';
    const ps1_id = 'Vendor: 0810 Product: 0001';
    const gcn_id = 'Vendor: 0079 Product: 1846';
    const ps5_id = 'Vendor: 054c Product: 0ce6';

    let _eventBt = null;
    let _lastState = false;
    let _images = [];
    let _active = null;
    let _activeGpIdx = -1;
    let _interval = null;
    let _defaultSkin = null;

    let _xboxSkin = {
        'released': {
            'src': '/img/gamepad/x360/released_buttons.png',
            'width': 190,
            'left': 18,
            'top': 18
        },
        'button': {
            'l2': {
                'src': '/img/gamepad/x360/pressed/l2.png',
                'width': 42,
                'height': 20,
                'left': 22,
                'top': 18,
                'button': 6
            },
            'l1': {
                'src': '/img/gamepad/x360/pressed/l1.png',
                'width': 34,
                'height': 14,
                'left': 62,
                'top': 28,
                'button': 4
            },
            'r2': {
                'src': '/img/gamepad/x360/pressed/r2.png',
                'width': 42,
                'height': 20,
                'left': 156,
                'top': 18,
                'button': 7
            },
            'r1': {
                'src': '/img/gamepad/x360/pressed/r1.png',
                'width': 34,
                'height': 14,
                'left': 122,
                'top': 28,
                'button': 5
            },
            'home': {
                'src': '/img/gamepad/x360/pressed/home.png',
                'width': 26,
                'height': 26,
                'left': 96,
                'top': 46,
                'button': 16
            },
            'select': {
                'src': '/img/gamepad/x360/pressed/select.png',
                'width': 24,
                'height': 14,
                'left': 74,
                'top': 70,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/x360/pressed/start.png',
                'width': 24,
                'height': 14,
                'left': 120,
                'top': 70,
                'button': 9
            },
            'x': {
                'src': '/img/gamepad/x360/pressed/x.png',
                'width': 26,
                'height': 26,
                'left': 146,
                'top': 66,
                'button': 2
            },
            'y': {
                'src': '/img/gamepad/x360/pressed/y.png',
                'width': 26,
                'height': 26,
                'left': 168,
                'top': 44,
                'button': 3
            },
            'a': {
                'src': '/img/gamepad/x360/pressed/a.png',
                'width': 26,
                'height': 26,
                'left': 160,
                'top': 92,
                'button': 0
            },
            'b': {
                'src': '/img/gamepad/x360/pressed/b.png',
                'width': 26,
                'height': 26,
                'left': 182,
                'top': 70,
                'button': 1
            },
            'up': {
                'src': '/img/gamepad/x360/pressed/up.png',
                'width': 14,
                'height': 14,
                'left': 82,
                'top': 102,
                'button': 12
            },
            'down': {
                'src': '/img/gamepad/x360/pressed/down.png',
                'width': 14,
                'height': 14,
                'left': 82,
                'top': 126,
                'button': 13
            },
            'left': {
                'src': '/img/gamepad/x360/pressed/left.png',
                'width': 14,
                'height': 14,
                'left': 70,
                'top': 114,
                'button': 14
            },
            'right': {
                'src': '/img/gamepad/x360/pressed/right.png',
                'width': 14,
                'height': 14,
                'left': 94,
                'top': 114,
                'button': 15
            }
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/x360/stick.png',
                'width': 18,
                'height': 18,
                'left': 36,
                'top': 80,
                'hw': 25,
                'hh': 25,
                'hor': 0,
                'ver': 1
            },
            'rstick': {
                'src': '/img/gamepad/x360/stick.png',
                'width': 18,
                'height': 18,
                'left': 126,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 2,
                'ver': 3
            }
        }
    };
    let _ps1Skin = {
        'released': {
            'src': '/img/gamepad/ps1/released.png',
            'width': 190,
            'left': 18,
            'top': 18
        },
        'button': {
            'l2': {
                'src': '/img/gamepad/ps1/pressed/l2.png',
                'width': 36,
                'height': 20,
                'left': 32,
                'top': 18,
                'button': 4
            },
            'l1': {
                'src': '/img/gamepad/ps1/pressed/l1.png',
                'width': 34,
                'height': 16,
                'left': 68,
                'top': 38,
                'button': 6
            },
            'r2': {
                'src': '/img/gamepad/ps1/pressed/r2.png',
                'width': 36,
                'height': 20,
                'left': 156,
                'top': 18,
                'button': 5
            },
            'r1': {
                'src': '/img/gamepad/ps1/pressed/r1.png',
                'width': 34,
                'height': 16,
                'left': 124,
                'top': 38,
                'button': 7
            },
            'select': {
                'src': '/img/gamepad/ps1/pressed/select.png',
                'width': 24,
                'height': 12,
                'left': 84,
                'top': 68,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/ps1/pressed/start.png',
                'width': 24,
                'height': 12,
                'left': 114,
                'top': 68,
                'button': 9
            },
            'square': {
                'src': '/img/gamepad/ps1/pressed/square.png',
                'width': 28,
                'height': 28,
                'left': 140,
                'top': 72,
                'button': 3
            },
            'triangle': {
                'src': '/img/gamepad/ps1/pressed/triangle.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 52,
                'button': 0
            },
            'cross': {
                'src': '/img/gamepad/ps1/pressed/cross.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 92,
                'button': 2
            },
            'circle': {
                'src': '/img/gamepad/ps1/pressed/circle.png',
                'width': 28,
                'height': 28,
                'left': 180,
                'top': 72,
                'button': 1
            },
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 80,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 4,
                'ver': 5
            },
            'rstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 126,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 3,
                'ver': 2
            }
        },
        'hat': {
            'up': {
                'src': '/img/gamepad/ps1/pressed/up.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 52,
                'button': 12,
                'idx': 1,
                'val': -0.75
            },
            'down': {
                'src': '/img/gamepad/ps1/pressed/down.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 90,
                'button': 13,
                'idx': 1,
                'val': 0.75
            },
            'left': {
                'src': '/img/gamepad/ps1/pressed/left.png',
                'width': 26,
                'height': 24,
                'left': 18,
                'top': 72,
                'button': 14,
                'idx': 0,
                'val': -0.75
            },
            'right': {
                'src': '/img/gamepad/ps1/pressed/right.png',
                'width': 26,
                'height': 24,
                'left': 56,
                'top': 72,
                'button': 15,
                'idx': 0,
                'val': 0.75
            }
        }
    };
    let _ps1AnalogSkin = {
        'released': {
            'src': '/img/gamepad/ps1/released.png',
            'width': 190,
            'left': 18,
            'top': 18
        },
        'button': {
            'l2': {
                'src': '/img/gamepad/ps1/pressed/l2.png',
                'width': 36,
                'height': 20,
                'left': 32,
                'top': 18,
                'button': 4
            },
            'l1': {
                'src': '/img/gamepad/ps1/pressed/l1.png',
                'width': 34,
                'height': 16,
                'left': 68,
                'top': 38,
                'button': 6
            },
            'r2': {
                'src': '/img/gamepad/ps1/pressed/r2.png',
                'width': 36,
                'height': 20,
                'left': 156,
                'top': 18,
                'button': 5
            },
            'r1': {
                'src': '/img/gamepad/ps1/pressed/r1.png',
                'width': 34,
                'height': 16,
                'left': 124,
                'top': 38,
                'button': 7
            },
            'select': {
                'src': '/img/gamepad/ps1/pressed/select.png',
                'width': 24,
                'height': 12,
                'left': 84,
                'top': 68,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/ps1/pressed/start.png',
                'width': 24,
                'height': 12,
                'left': 114,
                'top': 68,
                'button': 9
            },
            'square': {
                'src': '/img/gamepad/ps1/pressed/square.png',
                'width': 28,
                'height': 28,
                'left': 140,
                'top': 72,
                'button': 3
            },
            'triangle': {
                'src': '/img/gamepad/ps1/pressed/triangle.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 52,
                'button': 0
            },
            'cross': {
                'src': '/img/gamepad/ps1/pressed/cross.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 92,
                'button': 2
            },
            'circle': {
                'src': '/img/gamepad/ps1/pressed/circle.png',
                'width': 28,
                'height': 28,
                'left': 180,
                'top': 72,
                'button': 1
            },
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 80,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 0,
                'ver': 1
            },
            'rstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 126,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 3,
                'ver': 2
            }
        },
        'hat': {
            'up': {
                'src': '/img/gamepad/ps1/pressed/up.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 52,
                'button': 12,
                'idx': 5,
                'val': -0.75
            },
            'down': {
                'src': '/img/gamepad/ps1/pressed/down.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 90,
                'button': 13,
                'idx': 5,
                'val': 0.75
            },
            'left': {
                'src': '/img/gamepad/ps1/pressed/left.png',
                'width': 26,
                'height': 24,
                'left': 18,
                'top': 72,
                'button': 14,
                'idx': 4,
                'val': -0.75
            },
            'right': {
                'src': '/img/gamepad/ps1/pressed/right.png',
                'width': 26,
                'height': 24,
                'left': 56,
                'top': 72,
                'button': 15,
                'idx': 4,
                'val': 0.75
            }
        }
    };
    let _ps1AnalogWinSkin = {
        'released': {
            'src': '/img/gamepad/ps1/released.png',
            'width': 190,
            'left': 18,
            'top': 18
        },
        'button': {
            'l2': {
                'src': '/img/gamepad/ps1/pressed/l2.png',
                'width': 36,
                'height': 20,
                'left': 32,
                'top': 18,
                'button': 4
            },
            'l1': {
                'src': '/img/gamepad/ps1/pressed/l1.png',
                'width': 34,
                'height': 16,
                'left': 68,
                'top': 38,
                'button': 6
            },
            'r2': {
                'src': '/img/gamepad/ps1/pressed/r2.png',
                'width': 36,
                'height': 20,
                'left': 156,
                'top': 18,
                'button': 5
            },
            'r1': {
                'src': '/img/gamepad/ps1/pressed/r1.png',
                'width': 34,
                'height': 16,
                'left': 124,
                'top': 38,
                'button': 7
            },
            'select': {
                'src': '/img/gamepad/ps1/pressed/select.png',
                'width': 24,
                'height': 12,
                'left': 84,
                'top': 68,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/ps1/pressed/start.png',
                'width': 24,
                'height': 12,
                'left': 114,
                'top': 68,
                'button': 9
            },
            'square': {
                'src': '/img/gamepad/ps1/pressed/square.png',
                'width': 28,
                'height': 28,
                'left': 140,
                'top': 72,
                'button': 3
            },
            'triangle': {
                'src': '/img/gamepad/ps1/pressed/triangle.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 52,
                'button': 0
            },
            'cross': {
                'src': '/img/gamepad/ps1/pressed/cross.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 92,
                'button': 2
            },
            'circle': {
                'src': '/img/gamepad/ps1/pressed/circle.png',
                'width': 28,
                'height': 28,
                'left': 180,
                'top': 72,
                'button': 1
            },
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 80,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 0,
                'ver': 1
            },
            'rstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 126,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 5,
                'ver': 2
            }
        },
        'hat': {
            'up': {
                'src': '/img/gamepad/ps1/pressed/up.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 52,
                'button': 12,
                'idx': 9,
                'approx': [-1, 1, -0.714]
            },
            'down': {
                'src': '/img/gamepad/ps1/pressed/down.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 90,
                'button': 13,
                'idx': 9,
                'approx': [0.14, 0.42, -0.14]
            },
            'left': {
                'src': '/img/gamepad/ps1/pressed/left.png',
                'width': 26,
                'height': 24,
                'left': 18,
                'top': 72,
                'button': 14,
                'idx': 9,
                'approx': [0.714, 1, 0.42]
            },
            'right': {
                'src': '/img/gamepad/ps1/pressed/right.png',
                'width': 26,
                'height': 24,
                'left': 56,
                'top': 72,
                'button': 15,
                'idx': 9,
                'approx': [-0.42, -0.714, -0.14]
            }
        }
    };
    let _gcnSkin = {
        'released': {
            'src': '/img/gamepad/gcn/released_buttons.png',
            'width': 178,
            'left': 23,
            'top': 18
        },
        'button': {
            'l2': {
                /* NOTE: unused */
                'src': '/img/gamepad/x360/pressed/l2.png',
                'width': 42,
                'height': 20,
                'left': 22,
                'top': 18,
                'button': 6
            },
            'l1': {
                'src': '/img/gamepad/gcn/pressed/l.png',
                'width': 44,
                'height': 20,
                'left': 35,
                'top': 18,
                'button': 4
            },
            'r2': {
                'src': '/img/gamepad/gcn/pressed/z.png',
                'width': 32,
                'height': 16,
                'left': 111,
                'top': 24,
                'button': 7
            },
            'r1': {
                'src': '/img/gamepad/gcn/pressed/r.png',
                'width': 44,
                'height': 20,
                'left': 155,
                'top': 18,
                'button': 5
            },
            'home': {
                /* NOTE: unused */
                'src': '/img/gamepad/x360/pressed/home.png',
                'width': 26,
                'height': 26,
                'left': 96,
                'top': 46,
                'button': 16
            },
            'select': {
                /* NOTE: unused */
                'src': '/img/gamepad/x360/pressed/select.png',
                'width': 24,
                'height': 14,
                'left': 74,
                'top': 70,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/gcn/pressed/start.png',
                'width': 22,
                'height': 20,
                'left': 99,
                'top': 48,
                'button': 9
            },
            'x': {
                'src': '/img/gamepad/gcn/pressed/x.png',
                'width': 18,
                'height': 30,
                'left': 183,
                'top': 54,
                'button': 0
            },
            'y': {
                'src': '/img/gamepad/gcn/pressed/y.png',
                'width': 30,
                'height': 18,
                'left': 149,
                'top': 44,
                'button': 3
            },
            'a': {
                'src': '/img/gamepad/gcn/pressed/a.png',
                'width': 26,
                'height': 26,
                'left': 157,
                'top': 62,
                'button': 1
            },
            'b': {
                'src': '/img/gamepad/gcn/pressed/b.png',
                'width': 24,
                'height': 24,
                'left': 135,
                'top': 78,
                'button': 2
            },
            'up': {
                'src': '/img/gamepad/gcn/pressed/up.png',
                'width': 14,
                'height': 14,
                'left': 87,
                'top': 102,
                'button': 12
            },
            'down': {
                'src': '/img/gamepad/gcn/pressed/down.png',
                'width': 14,
                'height': 14,
                'left': 87,
                'top': 126,
                'button': 14
            },
            'left': {
                'src': '/img/gamepad/gcn/pressed/left.png',
                'width': 14,
                'height': 14,
                'left': 75,
                'top': 114,
                'button': 15
            },
            'right': {
                'src': '/img/gamepad/gcn/pressed/right.png',
                'width': 14,
                'height': 14,
                'left': 99,
                'top': 114,
                'button': 13
            }
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/gcn/l_stick.png',
                'width': 18,
                'height': 18,
                'left': 41,
                'top': 80,
                'hw': 25,
                'hh': 25,
                'hor': 0,
                'ver': 1,
            },
            'rstick': {
                'src': '/img/gamepad/gcn/c_stick.png',
                'width': 18,
                'height': 18,
                'left': 131,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 5,
                'ver': 2,
                'invert_x': true,
                'invert_y': true
            }
        }
    };
    let _ps5Skin = {
        'released': {
            'src': '/img/gamepad/ps1/released.png',
            'width': 190,
            'left': 18,
            'top': 18
        },
        'button': {
            'l1': {
                'src': '/img/gamepad/ps1/pressed/l1.png',
                'width': 34,
                'height': 16,
                'left': 68,
                'top': 38,
                'button': 4
            },
            'r1': {
                'src': '/img/gamepad/ps1/pressed/r1.png',
                'width': 34,
                'height': 16,
                'left': 124,
                'top': 38,
                'button': 5
            },
            'select': {
                'src': '/img/gamepad/ps1/pressed/select.png',
                'width': 24,
                'height': 12,
                'left': 84,
                'top': 68,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/ps1/pressed/start.png',
                'width': 24,
                'height': 12,
                'left': 114,
                'top': 68,
                'button': 9
            },
            'square': {
                'src': '/img/gamepad/ps1/pressed/square.png',
                'width': 28,
                'height': 28,
                'left': 140,
                'top': 72,
                'button': 3
            },
            'triangle': {
                'src': '/img/gamepad/ps1/pressed/triangle.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 52,
                'button': 0
            },
            'cross': {
                'src': '/img/gamepad/ps1/pressed/cross.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 92,
                'button': 2
            },
            'circle': {
                'src': '/img/gamepad/ps1/pressed/circle.png',
                'width': 28,
                'height': 28,
                'left': 180,
                'top': 72,
                'button': 1
            },
            'up': {
                'src': '/img/gamepad/ps1/pressed/up.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 52,
                'button': 12
            },
            'down': {
                'src': '/img/gamepad/ps1/pressed/down.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 90,
                'button': 13
            },
            'left': {
                'src': '/img/gamepad/ps1/pressed/left.png',
                'width': 26,
                'height': 24,
                'left': 18,
                'top': 72,
                'button': 14
            },
            'right': {
                'src': '/img/gamepad/ps1/pressed/right.png',
                'width': 26,
                'height': 24,
                'left': 56,
                'top': 72,
                'button': 15
            }
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 80,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 0,
                'ver': 1
            },
            'rstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 126,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 6, //button 6 (0: left, 1: right) and button 7 (0: top, 1: bot)
                'ver': 7,
                'from-button': true
            }
        },
        'hat': {
            'l2': {
                'src': '/img/gamepad/ps1/pressed/l2.png',
                'width': 36,
                'height': 20,
                'left': 32,
                'top': 18,
                'idx': 2,
                'val': 0.25
            },
            'r2': {
                'src': '/img/gamepad/ps1/pressed/r2.png',
                'width': 36,
                'height': 20,
                'left': 156,
                'top': 18,
                'idx': 3,
                'val': 0.25
            },
        }
    };
    let _ps5Skin_win = {
        'released': {
            'src': '/img/gamepad/ps1/released.png',
            'width': 190,
            'left': 18,
            'top': 18
        },
        'button': {
            'l1': {
                'src': '/img/gamepad/ps1/pressed/l1.png',
                'width': 34,
                'height': 16,
                'left': 68,
                'top': 38,
                'button': 4
            },
            'r1': {
                'src': '/img/gamepad/ps1/pressed/r1.png',
                'width': 34,
                'height': 16,
                'left': 124,
                'top': 38,
                'button': 5
            },
            'l2': {
                'src': '/img/gamepad/ps1/pressed/l2.png',
                'width': 36,
                'height': 20,
                'left': 32,
                'top': 18,
                'button': 6
            },
            'r2': {
                'src': '/img/gamepad/ps1/pressed/r2.png',
                'width': 36,
                'height': 20,
                'left': 156,
                'top': 18,
                'button': 7
            },
            'select': {
                'src': '/img/gamepad/ps1/pressed/select.png',
                'width': 24,
                'height': 12,
                'left': 84,
                'top': 68,
                'button': 8
            },
            'start': {
                'src': '/img/gamepad/ps1/pressed/start.png',
                'width': 24,
                'height': 12,
                'left': 114,
                'top': 68,
                'button': 9
            },
            'square': {
                'src': '/img/gamepad/ps1/pressed/square.png',
                'width': 28,
                'height': 28,
                'left': 140,
                'top': 72,
                'button': 0
            },
            'triangle': {
                'src': '/img/gamepad/ps1/pressed/triangle.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 52,
                'button': 3
            },
            'cross': {
                'src': '/img/gamepad/ps1/pressed/cross.png',
                'width': 28,
                'height': 28,
                'left': 160,
                'top': 92,
                'button': 1
            },
            'circle': {
                'src': '/img/gamepad/ps1/pressed/circle.png',
                'width': 28,
                'height': 28,
                'left': 180,
                'top': 72,
                'button': 2
            },
            'up': {
                'src': '/img/gamepad/ps1/pressed/up.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 52,
                'button': 12
            },
            'down': {
                'src': '/img/gamepad/ps1/pressed/down.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 90,
                'button': 13
            },
            'left': {
                'src': '/img/gamepad/ps1/pressed/left.png',
                'width': 26,
                'height': 24,
                'left': 18,
                'top': 72,
                'button': 14
            },
            'right': {
                'src': '/img/gamepad/ps1/pressed/right.png',
                'width': 26,
                'height': 24,
                'left': 56,
                'top': 72,
                'button': 15
            }
        },
        'axis': {
            'lstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 80,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 0,
                'ver': 1
            },
            'rstick': {
                'src': '/img/gamepad/ps1/stick.png',
                'width': 18,
                'height': 18,
                'left': 126,
                'top': 112,
                'hw': 17,
                'hh': 17,
                'hor': 2,
                'ver': 5
            }
        },
        'hat': {
            'up': {
                'src': '/img/gamepad/ps1/pressed/up.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 52,
                'button': 12,
                'idx': 9,
                'approx': [-1, 1, -0.714]
            },
            'down': {
                'src': '/img/gamepad/ps1/pressed/down.png',
                'width': 24,
                'height': 26,
                'left': 38,
                'top': 90,
                'button': 13,
                'idx': 9,
                'approx': [0.14, 0.42, -0.14]
            },
            'left': {
                'src': '/img/gamepad/ps1/pressed/left.png',
                'width': 26,
                'height': 24,
                'left': 18,
                'top': 72,
                'button': 14,
                'idx': 9,
                'approx': [0.714, 1, 0.42]
            },
            'right': {
                'src': '/img/gamepad/ps1/pressed/right.png',
                'width': 26,
                'height': 24,
                'left': 56,
                'top': 72,
                'button': 15,
                'idx': 9,
                'approx': [-0.42, -0.714, -0.14]
            }
        }
    };

    /** Reset every image, so a new skin may be loaded. */
    let resetGamepad = function() {
        for (let i in _images) {
            _images[i].style.visibility = 'hidden';
            _images[i].in_use = false;
        }
    };

    /**
     * Retrieve the current list of connected gamepads.
     *
     * From: https://developer.mozilla.org/en-US/docs/Web/API/Gamepad_API/Using_the_Gamepad_API
     */
    function getGamepadList() {
        if (navigator.getGamepads)
           return navigator.getGamepads();
        else if (navigator.webkitGetGamepads)
           return navigator.webkitGetGamepads()
        return [];
    }

    /** Set a callback for updating the buttons */
    function enableGamepad(gamepadId) {
        if (_interval === null && getGamepadCount(getGamepadList()) > 0) {
            /* Configure the gamepad skin */
            let skin = _defaultSkin;
            if (!skin) {
                if (gamepadId.startsWith(xbox_id) ||
                        gamepadId.indexOf(xbox_id_linux) != -1)
                    skin = 'xbox';
                else if (gamepadId.indexOf(ps1_id) != -1)
                    skin = 'ps1';
                else if (gamepadId.indexOf(gcn_id) != -1)
                    skin = 'gcn';
                else if (gamepadId.indexOf(ps5_id) != -1) {
                    if (navigator.userAgent.indexOf('Windows') != -1)
                        skin = 'ps5-win';
                    else
                        skin = 'ps5';
                }
            }
            let parentContent = document.getElementById('gamepad-view');
            gamepad.setup(parentContent, skin);
            /* No 'onpress' event, gotta pool... */
            _interval = setInterval(pollGamepad, 1000 / _fps);
        }
    }

    /** Remove the callback for updating the buttons */
    function disableGamepad(force = false) {
        if (force || (getGamepadCount(getGamepadList()) <= 0 &&
                      _interval !== null)) {
            clearInterval(_interval);
            _interval = null;
        }
    }

    /**
     * Retrieve the actual number of connected gamepads, given a list retrieved
     * from 'getGamepadList()'.
     */
    function getGamepadCount(gpList) {
        let count = 0;

        for (let i in gpList) {
            if (gpList[i])
                count++;
        }

        return count;
    }


    /**
     * Add a new <img> to the gamepad. 'obj' must have the following fields:
     *   - src
     *   - width
     *   - left
     *   - top
     * Optionally, 'obj' may have a 'height' attribute.
     *
     * @param{obj} The object to be added.
     * @return The newly created <img>.
     */
    let addObject = function(obj, parentContent) {
        let _new = null;
        for (let i in _images) {
            if (!_images[i].in_use) {
                _new = _images[i];
                break;
            }
        }
        if (!_new) {
            _new = document.createElement('img');
            _new.style.visibility = 'hidden';
            _new.style.position = 'absolute';
            /* Add it to the document. */
            parentContent.appendChild(_new);
            /* Cache it for later use. */
            _images.push(_new);
        }
        _new.src = obj.src;
        _new.style.width = obj.width + 'px';
        if ('height' in obj)
            _new.style.height = obj.height + 'px';
        _new.style.left = obj.left + 'px';
        _new.style.top = obj.top + 'px';
        _new.in_use = true;
        return _new;
    };

    /** Check whether 'b' is pressed */
    let buttonPressed = function(b) {
        if (typeof(b) == "object") {
            return b.pressed;
        }
        return b == 1.0;
    }

    /** Update the gamepad's objects */
    let pollGamepad = function() {
        if (!_active)
            return;
        let _gps = getGamepadList();
        /* Shouldn't happen... */
        if (!_gps || _gps.length <= 0)
            return;

        /** Check for the first gamepad with any button pressed */
_check_gp:
        if (_activeGpIdx == -1) {
            for (let i  = 0; i < _gps.length; i++) {
                if (!_gps[i])
                    continue;
                for (let bt in _gps[i].buttons) {
                    if (buttonPressed(_gps[i].buttons[bt])) {
                        _activeGpIdx = i;
                        break _check_gp;
                    }
                }
            }
        }
        if (_activeGpIdx == -1 || _activeGpIdx >= _gps.length ||
                !_gps[_activeGpIdx]) {
            /* Gamepad was just removed */
            _activeGpIdx = -1;
            return;
        }
        _gp = _gps[_activeGpIdx];

        for (let i in _active.button) {
            let _bt = _active.button[i];
            let idx = _bt.button;

            if (idx >= _gp.buttons.length)
                continue;

            let state = buttonPressed(_gp.buttons[idx]);
            if (state != _bt.visible) {
                _bt.visible = state;
                if (state) {
                    _bt.img.style.visibility = 'visible';
                }
                else {
                    _bt.img.style.visibility = 'hidden';
                }
            }

            gamepad.dispatchTimerEvent(i, state);
        }
        for (let i in _active.hat) {
            let _hat = _active.hat[i];

            let state = false;
            if ('approx' in _hat) {
                for (let j in _hat.approx) {
                    if (Math.abs(_hat.approx[j] - _gp.axes[_hat.idx]) < 0.01) {
                        state = true;
                        break;
                    }
                }
            }
            else {
                if (_hat.val > 0)
                    state = _gp.axes[_hat.idx] >= _hat.val;
                else
                    state = _gp.axes[_hat.idx] <= _hat.val;
            }
            if (state != _hat.visible) {
                _hat.visible = state;
                if (state) {
                    _hat.img.style.visibility = 'visible';
                }
                else {
                    _hat.img.style.visibility = 'hidden';
                }
            }

            gamepad.dispatchTimerEvent(i, state);
        }
        for (let i in _active.axis) {
            let _axis = _active.axis[i];
            let x = 0.0;
            let y = 0.0;
            if ('from-button' in _axis) {
                x = 2.0 * _gp.buttons[_axis.hor].value - 1.0;
                y = 2.0 * _gp.buttons[_axis.ver].value - 1.0;
            }
            else {
                x = _gp.axes[_axis.hor];
                y = _gp.axes[_axis.ver];
            }

            if (_axis.invert_x)
                x *= -1;
            if (_axis.invert_y)
                y *= -1;

            x *= _axis.hw;
            x += _axis.cx - _axis.img.width / 2;
            y *= _axis.hh;
            y += _axis.cy - _axis.img.height / 2;

            _axis.img.style.left = x + 'px';
            _axis.img.style.top = y + 'px';
        }
    };

    /** Configure the event listener for gamepads. */
    document.addEventListener("DOMContentLoaded", function (event) {
        /* Delay it unti ~1s after the page has loaded */
        window.setTimeout(function() {
            window.addEventListener("gamepadconnected", function(e) {
                _activeGpIdx = -1;
                enableGamepad(e.gamepad.id);
            });
            window.addEventListener("gamepaddisconnected", function(e) {
                disableGamepad();
            })
        }, 1000);
    });

    return {
        /**
         * Set a skin as the default one. Otherwise, a skin will be selected
         * from the gamepad id.
         *
         * @param{skin} The name of the skin.
         */
        setDefaultSkin(skin) {
            _defaultSkin = skin;
        },
        /**
         * Configure the view and bind the buttons according to a gamepad skin.
         *
         * @param{parentContent} Attach point for the images.
         * @param{skin} Name of the skin to be used.
         */
        setup(parentContent, skin='') {
            let _obj = null;
            switch (skin) {
            case 'xbox':
                _obj = _xboxSkin;
                break;
            case 'ps1':
                _obj = _ps1Skin;
                break;
            case 'ps1-analog':
                _obj = _ps1AnalogSkin;
                break;
            case 'ps1-analog-win':
                _obj = _ps1AnalogWinSkin;
                break;
            case 'gcn':
                _obj = _gcnSkin;
                break;
            case 'ps5':
                _obj = _ps5Skin;
                break;
            case 'ps5-win':
                _obj = _ps5Skin_win;
                break;
            case true:
                return;
            default:
                _obj = _xboxSkin;
                break;
            }

            resetGamepad();
            addObject(_obj.released, parentContent).style.visibility = 'visible';
            for (let i in _obj.button)
                _obj.button[i].img = addObject(_obj.button[i], parentContent);
            for (let i in _obj.axis) {
                _obj.axis[i].img = addObject(_obj.axis[i], parentContent);
                _obj.axis[i].img.style.visibility = 'visible';
                _obj.axis[i].cx = _obj.axis[i].img.offsetLeft;
                _obj.axis[i].cx += _obj.axis[i].img.width / 2;
                _obj.axis[i].cy = _obj.axis[i].img.offsetTop;
                _obj.axis[i].cy += _obj.axis[i].img.height / 2;
            }
            for (let i in _obj.hat)
                _obj.hat[i].img = addObject(_obj.hat[i], parentContent);

            _active = _obj;
        },
        /**
         * Try to dispatch a timer control event, based on a button press.
         *
         * @param{buttonName} Name of the button that changed states.
         * @param{state} State of the button (true = pressed, false = released).
         */
        dispatchTimerEvent: function(buttonName, state) {
            if (_eventBt != buttonName)
                return; /* Do nothing */
            else if (_lastState == state) {
                if (state)
                    document.dispatchEvent(new Event('timer-pressed'));
            }
            else if (state)
                document.dispatchEvent(new Event('timer-onpress'));
            else
                document.dispatchEvent(new Event('timer-onrelease'));
            _lastState = state;
        },
        /**
         * Configure a button to generate the following events:
         *   - 'timer-onpress'
         *   - 'timer-pressed'
         *   - 'timer-onrelease'
         *
         * @param{buttonName} Name of the button that shall trigger the events.
         */
        setTimerEventButton: function(buttonName) {
            _eventBt = buttonName;
            _lastState = false;
        },
        logKeys: function() {
            let _gps = getGamepadList();
            /* Shouldn't happen... */
            if (!_gps || _gps.length <= 0)
                return;
            let _gp = null;
            for (let i in _gps)
                if (_gps[i]) {
                    _gp = _gps[i];
                    break;
                }
            if (!_gp) {
                window.alert('No gamepad!');
                return;
            }
            let s = '';
            for (let i in _gp.buttons)
                s += 'b: ' + i + ' - ' + _gp.buttons[i].pressed + '</br>';
            for (let i in _gp.axes)
                s += 'x: ' + i + ' - ' + _gp.axes[i] + '</br>';
            let txt = document.getElementById('debug');
            if (!txt) {
                txt = document.createElement('p');
                txt.id = 'debug';
                txt.style.position = 'absolute';
                txt.style.zIndex = 1000;
                document.body.appendChild(txt);
            }
            txt.innerHTML = s;
        }
    };
}()
