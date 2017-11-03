package widget

import (
	"io/ioutil"
	"regexp"
	"strings"

	ui "github.com/gizak/termui"
	m2s "github.com/mitchellh/mapstructure"
	"github.com/qmu/mcc/widget/listable"
	// "github.com/k0kubun/pp"
)

// NoteWidget is a command launcher
type NoteWidget struct {
	options  *Option
	renderer *listable.ListWrapper
	isReady  bool
	disabled bool
}

// NewNoteWidget constructs a New NoteWidget
func NewNoteWidget(opt *Option) (n *NoteWidget, err error) {
	n = new(NoteWidget)
	n.options = opt
	var note string
	if n.options.Type == "text_file" {
		// for TextFile Widget
		var path string
		if n.options.Path[0:1] == "/" {
			path = n.options.Path
		} else {
			path = "./" + n.options.ExecPath + "/" + n.options.Path
		}
		b, err := ioutil.ReadFile(path)

		if err != nil {
			return n, err
		}
		note = string(b)
	} else {
		// for Note Widget
		if err := m2s.Decode(n.options.Content, &note); err != nil {
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
	lopt := &listable.ListWrapperOption{
		Title:      n.options.GetTitle(),
		RealHeight: n.options.GetHeight(),
		Body:       body,
	}
	n.renderer = listable.NewListWrapper(lopt)
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

// GetGridBufferers is the implementation of stack.Activate
func (n *NoteWidget) GetGridBufferers() []ui.GridBufferer {
	return []ui.GridBufferer{n.renderer.GetWidget()}
}

// Render is the implementation of stack.Render
func (n *NoteWidget) Render() (err error) {
	return
}

// GetWidth is the implementation of stack.Render
func (n *NoteWidget) GetWidth() int {
	return n.renderer.GetWidth()
}

// GetHeight is the implementation of stack.Render
func (n *NoteWidget) GetHeight() int {
	return n.renderer.GetHeight()
}

// Disable is
func (n *NoteWidget) Disable() {
}

// SetOption is
func (n *NoteWidget) SetOption(opt *AdditionalWidgetOption) {
}
