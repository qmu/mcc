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

	// remember widget walk
	lastWidget    *widget.WrapperWidget
	lastDirection string
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
		for i2, row := range tab.Rows {
			colWTotal := 0
			rowH := utils.Percentalize(windowH, row.Height)
			for i3, col := range row.Cols {
				stackHTotal := 0
				realWidth := windowW / len(row.Cols)
				if col.Width > 0 {
					realWidth = windowW / 12 * col.Width
				}
				for i4, stack := range col.Stacks {
					realHeight := 0
					// jusify height to fill at the very last of stacks
					if i4 == len(col.Stacks)-1 {
						// at the last row
						if i2 == len(tab.Rows)-1 {
							realHeight = windowH - rowHTotal - 3
						} else {
							realHeight = rowH - stackHTotal
						}
					} else {
						realHeight = utils.Percentalize(rowH, stack.Height)
					}
					collection.Register(&vector.RectangleOptions{
						Index:      idx,
						RowIndex:   i2,
						ColIndex:   i3,
						LastCol:    i3 == len(row.Cols)-1,
						FirstStack: i4 == 0,
						LastStack:  i4 == len(col.Stacks)-1,
						WindowW:    windowW,
						WindowH:    windowH,
						TabIndex:   i1,
						CenterX:    colWTotal + (realWidth / 2),
						CenterY:    rowHTotal + stackHTotal + (realHeight / 2),
						Width:      realWidth,
						Height:     realHeight,
					})
					idx++
					stackHTotal += realHeight
				}
				colWTotal += realWidth
			}
			rowHTotal += rowH
		}
		c.totalTab++
	}

	collection.CalcDistances()

	idx = 0
	for i1, tab := range c.config.Layout {
		for _, row := range tab.Rows {
			for _, col := range row.Cols {
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
					col.Widgets = append(col.Widgets, ew)
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

func (c *ViewManager) getWidgetByID(id string) (result *widgetNode, err error) {
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

func (c *ViewManager) renderTabPane(tab *tabNode) (err error) {
	// clear buffer
	ui.Clear()
	ui.Body.Rows = ui.Body.Rows[:0]

	if tab.initialized {
		ui.Body.AddRows(tab.renderedCells...)
		ui.Body.Align()
		ui.Render(ui.Body)
		return
	}
	tab.initialized = true

	screen := []*ui.Row{}
	tabs := make([]*ui.Row, len(c.config.Layout))
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
		tabs[i] = ui.NewCol(2, 0, tabP)
	}
	screen = append(screen, ui.NewRow(tabs...))

	// body
	newRows := make([]*ui.Row, len(tab.Rows))
	for i, row := range tab.Rows {
		newCols := make([]*ui.Row, len(row.Cols))
		for j, col := range row.Cols {
			var cols []ui.GridBufferer
			for _, w := range col.Widgets {
				w.Init()
				gw := w.GetGridBufferers()
				cols = append(cols, gw...)
			}
			colWidth := 12 / len(row.Cols)
			if col.Width > 0 {
				colWidth = col.Width
			}
			newCols[j] = ui.NewCol(colWidth, 0, cols...)
		}
		newRows[i] = ui.NewRow(newCols...)
	}
	screen = append(screen, newRows...)
	tab.renderedCells = screen

	ui.Body.AddRows(screen...)
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
			break
		}
	}
	from := c.getWidgetByIndex(c.activeWidgetIndex)
	if from != nil {
		from.Deactivate()
	}

	c.activateFirstWidgetOnTab(tabIdx)
	c.activeTabIndex = tabIdx
	c.lastWidget = nil
	c.lastDirection = ""
}

// NextWidget is
func (c *ViewManager) NextWidget(direction string) {
	from := c.getWidgetByIndex(c.activeWidgetIndex)
	toIdx := from.GetNeighborIndex(direction)
	to := c.getWidgetByIndex(toIdx)

	// if moving back to last widget, move remembered widget
	tb := c.lastDirection == "top" && direction == "bottom"
	bt := c.lastDirection == "bottom" && direction == "top"
	rl := c.lastDirection == "right" && direction == "left"
	lr := c.lastDirection == "left" && direction == "right"
	if tb || bt || rl || lr {
		c.moveActually(from, c.lastWidget, direction)
	} else if from != nil && to != nil {
		c.moveActually(from, to, direction)
	}
}

func (c *ViewManager) moveActually(from *widget.WrapperWidget, to *widget.WrapperWidget, direction string) {
	if !to.IsDisabled() && to.IsReady() {
		from.Deactivate()
		to.Activate()
		c.activeWidgetIndex = to.Index
		c.lastWidget = from
		c.lastDirection = direction
	} else if to.IsDisabled() {
		return
	} else if !to.IsReady() {
		// skip !IsReady widget
		skipIdx := to.GetNeighborIndex(direction)
		skip := c.getWidgetByIndex(skipIdx)
		if skipIdx != -1 && skip != nil {
			from.Deactivate()
			c.activeWidgetIndex = to.Index
			c.lastWidget = from
			c.lastDirection = direction
			c.NextWidget(direction)
		}
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
