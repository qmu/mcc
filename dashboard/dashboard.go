package dashboard

import (
	"path/filepath"
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/github"
)

// Dashboard controls termui's widget layout and keybindings
type Dashboard struct {
	execPath          string
	config            *Config
	widgetManager     *WidgetManager
	githubWidgets     []*GithubIssueWidget
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

	// init widgetManager
	d.widgetManager, err = NewWidgetManager(&WidgetManagerOptions{
		execPath:            d.execPath,
		configPath:          configPath,
		appVersion:          appVersion,
		configSchemaVersion: configSchemaVersion,
	})
	if err != nil {
		return
	}
	d.config = d.widgetManager.config

	// layout first tab
	d.switchTab(0)

	// init asynchronously
	if d.widgetManager.HasWidget("github_issue") {
		go d.renderGitHubIssueWidgets()
	}

	// init asynchronously
	if d.widgetManager.HasWidget("git_status") {
		go func() {
			for _, aw := range d.widgetManager.GetActiveWidgetsOf("git_status") {
				aw.Render()
			}
		}()
	}

	ui.Loop()
	return
}

func (d *Dashboard) switchTab(idx int) {
	if idx > len(d.config.Layout)-1 {
		return
	}
	tab := d.config.Layout[idx]
	d.activeTabIndex = idx
	// layout header and body
	if err := d.layout(tab); err != nil {
		panic(err)
	}
}

func (d *Dashboard) layout(tab Tab) (err error) {
	ui.Clear()
	ui.Body.Rows = ui.Body.Rows[:0]

	// header
	// header := ui.NewPar("Press q to quit, Press C-[j,k,h,l] to switch widget")
	// header.Height = 1
	// header.Border = false
	// header.TextFgColor = ui.ColorWhite

	tabs := []*ui.Row{}
	for i, t := range d.config.Layout {
		tabP := ui.NewList()
		color := "(fg-white,bg-default)"
		if tab.Name == t.Name {
			color = "(fg-white,bg-blue)"
		}
		space := strings.Repeat(" ", 500)
		tabP.Items = []string{"[ " + strconv.Itoa(i+1) + "." + t.Name + space + "]" + color}
		tabP.Height = 3
		tabP.Border = true
		tabP.BorderFg = ui.ColorBlue
		tabs = append(tabs, ui.NewCol(2, 0, tabP))
	}

	ui.Body.AddRows(
		ui.NewRow(tabs...))

	// body
	cnt := 0
	var newRows []*ui.Row
	d.githubWidgets = []*GithubIssueWidget{}
	for _, row := range tab.Rows {
		var newCols []*ui.Row
		for _, col := range row.Cols {
			var cols []ui.GridBufferer
			for _, w := range col.Widgets {
				gw := w.extendedWidget.GetGridBufferers()
				cols = append(cols, gw...)
				if w.Type == "github_issue" {
					d.githubWidgets = append(d.githubWidgets, w.extendedWidget.githubWidget)
				}
				cnt++
			}
			newCols = append(newCols, ui.NewCol(12/len(row.Cols), 0, cols...))
		}
		newRows = append(newRows, ui.NewRow(newCols...))
	}

	ui.Body.AddRows(newRows...)
	ui.Body.Align()
	ui.Render(ui.Body)

	// deactivate all, and activate first widget
	widgets := d.widgetManager.GetAllWidgets()
	for _, w := range widgets {
		w.Deactivate()
	}
	d.activateFirstWidgetOnTab(tab)

	return nil
}

func (d *Dashboard) renderGitHubIssueWidgets() {
	// initialize GitHub Client
	host := d.config.GitHubHost
	if d.config.GitHubHost == "" {
		host = "github.com"
	}
	c, err := github.NewClient(d.execPath, host)
	if err != nil {
		for _, w := range d.githubWidgets {
			w.Disable()
		}
	} else {
		if err = c.Init(); err != nil {
			return
		}
		d.client = c
	}

	for _, w := range d.githubWidgets {
		if !w.IsDisabled() {
			w.SetClient(d.client)
			w.Render()
		}
	}
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
		d.downerRowWidget()
	})
	// press Ctrl-k to previous widget
	ui.Handle("/sys/kbd/C-k", func(ui.Event) {
		d.upperRowWidget()
	})
	// press Ctrl-h to left col widget
	ui.Handle("/sys/kbd/<backspace>", func(ui.Event) {
		d.leftColWidget()
	})
	// press Ctrl-l to right col widget
	ui.Handle("/sys/kbd/C-l", func(ui.Event) {
		d.rightColWidget()
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
	if d.widgetManager.HasWidget("docker_status") {
		ui.Handle("/timer/1s", func(e ui.Event) {
			widgets := d.widgetManager.GetAllWidgets()
			for _, w := range widgets {
				if w.widgetType == "docker_status" && !w.IsDisabled() {
					w.Render()
				}
			}
		})
	}

	return nil
}

func (d *Dashboard) activateFirstWidgetOnTab(tab Tab) {
	w := d.widgetManager.GetWidgetByIndex(d.activeWidgetIndex)
	if w != nil {
		for _, r := range tab.Rows {
			for _, c := range r.Cols {
				for _, wi := range c.Widgets {
					if !wi.extendedWidget.IsDisabled() && wi.extendedWidget.IsReady() {
						w.Deactivate()
						d.activate(wi.extendedWidget)
						return
					}
				}
			}
		}
	}
}

func (d *Dashboard) activate(w *ExtendedWidget) {
	ui.ResetHandlers()
	w.Activate()
	d.setKeyBindings()
	d.activeWidgetIndex = w.index
}

func (d *Dashboard) downerRowWidget() {
	w := d.widgetManager.GetWidgetByIndex(d.activeWidgetIndex)
	if w != nil && w.bottomWidget != nil {
		w.Deactivate()
		d.activate(w.bottomWidget)
	}
}

func (d *Dashboard) upperRowWidget() {
	w := d.widgetManager.GetWidgetByIndex(d.activeWidgetIndex)
	if w != nil && w.topWidget != nil {
		w.Deactivate()
		d.activate(w.topWidget)
	}
}

func (d *Dashboard) rightColWidget() {
	w := d.widgetManager.GetWidgetByIndex(d.activeWidgetIndex)
	if w != nil && w.rightWidget != nil {
		w.Deactivate()
		d.activate(w.rightWidget)
	}
}

func (d *Dashboard) leftColWidget() {
	w := d.widgetManager.GetWidgetByIndex(d.activeWidgetIndex)
	if w != nil && w.leftWidget != nil {
		w.Deactivate()
		d.activate(w.leftWidget)
	}
}
