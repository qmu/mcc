package dashboard

import (
	"fmt"
	"io/ioutil"
	"os"

	version "github.com/hashicorp/go-version"
	yaml "gopkg.in/yaml.v1"
)

// ConfigManager load and unmarshal config file
type ConfigManager struct {
	config  *Config
	options *ConfigManagerOptions
}

// ConfigManagerOptions is some options for ConfigManager constructor
type ConfigManagerOptions struct {
	configPath          string
	appVersion          string
	configSchemaVersion string
}

// Config is the root schema of config file
type Config struct {
	SchemaVersion string `yaml:"schema_version"`
	Timezone      string
	GitHubHost    string `yaml:"github_url"`
	Envs          []map[string]string
	Widgets       []Widget
	Layout        []Tab
}

// Tab is the schema implements Config.Widgets.Section
type Tab struct {
	Section string
	Name    string
	Rows    []Row
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
	Stacks  []string
}

// Widget is the schema implements Config.Widgets
type Widget struct {
	ID             string
	Title          string
	Col            int
	Height         string // percent
	Type           string
	IssueRegex     string `yaml:"issue_regex"`
	Content        interface{}
	Path           string
	extendedWidget *ExtendedWidget
}

// Menu is the schema implements Config.Widgets.Menu
type Menu struct {
	Category    string
	Name        string
	Description string
	Command     string
}

// Container is the schema implements Config.Widgets.Conttainer
type Container struct {
	Metrics   string
	Name      string
	Container string
}

// NewConfigManager constructs a ConfigManager
func NewConfigManager(opt *ConfigManagerOptions) (c *ConfigManager, err error) {
	c = new(ConfigManager)
	c.options = opt
	file, err := ioutil.ReadFile(opt.configPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(file, &c.config)
	if err != nil {
		return
	}

	// check ConfigSchemaVersion
	if err = c.checkConfigScheme(); err != nil {
		return
	}

	return
}

func (c *ConfigManager) checkConfigScheme() (err error) {
	vApp, err := version.NewVersion(c.options.configSchemaVersion)
	vConfig, err := version.NewVersion(c.config.SchemaVersion)
	if err != nil {
		return
	}
	if vConfig.LessThan(vApp) {
		fmt.Printf("mcc %s supports schema_version %s but ths schema_version in %s seems to be %s\n", c.options.appVersion, vApp, c.options.configPath, vConfig)
		fmt.Printf("please upgrade mcc or %s first\n", c.options.configPath)
		os.Exit(1)
	}
	return
}

// GetConfig is
func (c *ConfigManager) GetConfig() *Config {
	return c.config
}
