package dashboard

import (
	"os"

	ui "github.com/gizak/termui"
	"github.com/hpcloud/tail"
	// "github.com/k0kubun/pp"
)

// TailFileWidget is a command launcher
type TailFileWidget struct {
	renderer *ListWrapper
	isReady  bool
	disabled bool
	path     string
}

// NewTailFileWidget constructs a New TailFileWidget
func NewTailFileWidget(wi Widget, execPath string) (n *TailFileWidget, err error) {
	n = new(TailFileWidget)
	if wi.Path[0:1] == "/" {
		n.path = wi.Path
	} else {
		n.path = "./" + execPath + "/" + wi.Path
	}
	opt := &ListWrapperOption{
		Title:      wi.Title,
		RealHeight: wi.RealHeight,
	}
	n.renderer = NewListWrapper(opt)
	n.isReady = true
	n.tail()

	return
}

func (n *TailFileWidget) tail() (err error) {
	go func() {
		f, _ := os.Open(n.path)
		defer f.Close()
		// check if the file exists
		var fi os.FileInfo
		if fi, err = f.Stat(); err != nil {
			n.renderer.AddBody(n.path + " does not exists")
			n.renderer.Render()
			return
		}
		// check the file size, if it's above 3KB
		// use Location option
		var loc tail.SeekInfo
		if fi.Size() > 3000 {
			loc = tail.SeekInfo{
				Offset: -2500,
				Whence: 2,
			}
		}

		t, err := tail.TailFile(n.path, tail.Config{
			Location: &loc,
			ReOpen:   true,
			Follow:   true,
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
			n.renderer.AddBody(" " + line.Text)
			n.renderer.Render()
			n.renderer.moveAndRender("bottom")
		}
	}()
	return
}

// Activate is the implementation of Widget.Activate
func (n *TailFileWidget) Activate() {
	n.renderer.Render()
}

// Deactivate is the implementation of Widget.Activate
func (n *TailFileWidget) Deactivate() {
	n.renderer.ResetRender()
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

// GetWidget is the implementation of widget.Activate
func (n *TailFileWidget) GetWidget() *ui.List {
	return n.renderer.GetWidget()
}
