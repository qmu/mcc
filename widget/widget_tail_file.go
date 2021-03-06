package widget

import (
	"os"

	ui "github.com/gizak/termui"
	"github.com/hpcloud/tail"
	"github.com/qmu/mcc/widget/listable"
	// "github.com/k0kubun/pp"
)

// TailFileWidget is a command launcher
type TailFileWidget struct {
	options  *Option
	renderer *listable.ListWrapper
	isReady  bool
	disabled bool
	path     string
}

// NewTailFileWidget constructs a New TailFileWidget
func NewTailFileWidget(opt *Option) (n *TailFileWidget, err error) {
	n = new(TailFileWidget)
	n.options = opt
	return
}

// Init is the implementation of stack.Init
func (n *TailFileWidget) Init() (err error) {
	if n.options.Path[0:1] == "/" {
		n.path = n.options.Path
	} else {
		n.path = "./" + n.options.ExecPath + "/" + n.options.Path
	}
	lopt := &listable.ListWrapperOption{
		Title:      n.options.GetTitle(),
		RealHeight: n.options.GetHeight(),
	}
	n.renderer = listable.NewListWrapper(lopt)
	n.isReady = true
	n.tail()
	return
}

func (n *TailFileWidget) tail() (err error) {
	go n.tailActually()
	return
}

func (n *TailFileWidget) tailActually() (err error) {
	f, _ := os.Open(n.path)
	defer f.Close()
	// check if the file exists
	var fi os.FileInfo
	if fi, err = f.Stat(); err != nil {
		n.renderer.AddBody(n.path + " does not exists")
		n.renderer.Activate()
		return
	}
	// check the file size, if it's above 3KB
	// use Location option
	var loc tail.SeekInfo
	if fi.Size() > 3000 {
		loc = tail.SeekInfo{
			Offset: -500,
			Whence: 2,
		}
	}

	t, err := tail.TailFile(n.path, tail.Config{
		Location: &loc,
		Follow:   true,
		Logger:   tail.DiscardingLogger,
	})
	if err != nil {
		return
	}
	first := true
	for line := range t.Lines {
		// if the Location option is enable
		// cut the first line
		if fi.Size() > 3000 && first == true {
			first = false
			continue
		}

		l := line.Text
		n.renderer.AddBody(" " + l)
		n.renderer.MoveCursor("bottom")
	}
	return
}

// Activate is the implementation of Widget.Activate
func (n *TailFileWidget) Activate() {
	n.renderer.Activate()
}

// Deactivate is the implementation of Widget.Activate
func (n *TailFileWidget) Deactivate() {
	n.renderer.Deactivate()
}

// IsDisabled is the implementation of Widget.IsDisabled
func (n *TailFileWidget) IsDisabled() bool {
	return n.disabled
}

// IsReady is the implementation of Widget.IsReady
func (n *TailFileWidget) IsReady() bool {
	return n.isReady
}

// GetHighlightenPos is the implementation of Widget.GetHighlightenPos
func (n *TailFileWidget) GetHighlightenPos() int {
	return n.renderer.GetCursor()
}

// GetGridBufferers is the implementation of stack.Activate
func (n *TailFileWidget) GetGridBufferers() []ui.GridBufferer {
	return []ui.GridBufferer{n.renderer.GetWidget()}
}

// GetWidth is the implementation of stack.Init
func (n *TailFileWidget) GetWidth() int {
	return n.renderer.GetWidth()
}

// GetHeight is the implementation of stack.Init
func (n *TailFileWidget) GetHeight() int {
	return n.renderer.GetHeight()
}

// Disable is
func (n *TailFileWidget) Disable() {
}

// SetOption is
func (n *TailFileWidget) SetOption(opt *AdditionalWidgetOption) {
}
