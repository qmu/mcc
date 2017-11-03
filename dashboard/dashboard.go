package dashboard

import (
	"path/filepath"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/collection"
	"github.com/qmu/mcc/collection/config"
	"github.com/qmu/mcc/github"
	"github.com/qmu/mcc/widget"
)

// Dashboard controls termui's widget layout and keybindings
type Dashboard struct {
	execPath          string
	widgetCollection  *collection.WidgetCollection
	client            *github.Client
	activeTabIndex    int
	activeWidgetIndex int
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

	d.widgetCollection, err = collection.NewWidgetCollection(&config.LoaderOption{
		ExecPath:            d.execPath,
		ConfigPath:          configPath,
		AppVersion:          appVersion,
		ConfigSchemaVersion: configSchemaVersion,
	})
	if err != nil {
		return
	}

	// deactivate all, and activate first widget
	err = d.widgetCollection.MapWidgets(func(w *widget.WrapperWidget) (err error) {
		w.Deactivate()
		return
	})
	if err != nil {
		return
	}

	// layout first tab
	d.switchTab(0)

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

func (d *Dashboard) switchTab(idx int) {
	if idx > d.widgetCollection.Count()-1 {
		return
	}
	tab := d.widgetCollection.GetTabByTabIndex(idx)

	d.activeTabIndex = idx

	dw := d.widgetCollection.GetWidgetByIndex(d.activeWidgetIndex)
	dw.Deactivate()
	// layout header and body
	if err := d.widgetCollection.Render(tab); err != nil {
		panic(err)
	}
	d.activateFirstWidgetOnTab(tab)
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
		}
		return
	})
}

func (d *Dashboard) setKeyBindings() error {
	// press q to quit
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	// press Ctrl-c to quit
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})
	// press Ctrl-j to next widget
	ui.Handle("/sys/kbd/C-j", func(ui.Event) {
		d.nextWidget("bottom")
	})
	// press Ctrl-k to previous widget
	ui.Handle("/sys/kbd/C-k", func(ui.Event) {
		d.nextWidget("top")
	})
	// press Ctrl-h to left col widget
	ui.Handle("/sys/kbd/<backspace>", func(ui.Event) {
		d.nextWidget("left")
	})
	// press Ctrl-l to right col widget
	ui.Handle("/sys/kbd/C-l", func(ui.Event) {
		d.nextWidget("right")
	})
	// press 1 to switch tab
	ui.Handle("/sys/kbd/1", func(ui.Event) {
		d.switchTab(0)
	})
	// press 2 to switch tab
	ui.Handle("/sys/kbd/2", func(ui.Event) {
		d.switchTab(1)
	})
	// press 3 to switch tab
	ui.Handle("/sys/kbd/3", func(ui.Event) {
		d.switchTab(2)
	})
	// press 4 to switch tab
	ui.Handle("/sys/kbd/4", func(ui.Event) {
		d.switchTab(3)
	})
	// press 5 to switch tab
	ui.Handle("/sys/kbd/5", func(ui.Event) {
		d.switchTab(4)
	})
	// press 6 to switch tab
	ui.Handle("/sys/kbd/6", func(ui.Event) {
		d.switchTab(5)
	})
	// press 7 to switch tab
	ui.Handle("/sys/kbd/7", func(ui.Event) {
		d.switchTab(6)
	})
	// press 8 to switch tab
	ui.Handle("/sys/kbd/8", func(ui.Event) {
		d.switchTab(7)
	})
	// press 9 to switch tab
	ui.Handle("/sys/kbd/9", func(ui.Event) {
		d.switchTab(8)
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

func (d *Dashboard) activateFirstWidgetOnTab(tab config.ConfTab) {
	w := d.widgetCollection.GetWidgetByIndex(d.activeWidgetIndex)
	if w != nil {
		for _, r := range tab.Rows {
			for _, c := range r.Cols {
				for _, wi := range c.Widgets {
					if !wi.IsDisabled() && wi.IsReady() {
						w.Deactivate()
						d.activate(wi)
						return
					}
				}
			}
		}
	}
}

func (d *Dashboard) activate(w *widget.WrapperWidget) {
	ui.ResetHandlers()
	w.Activate()
	d.setKeyBindings()
	d.activeWidgetIndex = w.Index
}

func (d *Dashboard) nextWidget(direction string) {
	from := d.widgetCollection.GetWidgetByIndex(d.activeWidgetIndex)
	toIdx := from.GetNeighborIndex(direction)
	to := d.widgetCollection.GetWidgetByIndex(toIdx)
	if from != nil && to != nil {
		from.Deactivate()
		d.activate(to)
	}
}
