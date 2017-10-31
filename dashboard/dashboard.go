package dashboard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	ui "github.com/gizak/termui"
	version "github.com/hashicorp/go-version"
	"github.com/qmu/mcc/github"
)

// Dashboard controls termui's widget layout and keybindings
type Dashboard struct {
	execPath            string
	configPath          string
	config              *Config
	widgets             []*widgetInfo
	githubWidgets       []*GithubIssueWidget
	active              int // active widget index of Dashboard.widgets
	client              *github.Client
	header              *ui.Par // for debug
	appVersion          string
	configSchemaVersion string
}

type widgetInfo struct {
	widgetType string
	row        int
	col        int
	stack      int
	vOffset    int
	height     int
	widgetItem WidgetItem
	widthFrom  int
	widthTo    int
}

// WidgetItem define common interface for each widgets
type WidgetItem interface {
	Activate()
	Deactivate()
	IsDisabled() bool
	IsReady() bool
	GetWidget() []ui.GridBufferer
	GetHighlightenPos() int
	Render() error
}

// NewDashboard constructs a new Dashboard
func NewDashboard(appVersion string, configSchemaVersion string, configPath string) (err error) {
	d := new(Dashboard)
	d.execPath = filepath.Dir(configPath)
	d.configPath = configPath
	d.appVersion = appVersion
	d.configSchemaVersion = configSchemaVersion

	// initialize termui
	if err = ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// load config file
	configManager, err := NewConfigManager(d.configPath)
	if err != nil {
		return
	}
	d.config = configManager.LoadedData

	// check ConfigSchemaVersion
	if err = d.checkConfigScheme(); err != nil {
		return
	}

	// initialize interface
	if err = d.prepareUI(); err != nil {
		return
	}
	return
}

