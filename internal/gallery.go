package internal

import (
	"fmt"
	"io"
	"log"
	"os"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

type galleryMode string

const (
	gallery  galleryMode = "gallery"
	enlarged galleryMode = "enlarged"
)

type model struct {
	paths         []string
	idx           int
	mode          galleryMode
	list          list.Model
	width         int
	height        int
	images        map[string]Image
	enlargedImage Image
	tmux          bool
}

type imageItem struct {
	path       string
	size       int64
	resolution string
}

func (i imageItem) Title() string       { return imageName(i.path) }
func (i imageItem) Description() string { return fmt.Sprintf("(%d kb)\n%s", i.size, i.resolution) }
func (i imageItem) FilterValue() string { return imageName(i.path) }

type imageDelegate struct {
	images       map[string]Image
	normalName   lipgloss.Style
	selectedName lipgloss.Style
	placeholder  lipgloss.Style
	elem         lipgloss.Style
	selectedElem lipgloss.Style
}

func newImageDelegate(images map[string]Image) imageDelegate {
	return imageDelegate{
		images: images,
		normalName: lipgloss.NewStyle().
			AlignVertical(lipgloss.Center).
			Foreground(lipgloss.White).
			Bold(true).
			PaddingLeft(2),

		selectedName: lipgloss.NewStyle().
			// Foreground(lipgloss.Color("#EE6FF8")).
			Bold(true).
			PaddingLeft(2),
		placeholder: lipgloss.NewStyle().
			Width(maxWidth).
			Height(maxHeight),
		elem: lipgloss.NewStyle().
			// Border(lipgloss.NormalBorder()).
			Padding(0, 1),
		selectedElem: lipgloss.NewStyle().
			// Border(lipgloss.NormalBorder()).
			Background(lipgloss.Color("#3e3385")).
			Padding(0, 1),
	}
}

func (d imageDelegate) Height() int { return maxHeight }

func (d imageDelegate) Spacing() int { return 1 }

func (d imageDelegate) Update(tea.Msg, *list.Model) tea.Cmd { return nil }

func (d imageDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	img, ok := item.(imageItem)
	if !ok {
		return
	}

	preview := d.placeholder.Render("")
	if im, ok := d.images[img.path]; ok && im.View != "" {
		preview = im.View
	}

	name := imageName(img.path)
	style := d.normalName
	selected_elem := d.elem
	if index == m.Index() && m.FilterState() != list.Filtering {
		name = imageName(img.path)
		// style = d.selectedName
		selected_elem = d.selectedElem
	}
	width, _, _ := term.GetSize(os.Stdout.Fd())

	fmt.Fprint(w, selected_elem.Width(width).AlignVertical(lipgloss.Center).Render(lipgloss.JoinHorizontal(lipgloss.Top, preview, style.Render(lipgloss.JoinVertical(lipgloss.Left, name, img.Description())))))
}

func InitModel(paths []string) *model {
	w, h, _ := term.GetSize(os.Stdout.Fd())
	images := make(map[string]Image)
	items := make([]list.Item, len(paths))
	for i, path := range paths {
		items[i] = imageItem{path: path, size: 0, resolution: ""}
	}
	delegate := newImageDelegate(images)
	imageList := list.New(items, delegate, w, h)
	imageList.SetShowTitle(false)
	imageList.SetStatusBarItemName("image", "images")
	return &model{
		paths:  paths,
		idx:    0,
		mode:   gallery,
		list:   imageList,
		width:  w,
		height: h,
		images: images,
		tmux:   os.Getenv("TMUX") != "",
	}
}

const (
	maxWidth  = 9
	maxHeight = 4
)

func (m *model) Init() tea.Cmd {
	return m.loadGalleryImages()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.mode == enlarged {
				m.mode = gallery
			}
			return m, nil
		case "right":
			if m.mode != gallery {
				if m.idx < len(m.paths)-1 {
					m.idx++
					return m, m.loadImage(m.idx)
				}
			}
		case "left":
			if m.mode != gallery {
				if m.idx > 0 {
					m.idx--
					return m, m.loadImage(m.idx)
				}
			}
		case "space":
			if m.mode == gallery {
				item, ok := m.list.SelectedItem().(imageItem)
				if !ok {
					return m, nil
				}
				m.idx = m.indexForPath(item.path)
				m.mode = enlarged
				return m, m.loadEnlargedImage()
			} else {
				m.mode = gallery
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height)
	}

	if m.mode == gallery {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *model) indexForPath(path string) int {
	for i, p := range m.paths {
		if p == path {
			return i
		}
	}
	return 0
}

func (m *model) loadGalleryImages() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.paths))
	for i := range m.paths {
		cmd := m.loadImage(i)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (m *model) View() tea.View {
	img := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	if m.mode == gallery {
		v := tea.View{Content: m.list.View()}
		v.AltScreen = true
		return v
	}

	return RenderEnlargedView(m, img)
}

func RenderEnlargedView(m *model, img lipgloss.Style) tea.View {
	// imgData := m.images[m.paths[m.idx]]
	// path := m.paths[m.idx]
	// im, err := ShowImage(path, 10, 30, m.idx+1000, m.tmux)
	// if err != nil {
	// log.Fatal(err)
	// }
	// m.enlargedImage = im
	v := tea.View{Content: imageName(m.enlargedImage.Path) + "\n" + img.Render(m.enlargedImage.View)}
	v.AltScreen = true
	return v
}

func (m *model) loadEnlargedImage() tea.Cmd {
	path := m.paths[m.idx]
	im, ok := m.images[path]
	if !ok {
		var err error
		im, err = ShowImage(path, 10, 30, m.idx+1000, m.tmux)
		if err != nil {
			log.Fatal(err)
		}
		m.images[path] = im
	}
	m.enlargedImage = im
	return tea.Raw(im.Raw)
}

func (m *model) loadImage(idx int) tea.Cmd {
	path := m.paths[idx]
	im, ok := m.images[path]
	if !ok {
		var err error
		im, err = ShowImage(path, maxHeight, maxWidth, idx+1, m.tmux)
		if err != nil {
			log.Fatal(err)
		}
		m.images[path] = im
		size := im.Size
		if size > 1000 {
			size = size / 1000
		}
		m.list.SetItem(idx, imageItem{path: path, size: size, resolution: im.Resolution})
	}
	if im.Raw == "" {
		return nil
	}
	return tea.Raw(im.Raw)
}

func Initialize(paths []string) {
	p := tea.NewProgram(InitModel(paths))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
