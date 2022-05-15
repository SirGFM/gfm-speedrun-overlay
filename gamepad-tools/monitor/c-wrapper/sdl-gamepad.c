#define SDL_MAIN_HANDLED
#include <SDL2/SDL.h>
#include <SDL2/SDL_atomic.h>
#include <SDL2/SDL_events.h>
#include <SDL2/SDL_joystick.h>
#include <SDL2/SDL_mutex.h>
#include <SDL2/SDL_thread.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#if defined(_WIN32)
#define  RET(type) __declspec(dllexport) type __stdcall
#else
#define  RET(type) type
#endif

#if defined(DEBUG)
#    include <stdio.h>
#    define LOG(fmt, ...) printf(fmt, ##__VA_ARGS__)
#else
#    define LOG(fmt, ...) do {} while (0)
#endif

#define ASSERT_RET(stmt, err) \
    do { \
        if (!(stmt)) { \
            LOG("Statement \"%s\" failed on line %d!\n", #stmt, __LINE__); \
            return err; \
        } \
    } while (0)

#define ASSERT_RV_GOTO(stmt, err, label) \
    do { \
        if (!(stmt)) { \
            LOG("Statement \"%s\" failed on line %d!\n", #stmt, __LINE__); \
            rv = err; \
            goto label; \
        } \
    } while (0)

/** =========================================================================
 *    Custom types
 *  ========================================================================= */

enum error {
    ERR_OK = 0,
    ERR_WAITING,
    ERR_NOT_STARTED,
    ERR_SDL_START_MUTEX,
    ERR_SDL_START_THREAD,
    ERR_SDL_INIT,
    ERR_SDL_INIT_BG_JOYSTICK,
    ERR_SDL_INIT_NO_SIGNALS,
    ERR_SDL_INIT_EVENT,
    ERR_SDL_INIT_JOYSTICK,
    ERR_SDL_WAIT,
    ERR_LOCK,
    ERR_UNLOCK,
    ERR_WAIT_LOCK,
    ERR_ERR_ADD_JOYSTICK_EVENT_1, /* DEPRECATED */
    ERR_ERR_ADD_JOYSTICK_EVENT_2,
    ERR_EXPAND_LIST,
    ERR_OPEN_JOYSTICK,
    ERR_ALLOC_JOY_NODE,
    ERR_WAIT_UNLOCK,
    ERR_GET_JOYSTICK_ID,
    ERR_ADD_JOYSTICK_EVENT_3,
};

/** Tracks "balls", a 2D trackpad?. */
struct joystick_ball {
    int32_t x, y;
};

/**
 * Accessors for the region used to track a SDL2 Joystick. Every field
 * points to a memory region right after this struct.
 */
struct joystick {
    /** The SDL2 Joystick tracked by this struct. */
    SDL_Joystick *sdlj;
    /** Points to the number of balls in the SDL2 Joystick. */
    uint8_t *num_balls;
    /** Points to the axes of balls in the SDL2 Joystick. */
    uint8_t *num_axes;
    /** Points to the buttons of balls in the SDL2 Joystick. */
    uint8_t *num_buttons;
    /** Points to the hats of balls in the SDL2 Joystick. */
    uint8_t *num_hats;
    /** Points to the length of the SDL2 Joystick's name. */
    uint32_t *name_len;
    /** Points to the length of the SDL2 Joystick's GUID. */
    uint32_t *guid_len;
    /** Points to the list of balls. */
    struct joystick_ball *balls;
    /** Points to the list of axes. */
    Sint16 *axes;
    /** Points to the list of buttons. */
    Uint8 *buttons;
    /** Points to the list of hats. */
    Uint8 *hats;
    /** Points to the joystick name. */
    char *name;
    /** Points to the device's GUID. */
    Uint8 *guid;
};

/** Track the number of fields in `struct joystick`. */
struct joystick_fields_num {
    /** The number of `balls`. */
    size_t balls;
    /** The number of `axes`. */
    size_t axes;
    /** The number of `buttons`. */
    size_t buttons;
    /** The number of `hats`. */
    size_t hats;
    /** The lenght of the name, with the trailing '\0'. */
    size_t name_len;
    /** The lenght of the device's GUID. */
    size_t guid_len;
};

/** =========================================================================
 *    Private, global variable
 *  ========================================================================= */

/** Length, in bytes, of the longest packed gamepad data. */
static size_t biggest_node;

/** ID of the latest gamepad to send an event. */
static int last_id;

/** Number of entries in `joy_list`. Entries are guaranteed to never be NULL. */
static int joy_count;

/** List of gamepads being tracked. */
static struct joystick **joy_list;

/** Thread for monitoring gamepad events. */
static SDL_Thread *thread;

/** Synchronizes access to the global joystick data. */
static SDL_mutex *mutex;

/** Controls whether the thread should keep running. */
static int running;

/** Signals that the BG thread is done initializing, or that some error happened */
static SDL_atomic_t ready;

/** =========================================================================
 *    Internal functions
 *  ========================================================================= */

/**
 * Synchronizes access to the library.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static enum error lock() {
    if (SDL_LockMutex(mutex) != 0)
        return ERR_LOCK;
    return ERR_OK;
}

/**
 * Reverts to unsynchronized access to the library.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static enum error unlock() {
    if (SDL_UnlockMutex(mutex) != 0)
        return ERR_UNLOCK;
    return ERR_OK;
}

/** Static name used when SDL2 fails to retrieve the Joystick's name. */
static const char *unknown_name = "Unknown";

/** Static name used when the Joystick doesn't have an assigned name. */
static const char *empty_name = "Empty";

/**
 * Retrieve the name for a SDL2 Joystick
 *
 * @param[in] sdlj: The SDL2 Joystick being queried.
 *
 * @return Either the Joystick's name or a (static) dummy name.
 */
static const char* get_sdljoy_name(SDL_Joystick *sdlj) {
    const char *name;

    name = SDL_JoystickName(sdlj);
    if (name == 0)
        return unknown_name;
    else if (strlen(name) == 0)
        return empty_name;
    return name;
}

/**
 * Retrieve the number of each field in a SDL2 Joystick.
 *
 * @param[in] sdlj: The SDL2 Joystick.
 *
 * @return A `struct joystick_fields_num` with the retrieved data.
 */
static struct joystick_fields_num get_num_joy_fields(SDL_Joystick *sdlj) {
    struct joystick_fields_num num_fields;

    num_fields.balls = SDL_JoystickNumBalls(sdlj);
    if (num_fields.balls < 0)
        num_fields.balls = 0;
    num_fields.axes = SDL_JoystickNumAxes(sdlj);
    if (num_fields.axes < 0)
        num_fields.axes = 0;
    num_fields.buttons = SDL_JoystickNumButtons(sdlj);
    if (num_fields.buttons < 0)
        num_fields.buttons = 0;
    num_fields.hats = SDL_JoystickNumHats(sdlj);
    if (num_fields.hats < 0)
        num_fields.hats = 0;
    num_fields.name_len = strlen(get_sdljoy_name(sdlj)) + 1;
    num_fields.guid_len = sizeof(SDL_JoystickGUID);

    return num_fields;
}

/**
 * Expand the list of tracked devices so the requested index may be
 * accessed.
 *
 * @param[in] idx: Index of the joystick.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static enum error expand_list(int idx) {
    struct joystick **new_joy_list;
    int new_joy_count = 0;

    if (joy_count > idx)
        return ERR_OK;

    new_joy_count = (idx + 1) * 2;
    new_joy_list = malloc(sizeof(struct joystick*) * new_joy_count);
    ASSERT_RET(new_joy_list, ERR_EXPAND_LIST);

    memset(new_joy_list, 0x0, sizeof(struct joystick*) * new_joy_count);
    memcpy(new_joy_list, joy_list, sizeof(struct joystick*) * joy_count);

    free(joy_list);
    joy_list = new_joy_list;
    joy_count = new_joy_count;

    return ERR_OK;
}

/**
 * Retrieve the number of each field in `node`, a `struct joystick`. `node`
 * may be NULL, in which case 0 will be returned for each field.
 *
 * @param[in] node: The `struct joystick` being coverted.
 *
 * @return The converted `struct joystick_fields_num`.
 */
static struct joystick_fields_num node_to_joy_fields(struct joystick *node) {
    struct joystick_fields_num num_fields;

    if (node) {
        num_fields.balls = *node->num_balls;
        num_fields.axes = *node->num_axes;
        num_fields.buttons = *node->num_buttons;
        num_fields.hats = *node->num_hats;
        num_fields.name_len = *node->name_len;
        num_fields.guid_len = *node->guid_len;
    }
    else {
        memset(&num_fields, 0x0, sizeof(num_fields));
    }

    return num_fields;
}

/**
 * @return The rounded-up size of a `struct joystick` aligned to 8 bytes.
 */
static size_t get_aligned8_joystick_size() {
    return (sizeof(struct joystick) & (~0x7)) + 0x8;
}

/**
 * Calculate the size of a given joystick, alongside all its extra data.
 *
 * @param[in] num_fields: The number of each field in the desired joystick.
 *
 * @return The size of the given joystick.
 */
static size_t get_joystick_data_size(struct joystick_fields_num num_fields) {
    size_t size;

    size = get_aligned8_joystick_size();
    size += sizeof(uint8_t) + num_fields.balls * sizeof(struct joystick_ball);
    size += sizeof(uint8_t) + num_fields.axes * sizeof(uint16_t);
    size += sizeof(uint8_t) + num_fields.buttons * sizeof(uint8_t);
    size += sizeof(uint8_t) + num_fields.hats * sizeof(uint8_t);
    size += sizeof(uint32_t) + num_fields.name_len;
    size += sizeof(uint32_t) + num_fields.guid_len;

    return size;
}

/**
 * Retrieve an address and advance it to the next field. Should be called
 * from the helper macro `get_ptr_and_move`.
 *
 * @param[in/out] src:        in:  The base address;
 *                            out: The address after the region.
 * @param[in]     field_size: The size of the given region.
 *
 * @return The starting address
 */
static void* _get_ptr_and_move(uintptr_t *src, size_t field_size) {
    uintptr_t dst = *src;
    *src += (uintptr_t)field_size;
    return (void*)dst;
}
/**
 * Retrieve an address of a given type, advancing it by a number of
 * elements of the same type.
 *
 * @param[in/out] src:  in:  The base address;
 *                      out: The address after the region.
 * @param[in]    type:  The type of the field.
 * @param[in]    nmemb: The number of members of the given type.
 *
 * @return An address cast to `type*`.
 */
#define get_ptr_and_move(src, type, nmemb) \
    (type*)_get_ptr_and_move(src, sizeof(type) * (nmemb))

/**
 * Configure a new joystick representing a SDL2 Joystick. This should be
 * called in response to a SDL_JOYDEVICEADDED event.
 *
 * @param[in]  idx: Index of the new SDL2 joystick.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static int new_joy_node(int idx) {
    struct joystick *node;
    SDL_Joystick *sdlj;
    struct joystick_fields_num num_fields;
    SDL_JoystickID jid;
    SDL_JoystickGUID guid;
    size_t size;
    uintptr_t ptr;
    int rv = 0;

    /* SDL_JOYDEVICEADDED indicates the index of the connected joystick.
     * However, further events use the joystick instance ID to identify the
     * joystick that generated the event.
     *
     * Therefore, not only the joystick ID must be used to index the data
     * in the list of joysticks, but last_id must also be updated
     * accordingly. */
    sdlj = SDL_JoystickOpen(idx);
    ASSERT_RV_GOTO(sdlj != 0, ERR_OPEN_JOYSTICK, open_joystick);

    jid = SDL_JoystickInstanceID(sdlj);
    ASSERT_RV_GOTO(jid >= 0, ERR_GET_JOYSTICK_ID, close_joystick);
    rv = expand_list(jid);
    ASSERT_RV_GOTO(rv == ERR_OK, ERR_ADD_JOYSTICK_EVENT_3, close_joystick);

    last_id = jid;

    num_fields = get_num_joy_fields(sdlj);
    size = get_joystick_data_size(num_fields);

    node = malloc(size);
    ASSERT_RV_GOTO(node != 0, ERR_ALLOC_JOY_NODE, close_joystick);
    memset(node, 0x0, size);

    node->sdlj = sdlj;
    ptr = (uintptr_t)node + (uintptr_t)get_aligned8_joystick_size();

    node->num_balls = get_ptr_and_move(&ptr, uint8_t, 1);
    node->num_axes = get_ptr_and_move(&ptr, uint8_t, 1);
    node->num_buttons = get_ptr_and_move(&ptr, uint8_t, 1);
    node->num_hats = get_ptr_and_move(&ptr, uint8_t, 1);
    node->name_len = get_ptr_and_move(&ptr, uint32_t, 1);
    node->guid_len = get_ptr_and_move(&ptr, uint32_t, 1);
    node->balls = get_ptr_and_move(&ptr, struct joystick_ball, num_fields.balls);
    node->axes = get_ptr_and_move(&ptr, uint16_t, num_fields.axes);
    node->buttons = get_ptr_and_move(&ptr, uint8_t, num_fields.buttons);
    node->hats = get_ptr_and_move(&ptr, uint8_t, num_fields.hats);
    node->name = get_ptr_and_move(&ptr, char, num_fields.name_len);
    node->guid = get_ptr_and_move(&ptr, char, num_fields.guid_len);

    *node->num_balls = num_fields.balls;
    *node->num_axes = num_fields.axes;
    *node->num_buttons = num_fields.buttons;
    *node->num_hats = num_fields.hats;
    *node->name_len = num_fields.name_len;
    *node->guid_len = num_fields.guid_len;

    memcpy(node->name, get_sdljoy_name(sdlj), num_fields.name_len);
    guid = SDL_JoystickGetGUID(sdlj);
    memcpy(node->guid, &guid, num_fields.guid_len);

    size -= get_aligned8_joystick_size();
    if (size > biggest_node)
        biggest_node = size;

    joy_list[jid] = node;
    LOG("Added joystick %d\n", jid);
    return ERR_OK;

close_joystick:
    SDL_JoystickClose(sdlj);
open_joystick:
    return rv;
}

/**
 * Retrieve a joystick from the list of devices being tracked. Returns NULL
 * if the device hasn't been configured yet.
 *
 * @param[in] idx: Index of the joystick.
 *
 * @return A pointer to the joystick.
 */
static struct joystick *get_node(int idx) {
    if (idx < joy_count)
        return joy_list[idx];
    return 0;
}

/**
 * Release a joystick from the list of devices being tracked.
 *
 * @param[in] idx: Index of the joystick.
 */
static void rem_joy_node(int idx) {
    if (idx > joy_count)
        return;
    free(get_node(idx));
    joy_list[idx] = 0;
}

/**
 * Track a SDL2 Joystick's axis from its event.
 *
 * @param[in] axis: The SDL2 event.
 */
static void set_axis(SDL_JoyAxisEvent *axis) {
    struct joystick *node = get_node(axis->which);
    if (node == 0)
        return;

    node->axes[axis->axis] = axis->value;
}

/**
 * Track a SDL2 Joystick's ball from its event.
 *
 * @param[in] ball: The SDL2 event.
 */
static void set_ball(SDL_JoyBallEvent *ball) {
    struct joystick *node = get_node(ball->which);
    if (node == 0)
        return;

    node->balls[ball->ball].x = (int32_t)ball->xrel;
    node->balls[ball->ball].y = (int32_t)ball->yrel;
}

/**
 * Track a SDL2 Joystick's button from its event.
 *
 * @param[in] button: The SDL2 event.
 * @param[in] state:  State of the button: released (0) or pressed (other).
 */
static void set_button(SDL_JoyButtonEvent *bt, int state) {
    struct joystick *node = get_node(bt->which);
    if (node == 0)
        return;

    node->buttons[bt->button] = state;
}

/**
 * Track a SDL2 Joystick's hat from its event.
 *
 * @param[in] hat: The SDL2 event.
 */
static void set_hat(SDL_JoyHatEvent *hat) {
    struct joystick *node = get_node(hat->which);
    if (node == 0)
        return;

    node->hats[hat->hat] = hat->value;
}

/**
 * Handle a SDL2 Joystick event. `ev` should have been previously filtered
 * by calling `is_joy_event()`.
 *
 * @param[in] ev: The SDL2 event.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static enum error handle_event(SDL_Event *ev) {
    switch (ev->type) {
    case SDL_JOYDEVICEADDED:
        ASSERT_RET(new_joy_node(ev->jdevice.which) == ERR_OK, ERR_ERR_ADD_JOYSTICK_EVENT_2);
        LOG("Added new joystick (%d)!\n", ev->jdevice.which);
        break;
    case SDL_JOYDEVICEREMOVED:
        rem_joy_node(ev->jdevice.which);
        LOG("Removed joystick (%d)!\n", ev->jdevice.which);
        break;
    case SDL_JOYAXISMOTION:
        set_axis(&ev->jaxis);
        break;
    case SDL_JOYBALLMOTION:
        set_ball(&ev->jball);
        break;
    case SDL_JOYBUTTONDOWN:
        set_button(&ev->jbutton, 1);
        break;
    case SDL_JOYBUTTONUP:
        set_button(&ev->jbutton, 0);
        break;
    case SDL_JOYHATMOTION:
        set_hat(&ev->jhat);
        break;
    }

    return ERR_OK;
}

/**
 * Check whether a given SDL2 event is for a Joystick.
 *
 * @param[in] ev: The SDL2 event.
 *
 * @return 1 for joystick events, 0 otherwise.
 */
static inline int is_joy_event(SDL_Event *ev) {
    switch (ev->type) {
    case SDL_JOYDEVICEADDED:
    case SDL_JOYDEVICEREMOVED:
    case SDL_JOYAXISMOTION:
    case SDL_JOYBALLMOTION:
    case SDL_JOYBUTTONDOWN:
    case SDL_JOYBUTTONUP:
    case SDL_JOYHATMOTION:
        return 1;
    default:
        return 0;
    }
}

/**
 * Wait and handle every queued event. It returns after handling every
 * queued event, so this should be called in a loop.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static enum error mainloop() {
    SDL_Event ev;
    const int timeout = 500; /* ms */
    int rv;
    int locked = 0;
    unsigned int time;

    ASSERT_RET(SDL_WaitEventTimeout(&ev, timeout) == 1, ERR_SDL_WAIT);

    time = SDL_GetTicks();
    do {
        LOG("Got event of type: %d\n", ev.type);
        if (is_joy_event(&ev)) {
            if (!locked) {
                /* Lock only on the first joystick event */
                ASSERT_RET(lock() == ERR_OK, ERR_WAIT_LOCK);
                locked = 1;
            }
            last_id = ev.jdevice.which;
            rv = handle_event(&ev);
        }
    } while (SDL_PollEvent(&ev) && SDL_GetTicks() - time < timeout);

    if (locked)
        ASSERT_RET(unlock() == ERR_OK, ERR_WAIT_UNLOCK);

    return rv;
}

