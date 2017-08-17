package dashboard

import (
	"fmt"
	"os"
	"path/filepath"

	ui "github.com/gizak/termui"
	version "github.com/hashicorp/go-version"
	"github.com/qmu/mcc/github"
	// "github.com/k0kubun/pp"
)

// WidgetItem define common interface for each widgets
type WidgetItem interface {
	Activate()
	Deactivate()
	IsDisabled() bool
	IsReady() bool
	GetWidget() *ui.List
	GetHighlightenPos() int
}

// Dashboard controls termui's widget layout and keybindings
type Dashboard struct {
	execPath        string
	configPath      string
	config          *Config
	widgetPositions []*widgetPosition
	githubWidgets   []*GithubIssueWidget
	active          int // active widget index of Dashboard.widgetPositions
	client          *github.Client
}

type widgetPosition struct {
	row        int
	col        int
	stack      int
	vOffset    int
	height     int
	widgetItem WidgetItem
}

// NewDashboard constructs a new Dashboard
func NewDashboard(appVersion string, configSchemaVersion string, configPath string) (err error) {
	d := new(Dashboard)
	d.execPath = filepath.Dir(configPath)
	d.configPath = configPath

	// initialize GitHub Client
	c, err := github.NewClient(d.execPath)
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

	// initialize termui
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// load config file
	configManager, err := NewConfigManager(d.configPath)
	if err != nil {
		return
	}
	d.config = configManager.LoadedData

	// check the ConfigSchemaVersion
	vApp, err := version.NewVersion(configSchemaVersion)
	vConfig, err := version.NewVersion(d.config.SchemaVersion)
	if err != nil {
		return
	}
	if vConfig.LessThan(vApp) {
		fmt.Printf("mcc %s supports schema_version %s but ths schema_version in %s seems to be %s\n", appVersion, vApp, configPath, vConfig)
		fmt.Printf("please upgrade mcc or %s first\n", configPath)
		os.Exit(1)
	}

	// initialize interface
	if err = d.prepareUI(); err != nil {
		return
	}
	return
}

func (d *Dashboard) prepareUI() (err error) {
	d.layoutHeader()
	d.layoutWidgets()
	for _, w := range d.widgetPositions {
		d.deactivateWidget(w.widgetItem)
	}
	d.activateWidget(d.widgetPositions[0].widgetItem)

	if d.hasGithubIssueWidget() {
		go func() {
			for _, w := range d.githubWidgets {
				if !w.IsDisabled() {
					w.Render(d.client)
				}
			}
		}()
	}

	ui.Loop()
	return nil
}

func (d *Dashboard) activateWidget(w WidgetItem) {
	w.Activate()
	d.setKeyBindings()
}

func (d *Dashboard) deactivateWidget(w WidgetItem) {
	w.Deactivate()
	ui.ResetHandlers()
}

func (d *Dashboard) hasGithubIssueWidget() bool {
	result := false
	for _, row := range d.config.Rows {
		for _, col := range row.Cols {
			for _, w := range col.Stacks {
				if w.Type == "github_issue" {
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

	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	return nil
}

func (d *Dashboard) downerRowWidget() {
	c := d.widgetPositions[d.active]
	num := d.active + 1
	for i := 0; i < len(d.widgetPositions)-1; i++ {
		w := d.widgetPositions[i]
		if c.row < w.row && w.stack == 0 && w.col <= c.col {
			num = i
		}
	}
	d.moveWidget(num)
}

func (d *Dashboard) upperRowWidget() {
	c := d.widgetPositions[d.active]
	if c.row > 0 && c.stack == 0 {
		for i := len(d.widgetPositions) - 1; i >= 0; i-- {
			var row int
			if i < d.active && d.widgetPositions[i].row <= row {
				if d.widgetPositions[i].col <= c.col {
					d.moveWidget(i)
					break
				}
			} else {
				row = d.widgetPositions[i].row
			}
		}
	} else {
		d.moveWidget(d.active - 1)
	}
}

func (d *Dashboard) rightColWidget() {
	col := 0
	c := d.widgetPositions[d.active]
	cursor := c.widgetItem.GetHighlightenPos()
	myPosition := c.vOffset + cursor
	for i, w := range d.widgetPositions {
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
	col := d.widgetPositions[len(d.widgetPositions)-1].col
	c := d.widgetPositions[d.active]
	cursor := c.widgetItem.GetHighlightenPos()
	myPosition := c.vOffset + cursor
	for i := len(d.widgetPositions) - 1; i >= 0; i-- {
		wp := d.widgetPositions[i]
		if i < d.active && col > wp.col {
			if myPosition >= wp.vOffset && myPosition <= wp.vOffset+wp.height {
				d.moveWidget(i)
				break
			}
		} else {
			col = wp.col
		}
	}
}

func (d *Dashboard) moveWidget(posIdx int) {
	w := d.getActiveWidget()
	d.deactivateWidget(w)
	ini := d.active
	if posIdx > len(d.widgetPositions)-1 {
		d.active = 0
	} else if posIdx < 0 {
		d.active = len(d.widgetPositions) - 1
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
	for k, w := range d.widgetPositions {
		if d.active == k {
			return w.widgetItem
		}
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
			for j, w := range col.Stacks {
				switch w.Type {
				case "menu":
					wi, err = NewMenuWidget(w)
					if err != nil {
						return err
					}
				case "note":
					wi, err = NewNoteWidget(w)
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
				}
				gw := wi.GetWidget()
				cols = append(cols, gw)
				d.widgetPositions = append(d.widgetPositions, &widgetPosition{
					row:        h,
					col:        i,
					stack:      j,
					vOffset:    offset,
					height:     gw.Height,
					widgetItem: wi,
				})
				offset += gw.Height
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
