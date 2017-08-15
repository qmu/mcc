package cmd

import (
	"github.com/qmu/mcc/dashboard"
	"github.com/spf13/cobra"
)

var config string

// RootCmd is the cobra's root command
var RootCmd = &cobra.Command{
	Use: "mcc",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dashboard.NewDashboard(config); err != nil {
			return
		}
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "mcc.yml", "Path of the mcc config file")
	cobra.OnInitialize()
}
