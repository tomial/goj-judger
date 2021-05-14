package judger

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

// compare the output of user's code with correct answer
func compare(volumeDir string, caseNum int) ResultCode {
	fmt.Println("正在对比运行结果")
	path := volumeDir + "rte"
	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		// 没有runtime error
	} else {
		// runtime error
		if fi.Size() > 0 {
			data, _ := ioutil.ReadFile(path)
			fmt.Println(string(data))
			fmt.Println("RUNTIME_ERROR")
			return RUNTIME_ERROR
		}
	}

	// RUNTIME ERROR

	for i := 1; i <= caseNum; i++ {
		caseNo := strconv.Itoa(i)

		cmd := exec.Command("diff", "-uZ", volumeDir+"output-"+caseNo, volumeDir+"res-"+caseNo)
		err := cmd.Run()
		if err != nil || cmd.ProcessState.ExitCode() > 0 {
			fmt.Println("答案错误")
			exec.Command("cp", volumeDir+"input-"+caseNo, volumeDir+"wa_input").Run()
			exec.Command("cp", volumeDir+"output-"+caseNo, volumeDir+"wa_output").Run()
			exec.Command("cp", volumeDir+"res-"+caseNo, volumeDir+"wa_expect").Run()
			return WRONG_ANSWER
		}
	}

	fmt.Println("提交通过")

	return PASS
}
