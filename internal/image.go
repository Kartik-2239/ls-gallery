package internal

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi/kitty"
)

type Image struct {
	Raw        string
	View       string
	IsVisible  bool
	Path       string
	Size       int64
	Resolution string
}

func ShowImage(path string, maxRows int, maxCols int, imageID int, tmux bool) (Image, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Image{}, err
	}
	i, f, err := image.Decode(bytes.NewReader(data))
	if f != "png" && f != "jpeg" && f != "gif" {
		return Image{}, fmt.Errorf("unsupported image format: %s", f)
	}
	if err != nil {
		return Image{}, err
	}
	if imageID <= 0 {
		imageID = 1
	}
	var transmit strings.Builder
	bounds := i.Bounds()
	options := &kitty.Options{
		Action:           kitty.TransmitAndPut,
		Format:           kitty.PNG,
		ID:               imageID,
		ImageWidth:       bounds.Dx(),
		ImageHeight:      bounds.Dy(),
		Transmission:     kitty.Direct,
		Chunk:            true,
		Quite:            2,
		VirtualPlacement: true,
	}
	if tmux {
		options.ChunkFormatter = tmuxPassthrough
	}
	if maxCols > 0 {
		options.Columns = maxCols
	}
	if maxRows > 0 {
		options.Rows = maxRows
	}
	err = kitty.EncodeGraphics(&transmit, i, options)
	if err != nil {
		return Image{}, err
	}
	return Image{
		Raw:        transmit.String(),
		View:       kittyPlaceholder(imageID, maxRows, maxCols),
		IsVisible:  true,
		Path:       path,
		Size:       int64(len(data)),
		Resolution: fmt.Sprintf("%d x %d", bounds.Dx(), bounds.Dy()),
	}, nil
}

func tmuxPassthrough(seq string) string {
	return "\x1bPtmux;" + strings.ReplaceAll(seq, "\x1b", "\x1b\x1b") + "\x1b\\"
}

func kittyPlaceholder(imageID int, rows int, cols int) string {
	if rows <= 0 {
		rows = 1
	}
	if cols <= 0 {
		cols = 1
	}

	var b strings.Builder
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			b.WriteString(fmt.Sprintf("\x1b[38;5;%dm%c", imageID%256, kitty.Placeholder))
			b.WriteRune(kitty.Diacritic(y))
			b.WriteRune(kitty.Diacritic(x))
		}
		if y < rows-1 {
			b.WriteByte('\n')
		}
	}
	b.WriteString("\x1b[39m")
	return b.String()
}
