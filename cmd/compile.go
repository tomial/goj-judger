package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomial/goj-judger/internal/compile"
	"github.com/tomial/goj-judger/internal/params"
)

var bp = params.NewBuild()

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:       "compile <language>",
	Short:     "compile source code in docker",
	Long:      "compile source code in docker",
	ValidArgs: []string{"go", "golang", "c", "c++", "cpp"},
	Args:      cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) <= 0 {
			fmt.Println(cmd.UsageString())
		} else {
			lang := strings.ToLower(args[0])

			if bp.SourcePath == "" {
				bp.SourcePath = "./source" + params.SourceSuffix[lang]
			}

			switch lang {
			case "go":
				fallthrough
			case "golang":
				fmt.Println("Compiling go code.")
				fmt.Println("[Params]")
				fmt.Printf("- Source Path: %10v\n", bp.SourcePath)
				fmt.Printf("- Size Limit: %10vMB\n", bp.SizeLimit)
				fmt.Printf("- Time Limit: %10v\n", bp.TimeLimit)
				compile.Golang(bp)
			case "c":
				fmt.Println("Compiling c code.")
			case "c++":
				fallthrough
			case "cpp":
				fmt.Println("Compiling c++ code.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)

	compileCmd.Flags().StringVarP(&bp.SourcePath, "source code path", "p", bp.SourcePath,
		"specify <path>/<name>.<suffix> for source code.")

	compileCmd.Flags().IntVarP(&bp.SizeLimit, "binary file size limit", "s", bp.SizeLimit,
		"specify size for compiled binary file in MB.")

	var compileTime int64
	compileCmd.Flags().Int64VarP(&compileTime, "compile time limit", "t", 5000,
		"specify compile time in ms.")

	bp.TimeLimit = time.Duration(compileTime) * time.Millisecond

	// "specify source code language. currently supported languages: go, c, c++")
}
