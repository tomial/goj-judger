package judger

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tomial/goj-judger/internal/params"
)

func JudgeGo(jr *params.JudgeRequest) {
	caseNum = jr.CaseNum
	tlimit = jr.TimeLimit
	rlimit = jr.RamLimit
	volumeDir = jr.VolumeDir
	containerName = "go-judger-container"

	imageName = "miata/goj-judger-go-img"
	buildCmd = "go build -o main source.go 2> build_result; exit $?"
	runCmd = "./xtime.sh " + strconv.Itoa(caseNum) + " " + strconv.Itoa(tlimit)

	// 准备Go判题镜像
	err := prepareImg()

	if err != nil {
		fmt.Println(err, "准备阶段发生错误")
		os.Exit(UNKNOWN_ERROR)
	}

	// 编译Go代码
	err = compile()
	if err != nil {
		fmt.Println(err, "编译错误")
		os.Exit(COMPILE_ERROR)
	}

	fmt.Println("正在运行用户Go语言代码")
	err = run()
	if err != nil {
		fmt.Println(err, "运行用户Go语言代码失败")
		os.Exit(UNKNOWN_ERROR)
	}

	result := compare(volumeDir, caseNum)

	os.Exit(result)
}
