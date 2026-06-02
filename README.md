# ls-gallery

A small terminal image gallery written in Go.

<div align="center">
<video src="https://github.com/user-attachments/assets/2e6e20d0-bbae-4cf0-8f4f-2c3372ecc29e" controls muted playsinline width="100%"></video>
</div>

## Requirements

- Go 1.25+
- A modern terminal with Kitty graphics protocol support like Kitty, Ghostty or Wezterm

## Install

```bash
go install github.com/Kartik-2239/ls-gallery/cmd/ls-gallery@latest
```

## Build

```bash
go build -o ls-gallery ./cmd
```

## Run

To list the current directory:

```bash
./ls-gallery
```

To list a specific folder:

```bash
./ls-gallery -path /path/to/images
```

Run without building:

```bash
go run ./cmd -path /path/to/images
```