/**
 * Initialize SDL2 and handle events until signaled to stop.
 *
 * @return An `enum error` indicating the result of the operation.
 */
static int runner(void *arg) {
    int rv;

    ASSERT_RV_GOTO(SDL_Init(0) == ERR_OK, ERR_SDL_INIT, no_sdl);
    ASSERT_RV_GOTO(SDL_SetHint("SDL_JOYSTICK_ALLOW_BACKGROUND_EVENTS", "1") == SDL_TRUE, ERR_SDL_INIT_BG_JOYSTICK, sdl_init);
    ASSERT_RV_GOTO(SDL_SetHint("SDL_HINT_NO_SIGNAL_HANDLERS", "0") == SDL_TRUE, ERR_SDL_INIT_NO_SIGNALS, sdl_init);
    ASSERT_RV_GOTO(SDL_InitSubSystem(SDL_INIT_EVENTS) == ERR_OK, ERR_SDL_INIT_EVENT, sdl_init);
    ASSERT_RV_GOTO(SDL_InitSubSystem(SDL_INIT_JOYSTICK) == ERR_OK, ERR_SDL_INIT_JOYSTICK, sdl_init);

    SDL_AtomicSet(&ready, ERR_OK);

    while (running) {
        rv = mainloop();
    }

    for (int i = 0; i < joy_count; i++) {
        free(joy_list[i]);
    }
    free(joy_list);

    rv = ERR_OK;
sdl_init:
    SDL_Quit();
no_sdl:
    SDL_AtomicSet(&ready, rv);
    return rv;
}

