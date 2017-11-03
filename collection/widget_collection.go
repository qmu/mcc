package collection

import (
	"errors"
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/collection/config"
	"github.com/qmu/mcc/collection/vector"
	"github.com/qmu/mcc/utils"
	"github.com/qmu/mcc/widget"
)

// WidgetCollection load and unmarshal config file
type WidgetCollection struct {
	config        *config.ConfRoot
	configManager *config.Loader
	execPath      string
	total         int
}

// NewWidgetCollection constructs a WidgetCollection
func NewWidgetCollection(opt *config.LoaderOption) (c *WidgetCollection, err error) {
	c = new(WidgetCollection)

	c.configManager, err = config.NewLoader(opt)
	if err != nil {
		return
	}
	c.config = c.configManager.GetConfig()

	c.execPath = opt.ExecPath
	err = c.buildCollection()
	if err != nil {
		return
	}

	return
}

func (c *WidgetCollection) buildCollection() (err error) {
	collection, err := vector.NewRectangleCollection()
	if err != nil {
		return
	}
	windowH := ui.TermHeight() - 1
	windowW := ui.TermWidth() - 1
	idx := 0
	for i1, tab := range c.config.Layout {
		rowHTotal := 0
		for _, row := range tab.Rows {
			rowH := utils.Percentalize(windowH, row.Height)
			for i3, col := range row.Cols {
				cntH := 0
				stackHTotal := 0
				for _, id := range col.Stacks {
					wi, err := c.getWidgetByID(id)
					if err != nil {
						return err
					}
					realWidth := windowW / len(row.Cols)
					realHeight := 0
					if len(col.Stacks)-1 == i3 {
						realHeight = rowH - cntH
					} else {
						realHeight = utils.Percentalize(rowH, wi.Height)
						cntH += realHeight
					}
					collection.Register(&vector.RectangleOptions{
						Index:    idx,
						WindowW:  windowW,
						WindowH:  windowH,
						TabIndex: i1,
						CenterX:  windowW/len(row.Cols)*i3 + (realWidth / 2),
						CenterY:  rowHTotal + stackHTotal + (realHeight / 2),
						Width:    realWidth,
						Height:   realHeight,
					})
					idx++
					stackHTotal += realHeight
				}
			}
			rowHTotal += rowH
		}
	}
	c.total = idx + 1

	collection.CalcDistances()

	idx = 0
	for i1, tab := range c.config.Layout {
		for i2, row := range tab.Rows {
			for i3, col := range row.Cols {
				for _, id := range col.Stacks {
					wi, err := c.getWidgetByID(id)
					if err != nil {
						return err
					}
					ew := &widget.WrapperWidget{
						Index:      idx,
						WidgetType: wi.Type,
						Tab:        i1,
						Title:      wi.Title,
						Rectangle:  collection.GetRectangle(idx),
						Envs:       c.config.Envs,
						ExecPath:   c.execPath,
						Timezone:   c.config.Timezone,
						Content:    wi.Content,
						IssueRegex: wi.IssueRegex,
						Type:       wi.Type,
						Path:       wi.Path,
					}
					if err != nil {
						return err
					}
					c.config.Layout[i1].Rows[i2].Cols[i3].Widgets = append(c.config.Layout[i1].Rows[i2].Cols[i3].Widgets, ew)
					idx++
				}
			}
		}
	}
	c.MapWidgets(func(wi *widget.WrapperWidget) (err error) {
		err = wi.Vary()
		return
	})

	return
}

func (c *WidgetCollection) getWidgetByID(id string) (result config.ConfWidget, err error) {
	for _, d := range c.config.Widgets {
		if d.ID == id {
			result = d
			return
		}
	}
	return result, errors.New("no widget named " + id)
}

// GetConfig is
func (c *WidgetCollection) GetConfig() *config.ConfRoot {
	return c.config
}

// HasWidget returns whether config contains 'widgetType' stack or not
func (c *WidgetCollection) HasWidget(widgetType string) bool {
	result := false
	c.MapWidgets(func(wi *widget.WrapperWidget) (err error) {
		if wi.WidgetType == widgetType {
			result = true
		}
		return
	})
	return result
}

// GetActiveWidgetsOf is
func (c *WidgetCollection) GetActiveWidgetsOf(name string) (result []*widget.WrapperWidget) {
	c.MapWidgets(func(wi *widget.WrapperWidget) (err error) {
		if wi.WidgetType == name && !wi.IsDisabled() {
			result = append(result, wi)
		}

		return
	})
	return
}

// GetWidgetByIndex is
func (c *WidgetCollection) GetWidgetByIndex(index int) (result *widget.WrapperWidget) {
	c.MapWidgets(func(wi *widget.WrapperWidget) (err error) {
		if index == wi.Index {
			result = wi
		}
		return
	})
	return
}

// GetTabByTabIndex is
func (c *WidgetCollection) GetTabByTabIndex(tabIndex int) (result config.ConfTab) {
	for i, tab := range c.config.Layout {
		if i == tabIndex {
			result = tab
			return
		}
	}
	return
}

// MapWidgets is
func (c *WidgetCollection) MapWidgets(fn func(*widget.WrapperWidget) error) (err error) {
	for _, tab := range c.config.Layout {
		for _, row := range tab.Rows {
			for _, col := range row.Cols {
				for _, wi := range col.Widgets {
					err = fn(wi)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// Count is
func (c *WidgetCollection) Count() int {
	return c.total
}

// Render is
func (c *WidgetCollection) Render(tab config.ConfTab) (err error) {
	ui.Clear()
	ui.Body.Rows = ui.Body.Rows[:0]

	tabs := []*ui.Row{}
	for i, t := range c.config.Layout {
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
	for _, row := range tab.Rows {
		var newCols []*ui.Row
		for _, col := range row.Cols {
			var cols []ui.GridBufferer
			for _, w := range col.Widgets {
				gw := w.GetGridBufferers()
				cols = append(cols, gw...)
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

// GetGithubHost is
func (c *WidgetCollection) GetGithubHost() string {
	return c.config.GitHubHost
}
