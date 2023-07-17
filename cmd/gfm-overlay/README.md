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
