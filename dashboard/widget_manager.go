package dashboard

import (
	"errors"
	"math"
	"strconv"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/utils"
)

// WidgetManager load and unmarshal config file
type WidgetManager struct {
	config        *Config
	configManager *ConfigManager
	exWidgets     []*ExtendedWidget
	options       *WidgetManagerOptions
}

// WidgetManagerOptions is some options for ConfigManager constructor
type WidgetManagerOptions struct {
	execPath            string
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
	err = w.loadWidgetsToLayout()
	if err != nil {
		return
	}
	w.constructRealWidgets()
	w.calcDistances()

	return
}

func (c *WidgetManager) loadWidgetsToLayout() (err error) {
	for i1, tab := range c.config.Layout {
		for i2, row := range tab.Rows {
			for i3, col := range row.Cols {
				for _, s := range col.Stacks {
					wi, err := c.getWidgetByID(s)
					if err != nil {
						return err
					}
					c.config.Layout[i1].Rows[i2].Cols[i3].Widgets = append(c.config.Layout[i1].Rows[i2].Cols[i3].Widgets, wi)
				}
			}
		}
	}
	return
}

func (c *WidgetManager) constructRealWidgets() (err error) {
	windowH := ui.TermHeight() - 1
	windowW := ui.TermWidth() - 1
	idx := 0
	for i1, tab := range c.config.Layout {
		rowHTotal := 0
		for i2, row := range tab.Rows {
			rowH := utils.Percentalize(windowH, row.Height)
			for i3, col := range row.Cols {
				cntH := 0
				offset := 1
				stackHTotal := 0
				for i4, w := range col.Widgets {
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
						tab:        i1,
						row:        i2,
						col:        i3,
						stack:      i4,
						vOffset:    offset,
						height:     realHeight,
						width:      realWidth,
						point: Point{
							x: windowW/len(row.Cols)*i3 + (realWidth / 2),
							y: rowHTotal + stackHTotal + (realHeight / 2),
						},
						widthFrom: 100 / len(row.Cols) * i3,
						widthTo:   100 / len(row.Cols) * (i3 + 1),
						title:     w.Title,
					}
					rw.title += "   x:" + strconv.Itoa(rw.point.x) + " y:" + strconv.Itoa(rw.point.y)

					err := rw.Vary(&WidgetOptions{
						envs:     c.config.Envs,
						execPath: c.options.execPath,
						timezone: c.config.Timezone,
					})
					if err != nil {
						return err
					}

					c.exWidgets = append(c.exWidgets, rw)
					c.config.Layout[i1].Rows[i2].Cols[i3].Widgets[i4].extendedWidget = rw
					offset += realHeight
					idx++
					stackHTotal += realHeight
				}
			}
			rowHTotal += rowH
		}
	}
	return
}

func (c *WidgetManager) calcDistances() {
	type distance struct {
		from    *ExtendedWidget
		to      *ExtendedWidget
		toX     int
		toY     int
		toTab   int
		toIndex int
		value   float64
	}
	distances := []*distance{}
	for i1, w1 := range c.exWidgets {
		for i2, w2 := range c.exWidgets {
			if i1 != i2 && i1 < i2 && w1.tab == w2.tab {
				distances = append(distances, &distance{
					from:    w1,
					to:      w2,
					toX:     w2.point.x,
					toY:     w2.point.y,
					toTab:   w2.tab,
					toIndex: w2.index,
					value:   c.vectorDistance(w1.point, w2.point),
				})
			}
		}
	}
	for _, w := range c.exWidgets {
		// bottom
		var nearest *distance
		for _, d := range distances {
			if w.tab == d.toTab && w.index != d.toIndex && w.point.y < d.toY && math.Abs(float64(d.to.point.x-w.point.x)) < float64(w.width/2) {
				if nearest == nil || nearest.value > d.value {
					nearest = d
				}
			}
		}
		if nearest != nil {
			w.bottomWidget = nearest.to
		}

		// top
		nearest = nil
		for _, d := range distances {
			if w.tab == d.toTab && w.index != d.toIndex && w.point.y > d.toY && math.Abs(float64(d.to.point.x-w.point.x)) < float64(w.width/2) {
				if nearest == nil || nearest.value > d.value {
					nearest = d
				}
			}
		}
		if nearest != nil {
			w.topWidget = nearest.to
		}

		// right
		nearest = nil
		for _, d := range distances {
			if w.tab == d.toTab && w.index != d.toIndex && w.point.x < d.toX {
				if nearest == nil || nearest.value > d.value {
					nearest = d
				}
			}
		}
		if nearest != nil {
			w.rightWidget = nearest.to
		}

		// left
		nearest = nil
		for _, d := range distances {
			if w.tab == d.toTab && w.index != d.toIndex && w.point.x > d.toX {
				if nearest == nil || nearest.value > d.value {
					nearest = d
				}
			}
		}
		if nearest != nil {
			w.leftWidget = nearest.to
		}
	}

}

func (c *WidgetManager) vectorDistance(fromPoint Point, toPoint Point) (distance float64) {
	x1 := fromPoint.x
	y1 := fromPoint.y
	x2 := toPoint.x
	y2 := toPoint.y
	distance = math.Pow(float64((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1)), 0.5)
	return
}

func (c *WidgetManager) getWidgetByID(id string) (result Widget, err error) {
	for _, d := range c.config.Widgets {
		if d.ID == id {
			result = d
			return
		}
	}
	return result, errors.New("no widget named " + id)
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
