package dashboard

import (
	"io/ioutil"

	ui "github.com/gizak/termui"
	"github.com/qmu/mcc/utils"
	yaml "gopkg.in/yaml.v1"
)

// ConfigManager load and unmarshal config file
type ConfigManager struct {
	LoadedData *Config
}

// Config is the root schema of config file
type Config struct {
	SchemaVersion string `yaml:"schema_version"`
	Timezone      string
	Envs          []map[string]string
	Rows          []Row
}

// Row is the schema implements Config.Widgets.Section
type Row struct {
	Section string
	Height  string // percent
	Cols    []Col
}

// Col is the schema implements Config.Widgets.Section
type Col struct {
	Section string
	Widgets []Widget
}

// Widget is the schema implements Config.Widgets
type Widget struct {
	Title      string
	Col        int
	Height     string // percent
	RealHeight int
	Type       string
	IssueRegex string `yaml:"issue_regex"`
	Content    interface{}
}

// Menu is the schema implements Config.Widgets.Menu
type Menu struct {
	Category    string
	Name        string
	Description string
	Command     string
}

// NewConfigManager constructs a ConfigManager
func NewConfigManager(path string) (c *ConfigManager, err error) {
	c = new(ConfigManager)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	if err = yaml.Unmarshal(file, &c.LoadedData); err != nil {
		return
	}
	windowH := ui.TermHeight() - 1
	for i1, row := range c.LoadedData.Rows {
		rowH := utils.Percentalize(windowH, row.Height)
		for i2, col := range row.Cols {
			cntH := 0
			for i3, w := range col.Widgets {
				if len(col.Widgets)-1 == i3 {
					c.LoadedData.Rows[i1].Cols[i2].Widgets[i3].RealHeight = rowH - cntH
				} else {
					widgetH := utils.Percentalize(rowH, w.Height)
					c.LoadedData.Rows[i1].Cols[i2].Widgets[i3].RealHeight = widgetH
					cntH += widgetH
				}
			}
		}
	}
	return
}
