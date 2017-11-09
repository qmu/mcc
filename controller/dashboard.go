package controller

import (
	"path/filepath"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/config"
	"github.com/qmu/mcc/github"
	"github.com/qmu/mcc/widget"
)

// Dashboard controls termui's widget layout and keybindings
type Dashboard struct {
	execPath         string
	widgetCollection *config.WidgetCollection
	client           *github.Client
}

// NewDashboard constructs a new Dashboard
func NewDashboard(appVersion string, configSchemaVersion string, configPath string) (err error) {
	d := new(Dashboard)
	d.execPath = filepath.Dir(configPath)

	// initialize termui
	if err = ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	d.widgetCollection, err = config.NewWidgetCollection(&config.LoaderOption{
		ExecPath:            d.execPath,
		ConfigPath:          configPath,
		AppVersion:          appVersion,
		ConfigSchemaVersion: configSchemaVersion,
	})
	if err != nil {
		return
	}

	// layout first tab
	d.widgetCollection.SwitchTab(0)
	d.setKeyBindings()

	// init asynchronously
	if d.widgetCollection.HasWidget("github_issue") {
		go d.renderGitHubIssueWidgets()
	}

	// init asynchronously
	if d.widgetCollection.HasWidget("git_status") {
		go func() {
			for _, aw := range d.widgetCollection.GetActiveWidgetsOf("git_status") {
				aw.Render()
			}
		}()
	}

	ui.Loop()
	return
}

func (d *Dashboard) setKeyBindings() error {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/kbd/C-j", func(ui.Event) {
		d.widgetCollection.NextWidget("bottom")
	})
	ui.Handle("/sys/kbd/C-k", func(ui.Event) {
		d.widgetCollection.NextWidget("top")
	})
	ui.Handle("/sys/kbd/<backspace>", func(ui.Event) {
		d.widgetCollection.NextWidget("left")
	})
	ui.Handle("/sys/kbd/C-l", func(ui.Event) {
		d.widgetCollection.NextWidget("right")
	})
	ui.Handle("/sys/kbd/1", func(ui.Event) {
		d.widgetCollection.SwitchTab(0)
	})
	ui.Handle("/sys/kbd/2", func(ui.Event) {
		d.widgetCollection.SwitchTab(1)
	})
	ui.Handle("/sys/kbd/3", func(ui.Event) {
		d.widgetCollection.SwitchTab(2)
	})
	ui.Handle("/sys/kbd/4", func(ui.Event) {
		d.widgetCollection.SwitchTab(3)
	})
	ui.Handle("/sys/kbd/5", func(ui.Event) {
		d.widgetCollection.SwitchTab(4)
	})
	ui.Handle("/sys/kbd/6", func(ui.Event) {
		d.widgetCollection.SwitchTab(5)
	})
	ui.Handle("/sys/kbd/7", func(ui.Event) {
		d.widgetCollection.SwitchTab(6)
	})
	ui.Handle("/sys/kbd/8", func(ui.Event) {
		d.widgetCollection.SwitchTab(7)
	})
	ui.Handle("/sys/kbd/9", func(ui.Event) {
		d.widgetCollection.SwitchTab(8)
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})
	if d.widgetCollection.HasWidget("docker_status") {
		ui.Handle("/timer/1s", func(e ui.Event) {
			err := d.widgetCollection.MapWidgets(func(w *widget.WrapperWidget) (err error) {
				if w.Is("docker_status") && !w.IsDisabled() {
					w.Render()
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

func (d *Dashboard) renderGitHubIssueWidgets() {
	// initialize GitHub Client
	host := d.widgetCollection.GetGithubHost()
	if host == "" {
		host = "github.com"
	}
	c, err := github.NewClient(d.execPath, host)
	if err != nil {
		err = d.widgetCollection.MapWidgets(func(w *widget.WrapperWidget) (err error) {
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
	err = d.widgetCollection.MapWidgets(func(w *widget.WrapperWidget) (err error) {
		if w.Is("github_issue") && !w.IsDisabled() {
			w.SetOption(&widget.AdditionalWidgetOption{
				GithubClient: d.client,
			})
			w.Activate()
			w.Render()
			w.Deactivate()
		}
		return
	})
}