/** =========================================================================
 *    SDL2-based gamepad monitor - Public API
 *  -------------------------------------------------------------------------
 *     This gamepad monitor runs in a background thread, detecting and
 *   recording gamepad events.
 *
 *     Call `start()` to initialize the monitor, and call 'clean()' when you are
 *   done with it, so it may stop running and its resources may be released.
 *   Before requesting gamepad data, be sure to wait until `check_ready()`
 *   returns `ERR_OK`.
 *
 *     Each gamepad data is in packed into an array of uint8_t, structured as
 *   follows:
 *
 *    struct ball {
 *        int32_t x;
 *        int32_t y;
 *    }
 *
 *    struct gamepad_data {
 *        uint8_t num_balls;
 *        uint8_t num_axes;
 *        uint8_t num_buttons;
 *        uint8_t num_hats;
 *        uint32_t name_len;
 *        uint32_t guid_len;
 *        struct ball balls[];
 *        int16_t axes[];
 *        uint8_t buttons[];
 *        uint8_t hats[];
 *        char *name;
 *        uint8_t guid[];
 *    }
 *
 *     The number of each field (`num_*`) is required, but the arrays are only
 *   present if the gamepad has at least one of that kind. The fields were
 *   sorted in such a way that no matter if a field is missing, the rest of
 *   the struct should be properly aligned (e.g., uint32_t are always 4 bytes
 *   aligned).
 *
 *     While running, all requests must be made between a matching pair of
 *   `begin_request()` and `end_request()` calls. These functions ensure the
 *   monitor won't modify any data while you are reading a gamepad's state.
 *   The buffer used to retrieve a gamepad data must be at least
 *   `get_node_size()` bytes long, and `get_last_id()` may be used to retrieve
 *   the ID of the last gamepad that sent any event. To check if a given
 *   gamepad has valid data, check whether or not its name is empty. Finally,
 *   to retrieve a gamepad state call `get_data()`.
 *  ========================================================================= */

