# BCoding Go Library (WIP)

This package is meant for encoding/decoding .torrent files according to the BitTorrent Protocol (https://www.bittorrent.org/beps/bep_0000.html)

## Prerequisites

- Golang 1.27.3+

## How to use (on Linux)

Clone this repository

Get some .torrent file, for instance Ubuntu's core server image at https://releases.ubuntu.com/26.04/ubuntu-26.04.1-live-server-amd64.iso.torrent

Put it in the projet directory, then change `var filename`'s value in `bcoding_test.go` to the .torrent filename.

Then run these commands:

```bash
cd go-bcoding
go test -v
```

You know the run passed if you see something like this:

```
--- PASS: TestBdecode (0.05s)$
PASS$
ok      bcoding 0.053s$
```
