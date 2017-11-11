package model

import (
	"fmt"
	"io/ioutil"
	"os"

	version "github.com/hashicorp/go-version"
	"github.com/qmu/mcc/widget"
	yaml "gopkg.in/yaml.v1"
)

// ConfigLoader load and unmarshal config file
type ConfigLoader struct {
	config  *ConfRoot
	options *ConfigLoaderOption
}

// ConfigLoaderOption is some options for ConfigLoader constructor
type ConfigLoaderOption struct {
	ExecPath            string
	ConfigPath          string
	AppVersion          string
	ConfigSchemaVersion string
}

// ConfRoot is the root schema of config file
type ConfRoot struct {
	SchemaVersion string `yaml:"schema_version"`
	Timezone      string
	GitHubHost    string `yaml:"github_url"`
	Envs          []map[string]string
	Widgets       []confWidget
	Layout        []confTab
}

// confTab is the schema implements ConfRoot.OriginalWidgets.Section
type confTab struct {
	Section string
	Index   int
	Name    string
	Rows    []confRow
}

// confRow is the schema implements ConfRoot.OriginalWidgets.Section
type confRow struct {
	Section string
	Height  string // percent
	Cols    []confCol
}

// confCol is the schema implements ConfRoot.OriginalWidgets.Section
type confCol struct {
	Section string
	Widgets []*widget.WrapperWidget
	Stacks  []confStack
	Width   int // grid system, accepts 1~12
}

// confCol is the schema implements ConfRoot.OriginalWidgets.Section
type confStack struct {
	ID     string
	Height string // percent
}

// confWidget is the schema implements ConfRoot.OriginalWidgets
type confWidget struct {
	ID         string
	Title      string
	Col        int
	Type       string
	IssueRegex string `yaml:"issue_regex"`
	Content    interface{}
	Path       string
}

// NewLoader constructs a ConfigLoader
func NewLoader(opt *ConfigLoaderOption) (c *ConfigLoader, err error) {
	c = new(ConfigLoader)
	c.options = opt
	file, err := ioutil.ReadFile(opt.ConfigPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(file, &c.config)
	if err != nil {
		return
	}
	err = c.checkConfigScheme()
	if err != nil {
		return
	}
	return
}

func (c *ConfigLoader) checkConfigScheme() (err error) {
	vApp, err := version.NewVersion(c.options.ConfigSchemaVersion)
	vConfig, err := version.NewVersion(c.config.SchemaVersion)
	if err != nil {
		return
	}
	if vConfig.LessThan(vApp) {
		fmt.Printf("mcc %s supports schema_version %s but ths schema_version in %s seems to be %s\n", c.options.AppVersion, vApp, c.options.ConfigPath, vConfig)
		fmt.Printf("please upgrade mcc or %s first\n", c.options.ConfigPath)
		os.Exit(1)
	}
	return
}

// GetConfig is
func (c *ConfigLoader) GetConfig() *ConfRoot {
	return c.config
}
