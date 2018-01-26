package cmd

import (
	"fmt"
	"os"

	"github.com/ans-ashkan/thc/config"

	"github.com/spf13/cobra"
)

// RootCmd of twh
var RootCmd = &cobra.Command{
	Use:   "twh",
	Short: "twh is a twitter cli helper",
	Run: func(cmd *cobra.Command, args []string) {
		// homeDir, err := homedir.Dir()
		cfg := config.GetConfig()
		if cfg == nil {
			panic("No config found")
		}

		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
