package widget

import (
	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/github"
)

// Option defines options for each type of widgets
type Option struct {
	Envs       []map[string]string
	ExecPath   string
	Timezone   string
	Content    interface{}
	IssueRegex string
	Height     int
	Width      int
	Title      string
	Type       string
	Path       string
}

// GetHeight is
func (w *Option) GetHeight() int {
	return w.Height
}

// GetWidth is
func (w *Option) GetWidth() int {
	return w.Width
}

// GetTitle is
func (w *Option) GetTitle() string {
	return w.Title
}

// Menu is the schema implements Config.Widgets.Menu
type Menu struct {
	Category    string
	Name        string
	Description string
	Command     string
}

// Container is the schema implements Config.Widgets.Conttainer
type Container struct {
	Metrics   string
	Name      string
	Container string
}

// AdditionalWidgetOption is
type AdditionalWidgetOption struct {
	done         chan bool
	GithubClient *github.Client
}

// Widgetter define common interface for each widgets
type Widgetter interface {
	Activate()
	Deactivate()
	SetOption(opt *AdditionalWidgetOption)
	IsReady() bool
	GetGridBufferers() []ui.GridBufferer
	Init() error
}
