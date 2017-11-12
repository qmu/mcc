package controller

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Version defined in Makefile
var Version string

// ConfigSchemaVersion defined in Makefile
var ConfigSchemaVersion string

// RootCmd is the cobra's root command
var RootCmd *cobra.Command

// Run runs mcc
func Run() (err error) {
	var config string
	RootCmd = &cobra.Command{
		Use: "mcc",
		Run: func(cmd *cobra.Command, args []string) {
			version, err := cmd.Flags().GetBool("version")
			if err == nil && version {
				fmt.Println("mcc version " + Version)
				os.Exit(0)
			}
			if config == "" {
				if _, err := os.Stat("./mcc.yml"); err == nil {
					config = "./mcc.yml"
				} else {
					fmt.Println("Error: check \"mcc.yml\" exists in the current directory, or use -c to set its path")
					os.Exit(1)
				}
			}
			if err := NewController(Version, ConfigSchemaVersion, config); err != nil {
				log.Panicln(err)
				os.Exit(1)
			}
		},
	}

	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "path to a yaml config")
	RootCmd.PersistentFlags().BoolP("version", "v", false, "print the version")
	cobra.OnInitialize()

	if err = RootCmd.Execute(); err != nil {
		return err
	}
	return err
}
