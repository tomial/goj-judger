package judger

import (
	"os"
	"os/exec"
)

type judger interface {
	Start() error
	prepare() error
	compile() error
	run() error
}

func NewGolang() judger {
	return &golang{}
}

func NewC() judger {
	return &c{}
}

func NewCPP() judger {
	return &cpp{}
}

// compare the output of user's code with correct answer
func compare(out []byte) bool {
	err := os.WriteFile("/home/lsxph/volumetest/output", out, 0666)
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("diff", "-uZ", "/home/lsxph/volumetest/output", "/home/lsxph/volumetest/res").CombinedOutput()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return false
		}
	}
	return true
}
