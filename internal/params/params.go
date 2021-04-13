package params

import "time"

type Build struct {
	WorkDir   string // path of source code
	TimeLimit time.Duration
	SizeLimit int
}

var SourceSuffix = map[string]string{
	"go":  ".go",
	"c":   ".c",
	"cpp": ".cpp",
	"c++": ".cpp",
}

func NewBuild() *Build {
	return &Build{
		WorkDir:   "",
		TimeLimit: 5000 * time.Microsecond,
		SizeLimit: 10,
	}
}
