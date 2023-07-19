# GFM's stream overlay backend

## Quick guide

First, build the latest frontend:

```sh
git clone git@github.com:SirGFM/react-overlay.git
cd ./react-overlay
npm run build
```

Then, copy the generated files into the backend:

```sh
cp -r ./react-overlay/build ./res
```

Finally, build the backend:

```sh
go build .
```

## Hotkeys

Hotkeys may be configured with a `config.ini` file.
Each `[segment]` in this file must define a hotkey,
with the name of the segment being the key and accepting the attributes:

- `run`: the action when a category was configured
- `timer`: the action used when there's no category

Also, the special segment `[config]` accepts the attribute `pool-rate`,
which determines how many times the keyboard is checked per second.

The config file must be specified on the argument `hotkey-config`,
and the list of keys may be obtained with the argument `print-keys`.

### Example

The config file bellow is mapped to the following actions:

```ini
[config]
	pool-rate = 10

[esc]
	timer = reset
	run = reset

[space]
	timer = stop
	run = split

[backspace]
	run = undo

[return]
	timer = start
	run = start

[s]
	run = save
```

| Key | Timer Action | Category Action |
| -- | -- | -- |
| Enter | Start/Continue | Start |
| Spacebar | Stop | Split the segment |
| Escape | Reset | Reset |
| S | n/a | Save the run |
| Backspace | n/a | Undo the segment |

## Pages

- [config](http://localhost:8080/gfm/config)
- [dashboard](http://localhost:8080/gfm/dashboard)
- [overlay](http://localhost:8080/gfm/stream-layout)
