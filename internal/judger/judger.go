package judger

import (
	"context"
	"os"
	"os/exec"
)

type Judger interface {
	Start() error
	prepare() error
	compile() error
	run() error
}

var ctx context.Context

func init() {
	ctx = context.Background()
}

func NewGolang() Judger {
	return &golang{}
}

func NewC() Judger {
	return &c{}
}

func NewCPP() Judger {
	return &cpp{}
}

// compare the output of user's code with correct answer
func compare(out []byte) bool {
	err := os.WriteFile("/home/lsxph/volume/output", out, 0666)
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("diff", "-uZ", "/home/lsxph/volume/output", "/home/lsxph/volume/res").CombinedOutput()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			return false
		}
	}
	return true
}
