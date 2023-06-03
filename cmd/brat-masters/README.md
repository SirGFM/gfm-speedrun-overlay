# Brat Master Overlay

## Building:

To build this overlay for Windows, run:

```bash
GOOS=windows GOARCH=amd64 go build .
```

### Changing icons/meta-data

To change icons/metadata, make changes to the configuration file ([winres/winres.json](winres/winres.json)) then run:

```bash
go run github.com/tc-hib/go-winres make
```
