package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "v1.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "This command show a current version",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(appVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}