package dashboard

import ui "github.com/gizak/termui"

// ListWrapper make a List widget which includes
// multi-line texts look like scrolled
type ListWrapper struct {
	widget       *ui.List
	gPressed     bool
	stepsByJump  int // how many steps to jump when C-d, C-u
	listRenderer *ListRenderer
}

// ListWrapperOption is the option argument for NewListWrapper
type ListWrapperOption struct {
	Title         string
	RealHeight    int
	Header        []string
	Body          []string
	LineHighLight bool
}

// NewListWrapper constructs a ListWrapper
func NewListWrapper(opt *ListWrapperOption) (l *ListWrapper) {
	l = new(ListWrapper)

	w := ui.NewList()
	w.Height = opt.RealHeight
	w.Items = []string{"loading..."}
	w.BorderLabel = opt.Title
	l.widget = w

	ropt := &ListRendererOption{
		Header:        opt.Header,
		Body:          opt.Body,
		MaxH:          w.Height,
		LineHighLight: opt.LineHighLight,
	}
	l.listRenderer = NewListRenderer(ropt)
	l.stepsByJump = 10

	return
}

func (l *ListWrapper) setKeyBindings() error {
	// move down cursor by j
	ui.Handle("/sys/kbd/j", func(ui.Event) {
		l.moveAndRender("down")
	})
	// move down cursor by down key
	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		l.moveAndRender("down")
	})
	// move up cursor by k
	ui.Handle("/sys/kbd/k", func(ui.Event) {
		l.moveAndRender("up")
	})
	// move up cursor by up key
	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		l.moveAndRender("up")
	})
	// skip up cursor by Ctrl + u
	ui.Handle("/sys/kbd/C-u", func(ui.Event) {
		for i := 0; i < l.stepsByJump; i++ {
			l.listRenderer.Move("up")
		}
		l.moveAndRender("up")
	})
	// skip down cursor by Ctrl + d
	ui.Handle("/sys/kbd/C-d", func(ui.Event) {
		l.gPressed = false // cancel gg to top
		for i := 0; i < l.stepsByJump; i++ {
			l.listRenderer.Move("down")
		}
		l.moveAndRender("down")
	})
	// cancel pressed g for moving top by gg
	ui.Handle("/sys/kbd", func(e ui.Event) {
		if l.gPressed && e.Path != "/sys/kbd/g" {
			l.gPressed = false
		}
	})
	// moving top by gg
	ui.Handle("/sys/kbd/g", func(ui.Event) {
		if l.gPressed {
			l.listRenderer.Move("top")
			l.moveAndRender("up")
		} else {
			l.gPressed = true
		}
	})
	// moving bottom by G
	ui.Handle("/sys/kbd/G", func(ui.Event) {
		l.listRenderer.Move("bottom")
		l.moveAndRender("down")
	})

	return nil
}

// Render display current *ui.List.Items
func (l *ListWrapper) Render() {
	l.setKeyBindings()
	l.widget.BorderLabelFg = ui.ColorGreen
	l.widget.BorderFg = ui.ColorGreen
	l.widget.Items = l.listRenderer.RenderActually()
	ui.Render(ui.Body)
}

// ResetRender returns a initial multi-line texts
func (l *ListWrapper) ResetRender() {
	l.widget.BorderLabelFg = ui.ColorWhite
	l.widget.BorderFg = ui.ColorWhite
	l.widget.Items = l.listRenderer.ResetRender()
	ui.Render(ui.Body)
}

// GetWidget returns the instance of ui.List
func (l *ListWrapper) GetWidget() *ui.List {
	return l.widget
}

// GetCursor returns ListWrapper.cursor
func (l *ListWrapper) GetCursor() int {
	return l.listRenderer.GetCursor()
}

// SetBody replace strings on ListWrapper.body
func (l *ListWrapper) SetBody(items []string) {
	l.listRenderer.SetBody(items)
	l.widget.Items = items
}

// AddBody add an another line of textto ListWrapper.body
func (l *ListWrapper) AddBody(line string) {
	l.listRenderer.AddBody(line)
	l.widget.Items = append(l.widget.Items, line)
}
func (l *ListWrapper) moveAndRender(direction string) {
	l.gPressed = false // cancel gg to top
	l.widget.Items = l.listRenderer.Move(direction)
	ui.Render(ui.Body)
}
