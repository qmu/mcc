package dashboard

import (
	"errors"

	ui "github.com/gizak/termui"
)

// ExtendedWidget is
type ExtendedWidget struct {
	index        int
	widget       Widget
	widgetType   string
	tab          int
	row          int
	col          int
	stack        int
	vOffset      int
	height       int
	width        int
	point        Point
	widgetter    Widgetter
	widthFrom    int // int of percentage e.g. 33(%)
	widthTo      int // int of percentage e.g. 66(%)
	title        string
	githubWidget *GithubIssueWidget

	bottomWidget *ExtendedWidget
	topWidget    *ExtendedWidget
	leftWidget   *ExtendedWidget
	rightWidget  *ExtendedWidget
}

// Point is
type Point struct {
	x int
	y int
}

// Widgetter define common interface for each widgets
type Widgetter interface {
	Activate()
	Deactivate()
	IsDisabled() bool
	IsReady() bool
	GetGridBufferers() []ui.GridBufferer
	GetHighlightenPos() int
	Render() error
}

// GetContent is
func (w *ExtendedWidget) GetContent() interface{} {
	return w.widget.Content
}

// Activate is
func (w *ExtendedWidget) Activate() {
	w.widgetter.Activate()
}

// Deactivate is
func (w *ExtendedWidget) Deactivate() {
	w.widgetter.Deactivate()
}

// IsDisabled is
func (w *ExtendedWidget) IsDisabled() bool {
	return w.widgetter.IsDisabled()
}

// IsReady is
func (w *ExtendedWidget) IsReady() bool {
	return w.widgetter.IsReady()
}

// GetGridBufferers is
func (w *ExtendedWidget) GetGridBufferers() []ui.GridBufferer {
	return w.widgetter.GetGridBufferers()
}

// GetHighlightenPos is
func (w *ExtendedWidget) GetHighlightenPos() int {
	return w.widgetter.GetHighlightenPos()
}

// Render is
func (w *ExtendedWidget) Render() error {
	return w.widgetter.Render()
}

// Vary is
func (w *ExtendedWidget) Vary(opt *WidgetOptions) (err error) {
	opt.extendedWidget = w
	var wi Widgetter
	switch w.widgetType {
	case "menu":
		wi, err = NewMenuWidget(opt)
	case "note":
		wi, err = NewNoteWidget(opt)
	case "github_issue":
		w.githubWidget, err = NewGithubIssueWidget(opt)
		if err != nil {
			return
		}
		wi = w.githubWidget
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
		return errors.New("Widget type \"" + w.widgetType + "\" is not supported")
	}
	w.widgetter = wi
	return
}
