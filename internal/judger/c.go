package judger

import (
	"fmt"
	"os"
	"strconv"

	"github.com/tomial/goj-judger/internal/params"
)

func JudgeC(jr *params.JudgeRequest) {
	caseNum = jr.CaseNum
	tlimit = jr.TimeLimit
	rlimit = jr.RamLimit
	volumeDir = jr.VolumeDir
	containerName = "c-cpp-judger-container"

	imageName = "miata/goj-judger-c-cpp-img"
	buildCmd = "gcc -o main source.c -Wfatal-errors 2> build_result; exit $?"
	runCmd = "./xtime.sh " + strconv.Itoa(caseNum) + " " + strconv.Itoa(tlimit)

	// 准备C判题镜像
	err := prepareImg()

	if err != nil {
		fmt.Println(err, "准备阶段发生错误")
		os.Exit(UNKNOWN_ERROR)
	}

	// 编译C代码
	err = compile()
	if err != nil {
		fmt.Println(err, "编译错误")
		os.Exit(COMPILE_ERROR)
	}

	fmt.Println("正在运行用户C代码")
	err = run()
	if err != nil {
		fmt.Println(err, "运行用户C代码失败")
		os.Exit(UNKNOWN_ERROR)
	}

	result := compare(volumeDir, caseNum)

	os.Exit(result)
}
