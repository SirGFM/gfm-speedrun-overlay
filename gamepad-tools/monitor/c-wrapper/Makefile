CUSTOM_CFLAGS := $(CFLAGS)
#CUSTOM_CFLAGS := -fsanitize=address -g -O1 $(CUSTOM_CFLAGS)
#CUSTOM_CFLAGS := -DDEBUG $(CUSTOM_CFLAGS)
CUSTOM_LDFLAGS := $(LDFLAGS) -shared
CUSTOM_LDLIBS := $(LDLIBS) -lSDL2

libsdl-gamepad.so: sdl-gamepad.c
	$(CC) $(CUSTOM_CFLAGS) -fPIC $(CUSTOM_LDFLAGS) $^ -o $@ $(CUSTOM_LDLIBS)

sdl-gamepad.dll: sdl-gamepad.c
	$(CC) $(CUSTOM_CFLAGS) $(CUSTOM_LDFLAGS) -Wl,-soname,$(TARGET).dll -Wl,-export-all-symbols $^ -o $@ $(CUSTOM_LDLIBS)
