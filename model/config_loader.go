package model

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	version "github.com/hashicorp/go-version"
	"github.com/qmu/mcc/widget"
	yaml "gopkg.in/yaml.v1"
)

// ConfRoot is the root schema of config file
type ConfRoot struct {
	SchemaVersion string `yaml:"schema_version"`
	Timezone      string
	GitHubHost    string `yaml:"github_url"`
	Envs          []map[string]string
	Widgets       []widgetNode
	Layout        []tabNode
}

// tabNode is the schema implements ConfRoot.OriginalWidgets.Section
type tabNode struct {
	Index int
	Name  string
	Rows  []rowNode
}

// rowNode is the schema implements ConfRoot.OriginalWidgets.Section
type rowNode struct {
	Height string // percent
	Cols   []colNode
}

// colNode is the schema implements ConfRoot.OriginalWidgets.Section
type colNode struct {
	Widgets []*widget.WrapperWidget
	Stacks  []stackNode
	Width   int // grid system, accepts 1~12
}

// colNode is the schema implements ConfRoot.OriginalWidgets.Section
type stackNode struct {
	ID     string
	Height string // percent
}

// widgetNode is the schema implements ConfRoot.OriginalWidgets
type widgetNode struct {
	ID         string
	Title      string
	Type       string
	IssueRegex string `yaml:"issue_regex"`
	Content    interface{}
	Path       string
}

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
	validator, err := NewConfigValidator()
	if err != nil {
		return
	}

	res, err := validator.validate(c.config)
	if err != nil {
		return
	}
	if res != nil {
		fmt.Println("==================================================================")
		fmt.Println("CONFIGURATION ERROR")
		fmt.Println("==================================================================")
		for i, r := range res {
			fmt.Println("No." + strconv.Itoa(i+1) + " : " + r.message + " (position : " + r.position + ")")
		}
		fmt.Println("==================================================================")
		os.Exit(0)
	}

	err = c.fixIncompleteParamas()
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

func (c *ConfigLoader) fixIncompleteParamas() (err error) {
	// for _, t := range c.config.Layout {
	// 	for _, r := range t.Rows {
	// 		pp.Println(r.Height)
	// 	}
	// }
	// os.Exit(0)
	return
}

// GetConfig is
func (c *ConfigLoader) GetConfig() *ConfRoot {
	return c.config
}
