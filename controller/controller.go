package controller

import (
	"path/filepath"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/github"
	"github.com/qmu/mcc/model"
	"github.com/qmu/mcc/widget"
)

// Controller controls termui's widget layout and keybindings
type Controller struct {
	execPath    string
	viewManager *model.ViewManager
	client      *github.Client
}

// NewController constructs a new Controller
func NewController(appVersion string, configSchemaVersion string, configPath string, debug bool) (err error) {
	d := new(Controller)
	d.execPath = filepath.Dir(configPath)

	// initialize termui
	if err = ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	d.viewManager, err = model.NewViewManager(&model.ConfigLoaderOption{
		ExecPath:            d.execPath,
		ConfigPath:          configPath,
		AppVersion:          appVersion,
		ConfigSchemaVersion: configSchemaVersion,
	})
	if err != nil {
		return
	}

	// layout first tab
	d.viewManager.SwitchTab(0)
	d.setKeyBindings()

	// init asynchronously
	if d.viewManager.HasWidget("github_issue") {
		go d.renderGitHubIssueWidgets()
	}

	if debug {
		return
	}
	ui.Loop()
	return
}

func (d *Controller) setKeyBindings() error {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-j", func(ui.Event) {
		d.viewManager.NextWidget("bottom")
	})
	ui.Handle("/sys/kbd/C-k", func(ui.Event) {
		d.viewManager.NextWidget("top")
	})
	ui.Handle("/sys/kbd/<backspace>", func(ui.Event) {
		d.viewManager.NextWidget("left")
	})
	ui.Handle("/sys/kbd/C-l", func(ui.Event) {
		d.viewManager.NextWidget("right")
	})
	ui.Handle("/sys/kbd/1", func(ui.Event) {
		d.viewManager.SwitchTab(0)
	})
	ui.Handle("/sys/kbd/2", func(ui.Event) {
		d.viewManager.SwitchTab(1)
	})
	ui.Handle("/sys/kbd/3", func(ui.Event) {
		d.viewManager.SwitchTab(2)
	})
	ui.Handle("/sys/kbd/4", func(ui.Event) {
		d.viewManager.SwitchTab(3)
	})
	ui.Handle("/sys/kbd/5", func(ui.Event) {
		d.viewManager.SwitchTab(4)
	})
	ui.Handle("/sys/kbd/6", func(ui.Event) {
		d.viewManager.SwitchTab(5)
	})
	ui.Handle("/sys/kbd/7", func(ui.Event) {
		d.viewManager.SwitchTab(6)
	})
	ui.Handle("/sys/kbd/8", func(ui.Event) {
		d.viewManager.SwitchTab(7)
	})
	ui.Handle("/sys/kbd/9", func(ui.Event) {
		d.viewManager.SwitchTab(8)
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})
	if d.viewManager.HasWidget("docker_status") {
		ui.Handle("/timer/1s", func(e ui.Event) {
			err := d.viewManager.MapWidgets(func(w *widget.WrapperWidget) (err error) {
				if w.Is("docker_status") {
					w.Activate()
				}
				return
			})
			if err != nil {
				return
			}
		})
	}

	return nil
}

func (d *Controller) renderGitHubIssueWidgets() {
	// initialize GitHub Client
	host := d.viewManager.GetGithubHost()
	if host == "" {
		host = "github.com"
	}
	c, err := github.NewClient(d.execPath, host)
	if err != nil {
		err = d.viewManager.MapWidgets(func(w *widget.WrapperWidget) (err error) {
			if w.Is("github_issue") {
				w.Disable()
			}
			return
		})
	} else {
		if err = c.Init(); err != nil {
			return
		}
		d.client = c
	}
	err = d.viewManager.MapWidgets(func(w *widget.WrapperWidget) (err error) {
		if w.Is("github_issue") && !w.IsDisabled() {
			w.SetOption(&widget.AdditionalWidgetOption{
				GithubClient: d.client,
			})
		}
		return
	})
}
