# ls-gallery

A small terminal image gallery written in Go.

It scans a directory for image files and opens a keyboard-driven gallery view in the terminal.

## Features

- Shows images from a directory
- Supports `.jpg`, `.jpeg`, `.png`, and `.gif`
- Toggle between gallery mode and enlarged image mode
- Arrow key navigation in enlarged mode
- Basic Kitty graphics support with tmux passthrough handling

## Requirements

- Go 1.25+
- A terminal with Kitty graphics protocol support

## Install

```bash
go install github.com/Kartik-2239/ls-gallery/cmd@latest
```

## Build

```bash
go build -o ls-gallery ./cmd
```

## Run

To scan the current directory:

```bash
./ls-gallery
```

To scan a specific folder:

```bash
./ls-gallery -path /path/to/images
```

You can also run it directly without building:

```bash
go run ./cmd -path /path/to/images
```

## Usage

- `left`: previous image in enlarged mode
- `right`: next image in enlarged mode
- `esc`: return to gallery
- `q` or `ctrl+c`: quit