func (d *Dashboard) prepareUI() (err error) {
	// layout header and body
	d.layoutHeader()
	if err = d.layoutWidgets(); err != nil {
		return
	}

	// deactivate all, and activate first widget
	for _, w := range d.widgets {
		d.deactivateWidget(w.widgetItem)
	}
	d.activateWidget(d.widgets[0].widgetItem)

	// init asynchronously
	if d.hasWidget("github_issue") {
		go d.renderGitHubIssueWidgets()
	}

	// init asynchronously
	if d.hasWidget("git_status") {
		go d.renderGitStatusWidgets()
	}

	ui.Loop()
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

func (d *Dashboard) renderGitStatusWidgets() {
	for _, w := range d.widgets {
		if w.widgetType == "git_status" && !w.widgetItem.IsDisabled() {
			w.widgetItem.Render()
		}
	}
}

func (d *Dashboard) checkConfigScheme() (err error) {
	vApp, err := version.NewVersion(d.configSchemaVersion)
	vConfig, err := version.NewVersion(d.config.SchemaVersion)
	if err != nil {
		return
	}
	if vConfig.LessThan(vApp) {
		fmt.Printf("mcc %s supports schema_version %s but ths schema_version in %s seems to be %s\n", d.appVersion, vApp, d.configPath, vConfig)
		fmt.Printf("please upgrade mcc or %s first\n", d.configPath)
		os.Exit(1)
	}
	return
}

func (d *Dashboard) activateWidget(w WidgetItem) {
	w.Activate()
	d.setKeyBindings()
}

func (d *Dashboard) deactivateWidget(w WidgetItem) {
	w.Deactivate()
	ui.ResetHandlers()
}

func (d *Dashboard) hasWidget(widgetType string) bool {
	result := false
	for _, row := range d.config.Rows {
		for _, col := range row.Cols {
			for _, w := range col.Widgets {
				if w.Type == widgetType {
					result = true
				}
			}
		}
	}
	return result
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
	// press tab to switch widget
	ui.Handle("/sys/kbd/<tab>", func(ui.Event) {
		d.nextWidget()
	})

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	if d.hasWidget("docker_status") {
		ui.Handle("/timer/1s", func(e ui.Event) {
			for _, w := range d.widgets {
				if w.widgetType == "docker_status" && !w.widgetItem.IsDisabled() {
					w.widgetItem.Render()
				}
			}
		})
	}

	return nil
}

func (d *Dashboard) nextWidget() {
	d.moveWidget(d.active + 1)
}

func (d *Dashboard) downerRowWidget() {
	c := d.widgets[d.active]
	num := d.active + 1
	for i := 0; i < len(d.widgets); i++ {
		w := d.widgets[i]
		cond1 := c.row == w.row && c.col == w.col && w.stack == c.stack+1
		cond2 := c.row < w.row && w.stack == 0 && w.widthFrom <= c.widthFrom && c.widthFrom <= w.widthTo
		if cond1 || cond2 {
			num = i
			break
		}
	}
	d.moveWidget(num)
}

func (d *Dashboard) upperRowWidget() {
	c := d.widgets[d.active]
	if c.row == 0 && c.stack == 0 {
		d.moveWidget(d.active - 1)
	}
	for i := len(d.widgets) - 1; i >= 0; i-- {
		w := d.widgets[i]
		if i < d.active {
			cond1 := w.row == c.row && w.col == c.col && w.stack == c.stack-1
			cond2 := w.row < c.row && w.widthFrom <= c.widthFrom && c.widthFrom <= w.widthTo
			if cond1 || cond2 {
				d.moveWidget(i)
				break
			}
		}
	}

}

func (d *Dashboard) rightColWidget() {
	col := 0
	c := d.widgets[d.active]
	cursor := c.widgetItem.GetHighlightenPos()
	myPosition := c.vOffset + cursor
	for i, w := range d.widgets {
		if w.row == c.row && i > d.active && col < w.col {
			if myPosition >= w.vOffset && myPosition <= w.vOffset+w.height {
				d.moveWidget(i)
				break
			}
		} else {
			col = w.col
		}
	}
}

func (d *Dashboard) leftColWidget() {
	col := d.widgets[len(d.widgets)-1].col
	c := d.widgets[d.active]
	cursor := c.widgetItem.GetHighlightenPos()
	myPosition := c.vOffset + cursor
	for i := len(d.widgets) - 1; i >= 0; i-- {
		w := d.widgets[i]
		if i < d.active && col > w.col && c.row == w.row {
			if myPosition >= w.vOffset && myPosition <= w.vOffset+w.height {
				d.moveWidget(i)
				break
			}
		} else {
			col = w.col
		}
	}
}

func (d *Dashboard) moveWidget(posIdx int) {
	w := d.getActiveWidget()
	d.deactivateWidget(w)
	ini := d.active
	if posIdx > len(d.widgets)-1 {
		d.active = 0
	} else if posIdx < 0 {
		d.active = len(d.widgets) - 1
	} else {
		d.active = posIdx
	}
	w = d.getActiveWidget()
	d.activateWidget(w)
	if ini < posIdx { // moving foward
		if !w.IsReady() || w.IsDisabled() {
			d.downerRowWidget()
		}
	} else { // moving backfoward
		if !w.IsReady() || w.IsDisabled() {
			d.upperRowWidget()
		}

	}
}

func (d *Dashboard) getActiveWidget() (w WidgetItem) {
	for k, w := range d.widgets {
		if d.active == k {
			return w.widgetItem
		}
	}
	return nil
}

func (d *Dashboard) layoutHeader() {
	header := ui.NewPar("Press q to quit, Press tab or C-[j,k,h,l] to switch widget, Press j or k to move cursor")
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
	for h, row := range d.config.Rows {
		var newCols []*ui.Row
		colsCnt := len(row.Cols)
		for i, col := range row.Cols {
			cnt++
			var cols []ui.GridBufferer
			var wi WidgetItem
			offset := 1
			for j, w := range col.Widgets {
				switch w.Type {
				case "menu":
					wi, err = NewMenuWidget(w, d.config.Envs)
					if err != nil {
						return err
					}
				case "note":
					wi, err = NewNoteWidget(w, d.execPath)
					if err != nil {
						return err
					}
				case "github_issue":
					var giw *GithubIssueWidget
					giw, err = NewGithubIssueWidget(w, d.config.Timezone)
					if err != nil {
						return err
					}
					wi = giw
					d.githubWidgets = append(d.githubWidgets, giw)
				case "text_file":
					wi, err = NewNoteWidget(w, d.execPath)
					if err != nil {
						return err
					}
				case "git_status":
					wi, err = NewGitStatusWidget(w, d.execPath, d.config.Envs)
					if err != nil {
						return err
					}
				case "tail_file":
					wi, err = NewTailFileWidget(w, d.execPath)
					if err != nil {
						return err
					}
				case "docker_status":
					wi, err = NewDockerStatusWidget(w)
					if err != nil {
						return err
					}
				}
				if wi == nil {
					return errors.New("Widget type \"" + w.Type + "\" is not supported")
				}
				gw := wi.GetWidget()
				cols = append(cols, gw...)
				d.widgets = append(d.widgets, &widgetInfo{
					widgetType: w.Type,
					row:        h,
					col:        i,
					stack:      j,
					vOffset:    offset,
					height:     w.RealHeight,
					widgetItem: wi,
					widthFrom:  100 / len(row.Cols) * i,
					widthTo:    100 / len(row.Cols) * (i + 1),
				})
				offset += w.RealHeight
			}
			newCols = append(newCols, ui.NewCol(12/colsCnt, 0, cols...))
		}
		newRows = append(newRows, ui.NewRow(newCols...))
	}
	ui.Body.AddRows(newRows...)
	ui.Body.Align()
	ui.Render(ui.Body)

	return nil
}
