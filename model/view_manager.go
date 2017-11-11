package model

import (
	"errors"
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/model/vector"
	"github.com/qmu/mcc/utils"
	"github.com/qmu/mcc/widget"
)

// ViewManager load and unmarshal config file
type ViewManager struct {
	config            *ConfRoot
	configManager     *ConfigLoader
	execPath          string
	totalTab          int
	activeWidgetIndex int
	activeTabIndex    int
}

// NewViewManager constructs a ViewManager
func NewViewManager(opt *ConfigLoaderOption) (c *ViewManager, err error) {
	c = new(ViewManager)
	c.configManager, err = NewLoader(opt)
	c.activeTabIndex = -1
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

func (c *ViewManager) buildCollection() (err error) {
	collection, err := vector.NewRectangleCollection()
	if err != nil {
		return
	}
	windowH := ui.TermHeight() - 1
	windowW := ui.TermWidth() - 1
	idx := 0
	for i1, tab := range c.config.Layout {
		rowHTotal := 0
		c.config.Layout[i1].Index = i1
		for _, row := range tab.Rows {
			rowH := utils.Percentalize(windowH, row.Height)
			for i3, col := range row.Cols {
				cntH := 0
				stackHTotal := 0
				for _, stack := range col.Stacks {
					realWidth := windowW / len(row.Cols)
					realHeight := 0
					if len(col.Stacks)-1 == i3 {
						realHeight = rowH - cntH
					} else {
						realHeight = utils.Percentalize(rowH, stack.Height)
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
		c.totalTab++
	}

	collection.CalcDistances()

	idx = 0
	for i1, tab := range c.config.Layout {
		for i2, row := range tab.Rows {
			for i3, col := range row.Cols {
				for _, stack := range col.Stacks {
					wi, err := c.getWidgetByID(stack.ID)
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

func (c *ViewManager) getWidgetByID(id string) (result confWidget, err error) {
	for _, d := range c.config.Widgets {
		if d.ID == id {
			result = d
			return
		}
	}
	return result, errors.New("no widget named " + id)
}

func (c *ViewManager) getWidgetByIndex(index int) (result *widget.WrapperWidget) {
	c.MapWidgets(func(wi *widget.WrapperWidget) (err error) {
		if index == wi.Index {
			result = wi
		}
		return
	})
	return
}

func (c *ViewManager) activateFirstWidgetOnTab(idx int) {
	for _, r := range c.config.Layout[idx].Rows {
		for _, cl := range r.Cols {
			for _, wi := range cl.Widgets {
				if !wi.IsDisabled() && wi.IsReady() {
					wi.Activate()
					c.activeWidgetIndex = wi.Index
					return
				}
			}
		}
	}
}

func (c *ViewManager) renderTabPane(tab confTab) (err error) {
	ui.Clear()
	ui.Body.Rows = ui.Body.Rows[:0]

	tabs := []*ui.Row{}
	for i, t := range c.config.Layout {
		tabP := ui.NewList()
		color := "(fg-white,bg-default)"
		if tab.Index == t.Index {
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

// SwitchTab is
func (c *ViewManager) SwitchTab(tabIdx int) {
	if tabIdx > c.totalTab-1 || tabIdx == c.activeTabIndex {
		return
	}
	// layout header and body
	for i, tab := range c.config.Layout {
		if i == tabIdx {
			if err := c.renderTabPane(tab); err != nil {
				panic(err)
			}
		}
	}
	// deactivate all, and activate first widget
	c.MapWidgets(func(w *widget.WrapperWidget) (err error) {
		if w.Tab == tabIdx {
			w.Deactivate()
		}
		return
	})

	c.activateFirstWidgetOnTab(tabIdx)
	c.activeTabIndex = tabIdx
}

// NextWidget is
func (c *ViewManager) NextWidget(direction string) {
	from := c.getWidgetByIndex(c.activeWidgetIndex)
	toIdx := from.GetNeighborIndex(direction)
	to := c.getWidgetByIndex(toIdx)
	if from != nil && to != nil && !to.IsDisabled() && to.IsReady() {
		from.Deactivate()
		to.Activate()
		c.activeWidgetIndex = toIdx
	}
}

// HasWidget returns whether config contains 'widgetType' stack or not
func (c *ViewManager) HasWidget(widgetType string) bool {
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
func (c *ViewManager) GetActiveWidgetsOf(name string) (result []*widget.WrapperWidget) {
	c.MapWidgets(func(wi *widget.WrapperWidget) (err error) {
		if wi.WidgetType == name && !wi.IsDisabled() {
			result = append(result, wi)
		}
		return
	})
	return
}

// MapWidgets is
func (c *ViewManager) MapWidgets(fn func(*widget.WrapperWidget) error) (err error) {
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

// GetGithubHost is
func (c *ViewManager) GetGithubHost() string {
	return c.config.GitHubHost
}