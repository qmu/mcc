package widget

import (
	"errors"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/config/vector"
)

// WrapperWidget is
type WrapperWidget struct {
	Index      int
	WidgetType string
	Tab        int
	Title      string
	Rectangle  *vector.Rectangle
	Envs       []map[string]string
	ExecPath   string
	Timezone   string
	Content    interface{}
	IssueRegex string
	Type       string
	Path       string
	widgetter  Widgetter
}

// Activate is
func (w *WrapperWidget) Activate() {
	w.widgetter.Activate()
}

// Deactivate is
func (w *WrapperWidget) Deactivate() {
	w.widgetter.Deactivate()
}

// Disable is
func (w *WrapperWidget) Disable() {
	w.widgetter.Disable()
}

// IsDisabled is
func (w *WrapperWidget) IsDisabled() bool {
	return w.widgetter.IsDisabled()
}

// IsReady is
func (w *WrapperWidget) IsReady() bool {
	return w.widgetter.IsReady()
}

// GetGridBufferers is
func (w *WrapperWidget) GetGridBufferers() []ui.GridBufferer {
	return w.widgetter.GetGridBufferers()
}

// GetHighlightenPos is
func (w *WrapperWidget) GetHighlightenPos() int {
	return w.widgetter.GetHighlightenPos()
}

// Render is
func (w *WrapperWidget) Render() error {
	return w.widgetter.Render()
}

// SetOption is
func (w *WrapperWidget) SetOption(opt *AdditionalWidgetOption) {
	w.widgetter.SetOption(opt)
}

// GetNeighborIndex is
func (w *WrapperWidget) GetNeighborIndex(direction string) (idx int) {
	if direction == "top" {
		idx = w.Rectangle.TopWidgetIndex
	} else if direction == "right" {
		idx = w.Rectangle.RightWidgetIndex
	} else if direction == "bottom" {
		idx = w.Rectangle.BottomWidgetIndex
	} else if direction == "left" {
		idx = w.Rectangle.LeftWidgetIndex
	} else {
		panic("direction only accepts 'top', 'right', 'bottom', 'left'")
	}
	return idx
}

// Is is
func (w *WrapperWidget) Is(wType string) bool {
	return w.WidgetType == wType
}

// Vary is
func (w *WrapperWidget) Vary() (err error) {
	var wi Widgetter
	opt := &Option{
		Envs:       w.Envs,
		ExecPath:   w.ExecPath,
		Timezone:   w.Timezone,
		Content:    w.Content,
		IssueRegex: w.IssueRegex,
		Height:     w.Rectangle.Height,
		Width:      w.Rectangle.Width,
		Title:      w.Title,
		Type:       w.Type,
		Path:       w.Path,
	}
	switch w.WidgetType {
	case "menu":
		wi, err = NewMenuWidget(opt)
	case "note":
		wi, err = NewNoteWidget(opt)
	case "github_issue":
		wi, err = NewGithubIssueWidget(opt)
	case "text_file":
		wi, err = NewNoteWidget(opt)
	case "git_status":
		wi, err = NewGitStatusWidget(opt)
	case "tail_file":
		wi, err = NewTailFileWidget(opt)
	case "docker_status":
		wi, err = NewDockerStatusWidget(opt)
	}
	if err != nil {
		return
	}
	if wi == nil {
		return errors.New("Widget type \"" + w.WidgetType + "\" is not supported")
	}
	w.widgetter = wi
	return
}
