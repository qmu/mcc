package listable

import "strings"

// ListRenderer make a List widget which includes
// multi-line texts look like scrolled
type ListRenderer struct {
	header        []string
	body          []string
	maxH          int
	top           int
	bottom        int
	lineHighLight bool
	cursor        int
	unfocused     bool
}

// ListRendererOption is the option argument for NewListWrapper
type ListRendererOption struct {
	MaxH          int
	Header        []string
	Body          []string
	LineHighLight bool
}

// NewListRenderer constructs a ListWrapper
func NewListRenderer(opt *ListRendererOption) (l *ListRenderer) {
	l = new(ListRenderer)
	if opt.Header == nil {
		l.header = []string{}
	} else {
		l.header = opt.Header
	}
	if opt.Body == nil {
		l.body = []string{}
	} else {
		l.body = opt.Body
	}
	l.maxH = opt.MaxH
	l.top = 0
	l.bottom = opt.MaxH - len(l.header) - 3
	l.lineHighLight = opt.LineHighLight

	return
}

// RenderActually renders list
func (l *ListRenderer) RenderActually() []string {
	lH := len(l.header)
	lB := 0
	for k := range l.body {
		if l.maxH-2 <= lH+lB {
			break
		}
		if l.top > k {
			continue
		}
		lB++
	}
	items := make([]string, lH+lB)
	for i, h := range l.header {
		items[i] = h
	}
	cidx := 0
	for k, v := range l.body {
		if l.maxH-2 <= lH+cidx {
			break
		}
		if l.top > k {
			continue
		}
		if k == l.cursor {
			if l.lineHighLight {
				v = strings.Replace(v, "[", "", -1)
				v = strings.Replace(v, "]", "", -1)
				v = strings.Replace(v, "(fg-blue)", "", -1)
				v = strings.Replace(v, "(fg-green)", "", -1)
				v = strings.Replace(v, "(fg-red)", "", -1)
				v = "[" + v + "](fg-black,bg-green)"
			} else {
				v2 := ""
				for i, c := range v {
					if i == 0 && l.unfocused == false {
						v2 += "[" + string(c) + "](bg-green)"
					} else {
						v2 += string(c)
					}
				}
				v = v2
			}
		}
		items[lH+cidx] = v
		cidx++
	}
	return items
}

// ResetRender returns a initial multi-line texts
func (l *ListRenderer) ResetRender() []string {
	lH := len(l.header)
	lB := 0
	for i := 0; i < len(l.body); i++ {
		if l.maxH-2 <= lH+lB {
			break
		}
		lB++
	}
	var items = make([]string, lH+lB)
	for i, h := range l.header {
		items[i] = h
	}
	cidx := 0
	for _, v := range l.body {
		if l.maxH-2 <= lH+cidx {
			break
		}
		items[lH+cidx] = v
		cidx++
	}
	return items
}

// MoveCursor moves cursor position to "direction" with no focuse
func (l *ListRenderer) MoveCursor(direction string) (items []string) {
	l.unfocused = true
	items = l.MoveCursorWithFocus(direction)
	l.unfocused = false
	return
}

// MoveCursorWithFocus moves cursor position to "direction"
// with a highlightened focus
func (l *ListRenderer) MoveCursorWithFocus(direction string) (items []string) {
	if direction == "up" {
		items = l.up()
	} else if direction == "down" {
		items = l.down()
	} else if direction == "top" {
		c := len(l.body)
		for i := 0; i < c; i++ {
			items = l.up()
		}
	} else if direction == "bottom" {
		c := len(l.body) - l.GetCursor() - 1
		for i := 0; i < c; i++ {
			items = l.down()
		}
	}
	return
}

func (l *ListRenderer) up() (items []string) {
	if l.cursor > 0 {
		l.cursor--
	}
	if l.cursor < l.top {
		l.top--
		l.bottom--
	}
	return l.RenderActually()
}

func (l *ListRenderer) down() (items []string) {
	if len(l.body)-1 > l.cursor {
		l.cursor++
	}
	if l.cursor > l.bottom {
		l.top++
		l.bottom++
	}
	return l.RenderActually()
}

// GetCursor returns ListWrapper.cursor
func (l *ListRenderer) GetCursor() int {
	return l.cursor
}

// SetBody replace strings on ListWrapper.body
func (l *ListRenderer) SetBody(items []string) {
	l.body = items
}

// AddBody add an another line of text to ListWrapper.body
func (l *ListRenderer) AddBody(line string) {
	l.body = append(l.body, line)
}
