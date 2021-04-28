package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tomial/goj-judger/internal/judger"
	"github.com/tomial/goj-judger/internal/params"
)

var bp = params.NewBuild()

// judgeCmd represents the compile command
var judgeCmd = &cobra.Command{
	Use:       "judge <language>",
	Short:     "judge source code in docker",
	Long:      "judge source code in docker",
	ValidArgs: []string{"go", "golang", "c", "c++", "cpp"},
	Args:      cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) <= 0 {
			fmt.Println(cmd.UsageString())
		} else {
			lang := strings.ToLower(args[0])

			if bp.WorkDir == "" {
				bp.WorkDir = "/build"
			}

			// fmt.Println("Compiling go code.")
			// fmt.Println("[Params]")
			// fmt.Printf("- Source Path: %10v\n", bp.WorkDir)
			// fmt.Printf("- Size Limit: %10vMB\n", bp.SizeLimit)
			// fmt.Printf("- Time Limit: %10v\n", bp.TimeLimit)

			var err error
			var j judger.Judger

			switch lang {
			case "go":
				fallthrough
			case "golang":
				fmt.Println("Judging go code.")
				j = judger.NewGolang()
				err = j.Start()
			case "c":
				fmt.Println("Judging c code.")
				j := judger.NewC()
				err = j.Start()
			case "c++":
				fallthrough
			case "cpp":
				fmt.Println("Judging c++ code.")
				j := judger.NewCPP()
				err = j.Start()
			}

			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(judgeCmd)

	judgeCmd.Flags().StringVarP(&bp.WorkDir, "work directory for container", "p", bp.WorkDir,
		"specify work directory for container")

	judgeCmd.Flags().IntVarP(&bp.SizeLimit, "binary file size limit", "s", bp.SizeLimit,
		"specify size for compiled binary file in MB.")

	var compileTime int64
	judgeCmd.Flags().Int64VarP(&compileTime, "compile time limit", "t", 5000,
		"specify compile time in ms.")

	bp.TimeLimit = time.Duration(compileTime) * time.Millisecond

	// "specify source code language. currently supported languages: go, c, c++")
}
