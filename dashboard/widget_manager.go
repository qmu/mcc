package dashboard

import "github.com/qmu/mcc/utils"
import ui "github.com/gizak/termui"

// WidgetManager load and unmarshal config file
type WidgetManager struct {
	config        *Config
	configManager *ConfigManager
	exWidgets     []*ExtendedWidget
	options       *WidgetManagerOptions
}

// WidgetManagerOptions is some options for ConfigManager constructor
type WidgetManagerOptions struct {
	configPath          string
	appVersion          string
	configSchemaVersion string
}

// NewWidgetManager constructs a ConfigManager
func NewWidgetManager(opt *WidgetManagerOptions) (w *WidgetManager, err error) {
	w = new(WidgetManager)
	w.options = opt

	// load config file
	w.configManager, err = NewConfigManager(&ConfigManagerOptions{
		configPath:          w.options.configPath,
		appVersion:          w.options.appVersion,
		configSchemaVersion: w.options.configSchemaVersion,
	})
	if err != nil {
		return
	}
	w.config = w.configManager.GetConfig()

	w.constructRealWidgets()

	return
}

func (c *WidgetManager) constructRealWidgets() {
	windowH := ui.TermHeight() - 1
	windowW := ui.TermWidth() - 1
	idx := 0
	for i1, row := range c.config.Rows {
		rowH := utils.Percentalize(windowH, row.Height)
		for i2, col := range row.Cols {
			cntH := 0
			offset := 1
			for i3, w := range col.Widgets {
				realWidth := windowW / len(row.Cols)
				realHeight := 0
				if len(col.Widgets)-1 == i3 {
					realHeight = rowH - cntH
				} else {
					realHeight = utils.Percentalize(rowH, w.Height)
					cntH += realHeight
				}
				rw := &ExtendedWidget{
					index:      idx,
					widget:     w,
					widgetType: w.Type,
					row:        i1,
					col:        i2,
					stack:      i3,
					vOffset:    offset,
					height:     realHeight,
					width:      realWidth,
					widthFrom:  100 / len(row.Cols) * i2,
					widthTo:    100 / len(row.Cols) * (i2 + 1),
					title:      w.Title,
				}
				c.exWidgets = append(c.exWidgets, rw)
				c.config.Rows[i1].Cols[i2].Widgets[i3].extendedWidget = rw
				offset += realHeight
				idx++
			}
		}
	}
}

// HasWidget returns whether config contains 'widgetType' stack or not
func (c *WidgetManager) HasWidget(widgetType string) bool {
	result := false
	for _, ew := range c.exWidgets {
		if ew.widgetType == widgetType {
			result = true
		}
	}
	return result
}

// GetDownerWidget is
func (c *WidgetManager) GetDownerWidget(from int) (result *ExtendedWidget) {
	cw := c.exWidgets[from]
	if cw.index == len(c.exWidgets)-1 {
		return c.exWidgets[0]
	}
	for i := 0; i < len(c.exWidgets); i++ {
		w := c.exWidgets[i]
		if !w.IsReady() || w.IsDisabled() {
			continue
		}
		cond1 := cw.row == w.row && cw.col == w.col && w.stack == cw.stack+1
		cond2 := cw.row < w.row && w.stack == 0 && w.widthFrom <= cw.widthFrom && cw.widthFrom <= w.widthTo
		if cond1 || cond2 {
			result = w
			return
		}
	}
	return
}

// GetUpperWidget is
func (c *WidgetManager) GetUpperWidget(from int) (result *ExtendedWidget) {
	cw := c.exWidgets[from]
	if cw.row == 0 && cw.stack == 0 {
		return c.exWidgets[len(c.exWidgets)-1]
	}
	for i := len(c.exWidgets) - 1; i >= 0; i-- {
		w := c.exWidgets[i]
		if !w.IsReady() || w.IsDisabled() {
			continue
		}
		if i < from {
			cond1 := w.row == cw.row && w.col == cw.col && w.stack == cw.stack-1
			cond2 := w.row < cw.row && w.widthFrom <= cw.widthFrom && cw.widthFrom <= w.widthTo
			if cond1 || cond2 {
				result = w
				return
			}
		}
	}
	return
}

// GetRightWidget is
func (c *WidgetManager) GetRightWidget(from int) (result *ExtendedWidget) {
	col := 0
	cw := c.exWidgets[from]
	cursor := cw.GetHighlightenPos()
	myPosition := cw.vOffset + cursor
	for i, w := range c.exWidgets {
		if !w.IsReady() || w.IsDisabled() {
			continue
		}
		if w.row == cw.row && i > from && col < w.col {
			if myPosition >= w.vOffset && myPosition <= w.vOffset+w.height {
				result = w
				return
			}
		} else {
			col = w.col
		}
	}
	return
}

// GetLeftWidget is
func (c *WidgetManager) GetLeftWidget(from int) (result *ExtendedWidget) {
	col := c.exWidgets[len(c.exWidgets)-1].col
	cw := c.exWidgets[from]
	cursor := cw.GetHighlightenPos()
	myPosition := cw.vOffset + cursor
	for i := len(c.exWidgets) - 1; i >= 0; i-- {
		w := c.exWidgets[i]
		if !w.IsReady() || w.IsDisabled() {
			continue
		}
		if i < from && col > w.col && cw.row == w.row {
			if myPosition >= w.vOffset && myPosition <= w.vOffset+w.height {
				result = w
				return
			}
		} else {
			col = w.col
		}
	}
	return
}

// GetWidgetByIndex is
func (c *WidgetManager) GetWidgetByIndex(index int) (w *ExtendedWidget) {
	for k, w := range c.exWidgets {
		if index == k {
			return w
		}
	}
	return nil
}

// GetActiveWidgetsOf is
func (c *WidgetManager) GetActiveWidgetsOf(name string) (result []*ExtendedWidget) {
	for _, w := range c.exWidgets {
		if w.widgetType == name && !w.IsDisabled() {
			result = append(result, w)
		}
	}
	return
}

// GetAllWidgets is
func (c *WidgetManager) GetAllWidgets() (result []*ExtendedWidget) {
	return c.exWidgets
}

// GetAllWidgetsCount is
func (c *WidgetManager) GetAllWidgetsCount() (result int) {
	return len(c.exWidgets)
}