/**
 * Start running the SDL2-based gamepad monitor.
 *
 * @return An `enum error` indicating the result of the operation.
 */
RET(int32_t) start() {
    SDL_AtomicSet(&ready, ERR_NOT_STARTED);

    /* Note: These may safely be called before SDL_Init() */
    mutex = SDL_CreateMutex();
    ASSERT_RET(mutex != 0, ERR_SDL_START_MUTEX);

    running = 1;
    thread = SDL_CreateThread(runner, "gamepad-overlay-main", 0);
    ASSERT_RET(thread != 0, ERR_SDL_START_THREAD);

    SDL_AtomicSet(&ready, ERR_WAITING);
    return ERR_OK;
}

/**
 * Check whether the monitor has finished starting.
 *
 * @return An `enum error` indicating the state of the monitor:
 *           - ERR_WAITING: The monitor is still starting without errors
 *           - ERR_OK: The monitor has started and is running in background
 *           - Anything else: The error detected while starting the monitor
 */
RET(int32_t) check_ready() {
    return SDL_AtomicGet(&ready);
}

/**
 * Stop the monitor and clean its resources.
 */
RET(void) clean() {
    if (thread != 0) {
        running = 0;
        SDL_WaitThread(thread, 0);
    }
    thread = 0;

    if (mutex != 0) {
        SDL_DestroyMutex(mutex);
    }
    mutex = 0;
}

