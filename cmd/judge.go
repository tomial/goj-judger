package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tomial/goj-judger/internal/judger"
	"github.com/tomial/goj-judger/internal/params"
)

var jr = params.NewJudgeRequest()

// judgeCmd represents the compile command
var judgeCmd = &cobra.Command{
	Use:       "judge <language>",
	Short:     "judge source code in docker",
	Long:      "judge source code in docker",
	ValidArgs: []string{"go", "golang", "c", "c++", "cpp"},
	Args:      cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) <= 0 || jr.CaseNum <= 0 || jr.RamLimit <= 0 || jr.TimeLimit <= 0 || jr.VolumeDir == "" {
			fmt.Println(cmd.UsageString())
		} else {
			lang := strings.ToLower(args[0])

			var err error

			fmt.Println("CaseNum:", jr.CaseNum)
			fmt.Println("Tlimit:", jr.TimeLimit)
			fmt.Println("Rlimit:", jr.RamLimit)
			fmt.Println("VolumeDir:", jr.VolumeDir)

			switch lang {
			case "go":
				fallthrough
			case "golang":
				fmt.Println("正在判断Go代码")
				judger.JudgeGo(jr)
			case "c":
				fmt.Println("正在判断C代码")
				judger.JudgeC(jr)
			case "c++":
				fallthrough
			case "cpp":
				fmt.Println("正在判断C++代码")
				judger.JudgeCPP(jr)
			}

			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(judgeCmd)

	judgeCmd.Flags().IntVarP(&jr.CaseNum, "测试用例数量", "n", 0,
		"指定测试用例数量")

	judgeCmd.Flags().StringVarP(&jr.VolumeDir, "volume挂载的目录", "d", "",
		"指定volume挂载的目录")

	judgeCmd.Flags().IntVarP(&jr.TimeLimit, "时间限制", "t", jr.TimeLimit,
		"指定运行时间限制")

	judgeCmd.Flags().IntVarP(&jr.RamLimit, "内存限制", "r", jr.RamLimit,
		"指定运行内存限制")
}
