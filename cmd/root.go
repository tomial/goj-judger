package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gj",
	Short: "goj-judger is a judger for goj",
	Long:  `goj-judger is a cli application, it compiles and run code in docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello Cobra CLI.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