/**
 * Block the monitor, so you may retrieve gamepad data.
 *
 * @return An `enum error` indicating the result of the operation.
 */
RET(int32_t) begin_request() {
    return (int32_t)lock();
}

/**
 * Unblock the monitor.
 *
 * @return An `enum error` indicating the result of the operation.
 */
RET(int32_t) end_request() {
    return (int32_t)unlock();
}

/**
 * Retrieve the greatest buffer length possibly needed by `get_data()`.
 *
 * Must be called within a block enclosed by a `begin_request()` and a
 * `end_request()` call.
 *
 * @return The expected buffer length.
 */
RET(uint32_t) get_node_size() {
    return (uint32_t)biggest_node;
}

/**
 * Retrieve the ID of the most recently updated gamepad, which may be
 * supplied to `get_data()`.
 *
 * Must be called within a block enclosed by a `begin_request()` and a
 * `end_request()` call.
 *
 * @return The gamepad ID.
 */
RET(uint32_t) get_last_id() {
    return (uint32_t)last_id;
}

/**
 * Retrieve the an upper bound for the number of gamepads being tracked.
 * This value is guaranteed to remain unchanged until the `end_request()`
 * call.
 *
 * Must be called within a block enclosed by a `begin_request()` and a
 * `end_request()` call.
 *
 * @return The maximum number of currently tracked gamepads.
 */
RET(uint32_t) get_num_gamepads() {
    return (uint32_t)joy_count;
}

/**
 * Retrieve a gamepad's state. Check the API description for instructions
 * to unpack this data.
 *
 * Must be called within a block enclosed by a `begin_request()` and a
 * `end_request()` call.
 *
 * @param[in]  idx:  The ID of the gamepad.
 * @param[out] data: The buffer to be filled with the gamepad data.
 */
RET(void) get_data(uint32_t idx, uint8_t *data) {
    struct joystick *node;
    struct joystick_fields_num num_fields;
    size_t size;

    node = get_node((int)idx);
    num_fields = node_to_joy_fields(node);
    size = get_joystick_data_size(num_fields);
    size -= get_aligned8_joystick_size();

    if (node) {
        uint8_t *joy_data;
        uintptr_t ptr;

        ptr = (uintptr_t)node + (uintptr_t)get_aligned8_joystick_size();
        joy_data = (uint8_t*)ptr;

        memcpy(data, joy_data, size);
    }
    else {
        memset(data, 0x0, size);
    }
}
