package dashboard

import (
	"io/ioutil"
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
	m2s "github.com/mitchellh/mapstructure"
	// "github.com/k0kubun/pp"
)

// NoteWidget is a command launcher
type NoteWidget struct {
	renderer *ListWrapper
	isReady  bool
	disabled bool
}

// NewNoteWidget constructs a New NoteWidget
func NewNoteWidget(wi Widget, execPath string) (n *NoteWidget, err error) {
	n = new(NoteWidget)

	var note string
	// TextFile Widget
	if wi.Type == "text_file" {
		b, err := ioutil.ReadFile("./" + execPath + "/" + wi.Path)
		if err != nil {
			return n, err
		}
		note = string(b)
		// Note Widget
	} else {
		if err := m2s.Decode(wi.Content, &note); err != nil {
			return nil, err
		}
	}

	items := strings.Split(note, "\n")
	var body []string
	for _, item := range items {
		rep := regexp.MustCompile(`(^#.*|^--*)`)
		item = rep.ReplaceAllString(item, "[$1](fg-blue)")
		body = append(body, " "+item)
	}
	opt := &ListWrapperOption{
		Title:      wi.Title,
		RealHeight: wi.RealHeight,
		Body:       body,
	}
	n.renderer = NewListWrapper(opt)
	n.isReady = true

	return
}

// Activate is the implementation of Widget.Activate
func (n *NoteWidget) Activate() {
	n.renderer.Render()
}

// Deactivate is the implementation of Widget.Activate
func (n *NoteWidget) Deactivate() {
	n.renderer.ResetRender()
}

// IsDisabled is the implementation of Widget.IsDisabled
func (n *NoteWidget) IsDisabled() bool {
	return n.disabled
}

// IsReady is the implementation of Widget.IsReady
func (n *NoteWidget) IsReady() bool {
	return n.isReady
}

// GetHighlightenPos is the implementation of Widget.GetHighlightenPos
func (n *NoteWidget) GetHighlightenPos() int {
	return n.renderer.GetCursor()
}

// GetWidget is the implementation of widget.Activate
func (n *NoteWidget) GetWidget() *ui.List {
	return n.renderer.GetWidget()
}
