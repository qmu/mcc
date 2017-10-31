package dashboard

import (
	"path/filepath"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/github"
)

// Dashboard controls termui's widget layout and keybindings
type Dashboard struct {
	execPath      string
	config        *Config
	widgetManager *WidgetManager
	githubWidgets []*GithubIssueWidget
	client        *github.Client
	header        *ui.Par // for debug
	active        int
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
		configPath:          configPath,
		appVersion:          appVersion,
		configSchemaVersion: configSchemaVersion,
	})
	if err != nil {
		return
	}
	d.config = d.widgetManager.config

	// initialize interface
	if err = d.prepareUI(); err != nil {
		return
	}

	ui.Loop()
	return
}

func (d *Dashboard) prepareUI() (err error) {
	// layout header and body
	d.layoutHeader()
	if err = d.layoutWidgets(); err != nil {
		return
	}

	// deactivate all, and activate first widget
	widgets := d.widgetManager.GetAllWidgets()
	for _, w := range widgets {
		d.deactivate(w)
	}
	d.activate(widgets[0])

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
	return nil
}

func (d *Dashboard) layoutHeader() {
	header := ui.NewPar("Press q to quit, Press C-[j,k,h,l] to switch widget, Press j or k to move cursor")
	header.Height = 1
	header.Border = false
	header.TextFgColor = ui.ColorWhite
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, header)))
	d.header = header
}

func (d *Dashboard) layoutWidgets() (err error) {
	cnt := 0
	var newRows []*ui.Row
	for _, row := range d.config.Rows {
		var newCols []*ui.Row
		for _, col := range row.Cols {
			var cols []ui.GridBufferer
			for _, w := range col.Widgets {
				ew := w.extendedWidget
				err = ew.Vary(&WidgetOptions{
					envs:     d.config.Envs,
					execPath: d.execPath,
					timezone: d.config.Timezone,
				})
				if err != nil {
					return err
				}
				gw := ew.GetGridBufferers()
				cols = append(cols, gw...)
				if w.Type == "github_issue" {
					d.githubWidgets = append(d.githubWidgets, ew.githubWidget)
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

func (d *Dashboard) downerRowWidget() {
	w1 := d.widgetManager.GetWidgetByIndex(d.active)
	w2 := d.widgetManager.GetDownerWidget(d.active)
	if w1 != nil && w2 != nil {
		d.deactivate(w1)
		d.activate(w2)
	}
}

func (d *Dashboard) upperRowWidget() {
	w1 := d.widgetManager.GetWidgetByIndex(d.active)
	w2 := d.widgetManager.GetUpperWidget(d.active)
	if w1 != nil && w2 != nil {
		d.deactivate(w1)
		d.activate(w2)
	}
}

func (d *Dashboard) rightColWidget() {
	w1 := d.widgetManager.GetWidgetByIndex(d.active)
	w2 := d.widgetManager.GetRightWidget(d.active)
	if w1 != nil && w2 != nil {
		d.deactivate(w1)
		d.activate(w2)
	}
}

func (d *Dashboard) leftColWidget() {
	w1 := d.widgetManager.GetWidgetByIndex(d.active)
	w2 := d.widgetManager.GetLeftWidget(d.active)
	if w1 != nil && w2 != nil {
		d.deactivate(w1)
		d.activate(w2)
	}
}

func (d *Dashboard) deactivate(w *ExtendedWidget) {
	w.Deactivate()
	ui.ResetHandlers()
}

func (d *Dashboard) activate(w *ExtendedWidget) {
	w.Activate()
	d.setKeyBindings()
	total := d.widgetManager.GetAllWidgetsCount()
	if w.index > total-1 {
		d.active = 0
	} else if w.index < 0 {
		d.active = total - 1
	} else {
		d.active = w.index
	}
}
